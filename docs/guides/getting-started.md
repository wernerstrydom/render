# Getting Started

This guide walks you through your first use of render.

## Prerequisites

- Go 1.24 or later (for installation)
- Basic familiarity with command-line tools

## Installation

```bash
go install github.com/wernerstrydom/render/cmd/render@latest
```

Verify the installation:

```bash
render --help
```

## Your First Template

Create a simple template file `greeting.tmpl`:

```
Hello, {{ .name }}!
Welcome to {{ .company }}.
```

Create a data file `data.json`:

```json
{
  "name": "Alice",
  "company": "Acme Corp"
}
```

Render the template:

```bash
render greeting.tmpl data.json -o greeting.txt
```

View the output:

```bash
cat greeting.txt
```

Output:
```
Hello, Alice!
Welcome to Acme Corp.
```

## Using YAML Data

render supports both JSON and YAML. Create `data.yaml`:

```yaml
name: Bob
company: Tech Inc
```

```bash
render greeting.tmpl data.yaml -o greeting.txt
```

## Template Functions

render provides many built-in functions. Update `greeting.tmpl`:

```
Hello, {{ .name | upper }}!
Welcome to {{ .company }}.
Your username is: {{ .name | lower | replace " " "_" }}
```

Output:
```
Hello, ALICE!
Welcome to Acme Corp.
Your username is: alice
```

## Conditional Logic

Create `status.tmpl`:

```
User: {{ .name }}
Status: {{ if .active }}Active{{ else }}Inactive{{ end }}
{{ if gt .loginCount 100 }}Power user!{{ end }}
```

Create `user.json`:

```json
{
  "name": "Alice",
  "active": true,
  "loginCount": 150
}
```

```bash
render status.tmpl user.json -o status.txt
```

## Rendering Multiple Items

Create `users.json`:

```json
{
  "users": [
    {"name": "alice", "role": "admin"},
    {"name": "bob", "role": "user"},
    {"name": "charlie", "role": "user"}
  ]
}
```

Create `user-profile.tmpl`:

```
Name: {{ .name }}
Role: {{ .role | title }}
```

Generate one file per user:

```bash
render user-profile.tmpl users.json --item-query '.users[]' -o '{{.name}}.txt'
```

This creates:
- `alice.txt`
- `bob.txt`
- `charlie.txt`

## Rendering a Directory

Create a template directory:

```
templates/
  config.yaml.tmpl
  README.md.tmpl
```

`templates/config.yaml.tmpl`:
```yaml
app:
  name: {{ .appName }}
  version: {{ .version }}
```

`templates/README.md.tmpl`:
```markdown
# {{ .appName }}

Version {{ .version }}
```

Create `project.json`:

```json
{
  "appName": "MyApp",
  "version": "1.0.0"
}
```

Render the entire directory:

```bash
render ./templates project.json -o ./output
```

Output structure:
```
output/
  config.yaml
  README.md
```

## Preview Changes

Use `--dry-run` to see what would be created:

```bash
render ./templates project.json -o ./output --dry-run
```

## Overwriting Files

By default, render refuses to overwrite existing files. Use `--force`:

```bash
render config.tmpl data.json -o config.yaml --force
```

## Next Steps

- Learn about [rendering modes](../concepts/modes.md)
- Explore [template functions](../reference/functions.md)
- See [jq query expressions](query-expressions.md)
- Configure [path mappings](control-files.md)
