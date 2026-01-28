# Go Workspace Makefile Generator

This example demonstrates generating a comprehensive Makefile for a Go workspace
with multiple modules. The Makefile includes targets for building, testing, linting,
and managing all modules consistently.

## Why Use Render?

Manually maintaining a Makefile for a multi-module Go workspace is tedious:
- Each module needs similar targets
- Dependencies between modules must be tracked
- Adding a new module requires updating multiple places

With `render`, regenerate the Makefile whenever modules change.

## Generated Targets

For a workspace with modules `api`, `cli`, and `worker`:

```makefile
# Global targets
make all          # Build all modules
make test         # Test all modules
make lint         # Lint all modules
make clean        # Clean all modules
make tidy         # Tidy all modules

# Per-module targets
make build-api    # Build api module
make test-api     # Test api module
make lint-api     # Lint api module

make build-cli    # Build cli module
make build-worker # Build worker module
# ... etc
```

## Usage

```bash
# Generate Makefile
render templates/Makefile.tmpl workspace.yaml -o Makefile

# Preview
render templates/Makefile.tmpl workspace.yaml -o Makefile --dry-run
```

## Template Features Demonstrated

- **File Mode**: Single template to single output file
- **Loops**: Generate targets per module
- **Conditionals**: Different build flags per module type
- **String Functions**: Path manipulation, naming conventions

## Real-World Use Case

An AI assistant asked to "add a new service module to the workspace" can:

1. Add the module to `workspace.yaml`
2. Regenerate the Makefile
3. Have all standard targets immediately available
