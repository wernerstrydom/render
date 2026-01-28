package acceptance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestConfigWithoutFile tests backwards compatibility - no config file.
func TestConfigWithoutFile(t *testing.T) {
	dir := createTempDir(t)

	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	writeFile(t, tmplDir, "config.yaml.tmpl", "name: {{ .name }}")
	writeFile(t, tmplDir, "static.txt", "static content")

	data := writeFile(t, dir, "data.json", `{"name": "TestApp"}`)
	outputDir := filepath.Join(dir, "output")

	stdout, stderr, err := runRender(t, tmplDir, data, "-o", outputDir)
	if err != nil {
		t.Fatalf("render failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	// Verify output
	if !fileExists(filepath.Join(outputDir, "config.yaml")) {
		t.Error("config.yaml not created")
	}
	if !fileExists(filepath.Join(outputDir, "static.txt")) {
		t.Error("static.txt not created")
	}
}

// TestConfigWithPathTransformation tests path transformation via .render.yaml.
func TestConfigWithPathTransformation(t *testing.T) {
	dir := createTempDir(t)

	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	// Create template file
	writeFile(t, tmplDir, "model.go.tmpl", `package {{ .package }}

type {{ .name | pascalCase }} struct {}`)

	// Create config file
	writeFile(t, tmplDir, ".render.yaml", `paths:
  "model.go.tmpl": "{{ .name | snakeCase }}.go"
`)

	data := writeFile(t, dir, "data.json", `{"name": "UserProfile", "package": "models"}`)
	outputDir := filepath.Join(dir, "output")

	stdout, stderr, err := runRender(t, tmplDir, data, "-o", outputDir)
	if err != nil {
		t.Fatalf("render failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	// Verify renamed output
	expectedPath := filepath.Join(outputDir, "user_profile.go")
	if !fileExists(expectedPath) {
		t.Errorf("Expected file %s not created", expectedPath)
	}

	// Verify content
	content := readFile(t, expectedPath)
	if !strings.Contains(content, "type UserProfile struct") {
		t.Errorf("Content missing expected struct: %s", content)
	}
}

// TestConfigWithDirPrefixMapping tests directory prefix mapping.
func TestConfigWithDirPrefixMapping(t *testing.T) {
	dir := createTempDir(t)

	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(filepath.Join(tmplDir, "server", "src", "main", "java"), 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	// Create template file in nested directory
	writeFile(t, tmplDir, "server/src/main/java/ServiceImpl.java.tmpl", `package {{ .package }};

public class {{ .displayName }}ServiceImpl {
    // {{ .displayName }} implementation
}`)

	// Create config with directory prefix mapping
	writeFile(t, tmplDir, ".render.yaml", `paths:
  "server/src/main/java/ServiceImpl.java.tmpl": "server/src/main/java/{{ .displayName }}ServiceImpl.java"
  "server/src/main/java": "server/src/main/java/{{ .package | replace \".\" \"/\" }}"
`)

	data := writeFile(t, dir, "data.json", `{
  "id": "taxonomy",
  "displayName": "Taxonomy",
  "package": "com.example.taxonomy.server"
}`)

	outputDir := filepath.Join(dir, "output")

	stdout, stderr, err := runRender(t, tmplDir, data, "-o", outputDir)
	if err != nil {
		t.Fatalf("render failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	// Verify output path was transformed
	expectedPath := filepath.Join(outputDir, "server", "src", "main", "java",
		"com", "example", "taxonomy", "server", "TaxonomyServiceImpl.java")
	if !fileExists(expectedPath) {
		t.Errorf("Expected file %s not created. Stdout: %s", expectedPath, stdout)
	}

	// Verify content
	if fileExists(expectedPath) {
		content := readFile(t, expectedPath)
		if !strings.Contains(content, "package com.example.taxonomy.server;") {
			t.Errorf("Content missing package: %s", content)
		}
		if !strings.Contains(content, "class TaxonomyServiceImpl") {
			t.Errorf("Content missing class: %s", content)
		}
	}
}

// TestConfigSkipped tests that .render.yaml is not copied to output.
func TestConfigSkipped(t *testing.T) {
	dir := createTempDir(t)

	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	writeFile(t, tmplDir, "file.txt", "content")
	writeFile(t, tmplDir, ".render.yaml", `paths: {}`)
	writeFile(t, tmplDir, ".render.yml", `paths: {}`)

	data := writeFile(t, dir, "data.json", `{}`)
	outputDir := filepath.Join(dir, "output")

	_, _, err := runRender(t, tmplDir, data, "-o", outputDir)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Verify config files are NOT copied
	if fileExists(filepath.Join(outputDir, ".render.yaml")) {
		t.Error(".render.yaml should not be copied to output")
	}
	if fileExists(filepath.Join(outputDir, ".render.yml")) {
		t.Error(".render.yml should not be copied to output")
	}

	// Verify regular file IS copied
	if !fileExists(filepath.Join(outputDir, "file.txt")) {
		t.Error("file.txt should be copied to output")
	}
}

// TestConfigUnknownKey tests error for unknown config keys.
func TestConfigUnknownKey(t *testing.T) {
	dir := createTempDir(t)

	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	writeFile(t, tmplDir, "file.txt", "content")
	writeFile(t, tmplDir, ".render.yaml", `files:
  "file.txt": "output.txt"
`)

	data := writeFile(t, dir, "data.json", `{}`)
	outputDir := filepath.Join(dir, "output")

	_, stderr, err := runRender(t, tmplDir, data, "-o", outputDir)
	if err == nil {
		t.Fatal("Expected error for unknown config key")
	}

	if !strings.Contains(stderr, "unknown key") {
		t.Errorf("Expected 'unknown key' in error: %s", stderr)
	}
}

// TestConfigInvalidTemplate tests error for invalid template syntax.
func TestConfigInvalidTemplate(t *testing.T) {
	dir := createTempDir(t)

	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	writeFile(t, tmplDir, "file.txt", "content")
	writeFile(t, tmplDir, ".render.yaml", `paths:
  "file.txt": "{{ .name | invalid"
`)

	data := writeFile(t, dir, "data.json", `{}`)
	outputDir := filepath.Join(dir, "output")

	_, stderr, err := runRender(t, tmplDir, data, "-o", outputDir)
	if err == nil {
		t.Fatal("Expected error for invalid template syntax")
	}

	if !strings.Contains(stderr, "invalid template syntax") && !strings.Contains(stderr, "template") {
		t.Errorf("Expected template error in stderr: %s", stderr)
	}
}

// TestConfigSourceNotExist tests error for non-existent source.
func TestConfigSourceNotExist(t *testing.T) {
	dir := createTempDir(t)

	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	writeFile(t, tmplDir, ".render.yaml", `paths:
  "missing.txt": "output.txt"
`)

	data := writeFile(t, dir, "data.json", `{}`)
	outputDir := filepath.Join(dir, "output")

	_, stderr, err := runRender(t, tmplDir, data, "-o", outputDir)
	if err == nil {
		t.Fatal("Expected error for non-existent source")
	}

	if !strings.Contains(stderr, "does not exist") {
		t.Errorf("Expected 'does not exist' in error: %s", stderr)
	}
}

// TestExplicitControlFile tests --control flag for explicit control file.
func TestExplicitControlFile(t *testing.T) {
	dir := createTempDir(t)

	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	writeFile(t, tmplDir, "model.go.tmpl", `package main`)

	// Create control file outside template directory
	controlFile := writeFile(t, dir, "custom-render.yaml", `paths:
  "model.go.tmpl": "custom_output.go"
`)

	data := writeFile(t, dir, "data.json", `{}`)
	outputDir := filepath.Join(dir, "output")

	_, _, err := runRender(t, tmplDir, data, "-o", outputDir, "--control", controlFile)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Verify renamed output using explicit control file
	expectedPath := filepath.Join(outputDir, "custom_output.go")
	if !fileExists(expectedPath) {
		t.Errorf("Expected file %s not created", expectedPath)
	}
}

// TestConfigOverwriteFalse tests that overwrite: false preserves existing files.
func TestConfigOverwriteFalse(t *testing.T) {
	dir := createTempDir(t)

	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	// Create template files
	writeFile(t, tmplDir, "regular.txt.tmpl", "Regular content: {{ .name }}")
	writeFile(t, tmplDir, "protected.txt.tmpl", "Protected content: {{ .name }}")

	// Create config with overwrite: false for protected file
	writeFile(t, tmplDir, ".render.yaml", `paths:
  "regular.txt.tmpl": "regular.txt"
  "protected.txt.tmpl":
    path: "protected.txt"
    overwrite: false
`)

	data := writeFile(t, dir, "data.json", `{"name": "TestApp"}`)
	outputDir := filepath.Join(dir, "output")

	// First render - both files created
	stdout, stderr, err := runRender(t, tmplDir, data, "-o", outputDir)
	if err != nil {
		t.Fatalf("First render failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	// Verify both files created
	regularPath := filepath.Join(outputDir, "regular.txt")
	protectedPath := filepath.Join(outputDir, "protected.txt")

	if !fileExists(regularPath) {
		t.Fatal("regular.txt not created")
	}
	if !fileExists(protectedPath) {
		t.Fatal("protected.txt not created")
	}

	// Modify the protected file (simulating user customization)
	if err := os.WriteFile(protectedPath, []byte("User customized content"), 0644); err != nil {
		t.Fatalf("Failed to write customized content: %v", err)
	}

	// Second render with --force - protected file should be preserved
	stdout, stderr, err = runRender(t, tmplDir, data, "-o", outputDir, "--force")
	if err != nil {
		t.Fatalf("Second render failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	// Verify protected file content unchanged
	protectedContent := readFile(t, protectedPath)
	if protectedContent != "User customized content" {
		t.Errorf("Protected file was overwritten. Content = %q, want %q", protectedContent, "User customized content")
	}

	// Verify stdout mentions skipped (case-insensitive check)
	lowerStdout := strings.ToLower(stdout)
	if !strings.Contains(lowerStdout, "skipped") || !strings.Contains(lowerStdout, "no-overwrite") {
		t.Errorf("Output should mention skipped file: %s", stdout)
	}
}

// TestConfigOverwriteFalseDryRun tests dry-run output for overwrite: false files.
func TestConfigOverwriteFalseDryRun(t *testing.T) {
	dir := createTempDir(t)

	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	writeFile(t, tmplDir, "protected.txt.tmpl", "Protected content")
	writeFile(t, tmplDir, ".render.yaml", `paths:
  "protected.txt.tmpl":
    path: "protected.txt"
    overwrite: false
`)

	data := writeFile(t, dir, "data.json", `{}`)
	outputDir := filepath.Join(dir, "output")

	// Create existing file
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}
	writeFile(t, outputDir, "protected.txt", "Existing content")

	// Dry run should show skip action
	stdout, stderr, err := runRender(t, tmplDir, data, "-o", outputDir, "--dry-run")
	if err != nil {
		t.Fatalf("Dry run failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	if !strings.Contains(stdout, "skip") && !strings.Contains(stdout, "no-overwrite") {
		t.Errorf("Dry run should show skip action: %s", stdout)
	}
}

// TestConfigOverwriteFalseJSON tests JSON output for overwrite: false files.
func TestConfigOverwriteFalseJSON(t *testing.T) {
	dir := createTempDir(t)

	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	writeFile(t, tmplDir, "protected.txt.tmpl", "Protected content")
	writeFile(t, tmplDir, ".render.yaml", `paths:
  "protected.txt.tmpl":
    path: "protected.txt"
    overwrite: false
`)

	data := writeFile(t, dir, "data.json", `{}`)
	outputDir := filepath.Join(dir, "output")

	// Create existing file
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}
	writeFile(t, outputDir, "protected.txt", "Existing content")

	// Render with --force --json should show skipped in JSON
	stdout, stderr, err := runRender(t, tmplDir, data, "-o", outputDir, "--force", "--json")
	if err != nil {
		t.Fatalf("JSON render failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	if !strings.Contains(stdout, "skipped") || !strings.Contains(stdout, "no-overwrite") {
		t.Errorf("JSON output should show skipped action: %s", stdout)
	}
}

// TestConfigOverwriteFalseNewFile tests that overwrite: false still creates new files.
func TestConfigOverwriteFalseNewFile(t *testing.T) {
	dir := createTempDir(t)

	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	writeFile(t, tmplDir, "protected.txt.tmpl", "Protected content: {{ .name }}")
	writeFile(t, tmplDir, ".render.yaml", `paths:
  "protected.txt.tmpl":
    path: "protected.txt"
    overwrite: false
`)

	data := writeFile(t, dir, "data.json", `{"name": "TestApp"}`)
	outputDir := filepath.Join(dir, "output")

	// Render - file doesn't exist, should be created
	stdout, stderr, err := runRender(t, tmplDir, data, "-o", outputDir)
	if err != nil {
		t.Fatalf("Render failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	protectedPath := filepath.Join(outputDir, "protected.txt")
	if !fileExists(protectedPath) {
		t.Fatal("Protected file should be created when it doesn't exist")
	}

	content := readFile(t, protectedPath)
	if !strings.Contains(content, "Protected content: TestApp") {
		t.Errorf("Content = %q, want to contain 'Protected content: TestApp'", content)
	}
}

// TestConfigMixedPathFormats tests mixing string and object path formats.
func TestConfigMixedPathFormats(t *testing.T) {
	dir := createTempDir(t)

	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	writeFile(t, tmplDir, "string-format.txt.tmpl", "String format")
	writeFile(t, tmplDir, "object-format.txt.tmpl", "Object format")
	writeFile(t, tmplDir, "explicit-true.txt.tmpl", "Explicit true")

	// Mix string and object formats in config
	writeFile(t, tmplDir, ".render.yaml", `paths:
  "string-format.txt.tmpl": "string-format.txt"
  "object-format.txt.tmpl":
    path: "object-format.txt"
    overwrite: false
  "explicit-true.txt.tmpl":
    path: "explicit-true.txt"
    overwrite: true
`)

	data := writeFile(t, dir, "data.json", `{}`)
	outputDir := filepath.Join(dir, "output")

	stdout, stderr, err := runRender(t, tmplDir, data, "-o", outputDir)
	if err != nil {
		t.Fatalf("Render failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	// Verify all files created
	for _, name := range []string{"string-format.txt", "object-format.txt", "explicit-true.txt"} {
		path := filepath.Join(outputDir, name)
		if !fileExists(path) {
			t.Errorf("File %s not created", name)
		}
	}
}
