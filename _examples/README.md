# Render Examples

This directory contains examples demonstrating how to use the `render` CLI tool
for various code generation scenarios. Each example shows a practical use case
where render provides significant advantages over manual file creation.

## Why Use Render?

When an AI assistant needs to generate repetitive code structures, using `render`
offers several advantages over writing files individually:

1. **Speed**: Generate dozens of files in a single command
2. **Consistency**: Templates enforce uniform patterns across files
3. **Maintainability**: Update the template once, regenerate everything
4. **Reduced Errors**: Less opportunity for copy-paste mistakes
5. **Context Efficiency**: Smaller token usage than writing each file

## Examples

| Example | Description | Key Features |
|---------|-------------|--------------|
| [go-workspace-grpc](./go-workspace-grpc/) | Multi-module Go workspace with gRPC | Each mode, path transformation |
| [maven-multimodule-grpc](./maven-multimodule-grpc/) | Maven multi-module Java project | Directory mode, package paths |
| [go-driver](./go-driver/) | Driver interface pattern | Two-pass, interface + implementations |
| [go-workspace-makefile](./go-workspace-makefile/) | Makefile for Go workspaces | File mode, loops, conditionals |
| [go-cli-scaffold](./go-cli-scaffold/) | CLI with Cobra/Viper | Command structure, flags |
| [python-fastapi-service](./python-fastapi-service/) | Python FastAPI Service | Two-pass rendering, Models/Routers |
| [docker-compose](./docker-compose/) | Docker Compose stacks | Complex config from simple data |
| [kubernetes-manifests](./kubernetes-manifests/) | K8s manifests per environment | Nested loops, environment config |
| [github-actions](./github-actions/) | CI/CD workflow generation | Conditionals, matrix builds |
| [terraform-modules](./terraform-modules/) | Terraform infrastructure | Modules, environments |

## Quick Start

```bash
# Navigate to any example directory
cd examples/go-cli-scaffold

# Preview what would be generated
render templates cli.yaml -o output --dry-run

# Generate the files
render templates cli.yaml -o output

# Explore the generated structure
tree output/
```

## Template Syntax

Render uses Go's `text/template` syntax with additional functions:

```go
// Variable interpolation
{{ .name }}

// Nested access
{{ .config.database.host }}

// Pipes and functions
{{ .name | pascalCase }}
{{ .items | join ", " }}

// Conditionals
{{ if .enabled }}...{{ end }}
{{ if eq .type "service" }}...{{ end }}

// Loops
{{ range .services }}
  {{ .name }}
{{ end }}

// With context
{{ with .database }}
  Host: {{ .host }}
{{ end }}
```

## Available Functions

Render includes 100+ template functions:

- **Casing**: `camelCase`, `pascalCase`, `snakeCase`, `kebabCase`
- **String**: `lower`, `upper`, `title`, `trim`, `replace`
- **Collections**: `first`, `last`, `join`, `split`, `len`
- **Math**: `add`, `sub`, `mul`, `div`
- **Logic**: `eq`, `ne`, `lt`, `gt`, `and`, `or`, `not`
- **JSON**: `toJson`, `toPrettyJson`, `fromJson`

See the [main documentation](../README.md) for the complete function reference.

## Modes of Operation

### File Mode
Single template to single output file.
```bash
render template.tmpl data.yaml -o output.txt
```

### Directory Mode
Template directory to output directory, preserving structure.
```bash
render templates/ data.yaml -o output/
```

### Each Mode
Template rendered for each array element with dynamic paths.
```bash
render template.tmpl items.yaml -o "{{.name}}/file.txt"
```

## Path Transformation

Use `.render.yaml` in your template directory to transform output paths:

```yaml
paths:
  # Rename files dynamically
  "model.go.tmpl": "{{ .name | snakeCase }}.go"

  # Transform directory prefixes
  "src/main/java": "src/main/java/{{ .package | replace \".\" \"/\" }}"
```

## When to Use Render vs. Write Tool

**Use Render when:**
- Generating multiple similar files
- Creating environment-specific configurations
- Scaffolding project structures
- Maintaining templates for reuse

**Use Write Tool when:**
- Creating a single unique file
- Making targeted edits to existing files
- Content doesn't follow a repeatable pattern
