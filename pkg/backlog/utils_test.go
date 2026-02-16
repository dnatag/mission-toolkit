// Package backlog provides tests for utility functions.
package backlog

import (
	"testing"
)

func TestFilterStringSlice(t *testing.T) {
	tests := []struct {
		name     string
		items    []string
		include  []string
		exclude  []string
		expected []string
	}{
		{
			name:     "No filters",
			items:    []string{"a", "b", "c"},
			include:  []string{},
			exclude:  []string{},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "Include only",
			items:    []string{"a", "b", "c", "d"},
			include:  []string{"a", "c"},
			exclude:  []string{},
			expected: []string{"a", "c"},
		},
		{
			name:     "Exclude only",
			items:    []string{"a", "b", "c", "d"},
			include:  []string{},
			exclude:  []string{"b", "d"},
			expected: []string{"a", "c"},
		},
		{
			name:     "Include and exclude",
			items:    []string{"a", "b", "c", "d", "e"},
			include:  []string{"a", "b", "c", "d"},
			exclude:  []string{"b"},
			expected: []string{"a", "c", "d"},
		},
		{
			name:     "Exclude takes precedence",
			items:    []string{"a", "b", "c"},
			include:  []string{"a", "b", "c"},
			exclude:  []string{"a", "b"},
			expected: []string{"c"},
		},
		{
			name:     "Empty items",
			items:    []string{},
			include:  []string{"a"},
			exclude:  []string{},
			expected: []string{},
		},
		{
			name:     "Non-matching include returns empty",
			items:    []string{"a", "b", "c"},
			include:  []string{"x", "y"},
			exclude:  []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterStringSlice(tt.items, tt.include, tt.exclude)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d items, got %d", len(tt.expected), len(result))
				return
			}
			for i, item := range result {
				if item != tt.expected[i] {
					t.Errorf("Item %d: expected %s, got %s", i, tt.expected[i], item)
				}
			}
		})
	}
}

func TestParseJSONSlice(t *testing.T) {
	type testItem struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}

	tests := []struct {
		name      string
		jsonInput string
		wantErr   bool
		expected  []testItem
	}{
		{
			name:      "Valid JSON array",
			jsonInput: `[{"id":"1","title":"First"},{"id":"2","title":"Second"}]`,
			wantErr:   false,
			expected:  []testItem{{ID: "1", Title: "First"}, {ID: "2", Title: "Second"}},
		},
		{
			name:      "Empty JSON array",
			jsonInput: `[]`,
			wantErr:   false,
			expected:  []testItem{},
		},
		{
			name:      "Invalid JSON",
			jsonInput: `{invalid json}`,
			wantErr:   true,
			expected:  nil,
		},
		{
			name:      "JSON object instead of array",
			jsonInput: `{"id":"1","title":"First"}`,
			wantErr:   true,
			expected:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseJSONSlice[testItem](tt.jsonInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseJSONSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(result) != len(tt.expected) {
					t.Errorf("Expected %d items, got %d", len(tt.expected), len(result))
					return
				}
				for i, item := range result {
					if item != tt.expected[i] {
						t.Errorf("Item %d: expected %+v, got %+v", i, tt.expected[i], item)
					}
				}
			}
		})
	}
}
