package backlog

import (
	"testing"
)

func TestGetSectionType(t *testing.T) {
	manager := NewManager(t.TempDir())

	tests := []struct {
		section  string
		expected string
	}{
		{"## DECOMPOSED INTENTS", "decomposed"},
		{"## REFACTORING OPPORTUNITIES", "refactor"},
		{"## FUTURE ENHANCEMENTS", "future"},
		{"## FEATURES", "feature"},
		{"## BUGFIXES", "bugfix"},
		{"## UNKNOWN", ""},
		{"", ""},
	}

	for _, tt := range tests {
		result := manager.getSectionType(tt.section)
		if result != tt.expected {
			t.Errorf("getSectionType(%q) = %q, want %q", tt.section, result, tt.expected)
		}
	}
}

func TestGetSectionHeader(t *testing.T) {
	manager := NewManager(t.TempDir())

	tests := []struct {
		itemType string
		expected string
	}{
		{"decomposed", "## DECOMPOSED INTENTS"},
		{"refactor", "## REFACTORING OPPORTUNITIES"},
		{"future", "## FUTURE ENHANCEMENTS"},
		{"feature", "## FEATURES"},
		{"bugfix", "## BUGFIXES"},
		{"unknown", ""},
		{"", ""},
	}

	for _, tt := range tests {
		result := manager.getSectionHeader(tt.itemType)
		if result != tt.expected {
			t.Errorf("getSectionHeader(%q) = %q, want %q", tt.itemType, result, tt.expected)
		}
	}
}

func TestIsInSection(t *testing.T) {
	manager := NewManager(t.TempDir())

	tests := []struct {
		sectionHeader string
		itemType      string
		expected      bool
	}{
		{"## FEATURES", "feature", true},
		{"## BUGFIXES", "bugfix", true},
		{"## FEATURES", "bugfix", false},
		{"## UNKNOWN", "feature", false},
	}

	for _, tt := range tests {
		result := manager.isInSection(tt.sectionHeader, tt.itemType)
		if result != tt.expected {
			t.Errorf("isInSection(%q, %q) = %v, want %v", tt.sectionHeader, tt.itemType, result, tt.expected)
		}
	}
}

func TestFindAndModifySection(t *testing.T) {
	manager := NewManager(t.TempDir())

	lines := []string{
		"# Backlog",
		"",
		"## FEATURES",
		"- [ ] Existing feature",
		"",
		"## BUGFIXES",
		"- [ ] Existing bug",
	}

	result, err := manager.findAndModifySection(lines, "## FEATURES", func() []string {
		return []string{"- [ ] New feature"}
	})

	if err != nil {
		t.Fatalf("findAndModifySection failed: %v", err)
	}

	// Verify new item was added
	found := false
	for _, line := range result {
		if line == "- [ ] New feature" {
			found = true
			break
		}
	}

	if !found {
		t.Error("new feature was not added to section")
	}

	// Verify existing content is preserved
	if !contains(result, "- [ ] Existing feature") {
		t.Error("existing feature was not preserved")
	}

	if !contains(result, "## BUGFIXES") {
		t.Error("other sections were not preserved")
	}
}

func TestFindAndModifySection_NotFound(t *testing.T) {
	manager := NewManager(t.TempDir())

	lines := []string{
		"# Backlog",
		"## FEATURES",
	}

	_, err := manager.findAndModifySection(lines, "## NONEXISTENT", func() []string {
		return []string{"- [ ] Item"}
	})

	if err == nil {
		t.Error("expected error for nonexistent section, got nil")
	}
}

func TestMatchesItemType(t *testing.T) {
	manager := NewManager(t.TempDir())

	tests := []struct {
		item     string
		itemType string
		expected bool
	}{
		{"- [x] Task (from Epic: Large Feature)", "decomposed", true},
		{"- [x] Regular task", "decomposed", false},
		{"- [x] Refactor authentication logic", "refactor", true},
		{"- [x] Extract helper function", "refactor", true},
		{"- [x] Add new feature", "refactor", false},
		{"- [x] Future enhancement", "future", false},
	}

	for _, tt := range tests {
		result := manager.matchesItemType(tt.item, tt.itemType)
		if result != tt.expected {
			t.Errorf("matchesItemType(%q, %q) = %v, want %v", tt.item, tt.itemType, result, tt.expected)
		}
	}
}
