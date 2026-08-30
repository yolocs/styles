package tui

import "testing"

func TestRegionContainsBoundary(t *testing.T) {
	t.Parallel()

	region := Region{X: 10, Y: 5, Width: 20, Height: 8}
	tests := []struct {
		name       string
		x, y       int
		wantInside bool
	}{
		{name: "upper left", x: 10, y: 5, wantInside: true},
		{name: "lower right", x: 29, y: 12, wantInside: true},
		{name: "past right", x: 30, y: 12, wantInside: false},
		{name: "past bottom", x: 29, y: 13, wantInside: false},
		{name: "before left", x: 9, y: 5, wantInside: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := region.Contains(test.x, test.y); got != test.wantInside {
				t.Fatalf("Region.Contains(%d, %d) = %v, want %v", test.x, test.y, got, test.wantInside)
			}
		})
	}
}

func TestLayoutForWideTerminalIncludesPickerAndColors(t *testing.T) {
	t.Parallel()

	layout := LayoutFor(120, 36)
	if layout.Small {
		t.Fatal("LayoutFor(120, 36) marked terminal small")
	}
	for _, key := range []string{
		"sv", "hue", "hex", "action:import", "action:save", "action:export",
		"color:colors.background", "color:palette.normal.red", "color:palette.bright.white",
	} {
		if region, ok := layout.Regions[key]; !ok || region.Width < 1 || region.Height < 1 {
			t.Errorf("layout region %q = %#v, present=%v", key, region, ok)
		}
	}
}

func TestLayoutForSmallTerminal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		width, height int
	}{
		{name: "too narrow", width: 99, height: 28},
		{name: "too short", width: 100, height: 27},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := LayoutFor(test.width, test.height); !got.Small {
				t.Fatalf("LayoutFor(%d, %d) did not use small layout", test.width, test.height)
			}
		})
	}
}
