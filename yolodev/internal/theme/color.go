package theme

import (
	"fmt"
	"math"
	"strconv"
)

type RGB struct {
	R uint8
	G uint8
	B uint8
}

type HSV struct {
	H float64
	S float64
	V float64
}

func ParseHex(value string) (RGB, error) {
	if len(value) != 7 || value[0] != '#' {
		return RGB{}, fmt.Errorf("color %q must use #RRGGBB", value)
	}
	n, err := strconv.ParseUint(value[1:], 16, 24)
	if err != nil {
		return RGB{}, fmt.Errorf("color %q must use hexadecimal digits: %w", value, err)
	}
	return RGB{R: uint8(n >> 16), G: uint8(n >> 8), B: uint8(n)}, nil
}

func (c RGB) Hex() string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

func RGBToHSV(c RGB) HSV {
	r := float64(c.R) / 255
	g := float64(c.G) / 255
	b := float64(c.B) / 255
	maxValue := math.Max(r, math.Max(g, b))
	minValue := math.Min(r, math.Min(g, b))
	delta := maxValue - minValue

	hue := 0.0
	if delta != 0 {
		switch maxValue {
		case r:
			hue = 60 * math.Mod((g-b)/delta, 6)
		case g:
			hue = 60 * ((b-r)/delta + 2)
		default:
			hue = 60 * ((r-g)/delta + 4)
		}
	}
	if hue < 0 {
		hue += 360
	}

	saturation := 0.0
	if maxValue != 0 {
		saturation = delta / maxValue
	}
	return HSV{H: hue, S: saturation, V: maxValue}
}

func HSVToRGB(c HSV) RGB {
	hue := math.Mod(c.H, 360)
	if hue < 0 {
		hue += 360
	}
	saturation := clampUnit(c.S)
	value := clampUnit(c.V)
	chroma := value * saturation
	x := chroma * (1 - math.Abs(math.Mod(hue/60, 2)-1))
	m := value - chroma

	var r, g, b float64
	switch {
	case hue < 60:
		r, g = chroma, x
	case hue < 120:
		r, g = x, chroma
	case hue < 180:
		g, b = chroma, x
	case hue < 240:
		g, b = x, chroma
	case hue < 300:
		r, b = x, chroma
	default:
		r, b = chroma, x
	}

	return RGB{
		R: uint8(math.Round((r + m) * 255)),
		G: uint8(math.Round((g + m) * 255)),
		B: uint8(math.Round((b + m) * 255)),
	}
}

func clampUnit(value float64) float64 {
	return math.Max(0, math.Min(1, value))
}
