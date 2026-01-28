# render

A CLI tool that uses Go text templates to generate output files from JSON or YAML data sources.

## Installation

```bash
go install github.com/wernerstrydom/render/cmd/render@latest
```

## Quick Start

```bash
# Render a single template
render config.tmpl values.json -o config.yaml

# Render a directory of templates
render ./templates data.yaml -o ./output

# Generate one file per item
render user.tmpl users.json --item-query '.users[]' -o '{{.username}}.txt'

# Preview without writing
render ./templates data.json -o ./dist --dry-run
```

## Usage

```
render <template-source> <data-source> -o <output> [OPTIONS]
```

### Arguments

- `template-source` - Path to template file or directory
- `data-source` - Path to JSON or YAML data file

### Options

| Flag | Description |
|------|-------------|
| `-o, --output` | Output path (required) |
| `-f, --force` | Overwrite existing files |
| `--query` | jq expression to transform data |
| `--item-query` | jq expression to extract items for iteration |
| `--control` | Path to control file for path mappings |
| `--dry-run` | Preview without writing files |
| `--json` | Machine-readable JSON output |

## Modes

render automatically selects the appropriate mode:

- **File mode**: Single template → single output
- **Directory mode**: Template directory → output directory
- **Each mode**: Template + item query → multiple outputs

## Documentation

See the [docs/](docs/) directory for comprehensive documentation:

- [Getting Started](docs/guides/getting-started.md)
- [Rendering Modes](docs/concepts/modes.md)
- [Template Functions](docs/reference/functions.md)
- [CLI Reference](docs/reference/cli.md)
- [Examples](docs/examples/)

## Template Functions

render provides 80+ template functions including:

- **Casing**: `lower`, `upper`, `camelCase`, `snakeCase`, `kebabCase`, `pascalCase`
- **String**: `trim`, `replace`, `split`, `join`, `contains`, `indent`
- **Conversion**: `toString`, `toInt`, `toJson`, `fromJson`
- **Collections**: `list`, `dict`, `keys`, `values`, `merge`, `sortAlpha`
- **Math**: `add`, `sub`, `mul`, `div`, `max`, `min`
- **Regex**: `regexMatch`, `regexReplace`, `regexFind`

See [Template Functions Reference](docs/reference/functions.md) for the complete list.

## Man Pages

Generate and view man pages:

```bash
render gen man ./man
man ./man/render.1
```

## License

MIT License - see LICENSE file for details.
