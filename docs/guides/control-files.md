# Control Files Guide

Control files configure how render maps template paths to output paths. This enables dynamic file and directory naming based on template data.

## Overview

A control file (`.render.yaml`, `.render.yml`, or `render.json`) placed in the template directory defines path mappings. Template source paths are mapped to output path templates.

## Basic Usage

Create `.render.yaml` in your template directory:

```yaml
paths:
  "model.go.tmpl": "{{ .modelName | snakeCase }}.go"
```

When rendering with data `{"modelName": "UserProfile"}`:
- Input: `model.go.tmpl`
- Output: `user_profile.go`

## File Naming

### Supported Files

render looks for control files in this order:
1. `.render.yaml`
2. `.render.yml`
3. `render.json`

Only the first found file is used.

### Explicit Path

Use `--control` to specify a control file:

```bash
render ./templates data.json -o ./output --control ./custom-config.yaml
```

This disables auto-discovery.

## Path Mapping Syntax

### File Mapping

Rename individual files:

```yaml
paths:
  "template.go.tmpl": "{{ .name }}.go"
  "config.yaml.tmpl": "{{ .appName | kebabCase }}-config.yaml"
```

### Directory Mapping

Rename directories:

```yaml
paths:
  "src": "{{ .packageName }}"
  "pkg/models": "{{ .moduleName }}/domain"
```

All files within the mapped directory use the new path.

### Combined Mappings

```yaml
paths:
  # Rename the directory
  "app": "{{ .appName }}"
  # Rename specific files
  "app/main.go.tmpl": "{{ .appName }}/{{ .appName | snakeCase }}.go"
```

## Template Syntax in Paths

Path templates support all render template functions:

### Casing Functions

```yaml
paths:
  "model.go.tmpl": "{{ .name | snakeCase }}.go"        # user_profile.go
  "model.go.tmpl": "{{ .name | kebabCase }}.go"        # user-profile.go
  "model.go.tmpl": "{{ .name | camelCase }}.go"        # userProfile.go
  "model.go.tmpl": "{{ .name | pascalCase }}.go"       # UserProfile.go
```

### String Manipulation

```yaml
paths:
  "service.go.tmpl": "{{ .name | trimSuffix \"Service\" | snakeCase }}_svc.go"
  "handler.go.tmpl": "{{ .endpoint | replace \"/\" \"_\" | trim \"_\" }}_handler.go"
```

### Conditionals

```yaml
paths:
  "entity.go.tmpl": "{{ if .isModel }}models{{ else }}entities{{ end }}/{{ .name | snakeCase }}.go"
```

## Example: Code Generator

Template structure:
```
templates/
  .render.yaml
  cmd/
    main.go.tmpl
  internal/
    model.go.tmpl
    repository.go.tmpl
  go.mod.tmpl
```

`.render.yaml`:
```yaml
paths:
  "cmd": "cmd/{{ .appName }}"
  "internal": "internal/{{ .appName }}"
  "internal/model.go.tmpl": "internal/{{ .appName }}/{{ .modelName | snakeCase }}.go"
  "internal/repository.go.tmpl": "internal/{{ .appName }}/{{ .modelName | snakeCase }}_repository.go"
```

Data:
```json
{
  "appName": "userservice",
  "modelName": "UserProfile"
}
```

Output structure:
```
output/
  cmd/
    userservice/
      main.go
  internal/
    userservice/
      user_profile.go
      user_profile_repository.go
  go.mod
```

## Example: Multi-Environment Config

Template structure:
```
templates/
  .render.yaml
  config.yaml.tmpl
```

`.render.yaml`:
```yaml
paths:
  "config.yaml.tmpl": "{{ .environment }}/config.yaml"
```

Render for each environment:

```bash
render ./templates dev.json -o ./output
render ./templates staging.json -o ./output
render ./templates prod.json -o ./output
```

Output:
```
output/
  development/
    config.yaml
  staging/
    config.yaml
  production/
    config.yaml
```

## Path Resolution

### Matching Order

Paths are matched from most specific to least specific:

```yaml
paths:
  "src/models/user.go.tmpl": "domain/user_entity.go"    # Most specific
  "src/models": "domain/models"                          # Directory
  "src": "lib"                                           # Least specific
```

### No Match

Files without a path mapping use their original name (minus `.tmpl` extension).

### The Control File Itself

The control file is never copied to output.

## JSON Format

For `render.json`:

```json
{
  "paths": {
    "model.go.tmpl": "{{ .name | snakeCase }}.go",
    "src": "{{ .package }}"
  }
}
```

## Error Handling

### Invalid Template

If a path template fails:

```bash
render ./templates data.json -o ./output
# Error: evaluating path mapping "model.go.tmpl": "unknownField" is not a field
```

### Missing Data

Use `default` for optional fields:

```yaml
paths:
  "config.yaml.tmpl": "{{ .env | default \"default\" }}/config.yaml"
```

## Best Practices

### Keep Mappings Simple

```yaml
# Good: clear, single transformation
paths:
  "model.go.tmpl": "{{ .name | snakeCase }}.go"

# Avoid: complex logic in paths
paths:
  "model.go.tmpl": "{{ if .legacy }}old{{ else }}new{{ end }}/{{ if .exported }}public{{ else }}internal{{ end }}/{{ .name | snakeCase }}.go"
```

### Document Expected Data

```yaml
# Expected data fields:
#   - appName: string (e.g., "myapp")
#   - modelName: string (e.g., "UserProfile")
paths:
  "internal": "internal/{{ .appName }}"
  "model.go.tmpl": "{{ .modelName | snakeCase }}.go"
```

### Test with Dry Run

```bash
render ./templates data.json -o ./output --dry-run
```

### Version Control the Control File

Include `.render.yaml` in your template repository so path mappings are always consistent.
