package mission

import (
	"testing"

	"github.com/spf13/afero"
)

func TestValidationService_CheckMissionState_NoExistingMission(t *testing.T) {
	fs := afero.NewMemMapFs()
	service := NewValidationService(fs, "/tmp")

	status, err := service.CheckMissionState()
	if err != nil {
		t.Fatalf("CheckMissionState() error = %v", err)
	}

	if status.HasActiveMission {
		t.Error("Expected no active mission")
	}

	if !status.Ready {
		t.Error("Expected ready status")
	}

	if status.MissionID == "" {
		t.Error("Expected mission ID to be generated")
	}
}

func TestValidationService_CheckMissionState_WithExistingMission(t *testing.T) {
	fs := afero.NewMemMapFs()
	service := NewValidationService(fs, "/tmp")

	// Create existing mission
	missionContent := `# MISSION

id: 20251231131228-1234
type: WET
status: active`

	err := afero.WriteFile(fs, "/tmp/mission.md", []byte(missionContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test mission: %v", err)
	}

	status, err := service.CheckMissionState()
	if err != nil {
		t.Fatalf("CheckMissionState() error = %v", err)
	}

	if !status.HasActiveMission {
		t.Error("Expected active mission to be detected")
	}

	if status.Ready {
		t.Error("Expected not ready when active mission exists")
	}

	if status.MissionStatus != "active" {
		t.Errorf("Expected status 'active', got '%s'", status.MissionStatus)
	}

	if status.MissionID != "20251231131228-1234" {
		t.Errorf("Expected mission ID '20251231131228-1234', got '%s'", status.MissionID)
	}
}

func TestValidationService_CleanupStaleArtifacts(t *testing.T) {
	fs := afero.NewMemMapFs()
	service := NewValidationService(fs, "/tmp")

	// Create stale artifacts
	err := afero.WriteFile(fs, "/tmp/id", []byte("old-id"), 0644)
	if err != nil {
		t.Fatalf("Failed to create stale id: %v", err)
	}

	err = afero.WriteFile(fs, "/tmp/execution.log", []byte("old log"), 0644)
	if err != nil {
		t.Fatalf("Failed to create stale log: %v", err)
	}

	err = afero.WriteFile(fs, "/tmp/plan.json", []byte("{}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create stale plan.json: %v", err)
	}

	status, err := service.CheckMissionState()
	if err != nil {
		t.Fatalf("CheckMissionState() error = %v", err)
	}

	// Check artifacts were cleaned
	if len(status.StaleArtifacts) != 3 {
		t.Errorf("Expected 3 stale artifacts cleaned, got %d", len(status.StaleArtifacts))
	}

	// Note: Files are cleaned up and new ID is generated, so we check the cleanup was recorded
	expectedArtifacts := map[string]bool{"id": false, "execution.log": false, "plan.json": false}
	for _, artifact := range status.StaleArtifacts {
		expectedArtifacts[artifact] = true
	}

	for artifact, found := range expectedArtifacts {
		if !found {
			t.Errorf("Expected artifact '%s' to be in cleanup list", artifact)
		}
	}
}

func TestMissionStatus_ToJSON(t *testing.T) {
	status := &MissionStatus{
		HasActiveMission: true,
		MissionStatus:    "active",
		MissionID:        "test-id",
		Ready:            false,
		Message:          "test message",
	}

	json, err := status.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() error = %v", err)
	}

	if json == "" {
		t.Error("Expected non-empty JSON")
	}
}
