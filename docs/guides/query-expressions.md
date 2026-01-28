# Query Expressions Guide

render supports jq expressions for data transformation via `--query` and `--item-query` flags.

## Overview

- `--query`: Transform the entire data before rendering
- `--item-query`: Extract items for iteration (enables each mode)

Both use [jq](https://jqlang.github.io/jq/manual/) syntax.

## Basic jq Syntax

### Identity

```bash
--query '.'          # Return data unchanged
```

### Field Access

```bash
--query '.name'           # Get "name" field
--query '.user.email'     # Nested field
--query '.users[0]'       # First array element
--query '.users[-1]'      # Last array element
```

### Array Iteration

```bash
--query '.users[]'        # Iterate over array (returns multiple values)
```

### Object Construction

```bash
--query '{name, email}'                    # Select specific fields
--query '{n: .name, e: .email}'            # Rename fields
--query '. + {computed: (.a + .b)}'        # Add computed field
```

## Using --query

Transform data before template rendering.

### Extract Nested Data

Data:
```json
{
  "result": {
    "config": {
      "database": {"host": "localhost", "port": 5432}
    }
  }
}
```

```bash
render db.tmpl data.json --query '.result.config.database' -o db.yaml
```

Template receives:
```json
{"host": "localhost", "port": 5432}
```

### Filter Arrays

Data:
```json
{
  "users": [
    {"name": "Alice", "active": true},
    {"name": "Bob", "active": false},
    {"name": "Charlie", "active": true}
  ]
}
```

```bash
render users.tmpl data.json --query '{users: [.users[] | select(.active)]}' -o active.txt
```

Template receives:
```json
{
  "users": [
    {"name": "Alice", "active": true},
    {"name": "Charlie", "active": true}
  ]
}
```

### Reshape Data

Data:
```json
{
  "firstName": "John",
  "lastName": "Doe",
  "emailAddress": "john@example.com"
}
```

```bash
render user.tmpl data.json --query '{name: (.firstName + " " + .lastName), email: .emailAddress}' -o user.txt
```

Template receives:
```json
{"name": "John Doe", "email": "john@example.com"}
```

## Using --item-query

Extract items for iteration in each mode.

### Basic Iteration

```bash
render user.tmpl data.json --item-query '.users[]' -o '{{.name}}.txt'
```

Each user becomes the root data for one template render.

### With Filtering

```bash
render user.tmpl data.json --item-query '.users[] | select(.role == "admin")' -o '{{.name}}.txt'
```

Only admin users are rendered.

### Transforming Items

```bash
render user.tmpl data.json --item-query '.users[] | {name, uppername: (.name | ascii_upcase)}' -o '{{.name}}.txt'
```

Each item is transformed before rendering.

## Common jq Operations

### Filtering

```bash
# By field value
--item-query '.items[] | select(.status == "active")'

# By existence
--item-query '.items[] | select(.email)'

# By comparison
--item-query '.items[] | select(.price > 100)'

# Multiple conditions
--item-query '.items[] | select(.active and .verified)'

# Negation
--item-query '.items[] | select(.status != "deleted")'
```

### Mapping

```bash
# Transform each item
--query '[.items[] | {name: .title, value: .amount}]'

# Add fields
--query '[.items[] | . + {processed: true}]'
```

### Sorting

```bash
--query '.items | sort_by(.name)'
--query '.items | sort_by(.date) | reverse'
```

### Grouping

```bash
--query '.items | group_by(.category)'
--query '[.items | group_by(.type) | .[] | {type: .[0].type, items: .}]'
```

### Aggregation

```bash
--query '{total: [.items[].price] | add}'
--query '{count: .items | length}'
--query '{avg: ([.items[].score] | add / length)}'
```

## Combining --query and --item-query

Apply transformations in sequence:

```bash
render item.tmpl data.json \
  --query '.response.data' \
  --item-query '.items[] | select(.active)' \
  -o '{{.id}}.txt'
```

Order of operations:
1. `--query` transforms root data
2. `--item-query` extracts items from result
3. Template renders once per item

## Debugging Queries

Test queries with jq directly:

```bash
# Test --query
jq '.result.config' data.json

# Test --item-query
jq '.users[] | select(.active)' data.json

# Test combined
jq '.result | .users[] | select(.active)' data.json
```

## Error Handling

### Invalid Query

```bash
render tmpl data.json --query '.invalid[' -o out.txt
# Error: invalid jq expression: unexpected end of input
```

### Missing Field

```bash
render tmpl data.json --query '.nonexistent' -o out.txt
# Returns null (not an error)
```

To fail on missing fields:

```bash
--query '.required // error("required field missing")'
```

### Empty Result

If `--item-query` returns no items, no files are created (not an error).

## Best Practices

### Start Simple

```bash
# Test field access first
--query '.users'

# Then add iteration
--item-query '.users[]'

# Then add filtering
--item-query '.users[] | select(.active)'
```

### Use Dry Run

```bash
render tmpl data.json --item-query '.users[]' -o '{{.name}}.txt' --dry-run
```

### Keep Queries Readable

For complex transformations, consider:
1. Pre-processing data with a separate jq command
2. Using multiple render invocations
3. Adding computed fields in templates instead

### Document Expectations

Comment your data requirements:

```bash
# Expects: {users: [{name, email, role}]}
render user.tmpl data.json --item-query '.users[]' -o '{{.name}}.txt'
```
