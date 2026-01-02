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

	// Create required files for validation
	err = afero.WriteFile(fs, "/tmp/id", []byte("20251231131228-1234"), 0644)
	if err != nil {
		t.Fatalf("Failed to create ID file: %v", err)
	}

	err = afero.WriteFile(fs, "/tmp/execution.log", []byte("test log"), 0644)
	if err != nil {
		t.Fatalf("Failed to create execution log: %v", err)
	}

	status, err := service.CheckMissionState()
	if err != nil {
		t.Fatalf("CheckMissionState() error = %v", err)
	}

	if !status.HasActiveMission {
		t.Error("Expected active mission")
	}

	if status.Ready {
		t.Error("Expected not ready status")
	}
}

func TestValidationService_ClarificationState(t *testing.T) {
	fs := afero.NewMemMapFs()
	service := NewValidationService(fs, "/tmp")

	// Create mission in clarifying state
	missionContent := `# MISSION

id: 20251231131228-5678
type: WET
status: clarifying

## INTENT
Test clarification intent`

	// Create required files
	err := afero.WriteFile(fs, "/tmp/mission.md", []byte(missionContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test mission: %v", err)
	}

	err = afero.WriteFile(fs, "/tmp/id", []byte("20251231131228-5678"), 0644)
	if err != nil {
		t.Fatalf("Failed to create ID file: %v", err)
	}

	err = afero.WriteFile(fs, "/tmp/execution.log", []byte("test log"), 0644)
	if err != nil {
		t.Fatalf("Failed to create execution log: %v", err)
	}

	status, err := service.CheckMissionState()
	if err != nil {
		t.Fatalf("CheckMissionState() error = %v", err)
	}

	if status.MissionStatus != "clarifying" {
		t.Errorf("Expected clarifying status, got %s", status.MissionStatus)
	}

	if status.NextStep != "Run the m.clarify prompt to resolve questions." {
		t.Errorf("Unexpected next step for clarifying state: %s", status.NextStep)
	}

	// Verify no cleanup occurred (stale artifacts should be empty)
	if len(status.StaleArtifacts) > 0 {
		t.Errorf("Expected no cleanup in clarification mode, but got: %v", status.StaleArtifacts)
	}
}

func TestValidationService_CheckMissionState_WithClarifyingMission(t *testing.T) {
	fs := afero.NewMemMapFs()
	service := NewValidationService(fs, "/tmp")

	// Create existing mission in clarifying state
	missionContent := `# MISSION

id: 20251231131228-1234
type: WET
status: clarifying`

	err := afero.WriteFile(fs, "/tmp/mission.md", []byte(missionContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test mission: %v", err)
	}

	// Create required files for validation
	err = afero.WriteFile(fs, "/tmp/id", []byte("20251231131228-1234"), 0644)
	if err != nil {
		t.Fatalf("Failed to create ID file: %v", err)
	}

	err = afero.WriteFile(fs, "/tmp/execution.log", []byte("test log"), 0644)
	if err != nil {
		t.Fatalf("Failed to create execution log: %v", err)
	}

	status, err := service.CheckMissionState()
	if err != nil {
		t.Fatalf("CheckMissionState() error = %v", err)
	}

	if !status.HasActiveMission {
		t.Error("Expected active mission to be detected")
	}

	if status.MissionStatus != "clarifying" {
		t.Errorf("Expected status 'clarifying', got '%s'", status.MissionStatus)
	}

	expectedNextStep := "Run the m.clarify prompt to resolve questions."
	if status.NextStep != expectedNextStep {
		t.Errorf("Expected next step '%s', got '%s'", expectedNextStep, status.NextStep)
	}
}

func TestValidationService_ClarificationModeSkipsCleanup(t *testing.T) {
	fs := afero.NewMemMapFs()
	service := NewValidationService(fs, "/tmp")

	// Create stale artifacts that would normally be cleaned up
	err := afero.WriteFile(fs, "/tmp/id", []byte("20251231131228-5678"), 0644)
	if err != nil {
		t.Fatalf("Failed to create ID file: %v", err)
	}

	err = afero.WriteFile(fs, "/tmp/execution.log", []byte("existing log"), 0644)
	if err != nil {
		t.Fatalf("Failed to create execution log: %v", err)
	}

	err = afero.WriteFile(fs, "/tmp/plan.json", []byte(`{"intent":"test"}`), 0644)
	if err != nil {
		t.Fatalf("Failed to create plan.json: %v", err)
	}

	// Create mission in clarifying state
	missionContent := `# MISSION

id: 20251231131228-5678
type: WET
status: clarifying

## INTENT
Test clarification intent`

	err = afero.WriteFile(fs, "/tmp/mission.md", []byte(missionContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test mission: %v", err)
	}

	status, err := service.CheckMissionState()
	if err != nil {
		t.Fatalf("CheckMissionState() error = %v", err)
	}

	// Verify no cleanup occurred in clarification mode
	if len(status.StaleArtifacts) > 0 {
		t.Errorf("Expected no cleanup in clarification mode, but got: %v", status.StaleArtifacts)
	}

	// Verify files still exist
	if exists, _ := afero.Exists(fs, "/tmp/id"); !exists {
		t.Error("Expected ID file to be preserved in clarification mode")
	}

	if exists, _ := afero.Exists(fs, "/tmp/execution.log"); !exists {
		t.Error("Expected execution log to be preserved in clarification mode")
	}

	if exists, _ := afero.Exists(fs, "/tmp/plan.json"); !exists {
		t.Error("Expected plan.json to be preserved in clarification mode")
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
