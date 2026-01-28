// Package engine provides the template rendering engine.
package engine

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/wernerstrydom/render/internal/funcs"
)

// Engine handles template parsing and execution.
type Engine struct {
	funcMap template.FuncMap
}

// New creates a new template engine with custom functions.
func New() *Engine {
	return &Engine{
		funcMap: funcs.Map(),
	}
}

// RenderString renders a template string with the given data.
func (e *Engine) RenderString(tmpl string, data any) (string, error) {
	t, err := template.New("template").Funcs(e.funcMap).Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// RenderFile renders a template file with the given data.
func (e *Engine) RenderFile(path string, data any) (string, error) {
	t, err := template.New("").Funcs(e.funcMap).ParseFiles(path)
	if err != nil {
		return "", fmt.Errorf("failed to parse template file %s: %w", path, err)
	}

	// ParseFiles creates a template with the base name of the file
	t = t.Lookup(filepath.Base(path))
	if t == nil {
		return "", fmt.Errorf("template not found after parsing: %s", path)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", path, err)
	}

	return buf.String(), nil
}
