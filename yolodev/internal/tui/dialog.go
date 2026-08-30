package tui

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/yolocs/styles/yolodev/internal/fileutil"
	"github.com/yolocs/styles/yolodev/internal/ghostty"
	"github.com/yolocs/styles/yolodev/internal/theme"
)

type DialogKind int

const (
	NoDialog DialogKind = iota
	ImportDialog
	ExportDialog
	ConfirmImportDialog
	ConfirmExportDialog
	ConfirmQuitDialog
	ErrorDialog
)

type Dialog struct {
	Kind    DialogKind
	Input   textinput.Model
	Message string
}

type fileLoadedMsg struct {
	Path  string
	Theme theme.Theme
	Err   error
}

type fileSavedMsg struct {
	Path  string
	Theme theme.Theme
	Err   error
}

type fileExportedMsg struct {
	Path string
	Err  error
}

func (m *Model) openDialog(kind DialogKind) {
	input := textinput.New()
	input.Prompt = "Path: "
	input.SetWidth(58)
	if kind == ImportDialog || kind == ExportDialog {
		input.Focus()
	}
	m.Dialog = Dialog{Kind: kind, Input: input}
	m.refreshRegions()
}

func (m *Model) openError(err error) {
	m.openDialog(ErrorDialog)
	m.Dialog.Message = err.Error()
}

func (m *Model) closeDialog() {
	m.Dialog = Dialog{}
	m.refreshRegions()
}

func (m *Model) installDialogRegions() {
	if m.Width < MinimumWidth || m.Height < MinimumHeight {
		return
	}
	boxWidth := 70
	x := (m.Width - boxWidth) / 2
	y := m.Height/2 - 3
	m.Regions["dialog:input"] = Region{X: x + 2, Y: y + 2, Width: boxWidth - 4, Height: 1}
	m.Regions["dialog:confirm"] = Region{X: x + 2, Y: y + 4, Width: 9, Height: 1}
	m.Regions["dialog:cancel"] = Region{X: x + 13, Y: y + 4, Width: 8, Height: 1}
}

func (m Model) updateDialog(message tea.Msg) (tea.Model, tea.Cmd) {
	if mouse, ok := message.(tea.MouseClickMsg); ok && mouse.Button == tea.MouseLeft {
		if region := m.Regions["dialog:cancel"]; region.Contains(mouse.X, mouse.Y) {
			m.closeDialog()
			return m, nil
		}
		if region := m.Regions["dialog:confirm"]; region.Contains(mouse.X, mouse.Y) {
			return m.confirmDialog()
		}
		if region := m.Regions["dialog:input"]; region.Contains(mouse.X, mouse.Y) {
			m.Dialog.Input.Focus()
		}
	}

	if key, ok := message.(tea.KeyPressMsg); ok {
		confirmation := m.Dialog.Kind == ConfirmImportDialog || m.Dialog.Kind == ConfirmExportDialog || m.Dialog.Kind == ConfirmQuitDialog || m.Dialog.Kind == ErrorDialog
		if key.Code == tea.KeyEscape || (confirmation && key.Keystroke() == "n") {
			m.closeDialog()
			return m, nil
		}
		if key.Code == tea.KeyEnter || (confirmation && key.Keystroke() == "y") {
			return m.confirmDialog()
		}
	}

	if m.Dialog.Kind == ImportDialog || m.Dialog.Kind == ExportDialog {
		var command tea.Cmd
		m.Dialog.Input, command = m.Dialog.Input.Update(message)
		return m, command
	}
	return m, nil
}

func (m Model) confirmDialog() (tea.Model, tea.Cmd) {
	switch m.Dialog.Kind {
	case ImportDialog:
		path := strings.TrimSpace(m.Dialog.Input.Value())
		if path == "" {
			m.Dialog.Message = "path must not be empty"
			return m, nil
		}
		if m.Dirty() {
			m.PendingPath = path
			m.openDialog(ConfirmImportDialog)
			return m, nil
		}
		m.closeDialog()
		return m, loadThemeCmd(path)
	case ExportDialog:
		path := strings.TrimSpace(m.Dialog.Input.Value())
		if path == "" {
			m.Dialog.Message = "path must not be empty"
			return m, nil
		}
		m.closeDialog()
		return m, exportThemeCmd(path, m.Theme, false)
	case ConfirmImportDialog:
		path := m.PendingPath
		m.closeDialog()
		return m, loadThemeCmd(path)
	case ConfirmExportDialog:
		path := m.PendingPath
		m.closeDialog()
		return m, exportThemeCmd(path, m.Theme, true)
	case ConfirmQuitDialog:
		m.closeDialog()
		return m, tea.Quit
	case ErrorDialog:
		m.closeDialog()
	}
	return m, nil
}

