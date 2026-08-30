package tui

import (
	"math"
	"testing"
)

func TestSVAtUsesClickedCell(t *testing.T) {
	region := Region{X: 10, Y: 5, Width: 11, Height: 11}
	tests := []struct {
		x, y int
		s, v float64
	}{
		{x: 10, y: 5, s: 0, v: 1},
		{x: 20, y: 5, s: 1, v: 1},
		{x: 20, y: 15, s: 1, v: 0},
		{x: 15, y: 10, s: 0.5, v: 0.5},
	}
	for _, test := range tests {
		gotS, gotV := svAt(region, test.x, test.y)
		if math.Abs(gotS-test.s) > 0.001 || math.Abs(gotV-test.v) > 0.001 {
			t.Errorf("svAt(%d,%d) = %.2f,%.2f, want %.2f,%.2f", test.x, test.y, gotS, gotV, test.s, test.v)
		}
	}
}

func TestHueAtClampsOutsideRegion(t *testing.T) {
	region := Region{X: 10, Y: 5, Width: 11, Height: 1}
	if got := hueAt(region, 0); got != 0 {
		t.Fatalf("hueAt(left) = %f, want 0", got)
	}
	if got := hueAt(region, 99); got != 359 {
		t.Fatalf("hueAt(right) = %f, want 359", got)
	}
}
