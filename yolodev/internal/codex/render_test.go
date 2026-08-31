package codex

import (
	"bytes"
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yolocs/styles/yolodev/internal/theme"
)

func TestRenderSignalGrid(t *testing.T) {
	t.Parallel()

	value := loadSignalGrid(t)
	var got bytes.Buffer
	if err := Render(&got, value); err != nil {
		t.Fatal(err)
	}

	want, err := os.ReadFile(filepath.Join("testdata", "signal-grid.tmTheme.golden"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got.Bytes(), want) {
		t.Fatalf("Render() output:\n%s\nwant:\n%s", got.Bytes(), want)
	}

	var document any
	if err := xml.Unmarshal(got.Bytes(), &document); err != nil {
		t.Fatalf("Render() produced invalid XML: %v", err)
	}
}

func TestRenderRejectsInvalidThemeWithoutWriting(t *testing.T) {
	t.Parallel()

	value := loadSignalGrid(t)
	value.Version++
	output := bytes.NewBufferString("untouched")
	err := Render(output, value)
	if err == nil || !strings.Contains(err.Error(), "version: must equal 1") {
		t.Fatalf("Render() error = %v, want invalid version", err)
	}
	if got := output.String(); got != "untouched" {
		t.Fatalf("Render() output = %q, want untouched", got)
	}
}

func loadSignalGrid(t *testing.T) theme.Theme {
	t.Helper()

	file, err := os.Open(filepath.Join("..", "..", "..", "themes", "yolodev", "signal-grid.toml"))
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
