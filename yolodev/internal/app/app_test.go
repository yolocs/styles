package app

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yolocs/styles/yolodev/internal/theme"
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

func TestExportGhosttyWritesStdout(t *testing.T) {
	path := writeTheme(t, readPlaceholder(t))
	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), strings.NewReader(""), &stdout, &stderr,
		[]string{"theme", "export", "ghostty", path})
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("code=%d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "palette = 15=#FFFFFF") ||
		!strings.Contains(stdout.String(), "cursor-color = #FFC777") {
		t.Fatalf("stdout does not contain complete Ghostty theme:\n%s", stdout.String())
	}
}

func TestExportGhosttyWritesNewOutputFile(t *testing.T) {
	themePath := writeTheme(t, readPlaceholder(t))
	outputPath := filepath.Join(t.TempDir(), "ghostty-theme")
	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), strings.NewReader(""), &stdout, &stderr,
		[]string{"theme", "export", "ghostty", "--output", outputPath, themePath})
	if code != 0 || stdout.Len() != 0 || stderr.Len() != 0 {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	got, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(got, []byte("background = #222436")) {
		t.Fatalf("output = %q, want background", got)
	}
}

func TestExportGhosttyRefusesExistingOutputFile(t *testing.T) {
	themePath := writeTheme(t, readPlaceholder(t))
	outputPath := filepath.Join(t.TempDir(), "ghostty-theme")
	if err := os.WriteFile(outputPath, []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), strings.NewReader(""), &stdout, &stderr,
		[]string{"theme", "export", "ghostty", "--output", outputPath, themePath})
	if code != 1 || !strings.Contains(stderr.String(), "destination exists") {
		t.Fatalf("code=%d stderr=%q", code, stderr.String())
	}
	got, _ := os.ReadFile(outputPath)
	if string(got) != "original" {
		t.Fatalf("existing output changed to %q", got)
	}
}

func TestEditCommandLoadsExplicitThemeAndRunsEditor(t *testing.T) {
	path := writeTheme(t, readPlaceholder(t))
	originalRunEditor := runEditor
	t.Cleanup(func() { runEditor = originalRunEditor })
	var gotTheme theme.Theme
	var gotPath string
	runEditor = func(value theme.Theme, path string) error {
		gotTheme = value
		gotPath = path
		return nil
	}

	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), strings.NewReader(""), &stdout, &stderr,
		[]string{"theme", "edit", path})
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("code=%d stderr=%q", code, stderr.String())
	}
	if gotPath != path || gotTheme.Theme.Name != "yolodev" {
		t.Fatalf("editor got path=%q theme=%#v", gotPath, gotTheme.Theme)
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
