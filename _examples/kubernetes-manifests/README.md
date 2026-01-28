# Kubernetes Manifests Generator

This example demonstrates generating Kubernetes manifests for multiple services
across multiple environments. Define your services once and generate consistent
deployments, services, and configmaps for each environment.

## Why Use Render?

Managing Kubernetes manifests for multiple environments involves:
- Copy-pasting YAML with subtle differences
- Risk of configuration drift between environments
- Difficulty keeping all environments in sync

With `render`, maintain a single template and generate environment-specific manifests.

## Structure Generated

```
k8s/
├── base/
│   └── namespace.yaml
└── overlays/
    ├── dev/
    │   └── {service}/
    │       ├── deployment.yaml
    │       ├── service.yaml
    │       └── configmap.yaml
    ├── staging/
    │   └── {service}/
    │       └── ...
    └── prod/
        └── {service}/
            └── ...
```

## Usage

```bash
# Generate manifests for all environments
render templates cluster.yaml -o k8s

# Generate for a specific environment (filter data first)
yq '.environments[] | select(.name == "dev")' cluster.yaml | \
  render templates/env - -o k8s/overlays/dev
```

## Template Features Demonstrated

- **Directory Mode**: Complete Kubernetes structure
- **Nested Loops**: Services × Environments
- **Path Transformation**: Environment and service in paths
- **Conditionals**: Different resource limits per environment

## Environment-Specific Configuration

- **Dev**: Lower resources, single replica, debug logging
- **Staging**: Moderate resources, 2 replicas, standard config
- **Prod**: High resources, 3+ replicas, production settings

## Real-World Use Case

An AI assistant asked to "deploy our services to a new staging environment"
can generate all required manifests with environment-appropriate settings in
a single command.
