package tui

import (
	"math"
	"testing"
)

func TestSVAtUsesClickedCell(t *testing.T) {
	t.Parallel()

	region := Region{X: 10, Y: 5, Width: 11, Height: 11}
	tests := []struct {
		name string
		x, y int
		s, v float64
	}{
		{name: "upper left", x: 10, y: 5, s: 0, v: 1},
		{name: "upper right", x: 20, y: 5, s: 1, v: 1},
		{name: "lower right", x: 20, y: 15, s: 1, v: 0},
		{name: "center", x: 15, y: 10, s: 0.5, v: 0.5},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			gotS, gotV := svAt(region, test.x, test.y)
			if math.Abs(gotS-test.s) > 0.001 || math.Abs(gotV-test.v) > 0.001 {
				t.Errorf("svAt(%d,%d) = %.2f,%.2f, want %.2f,%.2f", test.x, test.y, gotS, gotV, test.s, test.v)
			}
		})
	}
}

func TestHueAtClampsOutsideRegion(t *testing.T) {
	t.Parallel()

	region := Region{X: 10, Y: 5, Width: 11, Height: 1}
	tests := []struct {
		name string
		x    int
		want float64
	}{
		{name: "left", x: 0, want: 0},
		{name: "right", x: 99, want: 359},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := hueAt(region, test.x); got != test.want {
				t.Fatalf("hueAt(%d) = %f, want %f", test.x, got, test.want)
			}
		})
	}
}
