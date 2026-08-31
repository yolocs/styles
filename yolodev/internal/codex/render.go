// Package codex renders canonical terminal themes as Codex TextMate themes.
package codex

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/yolocs/styles/yolodev/internal/theme"
)

type style struct {
	name       string
	scope      string
	foreground string
}

// Render validates value and writes a deterministic Codex TextMate theme to w.
func Render(w io.Writer, value theme.Theme) error {
	diagnostics := theme.Validate(value)
	if theme.HasErrors(diagnostics) {
		var messages []string
		for _, diagnostic := range diagnostics {
			if diagnostic.Severity == theme.Error {
				messages = append(messages, diagnostic.Field+": "+diagnostic.Message)
			}
		}
		return fmt.Errorf("render Codex theme: invalid terminal theme: %s", strings.Join(messages, "; "))
	}

	value = theme.Normalize(value)
	styles := []style{
		{name: "Comments", scope: "comment, punctuation.definition.comment", foreground: value.Palette.Bright.Black},
		{name: "Strings", scope: "string, constant.other.symbol", foreground: value.Palette.Normal.Green},
		{name: "Numbers", scope: "constant.numeric", foreground: value.Palette.Normal.Yellow},
		{name: "Keywords", scope: "keyword, storage", foreground: value.Palette.Bright.Magenta},
		{name: "Functions", scope: "entity.name.function, support.function, variable.function", foreground: value.Palette.Normal.Cyan},
		{name: "Types", scope: "entity.name.type, entity.name.class, support.type, support.class", foreground: value.Palette.Bright.Blue},
		{name: "Invalid", scope: "invalid, invalid.illegal", foreground: value.Palette.Bright.Red},
		{name: "Deleted", scope: "markup.deleted, diff.deleted", foreground: value.Palette.Normal.Red},
		{name: "Inserted", scope: "markup.inserted, diff.inserted", foreground: value.Palette.Normal.Green},
		{name: "Changed", scope: "markup.changed, diff.changed", foreground: value.Palette.Normal.Yellow},
	}

	var output bytes.Buffer
	output.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	output.WriteString("<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"https://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n")
	output.WriteString("<plist version=\"1.0\">\n")
	output.WriteString("  <dict>\n")
	if err := writeKeyString(&output, 2, "name", value.Theme.Name+" "+value.Theme.Variant); err != nil {
		return fmt.Errorf("render Codex theme metadata: %w", err)
	}
	output.WriteString("    <key>settings</key>\n")
	output.WriteString("    <array>\n")
	output.WriteString("      <dict>\n")
	output.WriteString("        <key>settings</key>\n")
	output.WriteString("        <dict>\n")
	globalColors := []struct {
		key   string
		value string
	}{
		{key: "background", value: value.Colors.Background},
		{key: "foreground", value: value.Colors.Foreground},
		{key: "caret", value: value.Colors.Cursor},
		{key: "selection", value: value.Colors.SelectionBackground},
		{key: "selectionForeground", value: value.Colors.SelectionForeground},
	}
	for _, color := range globalColors {
		if err := writeKeyString(&output, 5, color.key, color.value); err != nil {
			return fmt.Errorf("render Codex theme global colors: %w", err)
		}
	}
	output.WriteString("        </dict>\n")
	output.WriteString("      </dict>\n")
	for _, entry := range styles {
		output.WriteString("      <dict>\n")
		if err := writeKeyString(&output, 4, "name", entry.name); err != nil {
			return fmt.Errorf("render Codex theme style name: %w", err)
		}
		if err := writeKeyString(&output, 4, "scope", entry.scope); err != nil {
			return fmt.Errorf("render Codex theme style scope: %w", err)
		}
		output.WriteString("        <key>settings</key>\n")
		output.WriteString("        <dict>\n")
		if err := writeKeyString(&output, 5, "foreground", entry.foreground); err != nil {
			return fmt.Errorf("render Codex theme style color: %w", err)
		}
		output.WriteString("        </dict>\n")
		output.WriteString("      </dict>\n")
	}
	output.WriteString("    </array>\n")
	output.WriteString("  </dict>\n")
	output.WriteString("</plist>\n")

	if _, err := io.Copy(w, &output); err != nil {
		return fmt.Errorf("write Codex theme: %w", err)
	}
	return nil
}

func writeKeyString(output *bytes.Buffer, indent int, key, value string) error {
	padding := strings.Repeat("  ", indent)
	output.WriteString(padding + "<key>")
	if err := xml.EscapeText(output, []byte(key)); err != nil {
		return err
	}
	output.WriteString("</key>\n")
	output.WriteString(padding + "<string>")
	if err := xml.EscapeText(output, []byte(value)); err != nil {
		return err
	}
	output.WriteString("</string>\n")
	return nil
}
