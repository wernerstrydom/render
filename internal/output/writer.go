// Package output provides file writing functionality with change detection.
package output

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Writer handles file output with optional overwrite protection and change detection.
type Writer struct {
	force bool
}

// New creates a new Writer.
func New(force bool) *Writer {
	return &Writer{force: force}
}

// Write writes content to a file.
// If force is false and the file exists, it returns an error.
// If the file exists and content is unchanged, it skips writing to preserve timestamps.
func (w *Writer) Write(path string, content []byte) error {
	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Check if file exists
	if info, err := os.Stat(path); err == nil {
		if !info.Mode().IsRegular() {
			return fmt.Errorf("path exists but is not a regular file: %s", path)
		}

		// File exists - check if we should overwrite
		if !w.force {
			return fmt.Errorf("file already exists (use --force to overwrite): %s", path)
		}

		// Check if content has changed
		existingContent, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read existing file %s: %w", path, err)
		}

		if bytes.Equal(existingContent, content) {
			// Content unchanged, skip writing
			return nil
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat file %s: %w", path, err)
	}

	// Write the file
	if err := os.WriteFile(path, content, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}

// WriteString writes string content to a file.
func (w *Writer) WriteString(path string, content string) error {
	return w.Write(path, []byte(content))
}

// Copy copies a file from src to dst, preserving source file permissions.
// If force is false and the destination exists, it returns an error.
// If the destination exists and content is unchanged, it skips copying.
func (w *Writer) Copy(src, dst string) error {
	// Get source file info for permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source file %s: %w", src, err)
	}

	// Read source file
	content, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read source file %s: %w", src, err)
	}

	// Use WriteWithPerm to preserve source permissions
	return w.WriteWithPerm(dst, content, srcInfo.Mode().Perm())
}

// WriteWithPerm writes content to a file with specified permissions.
func (w *Writer) WriteWithPerm(path string, content []byte, perm os.FileMode) error {
	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Check if file exists
	if info, err := os.Stat(path); err == nil {
		if !info.Mode().IsRegular() {
			return fmt.Errorf("path exists but is not a regular file: %s", path)
		}

		// File exists - check if we should overwrite
		if !w.force {
			return fmt.Errorf("file already exists (use --force to overwrite): %s", path)
		}

		// Check if content has changed
		existingContent, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read existing file %s: %w", path, err)
		}

		if bytes.Equal(existingContent, content) {
			// Content unchanged, skip writing
			return nil
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat file %s: %w", path, err)
	}

	// Write the file with specified permissions
	if err := os.WriteFile(path, content, perm); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}

// CopyReader copies content from a reader to a file.
func (w *Writer) CopyReader(r io.Reader, dst string) error {
	content, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read content: %w", err)
	}

	return w.Write(dst, content)
}

// WriteIfNotExists writes content to a file only if the file doesn't exist.
// Returns (true, nil) if the file already exists and was skipped.
// Returns (false, nil) if the file was written successfully.
// Returns (false, error) on failure.
func (w *Writer) WriteIfNotExists(path string, content []byte, perm os.FileMode) (skipped bool, err error) {
	// Check if file already exists
	if _, err := os.Stat(path); err == nil {
		// File exists, skip writing
		return true, nil
	} else if !os.IsNotExist(err) {
		return false, fmt.Errorf("failed to stat file %s: %w", path, err)
	}

	// File doesn't exist, create it
	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return false, fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	if err := os.WriteFile(path, content, perm); err != nil {
		return false, fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return false, nil
}

// CopyIfNotExists copies a file only if the destination doesn't exist.
// Returns (true, nil) if the destination already exists and was skipped.
// Returns (false, nil) if the file was copied successfully.
// Returns (false, error) on failure.
func (w *Writer) CopyIfNotExists(src, dst string) (skipped bool, err error) {
	// Check if destination already exists
	if _, err := os.Stat(dst); err == nil {
		// File exists, skip copying
		return true, nil
	} else if !os.IsNotExist(err) {
		return false, fmt.Errorf("failed to stat file %s: %w", dst, err)
	}

	// Get source file info for permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return false, fmt.Errorf("failed to stat source file %s: %w", src, err)
	}

	// Read source file
	content, err := os.ReadFile(src)
	if err != nil {
		return false, fmt.Errorf("failed to read source file %s: %w", src, err)
	}

	// Ensure parent directory exists
	dir := filepath.Dir(dst)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return false, fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write to destination
	if err := os.WriteFile(dst, content, srcInfo.Mode().Perm()); err != nil {
		return false, fmt.Errorf("failed to write file %s: %w", dst, err)
	}

	return false, nil
}

// Exists checks if a file exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDir checks if a path is a directory.
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
