package tui

import (
	"fmt"
	"image/color"
	"path/filepath"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/yolocs/styles/yolodev/internal/theme"
)

func renderMain(model Model) string {
	if model.Width < MinimumWidth || model.Height < MinimumHeight {
		return ""
	}
	value := theme.Normalize(model.Theme)
	background := lipgloss.Color(value.Colors.Background)
	foreground := lipgloss.Color(value.Colors.Foreground)
	muted := lipgloss.Color(value.Palette.Bright.Black)
	accent := lipgloss.Color(value.Colors.Cursor)
	base := lipgloss.NewStyle().Background(background).Foreground(foreground)

	state := "saved"
	if model.Dirty() {
		state = "modified"
	}
	warnings := 0
	for _, diagnostic := range theme.Validate(value) {
		if diagnostic.Severity == theme.Warning {
			warnings++
		}
	}
	warningLabel := fmt.Sprintf("%d warnings", warnings)
	if warnings == 1 {
		warningLabel = "1 warning"
	}
	headerText := fmt.Sprintf(" YOLODEV  %s • %s • %s", filepath.Base(model.Path), state, warningLabel)
	header := fitStyled(lipgloss.NewStyle().Bold(true).Background(lipgloss.Color(value.Palette.Normal.Black)).Foreground(accent).Render(headerText), model.Width)

	panelHeight := model.Height - 3
	dockWidth := 42
	previewWidth := model.Width - dockWidth - 1
	dock := make([]string, panelHeight)
	preview := make([]string, panelHeight)
	for index := range dock {
		dock[index] = ""
		preview[index] = ""
	}

	dock[0] = lipgloss.NewStyle().Bold(true).Foreground(accent).Render(" SEMANTIC")
	semantic := []ColorRef{
		{Group: "colors", Name: "background"},
		{Group: "colors", Name: "foreground"},
		{Group: "colors", Name: "cursor"},
		{Group: "colors", Name: "cursor_text"},
		{Group: "colors", Name: "selection_background"},
		{Group: "colors", Name: "selection_foreground"},
	}
	for index, ref := range semantic {
		dock[2+index] = renderColorRow(model, ref)
	}
	dock[9] = lipgloss.NewStyle().Bold(true).Foreground(accent).Render(" NORMAL")
	dock[10] = " " + renderPalette(value.Palette.Normal)
	dock[11] = lipgloss.NewStyle().Bold(true).Foreground(accent).Render(" BRIGHT")
	dock[12] = " " + renderPalette(value.Palette.Bright)
	dock[13] = lipgloss.NewStyle().Foreground(muted).Render(fmt.Sprintf(" %s  H %03.0f  S %03.0f%%  V %03.0f%%", model.Selected.Key(), model.Picker.HSV.H, model.Picker.HSV.S*100, model.Picker.HSV.V*100))

	layout := LayoutFor(model.Width, model.Height)
	sv := layout.Regions["sv"]
	for row := 0; row < sv.Height; row++ {
		var line strings.Builder
		line.WriteString(" ")
		for column := 0; column < sv.Width; column++ {
			saturation, valueLevel := svAt(sv, sv.X+column, sv.Y+row)
			color := theme.HSVToRGB(theme.HSV{H: model.Picker.HSV.H, S: saturation, V: valueLevel})
			cell := " "
			selectedX := int(model.Picker.HSV.S * float64(sv.Width-1))
			selectedY := int((1-model.Picker.HSV.V)*float64(sv.Height-1) + 0.5)
			if column == selectedX && row == selectedY {
				cell = "•"
			}
			line.WriteString(lipgloss.NewStyle().Background(lipgloss.Color(color.Hex())).Foreground(contrastColor(color)).Render(cell))
		}
		dock[sv.Y-1+row] = line.String()
	}
	hue := layout.Regions["hue"]
	var hueLine strings.Builder
	hueLine.WriteString(" ")
	for column := 0; column < hue.Width; column++ {
		color := theme.HSVToRGB(theme.HSV{H: hueAt(hue, hue.X+column), S: 1, V: 1})
		cell := " "
		selected := int(model.Picker.HSV.H / 359 * float64(hue.Width-1))
		if column == selected {
			cell = "◆"
		}
		hueLine.WriteString(lipgloss.NewStyle().Background(lipgloss.Color(color.Hex())).Foreground(contrastColor(color)).Render(cell))
	}
	dock[hue.Y-1] = hueLine.String()
	hex := " " + lipgloss.NewStyle().Bold(true).Foreground(accent).Render("HEX ")
	if model.HexFocused {
		hex += model.HexInput.View()
	} else {
		hex += model.color(model.Selected) + "  (click or Tab to edit)"
	}
	dock[layout.Regions["hex"].Y-1] = hex
	actionLine := " " + button("Import", accent, background) + "  " + button("Save TOML", accent, background) + "  " + button("Export Ghostty", accent, background)
	dock[layout.Regions["action:import"].Y-1] = actionLine

	preview[0] = lipgloss.NewStyle().Bold(true).Foreground(accent).Render(" TERMINAL PREVIEW")
	preview[2] = " " + colored(value.Palette.Normal.Green, "➜") + " " + colored(value.Palette.Normal.Blue, "~/styles") + " git status"
	preview[3] = " " + colored(value.Palette.Normal.Yellow, "M") + " themes/yolodev/placeholder.toml"
	preview[5] = " " + colored(value.Palette.Bright.Green, "✓ 42 tests passed")
	preview[6] = " " + colored(value.Palette.Normal.Yellow, "warning:") + " contrast is a design choice"
	preview[7] = " " + colored(value.Palette.Normal.Red, "error:") + " example failure text"
	preview[9] = " " + colored(value.Palette.Normal.Blue, "func") + " main() {"
	preview[10] = "     fmt.Println(" + colored(value.Palette.Normal.Green, `"hello, yolodev"`) + ")"
	preview[11] = " }"
	preview[13] = " " + lipgloss.NewStyle().Background(lipgloss.Color(value.Colors.SelectionBackground)).Foreground(lipgloss.Color(value.Colors.SelectionForeground)).Render(" selected text ")
	preview[15] = " cursor  " + lipgloss.NewStyle().Background(lipgloss.Color(value.Colors.Cursor)).Foreground(lipgloss.Color(value.Colors.CursorText)).Render(" A ")
	preview[17] = " normal  " + renderPalette(value.Palette.Normal)
	preview[19] = " bright  " + renderPalette(value.Palette.Bright)

	lines := []string{header}
	separator := lipgloss.NewStyle().Foreground(muted).Background(background).Render("│")
	for index := 0; index < panelHeight; index++ {
		left := base.Render(padStyled(dock[index], dockWidth))
		right := base.Render(padStyled(preview[index], previewWidth))
		lines = append(lines, left+separator+right)
	}
	footerText := " mouse: click/drag • Tab: hex • Ctrl+S: save • Ctrl+Q: quit"
	if model.Status != "" {
		footerText = " " + model.Status
	}
	lines = append(lines, fitStyled(lipgloss.NewStyle().Background(lipgloss.Color(value.Palette.Normal.Black)).Foreground(muted).Render(footerText), model.Width))
	return strings.Join(lines, "\n")
}

