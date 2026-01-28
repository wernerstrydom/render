// Package data provides JSON and YAML data loading functionality.
package data

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Load reads data from a file and returns it as a generic interface.
// It automatically detects the format based on file extension.
func Load(path string) (any, error) {
	format, err := detectFormat(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %w", err)
	}
	defer func() { _ = f.Close() }()

	return LoadReader(f, format)
}

// LoadReader reads data from a reader in the specified format.
func LoadReader(r io.Reader, format string) (any, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	return Parse(content, format)
}

// Parse parses data from bytes in the specified format.
func Parse(content []byte, format string) (any, error) {
	var data any

	switch format {
	case "json":
		if err := json.Unmarshal(content, &data); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
	case "yaml", "yml":
		if err := yaml.Unmarshal(content, &data); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}
		// Convert YAML maps to string-keyed maps for consistency with JSON
		data = normalizeYAML(data)
	default:
		return nil, fmt.Errorf("unsupported data format: %s", format)
	}

	return data, nil
}

// detectFormat determines the data format from the file extension.
// Returns an error for unrecognized file extensions.
func detectFormat(path string) (string, error) {
	lower := strings.ToLower(path)
	switch {
	case strings.HasSuffix(lower, ".json"):
		return "json", nil
	case strings.HasSuffix(lower, ".yaml"):
		return "yaml", nil
	case strings.HasSuffix(lower, ".yml"):
		return "yaml", nil
	default:
		return "", fmt.Errorf("unsupported file extension for %q: expected .json, .yaml, or .yml", path)
	}
}

// normalizeYAML converts YAML map[string]any and map[any]any to map[string]any
// for consistency with JSON parsing.
//
// Note: This function does not handle cyclic data structures. However, cycles
// cannot occur in practice because YAML parsers (yaml.v3) do not create cyclic
// structures when unmarshaling - YAML aliases are resolved to separate copies.
func normalizeYAML(v any) any {
	switch val := v.(type) {
	case map[string]any:
		result := make(map[string]any, len(val))
		for k, v := range val {
			result[k] = normalizeYAML(v)
		}
		return result
	case map[any]any:
		result := make(map[string]any, len(val))
		for k, v := range val {
			result[fmt.Sprintf("%v", k)] = normalizeYAML(v)
		}
		return result
	case []any:
		result := make([]any, len(val))
		for i, v := range val {
			result[i] = normalizeYAML(v)
		}
		return result
	default:
		return v
	}
}
