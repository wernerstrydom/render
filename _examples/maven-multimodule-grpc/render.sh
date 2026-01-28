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

echo "Generating Maven multi-module project..."
"$RENDER" "${SCRIPT_DIR}/templates" "${SCRIPT_DIR}/project.yaml" \
    -o "${OUTPUT_DIR}"

echo ""
echo "Output generated in: ${OUTPUT_DIR}"
echo ""
find "$OUTPUT_DIR" -type f
