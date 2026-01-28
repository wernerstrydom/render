package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPathMapper_NilMapper(t *testing.T) {
	var mapper *PathMapper

	result, err := mapper.TransformPath("src/main.go", nil)
	if err != nil {
		t.Fatalf("TransformPath failed: %v", err)
	}

	if result != "src/main.go" {
		t.Errorf("TransformPath = %q, want %q", result, "src/main.go")
	}
}

func TestPathMapper_NoMatch(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "other.txt", "content")

	content := []byte(`paths:
  "other.txt": "renamed.txt"
`)
	parsed, err := Parse(content, dir, ".render.yaml")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	mapper := NewPathMapper(parsed)

	result, err := mapper.TransformPath("unmatched.go", nil)
	if err != nil {
		t.Fatalf("TransformPath failed: %v", err)
	}

	if result != "unmatched.go" {
		t.Errorf("TransformPath = %q, want %q", result, "unmatched.go")
	}
}

func TestPathMapper_FileMatch(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "model.go.tmpl", "content")

	content := []byte(`paths:
  "model.go.tmpl": "{{ .name | snakeCase }}.go"
`)
	parsed, err := Parse(content, dir, ".render.yaml")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	mapper := NewPathMapper(parsed)

	data := map[string]any{"name": "UserProfile"}
	result, err := mapper.TransformPath("model.go.tmpl", data)
	if err != nil {
		t.Fatalf("TransformPath failed: %v", err)
	}

	if result != "user_profile.go" {
		t.Errorf("TransformPath = %q, want %q", result, "user_profile.go")
	}
}

func TestPathMapper_DirPrefixMatch(t *testing.T) {
	dir := t.TempDir()
	mkdirAll(t, dir, "server/src/main/java")
	writeFile(t, dir, "server/src/main/java/Service.java", "content")

	content := []byte(`paths:
  "server/src/main/java": "server/src/main/java/{{ .package | replace \".\" \"/\" }}"
`)
	parsed, err := Parse(content, dir, ".render.yaml")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	mapper := NewPathMapper(parsed)

	data := map[string]any{"package": "com.example.service"}
	result, err := mapper.TransformPath("server/src/main/java/Service.java", data)
	if err != nil {
		t.Fatalf("TransformPath failed: %v", err)
	}

	expected := "server/src/main/java/com/example/service/Service.java"
	if result != expected {
		t.Errorf("TransformPath = %q, want %q", result, expected)
	}
}

func TestPathMapper_FilePrecedence(t *testing.T) {
	dir := t.TempDir()
	mkdirAll(t, dir, "src")
	writeFile(t, dir, "src/special.go", "content")

	// File mapping should take precedence over dir prefix
	content := []byte(`paths:
  "src/special.go": "renamed/special.go"
  "src": "other/{{ .name }}"
`)
	parsed, err := Parse(content, dir, ".render.yaml")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	mapper := NewPathMapper(parsed)

	result, err := mapper.TransformPath("src/special.go", nil)
	if err != nil {
		t.Fatalf("TransformPath failed: %v", err)
	}

	if result != "renamed/special.go" {
		t.Errorf("TransformPath = %q, want %q", result, "renamed/special.go")
	}
}

func TestPathMapper_LongestPrefixWins(t *testing.T) {
	dir := t.TempDir()
	mkdirAll(t, dir, "src/main/java")
	writeFile(t, dir, "src/main/java/App.java", "content")

	content := []byte(`paths:
  "src": "pkg1"
  "src/main": "pkg2"
  "src/main/java": "pkg3"
`)
	parsed, err := Parse(content, dir, ".render.yaml")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	mapper := NewPathMapper(parsed)

	result, err := mapper.TransformPath("src/main/java/App.java", nil)
	if err != nil {
		t.Fatalf("TransformPath failed: %v", err)
	}

	// Longest prefix "src/main/java" should win
	if result != "pkg3/App.java" {
		t.Errorf("TransformPath = %q, want %q", result, "pkg3/App.java")
	}
}

func TestPathMapper_TraversalInRenderedPath(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "model.go", "content")

	content := []byte(`paths:
  "model.go": "{{ .path }}"
`)
	parsed, err := Parse(content, dir, ".render.yaml")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	mapper := NewPathMapper(parsed)

	// Attempt to render a path with traversal
	data := map[string]any{"path": "../escaped.go"}
	_, err = mapper.TransformPath("model.go", data)
	if err == nil {
		t.Fatal("Expected error for path traversal in rendered output")
	}
}

func TestNewPathMapper_EmptyConfig(t *testing.T) {
	mapper := NewPathMapper(nil)
	if mapper != nil {
		t.Error("Expected nil mapper for nil config")
	}

	emptyParsed := &ParsedConfig{}
	mapper = NewPathMapper(emptyParsed)
	if mapper != nil {
		t.Error("Expected nil mapper for empty config")
	}
}

func TestPathMapper_CanOverwrite_NilMapper(t *testing.T) {
	var mapper *PathMapper
	if !mapper.CanOverwrite("any/path") {
		t.Error("nil mapper should return true for CanOverwrite")
	}
}

func TestPathMapper_CanOverwrite_Default(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "model.go.tmpl", "content")

	content := []byte(`paths:
  "model.go.tmpl": "output.go"
`)
	parsed, err := Parse(content, dir, ".render.yaml")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	mapper := NewPathMapper(parsed)

	// String value should default to overwrite: true
	if !mapper.CanOverwrite("model.go.tmpl") {
		t.Error("CanOverwrite should return true for path with default overwrite")
	}

	// Unspecified paths should also return true
	if !mapper.CanOverwrite("other/path") {
		t.Error("CanOverwrite should return true for unspecified paths")
	}
}

func TestPathMapper_CanOverwrite_False(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "model.go.tmpl", "content")

	content := []byte(`paths:
  "model.go.tmpl":
    path: "output.go"
    overwrite: false
`)
	parsed, err := Parse(content, dir, ".render.yaml")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	mapper := NewPathMapper(parsed)

	// Explicit overwrite: false should return false
	if mapper.CanOverwrite("model.go.tmpl") {
		t.Error("CanOverwrite should return false for path with overwrite: false")
	}
}

func TestPathMapper_CanOverwrite_ExplicitTrue(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "model.go.tmpl", "content")

	content := []byte(`paths:
  "model.go.tmpl":
    path: "output.go"
    overwrite: true
`)
	parsed, err := Parse(content, dir, ".render.yaml")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	mapper := NewPathMapper(parsed)

	// Explicit overwrite: true should return true
	if !mapper.CanOverwrite("model.go.tmpl") {
		t.Error("CanOverwrite should return true for path with overwrite: true")
	}
}

// Helper function
func mkdirAll(t *testing.T, dir, path string) {
	t.Helper()
	fullPath := filepath.Join(dir, path)
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
}
