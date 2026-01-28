# Data Sources

render accepts data in JSON or YAML format. The data becomes available to templates as the root context (`.`).

## JSON Data

JSON files are parsed into Go data structures:

```json
{
  "name": "MyApp",
  "version": "1.0.0",
  "features": ["auth", "logging", "metrics"],
  "database": {
    "host": "localhost",
    "port": 5432
  }
}
```

Access in templates:

```
Name: {{ .name }}
Version: {{ .version }}
DB Host: {{ .database.host }}
Features:
{{ range .features }}
  - {{ . }}
{{ end }}
```

## YAML Data

YAML files are parsed identically to JSON:

```yaml
name: MyApp
version: 1.0.0
features:
  - auth
  - logging
  - metrics
database:
  host: localhost
  port: 5432
```

The same template works with both JSON and YAML:

```
Name: {{ .name }}
Version: {{ .version }}
```

## Data Types

### Strings

```json
{"greeting": "Hello, World!"}
```

```
{{ .greeting }}
{{ .greeting | upper }}
{{ .greeting | len }}
```

### Numbers

```json
{"count": 42, "price": 19.99}
```

```
Count: {{ .count }}
Price: ${{ .price }}
Doubled: {{ mul .count 2 }}
```

### Booleans

```json
{"enabled": true, "debug": false}
```

```
{{ if .enabled }}Enabled{{ else }}Disabled{{ end }}
{{ if .debug }}DEBUG MODE{{ end }}
```

### Arrays

```json
{"items": ["apple", "banana", "cherry"]}
```

```
{{ range .items }}
  - {{ . }}
{{ end }}

First: {{ .items | first }}
Last: {{ .items | last }}
Count: {{ .items | len }}
```

### Objects

```json
{
  "user": {
    "name": "Alice",
    "email": "alice@example.com"
  }
}
```

```
{{ with .user }}
  Name: {{ .name }}
  Email: {{ .email }}
{{ end }}

{{ .user.name }} <{{ .user.email }}>
```

### Nested Structures

```json
{
  "servers": [
    {
      "name": "web1",
      "config": {
        "port": 8080,
        "ssl": true
      }
    },
    {
      "name": "web2",
      "config": {
        "port": 8081,
        "ssl": false
      }
    }
  ]
}
```

```
{{ range .servers }}
Server: {{ .name }}
  Port: {{ .config.port }}
  SSL: {{ if .config.ssl }}enabled{{ else }}disabled{{ end }}
{{ end }}
```

## Data Transformation

### Using --query

Transform data before rendering with jq expressions:

```bash
# Extract a nested object
render template.tmpl data.json --query '.config.database' -o output.txt

# Filter an array
render template.tmpl data.json --query '.users | map(select(.active))' -o output.txt
```

Example data:
```json
{
  "config": {
    "database": {
      "host": "db.example.com",
      "port": 5432
    }
  }
}
```

With `--query '.config.database'`, the template receives:
```json
{
  "host": "db.example.com",
  "port": 5432
}
```

### Using --item-query

Extract items for iteration in each mode:

```bash
render user.tmpl data.json --item-query '.users[]' -o '{{.name}}.txt'
```

Each item from the query result becomes the root data for one template render.

## Special Considerations

### Null Values

JSON `null` becomes Go `nil`:

```json
{"value": null}
```

```
{{ if .value }}
  Value: {{ .value }}
{{ else }}
  No value provided
{{ end }}

{{ .value | default "fallback" }}
```

### Empty Arrays

```json
{"items": []}
```

```
{{ if .items }}
  {{ range .items }}{{ . }}{{ end }}
{{ else }}
  No items
{{ end }}
```

### Type Coercion

Some values may need conversion:

```json
{"port": "8080"}
```

```
{{ .port | toInt }}
{{ add (.port | toInt) 1 }}
```

### Unicode

Both JSON and YAML support Unicode:

```json
{"greeting": "Hello, \u4e16\u754c"}
```

```yaml
greeting: Hello, 世界
```

Templates handle Unicode correctly:
```
{{ .greeting }}          {{/* Hello, 世界 */}}
{{ .greeting | len }}    {{/* 9 (characters, not bytes) */}}
```

## Best Practices

1. **Use YAML for human-edited files**: YAML is more readable for configuration.

2. **Use JSON for machine-generated data**: JSON has strict syntax with no ambiguity.

3. **Validate data structure**: Use `--dry-run` to test templates against data.

4. **Handle missing fields**: Use `default` or conditionals for optional fields.

5. **Transform complex data**: Use `--query` to simplify data before template access.
