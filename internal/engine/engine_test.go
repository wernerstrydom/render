package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	eng := New()
	if eng == nil {
		t.Fatal("New() returned nil")
	}
	if eng.funcMap == nil {
		t.Fatal("New() funcMap is nil")
	}
}

func TestRenderString(t *testing.T) {
	eng := New()

	t.Run("simple template", func(t *testing.T) {
		tmpl := "Hello, {{ .name }}!"
		data := map[string]any{"name": "World"}

		result, err := eng.RenderString(tmpl, data)
		if err != nil {
			t.Fatalf("RenderString() error = %v", err)
		}
		if result != "Hello, World!" {
			t.Errorf("RenderString() = %q, want %q", result, "Hello, World!")
		}
	})

	t.Run("template with custom functions", func(t *testing.T) {
		tmpl := "{{ upper .name }}"
		data := map[string]any{"name": "hello"}

		result, err := eng.RenderString(tmpl, data)
		if err != nil {
			t.Fatalf("RenderString() error = %v", err)
		}
		if result != "HELLO" {
			t.Errorf("RenderString() = %q, want %q", result, "HELLO")
		}
	})

	t.Run("template with casing functions", func(t *testing.T) {
		tests := []struct {
			tmpl     string
			expected string
		}{
			{"{{ camelCase .name }}", "helloWorld"},
			{"{{ pascalCase .name }}", "HelloWorld"},
			{"{{ snakeCase .name }}", "hello_world"},
			{"{{ kebabCase .name }}", "hello-world"},
		}

		data := map[string]any{"name": "hello world"}

		for _, tt := range tests {
			result, err := eng.RenderString(tt.tmpl, data)
			if err != nil {
				t.Fatalf("RenderString(%q) error = %v", tt.tmpl, err)
			}
			if result != tt.expected {
				t.Errorf("RenderString(%q) = %q, want %q", tt.tmpl, result, tt.expected)
			}
		}
	})

	t.Run("template with range", func(t *testing.T) {
		tmpl := "{{ range .items }}{{ . }} {{ end }}"
		data := map[string]any{"items": []string{"a", "b", "c"}}

		result, err := eng.RenderString(tmpl, data)
		if err != nil {
			t.Fatalf("RenderString() error = %v", err)
		}
		if result != "a b c " {
			t.Errorf("RenderString() = %q, want %q", result, "a b c ")
		}
	})

	t.Run("template with conditionals", func(t *testing.T) {
		tmpl := "{{ if .enabled }}yes{{ else }}no{{ end }}"

		result1, _ := eng.RenderString(tmpl, map[string]any{"enabled": true})
		if result1 != "yes" {
			t.Errorf("RenderString(true) = %q, want 'yes'", result1)
		}

		result2, _ := eng.RenderString(tmpl, map[string]any{"enabled": false})
		if result2 != "no" {
			t.Errorf("RenderString(false) = %q, want 'no'", result2)
		}
	})

	t.Run("template with math", func(t *testing.T) {
		tmpl := "{{ add .a .b }}"
		data := map[string]any{"a": 5, "b": 3}

		result, err := eng.RenderString(tmpl, data)
		if err != nil {
			t.Fatalf("RenderString() error = %v", err)
		}
		if result != "8" {
			t.Errorf("RenderString() = %q, want '8'", result)
		}
	})

	t.Run("template with nested data", func(t *testing.T) {
		tmpl := "{{ .user.name }} - {{ .user.email }}"
		data := map[string]any{
			"user": map[string]any{
				"name":  "Alice",
				"email": "alice@example.com",
			},
		}

		result, err := eng.RenderString(tmpl, data)
		if err != nil {
			t.Fatalf("RenderString() error = %v", err)
		}
		if result != "Alice - alice@example.com" {
			t.Errorf("RenderString() = %q", result)
		}
	})

	t.Run("invalid template syntax", func(t *testing.T) {
		tmpl := "{{ .name"
		data := map[string]any{"name": "test"}

		_, err := eng.RenderString(tmpl, data)
		if err == nil {
			t.Error("RenderString() should return error for invalid syntax")
		}
	})

	t.Run("missing field", func(t *testing.T) {
		tmpl := "{{ .missing }}"
		data := map[string]any{"name": "test"}

		result, err := eng.RenderString(tmpl, data)
		if err != nil {
			t.Fatalf("RenderString() error = %v", err)
		}
		// Go templates output empty string for missing fields
		if result != "<no value>" {
			t.Errorf("RenderString() = %q, want '<no value>'", result)
		}
	})

	t.Run("template with default function", func(t *testing.T) {
		tmpl := `{{ default "default_value" .missing }}`
		data := map[string]any{}

		result, err := eng.RenderString(tmpl, data)
		if err != nil {
			t.Fatalf("RenderString() error = %v", err)
		}
		if result != "default_value" {
			t.Errorf("RenderString() = %q, want 'default_value'", result)
		}
	})

	t.Run("template with json function", func(t *testing.T) {
		tmpl := "{{ toJson .data }}"
		data := map[string]any{
			"data": map[string]any{"key": "value"},
		}

		result, err := eng.RenderString(tmpl, data)
		if err != nil {
			t.Fatalf("RenderString() error = %v", err)
		}
		if result != `{"key":"value"}` {
			t.Errorf("RenderString() = %q", result)
		}
	})

	t.Run("template with regex", func(t *testing.T) {
		tmpl := `{{ regexReplace "[0-9]+" "X" .text }}`
		data := map[string]any{"text": "abc123def456"}

		result, err := eng.RenderString(tmpl, data)
		if err != nil {
			t.Fatalf("RenderString() error = %v", err)
		}
		if result != "abcXdefX" {
			t.Errorf("RenderString() = %q, want 'abcXdefX'", result)
		}
	})
}