func loadThemeCmd(path string) tea.Cmd {
	return func() tea.Msg {
		file, err := os.Open(path)
		if err != nil {
			return fileLoadedMsg{Path: path, Err: fmt.Errorf("open %s: %w", path, err)}
		}
		value, decodeErr := theme.Decode(file)
		closeErr := file.Close()
		if decodeErr != nil {
			return fileLoadedMsg{Path: path, Err: decodeErr}
		}
		if closeErr != nil {
			return fileLoadedMsg{Path: path, Err: fmt.Errorf("close %s: %w", path, closeErr)}
		}
		diagnostics := theme.Validate(value)
		if theme.HasErrors(diagnostics) {
			return fileLoadedMsg{Path: path, Err: diagnosticsError(diagnostics)}
		}
		return fileLoadedMsg{Path: path, Theme: theme.Normalize(value)}
	}
}

func saveThemeCmd(path string, value theme.Theme) tea.Cmd {
	return func() tea.Msg {
		if strings.TrimSpace(path) == "" {
			return fileSavedMsg{Path: path, Err: errors.New("cannot save: theme path is empty")}
		}
		value = theme.Normalize(value)
		diagnostics := theme.Validate(value)
		if theme.HasErrors(diagnostics) {
			return fileSavedMsg{Path: path, Err: diagnosticsError(diagnostics)}
		}
		var output bytes.Buffer
		if err := theme.Encode(&output, value); err != nil {
			return fileSavedMsg{Path: path, Err: err}
		}
		if err := fileutil.WriteAtomic(path, output.Bytes(), true); err != nil {
			return fileSavedMsg{Path: path, Err: err}
		}
		return fileSavedMsg{Path: path, Theme: value}
	}
}

func exportThemeCmd(path string, value theme.Theme, overwrite bool) tea.Cmd {
	return func() tea.Msg {
		var output bytes.Buffer
		if err := ghostty.Render(&output, value); err != nil {
			return fileExportedMsg{Path: path, Err: err}
		}
		if err := fileutil.WriteAtomic(path, output.Bytes(), overwrite); err != nil {
			return fileExportedMsg{Path: path, Err: err}
		}
		return fileExportedMsg{Path: path}
	}
}

func (m Model) handleFileLoaded(message fileLoadedMsg) Model {
	if message.Err != nil {
		m.openError(message.Err)
		return m
	}
	m.Theme = message.Theme
	m.Clean = message.Theme
	m.Path = message.Path
	m.Selected = ColorRef{Group: "colors", Name: "background"}
	m.PendingPath = ""
	m.closeDialog()
	m.syncPicker()
	m.Status = "imported " + message.Path
	return m
}

func (m Model) handleFileSaved(message fileSavedMsg) Model {
	if message.Err != nil {
		m.openError(message.Err)
		return m
	}
	m.Theme = message.Theme
	m.Clean = message.Theme
	m.Path = message.Path
	m.Status = "saved " + message.Path
	m.syncPicker()
	return m
}

func (m Model) handleFileExported(message fileExportedMsg) Model {
	if errors.Is(message.Err, fileutil.ErrExists) {
		m.PendingPath = message.Path
		m.openDialog(ConfirmExportDialog)
		return m
	}
	if message.Err != nil {
		m.openError(message.Err)
		return m
	}
	m.PendingPath = ""
	m.closeDialog()
	m.Status = "exported " + message.Path
	return m
}

func diagnosticsError(diagnostics []theme.Diagnostic) error {
	var messages []string
	for _, diagnostic := range diagnostics {
		if diagnostic.Severity == theme.Error {
			messages = append(messages, diagnostic.Field+": "+diagnostic.Message)
		}
	}
	return errors.New(strings.Join(messages, "; "))
}
