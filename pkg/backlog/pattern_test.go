package backlog

import (
	"testing"
)

func TestGetPatternCount(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	// Create backlog with pattern
	if err := manager.ensureBacklogExists(); err != nil {
		t.Fatalf("ensureBacklogExists failed: %v", err)
	}

	// Add pattern with count
	if err := manager.AddWithPattern("Extract validation logic", "refactor", "validation-extract"); err != nil {
		t.Fatalf("AddWithPattern failed: %v", err)
	}

	// Get count
	count, err := manager.GetPatternCount("validation-extract")
	if err != nil {
		t.Fatalf("GetPatternCount failed: %v", err)
	}

	if count != 2 {
		t.Errorf("expected count 2, got %d", count)
	}
}

func TestGetPatternCount_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	if err := manager.ensureBacklogExists(); err != nil {
		t.Fatalf("ensureBacklogExists failed: %v", err)
	}

	count, err := manager.GetPatternCount("nonexistent")
	if err != nil {
		t.Fatalf("GetPatternCount failed: %v", err)
	}

	if count != 0 {
		t.Errorf("expected count 0 for nonexistent pattern, got %d", count)
	}
}

func TestIncrementPatternCount(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	if err := manager.ensureBacklogExists(); err != nil {
		t.Fatalf("ensureBacklogExists failed: %v", err)
	}

	// Add initial pattern
	if err := manager.AddWithPattern("Extract helper", "refactor", "helper-extract"); err != nil {
		t.Fatalf("AddWithPattern failed: %v", err)
	}

	// Verify initial count
	count, _ := manager.GetPatternCount("helper-extract")
	if count != 2 {
		t.Fatalf("expected initial count 2, got %d", count)
	}

	// Increment
	if err := manager.incrementPatternCount("helper-extract"); err != nil {
		t.Fatalf("incrementPatternCount failed: %v", err)
	}

	// Verify incremented count
	count, err := manager.GetPatternCount("helper-extract")
	if err != nil {
		t.Fatalf("GetPatternCount failed: %v", err)
	}

	if count != 3 {
		t.Errorf("expected count 3 after increment, got %d", count)
	}
}

func TestIncrementPatternCount_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	if err := manager.ensureBacklogExists(); err != nil {
		t.Fatalf("ensureBacklogExists failed: %v", err)
	}

	err := manager.incrementPatternCount("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent pattern, got nil")
	}

	if err != nil && err.Error() != "pattern not found: nonexistent" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestAddWithPattern_IncrementExisting(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	if err := manager.ensureBacklogExists(); err != nil {
		t.Fatalf("ensureBacklogExists failed: %v", err)
	}

	// Add first occurrence
	if err := manager.AddWithPattern("Extract config", "refactor", "config-extract"); err != nil {
		t.Fatalf("AddWithPattern failed: %v", err)
	}

	count, _ := manager.GetPatternCount("config-extract")
	if count != 2 {
		t.Fatalf("expected count 2 after first add, got %d", count)
	}

	// Add second occurrence - should increment
	if err := manager.AddWithPattern("Extract config again", "refactor", "config-extract"); err != nil {
		t.Fatalf("AddWithPattern failed: %v", err)
	}

	count, err := manager.GetPatternCount("config-extract")
	if err != nil {
		t.Fatalf("GetPatternCount failed: %v", err)
	}

	if count != 3 {
		t.Errorf("expected count 3 after second add, got %d", count)
	}
}
