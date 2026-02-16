package backlog

import (
	"errors"
	"testing"
)

// mockProvider implements BacklogProvider for testing
type mockProvider struct {
	items     []string
	err       error
	added     []string
	completed []string
}

func (m *mockProvider) List(include, exclude []string) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.items, nil
}

func (m *mockProvider) Add(description, itemType string) error {
	if m.err != nil {
		return m.err
	}
	m.added = append(m.added, description)
	return nil
}

func (m *mockProvider) AddWithPattern(description, itemType, patternID string) error {
	return m.Add(description, itemType)
}

func (m *mockProvider) AddMultiple(descriptions []string, itemType string) error {
	for _, desc := range descriptions {
		if err := m.Add(desc, itemType); err != nil {
			return err
		}
	}
	return nil
}

func (m *mockProvider) Complete(itemText string) error {
	if m.err != nil {
		return m.err
	}
	m.completed = append(m.completed, itemText)
	return nil
}

func (m *mockProvider) Cleanup(itemType string) (int, error) {
	return 0, nil
}

func (m *mockProvider) GetPatternCount(patternID string) (int, error) {
	return 0, nil
}

func (m *mockProvider) Decompose(jsonInput string) error {
	return nil
}

func TestNewReconciler(t *testing.T) {
	fileProvider := &mockProvider{}
	beadsProvider := &mockProvider{}

	reconciler := NewReconciler(fileProvider, beadsProvider)

	if reconciler == nil {
		t.Fatal("Expected reconciler to be created, got nil")
	}
	if reconciler.fileProvider != fileProvider {
		t.Error("File provider not set correctly")
	}
	if reconciler.beadsProvider != beadsProvider {
		t.Error("Beads provider not set correctly")
	}
}

func TestDetectRogueEdits_NoEdits(t *testing.T) {
	fileProvider := &mockProvider{
		items: []string{
			"- [ ] Add new feature",
			"- [ ] Fix login bug",
		},
	}
	beadsProvider := &mockProvider{
		items: []string{
			"- [ ] Add new feature",
			"- [ ] Fix login bug",
		},
	}

	reconciler := NewReconciler(fileProvider, beadsProvider)
	edits, err := reconciler.DetectRogueEdits()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(edits) != 0 {
		t.Errorf("Expected 0 rogue edits, got %d", len(edits))
	}
}

func TestDetectRogueEdits_WithRogueEdits(t *testing.T) {
	fileProvider := &mockProvider{
		items: []string{
			"- [ ] Add new feature",
			"- [ ] Fix login bug",
			"- [ ] Refactor database layer",
		},
	}
	beadsProvider := &mockProvider{
		items: []string{
			"- [ ] Add new feature",
		},
	}

	reconciler := NewReconciler(fileProvider, beadsProvider)
	edits, err := reconciler.DetectRogueEdits()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(edits) != 2 {
		t.Errorf("Expected 2 rogue edits, got %d", len(edits))
	}

	// Verify the edits contain the missing items
	foundFix := false
	foundRefactor := false
	for _, edit := range edits {
		if edit.Description == "Fix login bug" {
			foundFix = true
			if edit.Type != ItemTypeBugfix {
				t.Errorf("Expected type %s, got %s", ItemTypeBugfix, edit.Type)
			}
		}
		if edit.Description == "Refactor database layer" {
			foundRefactor = true
			if edit.Type != ItemTypeRefactor {
				t.Errorf("Expected type %s, got %s", ItemTypeRefactor, edit.Type)
			}
		}
	}

	if !foundFix {
		t.Error("Missing 'Fix login bug' from rogue edits")
	}
	if !foundRefactor {
		t.Error("Missing 'Refactor database layer' from rogue edits")
	}
}

func TestDetectRogueEdits_WithCompletedItems(t *testing.T) {
	fileProvider := &mockProvider{
		items: []string{
			"- [ ] Add new feature",
			"- [x] Completed task (Completed: 2026-02-16)",
		},
	}
	beadsProvider := &mockProvider{
		items: []string{
			"- [ ] Add new feature",
		},
	}

	reconciler := NewReconciler(fileProvider, beadsProvider)
	edits, err := reconciler.DetectRogueEdits()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(edits) != 1 {
		t.Errorf("Expected 1 rogue edit, got %d", len(edits))
	}

	if edits[0].Description != "Completed task" {
		t.Errorf("Expected 'Completed task', got '%s'", edits[0].Description)
	}
	if edits[0].Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", edits[0].Status)
	}
}

func TestDetectRogueEdits_FileProviderError(t *testing.T) {
	fileProvider := &mockProvider{
		err: errors.New("file read error"),
	}
	beadsProvider := &mockProvider{}

	reconciler := NewReconciler(fileProvider, beadsProvider)
	_, err := reconciler.DetectRogueEdits()

	if err == nil {
		t.Error("Expected error from file provider")
	}
	if !errors.Is(err, fileProvider.err) {
		t.Errorf("Expected file provider error, got: %v", err)
	}
}

