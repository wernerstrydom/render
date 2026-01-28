# render

A CLI tool that uses Go text templates to generate output files from JSON or YAML data sources. Useful for AI agents and humans generating repetitive code, configuration files, and documentation.

## Installation

```bash
go install github.com/wernerstrydom/render/cmd/render@latest
```

Or build from source:

```bash
git clone https://github.com/wernerstrydom/render.git
cd render
go build -o bin/render ./cmd/render
```

## Usage

render supports three modes of operation:

### File Mode

Render a single template file:

```bash
render file -t template.txt -d data.json -o output.txt
render file --template config.yaml.tmpl --data values.json --output config.yaml
```

### Directory Mode

Render a directory of templates. Files with `.tmpl` extension are rendered (extension stripped), other files are copied verbatim:

```bash
render dir -t ./templates -d data.json -o ./output
render dir --template ./src --data config.yaml --output ./dist
```

### Each Mode

Render templates for each element in an array selected via jq query. The output path is itself a Go template:

```bash
render each -t user.tmpl -d users.json -q '.users[]' -o '{{.username}}.txt'
render each -t ./user-templates -d data.yaml -q '.items[]' -o './output/{{.id}}'
```

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--template` | `-t` | Path to template file or directory |
| `--data` | `-d` | Path to JSON or YAML data file |
| `--output` | `-o` | Output file or directory path |
| `--query` | `-q` | jq query to select elements (each mode only) |
| `--force` | `-f` | Overwrite existing output files |
| `--dry-run` | | Show what would be written without writing (dir mode only) |

## Configuration File

When using `dir` or `each` mode with a template directory, you can create a `.render.yaml` file (or `.render.yml`, `render.json`) to control how output paths are transformed. This is useful for renaming files or directories based on template data.

### Path Mappings

The `paths` key maps source paths to output path templates:

```yaml
# .render.yaml
paths:
  # Rename a specific file using template data
  "model.go.tmpl": "{{ .name | snakeCase }}.go"

  # Rename a directory prefix
  "src": "{{ .package }}"
```

Path templates have access to:
- All data from your JSON/YAML data file
- All template functions (snakeCase, pascalCase, etc.)

### Examples

**Rename a file based on data:**

```yaml
# Template: user.go.tmpl
# Data: {"name": "UserProfile"}
# .render.yaml
paths:
  "user.go.tmpl": "{{ .name | snakeCase }}.go"
# Output: user_profile.go
```

**Rename a directory:**

```yaml
# Template directory contains: src/main.go.tmpl
# Data: {"package": "myapp"}
# .render.yaml
paths:
  "src": "{{ .package }}"
# Output: myapp/main.go
```

**Multiple mappings:**

```yaml
paths:
  "model.go.tmpl": "{{ .name | snakeCase }}.go"
  "templates": "{{ .outputDir }}"
