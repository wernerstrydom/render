# Directory Mode Guide

Directory mode renders an entire template directory to an output directory, preserving the structure.

## Basic Usage

```bash
render ./templates data.json -o ./output
```

## How It Works

1. render walks the template directory recursively
2. Files with `.tmpl` extension are processed as templates
3. Other files are copied verbatim
4. Directory structure is preserved
5. The `.tmpl` extension is stripped from output filenames

## Example

Given this template structure:

```
templates/
  config/
    app.yaml.tmpl
    database.yaml.tmpl
  scripts/
    setup.sh.tmpl
    deploy.sh
  assets/
    logo.png
  .render.yaml
```

And data:

```json
{
  "app": {
    "name": "myapp",
    "port": 8080
  },
  "database": {
    "host": "localhost",
    "port": 5432
  }
}
```

Running:

```bash
render ./templates data.json -o ./output
```

Produces:

```
output/
  config/
    app.yaml
    database.yaml
  scripts/
    setup.sh
    deploy.sh
  assets/
    logo.png
```

Notes:
- `.tmpl` files are rendered and have extension stripped
- `deploy.sh` is copied without modification (no `.tmpl` extension)
- `logo.png` is copied verbatim (binary file)
- `.render.yaml` is not copied (control file)

## Template Processing

Each `.tmpl` file receives the full data object:

`templates/config/app.yaml.tmpl`:
```yaml
name: {{ .app.name }}
port: {{ .app.port }}
environment: {{ .env | default "development" }}
```

`templates/scripts/setup.sh.tmpl`:
```bash
#!/bin/bash
echo "Setting up {{ .app.name }}..."
createdb -h {{ .database.host }} -p {{ .database.port }} {{ .app.name }}
```

## Control Files

Use `.render.yaml` to customize output paths:

```yaml
paths:
  "src": "{{ .app.name }}"
  "model.go.tmpl": "{{ .modelName | snakeCase }}.go"
```

See [Control Files Guide](control-files.md) for details.

## Previewing Output

Use `--dry-run` to see what files would be created:

```bash
render ./templates data.json -o ./output --dry-run
```

Output:
```
Would write: output/config/app.yaml
Would write: output/config/database.yaml
Would write: output/scripts/setup.sh
Would copy: output/scripts/deploy.sh
Would copy: output/assets/logo.png
```

## Overwriting Existing Files

By default, render refuses to overwrite existing files:

```bash
render ./templates data.json -o ./output
# Error: output/config/app.yaml exists, use --force to overwrite
```

Use `--force` to overwrite:

```bash
render ./templates data.json -o ./output --force
```

## Ignored Files

These files are automatically excluded from output:
- `.render.yaml`, `.render.yml`, `render.json` (control files)

## Machine-Readable Output

Use `--json` for scripting:

```bash
render ./templates data.json -o ./output --json
```

Output (one JSON object per line):
```json
{"path":"output/config/app.yaml","action":"write"}
{"path":"output/config/database.yaml","action":"write"}
{"path":"output/scripts/setup.sh","action":"write"}
{"path":"output/scripts/deploy.sh","action":"copy"}
{"path":"output/assets/logo.png","action":"copy"}
```

## Best Practices

### Organize Templates Logically

```
templates/
  src/           # Source code templates
  config/        # Configuration files
  docs/          # Documentation templates
  scripts/       # Build/deployment scripts
```

### Use Consistent Naming

```
templates/
  user-service.go.tmpl
  user-service_test.go.tmpl
  user-repository.go.tmpl
```

### Separate Static from Dynamic

```
templates/
  static/        # Files copied verbatim
    images/
    fonts/
  dynamic/       # Template files
    index.html.tmpl
    config.js.tmpl
```

### Include Sample Data

```
project/
  templates/
    config.yaml.tmpl
  sample-data.json      # Example data for testing
  .render.yaml          # Path mappings
```

Test with:
```bash
render ./templates sample-data.json -o ./test-output --dry-run
```
