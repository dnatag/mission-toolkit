package mission

import (
	"strings"
	"testing"

	"github.com/dnatag/mission-toolkit/internal/plan"
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

func TestWriter_CreateFromPlan(t *testing.T) {
	fs := afero.NewMemMapFs()
	writer := NewWriter(fs)

	spec := &plan.PlanSpec{
		Intent: "Test mission intent",
		Type:   "WET",
		Scope: []string{
			"file1.go",
			"file2.go",
		},
		Plan: []string{
			"1. Step one",
			"2. Step two",
		},
		Verification: "go test ./...",
	}

	path := "mission.md"
	err := writer.CreateFromPlan(path, "test-123", 2, spec)
	if err != nil {
		t.Fatalf("CreateFromPlan() error = %v", err)
	}

	// Read back and verify
	reader := NewReader(fs)
	mission, err := reader.Read(path)
	if err != nil {
		t.Fatalf("Failed to read created mission: %v", err)
	}

	if mission.ID != "test-123" {
		t.Errorf("CreateFromPlan() ID = %v, want test-123", mission.ID)
	}
	if mission.Type != "WET" {
		t.Errorf("CreateFromPlan() Type = %v, want WET", mission.Type)
	}
	if mission.Track != 2 {
		t.Errorf("CreateFromPlan() Track = %v, want 2", mission.Track)
	}
	if mission.Status != "planned" {
		t.Errorf("CreateFromPlan() Status = %v, want planned", mission.Status)
	}
	if !strings.Contains(mission.Body, "Test mission intent") {
		t.Error("CreateFromPlan() should include intent in body")
	}
	if !strings.Contains(mission.Body, "file1.go") {
		t.Error("CreateFromPlan() should include scope files in body")
	}
	if !strings.Contains(mission.Body, "go test") {
		t.Error("CreateFromPlan() should include verification in body")
	}
}
