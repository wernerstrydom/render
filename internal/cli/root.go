// Package cli provides the command-line interface for render.
package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "render <template-source> <data-source> -o <output>",
	Short: "Generate files from Go templates and JSON/YAML data",
	Long: `NAME
       render - generate files from Go templates and JSON/YAML data

SYNOPSIS
       render <template-source> <data-source> -o <output> [OPTIONS]

DESCRIPTION
       render uses Go text templates to generate output files from JSON or
       YAML data sources. It automatically detects the rendering mode based
       on the template type and output path.

       Template files with a .tmpl extension are processed as Go templates.
       Other files are copied verbatim (in directory mode). The .tmpl
       extension is stripped from output filenames.

MODES
       render operates in one of three modes, determined automatically:

       File mode
              Single template file renders to a single output file.
              Triggered when template-source is a file and output is a
              static path (no Go template syntax).

       Directory mode
              Template directory renders to mirrored output directory.
              Files with .tmpl extension are rendered; others are copied.
              Directory structure is preserved. Triggered when template-source
              is a directory.

       Each mode
              Template renders once per item extracted from the data.
              Triggered when output path contains Go template syntax ({{...}}).
              Use --item-query to specify which items to iterate over.

OPTIONS
       -o, --output <path>
              Output file or directory path. Required.
              - Static path: file or directory mode
              - Dynamic path with {{...}}: each mode
              - Trailing slash: file rendered into directory

       -f, --force
              Overwrite existing files without prompting. By default,
              render refuses to overwrite existing files.

       --query <jq-expression>
              Transform the entire data using a jq expression before
              rendering. The transformed data becomes the root object
              available to templates. Applied before --item-query.
              Example: --query '.config.database'

       --item-query <jq-expression>
              Extract items from data for iteration in each mode.
              Each extracted item becomes available to the template.
              Enables each mode even without dynamic output path.
              Example: --item-query '.users[] | select(.active)'

       --control <path>
              Explicit path to control file (.render.yaml) for path
              mappings. Disables auto-discovery of control files.

       --dry-run
              Show what files would be written without writing them.
              Useful for previewing output before committing changes.

       --json
              Output results in machine-readable JSON format.
              Each line is a JSON object with file operation details.

EXIT STATUS
       0      Success - all files rendered successfully
       1      Runtime error during template rendering
       2      Usage error - invalid command-line arguments
       3      Input validation error - missing files, malformed data
       4      Permission denied - filesystem permission error
       5      Output conflict - file exists and --force not specified
       6      Safety violation - path traversal or symlink attack detected

SEE ALSO
       Documentation: https://github.com/wernerstrydom/render/tree/main/docs
       Go templates:  https://pkg.go.dev/text/template
       jq manual:     https://jqlang.github.io/jq/manual/`,
	Example: `  # Basic file rendering
  render config.tmpl values.json -o config.yaml

  # Directory of templates
  render ./templates data.yaml -o ./output

  # Generate one file per user
  render user.tmpl users.json --item-query '.users[]' -o '{{.username}}.txt'

  # Transform data before rendering
  render template.tmpl config.json --query '.database' -o db-config.txt

  # Preview without writing
  render ./templates data.json -o ./dist --dry-run

  # Force overwrite existing files
  render config.tmpl values.json -o config.yaml --force

  # Machine-readable output for scripting
  render ./templates data.json -o ./dist --json`,
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
	rootCmd.Flags().StringVar(&flags.query, "query", "", "jq expression to transform data before rendering")
	rootCmd.Flags().StringVar(&flags.itemQuery, "item-query", "", "jq expression to extract items for iteration")

	if err := rootCmd.MarkFlagRequired("output"); err != nil {
		panic(err)
	}
}
