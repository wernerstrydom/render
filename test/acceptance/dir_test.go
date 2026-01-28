package acceptance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDirBasic tests basic directory rendering.
func TestDirBasic(t *testing.T) {
	dir := createTempDir(t)

	// Create template directory structure
	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	writeFile(t, tmplDir, "config.yaml.tmpl", `name: {{ .name }}
version: {{ .version }}`)

	writeFile(t, tmplDir, "readme.md.tmpl", `# {{ .name }}

Version: {{ .version }}`)

	// Create data file
	data := writeFile(t, dir, "data.json", `{
		"name": "MyProject",
		"version": "1.0.0"
	}`)

	// Output directory
	outputDir := filepath.Join(dir, "output")

	// Run render
	stdout, stderr, err := runRender(t, tmplDir, data, "-o", outputDir)

	if err != nil {
		t.Fatalf("render failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	// Verify config.yaml (note: .tmpl extension stripped)
	configPath := filepath.Join(outputDir, "config.yaml")
	if !fileExists(configPath) {
		t.Errorf("config.yaml not created")
	} else {
		content := readFile(t, configPath)
		if !strings.Contains(content, "name: MyProject") {
			t.Errorf("config.yaml missing name: %s", content)
		}
		if !strings.Contains(content, "version: 1.0.0") {
			t.Errorf("config.yaml missing version: %s", content)
		}
	}

	// Verify readme.md
	readmePath := filepath.Join(outputDir, "readme.md")
	if !fileExists(readmePath) {
		t.Errorf("readme.md not created")
	} else {
		content := readFile(t, readmePath)
		if !strings.Contains(content, "# MyProject") {
			t.Errorf("readme.md missing title: %s", content)
		}
	}
}

// TestDirCopiesNonTemplates tests that non-.tmpl files are copied verbatim.
func TestDirCopiesNonTemplates(t *testing.T) {
	dir := createTempDir(t)

	// Create template directory with mixed files
	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	// Template file
	writeFile(t, tmplDir, "config.txt.tmpl", "name={{ .name }}")

	// Non-template files (should be copied verbatim)
	writeFile(t, tmplDir, "static.txt", "This is {{ .name }} static content")
	writeFile(t, tmplDir, "binary.bin", "\x00\x01\x02\x03")

	data := writeFile(t, dir, "data.json", `{"name": "Test"}`)
	outputDir := filepath.Join(dir, "output")

	_, _, err := runRender(t, tmplDir, data, "-o", outputDir)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Verify template was rendered
	configContent := readFile(t, filepath.Join(outputDir, "config.txt"))
	if configContent != "name=Test" {
		t.Errorf("Template not rendered: %q", configContent)
	}

	// Verify static file was copied verbatim (not rendered)
	staticContent := readFile(t, filepath.Join(outputDir, "static.txt"))
	if staticContent != "This is {{ .name }} static content" {
		t.Errorf("Static file was modified: %q", staticContent)
	}

	// Verify binary file was copied
	binaryContent := readFile(t, filepath.Join(outputDir, "binary.bin"))
	if binaryContent != "\x00\x01\x02\x03" {
		t.Errorf("Binary file was modified")
	}
}

// TestDirPreservesStructure tests that directory structure is preserved.
func TestDirPreservesStructure(t *testing.T) {
	dir := createTempDir(t)

	// Create nested template directory structure
	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(filepath.Join(tmplDir, "src", "components"), 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(tmplDir, "config"), 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	writeFile(t, tmplDir, "root.txt.tmpl", "root: {{ .name }}")
	writeFile(t, tmplDir, "src/main.go.tmpl", "package {{ .package }}")
	writeFile(t, tmplDir, "src/components/button.tsx.tmpl", "// {{ .name }} button")
	writeFile(t, tmplDir, "config/settings.json.tmpl", `{"app": "{{ .name }}"}`)

	data := writeFile(t, dir, "data.json", `{
		"name": "MyApp",
		"package": "main"
	}`)

	outputDir := filepath.Join(dir, "output")

	_, _, err := runRender(t, tmplDir, data, "-o", outputDir)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Verify all files exist in correct locations
	expectedFiles := []string{
		"root.txt",
		"src/main.go",
		"src/components/button.tsx",
		"config/settings.json",
	}

	for _, f := range expectedFiles {
		path := filepath.Join(outputDir, f)
		if !fileExists(path) {
			t.Errorf("File not created: %s", f)
		}
	}

	// Verify content of nested file
	mainContent := readFile(t, filepath.Join(outputDir, "src", "main.go"))
	if mainContent != "package main" {
		t.Errorf("Nested template not rendered: %q", mainContent)
	}
}

// TestDirWithYAML tests dir rendering with YAML data.
func TestDirWithYAML(t *testing.T) {
	dir := createTempDir(t)

	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	writeFile(t, tmplDir, "output.txt.tmpl", "{{ .app.name }} v{{ .app.version }}")

	data := writeFile(t, dir, "data.yaml", `app:
  name: TestApp
  version: 2.0.0`)

	outputDir := filepath.Join(dir, "output")

	_, _, err := runRender(t, tmplDir, data, "-o", outputDir)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	content := readFile(t, filepath.Join(outputDir, "output.txt"))
	if content != "TestApp v2.0.0" {
		t.Errorf("Output = %q, want %q", content, "TestApp v2.0.0")
	}
}

// TestDirOverwriteProtection tests overwrite protection for directories.
func TestDirOverwriteProtection(t *testing.T) {
	dir := createTempDir(t)

	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}
	writeFile(t, tmplDir, "file.txt.tmpl", "new content")

	data := writeFile(t, dir, "data.json", `{}`)

	outputDir := filepath.Join(dir, "output")
	writeFile(t, outputDir, "file.txt", "original content")

	// Run without --force
	_, _, err := runRender(t, tmplDir, data, "-o", outputDir)
	if err == nil {
		t.Fatal("render should fail without --force when files exist")
	}

	// Verify original unchanged
	content := readFile(t, filepath.Join(outputDir, "file.txt"))
	if content != "original content" {
		t.Errorf("File was modified without --force")
	}
}

// TestDirForceOverwrite tests --force with directories.
func TestDirForceOverwrite(t *testing.T) {
	dir := createTempDir(t)

	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}
	writeFile(t, tmplDir, "file.txt.tmpl", "new content")

	data := writeFile(t, dir, "data.json", `{}`)

	outputDir := filepath.Join(dir, "output")
	writeFile(t, outputDir, "file.txt", "original content")

	// Run with --force
	_, _, err := runRender(t, tmplDir, data, "-o", outputDir, "--force")
	if err != nil {
		t.Fatalf("render --force failed: %v", err)
	}

	content := readFile(t, filepath.Join(outputDir, "file.txt"))
	if content != "new content" {
		t.Errorf("File content = %q, want %q", content, "new content")
	}
}

// TestDirNotDirectory tests error when template is not a directory.
func TestDirNotDirectory(t *testing.T) {
	dir := createTempDir(t)

	// Create a file, not a directory - this will be treated as file mode
	tmpl := writeFile(t, dir, "notdir.txt", "content")
	data := writeFile(t, dir, "data.json", `{}`)
	output := filepath.Join(dir, "output.txt")

	// This will work as file mode since template is a file
	_, _, err := runRender(t, tmpl, data, "-o", output)
	if err != nil {
		t.Fatalf("render should work with file template: %v", err)
	}
}

// TestDirEmptyDirectory tests rendering an empty template directory.
func TestDirEmptyDirectory(t *testing.T) {
	dir := createTempDir(t)

	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	data := writeFile(t, dir, "data.json", `{}`)
	outputDir := filepath.Join(dir, "output")

	// Should succeed even with empty directory
	_, _, err := runRender(t, tmplDir, data, "-o", outputDir)
	if err != nil {
		t.Fatalf("render failed with empty directory: %v", err)
	}
}

// TestIterativeDirectoryRendering tests each mode with directory templates.
func TestIterativeDirectoryRendering(t *testing.T) {
	dir := createTempDir(t)

	// Create template directory
	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	writeFile(t, tmplDir, "config.txt.tmpl", "name={{ .name }}")
	writeFile(t, tmplDir, "static.txt", "static")

	data := writeFile(t, dir, "data.json", `[{"id": "1", "name": "Alice"}, {"id": "2", "name": "Bob"}]`)

	// Dynamic output path for each mode
	outputPattern := filepath.Join(dir, "output-{{.id}}")

	stdout, stderr, err := runRender(t, tmplDir, data, "-o", outputPattern)
	if err != nil {
		t.Fatalf("render failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	// Verify directories were created
	for _, id := range []string{"1", "2"} {
		outDir := filepath.Join(dir, "output-"+id)
		if !fileExists(filepath.Join(outDir, "config.txt")) {
			t.Errorf("config.txt not created in %s", outDir)
		}
		if !fileExists(filepath.Join(outDir, "static.txt")) {
			t.Errorf("static.txt not created in %s", outDir)
		}
	}

	// Verify content
	content1 := readFile(t, filepath.Join(dir, "output-1", "config.txt"))
	if content1 != "name=Alice" {
		t.Errorf("Content = %q, want %q", content1, "name=Alice")
	}
	content2 := readFile(t, filepath.Join(dir, "output-2", "config.txt"))
	if content2 != "name=Bob" {
		t.Errorf("Content = %q, want %q", content2, "name=Bob")
	}
}

// TestDirWithDryRun tests dry-run mode with directories.
func TestDirWithDryRun(t *testing.T) {
	dir := createTempDir(t)

	tmplDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}
	writeFile(t, tmplDir, "file.txt.tmpl", "content")

	data := writeFile(t, dir, "data.json", `{}`)
	outputDir := filepath.Join(dir, "output")

	stdout, _, err := runRender(t, tmplDir, data, "-o", outputDir, "--dry-run")
	if err != nil {
		t.Fatalf("render --dry-run failed: %v", err)
	}

	// Output directory should NOT be created
	if fileExists(outputDir) {
		t.Error("Output directory was created in dry-run mode")
	}

	if !strings.Contains(stdout, "Dry run") && !strings.Contains(stdout, "render") {
		t.Errorf("Dry run output should indicate what would happen: %s", stdout)
	}
}