func TestRenderFile(t *testing.T) {
	eng := New()

	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "render-engine-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	t.Run("render file", func(t *testing.T) {
		path := filepath.Join(tmpDir, "test.tmpl")
		content := []byte("Hello, {{ .name }}!")
		if err := os.WriteFile(path, content, 0644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		data := map[string]any{"name": "World"}
		result, err := eng.RenderFile(path, data)
		if err != nil {
			t.Fatalf("RenderFile() error = %v", err)
		}
		if result != "Hello, World!" {
			t.Errorf("RenderFile() = %q, want %q", result, "Hello, World!")
		}
	})

	t.Run("render file with custom functions", func(t *testing.T) {
		path := filepath.Join(tmpDir, "funcs.tmpl")
		content := []byte("{{ upper .name }} {{ snakeCase .title }}")
		if err := os.WriteFile(path, content, 0644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		data := map[string]any{"name": "hello", "title": "MyProject"}
		result, err := eng.RenderFile(path, data)
		if err != nil {
			t.Fatalf("RenderFile() error = %v", err)
		}
		if result != "HELLO my_project" {
			t.Errorf("RenderFile() = %q", result)
		}
	})

	t.Run("render non-existent file", func(t *testing.T) {
		_, err := eng.RenderFile(filepath.Join(tmpDir, "nonexistent.tmpl"), nil)
		if err == nil {
			t.Error("RenderFile() should return error for non-existent file")
		}
	})

	t.Run("render file with invalid template", func(t *testing.T) {
		path := filepath.Join(tmpDir, "invalid.tmpl")
		content := []byte("{{ .name")
		if err := os.WriteFile(path, content, 0644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		_, err := eng.RenderFile(path, map[string]any{"name": "test"})
		if err == nil {
			t.Error("RenderFile() should return error for invalid template")
		}
	})
}

func TestComplexTemplates(t *testing.T) {
	eng := New()

	t.Run("generate code", func(t *testing.T) {
		tmpl := `package {{ .package }}

type {{ pascalCase .name }} struct {
{{- range .fields }}
	{{ pascalCase .name }} {{ .type }}
{{- end }}
}
`
		data := map[string]any{
			"package": "models",
			"name":    "user profile",
			"fields": []map[string]any{
				{"name": "user name", "type": "string"},
				{"name": "email address", "type": "string"},
				{"name": "age", "type": "int"},
			},
		}

		result, err := eng.RenderString(tmpl, data)
		if err != nil {
			t.Fatalf("RenderString() error = %v", err)
		}

		expected := `package models

type UserProfile struct {
	UserName string
	EmailAddress string
	Age int
}
`
		if result != expected {
			t.Errorf("RenderString() =\n%s\nwant:\n%s", result, expected)
		}
	})

	t.Run("generate config", func(t *testing.T) {
		tmpl := `database:
  host: {{ .db.host }}
  port: {{ .db.port }}
  name: {{ .db.name }}
  connection_string: "postgresql://{{ .db.host }}:{{ .db.port }}/{{ .db.name }}"
`
		data := map[string]any{
			"db": map[string]any{
				"host": "localhost",
				"port": 5432,
				"name": "myapp",
			},
		}

		result, err := eng.RenderString(tmpl, data)
		if err != nil {
			t.Fatalf("RenderString() error = %v", err)
		}

		if !contains(result, "connection_string: \"postgresql://localhost:5432/myapp\"") {
			t.Errorf("RenderString() missing connection string:\n%s", result)
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr, 0))
}

func containsAt(s, substr string, start int) bool {
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
