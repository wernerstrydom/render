package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wernerstrydom/render/internal/config"
	"github.com/wernerstrydom/render/internal/data"
	"github.com/wernerstrydom/render/internal/engine"
	"github.com/wernerstrydom/render/internal/output"
	"github.com/wernerstrydom/render/internal/render"
)

// renderFlags holds the command-line flags for the render command.
type renderFlags struct {
	output  string
	force   bool
	dryRun  bool
	control string
	jsonOut bool
}

var flags renderFlags

// renderResult represents the JSON output format.
type renderResult struct {
	Status string       `json:"status"`
	Files  []fileAction `json:"files,omitempty"`
	Error  string       `json:"error,omitempty"`
}

type fileAction struct {
	Path   string `json:"path"`
	Action string `json:"action"`
}

// renderMode represents the detected rendering mode.
type renderMode int

const (
	modeFile renderMode = iota
	modeFileIntoDir
	modeDirectory
	modeEachFile
	modeEachDirectory
)

func (m renderMode) String() string {
	switch m {
	case modeFile:
		return "file"
	case modeFileIntoDir:
		return "file-into-dir"
	case modeDirectory:
		return "directory"
	case modeEachFile:
		return "each-file"
	case modeEachDirectory:
		return "each-directory"
	default:
		return "unknown"
	}
}

// runRenderCmd executes the unified render command.
func runRenderCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return &exitError{
			code: ExitUsageError,
			msg:  "requires exactly 2 arguments: <template-source> <data-source>",
		}
	}

	templatePath := args[0]
	dataPath := args[1]

	// Validate output flag
	if flags.output == "" {
		return &exitError{
			code: ExitUsageError,
			msg:  "required flag --output/-o not set",
		}
	}

	// Check for symlinks in template source
	if err := checkForSymlinks(templatePath); err != nil {
		return &exitError{code: ExitSafetyViolation, msg: err.Error()}
	}

	// Load data
	d, err := data.Load(dataPath)
	if err != nil {
		return &exitError{
			code: ExitInputValidation,
			msg:  fmt.Sprintf("failed to load data: %v", err),
		}
	}

	// Determine template type
	tmplInfo, err := os.Stat(templatePath)
	if err != nil {
		return &exitError{
			code: ExitInputValidation,
			msg:  fmt.Sprintf("failed to access template: %v", err),
		}
	}

	// Determine rendering mode
	mode := inferMode(tmplInfo.IsDir(), flags.output, d)

	// Execute based on mode
	switch mode {
	case modeFile:
		return executeFileMode(cmd, templatePath, d)
	case modeFileIntoDir:
		return executeFileIntoDirMode(cmd, templatePath, d)
	case modeDirectory:
		return executeDirectoryMode(cmd, templatePath, d)
	case modeEachFile:
		return executeEachFileMode(cmd, templatePath, d)
	case modeEachDirectory:
		return executeEachDirectoryMode(cmd, templatePath, d)
	default:
		return &exitError{
			code: ExitRuntimeError,
			msg:  fmt.Sprintf("unknown mode: %v", mode),
		}
	}
}

// inferMode determines the rendering mode based on inputs.
func inferMode(isDir bool, outputPath string, _ any) renderMode {
	isDynamic := strings.Contains(outputPath, "{{") && strings.Contains(outputPath, "}}")
	hasTrailingSlash := strings.HasSuffix(outputPath, "/") || strings.HasSuffix(outputPath, string(os.PathSeparator))

	if isDir {
		if isDynamic {
			return modeEachDirectory
		}
		return modeDirectory
	}

	// File template
	if isDynamic {
		return modeEachFile
	}
	if hasTrailingSlash {
		return modeFileIntoDir
	}
	return modeFile
}

