package ghostty

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yolocs/styles/yolodev/internal/theme"
)

func TestRenderPlaceholder(t *testing.T) {
	t.Parallel()

	value := loadPlaceholder(t)
	var got bytes.Buffer
	if err := Render(&got, value); err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join("testdata", "placeholder.golden"))
	if err != nil {
		t.Fatal(err)
	}
	if got.String() != string(want) {
		t.Fatalf("Ghostty output mismatch:\n--- got ---\n%s--- want ---\n%s", got.String(), want)
	}
}

func TestRenderRejectsInvalidTheme(t *testing.T) {
	t.Parallel()

	value := loadPlaceholder(t)
	value.Version = 2
	var output bytes.Buffer
	if err := Render(&output, value); err == nil || !strings.Contains(err.Error(), "version") {
		t.Fatalf("Render() error = %v, want version validation error", err)
	}
	if output.Len() != 0 {
		t.Fatalf("Render() wrote partial output %q", output.String())
	}
}

func loadPlaceholder(t *testing.T) theme.Theme {
	t.Helper()
	file, err := os.Open(filepath.Join("..", "..", "..", "themes", "yolodev", "placeholder.toml"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	value, err := theme.Decode(file)
	if err != nil {
		t.Fatal(err)
	}
	return value
}
