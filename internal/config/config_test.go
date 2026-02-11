package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParse_ValidConfig(t *testing.T) {
	// Create temp directory with test files
	dir := t.TempDir()
	writeFile(t, dir, "model.go.tmpl", "package {{ .package }}")
	mkdir(t, dir, "src")
	writeFile(t, dir, "src/main.go", "package main")

	content := []byte(`paths:
  "model.go.tmpl": "{{ .name | snakeCase }}.go"
  "src": "pkg/{{ .name }}"
`)

	parsed, err := Parse(content, dir, ".render.yaml")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(parsed.fileTemplates) != 1 {
		t.Errorf("Expected 1 file template, got %d", len(parsed.fileTemplates))
	}

	if len(parsed.dirMappings) != 1 {
		t.Errorf("Expected 1 dir mapping, got %d", len(parsed.dirMappings))
	}
}

func TestParse_EmptyConfig(t *testing.T) {
	dir := t.TempDir()
	content := []byte(`paths: {}`)

	parsed, err := Parse(content, dir, ".render.yaml")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if !parsed.IsEmpty() {
		t.Error("Expected empty config")
	}
}

func TestParse_UnknownKey(t *testing.T) {
	dir := t.TempDir()
	content := []byte(`files:
  "model.go": "output.go"
`)

	_, err := Parse(content, dir, ".render.yaml")
	if err == nil {
		t.Fatal("Expected error for unknown key")
	}

	if want := `unknown key "files"`; !containsString(err.Error(), want) {
		t.Errorf("Error %q should contain %q", err.Error(), want)
	}
}

func TestParse_SourceNotExist(t *testing.T) {
	dir := t.TempDir()
	content := []byte(`paths:
  "missing.go": "output.go"
`)

	_, err := Parse(content, dir, ".render.yaml")
	if err == nil {
		t.Fatal("Expected error for missing source")
	}

	if want := "does not exist"; !containsString(err.Error(), want) {
		t.Errorf("Error %q should contain %q", err.Error(), want)
	}
}

func TestParse_InvalidTemplateSyntax(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "model.go", "content")

	content := []byte(`paths:
  "model.go": "{{ .name | invalid"
`)

	_, err := Parse(content, dir, ".render.yaml")
	if err == nil {
		t.Fatal("Expected error for invalid template")
	}

	if want := "invalid template syntax"; !containsString(err.Error(), want) {
		t.Errorf("Error %q should contain %q", err.Error(), want)
	}
}

func TestValidateSourcePath_AbsolutePath(t *testing.T) {
	// Use filepath.Abs to get a cross-platform absolute path
	absPath, err := filepath.Abs("somepath")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	err = validateSourcePath(absPath)
	if err == nil {
		t.Fatal("Expected error for absolute path")
	}

	if want := "must be relative"; !containsString(err.Error(), want) {
		t.Errorf("Error %q should contain %q", err.Error(), want)
	}
}

func TestValidateSourcePath_Traversal(t *testing.T) {
	err := validateSourcePath("../secret.txt")
	if err == nil {
		t.Fatal("Expected error for path traversal")
	}

	if want := "directory traversal"; !containsString(err.Error(), want) {
		t.Errorf("Error %q should contain %q", err.Error(), want)
	}
}

func TestValidateSourcePath_NullByte(t *testing.T) {
	err := validateSourcePath("file\x00.txt")
	if err == nil {
		t.Fatal("Expected error for null byte")
	}

	if want := "null byte"; !containsString(err.Error(), want) {
		t.Errorf("Error %q should contain %q", err.Error(), want)
	}
}

func TestValidateSourcePath_Valid(t *testing.T) {
	tests := []string{
		"file.txt",
		"dir/file.txt",
		"a/b/c/d.go",
		"src/main/java",
	}

	for _, path := range tests {
		err := validateSourcePath(path)
		if err != nil {
			t.Errorf("validateSourcePath(%q) = %v, want nil", path, err)
		}
	}
}

func TestValidateRenderedPath_Traversal(t *testing.T) {
	err := ValidateRenderedPath("../escaped")
	if err == nil {
		t.Fatal("Expected error for path traversal")
	}

	if want := "directory traversal"; !containsString(err.Error(), want) {
		t.Errorf("Error %q should contain %q", err.Error(), want)
	}
}

func TestValidateRenderedPath_Valid(t *testing.T) {
	tests := []string{
		"output.txt",
		"pkg/model/user.go",
		"com/example/service",
	}

	for _, path := range tests {
		err := ValidateRenderedPath(path)
		if err != nil {
			t.Errorf("ValidateRenderedPath(%q) = %v, want nil", path, err)
		}
	}
}

