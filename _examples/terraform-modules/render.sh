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

echo "Generating Terraform modules..."
echo ""
echo "Note: This example demonstrates templates for modules and environments."
echo "Full generation requires combining module and environment data."
echo "See README.md for details."
echo ""

# For demonstration, show what the module template would produce
echo "Module template preview (vpc):"
cat "${SCRIPT_DIR}/module-template/main.tf.tmpl" | head -20

echo ""
echo "Environment template preview:"
cat "${SCRIPT_DIR}/env-template/main.tf.tmpl" | head -20

echo ""
echo "To generate, prepare per-module YAML files and run:"
echo "  render module-template/main.tf.tmpl vpc.yaml -o output/modules/vpc/main.tf"
