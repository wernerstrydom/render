// Package cli provides the command-line interface for render.
package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate documentation for render",
	Long: `Generate documentation for render in various formats.

Available subcommands:
  man       Generate man pages
  markdown  Generate markdown documentation`,
}

var genManCmd = &cobra.Command{
	Use:   "man <output-dir>",
	Short: "Generate man pages",
	Long: `Generate man pages for render.

The output directory will be created if it doesn't exist.
Man pages are generated in troff format suitable for the
man(1) command.

After generation, you can view the man page with:
  man ./output-dir/render.1`,
	Example: `  # Generate man pages to ./man directory
  render gen man ./man

  # View the generated man page
  man ./man/render.1`,
	Args: cobra.ExactArgs(1),
	RunE: runGenMan,
}

var genMarkdownCmd = &cobra.Command{
	Use:   "markdown <output-dir>",
	Short: "Generate markdown documentation",
	Long: `Generate markdown documentation for render.

The output directory will be created if it doesn't exist.
Documentation is generated as markdown files suitable for
viewing on GitHub or other markdown renderers.`,
	Example: `  # Generate markdown docs to ./docs/cli directory
  render gen markdown ./docs/cli`,
	Args: cobra.ExactArgs(1),
	RunE: runGenMarkdown,
}

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.AddCommand(genManCmd)
	genCmd.AddCommand(genMarkdownCmd)
}

func runGenMan(cmd *cobra.Command, args []string) error {
	outputDir := args[0]

	// Create output directory
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate man pages
	header := &doc.GenManHeader{
		Title:   "RENDER",
		Section: "1",
		Date:    timePtr(time.Now()),
		Source:  "render",
		Manual:  "User Commands",
	}

	if err := doc.GenManTree(rootCmd, header, outputDir); err != nil {
		return fmt.Errorf("failed to generate man pages: %w", err)
	}

	fmt.Printf("Man pages generated in %s\n", outputDir)
	fmt.Printf("View with: man %s\n", filepath.Join(outputDir, "render.1"))
	return nil
}

func runGenMarkdown(cmd *cobra.Command, args []string) error {
	outputDir := args[0]

	// Create output directory
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate markdown docs
	if err := doc.GenMarkdownTree(rootCmd, outputDir); err != nil {
		return fmt.Errorf("failed to generate markdown docs: %w", err)
	}

	fmt.Printf("Markdown documentation generated in %s\n", outputDir)
	return nil
}

func timePtr(t time.Time) *time.Time {
	return &t
}
