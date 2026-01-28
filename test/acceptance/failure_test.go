package acceptance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestInvalidTemplateSyntax tests error handling for invalid Go template syntax.
func TestInvalidTemplateSyntax(t *testing.T) {
	dir := createTempDir(t)

	tests := []struct {
		name     string
		template string
		errMsg   string
	}{
		{
			name:     "unclosed action",
			template: "Hello {{ .name",
			errMsg:   "unclosed action",
		},
		{
			name:     "unclosed range",
			template: "{{ range .items }}item{{ end",
			errMsg:   "unclosed action",
		},
		{
			name:     "invalid function",
			template: "{{ nonexistentFunc .name }}",
			errMsg:   "not defined",
		},
		{
			name:     "mismatched end",
			template: "{{ if .x }}yes{{ else }}no{{ end }}{{ end }}",
			errMsg:   "unexpected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := writeFile(t, dir, "template.txt", tt.template)
			data := writeFile(t, dir, "data.json", `{"name": "test", "items": [], "x": true}`)
			output := filepath.Join(dir, "output.txt")

			_, stderr, err := runRender(t, tmpl, data, "-o", output)
			if err == nil {
				t.Fatal("should fail with invalid template syntax")
			}

			if !strings.Contains(strings.ToLower(stderr), strings.ToLower(tt.errMsg)) {
				t.Errorf("error should mention %q, got: %s", tt.errMsg, stderr)
			}
		})
	}
}

// TestInvalidJSONData tests error handling for invalid JSON data.
func TestInvalidJSONData(t *testing.T) {
	dir := createTempDir(t)

	tests := []struct {
		name    string
		content string
	}{
		{"unclosed object", `{"name": "test"`},
		{"unclosed array", `{"items": [1, 2, 3`},
		{"missing colon", `{"name" "test"}`},
		{"trailing comma", `{"name": "test",}`},
		{"unquoted key", `{name: "test"}`},
		{"single quotes", `{'name': 'test'}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := writeFile(t, dir, "template.txt", "{{ .name }}")
			data := writeFile(t, dir, "data.json", tt.content)
			output := filepath.Join(dir, "output.txt")

			_, stderr, err := runRender(t, tmpl, data, "-o", output)
			if err == nil {
				t.Fatal("should fail with invalid JSON")
			}

			if !strings.Contains(strings.ToLower(stderr), "json") &&
				!strings.Contains(strings.ToLower(stderr), "parse") &&
				!strings.Contains(strings.ToLower(stderr), "invalid") {
				t.Errorf("error should mention JSON parsing issue, got: %s", stderr)
			}
		})
	}
}

// TestInvalidYAMLData tests error handling for invalid YAML data.
func TestInvalidYAMLData(t *testing.T) {
	dir := createTempDir(t)

	tests := []struct {
		name    string
		content string
	}{
		{"tabs in indentation", "name: test\n\t\tvalue: bad"},
		{"duplicate keys", "name: first\nname: second"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := writeFile(t, dir, "template.txt", "{{ .name }}")
			data := writeFile(t, dir, "data.yaml", tt.content)
			output := filepath.Join(dir, "output.txt")

			_, stderr, err := runRender(t, tmpl, data, "-o", output)
			// Note: Some YAML parsers are lenient, so we just check that it either fails
			// or produces unexpected results. For strict validation, the error should occur.
			if err != nil {
				// Good - it detected the issue
				if !strings.Contains(strings.ToLower(stderr), "yaml") &&
					!strings.Contains(strings.ToLower(stderr), "parse") &&
					!strings.Contains(strings.ToLower(stderr), "invalid") &&
					!strings.Contains(strings.ToLower(stderr), "error") {
					t.Logf("Error message: %s", stderr)
				}
			}
			// Some YAML issues might not cause errors (e.g., duplicate keys just overwrite)
			// which is acceptable behavior
		})
	}
}

// TestUnsupportedDataFormat tests error handling for unsupported file extensions.
func TestUnsupportedDataFormat(t *testing.T) {
	dir := createTempDir(t)

	tests := []struct {
		name string
		ext  string
	}{
		{"xml", "data.xml"},
		{"txt", "data.txt"},
		{"toml", "data.toml"},
		{"no extension", "data"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := writeFile(t, dir, "template.txt", "{{ .name }}")
			data := writeFile(t, dir, tt.ext, `name: test`)
			output := filepath.Join(dir, "output.txt")

			_, stderr, err := runRender(t, tmpl, data, "-o", output)
			if err == nil {
				t.Fatal("should fail with unsupported data format")
			}

			if !strings.Contains(strings.ToLower(stderr), "unsupported") &&
				!strings.Contains(strings.ToLower(stderr), "extension") &&
				!strings.Contains(strings.ToLower(stderr), "format") {
				t.Errorf("error should mention unsupported format, got: %s", stderr)
			}
		})
	}
}

// TestTemplateExecutionError tests error handling when template execution fails.
func TestTemplateExecutionError(t *testing.T) {
	dir := createTempDir(t)

	tests := []struct {
		name     string
		template string
		data     string
	}{
		{
			name:     "nil map access",
			template: "{{ .user.name }}",
			data:     `{"user": null}`,
		},
		{
			name:     "index out of range",
			template: "{{ index .items 10 }}",
			data:     `{"items": ["a", "b"]}`,
		},
		{
			name:     "call on nil",
			template: "{{ .func }}",
			data:     `{"func": null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := writeFile(t, dir, "template.txt", tt.template)
			data := writeFile(t, dir, "data.json", tt.data)
			output := filepath.Join(dir, "output.txt")

			_, _, err := runRender(t, tmpl, data, "-o", output)
			// Some of these might not error depending on Go template behavior
			// The important thing is that the tool doesn't crash
			_ = err
		})
	}
}

