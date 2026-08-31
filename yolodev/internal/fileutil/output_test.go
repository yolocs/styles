package fileutil

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveOutputPathExpandsHomeDirectory(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	got, err := ResolveOutputPath("~", "signal-grid", ".tmTheme")
	if err != nil {
		t.Fatalf("ResolveOutputPath() error = %v", err)
	}
	want := filepath.Join(home, "signal-grid.tmTheme")
	if got != want {
		t.Fatalf("ResolveOutputPath() = %q, want %q", got, want)
	}
}

func TestResolveOutputPathExpandsHomePrefixForExplicitFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	got, err := ResolveOutputPath("~/exports/custom.tmTheme", "signal-grid", ".tmTheme")
	if err != nil {
		t.Fatalf("ResolveOutputPath() error = %v", err)
	}
	want := filepath.Join(home, "exports", "custom.tmTheme")
	if got != want {
		t.Fatalf("ResolveOutputPath() = %q, want %q", got, want)
	}
}

func TestResolveOutputPathUsesMissingExtensionlessPathAsDirectory(t *testing.T) {
	t.Parallel()

	directory := filepath.Join(t.TempDir(), "nested", "themes")
	got, err := ResolveOutputPath(directory, "signal-grid", ".tmTheme")
	if err != nil {
		t.Fatalf("ResolveOutputPath() error = %v", err)
	}
	want := filepath.Join(directory, "signal-grid.tmTheme")
	if got != want {
		t.Fatalf("ResolveOutputPath() = %q, want %q", got, want)
	}
}

func TestResolveOutputPathRejectsPathSeparatorInStem(t *testing.T) {
	t.Parallel()

	_, err := ResolveOutputPath(t.TempDir(), "../outside", ".tmTheme")
	if err == nil || !strings.Contains(err.Error(), "filename") {
		t.Fatalf("ResolveOutputPath() error = %v, want invalid filename", err)
	}
}
