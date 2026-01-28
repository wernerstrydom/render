#!/bin/bash

# Ensure we are in the example directory
cd "$(dirname "$0")"

# Path to the render binary
RENDER="../../bin/render"

# Build render if it doesn't exist
if [ ! -f "$RENDER" ]; then
    echo "Building render..."
    (cd ../.. && go build -o bin/render ./cmd/render/main.go)
fi

echo "Generating docker-compose.yaml..."
$RENDER templates stack.yaml -o output --force

echo "Generation complete! Check the 'output' directory."
