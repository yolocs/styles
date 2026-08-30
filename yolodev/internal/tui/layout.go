package tui

const (
	MinimumWidth  = 100
	MinimumHeight = 28
)

type Region struct {
	X      int
	Y      int
	Width  int
	Height int
}

func (r Region) Contains(x, y int) bool {
	return r.Width > 0 && r.Height > 0 && x >= r.X && x < r.X+r.Width && y >= r.Y && y < r.Y+r.Height
}

type Layout struct {
	Small   bool
	Dock    Region
	Preview Region
	Regions map[string]Region
}

func LayoutFor(width, height int) Layout {
	if width < MinimumWidth || height < MinimumHeight {
		return Layout{Small: true, Regions: map[string]Region{}}
	}

	dockWidth := 42
	layout := Layout{
		Dock:    Region{X: 0, Y: 1, Width: dockWidth, Height: height - 3},
		Preview: Region{X: dockWidth + 1, Y: 1, Width: width - dockWidth - 1, Height: height - 3},
		Regions: make(map[string]Region),
	}

	semanticNames := []string{"background", "foreground", "cursor", "cursor_text", "selection_background", "selection_foreground"}
	for index, name := range semanticNames {
		layout.Regions["color:colors."+name] = Region{X: 1, Y: 3 + index, Width: dockWidth - 2, Height: 1}
	}
	ansiNames := []string{"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white"}
	for index, name := range ansiNames {
		x := 1 + index*5
		layout.Regions["color:palette.normal."+name] = Region{X: x, Y: 11, Width: 4, Height: 1}
		layout.Regions["color:palette.bright."+name] = Region{X: x, Y: 13, Width: 4, Height: 1}
	}

	pickerHeight := height - MinimumHeight + 6
	if pickerHeight > 10 {
		pickerHeight = 10
	}
	layout.Regions["sv"] = Region{X: 1, Y: 15, Width: dockWidth - 2, Height: pickerHeight}
	layout.Regions["hue"] = Region{X: 1, Y: 16 + pickerHeight, Width: dockWidth - 2, Height: 1}
	layout.Regions["hex"] = Region{X: 1, Y: 18 + pickerHeight, Width: 30, Height: 1}
	actionY := height - 3
	layout.Regions["action:import"] = Region{X: 1, Y: actionY, Width: 8, Height: 1}
	layout.Regions["action:save"] = Region{X: 11, Y: actionY, Width: 11, Height: 1}
	layout.Regions["action:export"] = Region{X: 24, Y: actionY, Width: 16, Height: 1}
	return layout
}
