package tui

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/yolocs/styles/yolodev/internal/theme"
)

func TestImportActionOpensPathDialog(t *testing.T) {
	t.Parallel()

	model := testModel(t)
	model.Regions = map[string]Region{"action:import": {X: 2, Y: 2, Width: 10, Height: 1}}
	updated := updateModel(t, model, tea.MouseClickMsg{X: 3, Y: 2, Button: tea.MouseLeft})
	if updated.Dialog.Kind != ImportDialog || !updated.Dialog.Input.Focused() {
		t.Fatalf("dialog = %#v, want focused import dialog", updated.Dialog)
	}
}

func TestExportActionOpensFormatPopup(t *testing.T) {
	t.Parallel()

	model := updateModel(t, testModel(t), tea.WindowSizeMsg{Width: 120, Height: 36})
	region := model.Regions["action:export"]
	updated := updateModel(t, model, tea.MouseClickMsg{X: region.X, Y: region.Y, Button: tea.MouseLeft})
	if updated.Dialog.Kind != ExportDialog {
		t.Fatalf("dialog kind = %v, want export dialog", updated.Dialog.Kind)
	}
	if updated.Dialog.Input.Focused() {
		t.Fatal("export path is focused before a format is chosen")
	}
	content := updated.View().Content
	for _, text := range []string{"EXPORT THEME", "Ghostty", "Codex", "←/→ choose"} {
		if !strings.Contains(content, text) {
			t.Errorf("export popup missing %q", text)
		}
	}
}

func TestExportPopupWritesSelectedCodexFormat(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "signal-grid.tmTheme")
	model := updateModel(t, testModel(t), tea.WindowSizeMsg{Width: 120, Height: 36})
	action := model.Regions["action:export"]
	model = updateModel(t, model, tea.MouseClickMsg{X: action.X, Y: action.Y, Button: tea.MouseLeft})
	codex := model.Regions["dialog:format:codex"]
	model = updateModel(t, model, tea.MouseClickMsg{X: codex.X, Y: codex.Y, Button: tea.MouseLeft})
	input := model.Regions["dialog:input"]
	model = updateModel(t, model, tea.MouseClickMsg{X: input.X, Y: input.Y, Button: tea.MouseLeft})
	model.Dialog.Input.SetValue(path)

	updated, command := updateModelCommand(t, model, tea.KeyPressMsg{Code: tea.KeyEnter})
	if command == nil {
		t.Fatal("Codex export returned nil command")
	}
	updated = updateModel(t, updated, command())
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(data, []byte("<?xml")) {
		t.Fatalf("selected Codex export wrote unexpected data:\n%s", data)
	}
	if updated.Dialog.Kind != NoDialog {
		t.Fatalf("successful Codex export left dialog %v open", updated.Dialog.Kind)
	}
}

func TestExportPopupArrowKeysSelectFormat(t *testing.T) {
	t.Parallel()

	model := testModel(t)
	model.openDialog(ExportDialog)
	model = updateModel(t, model, tea.KeyPressMsg{Code: tea.KeyRight})
	if model.Dialog.ExportFormat != "codex" {
		t.Fatalf("right arrow selected %q, want codex", model.Dialog.ExportFormat)
	}
	model = updateModel(t, model, tea.KeyPressMsg{Code: tea.KeyLeft})
	if model.Dialog.ExportFormat != "ghostty" {
		t.Fatalf("left arrow selected %q, want ghostty", model.Dialog.ExportFormat)
	}
}

func TestExportPopupTabSwitchesFocus(t *testing.T) {
	t.Parallel()

	model := testModel(t)
	model.openDialog(ExportDialog)
	model = updateModel(t, model, tea.KeyPressMsg{Code: tea.KeyTab})
	if !model.Dialog.Input.Focused() {
		t.Fatal("Tab did not focus the export path")
	}
	model = updateModel(t, model, tea.KeyPressMsg{Code: tea.KeyTab})
	if model.Dialog.Input.Focused() {
		t.Fatal("second Tab did not return focus to the format selector")
	}
}

