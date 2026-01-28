# Rendering Modes

render operates in one of five modes, determined automatically based on the template source and output path.

## File Mode

Renders a single template file to a single output file.

**When activated:**
- Template source is a file (not a directory)
- Output path is static (no Go template syntax)

**Behavior:**
- Template file is processed with the data
- Output is written to the exact path specified with `-o`

**Example:**
```bash
render config.yaml.tmpl values.json -o config.yaml
```

Given `config.yaml.tmpl`:
```yaml
database:
  host: {{ .database.host }}
  port: {{ .database.port }}
```

And `values.json`:
```json
{
  "database": {
    "host": "localhost",
    "port": 5432
  }
}
```

Produces `config.yaml`:
```yaml
database:
  host: localhost
  port: 5432
```

## Directory Mode

Renders a directory of templates to a mirrored output directory.

**When activated:**
- Template source is a directory

**Behavior:**
- Directory structure is preserved in output
- Files with `.tmpl` extension are processed as templates
- Files without `.tmpl` extension are copied verbatim
- The `.tmpl` extension is stripped from output filenames
- Control files (`.render.yaml`, etc.) are not copied

**Example:**
```bash
render ./templates data.json -o ./output
```

Given this structure:
```
templates/
  config.yaml.tmpl
  scripts/
    setup.sh.tmpl
  static/
    logo.png
```

Produces:
```
output/
  config.yaml
  scripts/
    setup.sh
  static/
    logo.png
```

## Each Mode

Renders a template once per item extracted from the data.

**When activated:**
- Output path contains Go template syntax (`{{...}}`)
- Or `--item-query` flag is specified

**Behavior:**
- Items are extracted using `--item-query` (or the root data if not specified)
- Template is rendered once per item
- Output filename is rendered as a template for each item
- Each item becomes the root data for its template invocation

**Example:**
```bash
render user.tmpl users.json --item-query '.users[]' -o '{{.username}}.txt'
```

Given `user.tmpl`:
```
Name: {{ .name }}
Email: {{ .email }}
Role: {{ .role }}
```

And `users.json`:
```json
{
  "users": [
    {"username": "alice", "name": "Alice Smith", "email": "alice@example.com", "role": "admin"},
    {"username": "bob", "name": "Bob Jones", "email": "bob@example.com", "role": "user"}
  ]
}
```

Produces:
- `alice.txt`:
  ```
  Name: Alice Smith
  Email: alice@example.com
  Role: admin
  ```
- `bob.txt`:
  ```
  Name: Bob Jones
  Email: bob@example.com
  Role: user
  ```

## Mode Selection Summary

| Template Source | Output Path | Mode |
|----------------|-------------|------|
| File | Static path | File mode |
| File | Path with trailing `/` | File-into-dir mode |
| File | Path with `{{...}}` | Each-file mode |
| Directory | Static path | Directory mode |
| Directory | Path with `{{...}}` | Each-directory mode |

## Output Path Formats

The output path (`-o`) determines both the mode and how files are named:

| Output Path | Meaning |
|------------|---------|
| `config.yaml` | Single file named `config.yaml` |
| `./output/` | Directory (trailing slash) |
| `./output` | File or directory (depends on template source) |
| `{{.name}}.txt` | Multiple files, one per item (each mode) |
| `./output/{{.id}}/config.yaml` | Multiple files in subdirectories |
