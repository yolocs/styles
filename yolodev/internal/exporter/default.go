package exporter

import (
	"github.com/yolocs/styles/yolodev/internal/codex"
	"github.com/yolocs/styles/yolodev/internal/ghostty"
)

// NewDefaultRegistry returns the export formats built into yolodev.
func NewDefaultRegistry() *Registry {
	return &Registry{
		formats: map[string]registeredFormat{
			"codex":   {extension: ".tmTheme", render: codex.Render},
			"ghostty": {extension: ".ghostty", render: ghostty.Render},
		},
	}
}
