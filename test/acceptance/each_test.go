package acceptance

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestEachBasicFile tests each mode with a file template and array data.
func TestEachBasicFile(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "user.txt.tmpl", `Name: {{ .name }}
ID: {{ .id }}`)

	// Create data with array (directly, no jq query needed)
	data := writeFile(t, dir, "users.json", `[
		{"id": 1, "name": "Alice"},
		{"id": 2, "name": "Bob"},
		{"id": 3, "name": "Charlie"}
	]`)

	// Output pattern with template expression (triggers each mode)
	outputPattern := filepath.Join(dir, "output", "user-{{.id}}.txt")

	stdout, stderr, err := runRender(t, tmpl, data, "-o", outputPattern)
	if err != nil {
		t.Fatalf("render failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	// Verify files were created
	for i := 1; i <= 3; i++ {
		path := filepath.Join(dir, "output", fmt.Sprintf("user-%d.txt", i))
		if !fileExists(path) {
			t.Errorf("File not created: %s", path)
		}
	}

	// Verify content
	content := readFile(t, filepath.Join(dir, "output", "user-1.txt"))
	if !strings.Contains(content, "Name: Alice") || !strings.Contains(content, "ID: 1") {
		t.Errorf("Unexpected content: %s", content)
	}
}

// TestEachWithCasingFunctions tests each mode with casing functions.
func TestEachWithCasingFunctions(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "component.txt.tmpl", `{{ pascalCase .name }}Component`)

	data := writeFile(t, dir, "components.json", `[
		{"name": "user_profile"},
		{"name": "shopping-cart"},
		{"name": "navigation menu"}
	]`)

	outputPattern := filepath.Join(dir, "output", "{{kebabCase .name}}.txt")

	_, _, err := runRender(t, tmpl, data, "-o", outputPattern)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Verify files with kebab-case names
	expectedFiles := []string{
		"user-profile.txt",
		"shopping-cart.txt",
		"navigation-menu.txt",
	}

	for _, f := range expectedFiles {
		path := filepath.Join(dir, "output", f)
		if !fileExists(path) {
			t.Errorf("File not created: %s", f)
		}
	}
}

// TestEachWithDirectory tests each mode with a directory template.
func TestEachWithDirectory(t *testing.T) {
	dir := createTempDir(t)

	// Create template directory
	tmplDir := filepath.Join(dir, "service-template")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	writeFile(t, tmplDir, "main.go.tmpl", `package {{ .name }}`)
	writeFile(t, tmplDir, "README.md.tmpl", `# {{ .name }} Service`)
	writeFile(t, tmplDir, "config.json", `{"service": "{{ .name }}"}`)

	data := writeFile(t, dir, "services.json", `[
		{"name": "auth"},
		{"name": "users"}
	]`)

	// Dynamic output path (triggers each-directory mode)
	outputPattern := filepath.Join(dir, "output", "{{.name}}")

	stdout, stderr, err := runRender(t, tmplDir, data, "-o", outputPattern)
	if err != nil {
		t.Fatalf("render failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	// Verify directory structure for each service
	for _, name := range []string{"auth", "users"} {
		serviceDir := filepath.Join(dir, "output", name)

		mainPath := filepath.Join(serviceDir, "main.go")
		if !fileExists(mainPath) {
			t.Errorf("main.go not created in %s", name)
		} else {
			content := readFile(t, mainPath)
			if !strings.Contains(content, "package "+name) {
				t.Errorf("main.go content wrong: %s", content)
			}
		}

		readmePath := filepath.Join(serviceDir, "README.md")
		if !fileExists(readmePath) {
			t.Errorf("README.md not created in %s", name)
		}
	}
}

// TestEachWithYAML tests each mode with YAML data.
func TestEachWithYAML(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "config.txt.tmpl", `name: {{ .name }}
port: {{ .port }}`)

	data := writeFile(t, dir, "services.yaml", `- name: web
  port: 8080
- name: api
  port: 3000`)

	outputPattern := filepath.Join(dir, "output", "{{.name}}.yaml")

	_, _, err := runRender(t, tmpl, data, "-o", outputPattern)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Verify files
	webContent := readFile(t, filepath.Join(dir, "output", "web.yaml"))
	if !strings.Contains(webContent, "port: 8080") {
		t.Errorf("Wrong content: %s", webContent)
	}

	apiContent := readFile(t, filepath.Join(dir, "output", "api.yaml"))
	if !strings.Contains(apiContent, "port: 3000") {
		t.Errorf("Wrong content: %s", apiContent)
	}
}

// TestEachEmptyArray tests each mode with empty array.
func TestEachEmptyArray(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "item.txt.tmpl", "{{ .name }}")
	data := writeFile(t, dir, "data.json", `[]`)

	outputPattern := filepath.Join(dir, "output", "{{.name}}.txt")

	// Should succeed but create no files
	_, _, err := runRender(t, tmpl, data, "-o", outputPattern)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Output directory might not exist or be empty
	if fileExists(filepath.Join(dir, "output")) {
		entries, err := os.ReadDir(filepath.Join(dir, "output"))
		if err == nil && len(entries) > 0 {
			t.Errorf("Expected no output files, found %d", len(entries))
		}
	}
}

