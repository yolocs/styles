package tui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/yolocs/styles/yolodev/internal/theme"
)

func TestMouseClickAndDragUpdateSelectedColor(t *testing.T) {
	t.Parallel()

	model := testModel(t)
	model.Regions = map[string]Region{
		"sv":  {X: 10, Y: 5, Width: 11, Height: 11},
		"hue": {X: 10, Y: 17, Width: 11, Height: 1},
	}
	original := model.Theme.Colors.Background

	updated := updateModel(t, model, tea.MouseClickMsg{X: 20, Y: 5, Button: tea.MouseLeft})
	if updated.Picker.Drag != DragSV || updated.Theme.Colors.Background == original {
		t.Fatalf("click did not begin SV drag: %#v", updated.Picker)
	}

	updated = updateModel(t, updated, tea.MouseMotionMsg{X: 15, Y: 10, Button: tea.MouseLeft})
	hsv := theme.RGBToHSV(mustRGB(t, updated.Theme.Colors.Background))
	if hsv.S < 0.49 || hsv.S > 0.51 || hsv.V < 0.49 || hsv.V > 0.51 {
		t.Fatalf("dragged HSV = %#v, want S/V near 0.5", hsv)
	}

	updated = updateModel(t, updated, tea.MouseReleaseMsg{X: 15, Y: 10, Button: tea.MouseLeft})
	if updated.Picker.Drag != DragNone {
		t.Fatalf("release left drag active: %v", updated.Picker.Drag)
	}
}

func TestEscapeCancelsActiveDrag(t *testing.T) {
	t.Parallel()

	model := testModel(t)
	model.Regions = map[string]Region{"hue": {X: 10, Y: 5, Width: 11, Height: 1}}
	original := model.Theme.Colors.Background
	updated := updateModel(t, model, tea.MouseClickMsg{X: 20, Y: 5, Button: tea.MouseLeft})
	updated = updateModel(t, updated, tea.KeyPressMsg{Code: tea.KeyEscape})
	if updated.Theme.Colors.Background != original || updated.Picker.Drag != DragNone {
		t.Fatalf("Escape left color=%s drag=%v, want %s/no drag", updated.Theme.Colors.Background, updated.Picker.Drag, original)
	}
}

func TestClickingSwatchChangesSelection(t *testing.T) {
	t.Parallel()

	model := testModel(t)
	model.Regions = map[string]Region{
		"color:palette.normal.red": {X: 2, Y: 2, Width: 3, Height: 1},
	}
	updated := updateModel(t, model, tea.MouseClickMsg{X: 3, Y: 2, Button: tea.MouseLeft})
	if updated.Selected != (ColorRef{Group: "palette.normal", Name: "red"}) {
		t.Fatalf("Selected = %#v", updated.Selected)
	}
}

func TestWindowSizeRecalculatesLayoutAndViewEnablesMouse(t *testing.T) {
	t.Parallel()

	model := updateModel(t, testModel(t), tea.WindowSizeMsg{Width: 120, Height: 36})
	if model.Width != 120 || model.Height != 36 || len(model.Regions) == 0 {
		t.Fatalf("window update = %dx%d, regions=%d", model.Width, model.Height, len(model.Regions))
	}
	view := model.View()
	if !view.AltScreen || view.MouseMode != tea.MouseModeCellMotion {
		t.Fatalf("View flags: alt=%v mouse=%v", view.AltScreen, view.MouseMode)
	}
}

func TestDirtyComparesCurrentThemeToCleanSnapshot(t *testing.T) {
	t.Parallel()

	model := testModel(t)
	if model.Dirty() {
		t.Fatal("new model is dirty")
	}
	model.Theme.Colors.Foreground = "#FFFFFF"
	if !model.Dirty() {
		t.Fatal("changed model is not dirty")
	}
}

func testModel(t *testing.T) Model {
	t.Helper()
	value := theme.Theme{
		Format:  theme.Format,
		Version: theme.Version,
		Theme:   theme.Metadata{Name: "test", Variant: "test", Appearance: "dark", Description: "test theme"},
		Colors: theme.SemanticColors{
			Background: "#222436", Foreground: "#C8D3F5", Cursor: "#FFC777", CursorText: "#222436",
			SelectionBackground: "#444A73", SelectionForeground: "#C8D3F5",
		},
		Palette: theme.Palette{
			Normal: theme.ANSIColors{Black: "#000000", Red: "#AA0000", Green: "#00AA00", Yellow: "#AAAA00", Blue: "#0000AA", Magenta: "#AA00AA", Cyan: "#00AAAA", White: "#AAAAAA"},
			Bright: theme.ANSIColors{Black: "#555555", Red: "#FF5555", Green: "#55FF55", Yellow: "#FFFF55", Blue: "#5555FF", Magenta: "#FF55FF", Cyan: "#55FFFF", White: "#FFFFFF"},
		},
	}
	return New(value, "theme.toml")
}

func updateModel(t *testing.T, model Model, message tea.Msg) Model {
	t.Helper()
	updated, _ := model.Update(message)
	got, ok := updated.(Model)
	if !ok {
		t.Fatalf("Update returned %T, want Model", updated)
	}
	return got
}

func mustRGB(t *testing.T, value string) theme.RGB {
	t.Helper()
	got, err := theme.ParseHex(value)
	if err != nil {
		t.Fatal(err)
	}
	return got
}
