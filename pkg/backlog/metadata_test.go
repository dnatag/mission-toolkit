package backlog

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestReadBacklogWithMetadata(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	// Create backlog with metadata
	content := "## FEATURES\n- [ ] Test feature"
	if err := manager.writeBacklogWithMetadata(content, "Test action"); err != nil {
		t.Fatalf("writeBacklogWithMetadata failed: %v", err)
	}

	// Read back
	body, metadata, err := manager.readBacklogWithMetadata()
	if err != nil {
		t.Fatalf("readBacklogWithMetadata failed: %v", err)
	}

	if body != content {
		t.Errorf("expected body %q, got %q", content, body)
	}

	if metadata.LastAction != "Test action" {
		t.Errorf("expected last_action 'Test action', got %q", metadata.LastAction)
	}

	if metadata.LastUpdated.IsZero() {
		t.Error("expected last_updated to be set")
	}
}

func TestReadBacklogWithMetadata_LegacyFormat(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	// Create legacy backlog without frontmatter
	content := "# Mission Backlog\n\n## FEATURES\n- [ ] Test feature"
	if err := os.WriteFile(manager.backlogPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write legacy backlog: %v", err)
	}

	// Read back - should handle missing metadata gracefully
	body, metadata, err := manager.readBacklogWithMetadata()
	if err != nil {
		t.Fatalf("readBacklogWithMetadata failed: %v", err)
	}

	if body != content {
		t.Errorf("expected body to match original content")
	}

	if !metadata.LastUpdated.IsZero() {
		t.Error("expected last_updated to be zero for legacy format")
	}

	if metadata.LastAction != "" {
		t.Error("expected last_action to be empty for legacy format")
	}
}

func TestWriteBacklogWithMetadata(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	content := "## FEATURES\n- [ ] Test feature"
	action := "Added test feature"

	if err := manager.writeBacklogWithMetadata(content, action); err != nil {
		t.Fatalf("writeBacklogWithMetadata failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(manager.backlogPath); os.IsNotExist(err) {
		t.Fatal("backlog file was not created")
	}

	// Read back and verify
	body, metadata, err := manager.readBacklogWithMetadata()
	if err != nil {
		t.Fatalf("readBacklogWithMetadata failed: %v", err)
	}

	if body != content {
		t.Errorf("expected body %q, got %q", content, body)
	}

	if metadata.LastAction != action {
		t.Errorf("expected last_action %q, got %q", action, metadata.LastAction)
	}

	// Verify timestamp is recent (within last minute)
	if time.Since(metadata.LastUpdated) > time.Minute {
		t.Errorf("last_updated timestamp is too old: %v", metadata.LastUpdated)
	}
}

func TestWriteBacklogWithMetadata_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	nestedDir := filepath.Join(tmpDir, "nested", "path")
	manager := NewManager(nestedDir)

	content := "## FEATURES\n- [ ] Test"
	if err := manager.writeBacklogWithMetadata(content, "Test"); err != nil {
		t.Fatalf("writeBacklogWithMetadata failed: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(nestedDir); os.IsNotExist(err) {
		t.Error("nested directory was not created")
	}

	// Verify file exists
	if _, err := os.Stat(manager.backlogPath); os.IsNotExist(err) {
		t.Error("backlog file was not created")
	}
}
