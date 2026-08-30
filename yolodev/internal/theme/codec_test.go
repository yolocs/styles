package theme

import (
	"bytes"
	"strings"
	"testing"
)

const validTOML = `format = "yolodev-terminal-theme"
version = 1

[theme]
name = "yolodev"
variant = "placeholder"
appearance = "dark"
description = "Temporary palette used while yolodev is bootstrapped."

[colors]
background = "#222436"
foreground = "#C8D3F5"
cursor = "#FFC777"
cursor_text = "#222436"
selection_background = "#444A73"
selection_foreground = "#C8D3F5"

[palette.normal]
black = "#1B1D2B"
red = "#FF757F"
green = "#C3E88D"
yellow = "#FFC777"
blue = "#82AAFF"
magenta = "#C099FF"
cyan = "#86E1FC"
white = "#C8D3F5"

[palette.bright]
black = "#444A73"
red = "#FF757F"
green = "#C3E88D"
yellow = "#FFC777"
blue = "#82AAFF"
magenta = "#C099FF"
cyan = "#86E1FC"
white = "#FFFFFF"
`

func TestDecodeValidTheme(t *testing.T) {
	t.Parallel()

	got, err := Decode(strings.NewReader(validTOML))
	if err != nil {
		t.Fatal(err)
	}
	if got.Format != Format || got.Version != Version {
		t.Fatalf("format/version = %q/%d, want %q/%d", got.Format, got.Version, Format, Version)
	}
	if got.Theme.Name != "yolodev" || got.Palette.Bright.Blue != "#82AAFF" {
		t.Fatalf("decoded unexpected theme: %#v", got)
	}
}

func TestDecodeRejectsInvalidDocument(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		old, new string
	}{
		{name: "unknown field", old: "appearance = \"dark\"", new: "appearance = \"dark\"\nunknown = true"},
		{name: "duplicate field", old: "version = 1", new: "version = 1\nversion = 1"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			input := strings.Replace(validTOML, test.old, test.new, 1)
			if _, err := Decode(strings.NewReader(input)); err == nil {
				t.Fatal("Decode() error = nil")
			}
		})
	}
}

func TestEncodeNormalizesHex(t *testing.T) {
	t.Parallel()

	value, err := Decode(strings.NewReader(strings.ToLower(validTOML)))
	if err != nil {
		t.Fatal(err)
	}

	var got bytes.Buffer
	if err := Encode(&got, value); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got.String(), `foreground = '#C8D3F5'`) &&
		!strings.Contains(got.String(), `foreground = "#C8D3F5"`) {
		t.Fatalf("encoded theme did not normalize foreground:\n%s", got.String())
	}

	roundTrip, err := Decode(strings.NewReader(got.String()))
	if err != nil {
		t.Fatal(err)
	}
	if roundTrip != Normalize(value) {
		t.Fatalf("round trip mismatch:\n got %#v\nwant %#v", roundTrip, Normalize(value))
	}
}
