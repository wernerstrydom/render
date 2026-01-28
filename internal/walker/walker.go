// Package walker provides directory traversal and template rendering functionality.
package walker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wernerstrydom/render/internal/engine"
	"github.com/wernerstrydom/render/internal/output"
)

// Config holds configuration for directory walking and rendering.
type Config struct {
	// TemplateDir is the source directory containing templates.
	TemplateDir string

	// OutputDir is the destination directory for rendered output.
	OutputDir string

	// Data is the template data to use for rendering.
	Data any

	// Engine is the template rendering engine.
	Engine *engine.Engine

	// Writer is the output writer.
	Writer *output.Writer

	// TransformPath is an optional callback to transform relative paths.
	// If nil, paths are used as-is.
	// The function receives the relative path and data, and returns the transformed path.
	TransformPath func(relPath string, data any) (string, error)

	// ValidateSymlinks enables symlink security validation.
	// When true, symlinks that resolve outside the template directory cause an error.
	ValidateSymlinks bool

	// OnRendered is called after a template file is rendered.
	// Parameters: source relative path, destination absolute path.
	OnRendered func(srcRel, dstAbs string)

	// OnCopied is called after a non-template file is copied.
	// Parameters: source relative path, destination absolute path.
	OnCopied func(srcRel, dstAbs string)
}

// Walk traverses a template directory, rendering .tmpl files and copying others.
func Walk(cfg Config) error {
	// Resolve the template directory to an absolute path
	tmplDirAbs, err := filepath.Abs(cfg.TemplateDir)
	if err != nil {
		return fmt.Errorf("failed to resolve template directory path: %w", err)
	}

	// Verify template directory exists
	info, err := os.Stat(tmplDirAbs)
	if err != nil {
		return fmt.Errorf("failed to access template directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("template path is not a directory: %s", cfg.TemplateDir)
	}

	// Resolve symlinks for the template directory for consistent comparison
	var tmplDirReal string
	if cfg.ValidateSymlinks {
		tmplDirReal, err = filepath.EvalSymlinks(tmplDirAbs)
		if err != nil {
			return fmt.Errorf("failed to resolve template directory symlinks: %w", err)
		}
	}

	// Resolve output directory to absolute path
	outDirAbs, err := filepath.Abs(cfg.OutputDir)
	if err != nil {
		return fmt.Errorf("failed to resolve output directory path: %w", err)
	}

	return filepath.Walk(tmplDirAbs, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Security: Resolve symlinks and verify the real path is within the template directory
		if cfg.ValidateSymlinks {
			realPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				return fmt.Errorf("failed to resolve symlink %s: %w", path, err)
			}
			if !strings.HasPrefix(realPath, tmplDirReal) && realPath != tmplDirReal {
				return fmt.Errorf("security error: path %s resolves outside template directory", path)
			}
		}

		// Get relative path from template directory
		relPath, err := filepath.Rel(tmplDirAbs, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Skip the root directory
		if relPath == "." {
			return nil
		}

		// Security: Ensure relative path doesn't escape output directory
		if strings.Contains(relPath, "..") {
			return fmt.Errorf("security error: path contains directory traversal: %s", relPath)
		}

		// Transform path if callback is provided
		outputRelPath := relPath
		if cfg.TransformPath != nil {
			outputRelPath, err = cfg.TransformPath(relPath, cfg.Data)
			if err != nil {
				return fmt.Errorf("failed to transform path %s: %w", relPath, err)
			}
		}

		// Calculate output path
		outPath := filepath.Join(outDirAbs, outputRelPath)

		// Security: Verify output path is within output directory
		if cfg.ValidateSymlinks && !strings.HasPrefix(outPath, outDirAbs) {
			return fmt.Errorf("security error: output path %s is outside output directory", outPath)
		}

		// Handle directories
		if info.IsDir() {
			if err := os.MkdirAll(outPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", outPath, err)
			}
			return nil
		}

		// Check if file is a template
		if strings.HasSuffix(path, ".tmpl") {
			// Strip .tmpl extension for output
			outPath = strings.TrimSuffix(outPath, ".tmpl")

			// Read and render template
			tmplContent, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read template %s: %w", path, err)
			}

			result, err := cfg.Engine.RenderString(string(tmplContent), cfg.Data)
			if err != nil {
				return fmt.Errorf("failed to render template %s: %w", path, err)
			}

			if err := cfg.Writer.WriteString(outPath, result); err != nil {
				return err
			}

			if cfg.OnRendered != nil {
				cfg.OnRendered(relPath, outPath)
			}
		} else {
			// Copy file verbatim
			if err := cfg.Writer.Copy(path, outPath); err != nil {
				return err
			}

			if cfg.OnCopied != nil {
				cfg.OnCopied(relPath, outPath)
			}
		}

		return nil
	})
}
