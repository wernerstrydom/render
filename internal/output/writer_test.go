package output

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	w := New(false)
	if w == nil {
		t.Fatal("New() returned nil")
	}
	if w.force {
		t.Error("New(false) should have force=false")
	}

	w = New(true)
	if !w.force {
		t.Error("New(true) should have force=true")
	}
}

func TestWrite(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "render-output-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	t.Run("write new file", func(t *testing.T) {
		w := New(false)
		path := filepath.Join(tmpDir, "new.txt")
		content := []byte("hello world")

		err := w.Write(path, content)
		if err != nil {
			t.Fatalf("Write() error = %v", err)
		}

		written, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("Failed to read written file: %v", err)
		}
		if string(written) != "hello world" {
			t.Errorf("Write() content = %q, want %q", string(written), "hello world")
		}
	})

	t.Run("write creates parent directories", func(t *testing.T) {
		w := New(false)
		path := filepath.Join(tmpDir, "nested", "dir", "file.txt")
		content := []byte("nested content")

		err := w.Write(path, content)
		if err != nil {
			t.Fatalf("Write() error = %v", err)
		}

		written, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("Failed to read written file: %v", err)
		}
		if string(written) != "nested content" {
			t.Errorf("Write() content = %q", string(written))
		}
	})

	t.Run("write fails without force on existing file", func(t *testing.T) {
		w := New(false)
		path := filepath.Join(tmpDir, "existing.txt")

		// Create existing file
		if err := os.WriteFile(path, []byte("original"), 0644); err != nil {
			t.Fatalf("Failed to create existing file: %v", err)
		}

		err := w.Write(path, []byte("new content"))
		if err == nil {
			t.Error("Write() should return error when file exists and force=false")
		}

		// Verify original content unchanged
		content, _ := os.ReadFile(path)
		if string(content) != "original" {
			t.Errorf("Write() modified file content: %q", string(content))
		}
	})

	t.Run("write succeeds with force on existing file", func(t *testing.T) {
		w := New(true)
		path := filepath.Join(tmpDir, "force.txt")

		// Create existing file
		if err := os.WriteFile(path, []byte("original"), 0644); err != nil {
			t.Fatalf("Failed to create existing file: %v", err)
		}

		err := w.Write(path, []byte("new content"))
		if err != nil {
			t.Fatalf("Write() error = %v", err)
		}

		content, _ := os.ReadFile(path)
		if string(content) != "new content" {
			t.Errorf("Write() content = %q, want 'new content'", string(content))
		}
	})

	t.Run("write skips unchanged content", func(t *testing.T) {
		w := New(true)
		path := filepath.Join(tmpDir, "unchanged.txt")

		// Create existing file
		originalContent := []byte("same content")
		if err := os.WriteFile(path, originalContent, 0644); err != nil {
			t.Fatalf("Failed to create existing file: %v", err)
		}

		// Get original modification time
		origInfo, _ := os.Stat(path)
		origModTime := origInfo.ModTime()

		// Wait a bit to ensure time would change if file is written
		time.Sleep(10 * time.Millisecond)

		// Write same content
		err := w.Write(path, originalContent)
		if err != nil {
			t.Fatalf("Write() error = %v", err)
		}

		// Check modification time unchanged
		newInfo, _ := os.Stat(path)
		if !newInfo.ModTime().Equal(origModTime) {
			t.Error("Write() should not modify file when content unchanged")
		}
	})

	t.Run("write updates changed content", func(t *testing.T) {
		w := New(true)
		path := filepath.Join(tmpDir, "changed.txt")

		// Create existing file
		if err := os.WriteFile(path, []byte("original"), 0644); err != nil {
			t.Fatalf("Failed to create existing file: %v", err)
		}

		// Get original modification time
		origInfo, _ := os.Stat(path)
		origModTime := origInfo.ModTime()

		// Wait a bit
		time.Sleep(10 * time.Millisecond)

		// Write different content
		err := w.Write(path, []byte("changed"))
		if err != nil {
			t.Fatalf("Write() error = %v", err)
		}

		// Check content changed
		content, _ := os.ReadFile(path)
		if string(content) != "changed" {
			t.Errorf("Write() content = %q, want 'changed'", string(content))
		}

		// Check modification time changed
		newInfo, _ := os.Stat(path)
		if newInfo.ModTime().Equal(origModTime) {
			t.Error("Write() should modify file when content changed")
		}
	})
}

