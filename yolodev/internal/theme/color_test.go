package theme

import (
	"math"
	"testing"
)

func TestParseHex(t *testing.T) {
	t.Parallel()

	got, err := ParseHex("#82AAFF")
	if err != nil {
		t.Fatal(err)
	}
	want := RGB{R: 0x82, G: 0xAA, B: 0xFF}
	if got != want || got.Hex() != "#82AAFF" {
		t.Fatalf("ParseHex() = %#v (%s), want %#v (#82AAFF)", got, got.Hex(), want)
	}
}

func TestParseHexRejectsNonCanonicalInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name, input string
	}{
		{name: "missing hash", input: "82AAFF"},
		{name: "short form", input: "#fff"},
		{name: "non-hex digits", input: "#GGAAFF"},
		{name: "alpha channel", input: "#82AAFF00"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if _, err := ParseHex(test.input); err == nil {
				t.Errorf("ParseHex(%q) error = nil", test.input)
			}
		})
	}
}

func TestHSVRoundTrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input RGB
	}{
		{name: "black", input: RGB{R: 0, G: 0, B: 0}},
		{name: "white", input: RGB{R: 255, G: 255, B: 255}},
		{name: "red", input: RGB{R: 255, G: 0, B: 0}},
		{name: "blue", input: RGB{R: 130, G: 170, B: 255}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := HSVToRGB(RGBToHSV(test.input))
			if channelDistance(got, test.input) > 1 {
				t.Errorf("HSV round trip of %#v = %#v", test.input, got)
			}
		})
	}
}

func channelDistance(a, b RGB) int {
	return int(math.Abs(float64(a.R)-float64(b.R)) +
		math.Abs(float64(a.G)-float64(b.G)) +
		math.Abs(float64(a.B)-float64(b.B)))
}