// TestEachForceOverwrite tests --force with each mode.
func TestEachForceOverwrite(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "item.txt.tmpl", "new: {{ .name }}")
	data := writeFile(t, dir, "data.json", `[{"id": "1", "name": "Alice"}]`)

	// Pre-create the output file with different content
	outputPattern := filepath.Join(dir, "{{.id}}.txt")
	writeFile(t, dir, "1.txt", "original content")

	// Without force should fail
	_, _, err := runRender(t, tmpl, data, "-o", outputPattern)
	if err == nil {
		t.Fatal("render should fail without --force")
	}

	// With force should succeed
	_, _, err = runRender(t, tmpl, data, "-o", outputPattern, "--force")
	if err != nil {
		t.Fatalf("render --force failed: %v", err)
	}

	content := readFile(t, filepath.Join(dir, "1.txt"))
	if !strings.Contains(content, "new: Alice") {
		t.Errorf("Content not updated: %s", content)
	}
}

// TestEachWithObjectData tests each mode with object data (treated as single item).
func TestEachWithObjectData(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "item.txt.tmpl", "{{ .name }}")
	// Object data - will be treated as single item array
	data := writeFile(t, dir, "data.json", `{"id": "only", "name": "Single"}`)

	outputPattern := filepath.Join(dir, "{{.id}}.txt")

	_, _, err := runRender(t, tmpl, data, "-o", outputPattern)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Should create exactly one file
	if !fileExists(filepath.Join(dir, "only.txt")) {
		t.Error("File not created for single object")
	}

	content := readFile(t, filepath.Join(dir, "only.txt"))
	if content != "Single" {
		t.Errorf("Content = %q, want %q", content, "Single")
	}
}

// TestEachDryRun tests dry-run with each mode.
func TestEachDryRun(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "item.txt.tmpl", "{{ .name }}")
	data := writeFile(t, dir, "data.json", `[{"id": "1", "name": "Alice"}, {"id": "2", "name": "Bob"}]`)

	outputPattern := filepath.Join(dir, "output", "{{.id}}.txt")

	stdout, _, err := runRender(t, tmpl, data, "-o", outputPattern, "--dry-run")
	if err != nil {
		t.Fatalf("render --dry-run failed: %v", err)
	}

	// No files should be created
	if fileExists(filepath.Join(dir, "output")) {
		t.Error("Output directory should not exist in dry-run mode")
	}

	if !strings.Contains(stdout, "Dry run") {
		t.Errorf("Output should indicate dry run: %s", stdout)
	}
}

// TestEachInternalCollision tests that internal collision is detected.
func TestEachInternalCollision(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "item.txt.tmpl", "{{ .name }}")
	// Both items produce the same output path
	data := writeFile(t, dir, "data.json", `[{"id": "same", "name": "Alice"}, {"id": "same", "name": "Bob"}]`)

	outputPattern := filepath.Join(dir, "{{.id}}.txt")

	_, stderr, err := runRender(t, tmpl, data, "-o", outputPattern)
	if err == nil {
		t.Fatal("render should fail with internal path collision")
	}

	if !strings.Contains(stderr, "collision") {
		t.Errorf("Error should mention collision: %s", stderr)
	}
}

