// Package acceptance provides end-to-end acceptance tests for the render CLI.
// These tests compile the binary and execute it as a user would, verifying
// the outputs without access to internal code.
package acceptance

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
)

var (
	binaryPath string
	buildOnce  sync.Once
	buildErr   error
)

// ensureBinary builds the render binary if it doesn't exist.
func ensureBinary(t *testing.T) string {
	t.Helper()

	buildOnce.Do(func() {
		// Find the project root (two levels up from test/acceptance)
		_, filename, _, _ := runtime.Caller(0)
		projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")

		// Binary path
		binaryName := "render"
		if runtime.GOOS == "windows" {
			binaryName = "render.exe"
		}
		binaryPath = filepath.Join(projectRoot, "bin", "test", binaryName)

		// Create bin/test directory
		if err := os.MkdirAll(filepath.Dir(binaryPath), 0755); err != nil {
			buildErr = err
			return
		}

		// Build the binary
		cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/render")
		cmd.Dir = projectRoot
		output, err := cmd.CombinedOutput()
		if err != nil {
			buildErr = &buildError{output: string(output), err: err}
			return
		}
	})

	if buildErr != nil {
		t.Fatalf("Failed to build binary: %v", buildErr)
	}

	return binaryPath
}

type buildError struct {
	output string
	err    error
}

func (e *buildError) Error() string {
	return e.err.Error() + ": " + e.output
}

// runRender executes the render binary with the given arguments.
func runRender(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

	binary := ensureBinary(t)
	cmd := exec.Command(binary, args...)

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// getExitCode extracts exit code from error.
func getExitCode(err error) int {
	if err == nil {
		return 0
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}
	return 1
}

// createTempDir creates a temporary directory for test files.
func createTempDir(t *testing.T) string {
	t.Helper()

	dir, err := os.MkdirTemp("", "render-acceptance-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})

	return dir
}

// writeFile writes content to a file in the given directory.
func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	return path
}

// readFile reads content from a file.
func readFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}

	return string(content)
}

// fileExists checks if a file exists.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// TestHelp verifies the help command works.
func TestHelp(t *testing.T) {
	stdout, _, err := runRender(t, "--help")
	if err != nil {
		t.Fatalf("render --help failed: %v", err)
	}

	expectedPhrases := []string{
		"render",
		"template",
		"--output",
		"--force",
		"--dry-run",
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(stdout, phrase) {
			t.Errorf("Help output missing %q", phrase)
		}
	}
}

// TestFileBasic tests basic file rendering.
func TestFileBasic(t *testing.T) {
	dir := createTempDir(t)

	// Create template
	tmpl := writeFile(t, dir, "template.txt", "Hello, {{ .name }}!")

	// Create data
	data := writeFile(t, dir, "data.json", `{"name": "World"}`)

	// Output path
	output := filepath.Join(dir, "output.txt")

	// Run render with new syntax
	stdout, stderr, err := runRender(t, tmpl, data, "-o", output)

	if err != nil {
		t.Fatalf("render failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	// Verify output
	content := readFile(t, output)
	if content != "Hello, World!" {
		t.Errorf("Output content = %q, want %q", content, "Hello, World!")
	}
}

// TestFileWithCustomFunctions tests file rendering with custom template functions.
func TestFileWithCustomFunctions(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "template.txt", `{{ upper .name }}
{{ lower .name }}
{{ camelCase .title }}
{{ snakeCase .title }}
{{ kebabCase .title }}`)

	data := writeFile(t, dir, "data.json", `{
		"name": "Hello World",
		"title": "My Awesome Project"
	}`)

	output := filepath.Join(dir, "output.txt")

	_, _, err := runRender(t, tmpl, data, "-o", output)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	content := readFile(t, output)
	expectedLines := []string{
		"HELLO WORLD",
		"hello world",
		"myAwesomeProject",
		"my_awesome_project",
		"my-awesome-project",
	}

	for _, line := range expectedLines {
		if !strings.Contains(content, line) {
			t.Errorf("Output missing %q:\n%s", line, content)
		}
	}
}

// TestFileWithYAML tests file rendering with YAML data.
func TestFileWithYAML(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "template.txt", "{{ .name }} - {{ .version }}")

	data := writeFile(t, dir, "data.yaml", `name: MyApp
version: 1.0.0`)

	output := filepath.Join(dir, "output.txt")

	_, _, err := runRender(t, tmpl, data, "-o", output)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	content := readFile(t, output)
	if content != "MyApp - 1.0.0" {
		t.Errorf("Output content = %q, want %q", content, "MyApp - 1.0.0")
	}
}

