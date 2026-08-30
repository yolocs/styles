package theme

import (
	"math"
	"strings"
	"testing"
)

func TestValidateAcceptsPlaceholder(t *testing.T) {
	value, err := Decode(strings.NewReader(validTOML))
	if err != nil {
		t.Fatal(err)
	}
	got := Validate(value)
	if HasErrors(got) {
		t.Fatalf("Validate() returned errors: %#v", got)
	}
}

func TestValidateRejectsUnsupportedVersion(t *testing.T) {
	value := decodedValidTheme(t)
	value.Version = 2
	got := Validate(value)
	if !hasDiagnostic(got, Error, "version") {
		t.Fatalf("diagnostics = %#v, want version error", got)
	}
}

func TestValidateRejectsMissingMetadataAndMalformedColor(t *testing.T) {
	value := decodedValidTheme(t)
	value.Theme.Description = " "
	value.Palette.Normal.Red = "red"
	got := Validate(value)
	if !hasDiagnostic(got, Error, "theme.description") {
		t.Errorf("diagnostics = %#v, want description error", got)
	}
	if !hasDiagnostic(got, Error, "palette.normal.red") {
		t.Errorf("diagnostics = %#v, want red error", got)
	}
}

func TestValidateContrastWarningsDoNotBecomeErrors(t *testing.T) {
	value := decodedValidTheme(t)
	value.Colors.Foreground = value.Colors.Background
	value.Colors.SelectionForeground = value.Colors.SelectionBackground
	value.Colors.Cursor = value.Colors.Background

	got := Validate(value)
	for _, field := range []string{"colors.foreground", "colors.selection_foreground", "colors.cursor"} {
		if !hasDiagnostic(got, Warning, field) {
			t.Errorf("diagnostics = %#v, want warning for %s", got, field)
		}
	}
	if HasErrors(got) {
		t.Fatalf("contrast diagnostics unexpectedly contain errors: %#v", got)
	}
}

func TestContrastRatioBlackAndWhite(t *testing.T) {
	got := ContrastRatio(RGB{}, RGB{R: 255, G: 255, B: 255})
	if math.Abs(got-21) > 0.001 {
		t.Fatalf("ContrastRatio(black, white) = %f, want 21", got)
	}
}

func decodedValidTheme(t *testing.T) Theme {
	t.Helper()
	value, err := Decode(strings.NewReader(validTOML))
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func hasDiagnostic(values []Diagnostic, severity Severity, field string) bool {
	for _, value := range values {
		if value.Severity == severity && value.Field == field {
			return true
		}
	}
	return false
}
