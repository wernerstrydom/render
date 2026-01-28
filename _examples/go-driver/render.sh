#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RENDER="${SCRIPT_DIR}/../../render"
OUTPUT_DIR="${SCRIPT_DIR}/output"

# Build render if needed
if [ ! -f "$RENDER" ]; then
    echo "Building render..."
    (cd "${SCRIPT_DIR}/../.." && go build -o render ./cmd/render)
fi

# Clean and create output directory
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

echo "=== Pass 1: Generating core package (facade, interface, registration) ==="
"$RENDER" "${SCRIPT_DIR}/templates" "${SCRIPT_DIR}/drivers.yaml" \
    -o "${OUTPUT_DIR}/pkg/storage"

echo ""
echo "=== Pass 2: Generating driver implementations ==="
"$RENDER" "${SCRIPT_DIR}/driver-template/driver_impl.go.tmpl" "${SCRIPT_DIR}/drivers-list.yaml" \
    -o "${OUTPUT_DIR}/pkg/storage/drivers/{{.name}}/{{.name}}.go"

echo ""
echo "Output generated in: ${OUTPUT_DIR}"
echo ""
find "$OUTPUT_DIR" -type f | sort
