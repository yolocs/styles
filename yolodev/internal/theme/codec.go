package theme

import (
	"fmt"
	"io"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

func Decode(r io.Reader) (Theme, error) {
	var value Theme
	decoder := toml.NewDecoder(r)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&value); err != nil {
		return Theme{}, fmt.Errorf("decode terminal theme: %w", err)
	}
	return value, nil
}

func Encode(w io.Writer, value Theme) error {
	if err := toml.NewEncoder(w).Encode(Normalize(value)); err != nil {
		return fmt.Errorf("encode terminal theme: %w", err)
	}
	return nil
}

func Normalize(value Theme) Theme {
	value.Theme.Name = strings.TrimSpace(value.Theme.Name)
	value.Theme.Variant = strings.TrimSpace(value.Theme.Variant)
	value.Theme.Appearance = strings.ToLower(strings.TrimSpace(value.Theme.Appearance))
	value.Theme.Description = strings.TrimSpace(value.Theme.Description)

	value.Colors.Background = normalizeHex(value.Colors.Background)
	value.Colors.Foreground = normalizeHex(value.Colors.Foreground)
	value.Colors.Cursor = normalizeHex(value.Colors.Cursor)
	value.Colors.CursorText = normalizeHex(value.Colors.CursorText)
	value.Colors.SelectionBackground = normalizeHex(value.Colors.SelectionBackground)
	value.Colors.SelectionForeground = normalizeHex(value.Colors.SelectionForeground)
	value.Palette.Normal = normalizeANSI(value.Palette.Normal)
	value.Palette.Bright = normalizeANSI(value.Palette.Bright)
	return value
}

func normalizeANSI(value ANSIColors) ANSIColors {
	value.Black = normalizeHex(value.Black)
	value.Red = normalizeHex(value.Red)
	value.Green = normalizeHex(value.Green)
	value.Yellow = normalizeHex(value.Yellow)
	value.Blue = normalizeHex(value.Blue)
	value.Magenta = normalizeHex(value.Magenta)
	value.Cyan = normalizeHex(value.Cyan)
	value.White = normalizeHex(value.White)
	return value
}

func normalizeHex(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}
