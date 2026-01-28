# Docker Compose Example

This example demonstrates how to generate a `docker-compose.yaml` file from a simplified stack definition.

## Why Use Render?

Hand-writing `docker-compose.yaml` for large stacks is tedious and error-prone, especially when managing environment variables and volumes across multiple services. `render` allows you to:
- **Simplify Configuration**: Define your stack in a clean YAML format.
- **Enforce Best Practices**: Your template can include standard logging, restart policies, or networks that apply to all services.

## Usage

```bash
./render.sh
```

## Structure

- `stack.yaml`: The high-level definition of your services, networks, and volumes.
- `templates/docker-compose.yaml.tmpl`: The template that transforms the stack definition into a valid Docker Compose file.