func renderDialog(model Model) string {
	title, message := dialogCopy(model)
	boxWidth := 70
	x := (model.Width - boxWidth) / 2
	y := model.Height/2 - 3
	background := lipgloss.Color(model.Theme.Colors.Background)
	foreground := lipgloss.Color(model.Theme.Colors.Foreground)
	accent := lipgloss.Color(model.Theme.Colors.Cursor)
	lines := make([]string, model.Height-1)
	for index := range lines {
		lines[index] = strings.Repeat(" ", model.Width)
	}
	put := func(row int, content string) {
		if row >= 0 && row < len(lines) {
			lines[row] = strings.Repeat(" ", x) + fitStyled(lipgloss.NewStyle().Background(background).Foreground(foreground).Render(content), boxWidth) + strings.Repeat(" ", model.Width-x-boxWidth)
		}
	}
	put(y, lipgloss.NewStyle().Bold(true).Foreground(accent).Render("  "+title))
	if model.Dialog.Kind == ImportDialog || model.Dialog.Kind == ExportDialog {
		put(y+2, "  "+model.Dialog.Input.View())
	} else {
		put(y+2, "  "+truncatePlain(message, boxWidth-4))
	}
	if model.Dialog.Message != "" && (model.Dialog.Kind == ImportDialog || model.Dialog.Kind == ExportDialog) {
		put(y+3, "  "+truncatePlain(model.Dialog.Message, boxWidth-4))
	}
	put(y+4, "  "+button("Confirm", accent, background)+"  "+button("Cancel", accent, background))
	return strings.Join(lines, "\n")
}

func dialogCopy(model Model) (string, string) {
	switch model.Dialog.Kind {
	case ImportDialog:
		return "IMPORT THEME", ""
	case ExportDialog:
		return "EXPORT GHOSTTY", ""
	case ConfirmImportDialog:
		return "DISCARD UNSAVED CHANGES?", "Import " + model.PendingPath + " and replace current edits?"
	case ConfirmExportDialog:
		return "REPLACE EXPORT?", model.PendingPath + " already exists. Replace it?"
	case ConfirmQuitDialog:
		return "QUIT WITHOUT SAVING?", "Current theme has unsaved changes."
	case ErrorDialog:
		return "ERROR", model.Dialog.Message
	default:
		return "YOLODEV", ""
	}
}

func renderColorRow(model Model, ref ColorRef) string {
	marker := " "
	if model.Selected == ref {
		marker = "›"
	}
	color, _ := theme.ParseHex(model.color(ref))
	swatch := lipgloss.NewStyle().Background(lipgloss.Color(color.Hex())).Render("  ")
	return fmt.Sprintf(" %s %s %-22s %s", marker, swatch, ref.Name, color.Hex())
}

func renderPalette(value theme.ANSIColors) string {
	var output strings.Builder
	for _, color := range value.Values() {
		output.WriteString(lipgloss.NewStyle().Background(lipgloss.Color(color)).Render("    "))
		output.WriteString(" ")
	}
	return strings.TrimSuffix(output.String(), " ")
}

func colored(color, text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(text)
}

func contrastColor(value theme.RGB) color.Color {
	if theme.ContrastRatio(value, theme.RGB{}) > theme.ContrastRatio(value, theme.RGB{R: 255, G: 255, B: 255}) {
		return lipgloss.Color("#000000")
	}
	return lipgloss.Color("#FFFFFF")
}

func button(label string, foreground, background color.Color) string {
	return lipgloss.NewStyle().Bold(true).Foreground(background).Background(foreground).Render("[" + label + "]")
}

func padStyled(value string, width int) string {
	missing := width - lipgloss.Width(value)
	if missing <= 0 {
		return value
	}
	return value + strings.Repeat(" ", missing)
}

func fitStyled(value string, width int) string {
	if lipgloss.Width(value) > width {
		return lipgloss.NewStyle().MaxWidth(width).Render(value)
	}
	return padStyled(value, width)
}

func truncatePlain(value string, width int) string {
	runes := []rune(value)
	if len(runes) <= width {
		return value
	}
	if width <= 1 {
		return string(runes[:width])
	}
	return string(runes[:width-1]) + "…"
}
