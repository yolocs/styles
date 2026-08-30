// Package exporter dispatches canonical terminal themes to format renderers.
package exporter

import (
	"errors"
	"fmt"
	"io"

	"github.com/yolocs/styles/yolodev/internal/theme"
)

type renderFunc func(io.Writer, theme.Theme) error

// Registry maps export format names to renderers.
type Registry struct {
	renderers map[string]renderFunc
}

// NewRegistry returns an empty exporter registry.
func NewRegistry() *Registry {
	return &Registry{renderers: make(map[string]renderFunc)}
}

// Register adds a renderer for a format name.
func (r *Registry) Register(format string, render func(io.Writer, theme.Theme) error) error {
	if format == "" {
		return errors.New("register export format: format must not be empty")
	}
	if render == nil {
		return fmt.Errorf("register export format %q: renderer must not be nil", format)
	}
	if _, exists := r.renderers[format]; exists {
		return fmt.Errorf("register export format %q: already registered", format)
	}
	r.renderers[format] = render
	return nil
}

// Export renders a theme in the requested format.
func (r *Registry) Export(format string, w io.Writer, value theme.Theme) error {
	render, exists := r.renderers[format]
	if !exists {
		return fmt.Errorf("unsupported export format %q", format)
	}
	if err := render(w, value); err != nil {
		return fmt.Errorf("export %s: %w", format, err)
	}
	return nil
}
