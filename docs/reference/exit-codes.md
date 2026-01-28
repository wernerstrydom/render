# Exit Codes

render uses specific exit codes to indicate different error conditions. This enables scripting and automation.

## Exit Code Reference

| Code | Constant | Description |
|------|----------|-------------|
| 0 | `ExitSuccess` | Command completed successfully |
| 1 | `ExitRuntimeError` | Runtime error during template rendering |
| 2 | `ExitUsageError` | Invalid command-line arguments |
| 3 | `ExitInputValidation` | Input validation failed |
| 4 | `ExitPermissionDenied` | Filesystem permission error |
| 5 | `ExitOutputConflict` | Output file exists, --force not specified |
| 6 | `ExitSafetyViolation` | Security issue detected |

## Detailed Descriptions

### 0 - Success

All operations completed successfully. All files were written.

### 1 - Runtime Error

An error occurred during template execution. Common causes:

- Template syntax error
- Missing field in data
- Function error (e.g., division by zero)
- Template rendering failure

Example:
```bash
render bad.tmpl data.json -o out.txt
# Error: template: bad.tmpl:5: unexpected "}" in operand
# Exit code: 1
```

### 2 - Usage Error

Invalid command-line arguments or flags. Common causes:

- Missing required arguments
- Invalid flag values
- Unknown flags

Example:
```bash
render template.tmpl
# Error: required flag(s) "output" not set
# Exit code: 2

render template.tmpl data.json -o out.txt --unknown
# Error: unknown flag: --unknown
# Exit code: 2
```

### 3 - Input Validation Error

Input files failed validation. Common causes:

- Template file not found
- Data file not found
- Malformed JSON or YAML
- Invalid jq query expression

Example:
```bash
render missing.tmpl data.json -o out.txt
# Error: template file not found: missing.tmpl
# Exit code: 3

render template.tmpl bad.json -o out.txt
# Error: failed to parse data file: invalid JSON
# Exit code: 3
```

### 4 - Permission Denied

Filesystem permission error. Common causes:

- Cannot read template file
- Cannot read data file
- Cannot write output file
- Cannot create output directory

Example:
```bash
render template.tmpl data.json -o /root/out.txt
# Error: permission denied: cannot write to /root/out.txt
# Exit code: 4
```

### 5 - Output Conflict

Output file already exists and `--force` was not specified.

Example:
```bash
render template.tmpl data.json -o existing.txt
# Error: output file exists: existing.txt (use --force to overwrite)
# Exit code: 5
```

Solution:
```bash
render template.tmpl data.json -o existing.txt --force
```

### 6 - Safety Violation

A security issue was detected. Common causes:

- Output path traversal (e.g., `../../../etc/passwd`)
- Symlink to outside output directory
- Attempt to write outside output directory

Example:
```bash
render template.tmpl data.json -o '{{.path}}.txt'
# (where .path = "../../../etc/passwd")
# Error: path traversal detected: ../../../etc/passwd.txt
# Exit code: 6
```

## Scripting with Exit Codes

### Bash

```bash
#!/bin/bash

render template.tmpl data.json -o output.txt
case $? in
    0)
        echo "Success"
        ;;
    1)
        echo "Template error"
        exit 1
        ;;
    3)
        echo "Invalid input"
        exit 1
        ;;
    5)
        echo "File exists, retrying with --force"
        render template.tmpl data.json -o output.txt --force
        ;;
    *)
        echo "Unknown error"
        exit 1
        ;;
esac
```

### Check for Success

```bash
if render template.tmpl data.json -o output.txt; then
    echo "Generated successfully"
else
    echo "Generation failed"
    exit 1
fi
```

### Ignore Specific Errors

```bash
render template.tmpl data.json -o output.txt || {
    if [ $? -eq 5 ]; then
        echo "File exists, skipping"
    else
        exit 1
    fi
}
```

## Machine-Readable Output

Combine exit codes with `--json` for detailed error handling:

```bash
output=$(render template.tmpl data.json -o output.txt --json 2>&1)
status=$?

if [ $status -eq 0 ]; then
    echo "$output" | jq '.'
else
    echo "Error (exit code $status): $output"
fi
```

## CI/CD Integration

Exit codes integrate with CI/CD systems:

```yaml
# GitHub Actions
- name: Generate configs
  run: render ./templates config.json -o ./dist
  # Non-zero exit code fails the step
```

```yaml
# GitLab CI
generate:
  script:
    - render ./templates config.json -o ./dist
  # Non-zero exit code fails the job
```

## Debugging

When troubleshooting, check the exit code:

```bash
render template.tmpl data.json -o out.txt
echo "Exit code: $?"
```

Use `--dry-run` to preview without risk:

```bash
render template.tmpl data.json -o out.txt --dry-run
```
