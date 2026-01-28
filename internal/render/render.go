// Package render provides two-phase template rendering with validation.
package render

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wernerstrydom/render/internal/config"
	"github.com/wernerstrydom/render/internal/engine"
	"github.com/wernerstrydom/render/internal/output"
)

// Output represents a single file to be written.
type Output struct {
	SourcePath  string      // Relative path in template dir
	OutputPath  string      // Absolute path in output dir
	Content     []byte      // Rendered content (nil for copied files)
	CopyFrom    string      // Source path if copying verbatim (empty if rendered)
	Permissions os.FileMode // File permissions to apply
	Overwrite   bool        // Whether to overwrite existing files (default true)
}

// Plan represents the complete rendering operation.
type Plan struct {
	Outputs []Output
}

// CollectConfig configures the Collect function.
type CollectConfig struct {
	TemplateDir string
	OutputDir   string
	Data        any
	Config      *config.ParsedConfig // nil = no path transformation
	Engine      *engine.Engine
}

// Collect walks the template directory and builds a Plan.
// It collects all outputs into memory for validation before any writes.
func Collect(cfg CollectConfig) (*Plan, error) {
	// Resolve template directory to absolute path
	tmplDirAbs, err := filepath.Abs(cfg.TemplateDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve template directory: %w", err)
	}

	// Verify template directory exists
	info, err := os.Stat(tmplDirAbs)
	if err != nil {
		return nil, fmt.Errorf("failed to access template directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("template path is not a directory: %s", cfg.TemplateDir)
	}

	// Resolve output directory to absolute path
	outDirAbs, err := filepath.Abs(cfg.OutputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve output directory: %w", err)
	}

	// Create path mapper if config exists
	mapper := config.NewPathMapper(cfg.Config)

	plan := &Plan{
		Outputs: make([]Output, 0),
	}

	err = filepath.Walk(tmplDirAbs, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
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

		// Skip config files
		if config.ShouldSkipConfigFile(relPath) {
			return nil
		}

		// Security: Ensure relative path doesn't escape
		if strings.Contains(relPath, "..") {
			return fmt.Errorf("security error: path contains directory traversal: %s", relPath)
		}

		// Transform path using config
		outputRelPath := relPath
		if mapper != nil {
			outputRelPath, err = mapper.TransformPath(relPath, cfg.Data)
			if err != nil {
				return fmt.Errorf("failed to transform path %s: %w", relPath, err)
			}
		}

		// Calculate output path
		outPath := filepath.Join(outDirAbs, outputRelPath)

		// Security: Verify output path is within output directory
		if !strings.HasPrefix(outPath, outDirAbs) {
			return fmt.Errorf("security error: output path %s is outside output directory", outPath)
		}

		// Handle directories - just ensure they'll be created
		if info.IsDir() {
			return nil
		}

		// Determine if this file can overwrite existing files
		canOverwrite := true
		if mapper != nil {
			canOverwrite = mapper.CanOverwrite(relPath)
		}

		// Check if file is a template
		if strings.HasSuffix(path, ".tmpl") {
			// Strip .tmpl extension for output
			outPath = strings.TrimSuffix(outPath, ".tmpl")

			// Read and render template
			tmplContent, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read template %s: %w", relPath, err)
			}

			result, err := cfg.Engine.RenderString(string(tmplContent), cfg.Data)
			if err != nil {
				return fmt.Errorf("failed to render template %s: %w", relPath, err)
			}

			plan.Outputs = append(plan.Outputs, Output{
				SourcePath:  relPath,
				OutputPath:  outPath,
				Content:     []byte(result),
				Permissions: 0644,
				Overwrite:   canOverwrite,
			})
		} else {
			// Non-template file - will be copied
			srcInfo, err := os.Stat(path)
			if err != nil {
				return fmt.Errorf("failed to stat source file %s: %w", relPath, err)
			}

			plan.Outputs = append(plan.Outputs, Output{
				SourcePath:  relPath,
				OutputPath:  outPath,
				CopyFrom:    path,
				Permissions: srcInfo.Mode().Perm(),
				Overwrite:   canOverwrite,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return plan, nil
}

// Validate checks the Plan for collisions and security issues.
// Returns all errors found (doesn't stop at first error).
func (p *Plan) Validate() []error {
	var errs []error
	seen := make(map[string]string) // outputPath → sourcePath

	for _, out := range p.Outputs {
		// Check for collisions
		if existing, ok := seen[out.OutputPath]; ok {
			errs = append(errs, fmt.Errorf(
				"output path collision: %q produced by both:\n  - %s\n  - %s",
				out.OutputPath, existing, out.SourcePath))
		}
		seen[out.OutputPath] = out.SourcePath
	}

	return errs
}

// Preview returns a human-readable summary of planned outputs.
func (p *Plan) Preview() string {
	var sb strings.Builder
	sb.WriteString("Planned outputs:\n")
	for _, out := range p.Outputs {
		action := "render"
		if out.CopyFrom != "" {
			action = "copy"
		}
		sb.WriteString(fmt.Sprintf("  [%s] %s → %s\n",
			action, out.SourcePath, out.OutputPath))
	}
	return sb.String()
}

// ExecuteResult contains information about executed outputs.
type ExecuteResult struct {
	Skipped map[string]bool // Paths that were skipped due to no-overwrite
}

// Execute writes all files in the Plan.
// Returns ExecuteResult with information about skipped files.
func (p *Plan) Execute(writer *output.Writer) (*ExecuteResult, error) {
	result := &ExecuteResult{
		Skipped: make(map[string]bool),
	}

	for _, out := range p.Outputs {
		// Ensure parent directory exists
		dir := filepath.Dir(out.OutputPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		if out.CopyFrom != "" {
			// Copy file
			if !out.Overwrite {
				skipped, err := writer.CopyIfNotExists(out.CopyFrom, out.OutputPath)
				if err != nil {
					return nil, err
				}
				if skipped {
					result.Skipped[out.OutputPath] = true
				}
			} else {
				if err := writer.Copy(out.CopyFrom, out.OutputPath); err != nil {
					return nil, err
				}
			}
		} else {
			// Write rendered content
			if !out.Overwrite {
				skipped, err := writer.WriteIfNotExists(out.OutputPath, out.Content, out.Permissions)
				if err != nil {
					return nil, err
				}
				if skipped {
					result.Skipped[out.OutputPath] = true
				}
			} else {
				if err := writer.WriteWithPerm(out.OutputPath, out.Content, out.Permissions); err != nil {
					return nil, err
				}
			}
		}
	}

	return result, nil
}

// OutputCount returns the number of outputs in the plan.
func (p *Plan) OutputCount() int {
	return len(p.Outputs)
}
