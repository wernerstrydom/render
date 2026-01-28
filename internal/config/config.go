// Package config provides configuration loading for .render.yaml files.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"text/template"

	"github.com/wernerstrydom/render/internal/funcs"
	"gopkg.in/yaml.v3"
)

// PathMapping represents a path mapping which can be either a simple string
// or an object with path and overwrite options.
type PathMapping struct {
	Path      string `json:"path" yaml:"path"`
	Overwrite *bool  `json:"overwrite" yaml:"overwrite"` // nil = true (default)
}

// UnmarshalYAML implements custom YAML unmarshaling to support both string
// and object formats for path mappings.
func (p *PathMapping) UnmarshalYAML(value *yaml.Node) error {
	// Try string first
	if value.Kind == yaml.ScalarNode {
		p.Path = value.Value
		p.Overwrite = nil // default to true
		return nil
	}

	// Try object format
	if value.Kind == yaml.MappingNode {
		type rawPathMapping struct {
			Path      string `yaml:"path"`
			Overwrite *bool  `yaml:"overwrite"`
		}
		var raw rawPathMapping
		if err := value.Decode(&raw); err != nil {
			return err
		}
		if raw.Path == "" {
			return fmt.Errorf("path mapping object must have 'path' field")
		}
		p.Path = raw.Path
		p.Overwrite = raw.Overwrite
		return nil
	}

	return fmt.Errorf("path mapping must be a string or object, got %v", value.Kind)
}

// Config represents the raw .render.yaml configuration.
type Config struct {
	Paths map[string]PathMapping `json:"paths" yaml:"paths"`
}

// dirMapping holds a directory prefix mapping with its parsed template.
type dirMapping struct {
	prefix string
	tmpl   *template.Template
}

// ParsedConfig holds validated, pre-parsed configuration.
type ParsedConfig struct {
	fileTemplates map[string]*template.Template // Exact file mappings
	dirMappings   []dirMapping                  // Prefix mappings, sorted longest first
	noOverwrite   map[string]bool               // Source paths with overwrite: false
}

// configFileNames lists the supported config file names in priority order.
var configFileNames = []string{".render.yaml", ".render.yml", "render.json"}

// Load finds and loads a render config from the template directory.
// Returns nil (not an error) if no config file exists.
func Load(tmplDir string) (*ParsedConfig, error) {
	// Try each config file name in order
	var configPath string
	for _, name := range configFileNames {
		path := filepath.Join(tmplDir, name)
		if _, err := os.Stat(path); err == nil {
			configPath = path
			break
		}
	}

	// No config file found - this is not an error
	if configPath == "" {
		return nil, nil
	}

	return LoadFile(configPath, tmplDir)
}

// LoadFile loads and parses a config file.
func LoadFile(configPath, tmplDir string) (*ParsedConfig, error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return Parse(content, tmplDir, filepath.Base(configPath))
}

// Parse parses config content and validates it against the template directory.
func Parse(content []byte, tmplDir, filename string) (*ParsedConfig, error) {
	// First, validate schema by parsing into raw map
	var raw map[string]any
	if err := yaml.Unmarshal(content, &raw); err != nil {
		return nil, fmt.Errorf("%s: invalid YAML: %w", filename, err)
	}

	// Check for unknown keys
	for key := range raw {
		if key != "paths" {
			return nil, fmt.Errorf("%s: unknown key %q (only 'paths' allowed)", filename, key)
		}
	}

	// Parse into typed config
	var cfg Config
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return nil, fmt.Errorf("%s: failed to parse config: %w", filename, err)
	}

	// Empty config is valid but has nothing to transform
	if len(cfg.Paths) == 0 {
		return &ParsedConfig{
			fileTemplates: make(map[string]*template.Template),
			dirMappings:   nil,
			noOverwrite:   make(map[string]bool),
		}, nil
	}

	// Validate and parse all path mappings
	parsed := &ParsedConfig{
		fileTemplates: make(map[string]*template.Template),
		dirMappings:   nil,
		noOverwrite:   make(map[string]bool),
	}

	funcMap := funcs.Map()

	for src, mapping := range cfg.Paths {
		// Validate source path
		if err := validateSourcePath(src); err != nil {
			return nil, fmt.Errorf("%s: paths[%q]: %w", filename, src, err)
		}

		// Check if source exists and determine if it's a file or directory
		srcPath := filepath.Join(tmplDir, src)
		info, err := os.Stat(srcPath)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("%s: paths[%q]: source does not exist in template directory", filename, src)
			}
			return nil, fmt.Errorf("%s: paths[%q]: %w", filename, src, err)
		}

		// Parse the destination template
		tmpl, err := template.New(src).Funcs(funcMap).Parse(mapping.Path)
		if err != nil {
			return nil, fmt.Errorf("%s: paths[%q]: invalid template syntax: %w", filename, src, err)
		}

		// Track overwrite setting (nil or true means overwrite allowed; false means no overwrite)
		if mapping.Overwrite != nil && !*mapping.Overwrite {
			parsed.noOverwrite[src] = true
		}

		if info.IsDir() {
			// Directory prefix mapping
			parsed.dirMappings = append(parsed.dirMappings, dirMapping{
				prefix: src,
				tmpl:   tmpl,
			})
		} else {
			// File mapping
			parsed.fileTemplates[src] = tmpl
		}
	}

	// Sort directory mappings by prefix length (longest first)
	sort.Slice(parsed.dirMappings, func(i, j int) bool {
		return len(parsed.dirMappings[i].prefix) > len(parsed.dirMappings[j].prefix)
	})

	return parsed, nil
}

// validateSourcePath checks that a source path is safe.
func validateSourcePath(path string) error {
	// Check for absolute path
	if filepath.IsAbs(path) {
		return fmt.Errorf("source path must be relative (got absolute path)")
	}

	// Check for null bytes
	if strings.ContainsRune(path, '\x00') {
		return fmt.Errorf("source path contains null byte")
	}

	// Check for path traversal
	parts := strings.FieldsFunc(path, func(r rune) bool {
		return r == '/' || r == filepath.Separator
	})
	if slices.Contains(parts, "..") {
		return fmt.Errorf("source path contains '..' (directory traversal)")
	}

	return nil
}

// ValidateRenderedPath checks that a rendered output path is safe.
func ValidateRenderedPath(path string) error {
	// Check for path traversal in rendered output
	parts := strings.FieldsFunc(path, func(r rune) bool {
		return r == '/' || r == filepath.Separator
	})
	if slices.Contains(parts, "..") {
		return fmt.Errorf("rendered path contains '..' (directory traversal)")
	}

	return nil
}

// IsEmpty returns true if the config has no path mappings.
func (p *ParsedConfig) IsEmpty() bool {
	return p == nil || (len(p.fileTemplates) == 0 && len(p.dirMappings) == 0)
}

// HasFileMappings returns true if there are exact file mappings.
func (p *ParsedConfig) HasFileMappings() bool {
	return p != nil && len(p.fileTemplates) > 0
}

// HasDirMappings returns true if there are directory prefix mappings.
func (p *ParsedConfig) HasDirMappings() bool {
	return p != nil && len(p.dirMappings) > 0
}
