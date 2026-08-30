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
	model := testModel(t)
	model.Regions = map[string]Region{"action:import": {X: 2, Y: 2, Width: 10, Height: 1}}
	updated := updateModel(t, model, tea.MouseClickMsg{X: 3, Y: 2, Button: tea.MouseLeft})
	if updated.Dialog.Kind != ImportDialog || !updated.Dialog.Input.Focused() {
		t.Fatalf("dialog = %#v, want focused import dialog", updated.Dialog)
	}
}

func TestPathDialogAcceptsConfirmationShortcutLettersAsText(t *testing.T) {
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
	path := filepath.Join(t.TempDir(), "ghostty")
	if err := os.WriteFile(path, []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}
	model := testModel(t)
	model.openDialog(ExportDialog)
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

func TestDirtyQuitRequiresConfirmation(t *testing.T) {
	model := testModel(t)
	model.Theme.Colors.Foreground = "#FFFFFF"
	updated, command := updateModelCommand(t, model, tea.KeyPressMsg{Code: 'q', Mod: tea.ModCtrl})
	if command != nil || updated.Dialog.Kind != ConfirmQuitDialog {
		t.Fatalf("dirty quit = dialog %v cmd %v", updated.Dialog.Kind, command)
	}
}

func TestHexEntryUpdatesSelectedColor(t *testing.T) {
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
