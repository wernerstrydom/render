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

echo "=== Pass 1: Generating base CLI structure ==="
"$RENDER" "${SCRIPT_DIR}/templates" "${SCRIPT_DIR}/cli.yaml" \
    -o "${OUTPUT_DIR}"

echo ""
echo "=== Pass 2: Generating command files ==="
# Extract commands array and render each command
# The commands array needs to be at the root for each mode
"$RENDER" "${SCRIPT_DIR}/command-template/command.go.tmpl" "${SCRIPT_DIR}/commands.yaml" \
    -o "${OUTPUT_DIR}/cmd/{{.name}}.go"

echo ""
echo "Output generated in: ${OUTPUT_DIR}"
echo ""
find "$OUTPUT_DIR" -type f | sort
