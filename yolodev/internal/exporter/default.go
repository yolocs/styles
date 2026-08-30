package exporter

import "github.com/yolocs/styles/yolodev/internal/ghostty"

// NewDefaultRegistry returns the export formats built into yolodev.
func NewDefaultRegistry() *Registry {
	return &Registry{
		renderers: map[string]renderFunc{
			"ghostty": ghostty.Render,
		},
	}
}