func TestLoad_NoConfig(t *testing.T) {
	dir := t.TempDir()

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg != nil {
		t.Error("Expected nil config when no config file exists")
	}
}

func TestLoad_YAMLConfig(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "model.go.tmpl", "content")
	writeFile(t, dir, ".render.yaml", `paths:
  "model.go.tmpl": "{{ .name }}.go"
`)

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("Expected config to be loaded")
	}

	if !cfg.HasFileMappings() {
		t.Error("Expected file mappings")
	}
}

func TestLoad_YMLConfig(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "model.go.tmpl", "content")
	writeFile(t, dir, ".render.yml", `paths:
  "model.go.tmpl": "{{ .name }}.go"
`)

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("Expected config to be loaded")
	}
}

func TestParsedConfig_IsEmpty(t *testing.T) {
	var nilConfig *ParsedConfig
	if !nilConfig.IsEmpty() {
		t.Error("nil config should be empty")
	}

	emptyConfig := &ParsedConfig{}
	if !emptyConfig.IsEmpty() {
		t.Error("empty config should be empty")
	}
}

func TestShouldSkipConfigFile(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{".render.yaml", true},
		{".render.yml", true},
		{"render.json", true},
		{"other.yaml", false},
		{"config.yaml", false},
		{"src/.render.yaml", false}, // Only top-level
	}

	for _, tt := range tests {
		got := ShouldSkipConfigFile(tt.path)
		if got != tt.want {
			t.Errorf("ShouldSkipConfigFile(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestParse_PathMappingStringValue(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "model.go.tmpl", "content")

	content := []byte(`paths:
  "model.go.tmpl": "{{ .name }}.go"
`)

	parsed, err := Parse(content, dir, ".render.yaml")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(parsed.fileTemplates) != 1 {
		t.Errorf("Expected 1 file template, got %d", len(parsed.fileTemplates))
	}

	// String values should default to overwrite: true (not in noOverwrite set)
	if len(parsed.noOverwrite) != 0 {
		t.Errorf("Expected 0 noOverwrite entries, got %d", len(parsed.noOverwrite))
	}
}

func TestParse_PathMappingObjectValue(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "model.go.tmpl", "content")

	content := []byte(`paths:
  "model.go.tmpl":
    path: "{{ .name }}.go"
    overwrite: false
`)

	parsed, err := Parse(content, dir, ".render.yaml")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(parsed.fileTemplates) != 1 {
		t.Errorf("Expected 1 file template, got %d", len(parsed.fileTemplates))
	}

	// Should have noOverwrite set for this path
	if !parsed.noOverwrite["model.go.tmpl"] {
		t.Error("Expected model.go.tmpl in noOverwrite set")
	}
}

func TestParse_PathMappingMixedValues(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "model.go.tmpl", "content")
	writeFile(t, dir, "handler.go.tmpl", "content")
	writeFile(t, dir, "config.yaml.tmpl", "content")

	content := []byte(`paths:
  "model.go.tmpl": "{{ .name }}.go"
  "handler.go.tmpl":
    path: "internal/handler.go"
    overwrite: false
  "config.yaml.tmpl":
    path: "config.yaml"
    overwrite: true
`)

	parsed, err := Parse(content, dir, ".render.yaml")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(parsed.fileTemplates) != 3 {
		t.Errorf("Expected 3 file templates, got %d", len(parsed.fileTemplates))
	}

	// Only handler.go.tmpl should be in noOverwrite set
	if len(parsed.noOverwrite) != 1 {
		t.Errorf("Expected 1 noOverwrite entry, got %d", len(parsed.noOverwrite))
	}

	if !parsed.noOverwrite["handler.go.tmpl"] {
		t.Error("Expected handler.go.tmpl in noOverwrite set")
	}

	if parsed.noOverwrite["model.go.tmpl"] {
		t.Error("model.go.tmpl should not be in noOverwrite set")
	}

	if parsed.noOverwrite["config.yaml.tmpl"] {
		t.Error("config.yaml.tmpl should not be in noOverwrite set (overwrite: true)")
	}
}

func TestParse_PathMappingObjectMissingPath(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "model.go.tmpl", "content")

	content := []byte(`paths:
  "model.go.tmpl":
    overwrite: false
`)

	_, err := Parse(content, dir, ".render.yaml")
	if err == nil {
		t.Fatal("Expected error for object without path field")
	}

	if want := "must have 'path' field"; !containsString(err.Error(), want) {
		t.Errorf("Error %q should contain %q", err.Error(), want)
	}
}

// Helper functions

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}
}

func mkdir(t *testing.T, dir, name string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStringHelper(s, substr))
}

func containsStringHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