// executeFileMode renders a single template file to a single output file.
func executeFileMode(cmd *cobra.Command, templatePath string, d any) error {
	eng := engine.New()

	// Read template
	tmplContent, err := os.ReadFile(templatePath)
	if err != nil {
		return &exitError{
			code: ExitInputValidation,
			msg:  fmt.Sprintf("failed to read template: %v", err),
		}
	}

	// Render template
	result, err := eng.RenderString(string(tmplContent), d)
	if err != nil {
		return &exitError{
			code: ExitRuntimeError,
			msg:  fmt.Sprintf("failed to render template: %v", err),
		}
	}

	// Check for collision
	collision, err := checkCollision(flags.output, []byte(result))
	if err != nil {
		return err
	}

	// Dry run - just report what would happen
	if flags.dryRun {
		return reportDryRun(cmd, []fileAction{{Path: flags.output, Action: "create"}})
	}

	// Skip if content is identical (idempotency)
	if collision == collisionIdentical {
		return reportSuccess(cmd, []fileAction{{Path: flags.output, Action: "skipped (identical)"}})
	}

	// Write output
	writer := output.New(flags.force)
	if err := writer.WriteString(flags.output, result); err != nil {
		return wrapWriteError(err, flags.output)
	}

	return reportSuccess(cmd, []fileAction{{Path: flags.output, Action: "created"}})
}

// executeFileIntoDirMode renders a template file into a target directory.
func executeFileIntoDirMode(cmd *cobra.Command, templatePath string, d any) error {
	eng := engine.New()

	// Read template
	tmplContent, err := os.ReadFile(templatePath)
	if err != nil {
		return &exitError{
			code: ExitInputValidation,
			msg:  fmt.Sprintf("failed to read template: %v", err),
		}
	}

	// Render template
	result, err := eng.RenderString(string(tmplContent), d)
	if err != nil {
		return &exitError{
			code: ExitRuntimeError,
			msg:  fmt.Sprintf("failed to render template: %v", err),
		}
	}

	// Determine output filename: strip .tmpl if present
	baseName := strings.TrimSuffix(filepath.Base(templatePath), ".tmpl")
	outputPath := filepath.Join(strings.TrimSuffix(flags.output, "/"), baseName)

	// Check for collision
	collision, err := checkCollision(outputPath, []byte(result))
	if err != nil {
		return err
	}

	if flags.dryRun {
		return reportDryRun(cmd, []fileAction{{Path: outputPath, Action: "create"}})
	}

	// Skip if content is identical (idempotency)
	if collision == collisionIdentical {
		return reportSuccess(cmd, []fileAction{{Path: outputPath, Action: "skipped (identical)"}})
	}

	// Write output
	writer := output.New(flags.force)
	if err := writer.WriteString(outputPath, result); err != nil {
		return wrapWriteError(err, outputPath)
	}

	return reportSuccess(cmd, []fileAction{{Path: outputPath, Action: "created"}})
}

// executeDirectoryMode renders a directory of templates.
func executeDirectoryMode(cmd *cobra.Command, templatePath string, d any) error {
	eng := engine.New()

	// Load render config
	var cfg *config.ParsedConfig
	var err error
	if flags.control != "" {
		cfg, err = config.LoadFile(flags.control, templatePath)
	} else {
		cfg, err = config.Load(templatePath)
	}
	if err != nil {
		return &exitError{
			code: ExitInputValidation,
			msg:  fmt.Sprintf("failed to load render config: %v", err),
		}
	}

	// Check for symlinks in template directory
	if err := checkDirForSymlinks(templatePath); err != nil {
		return &exitError{code: ExitSafetyViolation, msg: err.Error()}
	}

	// Collect all outputs
	plan, err := render.Collect(render.CollectConfig{
		TemplateDir: templatePath,
		OutputDir:   flags.output,
		Data:        d,
		Config:      cfg,
		Engine:      eng,
	})
	if err != nil {
		return &exitError{
			code: ExitRuntimeError,
			msg:  fmt.Sprintf("failed to collect outputs: %v", err),
		}
	}

	// Validate
	if errs := plan.Validate(); len(errs) > 0 {
		for _, e := range errs {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", e)
		}
		return &exitError{
			code: ExitRuntimeError,
			msg:  fmt.Sprintf("validation failed with %d error(s)", len(errs)),
		}
	}

	// Check for collisions (skipping identical content and no-overwrite files)
	for _, out := range plan.Outputs {
		// Skip collision check for no-overwrite files - they're allowed to exist
		if !out.Overwrite {
			continue
		}
		_, err := checkCollision(out.OutputPath, out.Content)
		if err != nil {
			return err
		}
	}

	if flags.dryRun {
		actions := make([]fileAction, len(plan.Outputs))
		for i, out := range plan.Outputs {
			action := "render"
			if out.CopyFrom != "" {
				action = "copy"
			}
			if !out.Overwrite {
				// Check if file exists to determine skip action
				if _, err := os.Stat(out.OutputPath); err == nil {
					action = "skip (exists, no-overwrite)"
				}
			}
			actions[i] = fileAction{Path: out.OutputPath, Action: action}
		}
		return reportDryRun(cmd, actions)
	}

	// Execute
	result, err := plan.Execute(output.New(flags.force))
	if err != nil {
		return wrapWriteError(err, "")
	}

	// Report what was written
	actions := make([]fileAction, len(plan.Outputs))
	for i, out := range plan.Outputs {
		if result.Skipped[out.OutputPath] {
			actions[i] = fileAction{Path: out.OutputPath, Action: "skipped (exists, no-overwrite)"}
		} else if out.CopyFrom != "" {
			actions[i] = fileAction{Path: out.OutputPath, Action: "copied"}
		} else {
			actions[i] = fileAction{Path: out.OutputPath, Action: "rendered"}
		}
	}

	return reportSuccess(cmd, actions)
}

