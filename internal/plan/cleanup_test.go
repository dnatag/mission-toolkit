package plan

import (
	"testing"

	"github.com/spf13/afero"
)

func TestCleanupStalePlanFiles(t *testing.T) {
	// Create in-memory filesystem
	fs := afero.NewMemMapFs()
	tmpDir := "/tmp"
	planFile := "/tmp/plan.json"

	// Test cleanup when no file exists
	err := CleanupStalePlanFiles(fs, tmpDir)
	if err != nil {
		t.Errorf("CleanupStalePlanFiles() should not error when no file exists: %v", err)
	}

	// Create a stale plan file
	if err := afero.WriteFile(fs, planFile, []byte(`{"test": "data"}`), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Verify file exists
	if !HasStalePlanFile(fs, tmpDir) {
		t.Error("HasStalePlanFile() should return true when file exists")
	}

	// Test cleanup
	err = CleanupStalePlanFiles(fs, tmpDir)
	if err != nil {
		t.Errorf("CleanupStalePlanFiles() error = %v", err)
	}

	// Verify file is removed
	if HasStalePlanFile(fs, tmpDir) {
		t.Error("HasStalePlanFile() should return false after cleanup")
	}
}

func TestHasStalePlanFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	tmpDir := "/tmp"
	planFile := "/tmp/plan.json"

	// Test when file doesn't exist
	if HasStalePlanFile(fs, tmpDir) {
		t.Error("HasStalePlanFile() should return false when file doesn't exist")
	}

	// Create file
	if err := afero.WriteFile(fs, planFile, []byte(`{"test": "data"}`), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test when file exists
	if !HasStalePlanFile(fs, tmpDir) {
		t.Error("HasStalePlanFile() should return true when file exists")
	}
}
