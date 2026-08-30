package theme

import (
	"math"
	"strings"
	"testing"
)

func TestValidateAcceptsPlaceholder(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	value := decodedValidTheme(t)
	value.Version = 2
	got := Validate(value)
	if !hasDiagnostic(got, Error, "version") {
		t.Fatalf("diagnostics = %#v, want version error", got)
	}
}

func TestValidateRejectsInvalidField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		field  string
		mutate func(*Theme)
	}{
		{
			name:  "missing description",
			field: "theme.description",
			mutate: func(value *Theme) {
				value.Theme.Description = " "
			},
		},
		{
			name:  "malformed color",
			field: "palette.normal.red",
			mutate: func(value *Theme) {
				value.Palette.Normal.Red = "red"
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			value := decodedValidTheme(t)
			test.mutate(&value)
			got := Validate(value)
			if !hasDiagnostic(got, Error, test.field) {
				t.Errorf("diagnostics = %#v, want error for %s", got, test.field)
			}
		})
	}
}

func TestValidateContrastWarningsDoNotBecomeErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		field  string
		mutate func(*Theme)
	}{
		{
			name:  "foreground and background",
			field: "colors.foreground",
			mutate: func(value *Theme) {
				value.Colors.Foreground = value.Colors.Background
			},
		},
		{
			name:  "selection colors",
			field: "colors.selection_foreground",
			mutate: func(value *Theme) {
				value.Colors.SelectionForeground = value.Colors.SelectionBackground
			},
		},
		{
			name:  "cursor and background",
			field: "colors.cursor",
			mutate: func(value *Theme) {
				value.Colors.Cursor = value.Colors.Background
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			value := decodedValidTheme(t)
			test.mutate(&value)
			got := Validate(value)
			if !hasDiagnostic(got, Warning, test.field) {
				t.Errorf("diagnostics = %#v, want warning for %s", got, test.field)
			}
			if HasErrors(got) {
				t.Fatalf("contrast diagnostics unexpectedly contain errors: %#v", got)
			}
		})
	}
}

func TestContrastRatioBlackAndWhite(t *testing.T) {
	t.Parallel()

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