// executeEachFileMode renders a template for each item in an array.
func executeEachFileMode(cmd *cobra.Command, templatePath string, d any) error {
	eng := engine.New()

	// Read template
	tmplContent, err := os.ReadFile(templatePath)
	if err != nil {
		return &exitError{
			code: ExitInputValidation,
			msg:  fmt.Sprintf("failed to read template: %v", err),
		}
	}

	// Get items to iterate over
	items := getIterableItems(d)

	// Pre-flight: collect all outputs to check for collisions
	type plannedOutput struct {
		path    string
		content string
	}
	planned := make([]plannedOutput, 0, len(items))
	seenPaths := make(map[string]int) // path -> index in items

	for i, item := range items {
		// Render output path
		outPath, err := eng.RenderString(flags.output, item)
		if err != nil {
			return &exitError{
				code: ExitRuntimeError,
				msg:  fmt.Sprintf("failed to render output path: %v", err),
			}
		}
		outPath = strings.TrimSpace(outPath)

		// Validate output path
		if err := validateOutputPath(outPath); err != nil {
			return &exitError{code: ExitSafetyViolation, msg: err.Error()}
		}

		// Check for internal collision
		if prevIdx, exists := seenPaths[outPath]; exists {
			return &exitError{
				code: ExitRuntimeError,
				msg:  fmt.Sprintf("internal collision: items at index %d and %d both produce path %q", prevIdx, i, outPath),
			}
		}
		seenPaths[outPath] = i

		// Render template
		result, err := eng.RenderString(string(tmplContent), item)
		if err != nil {
			return &exitError{
				code: ExitRuntimeError,
				msg:  fmt.Sprintf("failed to render template: %v", err),
			}
		}

		planned = append(planned, plannedOutput{path: outPath, content: result})
	}

	// Check for filesystem collisions and track which files can be skipped
	skipMap := make(map[int]bool)
	for i, p := range planned {
		collision, err := checkCollision(p.path, []byte(p.content))
		if err != nil {
			return err
		}
		if collision == collisionIdentical {
			skipMap[i] = true
		}
	}

	if flags.dryRun {
		actions := make([]fileAction, len(planned))
		for i, p := range planned {
			actions[i] = fileAction{Path: p.path, Action: "create"}
		}
		return reportDryRun(cmd, actions)
	}

	// Write all outputs (skipping identical content)
	writer := output.New(flags.force)
	var actions []fileAction
	for i, p := range planned {
		if skipMap[i] {
			actions = append(actions, fileAction{Path: p.path, Action: "skipped (identical)"})
			continue
		}
		if err := writer.WriteString(p.path, p.content); err != nil {
			return wrapWriteError(err, p.path)
		}
		actions = append(actions, fileAction{Path: p.path, Action: "created"})
	}

	return reportSuccess(cmd, actions)
}

