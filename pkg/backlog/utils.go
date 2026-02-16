// Package backlog provides utility functions for backlog management.
package backlog

import (
	"encoding/json"
	"fmt"
)

// filterStringSlice filters a slice of strings based on include/exclude parameters.
// Uses O(n) lookup with maps for better performance with larger filter lists.
//
// Parameters:
//   - items: The slice of strings to filter
//   - include: Only include strings in this list (if empty, include all)
//   - exclude: Exclude strings in this list (takes precedence over include)
//
// Returns:
//   - A filtered slice of strings
func filterStringSlice(items []string, include []string, exclude []string) []string {
	if len(include) == 0 && len(exclude) == 0 {
		return items
	}

	// Convert exclude/include slices to maps for O(1) lookup
	excludeMap := make(map[string]bool, len(exclude))
	for _, ex := range exclude {
		excludeMap[ex] = true
	}

	includeMap := make(map[string]bool, len(include))
	for _, inc := range include {
		includeMap[inc] = true
	}

	filtered := make([]string, 0, len(items))
	for _, item := range items {
		// Check exclude first
		if excludeMap[item] {
			continue
		}

		// If include is specified, only include matching items
		if len(include) > 0 && !includeMap[item] {
			continue
		}

		filtered = append(filtered, item)
	}

	return filtered
}

// parseJSONSlice parses JSON output into a slice of structs.
// This is a generic helper for parsing JSON arrays from external commands.
//
// The type parameter T must be a struct type that can be unmarshaled from JSON.
// The JSON input should be an array of T objects.
//
// Parameters:
//   - jsonInput: The JSON string to parse (should be an array)
//
// Returns:
//   - A slice of T structs
//   - An error if parsing fails
func parseJSONSlice[T any](jsonInput string) ([]T, error) {
	var items []T
	if err := json.Unmarshal([]byte(jsonInput), &items); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return items, nil
}
