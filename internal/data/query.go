// Package data provides JSON and YAML data loading functionality.
package data

import (
	"fmt"

	"github.com/itchyny/gojq"
)

// Query executes a jq expression against data and returns the first result.
// This is useful for transforming data (e.g., extracting a nested object).
func Query(data any, expression string) (any, error) {
	query, err := gojq.Parse(expression)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query expression: %w", err)
	}

	iter := query.Run(data)
	result, ok := iter.Next()
	if !ok {
		return nil, nil // No results
	}

	if err, isErr := result.(error); isErr {
		return nil, fmt.Errorf("query execution error: %w", err)
	}

	return result, nil
}

// QueryAll executes a jq expression against data and collects all results.
// This is useful for extracting multiple items from data (e.g., ".items[]").
func QueryAll(data any, expression string) ([]any, error) {
	query, err := gojq.Parse(expression)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query expression: %w", err)
	}

	var results []any
	iter := query.Run(data)

	for {
		result, ok := iter.Next()
		if !ok {
			break
		}

		if err, isErr := result.(error); isErr {
			return nil, fmt.Errorf("query execution error: %w", err)
		}

		results = append(results, result)
	}

	return results, nil
}