// executeEachDirectoryMode renders a directory template for each item in an array.
func executeEachDirectoryMode(cmd *cobra.Command, templatePath string, d any) error {
	eng := engine.New()

	// Load render config
	var cfg *config.ParsedConfig
	var err error
	if flags.control != "" {
		cfg, err = config.LoadFile(flags.control, templatePath)
	} else {
		cfg, err = config.Load(templatePath)
	}
	if err != nil {
		return &exitError{
			code: ExitInputValidation,
			msg:  fmt.Sprintf("failed to load render config: %v", err),
		}
	}

	// Check for symlinks
	if err := checkDirForSymlinks(templatePath); err != nil {
		return &exitError{code: ExitSafetyViolation, msg: err.Error()}
	}

	// Get items to iterate over
	items := getIterableItems(d)

	// Pre-flight: collect all outputs to check for collisions
	type plannedDir struct {
		outputDir string
		plan      *render.Plan
	}
	allPlanned := make([]plannedDir, 0, len(items))
	seenPaths := make(map[string]int) // output path -> item index

	for i, item := range items {
		// Render output directory path
		outDir, err := eng.RenderString(flags.output, item)
		if err != nil {
			return &exitError{
				code: ExitRuntimeError,
				msg:  fmt.Sprintf("failed to render output path: %v", err),
			}
		}
		outDir = strings.TrimSpace(outDir)

		// Validate output path
		if err := validateOutputPath(outDir); err != nil {
			return &exitError{code: ExitSafetyViolation, msg: err.Error()}
		}

		// Collect outputs for this item
		plan, err := render.Collect(render.CollectConfig{
			TemplateDir: templatePath,
			OutputDir:   outDir,
			Data:        item,
			Config:      cfg,
			Engine:      eng,
		})
		if err != nil {
			return &exitError{
				code: ExitRuntimeError,
				msg:  fmt.Sprintf("failed to collect outputs: %v", err),
			}
		}

		// Check for internal collisions across all items
		for _, out := range plan.Outputs {
			if prevIdx, exists := seenPaths[out.OutputPath]; exists {
				return &exitError{
					code: ExitRuntimeError,
					msg:  fmt.Sprintf("internal collision: items at index %d and %d both produce path %q", prevIdx, i, out.OutputPath),
				}
			}
			seenPaths[out.OutputPath] = i
		}

		// Validate within item
		if errs := plan.Validate(); len(errs) > 0 {
			return &exitError{
				code: ExitRuntimeError,
				msg:  fmt.Sprintf("validation failed: %v", errs[0]),
			}
		}

		allPlanned = append(allPlanned, plannedDir{outputDir: outDir, plan: plan})
	}

	// Check for filesystem collisions (skipping identical content and no-overwrite files)
	for _, pd := range allPlanned {
		for _, out := range pd.plan.Outputs {
			// Skip collision check for no-overwrite files - they're allowed to exist
			if !out.Overwrite {
				continue
			}
			_, err := checkCollision(out.OutputPath, out.Content)
			if err != nil {
				return err
			}
		}
	}

	if flags.dryRun {
		var actions []fileAction
		for _, pd := range allPlanned {
			for _, out := range pd.plan.Outputs {
				action := "render"
				if out.CopyFrom != "" {
					action = "copy"
				}
				if !out.Overwrite {
					// Check if file exists to determine skip action
					if _, err := os.Stat(out.OutputPath); err == nil {
						action = "skip (exists, no-overwrite)"
					}
				}
				actions = append(actions, fileAction{Path: out.OutputPath, Action: action})
			}
		}
		return reportDryRun(cmd, actions)
	}

	// Execute all plans
	writer := output.New(flags.force)
	var actions []fileAction
	for _, pd := range allPlanned {
		result, err := pd.plan.Execute(writer)
		if err != nil {
			return wrapWriteError(err, "")
		}
		for _, out := range pd.plan.Outputs {
			if result.Skipped[out.OutputPath] {
				actions = append(actions, fileAction{Path: out.OutputPath, Action: "skipped (exists, no-overwrite)"})
			} else if out.CopyFrom != "" {
				actions = append(actions, fileAction{Path: out.OutputPath, Action: "copied"})
			} else {
				actions = append(actions, fileAction{Path: out.OutputPath, Action: "rendered"})
			}
		}
	}

	return reportSuccess(cmd, actions)
}

