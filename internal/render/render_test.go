package render

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wernerstrydom/render/internal/config"
	"github.com/wernerstrydom/render/internal/engine"
	"github.com/wernerstrydom/render/internal/output"
)

func TestCollect_BasicTemplateDir(t *testing.T) {
	dir := t.TempDir()

	// Create template directory
	tmplDir := filepath.Join(dir, "templates")
	mkdir(t, tmplDir)
	writeFile(t, tmplDir, "config.yaml.tmpl", "name: {{ .name }}")
	writeFile(t, tmplDir, "static.txt", "static content")

	outDir := filepath.Join(dir, "output")

	plan, err := Collect(CollectConfig{
		TemplateDir: tmplDir,
		OutputDir:   outDir,
		Data:        map[string]any{"name": "TestApp"},
		Engine:      engine.New(),
	})
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if len(plan.Outputs) != 2 {
		t.Errorf("Expected 2 outputs, got %d", len(plan.Outputs))
	}

	// Check that template was rendered
	var foundTemplate, foundStatic bool
	for _, out := range plan.Outputs {
		if strings.HasSuffix(out.OutputPath, "config.yaml") {
			foundTemplate = true
			if out.CopyFrom != "" {
				t.Error("Template should have Content, not CopyFrom")
			}
			if string(out.Content) != "name: TestApp" {
				t.Errorf("Content = %q, want %q", string(out.Content), "name: TestApp")
			}
		}
		if strings.HasSuffix(out.OutputPath, "static.txt") {
			foundStatic = true
			if out.CopyFrom == "" {
				t.Error("Static file should have CopyFrom")
			}
		}
	}

	if !foundTemplate {
		t.Error("Template output not found")
	}
	if !foundStatic {
		t.Error("Static output not found")
	}
}

func TestCollect_WithConfig(t *testing.T) {
	dir := t.TempDir()

	// Create template directory with config
	tmplDir := filepath.Join(dir, "templates")
	mkdir(t, tmplDir)
	writeFile(t, tmplDir, "model.go.tmpl", "package {{ .package }}")
	writeFile(t, tmplDir, ".render.yaml", `paths:
  "model.go.tmpl": "{{ .name | snakeCase }}.go"
`)

	outDir := filepath.Join(dir, "output")

	cfg, err := config.Load(tmplDir)
	if err != nil {
		t.Fatalf("Config load failed: %v", err)
	}

	plan, err := Collect(CollectConfig{
		TemplateDir: tmplDir,
		OutputDir:   outDir,
		Data:        map[string]any{"name": "UserProfile", "package": "models"},
		Config:      cfg,
		Engine:      engine.New(),
	})
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	// Should have only 1 output (config file is skipped)
	if len(plan.Outputs) != 1 {
		t.Errorf("Expected 1 output, got %d", len(plan.Outputs))
	}

	// Check the output path was transformed
	if len(plan.Outputs) > 0 {
		out := plan.Outputs[0]
		if !strings.HasSuffix(out.OutputPath, "user_profile.go") {
			t.Errorf("OutputPath = %q, should end with user_profile.go", out.OutputPath)
		}
	}
}

func TestCollect_SkipsConfigFile(t *testing.T) {
	dir := t.TempDir()

	tmplDir := filepath.Join(dir, "templates")
	mkdir(t, tmplDir)
	writeFile(t, tmplDir, "file.txt", "content")
	writeFile(t, tmplDir, ".render.yaml", "paths: {}")

	outDir := filepath.Join(dir, "output")

	plan, err := Collect(CollectConfig{
		TemplateDir: tmplDir,
		OutputDir:   outDir,
		Data:        nil,
		Engine:      engine.New(),
	})
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	// Should only have 1 output (file.txt), not .render.yaml
	if len(plan.Outputs) != 1 {
		t.Errorf("Expected 1 output, got %d", len(plan.Outputs))
	}

	for _, out := range plan.Outputs {
		if strings.Contains(out.SourcePath, ".render.yaml") {
			t.Error("Config file should be skipped")
		}
	}
}

