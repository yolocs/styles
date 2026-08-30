package exporter

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/yolocs/styles/yolodev/internal/theme"
)

func TestRegistryExportsWithRequestedFormat(t *testing.T) {
	registry := NewRegistry()
	err := registry.Register("custom", func(w io.Writer, value theme.Theme) error {
		_, writeErr := io.WriteString(w, "custom:"+value.Theme.Name)
		return writeErr
	})
	if err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer
	value := theme.Theme{Theme: theme.Metadata{Name: "yolodev"}}
	if err := registry.Export("custom", &output, value); err != nil {
		t.Fatal(err)
	}
	if got, want := output.String(), "custom:yolodev"; got != want {
		t.Fatalf("Export() output = %q, want %q", got, want)
	}
}

func TestRegistryRejectsUnknownFormatWithoutWriting(t *testing.T) {
	registry := NewRegistry()
	var output bytes.Buffer
	err := registry.Export("missing", &output, theme.Theme{})
	if err == nil || !strings.Contains(err.Error(), `unsupported export format "missing"`) {
		t.Fatalf("Export() error = %v, want unsupported format", err)
	}
	if output.Len() != 0 {
		t.Fatalf("Export() wrote %q for unknown format", output.String())
	}
}

func TestRegistryRejectsDuplicateFormat(t *testing.T) {
	registry := NewRegistry()
	render := func(io.Writer, theme.Theme) error { return nil }
	if err := registry.Register("custom", render); err != nil {
		t.Fatal(err)
	}
	if err := registry.Register("custom", render); err == nil {
		t.Fatal("Register() duplicate format returned nil error")
	}
}

func TestRegistryRejectsEmptyFormat(t *testing.T) {
	registry := NewRegistry()
	render := func(io.Writer, theme.Theme) error { return nil }
	if err := registry.Register("", render); err == nil {
		t.Fatal("Register() empty format returned nil error")
	}
}

func TestRegistryRejectsNilRenderer(t *testing.T) {
	registry := NewRegistry()
	if err := registry.Register("custom", nil); err == nil {
		t.Fatal("Register() nil renderer returned nil error")
	}
}
