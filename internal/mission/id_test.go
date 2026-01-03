package mission

import (
	"testing"

	"github.com/spf13/afero"
)

func TestIDService_GetOrCreateID(t *testing.T) {
	// Create in-memory filesystem
	fs := afero.NewMemMapFs()
	service := NewIDService(fs, "/tmp")

	// First call should create new ID
	id1, err := service.GetOrCreateID()
	if err != nil {
		t.Fatalf("GetOrCreateID() error = %v", err)
	}

	if id1 == "" {
		t.Error("Expected non-empty ID")
	}

	// Second call should return same ID
	id2, err := service.GetOrCreateID()
	if err != nil {
		t.Fatalf("GetOrCreateID() error = %v", err)
	}

	if id1 != id2 {
		t.Errorf("Expected same ID, got %s and %s", id1, id2)
	}
}

func TestIDService_CleanupStaleID(t *testing.T) {
	fs := afero.NewMemMapFs()
	service := NewIDService(fs, "/tmp")

	// Create stale ID file
	idPath := "/tmp/id"
	err := afero.WriteFile(fs, idPath, []byte("20251231114608-1234"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test ID file: %v", err)
	}

	// Cleanup should remove stale ID
	err = service.CleanupStaleID()
	if err != nil {
		t.Errorf("CleanupStaleID() error = %v", err)
	}

	// Verify ID file is removed
	exists, _ := afero.Exists(fs, idPath)
	if exists {
		t.Error("Expected ID file to be removed")
	}
}

func TestIDService_GetCurrentID(t *testing.T) {
	fs := afero.NewMemMapFs()
	service := NewIDService(fs, "/tmp")

	// Create mission.md with ID in YAML frontmatter
	missionPath := "/tmp/mission.md"
	missionContent := `---
id: 20251231114608-5678
type: WET
track: 2
iteration: 1
status: active
---

## INTENT
Test mission`

	err := afero.WriteFile(fs, missionPath, []byte(missionContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test mission file: %v", err)
	}

	// Should get ID from mission.md
	id, err := service.GetCurrentID()
	if err != nil {
		t.Errorf("GetCurrentID() error = %v", err)
	}

	if id != "20251231114608-5678" {
		t.Errorf("Expected ID 20251231114608-5678, got %s", id)
	}
}

func TestIDService_isValidID(t *testing.T) {
	fs := afero.NewMemMapFs()
	service := NewIDService(fs, "")

	testCases := []struct {
		id    string
		valid bool
	}{
		{"20251231114608-1234", true},
		{"20251231114608-12", false},  // too short random part
		{"2025123111460-1234", false}, // too short timestamp
		{"20251231114608", false},     // missing random part
		{"invalid-id", false},         // invalid format
		{"", false},                   // empty
	}

	for _, tc := range testCases {
		t.Run(tc.id, func(t *testing.T) {
			result := service.isValidID(tc.id)
			if result != tc.valid {
				t.Errorf("isValidID(%s) = %v, want %v", tc.id, result, tc.valid)
			}
		})
	}
}

func TestIDService_generateID(t *testing.T) {
	fs := afero.NewMemMapFs()
	service := NewIDService(fs, "")

	id := service.generateID()

	// Check format
	if !service.isValidID(id) {
		t.Errorf("Generated ID %s is not valid", id)
	}
}