func TestWriteString(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "render-output-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	w := New(false)
	path := filepath.Join(tmpDir, "string.txt")

	err = w.WriteString(path, "string content")
	if err != nil {
		t.Fatalf("WriteString() error = %v", err)
	}

	content, _ := os.ReadFile(path)
	if string(content) != "string content" {
		t.Errorf("WriteString() content = %q", string(content))
	}
}

func TestCopy(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "render-output-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	t.Run("copy file", func(t *testing.T) {
		w := New(false)
		src := filepath.Join(tmpDir, "source.txt")
		dst := filepath.Join(tmpDir, "dest.txt")

		// Create source file
		if err := os.WriteFile(src, []byte("source content"), 0644); err != nil {
			t.Fatalf("Failed to create source file: %v", err)
		}

		err := w.Copy(src, dst)
		if err != nil {
			t.Fatalf("Copy() error = %v", err)
		}

		content, _ := os.ReadFile(dst)
		if string(content) != "source content" {
			t.Errorf("Copy() content = %q", string(content))
		}
	})

	t.Run("copy non-existent file", func(t *testing.T) {
		w := New(false)
		src := filepath.Join(tmpDir, "nonexistent.txt")
		dst := filepath.Join(tmpDir, "dest2.txt")

		err := w.Copy(src, dst)
		if err == nil {
			t.Error("Copy() should return error for non-existent source")
		}
	})

	t.Run("copy fails without force", func(t *testing.T) {
		w := New(false)
		src := filepath.Join(tmpDir, "src2.txt")
		dst := filepath.Join(tmpDir, "dst2.txt")

		// Create both files
		if err := os.WriteFile(src, []byte("source"), 0644); err != nil {
			t.Fatalf("Failed to create source file: %v", err)
		}
		if err := os.WriteFile(dst, []byte("dest"), 0644); err != nil {
			t.Fatalf("Failed to create dest file: %v", err)
		}

		err := w.Copy(src, dst)
		if err == nil {
			t.Error("Copy() should return error when dest exists and force=false")
		}
	})

	t.Run("copy succeeds with force", func(t *testing.T) {
		w := New(true)
		src := filepath.Join(tmpDir, "src3.txt")
		dst := filepath.Join(tmpDir, "dst3.txt")

		// Create both files
		if err := os.WriteFile(src, []byte("new source"), 0644); err != nil {
			t.Fatalf("Failed to create source file: %v", err)
		}
		if err := os.WriteFile(dst, []byte("old dest"), 0644); err != nil {
			t.Fatalf("Failed to create dest file: %v", err)
		}

		err := w.Copy(src, dst)
		if err != nil {
			t.Fatalf("Copy() error = %v", err)
		}

		content, _ := os.ReadFile(dst)
		if string(content) != "new source" {
			t.Errorf("Copy() content = %q", string(content))
		}
	})
}

func TestExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "render-output-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	path := filepath.Join(tmpDir, "exists.txt")

	if Exists(path) {
		t.Error("Exists() should return false for non-existent file")
	}

	if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if !Exists(path) {
		t.Error("Exists() should return true for existing file")
	}
}

func TestIsDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "render-output-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	if !IsDir(tmpDir) {
		t.Error("IsDir() should return true for directory")
	}

	file := filepath.Join(tmpDir, "file.txt")
	if err := os.WriteFile(file, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if IsDir(file) {
		t.Error("IsDir() should return false for file")
	}

	if IsDir(filepath.Join(tmpDir, "nonexistent")) {
		t.Error("IsDir() should return false for non-existent path")
	}
}

func TestWriteIfNotExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "render-output-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	t.Run("creates file if not exists", func(t *testing.T) {
		w := New(false)
		path := filepath.Join(tmpDir, "newfile.txt")
		content := []byte("new content")

		skipped, err := w.WriteIfNotExists(path, content, 0644)
		if err != nil {
			t.Fatalf("WriteIfNotExists error = %v", err)
		}
		if skipped {
			t.Error("WriteIfNotExists should not skip for new file")
		}

		written, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("Failed to read written file: %v", err)
		}
		if string(written) != "new content" {
			t.Errorf("Content = %q, want %q", string(written), "new content")
		}
	})

	t.Run("skips if file exists", func(t *testing.T) {
		w := New(false)
		path := filepath.Join(tmpDir, "existing.txt")

		// Create existing file
		if err := os.WriteFile(path, []byte("original"), 0644); err != nil {
			t.Fatalf("Failed to create existing file: %v", err)
		}

		skipped, err := w.WriteIfNotExists(path, []byte("new content"), 0644)
		if err != nil {
			t.Fatalf("WriteIfNotExists error = %v", err)
		}
		if !skipped {
			t.Error("WriteIfNotExists should skip for existing file")
		}

		// Verify content unchanged
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}
		if string(content) != "original" {
			t.Errorf("Content = %q, want %q (should be unchanged)", string(content), "original")
		}
	})

	t.Run("creates parent directories", func(t *testing.T) {
		w := New(false)
		path := filepath.Join(tmpDir, "nested", "dir", "file.txt")
		content := []byte("nested content")

		skipped, err := w.WriteIfNotExists(path, content, 0644)
		if err != nil {
			t.Fatalf("WriteIfNotExists error = %v", err)
		}
		if skipped {
			t.Error("WriteIfNotExists should not skip for new file")
		}

		written, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("Failed to read written file: %v", err)
		}
		if string(written) != "nested content" {
			t.Errorf("Content = %q", string(written))
		}
	})
}

func TestCopyIfNotExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "render-output-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	t.Run("copies file if dest not exists", func(t *testing.T) {
		w := New(false)
		src := filepath.Join(tmpDir, "source.txt")
		dst := filepath.Join(tmpDir, "dest.txt")

		if err := os.WriteFile(src, []byte("source content"), 0644); err != nil {
			t.Fatalf("Failed to create source: %v", err)
		}

		skipped, err := w.CopyIfNotExists(src, dst)
		if err != nil {
			t.Fatalf("CopyIfNotExists error = %v", err)
		}
		if skipped {
			t.Error("CopyIfNotExists should not skip for new file")
		}

		content, err := os.ReadFile(dst)
		if err != nil {
			t.Fatalf("Failed to read dest: %v", err)
		}
		if string(content) != "source content" {
			t.Errorf("Content = %q, want %q", string(content), "source content")
		}
	})

	t.Run("skips if dest exists", func(t *testing.T) {
		w := New(false)
		src := filepath.Join(tmpDir, "src2.txt")
		dst := filepath.Join(tmpDir, "dst2.txt")

		if err := os.WriteFile(src, []byte("source"), 0644); err != nil {
			t.Fatalf("Failed to create source: %v", err)
		}
		if err := os.WriteFile(dst, []byte("original dest"), 0644); err != nil {
			t.Fatalf("Failed to create dest: %v", err)
		}

		skipped, err := w.CopyIfNotExists(src, dst)
		if err != nil {
			t.Fatalf("CopyIfNotExists error = %v", err)
		}
		if !skipped {
			t.Error("CopyIfNotExists should skip for existing file")
		}

		// Verify dest unchanged
		content, err := os.ReadFile(dst)
		if err != nil {
			t.Fatalf("Failed to read dest: %v", err)
		}
		if string(content) != "original dest" {
			t.Errorf("Content = %q, want %q (should be unchanged)", string(content), "original dest")
		}
	})
}