func TestPlan_Validate_NoCollisions(t *testing.T) {
	plan := &Plan{
		Outputs: []Output{
			{SourcePath: "a.go", OutputPath: "/out/a.go"},
			{SourcePath: "b.go", OutputPath: "/out/b.go"},
		},
	}

	errs := plan.Validate()
	if len(errs) != 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
}

func TestPlan_Validate_Collision(t *testing.T) {
	plan := &Plan{
		Outputs: []Output{
			{SourcePath: "a.go", OutputPath: "/out/same.go"},
			{SourcePath: "b.go", OutputPath: "/out/same.go"},
		},
	}

	errs := plan.Validate()
	if len(errs) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(errs))
	}

	if !strings.Contains(errs[0].Error(), "collision") {
		t.Errorf("Error should mention collision: %v", errs[0])
	}
}

func TestPlan_Preview(t *testing.T) {
	plan := &Plan{
		Outputs: []Output{
			{SourcePath: "template.go.tmpl", OutputPath: "/out/model.go", Content: []byte("content")},
			{SourcePath: "static.txt", OutputPath: "/out/static.txt", CopyFrom: "/src/static.txt"},
		},
	}

	preview := plan.Preview()

	if !strings.Contains(preview, "[render]") {
		t.Error("Preview should show render action for template")
	}
	if !strings.Contains(preview, "[copy]") {
		t.Error("Preview should show copy action for static file")
	}
	if !strings.Contains(preview, "template.go.tmpl") {
		t.Error("Preview should show source path")
	}
}

func TestPlan_Execute(t *testing.T) {
	dir := t.TempDir()

	// Create source file for copy
	srcFile := filepath.Join(dir, "source.txt")
	if err := os.WriteFile(srcFile, []byte("source content"), 0644); err != nil {
		t.Fatalf("Failed to write source: %v", err)
	}

	outDir := filepath.Join(dir, "output")

	plan := &Plan{
		Outputs: []Output{
			{
				SourcePath:  "template.go.tmpl",
				OutputPath:  filepath.Join(outDir, "rendered.go"),
				Content:     []byte("rendered content"),
				Permissions: 0644,
			},
			{
				SourcePath:  "static.txt",
				OutputPath:  filepath.Join(outDir, "copied.txt"),
				CopyFrom:    srcFile,
				Permissions: 0644,
			},
		},
	}

	writer := output.New(true)
	if _, err := plan.Execute(writer); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Verify rendered file
	rendered, err := os.ReadFile(filepath.Join(outDir, "rendered.go"))
	if err != nil {
		t.Fatalf("Failed to read rendered: %v", err)
	}
	if string(rendered) != "rendered content" {
		t.Errorf("Rendered content = %q, want %q", string(rendered), "rendered content")
	}

	// Verify copied file
	copied, err := os.ReadFile(filepath.Join(outDir, "copied.txt"))
	if err != nil {
		t.Fatalf("Failed to read copied: %v", err)
	}
	if string(copied) != "source content" {
		t.Errorf("Copied content = %q, want %q", string(copied), "source content")
	}
}

func TestCollect_PreservesDirectoryStructure(t *testing.T) {
	dir := t.TempDir()

	tmplDir := filepath.Join(dir, "templates")
	mkdir(t, filepath.Join(tmplDir, "a", "b", "c"))
	writeFile(t, tmplDir, "a/b/c/deep.txt", "deep")

	outDir := filepath.Join(dir, "output")

	plan, err := Collect(CollectConfig{
		TemplateDir: tmplDir,
		OutputDir:   outDir,
		Data:        nil,
		Engine:      engine.New(),
	})
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if len(plan.Outputs) != 1 {
		t.Fatalf("Expected 1 output, got %d", len(plan.Outputs))
	}

	// Check that nested structure is preserved
	out := plan.Outputs[0]
	expected := filepath.Join(outDir, "a", "b", "c", "deep.txt")
	if out.OutputPath != expected {
		t.Errorf("OutputPath = %q, want %q", out.OutputPath, expected)
	}
}

