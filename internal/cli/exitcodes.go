// Package cli provides the command-line interface for render.
package cli

// Exit codes for the render CLI.
const (
	// ExitSuccess indicates the command completed successfully.
	ExitSuccess = 0

	// ExitRuntimeError indicates a runtime error during template rendering.
	ExitRuntimeError = 1

	// ExitUsageError indicates invalid command-line arguments.
	ExitUsageError = 2

	// ExitInputValidation indicates input validation failed (missing files, malformed data).
	ExitInputValidation = 3

	// ExitPermissionDenied indicates a filesystem permission error.
	ExitPermissionDenied = 4

	// ExitOutputConflict indicates output file exists and --force was not specified.
	ExitOutputConflict = 5

	// ExitSafetyViolation indicates a security issue (path traversal, symlinks).
	ExitSafetyViolation = 6
)