// TestEachWithJSON tests JSON output mode.
func TestEachWithJSON(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "item.txt.tmpl", "{{ .name }}")
	data := writeFile(t, dir, "data.json", `[{"id": "1", "name": "Alice"}]`)

	outputPattern := filepath.Join(dir, "{{.id}}.txt")

	stdout, _, err := runRender(t, tmpl, data, "-o", outputPattern, "--json")
	if err != nil {
		t.Fatalf("render --json failed: %v", err)
	}

	if !strings.Contains(stdout, `"status"`) || !strings.Contains(stdout, `"success"`) {
		t.Errorf("JSON output missing expected fields: %s", stdout)
	}
}

// TestItemQueryBasic tests --item-query flag to extract items from nested data.
func TestItemQueryBasic(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "user.txt.tmpl", `Name: {{ .name }}
ID: {{ .id }}`)

	// Data with nested array (not top-level)
	data := writeFile(t, dir, "data.json", `{
		"metadata": {"version": "1.0"},
		"users": [
			{"id": 1, "name": "Alice"},
			{"id": 2, "name": "Bob"},
			{"id": 3, "name": "Charlie"}
		]
	}`)

	outputPattern := filepath.Join(dir, "output", "user-{{.id}}.txt")

	stdout, stderr, err := runRender(t, tmpl, data, "-o", outputPattern, "--item-query", ".users[]")
	if err != nil {
		t.Fatalf("render failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	// Verify files were created
	for i := 1; i <= 3; i++ {
		path := filepath.Join(dir, "output", fmt.Sprintf("user-%d.txt", i))
		if !fileExists(path) {
			t.Errorf("File not created: %s", path)
		}
	}

	// Verify content
	content := readFile(t, filepath.Join(dir, "output", "user-1.txt"))
	if !strings.Contains(content, "Name: Alice") || !strings.Contains(content, "ID: 1") {
		t.Errorf("Unexpected content: %s", content)
	}
}

// TestItemQueryWithFilter tests --item-query with jq select filter.
func TestItemQueryWithFilter(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "item.txt.tmpl", "{{ .name }}: active={{ .active }}")

	data := writeFile(t, dir, "data.json", `{
		"items": [
			{"id": "a", "name": "Alpha", "active": true},
			{"id": "b", "name": "Beta", "active": false},
			{"id": "c", "name": "Charlie", "active": true}
		]
	}`)

	outputPattern := filepath.Join(dir, "output", "{{.id}}.txt")

	// Only active items
	_, _, err := runRender(t, tmpl, data, "-o", outputPattern, "--item-query", ".items[] | select(.active)")
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Should create files for Alpha and Charlie, but not Beta
	if !fileExists(filepath.Join(dir, "output", "a.txt")) {
		t.Error("File for Alpha not created")
	}
	if fileExists(filepath.Join(dir, "output", "b.txt")) {
		t.Error("File for Beta should not be created (not active)")
	}
	if !fileExists(filepath.Join(dir, "output", "c.txt")) {
		t.Error("File for Charlie not created")
	}
}

// TestQueryBasic tests --query flag to transform data before rendering.
func TestQueryBasic(t *testing.T) {
	dir := createTempDir(t)

	// Template expects transformed data structure
	tmpl := writeFile(t, dir, "report.txt.tmpl", `Total items: {{ .total }}
First item: {{ .first }}`)

	data := writeFile(t, dir, "data.json", `{
		"items": ["apple", "banana", "cherry"]
	}`)

	output := filepath.Join(dir, "report.txt")

	// Transform data: extract total count and first item
	_, _, err := runRender(t, tmpl, data, "-o", output, "--query", "{total: (.items | length), first: .items[0]}")
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	content := readFile(t, output)
	if !strings.Contains(content, "Total items: 3") {
		t.Errorf("Missing total count in output: %s", content)
	}
	if !strings.Contains(content, "First item: apple") {
		t.Errorf("Missing first item in output: %s", content)
	}
}

