// Package exporter dispatches canonical terminal themes to format renderers.
package exporter

import (
	"errors"
	"fmt"
	"io"

	"github.com/yolocs/styles/yolodev/internal/theme"
)

type renderFunc func(io.Writer, theme.Theme) error

type registeredFormat struct {
	extension string
	render    renderFunc
}

// Registry maps export format names to renderers.
type Registry struct {
	formats map[string]registeredFormat
}

// NewRegistry returns an empty exporter registry.
func NewRegistry() *Registry {
	return &Registry{formats: make(map[string]registeredFormat)}
}

// Register adds a renderer for a format name.
func (r *Registry) Register(format string, render func(io.Writer, theme.Theme) error) error {
	if format == "" {
		return errors.New("register export format: format must not be empty")
	}
	if render == nil {
		return fmt.Errorf("register export format %q: renderer must not be nil", format)
	}
	if _, exists := r.formats[format]; exists {
		return fmt.Errorf("register export format %q: already registered", format)
	}
	r.formats[format] = registeredFormat{render: render}
	return nil
}

// Extension returns the preferred file extension for an export format.
func (r *Registry) Extension(format string) (string, error) {
	registered, exists := r.formats[format]
	if !exists {
		return "", fmt.Errorf("unsupported export format %q", format)
	}
	if registered.extension == "" {
		return "", fmt.Errorf("export format %q has no file extension", format)
	}
	return registered.extension, nil
}

// Export renders a theme in the requested format.
func (r *Registry) Export(format string, w io.Writer, value theme.Theme) error {
	registered, exists := r.formats[format]
	if !exists {
		return fmt.Errorf("unsupported export format %q", format)
	}
	if err := registered.render(w, value); err != nil {
		return fmt.Errorf("export %s: %w", format, err)
	}
	return nil
}
