# render Documentation

render is a CLI tool that uses Go text templates to generate output files from JSON or YAML data sources.

## Quick Start

```bash
# Install
go install github.com/wernerstrydom/render/cmd/render@latest

# Basic usage: render a single template
render config.tmpl values.json -o config.yaml

# Render a directory of templates
render ./templates data.yaml -o ./output

# Generate one file per item
render user.tmpl users.json --item-query '.users[]' -o '{{.username}}.txt'
```

## Documentation

### Concepts

- [Rendering Modes](concepts/modes.md) - File, directory, and each mode
- [Go Templates](concepts/templates.md) - Template syntax and features
- [Data Sources](concepts/data-sources.md) - JSON and YAML data handling

### Guides

- [Getting Started](guides/getting-started.md) - First steps with render
- [Directory Mode](guides/directory-mode.md) - Rendering template directories
- [Each Mode](guides/each-mode.md) - Generating multiple files from arrays
- [Query Expressions](guides/query-expressions.md) - Using jq to transform data
- [Control Files](guides/control-files.md) - Configuring path mappings

### Reference

- [CLI Reference](reference/cli.md) - Complete command-line reference
- [Template Functions](reference/functions.md) - All available template functions
- [Exit Codes](reference/exit-codes.md) - Exit codes and error handling

### Examples

- [Config Generation](examples/config-generation.md) - Generating configuration files
- [Code Generation](examples/code-generation.md) - Generating source code
- [Multi-File Output](examples/multi-file-output.md) - Creating multiple files from data

## Key Features

- **Three rendering modes**: File, directory, and each mode for different use cases
- **Go templates**: Full Go text/template support with 80+ helper functions
- **jq queries**: Transform data before rendering with jq expressions
- **Control files**: Configure output path mappings with .render.yaml
- **Safety**: Prevents path traversal attacks and symlink exploits
- **Dry run**: Preview changes before writing files

## Getting Help

```bash
# Show help
render --help

# Generate man pages
render gen man ./man
man ./man/render.1
```
