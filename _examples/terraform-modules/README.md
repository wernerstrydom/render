# Terraform Module Generator

This example demonstrates generating Terraform modules for infrastructure
components across multiple environments. Define your infrastructure once and
generate consistent, environment-specific configurations.

## Why Use Render?

Managing Terraform across environments involves:
- Duplicating modules with environment-specific values
- Risk of drift between environment configurations
- Difficulty maintaining DRY principles

With `render`, define infrastructure patterns once and generate for all environments.

## Structure Generated

```
terraform/
├── modules/
│   ├── vpc/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   ├── eks/
│   │   └── ...
│   └── rds/
│       └── ...
└── environments/
    ├── dev/
    │   ├── main.tf
    │   ├── variables.tf
    │   └── terraform.tfvars
    ├── staging/
    │   └── ...
    └── prod/
        └── ...
```

## Usage

```bash
# Generate all Terraform files
render templates infrastructure.yaml -o terraform

# Preview
render templates infrastructure.yaml -o terraform --dry-run
```

## Template Features Demonstrated

- **Directory Mode**: Complete Terraform structure
- **Nested Loops**: Modules × Environments
- **Environment Variables**: Correct sizing per environment
- **Module References**: Proper dependency handling

## Security Features

- Encryption at rest for all data stores
- Private subnets for workloads
- Security groups with minimal access
- KMS keys for sensitive data

## Real-World Use Case

An AI assistant asked to "create infrastructure for a new microservice with
database and caching" can generate complete, secure Terraform configurations
following organizational standards.