func TestCollect_OverwriteField(t *testing.T) {
	dir := t.TempDir()

	// Create template directory with config
	tmplDir := filepath.Join(dir, "templates")
	mkdir(t, tmplDir)
	writeFile(t, tmplDir, "regular.go.tmpl", "package main")
	writeFile(t, tmplDir, "protected.go.tmpl", "package protected")
	writeFile(t, tmplDir, ".render.yaml", `paths:
  "regular.go.tmpl": "regular.go"
  "protected.go.tmpl":
    path: "protected.go"
    overwrite: false
`)

	outDir := filepath.Join(dir, "output")

	cfg, err := config.Load(tmplDir)
	if err != nil {
		t.Fatalf("Config load failed: %v", err)
	}

	plan, err := Collect(CollectConfig{
		TemplateDir: tmplDir,
		OutputDir:   outDir,
		Data:        map[string]any{},
		Config:      cfg,
		Engine:      engine.New(),
	})
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if len(plan.Outputs) != 2 {
		t.Fatalf("Expected 2 outputs, got %d", len(plan.Outputs))
	}

	// Check Overwrite field is set correctly
	for _, out := range plan.Outputs {
		if strings.HasSuffix(out.OutputPath, "regular.go") {
			if !out.Overwrite {
				t.Error("regular.go should have Overwrite=true")
			}
		}
		if strings.HasSuffix(out.OutputPath, "protected.go") {
			if out.Overwrite {
				t.Error("protected.go should have Overwrite=false")
			}
		}
	}
}

func TestPlan_Execute_NoOverwrite(t *testing.T) {
	dir := t.TempDir()
	outDir := filepath.Join(dir, "output")
	mkdir(t, outDir)

	// Create existing file that should be protected
	protectedPath := filepath.Join(outDir, "protected.go")
	if err := os.WriteFile(protectedPath, []byte("original content"), 0644); err != nil {
		t.Fatalf("Failed to write protected file: %v", err)
	}

	plan := &Plan{
		Outputs: []Output{
			{
				SourcePath:  "regular.go.tmpl",
				OutputPath:  filepath.Join(outDir, "regular.go"),
				Content:     []byte("new content"),
				Permissions: 0644,
				Overwrite:   true,
			},
			{
				SourcePath:  "protected.go.tmpl",
				OutputPath:  protectedPath,
				Content:     []byte("new content"),
				Permissions: 0644,
				Overwrite:   false, // Should not overwrite existing file
			},
		},
	}

	writer := output.New(true) // force=true, but Overwrite=false should still skip
	result, err := plan.Execute(writer)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Check that protected file was skipped
	if !result.Skipped[protectedPath] {
		t.Error("Protected file should be in Skipped map")
	}

	// Verify protected file content unchanged
	content, err := os.ReadFile(protectedPath)
	if err != nil {
		t.Fatalf("Failed to read protected file: %v", err)
	}
	if string(content) != "original content" {
		t.Errorf("Protected file content = %q, want %q", string(content), "original content")
	}

	// Verify regular file was created
	regularContent, err := os.ReadFile(filepath.Join(outDir, "regular.go"))
	if err != nil {
		t.Fatalf("Failed to read regular file: %v", err)
	}
	if string(regularContent) != "new content" {
		t.Errorf("Regular file content = %q, want %q", string(regularContent), "new content")
	}
}

func TestPlan_Execute_NoOverwrite_NewFile(t *testing.T) {
	dir := t.TempDir()
	outDir := filepath.Join(dir, "output")
	mkdir(t, outDir)

	// Protected file doesn't exist yet
	protectedPath := filepath.Join(outDir, "protected.go")

	plan := &Plan{
		Outputs: []Output{
			{
				SourcePath:  "protected.go.tmpl",
				OutputPath:  protectedPath,
				Content:     []byte("new content"),
				Permissions: 0644,
				Overwrite:   false, // But file doesn't exist, so should create
			},
		},
	}

	writer := output.New(true)
	result, err := plan.Execute(writer)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Should NOT be in skipped (file was created since it didn't exist)
	if result.Skipped[protectedPath] {
		t.Error("Protected file should NOT be in Skipped map when it didn't exist")
	}

	// Verify file was created
	content, err := os.ReadFile(protectedPath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(content) != "new content" {
		t.Errorf("Content = %q, want %q", string(content), "new content")
	}
}

// Helper functions

func mkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("Failed to create parent: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write: %v", err)
	}
}
