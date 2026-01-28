# CLI Reference

Complete command-line reference for render.

## Synopsis

```
render <template-source> <data-source> -o <output> [OPTIONS]
```

## Arguments

### template-source

Path to a template file or directory.

- **File**: Single template file (e.g., `config.tmpl`)
- **Directory**: Directory containing templates (e.g., `./templates`)

Files with `.tmpl` extension are processed as Go templates. Other files in directories are copied verbatim.

### data-source

Path to a JSON or YAML data file.

- **JSON**: `data.json`, `config.json`
- **YAML**: `data.yaml`, `data.yml`, `config.yaml`

The file format is detected by extension.

## Required Flags

### -o, --output

Output path. Required.

```bash
-o output.txt          # Single file
-o ./output            # Directory
-o ./output/           # Directory (explicit with trailing slash)
-o '{{.name}}.txt'     # Dynamic path (each mode)
```

## Optional Flags

### -f, --force

Overwrite existing files without prompting.

```bash
render config.tmpl data.json -o config.yaml --force
```

Default: `false` (refuse to overwrite existing files)

### --query

Transform data using a jq expression before rendering.

```bash
render tmpl data.json --query '.config.database' -o out.txt
render tmpl data.json --query '.users | map(select(.active))' -o out.txt
```

The transformed data becomes the root object available to templates.

### --item-query

Extract items from data for iteration. Enables each mode.

```bash
render user.tmpl data.json --item-query '.users[]' -o '{{.name}}.txt'
render user.tmpl data.json --item-query '.users[] | select(.active)' -o '{{.name}}.txt'
```

Each extracted item becomes the root data for one template render.

### --control

Explicit path to a control file for path mappings.

```bash
render ./templates data.json -o ./output --control ./custom.yaml
```

Disables auto-discovery of `.render.yaml`, `.render.yml`, and `render.json`.

### --dry-run

Show what files would be written without writing them.

```bash
render ./templates data.json -o ./output --dry-run
```

Output shows planned file operations.

### --json

Output results in machine-readable JSON format.

```bash
render ./templates data.json -o ./output --json
```

Each line is a JSON object:

```json
{"path":"output/config.yaml","action":"write"}
{"path":"output/logo.png","action":"copy"}
```

## Subcommands

### render gen man

Generate man pages.

```bash
render gen man <output-dir>
```

Example:
```bash
render gen man ./man
man ./man/render.1
```

### render gen markdown

Generate markdown documentation.

```bash
render gen markdown <output-dir>
```

Example:
```bash
render gen markdown ./docs/cli
```

### render completion

Generate shell completion scripts.

```bash
render completion bash
render completion zsh
render completion fish
render completion powershell
```

## Exit Codes

| Code | Name | Description |
|------|------|-------------|
| 0 | Success | All files rendered successfully |
| 1 | RuntimeError | Error during template rendering |
| 2 | UsageError | Invalid command-line arguments |
| 3 | InputValidation | Missing files or malformed data |
| 4 | PermissionDenied | Filesystem permission error |
| 5 | OutputConflict | File exists and --force not specified |
| 6 | SafetyViolation | Path traversal or symlink attack detected |

## Environment Variables

render does not use environment variables directly, but templates can access them via shell expansion in data files:

```bash
# In shell
export DB_HOST=localhost
envsubst < data.template.json > data.json
render config.tmpl data.json -o config.yaml
```

## Examples

### Basic File Rendering

```bash
render config.tmpl values.json -o config.yaml
```

### Directory of Templates

```bash
render ./templates data.yaml -o ./output
```

### Generate One File Per Item

```bash
render user.tmpl users.json --item-query '.users[]' -o '{{.username}}.txt'
```

### Transform Data Before Rendering

```bash
render template.tmpl config.json --query '.database' -o db-config.txt
```

### Preview Without Writing

```bash
render ./templates data.json -o ./dist --dry-run
```

### Force Overwrite

```bash
render config.tmpl values.json -o config.yaml --force
```

### Machine-Readable Output

```bash
render ./templates data.json -o ./dist --json
```

### Filter and Render

```bash
render user.tmpl data.json \
  --query '.response.data' \
  --item-query '.users[] | select(.verified)' \
  -o 'users/{{.id}}.txt'
```

### With Custom Control File

```bash
render ./templates data.json -o ./output --control ./mappings.yaml
```

## See Also

- [Rendering Modes](../concepts/modes.md)
- [Template Functions](functions.md)
- [Exit Codes](exit-codes.md)
- [jq Manual](https://jqlang.github.io/jq/manual/)
- [Go Templates](https://pkg.go.dev/text/template)
