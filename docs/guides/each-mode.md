# Each Mode Guide

Each mode renders a template once per item extracted from the data, generating multiple output files.

## When to Use Each Mode

- Generate one file per user, product, or entity
- Create configuration files for multiple environments
- Produce documentation pages from a data array
- Generate code files from a schema definition

## Triggering Each Mode

Each mode is triggered by either:

1. **Dynamic output path**: Output contains `{{...}}`
2. **Item query**: Using `--item-query` flag

## Basic Usage

### Dynamic Output Path

```bash
render user.tmpl users.json -o '{{.username}}.txt'
```

The output path is itself a Go template, rendered for each item.

### With Item Query

```bash
render user.tmpl data.json --item-query '.users[]' -o '{{.username}}.txt'
```

The `--item-query` extracts items from the data using a jq expression.

## Example Walkthrough

### Data

`users.json`:
```json
{
  "company": "Acme Corp",
  "users": [
    {
      "username": "alice",
      "name": "Alice Smith",
      "role": "admin",
      "email": "alice@acme.com"
    },
    {
      "username": "bob",
      "name": "Bob Jones",
      "role": "developer",
      "email": "bob@acme.com"
    },
    {
      "username": "charlie",
      "name": "Charlie Brown",
      "role": "developer",
      "email": "charlie@acme.com"
    }
  ]
}
```

### Template

`user-profile.tmpl`:
```
Username: {{ .username }}
Name: {{ .name }}
Role: {{ .role | title }}
Email: {{ .email }}
```

### Rendering

```bash
render user-profile.tmpl users.json --item-query '.users[]' -o '{{.username}}.txt'
```

### Output Files

`alice.txt`:
```
Username: alice
Name: Alice Smith
Role: Admin
Email: alice@acme.com
```

`bob.txt`:
```
Username: bob
Name: Bob Jones
Role: Developer
Email: bob@acme.com
```

`charlie.txt`:
```
Username: charlie
Name: Charlie Brown
Role: Developer
Email: charlie@acme.com
```

## Output Path Patterns

### Simple Filename

```bash
-o '{{.name}}.txt'
```

### With Directory

```bash
-o 'users/{{.username}}/profile.txt'
```

### Using Functions

```bash
-o '{{.name | snakeCase}}.yaml'
-o '{{.name | kebabCase}}.md'
-o '{{.type | lower}}/{{.name}}.go'
```

### Computed Values

```bash
-o '{{printf "%03d" .id}}-{{.name}}.txt'
```

## Item Query Expressions

The `--item-query` flag accepts jq expressions:

### Extract Array

```bash
--item-query '.users[]'
```

### Filter Items

```bash
--item-query '.users[] | select(.role == "admin")'
--item-query '.users[] | select(.active)'
--item-query '.items[] | select(.price > 100)'
```

### Transform Items

```bash
--item-query '.users[] | {name, email}'
--item-query '.users[] | . + {fullName: (.first + " " + .last)}'
```

### Flatten Nested Arrays

```bash
--item-query '.departments[].employees[]'
```

## Combining with --query

Use `--query` to transform data before `--item-query`:

```bash
render template.tmpl data.json \
  --query '.result.data' \
  --item-query '.items[]' \
  -o '{{.id}}.txt'
```

The `--query` runs first, then `--item-query` extracts items from the result.

## Directory Templates in Each Mode

Each mode also works with template directories:

```bash
render ./user-templates users.json --item-query '.users[]' -o './output/{{.username}}'
```

This creates:
```
output/
  alice/
    profile.txt
    config.yaml
  bob/
    profile.txt
    config.yaml
  charlie/
    profile.txt
    config.yaml
```

## Error Handling

### Empty Result

If `--item-query` returns no items, no files are created:

```bash
render user.tmpl users.json --item-query '.users[] | select(.role == "superadmin")' -o '{{.username}}.txt'
# No output (no users match)
```

### Invalid Output Path

If the output path template fails, the error includes the item:

```bash
render user.tmpl users.json --item-query '.users[]' -o '{{.missing}}.txt'
# Error: rendering output path for item 0: "missing" is not a field
```

## Best Practices

### Unique Output Paths

Ensure each item produces a unique output path:

```bash
# Good: username is unique
-o '{{.username}}.txt'

# Risk: role might not be unique
-o '{{.role}}.txt'  # May overwrite!

# Better: combine fields
-o '{{.role}}/{{.username}}.txt'
```

### Preview First

Always preview with `--dry-run`:

```bash
render user.tmpl users.json --item-query '.users[]' -o '{{.username}}.txt' --dry-run
```

### Validate Data

Ensure required fields exist using conditionals:

```
{{ if not .username }}username is required{{ end }}
```

### Handle Missing Fields

Use `default` for optional fields:

```
Email: {{ .email | default "not provided" }}
```