func TestDetectRogueEdits_BeadsProviderError(t *testing.T) {
	fileProvider := &mockProvider{
		items: []string{"- [ ] Item"},
	}
	beadsProvider := &mockProvider{
		err: errors.New("beads error"),
	}

	reconciler := NewReconciler(fileProvider, beadsProvider)
	_, err := reconciler.DetectRogueEdits()

	if err == nil {
		t.Error("Expected error from beads provider")
	}
	if !errors.Is(err, beadsProvider.err) {
		t.Errorf("Expected beads provider error, got: %v", err)
	}
}

func TestImportToBeads_Success(t *testing.T) {
	beadsProvider := &mockProvider{}
	fileProvider := &mockProvider{}
	reconciler := NewReconciler(fileProvider, beadsProvider)

	edits := []RogueEdit{
		{Description: "Add new feature", Type: ItemTypeFeature, Status: "open"},
		{Description: "Fix bug in login", Type: ItemTypeBugfix, Status: "open"},
	}

	count, err := reconciler.ImportToBeads(edits)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected 2 imports, got %d", count)
	}
	if len(beadsProvider.added) != 2 {
		t.Errorf("Expected 2 items added, got %d", len(beadsProvider.added))
	}
}

func TestImportToBeads_WithDependencies(t *testing.T) {
	beadsProvider := &mockProvider{}
	fileProvider := &mockProvider{}
	reconciler := NewReconciler(fileProvider, beadsProvider)

	edits := []RogueEdit{
		{
			Description:  "Implement user service",
			Type:         ItemTypeDecomposed,
			Status:       "open",
			Dependencies: []string{"Create database schema", "Setup authentication"},
		},
	}

	count, err := reconciler.ImportToBeads(edits)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 import, got %d", count)
	}
	if len(beadsProvider.added) != 1 {
		t.Errorf("Expected 1 item added, got %d", len(beadsProvider.added))
	}

	// Verify dependencies were preserved
	expected := "Implement user service (depends on: Create database schema, Setup authentication)"
	if beadsProvider.added[0] != expected {
		t.Errorf("Expected '%s', got '%s'", expected, beadsProvider.added[0])
	}
}

func TestImportToBeads_CompletedItems(t *testing.T) {
	beadsProvider := &mockProvider{}
	fileProvider := &mockProvider{}
	reconciler := NewReconciler(fileProvider, beadsProvider)

	edits := []RogueEdit{
		{Description: "Completed task", Type: ItemTypeFeature, Status: "completed"},
	}

	count, err := reconciler.ImportToBeads(edits)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 import, got %d", count)
	}
	if len(beadsProvider.completed) != 1 {
		t.Errorf("Expected 1 item completed, got %d", len(beadsProvider.completed))
	}
}

func TestImportToBeads_AddError(t *testing.T) {
	beadsProvider := &mockProvider{
		err: errors.New("add error"),
	}
	fileProvider := &mockProvider{}
	reconciler := NewReconciler(fileProvider, beadsProvider)

	edits := []RogueEdit{
		{Description: "New feature", Type: ItemTypeFeature, Status: "open"},
	}

	count, err := reconciler.ImportToBeads(edits)

	if err == nil {
		t.Error("Expected error from add")
	}
	if count != 0 {
		t.Errorf("Expected 0 imports, got %d", count)
	}
}

