package app

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateCommandReturnsZeroWithWarnings(t *testing.T) {
	input := readPlaceholder(t)
	input = bytes.Replace(input, []byte(`foreground = "#C8D3F5"`), []byte(`foreground = "#222436"`), 1)
	path := writeTheme(t, input)

	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), strings.NewReader(""), &stdout, &stderr,
		[]string{"theme", "validate", path})
	if code != 0 {
		t.Fatalf("Run() code = %d, want 0; stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), ":colors.foreground: WARNING:") {
		t.Fatalf("stdout = %q, want foreground warning", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestValidateCommandReturnsOneForInvalidTheme(t *testing.T) {
	input := bytes.Replace(readPlaceholder(t), []byte("version = 1"), []byte("version = 2"), 1)
	path := writeTheme(t, input)

	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), strings.NewReader(""), &stdout, &stderr,
		[]string{"theme", "validate", path})
	if code != 1 {
		t.Fatalf("Run() code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), ":version: ERROR:") {
		t.Fatalf("stderr = %q, want version error", stderr.String())
	}
}

func TestUsageReturnsTwo(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), strings.NewReader(""), &stdout, &stderr, nil)
	if code != 2 {
		t.Fatalf("Run() code = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "usage:") {
		t.Fatalf("stderr = %q, want usage", stderr.String())
	}
}

func readPlaceholder(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "..", "themes", "yolodev", "placeholder.toml"))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func writeTheme(t *testing.T, data []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "theme.toml")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}