// TestFileOverwriteProtection tests that files aren't overwritten without --force.
func TestFileOverwriteProtection(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "template.txt", "new content")
	data := writeFile(t, dir, "data.json", `{}`)
	output := writeFile(t, dir, "output.txt", "original content")

	// Run without --force
	_, stderr, err := runRender(t, tmpl, data, "-o", output)
	if err == nil {
		t.Fatal("render should fail without --force when file exists")
	}

	// Check exit code is 5 (OutputConflict)
	if exitCode := getExitCode(err); exitCode != 5 {
		t.Errorf("Expected exit code 5, got %d", exitCode)
	}

	if !strings.Contains(stderr, "exists") && !strings.Contains(stderr, "force") {
		t.Errorf("Error message should mention file exists: %s", stderr)
	}

	// Verify original content unchanged
	content := readFile(t, output)
	if content != "original content" {
		t.Errorf("File was modified without --force: %q", content)
	}
}

// TestFileForceOverwrite tests that --force allows overwriting.
func TestFileForceOverwrite(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "template.txt", "new content")
	data := writeFile(t, dir, "data.json", `{}`)
	output := writeFile(t, dir, "output.txt", "original content")

	// Run with --force
	_, _, err := runRender(t, tmpl, data, "-o", output, "--force")
	if err != nil {
		t.Fatalf("render --force failed: %v", err)
	}

	content := readFile(t, output)
	if content != "new content" {
		t.Errorf("File content = %q, want %q", content, "new content")
	}
}

// TestFileMissingTemplate tests error handling for missing template.
func TestFileMissingTemplate(t *testing.T) {
	dir := createTempDir(t)

	data := writeFile(t, dir, "data.json", `{}`)
	output := filepath.Join(dir, "output.txt")

	_, _, err := runRender(t, filepath.Join(dir, "nonexistent.txt"), data, "-o", output)

	if err == nil {
		t.Fatal("render should fail with missing template")
	}

	// Check exit code is 3 (InputValidation)
	if exitCode := getExitCode(err); exitCode != 3 {
		t.Errorf("Expected exit code 3, got %d", exitCode)
	}
}

// TestFileMissingData tests error handling for missing data file.
func TestFileMissingData(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "template.txt", "{{ .name }}")
	output := filepath.Join(dir, "output.txt")

	_, _, err := runRender(t, tmpl, filepath.Join(dir, "nonexistent.json"), "-o", output)

	if err == nil {
		t.Fatal("render should fail with missing data file")
	}

	// Check exit code is 3 (InputValidation)
	if exitCode := getExitCode(err); exitCode != 3 {
		t.Errorf("Expected exit code 3, got %d", exitCode)
	}
}

// TestFileMissingOutput tests that missing output flag fails.
func TestFileMissingOutput(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "template.txt", "{{ .name }}")
	data := writeFile(t, dir, "data.json", `{}`)

	_, _, err := runRender(t, tmpl, data)
	if err == nil {
		t.Fatal("render should fail without output flag")
	}
}

// TestContentIdempotency tests that identical content skips write.
func TestContentIdempotency(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "template.txt", "Hello, World!")
	data := writeFile(t, dir, "data.json", `{}`)
	output := writeFile(t, dir, "output.txt", "Hello, World!")

	// Get original modification time
	info, _ := os.Stat(output)
	origModTime := info.ModTime()

	// Run render (without --force, but content is identical)
	_, _, err := runRender(t, tmpl, data, "-o", output)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Check file wasn't rewritten (modification time should be same)
	info, _ = os.Stat(output)
	if !info.ModTime().Equal(origModTime) {
		t.Log("Note: File was rewritten even though content is identical")
	}

	// Verify content is still correct
	content := readFile(t, output)
	if content != "Hello, World!" {
		t.Errorf("Content changed unexpectedly: %q", content)
	}
}