func TestNormalizeDescription(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "- [ ] Add new feature",
			expected: "Add new feature",
		},
		{
			input:    "- [x] Completed task",
			expected: "Completed task",
		},
		{
			input:    "- [ ] [PATTERN:helper-extract][COUNT:2] Extract helper function",
			expected: "Extract helper function",
		},
		{
			input:    "- [x] Task completed (Completed: 2026-02-16)",
			expected: "Task completed",
		},
		{
			input:    "  - [ ]   Multiple   spaces   ",
			expected: "Multiple spaces",
		},
	}

	for _, tt := range tests {
		result := normalizeDescription(tt.input)
		if result != tt.expected {
			t.Errorf("normalizeDescription(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestParseBacklogItem(t *testing.T) {
	tests := []struct {
		input          string
		expectedDesc   string
		expectedType   string
		expectedStatus string
		shouldParse    bool
	}{
		{
			input:          "- [ ] Add new feature",
			expectedDesc:   "Add new feature",
			expectedType:   ItemTypeFeature,
			expectedStatus: "open",
			shouldParse:    true,
		},
		{
			input:          "- [ ] Fix login bug",
			expectedDesc:   "Fix login bug",
			expectedType:   ItemTypeBugfix,
			expectedStatus: "open",
			shouldParse:    true,
		},
		{
			input:          "- [ ] Refactor database layer",
			expectedDesc:   "Refactor database layer",
			expectedType:   ItemTypeRefactor,
			expectedStatus: "open",
			shouldParse:    true,
		},
		{
			input:          "- [x] Completed feature",
			expectedDesc:   "Completed feature",
			expectedType:   ItemTypeFeature,
			expectedStatus: "completed",
			shouldParse:    true,
		},
		{
			input:          "- [ ] Task (from Epic: Main epic)",
			expectedDesc:   "Task (from Epic: Main epic)",
			expectedType:   ItemTypeDecomposed,
			expectedStatus: "open",
			shouldParse:    true,
		},
		{
			input:       "Not a backlog item",
			shouldParse: false,
		},
	}

	for _, tt := range tests {
		edit := parseBacklogItem(tt.input)
		if !tt.shouldParse {
			if edit != nil {
				t.Errorf("parseBacklogItem(%q) should return nil, got %+v", tt.input, edit)
			}
			continue
		}

		if edit == nil {
			t.Errorf("parseBacklogItem(%q) returned nil, expected valid edit", tt.input)
			continue
		}

		if edit.Description != tt.expectedDesc {
			t.Errorf("Description: got %q, want %q", edit.Description, tt.expectedDesc)
		}
		if edit.Type != tt.expectedType {
			t.Errorf("Type: got %q, want %q", edit.Type, tt.expectedType)
		}
		if edit.Status != tt.expectedStatus {
			t.Errorf("Status: got %q, want %q", edit.Status, tt.expectedStatus)
		}
	}
}

func TestExtractDependencies(t *testing.T) {
	tests := []struct {
		input        string
		expectedDesc string
		expectedDeps []string
	}{
		{
			input:        "Implement user service (depends on: Create database, Setup auth)",
			expectedDesc: "Implement user service",
			expectedDeps: []string{"Create database", "Setup auth"},
		},
		{
			input:        "Task without dependencies",
			expectedDesc: "Task without dependencies",
			expectedDeps: nil,
		},
		{
			input:        "Task (depends on: Single dep)",
			expectedDesc: "Task",
			expectedDeps: []string{"Single dep"},
		},
	}

	for _, tt := range tests {
		desc, deps := extractDependencies(tt.input)
		if desc != tt.expectedDesc {
			t.Errorf("Description: got %q, want %q", desc, tt.expectedDesc)
		}
		if len(deps) != len(tt.expectedDeps) {
			t.Errorf("Dependencies count: got %d, want %d", len(deps), len(tt.expectedDeps))
			continue
		}
		for i, dep := range deps {
			if dep != tt.expectedDeps[i] {
				t.Errorf("Dependency[%d]: got %q, want %q", i, dep, tt.expectedDeps[i])
			}
		}
	}
}

func TestInferTypeFromDescription(t *testing.T) {
	tests := []struct {
		description string
		expected    string
	}{
		{"Add new feature", ItemTypeFeature},
		{"Implement user service", ItemTypeFeature},
		{"Create new endpoint", ItemTypeFeature},
		{"Fix login bug", ItemTypeBugfix},
		{"Bug in database connection", ItemTypeBugfix},
		{"Error handling crash", ItemTypeBugfix},
		{"Refactor database layer", ItemTypeRefactor},
		{"Extract helper function", ItemTypeRefactor},
		{"Consolidate similar code", ItemTypeRefactor},
		{"Task (from Epic: Main epic)", ItemTypeDecomposed},
		{"Sub-intent for feature", ItemTypeDecomposed},
		{"Unclear future item", ItemTypeFuture},
		{"Generic task", ItemTypeFuture},
	}

	for _, tt := range tests {
		result := inferTypeFromDescription(tt.description)
		if result != tt.expected {
			t.Errorf("inferTypeFromDescription(%q) = %q, want %q", tt.description, result, tt.expected)
		}
	}
}

func TestRemovePatternMarker(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "[PATTERN:helper-extract][COUNT:2] Extract helper function",
			expected: "Extract helper function",
		},
		{
			input:    "No pattern marker here",
			expected: "No pattern marker here",
		},
		{
			input:    "[PATTERN:test][COUNT:1] Test item",
			expected: "Test item",
		},
	}

	for _, tt := range tests {
		result := removePatternMarker(tt.input)
		if result != tt.expected {
			t.Errorf("removePatternMarker(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestRemoveCompletionTimestamp(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "Completed task (Completed: 2026-02-16)",
			expected: "Completed task",
		},
		{
			input:    "Task without timestamp",
			expected: "Task without timestamp",
		},
		{
			input:    "Feature completed (Completed: 2025-12-31)",
			expected: "Feature completed",
		},
	}

	for _, tt := range tests {
		result := removeCompletionTimestamp(tt.input)
		if result != tt.expected {
			t.Errorf("removeCompletionTimestamp(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
