# GitHub Actions Workflow Generator

This example demonstrates generating GitHub Actions workflows for different
project types. Define your project structure and generate appropriate CI/CD
workflows with best practices built in.

## Why Use Render?

Setting up GitHub Actions involves:
- Understanding best practices for each language/framework
- Configuring caching, artifacts, and deployments correctly
- Maintaining consistency across multiple repositories

With `render`, generate proven workflow patterns instantly.

## Workflow Types

- **go**: Go project with testing, linting, and releases
- **node**: Node.js project with npm/yarn support
- **python**: Python project with pip, pytest, and typing
- **docker**: Multi-platform Docker builds with caching

## Usage

```bash
# Generate workflows for a Go project
render templates project.yaml -o .github/workflows

# Generate for multiple project types
render templates/go.yaml.tmpl config.yaml -o .github/workflows/go.yaml
render templates/docker.yaml.tmpl config.yaml -o .github/workflows/docker.yaml
```

## Template Features Demonstrated

- **File Mode**: Generate specific workflows
- **Conditionals**: Include/exclude steps based on config
- **Environment Variables**: Secure secrets handling
- **Matrix Builds**: Test across versions

## Workflow Features

- **Security**: Minimal permissions, dependency scanning
- **Caching**: Smart caching for faster builds
- **Releases**: Automatic versioning and changelog
- **Docker**: Multi-platform with layer caching
- **Notifications**: Slack/email on failure

## Real-World Use Case

An AI assistant asked to "set up CI/CD for our new Go service with Docker
deployment" can generate production-ready workflows with proper caching,
security scanning, and multi-platform Docker builds.
