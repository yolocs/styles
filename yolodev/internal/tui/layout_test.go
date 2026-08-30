package tui

import "testing"

func TestRegionContainsBoundary(t *testing.T) {
	region := Region{X: 10, Y: 5, Width: 20, Height: 8}
	if !region.Contains(10, 5) || !region.Contains(29, 12) {
		t.Fatal("Region.Contains rejected an inclusive boundary cell")
	}
	if region.Contains(30, 12) || region.Contains(29, 13) || region.Contains(9, 5) {
		t.Fatal("Region.Contains accepted a cell outside the region")
	}
}

func TestLayoutForWideTerminalIncludesPickerAndColors(t *testing.T) {
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
	if got := LayoutFor(99, 28); !got.Small {
		t.Fatal("99-column terminal must use small layout")
	}
	if got := LayoutFor(100, 27); !got.Small {
		t.Fatal("27-row terminal must use small layout")
	}
}
