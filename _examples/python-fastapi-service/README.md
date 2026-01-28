# Python FastAPI Service Example

This example demonstrates how to scaffold a Python FastAPI project with multiple resources using `render`.

It uses a **two-pass rendering** approach:
1. **Pass 1 (Directory Mode)**: Generates the project skeleton (FastAPI app, requirements, packages) using `project.yaml`.
2. **Pass 2 & 3 (Each Mode)**: Generates individual models and routers for each resource defined in `resources.yaml`.

## Why Use Render?

Generating a modern API service involves several repetitive files per resource (Model, Router, Tests). With `render`:
- **Consistency**: All routers follow the same pattern and injection style.
- **Speed**: Add a new resource to `resources.yaml` and `project.yaml`, run the script, and the boilerplate is ready.
- **Automation**: Easily integrable into CI/CD or CLI scaffolding tools.

## Usage

Run the render command:

```bash
./render.sh
```

Or manually:

```bash
# Generate skeleton
../../bin/render templates project.yaml -o output --force

# Generate resources
../../bin/render resource-template/model.py.tmpl resources.yaml -o "output/app/models/{{.name | snakeCase}}.py" --force
../../bin/render resource-template/router.py.tmpl resources.yaml -o "output/app/routers/{{.name | snakeCase}}.py" --force
```

## Structure

- `project.yaml`: Project metadata and list of resource names for imports.
- `resources.yaml`: Detailed resource definitions (fields, types).
- `templates/`: Project-wide skeleton.
- `resource-template/`: Templates for individual resource files.
