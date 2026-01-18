package mission

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestWriter_Write(t *testing.T) {
	fs := afero.NewMemMapFs()
	path := "test-mission.md"
	writer := NewWriterWithPath(fs, path)

	mission := &Mission{
		ID:        "test-456",
		Type:      "WET",
		Track:     2,
		Iteration: 1,
		Status:    "planned",
		Body:      "## INTENT\nTest body content\n",
	}

	err := writer.Write(mission)
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
	writer := NewWriterWithPath(fs, path)
	err := writer.UpdateStatus("active")
	if err != nil {
		t.Fatalf("UpdateStatus() error = %v", err)
	}

	// Read back and verify
	reader := NewReader(fs, path)
	mission, err := reader.Read()
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

	writer := NewWriterWithPath(fs, path)
	err := writer.UpdateStatus("completed")
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
	path := "/test/mission.md"
	writer := NewWriterWithPath(fs, path)

	if err := writer.CreateWithIntent("test-123", "Add user authentication"); err != nil {
		t.Fatalf("CreateWithIntent failed: %v", err)
	}

	mission, err := NewReader(fs, path).Read()
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
	path := "/test/mission.md"
	writer := NewWriterWithPath(fs, path)

	mission := &Mission{
		ID:        "test-123",
		Status:    "planning",
		Iteration: 1,
		Body:      "## INTENT\nOld intent\n\n## SCOPE\nfile.go\n",
	}

	if err := writer.Write(mission); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if err := writer.UpdateSection("intent", "New intent"); err != nil {
		t.Fatalf("UpdateSection failed: %v", err)
	}

	updated, err := NewReader(fs, path).Read()
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if !strings.Contains(updated.Body, "New intent") {
		t.Error("Body missing updated intent")
	}
	if strings.Contains(updated.Body, "Old intent") {
		t.Error("Body still contains old intent")
	}

	// Verify exactly one empty line between sections
	expected := "## INTENT\nNew intent\n\n## SCOPE\nfile.go\n"
	if updated.Body != expected {
		t.Errorf("Incorrect spacing. Expected:\n%q\nGot:\n%q", expected, updated.Body)
	}
}

func TestWriter_UpdateSectionCreatesNew(t *testing.T) {
	fs := afero.NewMemMapFs()
	path := "/test/mission.md"
	writer := NewWriterWithPath(fs, path)

	mission := &Mission{
		ID:        "test-123",
		Status:    "planning",
		Iteration: 1,
		Body:      "## INTENT\nTest intent\n\n## SCOPE\nfile.go\n",
	}

	if err := writer.Write(mission); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Update non-existent section - should create it
	if err := writer.UpdateSection("verification", "go test ./..."); err != nil {
		t.Fatalf("UpdateSection failed to create new section: %v", err)
	}

	updated, err := NewReader(fs, path).Read()
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if !strings.Contains(updated.Body, "## VERIFICATION") {
		t.Error("Body missing new VERIFICATION section header")
	}
	if !strings.Contains(updated.Body, "go test ./...") {
		t.Error("Body missing verification content")
	}
	// Should preserve existing content
	if !strings.Contains(updated.Body, "Test intent") {
		t.Error("Body missing original intent")
	}
	if !strings.Contains(updated.Body, "file.go") {
		t.Error("Body missing original scope")
	}
}

func TestWriter_UpdateList(t *testing.T) {
	fs := afero.NewMemMapFs()
	path := "/test/mission.md"
	writer := NewWriterWithPath(fs, path)

	mission := &Mission{
		ID:        "test-123",
		Status:    "planning",
		Iteration: 1,
		Body:      "## INTENT\nTest\n\n## SCOPE\n\n## PLAN\n- [ ] Step 1\n",
	}

	if err := writer.Write(mission); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	items := []string{"file1.go", "file2.go"}
	if err := writer.UpdateList("scope", items, false); err != nil {
		t.Fatalf("UpdateList failed: %v", err)
	}

	updated, err := NewReader(fs, path).Read()
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if !strings.Contains(updated.Body, "file1.go") {
		t.Error("Body missing file1.go")
	}
	if !strings.Contains(updated.Body, "file2.go") {
		t.Error("Body missing file2.go")
	}

	// Verify exactly one empty line between sections
	expected := "## INTENT\nTest\n\n## SCOPE\nfile1.go\nfile2.go\n\n## PLAN\n- [ ] Step 1\n"
	if updated.Body != expected {
		t.Errorf("Incorrect spacing. Expected:\n%q\nGot:\n%q", expected, updated.Body)
	}
}