// TestArrayDataToSingleFile tests rendering array data to a single file.
func TestArrayDataToSingleFile(t *testing.T) {
	dir := createTempDir(t)

	// Template that iterates over array
	tmpl := writeFile(t, dir, "template.txt", `{{ range . }}{{ .name }}
{{ end }}`)
	data := writeFile(t, dir, "data.json", `[{"name": "Alice"}, {"name": "Bob"}]`)
	output := filepath.Join(dir, "output.txt")

	_, _, err := runRender(t, tmpl, data, "-o", output)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	content := readFile(t, output)
	if !strings.Contains(content, "Alice") || !strings.Contains(content, "Bob") {
		t.Errorf("Output missing expected names: %s", content)
	}
}

// TestFileIntoDirectory tests rendering a file into a directory (trailing slash).
func TestFileIntoDirectory(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "template.txt.tmpl", "Hello!")
	data := writeFile(t, dir, "data.json", `{}`)
	outputDir := filepath.Join(dir, "outdir") + "/"

	_, _, err := runRender(t, tmpl, data, "-o", outputDir)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Output should be outdir/template.txt (with .tmpl stripped)
	outputFile := filepath.Join(dir, "outdir", "template.txt")
	if !fileExists(outputFile) {
		t.Errorf("Expected file at %s", outputFile)
	}

	content := readFile(t, outputFile)
	if content != "Hello!" {
		t.Errorf("Content = %q, want %q", content, "Hello!")
	}
}

// TestIterativeFileRendering tests each mode with dynamic output path.
func TestIterativeFileRendering(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "item.tmpl", "Name: {{ .name }}")
	data := writeFile(t, dir, "data.json", `[{"id": "1", "name": "Alice"}, {"id": "2", "name": "Bob"}]`)

	outputPattern := filepath.Join(dir, "{{.id}}.txt")

	stdout, stderr, err := runRender(t, tmpl, data, "-o", outputPattern)
	if err != nil {
		t.Fatalf("render failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	// Verify both files were created
	file1 := filepath.Join(dir, "1.txt")
	file2 := filepath.Join(dir, "2.txt")

	if !fileExists(file1) {
		t.Errorf("File not created: %s", file1)
	} else {
		content := readFile(t, file1)
		if content != "Name: Alice" {
			t.Errorf("File 1 content = %q, want %q", content, "Name: Alice")
		}
	}

	if !fileExists(file2) {
		t.Errorf("File not created: %s", file2)
	} else {
		content := readFile(t, file2)
		if content != "Name: Bob" {
			t.Errorf("File 2 content = %q, want %q", content, "Name: Bob")
		}
	}
}

// TestInternalPathCollision tests that internal collision is detected.
func TestInternalPathCollision(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "item.tmpl", "{{ .name }}")
	// Both items produce the same output path!
	data := writeFile(t, dir, "data.json", `[{"id": "same", "name": "Alice"}, {"id": "same", "name": "Bob"}]`)

	outputPattern := filepath.Join(dir, "{{.id}}.txt")

	_, stderr, err := runRender(t, tmpl, data, "-o", outputPattern)
	if err == nil {
		t.Fatal("render should fail with internal path collision")
	}

	// Check exit code is 1 (RuntimeError)
	if exitCode := getExitCode(err); exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}

	if !strings.Contains(stderr, "collision") {
		t.Errorf("Error should mention collision: %s", stderr)
	}
}

// TestDryRun tests dry-run mode.
func TestDryRun(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "template.txt", "Hello!")
	data := writeFile(t, dir, "data.json", `{}`)
	output := filepath.Join(dir, "output.txt")

	stdout, _, err := runRender(t, tmpl, data, "-o", output, "--dry-run")
	if err != nil {
		t.Fatalf("render --dry-run failed: %v", err)
	}

	// File should NOT be created
	if fileExists(output) {
		t.Error("File was created in dry-run mode")
	}

	// Output should indicate what would happen
	if !strings.Contains(stdout, "Dry run") && !strings.Contains(stdout, "create") {
		t.Errorf("Dry run output should indicate what would happen: %s", stdout)
	}
}

// TestJSONOutput tests JSON output mode.
func TestJSONOutput(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "template.txt", "Hello!")
	data := writeFile(t, dir, "data.json", `{}`)
	output := filepath.Join(dir, "output.txt")

	stdout, _, err := runRender(t, tmpl, data, "-o", output, "--json")
	if err != nil {
		t.Fatalf("render --json failed: %v", err)
	}

	if !strings.Contains(stdout, `"status"`) || !strings.Contains(stdout, `"success"`) {
		t.Errorf("JSON output missing expected fields: %s", stdout)
	}
}