// getIterableItems returns items to iterate over.
// If data is an array, returns the array elements.
// If data is an object, wraps it in a single-element array.
func getIterableItems(d any) []any {
	if arr, ok := d.([]any); ok {
		return arr
	}
	return []any{d}
}

// validateOutputPath checks that an output path is safe.
func validateOutputPath(path string) error {
	separator := func(r rune) bool {
		return r == '/' || r == '\\'
	}
	parts := strings.FieldsFunc(path, separator)
	if slices.Contains(parts, "..") {
		return fmt.Errorf("security error: output path contains directory traversal: %s", path)
	}
	return nil
}

// collisionResult represents the result of a collision check.
type collisionResult int

const (
	collisionNone      collisionResult = iota // No collision, proceed with write
	collisionIdentical                        // File exists but content is identical, skip write
)

// checkCollision checks if a file would collide with an existing file.
// Returns (collisionIdentical, nil) if file exists with identical content (skip write).
// Returns (collisionNone, nil) if file doesn't exist or force is enabled (proceed with write).
// Returns (_, error) if file exists with different content and force not enabled.
func checkCollision(path string, content []byte) (collisionResult, error) {
	if !flags.force {
		if info, err := os.Stat(path); err == nil {
			if info.Mode().IsRegular() {
				// Check if content is identical (idempotency)
				existing, err := os.ReadFile(path)
				if err != nil {
					return collisionNone, &exitError{
						code: ExitOutputConflict,
						msg:  fmt.Sprintf("file already exists (use --force to overwrite): %s", path),
					}
				}
				if string(existing) == string(content) {
					// Content is identical, skip the write
					return collisionIdentical, nil
				}
				return collisionNone, &exitError{
					code: ExitOutputConflict,
					msg:  fmt.Sprintf("file already exists (use --force to overwrite): %s", path),
				}
			}
		}
	}
	return collisionNone, nil
}

// checkForSymlinks checks if a path is a symlink.
func checkForSymlinks(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return nil // File doesn't exist, which is fine
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("template source contains symlink: %s", path)
	}
	return nil
}

// checkDirForSymlinks recursively checks a directory for symlinks.
func checkDirForSymlinks(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check lstat to detect symlinks
		linfo, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if linfo.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("template source contains symlink: %s", path)
		}
		return nil
	})
}

// wrapWriteError converts a write error to an appropriate exit error.
func wrapWriteError(err error, _ string) error {
	errStr := err.Error()
	if strings.Contains(errStr, "permission denied") {
		return &exitError{
			code: ExitPermissionDenied,
			msg:  errStr,
		}
	}
	if strings.Contains(errStr, "already exists") || strings.Contains(errStr, "force") {
		return &exitError{
			code: ExitOutputConflict,
			msg:  errStr,
		}
	}
	return &exitError{
		code: ExitRuntimeError,
		msg:  errStr,
	}
}

// reportDryRun reports what would be done in dry-run mode.
func reportDryRun(cmd *cobra.Command, actions []fileAction) error {
	if flags.jsonOut {
		result := renderResult{
			Status: "dry-run",
			Files:  actions,
		}
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Dry run - would perform:")
	for _, a := range actions {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  [%s] %s\n", a.Action, a.Path)
	}
	return nil
}

// reportSuccess reports successful completion.
func reportSuccess(cmd *cobra.Command, actions []fileAction) error {
	if flags.jsonOut {
		result := renderResult{
			Status: "success",
			Files:  actions,
		}
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	for _, a := range actions {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", capitalizeFirst(a.Action), a.Path)
	}
	return nil
}

// exitError represents an error with a specific exit code.
type exitError struct {
	code int
	msg  string
}

func (e *exitError) Error() string {
	return e.msg
}

// ExitCode returns the exit code for this error.
func (e *exitError) ExitCode() int {
	return e.code
}

// capitalizeFirst capitalizes the first letter of a string.
func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
