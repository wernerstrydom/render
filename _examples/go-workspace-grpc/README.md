# Go Workspace with gRPC Services

This example demonstrates how to quickly scaffold a Go workspace containing multiple
modules, each with a complete gRPC service architecture.

## Why Use Render?

When an AI assistant needs to create a multi-module Go workspace, writing each file
individually is slow and error-prone. With `render`:

1. **Speed**: Generate dozens of files in a single command
2. **Consistency**: All services follow the same structure
3. **Maintainability**: Update the template once, regenerate all services
4. **Customization**: Each service gets its own name, port, and configuration

## Structure Generated

For each service defined in `services.yaml`, render creates:

```
services/
└── {service-name}/
    ├── go.mod
    ├── proto/
    │   └── {service}.proto
    ├── cmd/
    │   ├── server/
    │   │   └── main.go
    │   └── cli/
    │       └── main.go
    ├── internal/
    │   ├── server/
    │   │   └── server.go
    │   └── client/
    │       └── client.go
    ├── api/
    │   └── bff/
    │       └── main.go
    └── web/
        ├── admin/
        │   └── index.html
        └── frontend/
            └── index.html
```

## Usage

```bash
# Generate all services
render templates services.yaml -o "services/{{.name | kebabCase}}"

# Preview what would be generated
render templates services.yaml -o "services/{{.name | kebabCase}}" --dry-run
```

## Template Features Demonstrated

- **Each Mode**: Generates a complete directory structure per service
- **Path Transformation**: Uses `.render.yaml` for dynamic file naming
- **Casing Functions**: `kebabCase`, `pascalCase`, `snakeCase` for idiomatic naming
- **Nested Data**: Services with ports, descriptions, and custom settings

## Real-World Use Case

An AI assistant asked to "create a microservices architecture with user, order, and
payment services" can:

1. Define the services in a YAML file
2. Run a single render command
3. Have a complete, consistent workspace ready for implementation
