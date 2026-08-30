package theme

import (
	"math"
	"testing"
)

func TestParseHex(t *testing.T) {
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
	for _, input := range []string{"82AAFF", "#fff", "#GGAAFF", "#82AAFF00"} {
		if _, err := ParseHex(input); err == nil {
			t.Errorf("ParseHex(%q) error = nil", input)
		}
	}
}

func TestHSVRoundTrip(t *testing.T) {
	inputs := []RGB{
		{R: 0, G: 0, B: 0},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 0, B: 0},
		{R: 130, G: 170, B: 255},
	}
	for _, input := range inputs {
		got := HSVToRGB(RGBToHSV(input))
		if channelDistance(got, input) > 1 {
			t.Errorf("HSV round trip of %#v = %#v", input, got)
		}
	}
}

func channelDistance(a, b RGB) int {
	return int(math.Abs(float64(a.R)-float64(b.R)) +
		math.Abs(float64(a.G)-float64(b.G)) +
		math.Abs(float64(a.B)-float64(b.B)))
}