```

The config file itself is never copied to the output directory.

## Template Functions

In addition to Go's standard template functions, render provides:

### Casing Functions

| Function | Description | Example |
|----------|-------------|---------|
| `lower` | Lowercase | `{{ lower "Hello" }}` → `hello` |
| `upper` | Uppercase | `{{ upper "Hello" }}` → `HELLO` |
| `title` | Title Case | `{{ title "hello world" }}` → `Hello World` |
| `camelCase` | camelCase | `{{ camelCase "hello world" }}` → `helloWorld` |
| `pascalCase` | PascalCase | `{{ pascalCase "hello world" }}` → `HelloWorld` |
| `snakeCase` | snake_case | `{{ snakeCase "HelloWorld" }}` → `hello_world` |
| `kebabCase` | kebab-case | `{{ kebabCase "HelloWorld" }}` → `hello-world` |
| `upperSnakeCase` | UPPER_SNAKE | `{{ upperSnakeCase "hello" }}` → `HELLO` |
| `upperKebabCase` | UPPER-KEBAB | `{{ upperKebabCase "hello" }}` → `HELLO` |

### String Functions

| Function | Description |
|----------|-------------|
| `trim` | Remove leading/trailing whitespace |
| `trimPrefix` | Remove prefix |
| `trimSuffix` | Remove suffix |
| `replace` | Replace all occurrences |
| `contains` | Check if string contains substring |
| `hasPrefix` | Check prefix |
| `hasSuffix` | Check suffix |
| `repeat` | Repeat string N times |
| `reverse` | Reverse string |
| `substr` | Substring (start, end) |
| `truncate` | Truncate to length |
| `padLeft` | Pad left to length |
| `padRight` | Pad right to length |
| `indent` | Indent with spaces |
| `nindent` | Newline + indent |

### Splitting and Joining

| Function | Description |
|----------|-------------|
| `split` | Split string by separator |
| `join` | Join array with separator |
| `lines` | Split into lines |
| `first` | First element |
| `last` | Last element |
| `rest` | All but first |
| `initial` | All but last |
| `nth` | Nth element |

### Conversion Functions

| Function | Description |
|----------|-------------|
| `toString` | Convert to string |
| `toInt` | Convert to int |
| `toFloat` | Convert to float |
| `toBool` | Convert to bool |
| `toJson` | Convert to JSON string |
| `toPrettyJson` | Convert to pretty JSON |
| `fromJson` | Parse JSON string |

### Unicode Functions

| Function | Description |
|----------|-------------|
| `nfc` | NFC normalization |
| `nfd` | NFD normalization |
| `nfkc` | NFKC normalization |
| `nfkd` | NFKD normalization |
| `ascii` | Convert to ASCII |
| `slug` | URL-safe slug |

### Logic and Comparison

| Function | Description |
|----------|-------------|
| `eq`, `ne` | Equal, not equal |
| `lt`, `le` | Less than, less or equal |
| `gt`, `ge` | Greater than, greater or equal |
| `and`, `or`, `not` | Boolean operations |
| `default` | Default value if empty |
| `empty` | Check if empty |
| `coalesce` | First non-empty value |
| `ternary` | Conditional value |

### Collection Functions

| Function | Description |
|----------|-------------|
| `list` | Create a list |
| `dict` | Create a dictionary |
| `keys` | Get map keys |
| `values` | Get map values |
| `hasKey` | Check if key exists |
| `get` | Get value by key |
| `merge` | Merge maps |
| `append` | Append to list |
| `uniq` | Remove duplicates |
| `sortAlpha` | Sort alphabetically |
| `len` | Length of collection |

### Math Functions

| Function | Description |
|----------|-------------|
| `add`, `sub` | Addition, subtraction |
| `mul`, `div` | Multiplication, division |
| `mod` | Modulo |
| `max`, `min` | Maximum, minimum |
| `floor`, `ceil`, `round` | Rounding |

### Regex Functions

| Function | Description |
|----------|-------------|
| `regexMatch` | Check if matches |
| `regexFind` | Find first match |
| `regexFindAll` | Find all matches |
| `regexReplace` | Replace matches |
| `regexSplit` | Split by pattern |

## Examples

### Generate Configuration Files

```bash
# values.yaml
database:
  host: localhost
  port: 5432
  name: myapp

# config.yaml.tmpl
database:
  connection_string: "postgresql://{{ .database.host }}:{{ .database.port }}/{{ .database.name }}"

# Generate
render file -t config.yaml.tmpl -d values.yaml -o config.yaml
```

### Generate Multiple Files from Array

```bash
# users.json
{
  "users": [
    {"name": "Alice", "role": "admin"},
    {"name": "Bob", "role": "user"}
  ]
}

# user.tmpl
Name: {{ .name }}
Role: {{ .role }}
Username: {{ snakeCase .name }}

# Generate one file per user
render each -t user.tmpl -d users.json -q '.users[]' -o '{{ snakeCase .name }}.txt'
```

## License

MIT License - see LICENSE file for details.
