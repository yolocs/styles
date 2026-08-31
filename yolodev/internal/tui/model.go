package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/yolocs/styles/yolodev/internal/theme"
)

type ColorRef struct {
	Group string
	Name  string
}

func (r ColorRef) Key() string {
	return r.Group + "." + r.Name
}

type Model struct {
	Theme         theme.Theme
	Clean         theme.Theme
	Path          string
	Width         int
	Height        int
	Layout        Layout
	Regions       map[string]Region
	Selected      ColorRef
	Picker        Picker
	HexInput      textinput.Model
	HexFocused    bool
	Dialog        Dialog
	PendingPath   string
	PendingFormat string
	Status        string
}

func New(value theme.Theme, path string) Model {
	value = theme.Normalize(value)
	hexInput := textinput.New()
	hexInput.Prompt = ""
	hexInput.CharLimit = 7
	hexInput.SetWidth(9)
	hexInput.SetValue(value.Colors.Background)
	model := Model{
		Theme:    value,
		Clean:    value,
		Path:     path,
		Regions:  make(map[string]Region),
		Selected: ColorRef{Group: "colors", Name: "background"},
		HexInput: hexInput,
	}
	model.syncPicker()
	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		m.Width = message.Width
		m.Height = message.Height
		m.refreshRegions()
		return m, nil
	case fileLoadedMsg:
		return m.handleFileLoaded(message), nil
	case fileSavedMsg:
		return m.handleFileSaved(message), nil
	case fileExportedMsg:
		return m.handleFileExported(message), nil
	}

	if m.Dialog.Kind != NoDialog {
		return m.updateDialog(message)
	}
	if m.HexFocused {
		return m.updateHex(message)
	}

	switch message := message.(type) {
	case tea.MouseClickMsg:
		if message.Button == tea.MouseLeft {
			return m, m.handleClick(message.X, message.Y)
		}
	case tea.MouseMotionMsg:
		if message.Button == tea.MouseLeft {
			m.handleDrag(message.X, message.Y)
		}
	case tea.MouseReleaseMsg:
		if message.Button == tea.MouseLeft {
			m.Picker.Drag = DragNone
			m.Picker.StartHex = ""
		}
	case tea.KeyPressMsg:
		switch message.Keystroke() {
		case "ctrl+s":
			return m, saveThemeCmd(m.Path, m.Theme)
		case "ctrl+q":
			if m.Dirty() {
				m.openDialog(ConfirmQuitDialog)
				return m, nil
			}
			return m, tea.Quit
		case "tab":
			m.beginHexEdit()
			return m, nil
		}
		if message.Code == tea.KeyEscape && m.Picker.Drag != DragNone {
			m.setColor(m.Picker.StartRef, m.Picker.StartHex)
			m.Selected = m.Picker.StartRef
			m.Picker.Drag = DragNone
			m.Picker.StartHex = ""
			m.syncPicker()
		}
		if message.Code == tea.KeyLeft || message.Code == tea.KeyRight || message.Code == tea.KeyUp || message.Code == tea.KeyDown {
			m.adjustPicker(message.Code)
		}
	}
	return m, nil
}

func (m Model) View() tea.View {
	content := renderMain(m)
	if m.Width < MinimumWidth || m.Height < MinimumHeight {
		content = fmt.Sprintf("yolodev requires %dx%d; current terminal is %dx%d\nCtrl+Q to quit", MinimumWidth, MinimumHeight, m.Width, m.Height)
	} else if m.Dialog.Kind != NoDialog {
		content = renderDialog(m)
	}
	view := tea.NewView(content)
	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion
	return view
}

func Run(value theme.Theme, path string) error {
	program := tea.NewProgram(New(value, path))
	if _, err := program.Run(); err != nil {
		return fmt.Errorf("run terminal-theme editor: %w", err)
	}
	return nil
}

func (m Model) Dirty() bool {
	return theme.Normalize(m.Theme) != theme.Normalize(m.Clean)
}

