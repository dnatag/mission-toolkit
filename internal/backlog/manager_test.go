package backlog

import (
	"os"
	"strings"
	"testing"
)

func TestBacklogManager_List(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	// Test with non-existent backlog (should create it)
	items, err := manager.List(false)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("Expected 0 items, got %d", len(items))
	}

	// Add some test items
	if err := manager.Add("Test item 1", "decomposed"); err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	if err := manager.Add("Test item 2", "refactor"); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// Test listing open items
	items, err = manager.List(false)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}

	// Complete one item
	if err := manager.Complete("Test item 1"); err != nil {
		t.Fatalf("Complete failed: %v", err)
	}

	// Test listing open items (should be 1)
	items, err = manager.List(false)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("Expected 1 open item, got %d", len(items))
	}

	// Test listing all items (should be 2)
	items, err = manager.List(true)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("Expected 2 total items, got %d", len(items))
	}
}

func TestBacklogManager_Add(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	tests := []struct {
		description string
		itemType    string
		wantErr     bool
	}{
		{"Test decomposed item", "decomposed", false},
		{"Test refactor item", "refactor", false},
		{"Test future item", "future", false},
		{"Test invalid type", "invalid", true},
	}

	for _, tt := range tests {
		err := manager.Add(tt.description, tt.itemType)
		if (err != nil) != tt.wantErr {
			t.Errorf("Add(%q, %q) error = %v, wantErr %v", tt.description, tt.itemType, err, tt.wantErr)
		}
	}

	// Verify items were added to correct sections
	content, err := manager.readBacklogContent()
	if err != nil {
		t.Fatalf("Failed to read backlog: %v", err)
	}

	if !strings.Contains(content, "- [ ] Test decomposed item") {
		t.Error("Decomposed item not found in backlog")
	}
	if !strings.Contains(content, "- [ ] Test refactor item") {
		t.Error("Refactor item not found in backlog")
	}
	if !strings.Contains(content, "- [ ] Test future item") {
		t.Error("Future item not found in backlog")
	}
}

func TestBacklogManager_Complete(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	// Add a test item
	if err := manager.Add("Test completion item", "decomposed"); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// Complete the item
	if err := manager.Complete("Test completion item"); err != nil {
		t.Fatalf("Complete failed: %v", err)
	}

	// Verify item was moved to completed section
	content, err := manager.readBacklogContent()
	if err != nil {
		t.Fatalf("Failed to read backlog: %v", err)
	}

	if strings.Contains(content, "- [ ] Test completion item") {
		t.Error("Item still appears as open")
	}
	if !strings.Contains(content, "- [x] Test completion item (Completed:") {
		t.Error("Item not found in completed section")
	}

	// Test completing non-existent item
	if err := manager.Complete("Non-existent item"); err == nil {
		t.Error("Expected error when completing non-existent item")
	}
}

func TestBacklogManager_ensureBacklogExists(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	// Backlog should not exist initially
	if _, err := os.Stat(manager.backlogPath); !os.IsNotExist(err) {
		t.Error("Backlog file should not exist initially")
	}

	// Call ensureBacklogExists
	if err := manager.ensureBacklogExists(); err != nil {
		t.Fatalf("ensureBacklogExists failed: %v", err)
	}

	// Backlog should now exist
	if _, err := os.Stat(manager.backlogPath); os.IsNotExist(err) {
		t.Error("Backlog file should exist after ensureBacklogExists")
	}

	// Verify content has the expected structure
	content, err := manager.readBacklogContent()
	if err != nil {
		t.Fatalf("Failed to read backlog: %v", err)
	}

	expectedSections := []string{
		"# Mission Backlog",
		"## DECOMPOSED INTENTS",
		"## REFACTORING OPPORTUNITIES", 
		"## FUTURE ENHANCEMENTS",
		"## COMPLETED",
	}

	for _, section := range expectedSections {
		if !strings.Contains(content, section) {
			t.Errorf("Expected section %q not found in backlog", section)
		}
	}
}

func TestBacklogManager_validateType(t *testing.T) {
	manager := NewManager("")

	validTypes := []string{"decomposed", "refactor", "future"}
	for _, validType := range validTypes {
		if err := manager.validateType(validType); err != nil {
			t.Errorf("validateType(%q) should be valid, got error: %v", validType, err)
		}
	}

	invalidTypes := []string{"invalid", "wrong", ""}
	for _, invalidType := range invalidTypes {
		if err := manager.validateType(invalidType); err == nil {
			t.Errorf("validateType(%q) should be invalid, got no error", invalidType)
		}
	}
}

func TestBacklogManager_getSectionHeader(t *testing.T) {
	manager := NewManager("")

	tests := []struct {
		itemType string
		expected string
	}{
		{"decomposed", "## DECOMPOSED INTENTS"},
		{"refactor", "## REFACTORING OPPORTUNITIES"},
		{"future", "## FUTURE ENHANCEMENTS"},
		{"invalid", ""},
	}

	for _, tt := range tests {
		result := manager.getSectionHeader(tt.itemType)
		if result != tt.expected {
			t.Errorf("getSectionHeader(%q) = %q, expected %q", tt.itemType, result, tt.expected)
		}
	}
}
