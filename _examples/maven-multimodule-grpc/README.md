# Maven Multi-Module Project with gRPC Services

This example demonstrates scaffolding a Maven multi-module project where each module
contains a complete gRPC service with CLI, admin UI, BFF, and browser frontend.

## Why Use Render?

Creating a Maven multi-module project manually involves:
- Writing parent POM with dependency management
- Creating module POMs with proper parent references
- Setting up consistent directory structures
- Configuring protobuf compilation plugins

With `render`, generate the entire structure from a single data file.

## Structure Generated

```
{project-name}/
├── pom.xml                           # Parent POM
└── modules/
    └── {service-name}/
        ├── pom.xml                   # Module POM
        ├── src/main/proto/
        │   └── {service}.proto
        ├── src/main/java/.../
        │   ├── server/
        │   │   └── {Service}Server.java
        │   ├── client/
        │   │   └── {Service}Client.java
        │   └── bff/
        │       └── {Service}BffApplication.java
        └── src/main/resources/
            └── static/
                ├── admin/
                │   └── index.html
                └── frontend/
                    └── index.html
```

## Usage

```bash
# Generate the complete project
render templates project.yaml -o output

# With custom package paths
render templates project.yaml -o output --force
```

## Template Features Demonstrated

- **Directory Mode**: Single invocation generates complete structure
- **Path Transformation**: `.render.yaml` maps Java package to directory structure
- **Nested Loops**: Services containing multiple entities
- **String Functions**: Package name to path conversion with `replace`

## Real-World Use Case

An AI assistant asked to "create a Spring Boot microservices project with three
services" can generate a production-ready Maven structure in seconds, including
proper dependency versions and plugin configurations.
