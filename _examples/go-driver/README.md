# Go Driver Interface Pattern

This example demonstrates how to generate a driver-based architecture in Go,
following the pattern used by `database/sql`. This pattern enables pluggable
implementations while maintaining a consistent API.

## Why Use Render?

The driver pattern requires multiple coordinated files:
- The public facade (user-facing API)
- The driver interface (implementor contract)
- A registration mechanism
- Multiple driver implementations

With `render`, generate a complete driver architecture using two passes.

## Structure Generated

```
pkg/storage/
├── storage.go            # Public facade API
├── driver.go             # Driver interface definition
├── register.go           # Driver registration mechanism
└── drivers/
    ├── s3/
    │   └── s3.go         # S3 driver implementation
    ├── gcs/
    │   └── gcs.go        # GCS driver implementation
    ├── azure/
    │   └── azure.go      # Azure driver implementation
    ├── filesystem/
    │   └── filesystem.go # Local filesystem driver
    └── memory/
        └── memory.go     # In-memory driver for testing
```

## Usage

This example demonstrates a **two-pass** approach:

```bash
# Pass 1: Generate core package (directory mode)
render templates drivers.yaml -o pkg/storage

# Pass 2: Generate driver implementations (each mode)
render driver-template/driver_impl.go.tmpl drivers-list.yaml \
    -o "pkg/storage/drivers/{{.name}}/{{.name}}.go"

# Or run both with the script
./render.sh
```

## Template Features Demonstrated

- **Two-Pass Rendering**: Core package + per-driver implementations
- **Directory Mode**: Generates facade, interface, and registration
- **Each Mode**: Generates one driver per entry in drivers-list.yaml
- **Separate Data Files**: drivers.yaml for core, drivers-list.yaml for each mode

## Pattern Benefits

1. **Decoupling**: Users depend on interface, not implementation
2. **Extensibility**: Add new drivers without changing core code
3. **Testability**: Mock drivers for unit testing
4. **Init-time Registration**: Drivers register via `init()` functions

## Real-World Use Case

An AI assistant asked to "create a storage abstraction supporting S3, GCS, and
local filesystem" can generate the complete driver architecture in seconds,
ensuring all implementations follow the same contract.