func TestExportPopupEnterMovesFromFormatToPath(t *testing.T) {
	t.Parallel()

	model := testModel(t)
	model.openDialog(ExportDialog)
	updated, command := updateModelCommand(t, model, tea.KeyPressMsg{Code: tea.KeyEnter})
	if command != nil || !updated.Dialog.Input.Focused() || updated.Dialog.Message != "" {
		t.Fatalf("Enter on format returned command=%v focused=%v message=%q", command, updated.Dialog.Input.Focused(), updated.Dialog.Message)
	}
}

func TestPathDialogAcceptsConfirmationShortcutLettersAsText(t *testing.T) {
	t.Parallel()

	model := testModel(t)
	model.openDialog(ImportDialog)
	updated := updateModel(t, model, tea.KeyPressMsg{Code: 'y', Text: "y"})
	if updated.Dialog.Kind != ImportDialog || updated.Dialog.Input.Value() != "y" {
		t.Fatalf("typing y produced dialog=%v value=%q", updated.Dialog.Kind, updated.Dialog.Input.Value())
	}
	updated = updateModel(t, updated, tea.KeyPressMsg{Code: 'n', Text: "n"})
	if updated.Dialog.Kind != ImportDialog || updated.Dialog.Input.Value() != "yn" {
		t.Fatalf("typing n produced dialog=%v value=%q", updated.Dialog.Kind, updated.Dialog.Input.Value())
	}
}

func TestInvalidImportLeavesCurrentThemeUntouched(t *testing.T) {
	t.Parallel()

	model := testModel(t)
	original := model.Theme
	path := filepath.Join(t.TempDir(), "invalid.toml")
	if err := os.WriteFile(path, []byte("not = [valid"), 0o600); err != nil {
		t.Fatal(err)
	}
	model.openDialog(ImportDialog)
	model.Dialog.Input.SetValue(path)

	updated, command := updateModelCommand(t, model, tea.KeyPressMsg{Code: tea.KeyEnter})
	if command == nil {
		t.Fatal("import Enter returned nil command")
	}
	updated = updateModel(t, updated, command())
	if updated.Theme != original || updated.Dialog.Kind != ErrorDialog {
		t.Fatalf("invalid import mutated theme or missed error: dialog=%#v", updated.Dialog)
	}
}

func TestDirtyImportRequiresConfirmation(t *testing.T) {
	t.Parallel()

	model := testModel(t)
	model.Theme.Colors.Foreground = "#FFFFFF"
	model.openDialog(ImportDialog)
	model.Dialog.Input.SetValue("replacement.toml")
	updated, command := updateModelCommand(t, model, tea.KeyPressMsg{Code: tea.KeyEnter})
	if command != nil || updated.Dialog.Kind != ConfirmImportDialog || updated.PendingPath != "replacement.toml" {
		t.Fatalf("dirty import = dialog %v path %q cmd %v", updated.Dialog.Kind, updated.PendingPath, command)
	}
}

func TestSaveWritesNormalizedThemeAndClearsDirty(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "theme.toml")
	model := testModel(t)
	model.Path = path
	model.Theme.Colors.Foreground = "#abcdef"
	updated, command := updateModelCommand(t, model, tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})
	if command == nil {
		t.Fatal("Ctrl+S returned nil command")
	}
	updated = updateModel(t, updated, command())
	if updated.Dirty() || updated.Theme.Colors.Foreground != "#ABCDEF" {
		t.Fatalf("saved model dirty=%v foreground=%s", updated.Dirty(), updated.Theme.Colors.Foreground)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte("#ABCDEF")) {
		t.Fatalf("saved file was not normalized:\n%s", data)
	}
}

func TestExportExistingFileRequiresConfirmation(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "ghostty")
	if err := os.WriteFile(path, []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}
	model := testModel(t)
	model.openDialog(ExportDialog)
	model = updateModel(t, model, tea.KeyPressMsg{Code: tea.KeyTab})
	model.Dialog.Input.SetValue(path)
	updated, command := updateModelCommand(t, model, tea.KeyPressMsg{Code: tea.KeyEnter})
	updated = updateModel(t, updated, command())
	if updated.Dialog.Kind != ConfirmExportDialog || updated.PendingPath != path {
		t.Fatalf("export existing = dialog %v path %q", updated.Dialog.Kind, updated.PendingPath)
	}

	updated, command = updateModelCommand(t, updated, tea.KeyPressMsg{Code: tea.KeyEnter})
	updated = updateModel(t, updated, command())
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "palette = 0=") || updated.Dialog.Kind != NoDialog {
		t.Fatalf("confirmed export = dialog %v data %q", updated.Dialog.Kind, data)
	}
}

