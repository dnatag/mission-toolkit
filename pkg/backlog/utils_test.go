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
		{"No filters", []string{"a", "b", "c"}, []string{}, []string{}, []string{"a", "b", "c"}},
		{"Include only", []string{"a", "b", "c", "d"}, []string{"a", "c"}, []string{}, []string{"a", "c"}},
		{"Exclude only", []string{"a", "b", "c", "d"}, []string{}, []string{"b", "d"}, []string{"a", "c"}},
		{"Include and exclude", []string{"a", "b", "c", "d", "e"}, []string{"a", "b", "c", "d"}, []string{"b"}, []string{"a", "c", "d"}},
		{"Exclude takes precedence", []string{"a", "b", "c"}, []string{"a", "b", "c"}, []string{"a", "b"}, []string{"c"}},
		{"Empty items", []string{}, []string{"a"}, []string{}, []string{}},
		{"Non-matching include returns empty", []string{"a", "b", "c"}, []string{"x", "y"}, []string{}, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterStringSlice(tt.items, tt.include, tt.exclude)
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

func TestContainsString(t *testing.T) {
	if !ContainsString([]string{"a", "b"}, "a") {
		t.Error("expected true")
	}
	if ContainsString([]string{"a", "b"}, "c") {
		t.Error("expected false")
	}
}

func TestRemoveString(t *testing.T) {
	result := RemoveString([]string{"a", "b", "c"}, "b")
	if len(result) != 2 || result[0] != "a" || result[1] != "c" {
		t.Errorf("expected [a c], got %v", result)
	}
}
