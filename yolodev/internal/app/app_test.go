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

func TestValidateCommand(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                   string
		old, new               []byte
		wantCode               int
		wantStdout, wantStderr string
	}{
		{
			name:       "returns zero with warnings",
			old:        []byte(`foreground = "#C8D3F5"`),
			new:        []byte(`foreground = "#222436"`),
			wantCode:   0,
			wantStdout: ":colors.foreground: WARNING:",
		},
		{
			name:       "returns one for invalid theme",
			old:        []byte("version = 1"),
			new:        []byte("version = 2"),
			wantCode:   1,
			wantStderr: ":version: ERROR:",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			input := bytes.Replace(readPlaceholder(t), test.old, test.new, 1)
			path := writeTheme(t, input)
			var stdout, stderr bytes.Buffer
			code := Run(context.Background(), strings.NewReader(""), &stdout, &stderr,
				[]string{"theme", "validate", path})
			if code != test.wantCode {
				t.Fatalf("Run() code = %d, want %d", code, test.wantCode)
			}
			assertOutput(t, "stdout", stdout.String(), test.wantStdout)
			assertOutput(t, "stderr", stderr.String(), test.wantStderr)
		})
	}
}

func TestUsageReturnsTwo(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), strings.NewReader(""), &stdout, &stderr, nil)
	if code != 2 {
		t.Fatalf("Run() code = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "Usage:") {
		t.Fatalf("stderr = %q, want usage", stderr.String())
	}
}

func TestExportGhosttyWritesStdout(t *testing.T) {
	t.Parallel()

	path := writeTheme(t, readPlaceholder(t))
	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), strings.NewReader(""), &stdout, &stderr,
		[]string{"theme", "export", "--format", "ghostty", path})
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("code=%d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "palette = 15=#FFFFFF") ||
		!strings.Contains(stdout.String(), "cursor-color = #FFC777") {
		t.Fatalf("stdout does not contain complete Ghostty theme:\n%s", stdout.String())
	}
}

func TestExportCodexWritesStdout(t *testing.T) {
	t.Parallel()

	path := writeTheme(t, readSignalGrid(t))
	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), strings.NewReader(""), &stdout, &stderr,
		[]string{"theme", "export", "--format", "codex", path})
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("code=%d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "<plist version=\"1.0\">") ||
		!strings.Contains(stdout.String(), "<string>#D092FF</string>") ||
		!strings.Contains(stdout.String(), "<string>entity.name.function, support.function, variable.function</string>") {
		t.Fatalf("stdout does not contain complete Codex theme:\n%s", stdout.String())
	}
}

func TestExportGhosttyToFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		initial    string
		wantCode   int
		wantStderr string
		wantOutput string
	}{
		{name: "writes new file", wantCode: 0, wantOutput: "background = #222436"},
		{name: "refuses existing file", initial: "original", wantCode: 1, wantStderr: "destination exists", wantOutput: "original"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			themePath := writeTheme(t, readPlaceholder(t))
			outputPath := filepath.Join(t.TempDir(), "ghostty-theme.ghostty")
			if test.initial != "" {
				if err := os.WriteFile(outputPath, []byte(test.initial), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			var stdout, stderr bytes.Buffer
			code := Run(context.Background(), strings.NewReader(""), &stdout, &stderr,
				[]string{"theme", "export", "--format", "ghostty", "--output", outputPath, themePath})
			if code != test.wantCode || stdout.Len() != 0 {
				t.Fatalf("code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
			}
			assertOutput(t, "stderr", stderr.String(), test.wantStderr)
			got, err := os.ReadFile(outputPath)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(got), test.wantOutput) {
				t.Fatalf("output = %q, want %q", got, test.wantOutput)
			}
		})
	}
}

func TestExportCodexToFolder(t *testing.T) {
	t.Parallel()

	themePath := writeTheme(t, readSignalGrid(t))
	outputDirectory := filepath.Join(t.TempDir(), "nested", "themes")
	var stdout, stderr bytes.Buffer
	code := Run(context.Background(), strings.NewReader(""), &stdout, &stderr,
		[]string{"theme", "export", "--format", "codex", "--output", outputDirectory, themePath})
	if code != 0 || stdout.Len() != 0 || stderr.Len() != 0 {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	outputPath := filepath.Join(outputDirectory, "signal-grid.tmTheme")
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(data, []byte("<?xml")) {
		t.Fatalf("output %s does not contain a Codex theme:\n%s", outputPath, data)
	}
}

func TestExportRejectsInvalidFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name, format string
		wantCode     int
		wantStderr   string
	}{
		{name: "missing format", wantCode: 2, wantStderr: "format"},
		{name: "unknown format", format: "unknown", wantCode: 1, wantStderr: `unsupported export format "unknown"`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			themePath := writeTheme(t, readPlaceholder(t))
			args := []string{"theme", "export"}
			if test.format != "" {
				args = append(args, "--format", test.format)
			}
			args = append(args, themePath)
			var stdout, stderr bytes.Buffer
			code := Run(context.Background(), strings.NewReader(""), &stdout, &stderr, args)
			if code != test.wantCode {
				t.Fatalf("Run() code = %d, want %d", code, test.wantCode)
			}
			if !strings.Contains(stderr.String(), test.wantStderr) {
				t.Fatalf("stderr = %q, want %q", stderr.String(), test.wantStderr)
			}
		})
	}
}

func TestEditCommandLoadsExplicitThemeAndRunsEditor(t *testing.T) {
	// This test replaces package-level state and must remain serial.
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

func readSignalGrid(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "..", "themes", "yolodev", "signal-grid.toml"))
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

func assertOutput(t *testing.T, name, got, want string) {
	t.Helper()
	if want == "" && got != "" {
		t.Fatalf("%s = %q, want empty", name, got)
	}
	if want != "" && !strings.Contains(got, want) {
		t.Fatalf("%s = %q, want %q", name, got, want)
	}
}