func TestCodexFormatSurvivesOverwriteConfirmation(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "existing.tmTheme")
	if err := os.WriteFile(path, []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}
	model := testModel(t)
	model.openDialog(ExportDialog)
	model = updateModel(t, model, tea.KeyPressMsg{Code: tea.KeyRight})
	model = updateModel(t, model, tea.KeyPressMsg{Code: tea.KeyTab})
	model.Dialog.Input.SetValue(path)

	updated, command := updateModelCommand(t, model, tea.KeyPressMsg{Code: tea.KeyEnter})
	if command == nil {
		t.Fatal("initial Codex export returned nil command")
	}
	updated = updateModel(t, updated, command())
	if updated.Dialog.Kind != ConfirmExportDialog {
		t.Fatalf("existing Codex export opened dialog %v, want confirmation", updated.Dialog.Kind)
	}
	updated, command = updateModelCommand(t, updated, tea.KeyPressMsg{Code: tea.KeyEnter})
	if command == nil {
		t.Fatal("confirmed Codex overwrite returned nil command")
	}
	updated = updateModel(t, updated, command())
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(data, []byte("<?xml")) || updated.Dialog.Kind != NoDialog {
		t.Fatalf("confirmed Codex export = dialog %v data %q", updated.Dialog.Kind, data)
	}
}

func TestExportFolderConfirmationUsesGeneratedFilename(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	outputPath := filepath.Join(directory, "test.tmTheme")
	if err := os.WriteFile(outputPath, []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}
	model := testModel(t)
	message, ok := exportThemeCmd("codex", directory, model.Theme, false)().(fileExportedMsg)
	if !ok {
		t.Fatal("exportThemeCmd() did not return fileExportedMsg")
	}
	updated := model.handleFileExported(message)
	if updated.Dialog.Kind != ConfirmExportDialog || updated.PendingPath != outputPath {
		t.Fatalf("export existing = dialog %v path %q, want path %q", updated.Dialog.Kind, updated.PendingPath, outputPath)
	}
}

func TestDirtyQuitRequiresConfirmation(t *testing.T) {
	t.Parallel()

	model := testModel(t)
	model.Theme.Colors.Foreground = "#FFFFFF"
	updated, command := updateModelCommand(t, model, tea.KeyPressMsg{Code: 'q', Mod: tea.ModCtrl})
	if command != nil || updated.Dialog.Kind != ConfirmQuitDialog {
		t.Fatalf("dirty quit = dialog %v cmd %v", updated.Dialog.Kind, command)
	}
}

func TestHexEntryUpdatesSelectedColor(t *testing.T) {
	t.Parallel()

	model := testModel(t)
	model.HexFocused = true
	model.HexInput.Focus()
	model.HexInput.SetValue("#FFFFFF")
	updated := updateModel(t, model, tea.KeyPressMsg{Code: tea.KeyEnter})
	if got := updated.Theme.Colors.Background; got != "#FFFFFF" || updated.HexFocused {
		t.Fatalf("hex commit = %s focused=%v", got, updated.HexFocused)
	}
}

func updateModelCommand(t *testing.T, model Model, message tea.Msg) (Model, tea.Cmd) {
	t.Helper()
	updated, command := model.Update(message)
	got, ok := updated.(Model)
	if !ok {
		t.Fatalf("Update returned %T, want Model", updated)
	}
	return got, command
}

func validEncodedTheme(t *testing.T, value theme.Theme) []byte {
	t.Helper()
	var output bytes.Buffer
	if err := theme.Encode(&output, value); err != nil {
		t.Fatal(err)
	}
	return output.Bytes()
}
