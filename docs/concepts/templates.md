# Go Templates

render uses Go's `text/template` package for template processing. This guide covers the template syntax and features.

## Basic Syntax

### Actions

Template actions are enclosed in double curly braces: `{{ action }}`.

```
Hello, {{ .name }}!
```

### Comments

Comments are ignored in output:

```
{{/* This is a comment */}}
```

### Whitespace Control

Trim whitespace around actions with `-`:

```
{{- .value -}}    {{/* Trims whitespace on both sides */}}
{{- .value }}     {{/* Trims whitespace on the left */}}
{{ .value -}}     {{/* Trims whitespace on the right */}}
```

## Accessing Data

### Dot (.)

The dot represents the current data context:

```
{{ . }}           {{/* The entire data object */}}
{{ .name }}       {{/* Field "name" of current object */}}
{{ .user.email }} {{/* Nested field access */}}
```

### Variables

Assign values to variables with `$`:

```
{{ $name := .user.name }}
Hello, {{ $name }}!

{{ $count := len .items }}
Total: {{ $count }} items
```

### Range

Iterate over arrays and maps:

```
{{ range .users }}
  Name: {{ .name }}
{{ end }}

{{ range $index, $user := .users }}
  {{ $index }}: {{ $user.name }}
{{ end }}

{{ range $key, $value := .config }}
  {{ $key }} = {{ $value }}
{{ end }}
```

### Conditionals

```
{{ if .enabled }}
  Feature is enabled
{{ end }}

{{ if .count }}
  Count: {{ .count }}
{{ else }}
  No items
{{ end }}

{{ if eq .status "active" }}
  Active
{{ else if eq .status "pending" }}
  Pending
{{ else }}
  Unknown
{{ end }}
```

### With

Change the dot context:

```
{{ with .database }}
  Host: {{ .host }}
  Port: {{ .port }}
{{ end }}

{{ with .optional }}
  Value: {{ . }}
{{ else }}
  No value provided
{{ end }}
```

## Pipelines

Chain values and functions using pipes:

```
{{ .name | upper }}
{{ .name | lower | title }}
{{ .items | len }}
{{ .text | replace "old" "new" }}
```

## Function Calls

Call functions with arguments:

```
{{ upper .name }}
{{ replace "old" "new" .text }}
{{ add 1 2 }}
```

## Built-in Functions

Go templates provide these built-in functions:

| Function | Description |
|----------|-------------|
| `call` | Calls a function |
| `html` | HTML escapes a string |
| `index` | Indexes arrays/slices/maps |
| `slice` | Slices arrays/slices/strings |
| `js` | JavaScript escapes a string |
| `print` | Alias for fmt.Sprint |
| `printf` | Alias for fmt.Sprintf |
| `println` | Alias for fmt.Sprintln |
| `urlquery` | URL query escapes a string |

> **Note:** render overrides several Go built-in functions (`and`, `or`, `not`, `eq`, `ne`, `lt`, `le`, `gt`, `ge`, `len`) with custom versions. For example, `and` and `or` return `bool` instead of the original Go behavior of returning the argument values. See the [Template Functions Reference](../reference/functions.md) for exact behavior.

## Comparison Functions

| Function | Description |
|----------|-------------|
| `eq` | Equal (deep comparison) |
| `ne` | Not equal |
| `lt` | Less than (numeric) |
| `le` | Less than or equal (numeric) |
| `gt` | Greater than (numeric) |
| `ge` | Greater than or equal (numeric) |

```
{{ if eq .status "active" }}Active{{ end }}
{{ if gt .count 10 }}Many items{{ end }}
```

## Custom Functions

render provides 80+ additional template functions. See [Template Functions Reference](../reference/functions.md) for the complete list.

Common categories:
- **Casing**: `lower`, `upper`, `camelCase`, `snakeCase`, etc.
- **String**: `trim`, `replace`, `contains`, `split`, `join`, etc.
- **Conversion**: `toString`, `toInt`, `toJson`, `fromJson`, etc.
- **Collections**: `list`, `dict`, `keys`, `values`, `merge`, etc.
- **Math**: `add`, `sub`, `mul`, `div`, `mod`, etc.
- **Regex**: `regexMatch`, `regexReplace`, `regexFind`, etc.

## Template Examples

### Config File Generation

```yaml
# config.yaml.tmpl
server:
  host: {{ .server.host | default "localhost" }}
  port: {{ .server.port | default 8080 }}

database:
  url: postgresql://{{ .db.host }}:{{ .db.port }}/{{ .db.name }}

{{ if .features }}
features:
{{ range .features }}
  - {{ . }}
{{ end }}
{{ end }}
```

### Code Generation

```go
// model.go.tmpl
package {{ .package }}

type {{ .name | pascalCase }} struct {
{{ range .fields }}
    {{ .name | pascalCase }} {{ .type }} `json:"{{ .name | snakeCase }}"`
{{ end }}
}
```

### Conditional Logic

```
{{ if and .enabled (gt .count 0) }}
  Feature enabled with {{ .count }} items
{{ else if .enabled }}
  Feature enabled but no items
{{ else }}
  Feature disabled
{{ end }}
```

### Working with Lists

```
Users: {{ .users | len }}

{{ range $i, $user := .users }}
{{ add $i 1 }}. {{ $user.name }} ({{ $user.email }})
{{ end }}

First: {{ .users | first | get "name" }}
Last: {{ .users | last | get "name" }}
```

## Common Patterns

### Default Values

```
{{ .value | default "fallback" }}
{{ if .value }}{{ .value }}{{ else }}fallback{{ end }}
```

### Conditional Output

```
{{ if .debug }}DEBUG=true{{ end }}
```

### JSON in Templates

```
config: {{ .settings | toJson }}
pretty:
{{ .settings | toPrettyJson | indent 2 }}
```

### String Manipulation

```
{{ .name | lower | replace " " "_" }}
{{ printf "%s-%s" .prefix .name }}
```
