package mission

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestWriter_Write(t *testing.T) {
	fs := afero.NewMemMapFs()
	writer := NewWriter(fs)

	mission := &Mission{
		ID:        "test-456",
		Type:      "WET",
		Track:     2,
		Iteration: 1,
		Status:    "planned",
		Body:      "## INTENT\nTest body content\n",
	}

	path := "test-mission.md"
	err := writer.Write(path, mission)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Read back and verify
	data, err := afero.ReadFile(fs, path)
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}

	content := string(data)
	if !strings.HasPrefix(content, "---\n") {
		t.Error("Write() content should start with YAML frontmatter delimiter")
	}
	if !strings.Contains(content, "id: test-456") {
		t.Error("Write() content should contain id field")
	}
	if !strings.Contains(content, "status: planned") {
		t.Error("Write() content should contain status field")
	}
	if !strings.Contains(content, "## INTENT") {
		t.Error("Write() content should preserve body")
	}
}

func TestWriter_UpdateStatus(t *testing.T) {
	fs := afero.NewMemMapFs()
	path := "mission.md"

	// Create initial mission file
	initialContent := `---
id: test-789
type: WET
track: 2
iteration: 1
status: planned
---

## INTENT
Original body content

## SCOPE
file1.go
`
	afero.WriteFile(fs, path, []byte(initialContent), 0644)

	// Update status
	writer := NewWriter(fs)
	err := writer.UpdateStatus(path, "active")
	if err != nil {
		t.Fatalf("UpdateStatus() error = %v", err)
	}

	// Read back and verify
	reader := NewReader(fs)
	mission, err := reader.Read(path)
	if err != nil {
		t.Fatalf("Failed to read updated mission: %v", err)
	}

	if mission.Status != "active" {
		t.Errorf("UpdateStatus() status = %v, want active", mission.Status)
	}
	if mission.ID != "test-789" {
		t.Errorf("UpdateStatus() should preserve ID, got %v", mission.ID)
	}
	if !strings.Contains(mission.Body, "Original body content") {
		t.Error("UpdateStatus() should preserve body content")
	}
}

func TestWriter_UpdateStatusPreservesBody(t *testing.T) {
	fs := afero.NewMemMapFs()
	path := "mission.md"

	initialContent := `---
id: preserve-test
type: DRY
track: 3
iteration: 2
status: planned
parent_mission: parent-123
---

## INTENT
Complex body with multiple sections

## SCOPE
- file1.go
- file2.go
- file3.go

## PLAN
- [ ] Step 1
- [ ] Step 2

## VERIFICATION
go test ./...
`
	afero.WriteFile(fs, path, []byte(initialContent), 0644)

	writer := NewWriter(fs)
	err := writer.UpdateStatus(path, "completed")
	if err != nil {
		t.Fatalf("UpdateStatus() error = %v", err)
	}

	// Verify all content preserved
	data, err := afero.ReadFile(fs, path)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	content := string(data)
	checks := []string{
		"status: completed",
		"id: preserve-test",
		"parent_mission: parent-123",
		"Complex body with multiple sections",
		"## SCOPE",
		"file1.go",
		"## PLAN",
		"Step 1",
		"## VERIFICATION",
		"go test",
	}

	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Errorf("UpdateStatus() should preserve '%s' in content", check)
		}
	}
}

func TestWriter_CreateWithIntent(t *testing.T) {
	fs := afero.NewMemMapFs()
	writer := NewWriter(fs)

	path := "/test/mission.md"
	if err := writer.CreateWithIntent(path, "test-123", "Add user authentication"); err != nil {
		t.Fatalf("CreateWithIntent failed: %v", err)
	}

	mission, err := NewReader(fs).Read(path)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if mission.ID != "test-123" {
		t.Errorf("Expected ID 'test-123', got '%s'", mission.ID)
	}
	if mission.Status != "planning" {
		t.Errorf("Expected status 'planning', got '%s'", mission.Status)
	}
	if !strings.Contains(mission.Body, "Add user authentication") {
		t.Error("Body missing intent text")
	}
}

