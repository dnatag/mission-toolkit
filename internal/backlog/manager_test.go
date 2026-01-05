package backlog

import (
	"fmt"
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

func TestBacklogManager_AddToEmptySection(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	// Create backlog with empty sections
	if err := manager.ensureBacklogExists(); err != nil {
		t.Fatalf("ensureBacklogExists failed: %v", err)
	}

	// Add item to empty decomposed section
	if err := manager.Add("First decomposed item", "decomposed"); err != nil {
		t.Fatalf("Add to empty section failed: %v", err)
	}

	// Verify item was added correctly
	items, err := manager.List(false)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(items))
	}
	if !strings.Contains(items[0], "First decomposed item") {
		t.Error("Item not found in list")
	}

	// Add second item to same section
	if err := manager.Add("Second decomposed item", "decomposed"); err != nil {
		t.Fatalf("Add second item failed: %v", err)
	}

	// Verify both items exist
	items, err = manager.List(false)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}
}

func TestBacklogManager_ExistingBacklogWithEmptySections(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	// Create backlog file with custom content (empty sections)
	customBacklog := `# Mission Backlog

## DECOMPOSED INTENTS
*This section lists atomic tasks that were broken down from a larger user-defined epic.*

## REFACTORING OPPORTUNITIES
*This section lists technical debt and refactoring opportunities identified by the AI during planning or execution.*

## FUTURE ENHANCEMENTS
*This section is for user-defined ideas and future feature requests.*

## COMPLETED
*History of completed backlog items.*
`
	if err := manager.writeBacklogContent(customBacklog); err != nil {
		t.Fatalf("Failed to write custom backlog: %v", err)
	}

	// Add items to each empty section
	testCases := []struct {
		description string
		itemType    string
	}{
		{"Test decomposed task", "decomposed"},
		{"Test refactor opportunity", "refactor"},
		{"Test future enhancement", "future"},
	}

	for _, tc := range testCases {
		if err := manager.Add(tc.description, tc.itemType); err != nil {
			t.Errorf("Add(%q, %q) failed: %v", tc.description, tc.itemType, err)
		}
	}

	// Verify all items were added
	items, err := manager.List(false)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
	}

	// Verify items are in correct sections
	content, err := manager.readBacklogContent()
	if err != nil {
		t.Fatalf("Failed to read backlog: %v", err)
	}

	expectedItems := []string{
		"- [ ] Test decomposed task",
		"- [ ] Test refactor opportunity",
		"- [ ] Test future enhancement",
	}

	for _, item := range expectedItems {
		if !strings.Contains(content, item) {
			t.Errorf("Expected item %q not found in backlog", item)
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

func TestBacklogManager_AddMultiple(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	// Test adding multiple items
	items := []string{"Multi item 1", "Multi item 2", "Multi item 3"}
	if err := manager.AddMultiple(items, "decomposed"); err != nil {
		t.Fatalf("AddMultiple failed: %v", err)
	}

	// Verify all items were added
	content, err := manager.readBacklogContent()
	if err != nil {
		t.Fatalf("Failed to read backlog: %v", err)
	}

	for _, item := range items {
		if !strings.Contains(content, fmt.Sprintf("- [ ] %s", item)) {
			t.Errorf("Expected item %q not found in backlog", item)
		}
	}

	// Verify count
	storedItems, err := manager.List(false)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(storedItems) != 3 {
		t.Errorf("Expected 3 items, got %d", len(storedItems))
	}

	// Test adding multiple items to different sections
	refactorItems := []string{"Refactor A", "Refactor B"}
	if err := manager.AddMultiple(refactorItems, "refactor"); err != nil {
		t.Fatalf("AddMultiple to refactor failed: %v", err)
	}

	// Verify refactor items
	content, err = manager.readBacklogContent()
	if err != nil {
		t.Fatalf("Failed to read backlog: %v", err)
	}

	for _, item := range refactorItems {
		if !strings.Contains(content, fmt.Sprintf("- [ ] %s", item)) {
			t.Errorf("Expected refactor item %q not found in backlog", item)
		}
	}

	// Test invalid type
	if err := manager.AddMultiple(items, "invalid"); err == nil {
		t.Error("AddMultiple with invalid type should fail")
	}

	// Test empty slice (should not fail but add nothing)
	if err := manager.AddMultiple([]string{}, "decomposed"); err != nil {
		t.Errorf("AddMultiple with empty slice should not fail: %v", err)
	}
}
