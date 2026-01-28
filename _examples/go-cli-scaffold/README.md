# Go CLI Application Scaffold

This example demonstrates generating a complete CLI application scaffold using
Cobra and Viper. Define your commands, flags, and structure in YAML and generate
a production-ready CLI.

## Why Use Render?

Building a CLI manually involves:
- Creating the root command with proper initialization
- Adding subcommands with consistent patterns
- Setting up configuration with environment variable support
- Wiring flags to configuration values

With `render`, define the CLI structure declaratively and generate it all at once.

## Structure Generated

```
{app}/
├── main.go
├── cmd/
│   ├── root.go           # Root command with config setup
│   ├── version.go        # Version command
│   └── {command}/
│       └── {command}.go  # Subcommand implementation
├── internal/
│   └── config/
│       └── config.go     # Configuration struct
└── .{app}.yaml           # Example config file
```

## Usage

This example demonstrates a **two-pass** approach:

```bash
# Pass 1: Generate base structure (directory mode)
render templates cli.yaml -o output

# Pass 2: Generate one file per command (each mode)
render command-template/command.go.tmpl commands.yaml -o "output/cmd/{{.name}}.go"

# Or run both with the script
./render.sh
```

## Template Features Demonstrated

- **Two-Pass Rendering**: Base structure + per-item files
- **Directory Mode**: Generates main.go, go.mod, cmd/root.go, cmd/version.go
- **Each Mode**: Generates cmd/get.go, cmd/create.go, cmd/delete.go, cmd/describe.go
- **Separate Data Files**: cli.yaml for base, commands.yaml for commands array

## CLI Features Generated

- Cobra command structure with help and completion
- Viper configuration with file, env, and flag binding
- Persistent and local flags per command
- Version command with build info
- Config file support with sensible defaults

## Real-World Use Case

An AI assistant asked to "create a CLI for managing Kubernetes resources with
get, create, delete, and describe commands" can generate a complete, consistent
CLI scaffold in seconds.
