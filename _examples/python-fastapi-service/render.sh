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

echo "Pass 1: Generating project skeleton..."
# Use project.yaml which is a map
$RENDER templates project.yaml -o output --force

echo "Pass 2: Generating resource models..."
# Use resources.yaml which is an array, triggers 'each' mode because of {{ }} in -o
$RENDER resource-template/model.py.tmpl resources.yaml -o "output/app/models/{{.name | snakeCase}}.py" --force

echo "Pass 3: Generating resource routers..."
$RENDER resource-template/router.py.tmpl resources.yaml -o "output/app/routers/{{.name | snakeCase}}.py" --force

echo "Generation complete! Check the 'output' directory."