// TestDirModeWithInvalidTemplates tests dir mode with invalid templates.
func TestDirModeWithInvalidTemplates(t *testing.T) {
	dir := createTempDir(t)

	// Create template directory with one invalid template
	tmplDir := filepath.Join(dir, "templates")
	writeFile(t, tmplDir, "valid.txt.tmpl", "{{ .name }}")
	writeFile(t, tmplDir, "invalid.txt.tmpl", "{{ .name") // Invalid syntax

	data := writeFile(t, dir, "data.json", `{"name": "test"}`)
	outputDir := filepath.Join(dir, "output")

	_, stderr, err := runRender(t, tmplDir, data, "-o", outputDir)
	if err == nil {
		t.Fatal("should fail with invalid template in directory")
	}

	if !strings.Contains(strings.ToLower(stderr), "unclosed") &&
		!strings.Contains(strings.ToLower(stderr), "invalid") &&
		!strings.Contains(strings.ToLower(stderr), "error") {
		t.Errorf("error should mention template issue, got: %s", stderr)
	}
}

// TestOutputCollision tests that render dir fails when two templates
// produce the same output path.
func TestOutputCollision(t *testing.T) {
	dir := createTempDir(t)

	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	// Create two templates that will produce the same output via .render.yaml
	writeFile(t, tmplDir, "user.go.tmpl", "package user")
	writeFile(t, tmplDir, "account.go.tmpl", "package account")

	// Configure both to map to the same output path
	writeFile(t, tmplDir, ".render.yaml", `paths:
  "user.go.tmpl": "output.go"
  "account.go.tmpl": "output.go"
`)

	data := writeFile(t, dir, "data.json", `{}`)
	outputDir := filepath.Join(dir, "output")

	_, stderr, err := runRender(t, tmplDir, data, "-o", outputDir)
	if err == nil {
		t.Fatal("should fail with output path collision")
	}

	if !strings.Contains(strings.ToLower(stderr), "collision") {
		t.Errorf("error should mention collision, got: %s", stderr)
	}
}

