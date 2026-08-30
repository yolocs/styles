package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func TestViewContainsPreviewFirstControlsAndRoles(t *testing.T) {
	t.Parallel()

	model := updateModel(t, testModel(t), tea.WindowSizeMsg{Width: 120, Height: 36})
	content := model.View().Content
	for _, text := range []string{
		"YOLODEV", "SEMANTIC", "NORMAL", "BRIGHT", "HEX", "Import", "Save TOML",
		"Export Ghostty", "TERMINAL PREVIEW", "warning:", "error:", "selected text",
	} {
		if !strings.Contains(content, text) {
			t.Errorf("view missing %q", text)
		}
	}
}

func TestViewNeverExceedsTerminalWidth(t *testing.T) {
	t.Parallel()

	model := updateModel(t, testModel(t), tea.WindowSizeMsg{Width: 100, Height: 28})
	for index, line := range strings.Split(model.View().Content, "\n") {
		if width := lipgloss.Width(line); width > 100 {
			t.Errorf("line %d width = %d, want <= 100", index, width)
		}
	}
}

func TestViewShowsDirtyMarkerAndDiagnostics(t *testing.T) {
	t.Parallel()

	model := updateModel(t, testModel(t), tea.WindowSizeMsg{Width: 120, Height: 36})
	model.Theme.Colors.Foreground = model.Theme.Colors.Background
	content := model.View().Content
	if !strings.Contains(content, "modified") || !strings.Contains(content, "1 warning") {
		t.Fatalf("dirty diagnostic view missing state:\n%s", content)
	}
}

func TestSmallViewReportsRequiredSize(t *testing.T) {
	t.Parallel()

	model := updateModel(t, testModel(t), tea.WindowSizeMsg{Width: 80, Height: 20})
	if got := model.View().Content; !strings.Contains(got, "requires 100x28") {
		t.Fatalf("small view = %q", got)
	}
}