func (m *Model) handleClick(x, y int) tea.Cmd {
	if region, ok := m.Regions["action:import"]; ok && region.Contains(x, y) {
		m.openDialog(ImportDialog)
		return nil
	}
	if region, ok := m.Regions["action:save"]; ok && region.Contains(x, y) {
		return saveThemeCmd(m.Path, m.Theme)
	}
	if region, ok := m.Regions["action:export"]; ok && region.Contains(x, y) {
		m.openDialog(ExportDialog)
		return nil
	}
	if region, ok := m.Regions["hex"]; ok && region.Contains(x, y) {
		m.beginHexEdit()
		return nil
	}
	for key, region := range m.Regions {
		if strings.HasPrefix(key, "color:") && region.Contains(x, y) {
			m.Selected = parseColorRef(strings.TrimPrefix(key, "color:"))
			m.syncPicker()
			return nil
		}
	}
	if region, ok := m.Regions["sv"]; ok && region.Contains(x, y) {
		m.beginDrag(DragSV)
		m.applySV(region, x, y)
		return nil
	}
	if region, ok := m.Regions["hue"]; ok && region.Contains(x, y) {
		m.beginDrag(DragHue)
		m.applyHue(region, x)
	}
	return nil
}

func (m *Model) handleDrag(x, y int) {
	switch m.Picker.Drag {
	case DragSV:
		m.applySV(m.Regions["sv"], x, y)
	case DragHue:
		m.applyHue(m.Regions["hue"], x)
	}
}

func (m *Model) beginDrag(target DragTarget) {
	m.Picker.Drag = target
	m.Picker.StartHex = m.color(m.Selected)
	m.Picker.StartRef = m.Selected
}

func (m *Model) applySV(region Region, x, y int) {
	m.Picker.HSV.S, m.Picker.HSV.V = svAt(region, x, y)
	m.setColor(m.Selected, theme.HSVToRGB(m.Picker.HSV).Hex())
}

func (m *Model) applyHue(region Region, x int) {
	m.Picker.HSV.H = hueAt(region, x)
	m.setColor(m.Selected, theme.HSVToRGB(m.Picker.HSV).Hex())
}

func (m *Model) syncPicker() {
	color, err := theme.ParseHex(m.color(m.Selected))
	if err != nil {
		return
	}
	m.Picker.HSV = theme.RGBToHSV(color)
	if !m.HexFocused {
		m.HexInput.SetValue(color.Hex())
	}
}

func (m *Model) refreshRegions() {
	m.Layout = LayoutFor(m.Width, m.Height)
	m.Regions = make(map[string]Region, len(m.Layout.Regions)+3)
	for key, region := range m.Layout.Regions {
		m.Regions[key] = region
	}
	if m.Dialog.Kind != NoDialog {
		m.installDialogRegions()
	}
}

func (m *Model) beginHexEdit() {
	m.HexInput.SetValue(m.color(m.Selected))
	m.HexInput.Focus()
	m.HexFocused = true
}

func (m Model) updateHex(message tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := message.(tea.KeyPressMsg); ok {
		switch key.Code {
		case tea.KeyEnter:
			color, err := theme.ParseHex(m.HexInput.Value())
			if err != nil {
				m.Status = err.Error()
				return m, nil
			}
			m.setColor(m.Selected, color.Hex())
			m.HexInput.Blur()
			m.HexFocused = false
			m.Status = "color updated"
			m.syncPicker()
			return m, nil
		case tea.KeyEscape:
			m.HexInput.SetValue(m.color(m.Selected))
			m.HexInput.Blur()
			m.HexFocused = false
			return m, nil
		}
	}
	var command tea.Cmd
	m.HexInput, command = m.HexInput.Update(message)
	return m, command
}

func (m *Model) adjustPicker(code rune) {
	switch code {
	case tea.KeyLeft:
		m.Picker.HSV.H -= 1
	case tea.KeyRight:
		m.Picker.HSV.H += 1
	case tea.KeyUp:
		m.Picker.HSV.V += 0.01
	case tea.KeyDown:
		m.Picker.HSV.V -= 0.01
	}
	m.setColor(m.Selected, theme.HSVToRGB(m.Picker.HSV).Hex())
	m.syncPicker()
}