// TestQueryAndItemQueryCombined tests using both --query and --item-query together.
func TestQueryAndItemQueryCombined(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "member.txt.tmpl", "{{ .name }} ({{ .role }})")

	// Deeply nested data
	data := writeFile(t, dir, "data.json", `{
		"organization": {
			"departments": {
				"engineering": {
					"team": {
						"members": [
							{"id": "e1", "name": "Alice", "role": "lead"},
							{"id": "e2", "name": "Bob", "role": "dev"}
						]
					}
				}
			}
		}
	}`)

	outputPattern := filepath.Join(dir, "output", "{{.id}}.txt")

	// --query narrows to the team object, --item-query extracts members
	_, _, err := runRender(t, tmpl, data, "-o", outputPattern,
		"--query", ".organization.departments.engineering.team",
		"--item-query", ".members[]")
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Verify files
	if !fileExists(filepath.Join(dir, "output", "e1.txt")) {
		t.Error("File for e1 not created")
	}
	if !fileExists(filepath.Join(dir, "output", "e2.txt")) {
		t.Error("File for e2 not created")
	}

	content := readFile(t, filepath.Join(dir, "output", "e1.txt"))
	if content != "Alice (lead)" {
		t.Errorf("Content = %q, want %q", content, "Alice (lead)")
	}
}

// TestItemQueryInvalidExpression tests error handling for invalid jq expression.
func TestItemQueryInvalidExpression(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "item.txt.tmpl", "{{ .name }}")
	data := writeFile(t, dir, "data.json", `{"items": []}`)

	outputPattern := filepath.Join(dir, "{{.id}}.txt")

	_, stderr, err := runRender(t, tmpl, data, "-o", outputPattern, "--item-query", ".items[")
	if err == nil {
		t.Fatal("render should fail with invalid jq expression")
	}

	if !strings.Contains(stderr, "query") || !strings.Contains(stderr, "parse") {
		t.Errorf("Error should mention query parsing: %s", stderr)
	}
}

// TestQueryInvalidExpression tests error handling for invalid --query expression.
func TestQueryInvalidExpression(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "item.txt.tmpl", "{{ .name }}")
	data := writeFile(t, dir, "data.json", `{"name": "test"}`)

	output := filepath.Join(dir, "output.txt")

	_, stderr, err := runRender(t, tmpl, data, "-o", output, "--query", ".name[")
	if err == nil {
		t.Fatal("render should fail with invalid jq expression")
	}

	if !strings.Contains(stderr, "query") || !strings.Contains(stderr, "parse") {
		t.Errorf("Error should mention query parsing: %s", stderr)
	}
}

// TestItemQueryWithDirectoryTemplate tests --item-query with directory template.
func TestItemQueryWithDirectoryTemplate(t *testing.T) {
	dir := createTempDir(t)

	// Create template directory
	tmplDir := filepath.Join(dir, "service-template")
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	writeFile(t, tmplDir, "main.go.tmpl", `package {{ .name }}`)
	writeFile(t, tmplDir, "README.md.tmpl", `# {{ .name }} Service`)

	// Data with nested services array
	data := writeFile(t, dir, "data.json", `{
		"config": {"version": "2.0"},
		"services": [
			{"name": "auth"},
			{"name": "users"}
		]
	}`)

	outputPattern := filepath.Join(dir, "output", "{{.name}}")

	_, _, err := runRender(t, tmplDir, data, "-o", outputPattern, "--item-query", ".services[]")
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Verify directory structure for each service
	for _, name := range []string{"auth", "users"} {
		serviceDir := filepath.Join(dir, "output", name)

		mainPath := filepath.Join(serviceDir, "main.go")
		if !fileExists(mainPath) {
			t.Errorf("main.go not created in %s", name)
		} else {
			content := readFile(t, mainPath)
			if !strings.Contains(content, "package "+name) {
				t.Errorf("main.go content wrong: %s", content)
			}
		}
	}
}

// TestItemQueryEmptyResult tests --item-query that returns no items.
func TestItemQueryEmptyResult(t *testing.T) {
	dir := createTempDir(t)

	tmpl := writeFile(t, dir, "item.txt.tmpl", "{{ .name }}")
	data := writeFile(t, dir, "data.json", `{"items": []}`)

	outputPattern := filepath.Join(dir, "output", "{{.id}}.txt")

	// Should succeed but create no files
	_, _, err := runRender(t, tmpl, data, "-o", outputPattern, "--item-query", ".items[]")
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// Output directory might not exist or be empty
	if fileExists(filepath.Join(dir, "output")) {
		entries, err := os.ReadDir(filepath.Join(dir, "output"))
		if err == nil && len(entries) > 0 {
			t.Errorf("Expected no output files, found %d", len(entries))
		}
	}
}
