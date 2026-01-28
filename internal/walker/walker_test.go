package walker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wernerstrydom/render/internal/engine"
	"github.com/wernerstrydom/render/internal/output"
)

func TestWalk(t *testing.T) {
	// Create temporary directories
	tmpDir := t.TempDir()
	tmplDir := filepath.Join(tmpDir, "templates")
	outDir := filepath.Join(tmpDir, "output")

	// Create template directory structure
	if err := os.MkdirAll(filepath.Join(tmplDir, "subdir"), 0755); err != nil {
		t.Fatalf("failed to create template directories: %v", err)
	}

	// Create a template file
	tmplContent := "Hello, {{.name}}!"
	if err := os.WriteFile(filepath.Join(tmplDir, "greeting.txt.tmpl"), []byte(tmplContent), 0644); err != nil {
		t.Fatalf("failed to write template file: %v", err)
	}

	// Create a non-template file
	staticContent := "This is static content"
	if err := os.WriteFile(filepath.Join(tmplDir, "static.txt"), []byte(staticContent), 0644); err != nil {
		t.Fatalf("failed to write static file: %v", err)
	}

	// Create a template in subdirectory
	subTmplContent := "Sub: {{.value}}"
	if err := os.WriteFile(filepath.Join(tmplDir, "subdir", "sub.txt.tmpl"), []byte(subTmplContent), 0644); err != nil {
		t.Fatalf("failed to write sub template file: %v", err)
	}

	// Track callbacks
	var rendered, copied []string

	// Walk
	cfg := Config{
		TemplateDir:      tmplDir,
		OutputDir:        outDir,
		Data:             map[string]any{"name": "World", "value": "test"},
		Engine:           engine.New(),
		Writer:           output.New(false),
		ValidateSymlinks: true,
		OnRendered: func(srcRel, dstAbs string) {
			rendered = append(rendered, srcRel)
		},
		OnCopied: func(srcRel, dstAbs string) {
			copied = append(copied, srcRel)
		},
	}

	if err := Walk(cfg); err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	// Verify rendered files
	greetingPath := filepath.Join(outDir, "greeting.txt")
	content, err := os.ReadFile(greetingPath)
	if err != nil {
		t.Errorf("failed to read greeting.txt: %v", err)
	} else if string(content) != "Hello, World!" {
		t.Errorf("greeting.txt content = %q, want %q", string(content), "Hello, World!")
	}

	subPath := filepath.Join(outDir, "subdir", "sub.txt")
	content, err = os.ReadFile(subPath)
	if err != nil {
		t.Errorf("failed to read subdir/sub.txt: %v", err)
	} else if string(content) != "Sub: test" {
		t.Errorf("subdir/sub.txt content = %q, want %q", string(content), "Sub: test")
	}

	// Verify copied files
	staticPath := filepath.Join(outDir, "static.txt")
	content, err = os.ReadFile(staticPath)
	if err != nil {
		t.Errorf("failed to read static.txt: %v", err)
	} else if string(content) != staticContent {
		t.Errorf("static.txt content = %q, want %q", string(content), staticContent)
	}

	// Verify callbacks
	if len(rendered) != 2 {
		t.Errorf("rendered count = %d, want 2", len(rendered))
	}
	if len(copied) != 1 {
		t.Errorf("copied count = %d, want 1", len(copied))
	}
}

func TestWalkWithPathTransform(t *testing.T) {
	// Create temporary directories
	tmpDir := t.TempDir()
	tmplDir := filepath.Join(tmpDir, "templates")
	outDir := filepath.Join(tmpDir, "output")

	// Create template directory
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("failed to create template directory: %v", err)
	}

	// Create a template file with templated filename
	tmplContent := "Content for {{.id}}"
	if err := os.WriteFile(filepath.Join(tmplDir, "{{.id}}.txt.tmpl"), []byte(tmplContent), 0644); err != nil {
		t.Fatalf("failed to write template file: %v", err)
	}

	eng := engine.New()

	cfg := Config{
		TemplateDir: tmplDir,
		OutputDir:   outDir,
		Data:        map[string]any{"id": "user-123"},
		Engine:      eng,
		Writer:      output.New(false),
		TransformPath: func(relPath string, data any) (string, error) {
			return eng.RenderString(relPath, data)
		},
	}

	if err := Walk(cfg); err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	// Verify the file was created with transformed name
	expectedPath := filepath.Join(outDir, "user-123.txt")
	content, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Errorf("failed to read user-123.txt: %v", err)
	} else if string(content) != "Content for user-123" {
		t.Errorf("user-123.txt content = %q, want %q", string(content), "Content for user-123")
	}
}

func TestWalkDirectoryTraversal(t *testing.T) {
	// Create temporary directories
	tmpDir := t.TempDir()
	tmplDir := filepath.Join(tmpDir, "templates")
	outDir := filepath.Join(tmpDir, "output")

	// Create template directory
	if err := os.MkdirAll(tmplDir, 0755); err != nil {
		t.Fatalf("failed to create template directory: %v", err)
	}

	// Create a simple template file
	if err := os.WriteFile(filepath.Join(tmplDir, "test.txt"), []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write template file: %v", err)
	}

	cfg := Config{
		TemplateDir:      tmplDir,
		OutputDir:        outDir,
		Data:             map[string]any{},
		Engine:           engine.New(),
		Writer:           output.New(false),
		ValidateSymlinks: true,
	}

	// This should work fine
	if err := Walk(cfg); err != nil {
		t.Fatalf("Walk failed: %v", err)
	}
}