func TestWriter_UpdateListAppend(t *testing.T) {
	fs := afero.NewMemMapFs()
	path := "/test/mission.md"
	writer := NewWriterWithPath(fs, path)

	mission := &Mission{
		ID:        "test-123",
		Status:    "planning",
		Iteration: 1,
		Body:      "## INTENT\nTest\n\n## SCOPE\nexisting1.go\nexisting2.go\n",
	}

	if err := writer.Write(mission); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Test append mode
	newItems := []string{"new1.go", "new2.go"}
	if err := writer.UpdateList("scope", newItems, true); err != nil {
		t.Fatalf("UpdateList append failed: %v", err)
	}

	updated, err := NewReader(fs, path).Read()
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

func TestWriter_UpdateList_PlanEntries(t *testing.T) {
	fs := afero.NewMemMapFs()
	path := "/test/mission.md"
	writer := NewWriterWithPath(fs, path)

	mission := &Mission{
		ID:        "test-123",
		Status:    "planning",
		Iteration: 1,
		Body:      "## INTENT\nTest\n\n## PLAN\n- [ ] 1. First step\n- [ ] 2. Second step\n\n## VERIFICATION\ngo test\n",
	}

	if err := writer.Write(mission); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Update plan with new items
	newItems := []string{"3. Third step", "4. Fourth step"}
	if err := writer.UpdateList("plan", newItems, false); err != nil {
		t.Fatalf("UpdateList plan failed: %v", err)
	}

	updated, err := NewReader(fs, path).Read()
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	// Should contain properly formatted plan entries
	if !strings.Contains(updated.Body, "- [ ] 3. Third step") {
		t.Error("Body missing properly formatted third step")
	}
	if !strings.Contains(updated.Body, "- [ ] 4. Fourth step") {
		t.Error("Body missing properly formatted fourth step")
	}

	// Should preserve VERIFICATION section
	if !strings.Contains(updated.Body, "## VERIFICATION") {
		t.Error("Body missing VERIFICATION section")
	}
	if !strings.Contains(updated.Body, "go test") {
		t.Error("Body missing verification content")
	}
}

func TestWriter_UpdateList_PreservesSubsequentSections(t *testing.T) {
	fs := afero.NewMemMapFs()
	path := "/test/mission.md"
	writer := NewWriterWithPath(fs, path)

	mission := &Mission{
		ID:        "test-123",
		Status:    "planning",
		Iteration: 1,
		Body:      "## INTENT\nTest intent\n\n## SCOPE\nold1.go\nold2.go\n\n## PLAN\n- [ ] Step 1\n- [ ] Step 2\n\n## VERIFICATION\ngo test ./...\n",
	}

	if err := writer.Write(mission); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Update scope section
	newItems := []string{"new1.go", "new2.go"}
	if err := writer.UpdateList("scope", newItems, false); err != nil {
		t.Fatalf("UpdateList scope failed: %v", err)
	}

	updated, err := NewReader(fs, path).Read()
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	// Should contain new scope items
	if !strings.Contains(updated.Body, "new1.go") {
		t.Error("Body missing new1.go")
	}
	if !strings.Contains(updated.Body, "new2.go") {
		t.Error("Body missing new2.go")
	}

	// Should preserve all subsequent sections intact
	if !strings.Contains(updated.Body, "## PLAN") {
		t.Error("Body missing PLAN section")
	}
	if !strings.Contains(updated.Body, "- [ ] Step 1") {
		t.Error("Body missing Step 1")
	}
	if !strings.Contains(updated.Body, "- [ ] Step 2") {
		t.Error("Body missing Step 2")
	}
	if !strings.Contains(updated.Body, "## VERIFICATION") {
		t.Error("Body missing VERIFICATION section")
	}
	if !strings.Contains(updated.Body, "go test ./...") {
		t.Error("Body missing verification content")
	}
}

func TestWriter_MarkPlanStepComplete(t *testing.T) {
	fs := afero.NewMemMapFs()
	path := "mission.md"
	writer := NewWriterWithPath(fs, path)

	mission := &Mission{
		ID:     "test-123",
		Status: "active",
		Body:   "## PLAN\n- [ ] Step 1\n- [ ] Step 2\n- [ ] Step 3",
	}

	if err := writer.Write(mission); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Test marking step 2 complete
	if err := writer.MarkPlanStepComplete(2, "", ""); err != nil {
		t.Fatalf("MarkPlanStepComplete failed: %v", err)
	}

	updated, err := NewReader(fs, path).Read()
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	lines := strings.Split(updated.Body, "\n")
	found := false
	for _, line := range lines {
		if strings.Contains(line, "- [x] Step 2") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Step 2 was not marked as complete")
	}

	// Verify other steps are untouched
	if !strings.Contains(updated.Body, "- [ ] Step 1") {
		t.Error("Step 1 should remain incomplete")
	}
	if !strings.Contains(updated.Body, "- [ ] Step 3") {
		t.Error("Step 3 should remain incomplete")
	}
}

func TestWriter_MarkPlanStepComplete_InvalidStep(t *testing.T) {
	fs := afero.NewMemMapFs()
	path := "mission.md"
	writer := NewWriterWithPath(fs, path)

	mission := &Mission{
		ID:     "test-123",
		Status: "active",
		Body:   "## PLAN\n- [ ] Step 1",
	}

	if err := writer.Write(mission); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if err := writer.MarkPlanStepComplete(99, "", ""); err == nil {
		t.Error("Expected error for invalid step, got nil")
	}
}

func TestWriter_MarkPlanStepComplete_AfterOtherStepsComplete(t *testing.T) {
	fs := afero.NewMemMapFs()
	path := "mission.md"
	writer := NewWriterWithPath(fs, path)

	mission := &Mission{
		ID:     "test-123",
		Status: "active",
		Body:   "## PLAN\n- [ ] Step 1\n- [ ] Step 2\n- [ ] Step 3",
	}

	if err := writer.Write(mission); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Mark step 1 complete first
	if err := writer.MarkPlanStepComplete(1, "", ""); err != nil {
		t.Fatalf("MarkPlanStepComplete(1) failed: %v", err)
	}

	// Now mark step 3 complete - this should work even though step 2 is not complete
	if err := writer.MarkPlanStepComplete(3, "", ""); err != nil {
		t.Fatalf("MarkPlanStepComplete(3) after step 1 complete failed: %v", err)
	}

	updated, err := NewReader(fs, path).Read()
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	lines := strings.Split(updated.Body, "\n")
	completeCount := 0
	for _, line := range lines {
		if strings.Contains(line, "- [x]") {
			completeCount++
		}
	}

	if completeCount != 2 {
		t.Errorf("Expected 2 completed steps, got %d", completeCount)
	}

	// Verify step 3 is marked complete
	foundStep3 := false
	for _, line := range lines {
		if strings.Contains(line, "- [x] Step 3") {
			foundStep3 = true
			break
		}
	}
	if !foundStep3 {
		t.Error("Step 3 was not marked as complete")
	}
}

func TestWriter_UpdateFrontmatter(t *testing.T) {
	fs := afero.NewMemMapFs()
	path := "/test/mission.md"
	writer := NewWriterWithPath(fs, path)

	mission := &Mission{
		ID:        "test-123",
		Status:    "planning",
		Track:     2,
		Type:      "WET",
		Iteration: 1,
		Body:      "## INTENT\nTest\n",
	}

	if err := writer.Write(mission); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	pairs := []string{"track=3", "type=DRY"}
	if err := writer.UpdateFrontmatter(pairs); err != nil {
		t.Fatalf("UpdateFrontmatter failed: %v", err)
	}

	updated, err := NewReader(fs, path).Read()
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

func TestWriter_WriteWithDomains(t *testing.T) {
	fs := afero.NewMemMapFs()
	path := "test-mission-with-domains.md"
	writer := NewWriterWithPath(fs, path)

	mission := &Mission{
		ID:        "test-789",
		Type:      "WET",
		Track:     2,
		Iteration: 1,
		Status:    "planned",
		Domains:   "security,performance",
		Body:      "## INTENT\nTest body content\n",
	}

	err := writer.Write(mission)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Read back and verify
	data, err := afero.ReadFile(fs, path)
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "domains: security,performance") {
		t.Error("Write() content should contain domains field")
	}

	// Verify that type field is still present and separate
	if !strings.Contains(content, "type: WET") {
		t.Error("Write() content should contain type field")
	}
}

func TestWriter_UpdateFrontmatter_Domains(t *testing.T) {
	fs := afero.NewMemMapFs()
	path := "/test/mission.md"
	writer := NewWriterWithPath(fs, path)

	mission := &Mission{
		ID:        "test-456",
		Status:    "planned",
		Track:     2,
		Type:      "WET",
		Iteration: 1,
		Body:      "## INTENT\nTest\n",
	}

	if err := writer.Write(mission); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Update domains field
	pairs := []string{"domains=security,performance"}
	if err := writer.UpdateFrontmatter(pairs); err != nil {
		t.Fatalf("UpdateFrontmatter failed: %v", err)
	}

	updated, err := NewReader(fs, path).Read()
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if updated.Domains != "security,performance" {
		t.Errorf("Expected domains 'security,performance', got '%s'", updated.Domains)
	}

	// Verify that type field is NOT affected by domains update
	if updated.Type != "WET" {
		t.Errorf("Expected type 'WET' (unchanged), got '%s'", updated.Type)
	}
}
