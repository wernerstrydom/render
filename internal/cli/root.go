// Package cli provides the command-line interface for render.
package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "render <template-source> <data-source>",
	Short: "Render text using Go templates and JSON/YAML data",
	Long: `render is a CLI tool that uses Go text templates to generate output files
from JSON or YAML data sources.

It automatically detects the rendering mode based on the template type and
output path:

  - File mode:      Single template file → single output file
  - Directory mode: Template directory → mirrored output directory
  - Each mode:      Template + dynamic output path → multiple output files

The output path determines the mode:
  - Static path (out.txt)       → file or directory mode
  - Dynamic path ({{.id}}.txt)  → each mode (iterates over data)
  - Trailing slash (output/)    → file rendered into directory

Example usage:
  render template.txt.tmpl data.json -o output.txt
  render ./templates data.yaml -o ./output
  render item.tmpl list.json -o "{{.id}}.txt"
  render ./templates data.json -o ./dist --control render.yaml`,
	Args:         cobra.ExactArgs(2),
	RunE:         runRenderCmd,
	SilenceUsage: true,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// Check if the error has an exit code
		if exitErr, ok := err.(*exitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&flags.output, "output", "o", "", "Output path (required)")
	rootCmd.Flags().BoolVarP(&flags.force, "force", "f", false, "Overwrite existing files")
	rootCmd.Flags().BoolVar(&flags.dryRun, "dry-run", false, "Show what would be written without writing")
	rootCmd.Flags().StringVar(&flags.control, "control", "", "Explicit path to control file (no auto-discovery)")
	rootCmd.Flags().BoolVar(&flags.jsonOut, "json", false, "Machine-readable JSON output")

	if err := rootCmd.MarkFlagRequired("output"); err != nil {
		panic(err)
	}
}
