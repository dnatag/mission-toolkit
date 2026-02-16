// Package backlog provides tests for shared type constants.
package backlog

import (
	"testing"
)

func TestItemTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"Feature constant", ItemTypeFeature, "feature"},
		{"Bugfix constant", ItemTypeBugfix, "bugfix"},
		{"Decomposed constant", ItemTypeDecomposed, "decomposed"},
		{"Refactor constant", ItemTypeRefactor, "refactor"},
		{"Future constant", ItemTypeFuture, "future"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestListItemFormatConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"Open format", ListItemOpenFormat, "- [ ] "},
		{"Closed format", ListItemClosedFormat, "- [x] "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}
