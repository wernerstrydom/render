package data

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	t.Run("parse JSON object", func(t *testing.T) {
		content := []byte(`{"name": "test", "value": 123}`)
		result, err := Parse(content, "json")
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		m, ok := result.(map[string]any)
		if !ok {
			t.Fatalf("Parse() result is not a map")
		}
		if m["name"] != "test" {
			t.Errorf("Parse() name = %v, want 'test'", m["name"])
		}
		if m["value"] != float64(123) {
			t.Errorf("Parse() value = %v, want 123", m["value"])
		}
	})

	t.Run("parse JSON array", func(t *testing.T) {
		content := []byte(`[1, 2, 3]`)
		result, err := Parse(content, "json")
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		arr, ok := result.([]any)
		if !ok {
			t.Fatalf("Parse() result is not an array")
		}
		if len(arr) != 3 {
			t.Errorf("Parse() len = %d, want 3", len(arr))
		}
	})

	t.Run("parse YAML object", func(t *testing.T) {
		content := []byte("name: test\nvalue: 123")
		result, err := Parse(content, "yaml")
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		m, ok := result.(map[string]any)
		if !ok {
			t.Fatalf("Parse() result is not a map")
		}
		if m["name"] != "test" {
			t.Errorf("Parse() name = %v, want 'test'", m["name"])
		}
		if m["value"] != 123 {
			t.Errorf("Parse() value = %v, want 123", m["value"])
		}
	})

	t.Run("parse YAML with nested structures", func(t *testing.T) {
		content := []byte(`
name: test
nested:
  key: value
items:
  - first
  - second
`)
		result, err := Parse(content, "yaml")
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		m, ok := result.(map[string]any)
		if !ok {
			t.Fatalf("Parse() result is not a map")
		}

		nested, ok := m["nested"].(map[string]any)
		if !ok {
			t.Fatalf("Parse() nested is not a map")
		}
		if nested["key"] != "value" {
			t.Errorf("Parse() nested.key = %v, want 'value'", nested["key"])
		}

		items, ok := m["items"].([]any)
		if !ok {
			t.Fatalf("Parse() items is not an array")
		}
		if len(items) != 2 {
			t.Errorf("Parse() items len = %d, want 2", len(items))
		}
	})

	t.Run("parse invalid JSON", func(t *testing.T) {
		content := []byte(`{invalid}`)
		_, err := Parse(content, "json")
		if err == nil {
			t.Error("Parse() should return error for invalid JSON")
		}
	})

	t.Run("parse invalid YAML", func(t *testing.T) {
		content := []byte(":\n  :\n    invalid")
		_, err := Parse(content, "yaml")
		if err == nil {
			t.Error("Parse() should return error for invalid YAML")
		}
	})

	t.Run("unsupported format", func(t *testing.T) {
		content := []byte(`test`)
		_, err := Parse(content, "xml")
		if err == nil {
			t.Error("Parse() should return error for unsupported format")
		}
	})
}

func TestLoadReader(t *testing.T) {
	t.Run("load JSON from reader", func(t *testing.T) {
		reader := strings.NewReader(`{"key": "value"}`)
		result, err := LoadReader(reader, "json")
		if err != nil {
			t.Fatalf("LoadReader() error = %v", err)
		}

		m, ok := result.(map[string]any)
		if !ok {
			t.Fatalf("LoadReader() result is not a map")
		}
		if m["key"] != "value" {
			t.Errorf("LoadReader() key = %v, want 'value'", m["key"])
		}
	})

	t.Run("load YAML from reader", func(t *testing.T) {
		reader := strings.NewReader("key: value")
		result, err := LoadReader(reader, "yaml")
		if err != nil {
			t.Fatalf("LoadReader() error = %v", err)
		}

		m, ok := result.(map[string]any)
		if !ok {
			t.Fatalf("LoadReader() result is not a map")
		}
		if m["key"] != "value" {
			t.Errorf("LoadReader() key = %v, want 'value'", m["key"])
		}
	})
}

func TestLoad(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "render-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	t.Run("load JSON file", func(t *testing.T) {
		path := filepath.Join(tmpDir, "test.json")
		content := []byte(`{"name": "test"}`)
		if err := os.WriteFile(path, content, 0644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		result, err := Load(path)
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		m, ok := result.(map[string]any)
		if !ok {
			t.Fatalf("Load() result is not a map")
		}
		if m["name"] != "test" {
			t.Errorf("Load() name = %v, want 'test'", m["name"])
		}
	})

	t.Run("load YAML file", func(t *testing.T) {
		path := filepath.Join(tmpDir, "test.yaml")
		content := []byte("name: test")
		if err := os.WriteFile(path, content, 0644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		result, err := Load(path)
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		m, ok := result.(map[string]any)
		if !ok {
			t.Fatalf("Load() result is not a map")
		}
		if m["name"] != "test" {
			t.Errorf("Load() name = %v, want 'test'", m["name"])
		}
	})

	t.Run("load YML file", func(t *testing.T) {
		path := filepath.Join(tmpDir, "test.yml")
		content := []byte("name: test")
		if err := os.WriteFile(path, content, 0644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		result, err := Load(path)
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		m, ok := result.(map[string]any)
		if !ok {
			t.Fatalf("Load() result is not a map")
		}
		if m["name"] != "test" {
			t.Errorf("Load() name = %v, want 'test'", m["name"])
		}
	})

	t.Run("load non-existent file", func(t *testing.T) {
		_, err := Load(filepath.Join(tmpDir, "nonexistent.json"))
		if err == nil {
			t.Error("Load() should return error for non-existent file")
		}
	})
}

func TestDetectFormat(t *testing.T) {
	validTests := []struct {
		path     string
		expected string
	}{
		{"file.json", "json"},
		{"file.JSON", "json"},
		{"file.yaml", "yaml"},
		{"file.YAML", "yaml"},
		{"file.yml", "yaml"},
		{"file.YML", "yaml"},
	}

	for _, tt := range validTests {
		t.Run(tt.path, func(t *testing.T) {
			result, err := detectFormat(tt.path)
			if err != nil {
				t.Fatalf("detectFormat(%q) unexpected error: %v", tt.path, err)
			}
			if result != tt.expected {
				t.Errorf("detectFormat(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}

	// Test unsupported extensions return errors
	unsupportedTests := []string{"file.txt", "file", "file.xml"}
	for _, path := range unsupportedTests {
		t.Run("unsupported_"+path, func(t *testing.T) {
			_, err := detectFormat(path)
			if err == nil {
				t.Errorf("detectFormat(%q) should return error for unsupported extension", path)
			}
		})
	}
}

func TestNormalizeYAML(t *testing.T) {
	t.Run("normalize map[string]any", func(t *testing.T) {
		input := map[string]any{
			"key": "value",
			"nested": map[string]any{
				"inner": "data",
			},
		}
		result := normalizeYAML(input).(map[string]any)
		if result["key"] != "value" {
			t.Errorf("normalizeYAML() key = %v, want 'value'", result["key"])
		}
	})

	t.Run("normalize array", func(t *testing.T) {
		input := []any{"a", "b", "c"}
		result := normalizeYAML(input).([]any)
		if len(result) != 3 {
			t.Errorf("normalizeYAML() len = %d, want 3", len(result))
		}
	})

	t.Run("normalize scalar", func(t *testing.T) {
		if result := normalizeYAML("test"); result != "test" {
			t.Errorf("normalizeYAML(string) = %v, want 'test'", result)
		}
		if result := normalizeYAML(123); result != 123 {
			t.Errorf("normalizeYAML(int) = %v, want 123", result)
		}
	})
}
