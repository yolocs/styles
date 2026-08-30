package tui

import "github.com/yolocs/styles/yolodev/internal/theme"

type DragTarget int

const (
	DragNone DragTarget = iota
	DragSV
	DragHue
)

type Picker struct {
	HSV      theme.HSV
	Drag     DragTarget
	StartHex string
	StartRef ColorRef
}

func svAt(region Region, x, y int) (saturation, value float64) {
	x = clampInt(x, region.X, region.X+region.Width-1)
	y = clampInt(y, region.Y, region.Y+region.Height-1)
	if region.Width > 1 {
		saturation = float64(x-region.X) / float64(region.Width-1)
	}
	value = 1
	if region.Height > 1 {
		value = 1 - float64(y-region.Y)/float64(region.Height-1)
	}
	return saturation, value
}

func hueAt(region Region, x int) float64 {
	x = clampInt(x, region.X, region.X+region.Width-1)
	if region.Width <= 1 {
		return 0
	}
	return float64(x-region.X) / float64(region.Width-1) * 359
}

func clampInt(value, minimum, maximum int) int {
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}
	return value
}
