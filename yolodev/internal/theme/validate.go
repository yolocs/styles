package theme

import (
	"fmt"
	"math"
	"strings"
)

type Severity string

const (
	Error   Severity = "error"
	Warning Severity = "warning"
)

type Diagnostic struct {
	Severity Severity
	Field    string
	Message  string
}

type colorField struct {
	path  string
	value string
}

func Validate(value Theme) []Diagnostic {
	var diagnostics []Diagnostic
	if value.Format != Format {
		diagnostics = append(diagnostics, Diagnostic{Error, "format", fmt.Sprintf("must equal %q", Format)})
	}
	if value.Version != Version {
		diagnostics = append(diagnostics, Diagnostic{Error, "version", fmt.Sprintf("must equal %d", Version)})
	}
	if strings.TrimSpace(value.Theme.Name) == "" {
		diagnostics = append(diagnostics, Diagnostic{Error, "theme.name", "must not be empty"})
	}
	if strings.TrimSpace(value.Theme.Variant) == "" {
		diagnostics = append(diagnostics, Diagnostic{Error, "theme.variant", "must not be empty"})
	}
	if value.Theme.Appearance != "dark" && value.Theme.Appearance != "light" {
		diagnostics = append(diagnostics, Diagnostic{Error, "theme.appearance", `must equal "dark" or "light"`})
	}
	if strings.TrimSpace(value.Theme.Description) == "" {
		diagnostics = append(diagnostics, Diagnostic{Error, "theme.description", "must not be empty"})
	}

	for _, field := range allColorFields(value) {
		if _, err := ParseHex(field.value); err != nil {
			diagnostics = append(diagnostics, Diagnostic{Error, field.path, err.Error()})
		}
	}

	diagnostics = appendContrastWarning(diagnostics, value.Colors.Foreground, value.Colors.Background, 4.5, "colors.foreground", "foreground/background")
	diagnostics = appendContrastWarning(diagnostics, value.Colors.SelectionForeground, value.Colors.SelectionBackground, 3.0, "colors.selection_foreground", "selection foreground/background")
	diagnostics = appendContrastWarning(diagnostics, value.Colors.Cursor, value.Colors.Background, 3.0, "colors.cursor", "cursor/background")
	return diagnostics
}

func HasErrors(diagnostics []Diagnostic) bool {
	for _, diagnostic := range diagnostics {
		if diagnostic.Severity == Error {
			return true
		}
	}
	return false
}

func ContrastRatio(a, b RGB) float64 {
	left := relativeLuminance(a)
	right := relativeLuminance(b)
	if left < right {
		left, right = right, left
	}
	return (left + 0.05) / (right + 0.05)
}

func relativeLuminance(c RGB) float64 {
	channel := func(value uint8) float64 {
		normalized := float64(value) / 255
		if normalized <= 0.04045 {
			return normalized / 12.92
		}
		return math.Pow((normalized+0.055)/1.055, 2.4)
	}
	return 0.2126*channel(c.R) + 0.7152*channel(c.G) + 0.0722*channel(c.B)
}

func appendContrastWarning(diagnostics []Diagnostic, foreground, background string, threshold float64, field, label string) []Diagnostic {
	fg, fgErr := ParseHex(foreground)
	bg, bgErr := ParseHex(background)
	if fgErr != nil || bgErr != nil {
		return diagnostics
	}
	ratio := ContrastRatio(fg, bg)
	if ratio < threshold {
		return append(diagnostics, Diagnostic{
			Severity: Warning,
			Field:    field,
			Message:  fmt.Sprintf("%s contrast is %.2f:1; recommended minimum is %.1f:1", label, ratio, threshold),
		})
	}
	return diagnostics
}

func allColorFields(value Theme) []colorField {
	fields := []colorField{
		{"colors.background", value.Colors.Background},
		{"colors.foreground", value.Colors.Foreground},
		{"colors.cursor", value.Colors.Cursor},
		{"colors.cursor_text", value.Colors.CursorText},
		{"colors.selection_background", value.Colors.SelectionBackground},
		{"colors.selection_foreground", value.Colors.SelectionForeground},
	}
	fields = append(fields, ansiColorFields("palette.normal", value.Palette.Normal)...)
	fields = append(fields, ansiColorFields("palette.bright", value.Palette.Bright)...)
	return fields
}

func ansiColorFields(prefix string, value ANSIColors) []colorField {
	return []colorField{
		{prefix + ".black", value.Black},
		{prefix + ".red", value.Red},
		{prefix + ".green", value.Green},
		{prefix + ".yellow", value.Yellow},
		{prefix + ".blue", value.Blue},
		{prefix + ".magenta", value.Magenta},
		{prefix + ".cyan", value.Cyan},
		{prefix + ".white", value.White},
	}
}
