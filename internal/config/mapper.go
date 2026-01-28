package config

import (
	"bytes"
	"slices"
	"strings"
)

// PathMapper transforms paths based on a ParsedConfig.
type PathMapper struct {
	parsed *ParsedConfig
}

// NewPathMapper creates a PathMapper from a ParsedConfig.
// Returns nil if parsed is nil or empty.
func NewPathMapper(parsed *ParsedConfig) *PathMapper {
	if parsed.IsEmpty() {
		return nil
	}
	return &PathMapper{parsed: parsed}
}

// TransformPath transforms a relative path using the config rules.
// Returns the transformed path, or the original path if no rule matches.
//
// Transformation order:
// 1. Apply exact file match if present → render the output template
// 2. Apply directory prefix match to the result (or original if no file match)
// 3. If no matches, return path unchanged
func (m *PathMapper) TransformPath(relPath string, data any) (string, error) {
	if m == nil || m.parsed == nil {
		return relPath, nil
	}

	result := relPath

	// Step 1: Check for exact file match first
	if tmpl, ok := m.parsed.fileTemplates[relPath]; ok {
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return "", err
		}
		result = buf.String()

		// Validate rendered path
		if err := ValidateRenderedPath(result); err != nil {
			return "", err
		}
	}

	// Step 2: Check for directory prefix match on the result
	// (longest prefix first due to sorting)
	for _, dm := range m.parsed.dirMappings {
		if strings.HasPrefix(result, dm.prefix+"/") || result == dm.prefix {
			// Render the prefix template
			var buf bytes.Buffer
			if err := dm.tmpl.Execute(&buf, data); err != nil {
				return "", err
			}
			newPrefix := buf.String()

			// Validate rendered prefix
			if err := ValidateRenderedPath(newPrefix); err != nil {
				return "", err
			}

			// Append the suffix (everything after the prefix)
			suffix := strings.TrimPrefix(result, dm.prefix)
			return newPrefix + suffix, nil
		}
	}

	return result, nil
}

// CanOverwrite returns true if the source path allows overwriting existing files.
// Returns true (default) if no explicit overwrite:false is set for this path.
func (m *PathMapper) CanOverwrite(sourcePath string) bool {
	if m == nil || m.parsed == nil {
		return true
	}
	return !m.parsed.noOverwrite[sourcePath]
}

// ShouldSkipConfigFile returns true if the path is a render config file
// that should not be copied to output.
func ShouldSkipConfigFile(relPath string) bool {
	return slices.Contains(configFileNames, relPath)
}