func TestWriter_UpdateSection(t *testing.T) {
	fs := afero.NewMemMapFs()
	writer := NewWriter(fs)

	mission := &Mission{
		ID:        "test-123",
		Status:    "planning",
		Iteration: 1,
		Body:      "## INTENT\nOld intent\n\n## SCOPE\nfile.go\n",
	}

	path := "/test/mission.md"
	if err := writer.Write(path, mission); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if err := writer.UpdateSection(path, "intent", "New intent"); err != nil {
		t.Fatalf("UpdateSection failed: %v", err)
	}

	updated, err := NewReader(fs).Read(path)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if !strings.Contains(updated.Body, "New intent") {
		t.Error("Body missing updated intent")
	}
	if strings.Contains(updated.Body, "Old intent") {
		t.Error("Body still contains old intent")
	}
}

func TestWriter_UpdateList(t *testing.T) {
	fs := afero.NewMemMapFs()
	writer := NewWriter(fs)

	mission := &Mission{
		ID:        "test-123",
		Status:    "planning",
		Iteration: 1,
		Body:      "## INTENT\nTest\n\n## SCOPE\n",
	}

	path := "/test/mission.md"
	if err := writer.Write(path, mission); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	items := []string{"file1.go", "file2.go"}
	if err := writer.UpdateList(path, "scope", items, false); err != nil {
		t.Fatalf("UpdateList failed: %v", err)
	}

	updated, err := NewReader(fs).Read(path)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if !strings.Contains(updated.Body, "file1.go") {
		t.Error("Body missing file1.go")
	}
	if !strings.Contains(updated.Body, "file2.go") {
		t.Error("Body missing file2.go")
	}
}

func TestWriter_UpdateListAppend(t *testing.T) {
	fs := afero.NewMemMapFs()
	writer := NewWriter(fs)

	mission := &Mission{
		ID:        "test-123",
		Status:    "planning",
		Iteration: 1,
		Body:      "## INTENT\nTest\n\n## SCOPE\nexisting1.go\nexisting2.go\n",
	}

	path := "/test/mission.md"
	if err := writer.Write(path, mission); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Test append mode
	newItems := []string{"new1.go", "new2.go"}
	if err := writer.UpdateList(path, "scope", newItems, true); err != nil {
		t.Fatalf("UpdateList append failed: %v", err)
	}

	updated, err := NewReader(fs).Read(path)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	// Should contain both existing and new items
	if !strings.Contains(updated.Body, "existing1.go") {
		t.Error("Body missing existing1.go")
	}
	if !strings.Contains(updated.Body, "existing2.go") {
		t.Error("Body missing existing2.go")
	}
	if !strings.Contains(updated.Body, "new1.go") {
		t.Error("Body missing new1.go")
	}
	if !strings.Contains(updated.Body, "new2.go") {
		t.Error("Body missing new2.go")
	}
}

func TestWriter_UpdateFrontmatter(t *testing.T) {
	fs := afero.NewMemMapFs()
	writer := NewWriter(fs)

	mission := &Mission{
		ID:        "test-123",
		Status:    "planning",
		Track:     2,
		Type:      "WET",
		Iteration: 1,
		Body:      "## INTENT\nTest\n",
	}

	path := "/test/mission.md"
	if err := writer.Write(path, mission); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	pairs := []string{"track=3", "type=DRY"}
	if err := writer.UpdateFrontmatter(path, pairs); err != nil {
		t.Fatalf("UpdateFrontmatter failed: %v", err)
	}

	updated, err := NewReader(fs).Read(path)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if updated.Track != 3 {
		t.Errorf("Expected track 3, got %d", updated.Track)
	}
	if updated.Type != "DRY" {
		t.Errorf("Expected type 'DRY', got '%s'", updated.Type)
	}
}