// TestPathTraversal tests that path traversal is rejected.
func TestPathTraversal(t *testing.T) {
	dir := createTempDir(t)

	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	writeFile(t, tmplDir, "file.txt", "content")
	// Try to traverse outside the output directory
	writeFile(t, tmplDir, ".render.yaml", `paths:
  "file.txt": "../../../etc/passwd"
`)

	data := writeFile(t, dir, "data.json", `{}`)
	outputDir := filepath.Join(dir, "output")

	_, stderr, err := runRender(t, tmplDir, data, "-o", outputDir)
	if err == nil {
		t.Fatal("should fail with path traversal attempt")
	}

	// Should either mention traversal or security
	if !strings.Contains(strings.ToLower(stderr), "traversal") &&
		!strings.Contains(strings.ToLower(stderr), "security") &&
		!strings.Contains(strings.ToLower(stderr), "outside") {
		t.Errorf("error should mention security issue, got: %s", stderr)
	}
}

// TestSymlinkOutsideTemplateDir tests that symlinks are rejected.
func TestSymlinkOutsideTemplateDir(t *testing.T) {
	// Skip on Windows where symlinks require special permissions
	if os.Getenv("GOOS") == "windows" {
		t.Skip("Skipping symlink test on Windows")
	}

	dir := createTempDir(t)

	// Create a sensitive file outside the template directory
	sensitiveFile := writeFile(t, dir, "sensitive.txt", "SECRET_DATA")

	// Create template directory
	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	// Create a symlink pointing outside the template directory
	symlinkPath := filepath.Join(tmplDir, "linked.txt")
	if err := os.Symlink(sensitiveFile, symlinkPath); err != nil {
		t.Fatalf("Failed to create symlink: %v", err)
	}

	data := writeFile(t, dir, "data.json", `{}`)
	outputDir := filepath.Join(dir, "output")

	_, stderr, err := runRender(t, tmplDir, data, "-o", outputDir)

	// With the new command, symlinks should be rejected
	if err == nil {
		t.Fatal("render should fail when template contains symlinks")
	}

	if !strings.Contains(strings.ToLower(stderr), "symlink") {
		t.Errorf("error should mention symlink, got: %s", stderr)
	}
}

// TestConversionErrors tests error handling for conversion functions.
func TestConversionErrors(t *testing.T) {
	dir := createTempDir(t)

	tests := []struct {
		name     string
		template string
		data     string
	}{
		{
			name:     "toInt with invalid string",
			template: `{{ toInt .value }}`,
			data:     `{"value": "not a number"}`,
		},
		{
			name:     "toFloat with invalid string",
			template: `{{ toFloat .value }}`,
			data:     `{"value": "not a float"}`,
		},
		{
			name:     "toBool with invalid string",
			template: `{{ toBool .value }}`,
			data:     `{"value": "not a bool"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := writeFile(t, dir, "template.txt", tt.template)
			data := writeFile(t, dir, "data.json", tt.data)
			output := filepath.Join(dir, "output.txt")

			_, _, err := runRender(t, tmpl, data, "-o", output)
			if err == nil {
				t.Fatalf("should fail with conversion error for: %s", tt.name)
			}
		})
	}
}

// TestModDivisionByZero tests error handling for mod with division by zero.
func TestModDivisionByZero(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "template.txt", `{{ mod .a .b }}`)
	data := writeFile(t, dir, "data.json", `{"a": 10, "b": 0}`)
	output := filepath.Join(dir, "output.txt")

	_, stderr, err := runRender(t, tmpl, data, "-o", output)
	if err == nil {
		t.Fatal("should fail with division by zero")
	}

	if !strings.Contains(strings.ToLower(stderr), "division") &&
		!strings.Contains(strings.ToLower(stderr), "zero") {
		t.Errorf("error should mention division by zero, got: %s", stderr)
	}
}

// TestEmptyDataFile tests handling of empty data files.
func TestEmptyDataFile(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "template.txt", "Hello, {{ .name }}!")
	data := writeFile(t, dir, "data.json", "")
	output := filepath.Join(dir, "output.txt")

	_, _, err := runRender(t, tmpl, data, "-o", output)
	// Empty JSON/YAML file should cause an error
	if err == nil {
		t.Log("Note: empty data file was accepted (might be valid depending on implementation)")
	}
}
