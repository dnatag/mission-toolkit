package plan

import (
	"testing"

	"github.com/spf13/afero"
)

func TestLoadPlanSpec(t *testing.T) {
	// Create in-memory filesystem
	fs := afero.NewMemMapFs()
	planFile := "/tmp/plan.json"

	// Create test plan file
	testJSON := `{
		"intent": "Test intent",
		"scope": ["file1.go", "file2.go"],
		"domain": ["security"],
		"plan": ["step 1", "step 2"],
		"verification": "go test"
	}`

	if err := afero.WriteFile(fs, planFile, []byte(testJSON), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test loading
	spec, err := LoadPlanSpec(fs, planFile)
	if err != nil {
		t.Fatalf("LoadPlanSpec() error = %v", err)
	}

	if spec.Intent != "Test intent" {
		t.Errorf("Expected intent 'Test intent', got '%s'", spec.Intent)
	}
	if len(spec.Scope) != 2 {
		t.Errorf("Expected 2 scope items, got %d", len(spec.Scope))
	}
}

func TestSavePlanSpec(t *testing.T) {
	fs := afero.NewMemMapFs()
	planFile := "/tmp/plan.json"

	spec := &PlanSpec{
		Intent:       "Test intent",
		Scope:        []string{"file1.go", "file2.go"},
		Domain:       []string{"security"},
		Plan:         []string{"step 1", "step 2"},
		Verification: "go test",
	}

	// Test saving
	if err := SavePlanSpec(fs, spec, planFile); err != nil {
		t.Fatalf("SavePlanSpec() error = %v", err)
	}

	// Verify file exists
	exists, err := afero.Exists(fs, planFile)
	if err != nil || !exists {
		t.Errorf("Plan file was not created")
	}

	// Test loading back
	loadedSpec, err := LoadPlanSpec(fs, planFile)
	if err != nil {
		t.Fatalf("Failed to load saved plan: %v", err)
	}

	if loadedSpec.Intent != spec.Intent {
		t.Errorf("Intent mismatch after save/load")
	}
}