func parseColorRef(value string) ColorRef {
	index := strings.LastIndex(value, ".")
	if index < 0 {
		return ColorRef{}
	}
	return ColorRef{Group: value[:index], Name: value[index+1:]}
}

func (m Model) color(ref ColorRef) string {
	switch ref.Group + "." + ref.Name {
	case "colors.background":
		return m.Theme.Colors.Background
	case "colors.foreground":
		return m.Theme.Colors.Foreground
	case "colors.cursor":
		return m.Theme.Colors.Cursor
	case "colors.cursor_text":
		return m.Theme.Colors.CursorText
	case "colors.selection_background":
		return m.Theme.Colors.SelectionBackground
	case "colors.selection_foreground":
		return m.Theme.Colors.SelectionForeground
	case "palette.normal.black":
		return m.Theme.Palette.Normal.Black
	case "palette.normal.red":
		return m.Theme.Palette.Normal.Red
	case "palette.normal.green":
		return m.Theme.Palette.Normal.Green
	case "palette.normal.yellow":
		return m.Theme.Palette.Normal.Yellow
	case "palette.normal.blue":
		return m.Theme.Palette.Normal.Blue
	case "palette.normal.magenta":
		return m.Theme.Palette.Normal.Magenta
	case "palette.normal.cyan":
		return m.Theme.Palette.Normal.Cyan
	case "palette.normal.white":
		return m.Theme.Palette.Normal.White
	case "palette.bright.black":
		return m.Theme.Palette.Bright.Black
	case "palette.bright.red":
		return m.Theme.Palette.Bright.Red
	case "palette.bright.green":
		return m.Theme.Palette.Bright.Green
	case "palette.bright.yellow":
		return m.Theme.Palette.Bright.Yellow
	case "palette.bright.blue":
		return m.Theme.Palette.Bright.Blue
	case "palette.bright.magenta":
		return m.Theme.Palette.Bright.Magenta
	case "palette.bright.cyan":
		return m.Theme.Palette.Bright.Cyan
	case "palette.bright.white":
		return m.Theme.Palette.Bright.White
	default:
		return ""
	}
}

func (m *Model) setColor(ref ColorRef, value string) {
	switch ref.Group + "." + ref.Name {
	case "colors.background":
		m.Theme.Colors.Background = value
	case "colors.foreground":
		m.Theme.Colors.Foreground = value
	case "colors.cursor":
		m.Theme.Colors.Cursor = value
	case "colors.cursor_text":
		m.Theme.Colors.CursorText = value
	case "colors.selection_background":
		m.Theme.Colors.SelectionBackground = value
	case "colors.selection_foreground":
		m.Theme.Colors.SelectionForeground = value
	case "palette.normal.black":
		m.Theme.Palette.Normal.Black = value
	case "palette.normal.red":
		m.Theme.Palette.Normal.Red = value
	case "palette.normal.green":
		m.Theme.Palette.Normal.Green = value
	case "palette.normal.yellow":
		m.Theme.Palette.Normal.Yellow = value
	case "palette.normal.blue":
		m.Theme.Palette.Normal.Blue = value
	case "palette.normal.magenta":
		m.Theme.Palette.Normal.Magenta = value
	case "palette.normal.cyan":
		m.Theme.Palette.Normal.Cyan = value
	case "palette.normal.white":
		m.Theme.Palette.Normal.White = value
	case "palette.bright.black":
		m.Theme.Palette.Bright.Black = value
	case "palette.bright.red":
		m.Theme.Palette.Bright.Red = value
	case "palette.bright.green":
		m.Theme.Palette.Bright.Green = value
	case "palette.bright.yellow":
		m.Theme.Palette.Bright.Yellow = value
	case "palette.bright.blue":
		m.Theme.Palette.Bright.Blue = value
	case "palette.bright.magenta":
		m.Theme.Palette.Bright.Magenta = value
	case "palette.bright.cyan":
		m.Theme.Palette.Bright.Cyan = value
	case "palette.bright.white":
		m.Theme.Palette.Bright.White = value
	}
}
