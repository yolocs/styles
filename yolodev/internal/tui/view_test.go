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
		"Export", "TERMINAL PREVIEW", "selected text",
	} {
		if !strings.Contains(content, text) {
			t.Errorf("view missing %q", text)
		}
	}
	if strings.Contains(content, "Export Ghostty") {
		t.Error("main export action still names a single format")
	}
}

func TestViewPreviewDemonstratesSyntaxRoles(t *testing.T) {
	t.Parallel()

	model := updateModel(t, testModel(t), tea.WindowSizeMsg{Width: 120, Height: 36})
	content := model.View().Content
	tests := []struct {
		name  string
		text  string
		color string
	}{
		{name: "comment", text: "// Go · editor state", color: "#555555"},
		{name: "keyword", text: "func", color: "#FF55FF"},
		{name: "function", text: "render", color: "#00AAAA"},
		{name: "type", text: "Theme", color: "#5555FF"},
		{name: "number", text: "42", color: "#AAAA00"},
		{name: "string", text: `"signal-grid"`, color: "#00AA00"},
		{name: "invalid", text: "invalid:", color: "#FF5555"},
		{name: "deleted", text: "- stale cache", color: "#AA0000"},
		{name: "inserted", text: "+ palette synced", color: "#00AA00"},
		{name: "changed", text: "~ preview updated", color: "#AAAA00"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			want := lipgloss.NewStyle().Foreground(lipgloss.Color(test.color)).Render(test.text)
			if !strings.Contains(content, want) {
				t.Errorf("preview missing %s syntax sample %q in %s", test.name, test.text, test.color)
			}
		})
	}
	for _, text := range []string{"// JSON", `"accent"`, "true", "selected text", "cursor", "normal", "bright"} {
		if !strings.Contains(content, text) {
			t.Errorf("preview missing %q", text)
		}
	}
}

func TestExportPopupShowsFocusedControl(t *testing.T) {
	t.Parallel()

	model := updateModel(t, testModel(t), tea.WindowSizeMsg{Width: 120, Height: 36})
	model.openDialog(ExportDialog)
	if content := model.View().Content; !strings.Contains(content, "› Format:") {
		t.Fatalf("format-focused popup missing focus marker:\n%s", content)
	}
	model = updateModel(t, model, tea.KeyPressMsg{Code: tea.KeyTab})
	if content := model.View().Content; !strings.Contains(content, "› Path:") {
		t.Fatalf("path-focused popup missing focus marker:\n%s", content)
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
