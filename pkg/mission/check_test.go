package mission

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestCheckService_CheckMissionState_NoMission(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	fs.MkdirAll(missionDir, 0755)

	service := NewCheckService(fs, missionDir)
	status, err := service.CheckMissionState()
	if err != nil {
		t.Fatalf("CheckMissionState() error = %v", err)
	}

	if status.HasActiveMission {
		t.Error("CheckMissionState() should not have active mission")
	}
	if !status.Ready {
		t.Error("CheckMissionState() should be ready")
	}
	if status.MissionID == "" {
		t.Error("CheckMissionState() should generate mission ID")
	}
}

func TestCheckService_CheckMissionState_WithMission(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	fs.MkdirAll(missionDir, 0755)

	missionContent := `---
id: test-123
type: WET
track: 2
iteration: 1
status: planned
---

## INTENT
Test mission intent

## SCOPE
file1.go
`
	afero.WriteFile(fs, missionDir+"/mission.md", []byte(missionContent), 0644)

	service := NewCheckService(fs, missionDir)
	status, err := service.CheckMissionState()
	if err != nil {
		t.Fatalf("CheckMissionState() error = %v", err)
	}

	if !status.HasActiveMission {
		t.Error("CheckMissionState() should have active mission")
	}
	if status.MissionID != "test-123" {
		t.Errorf("CheckMissionState() MissionID = %v, want test-123", status.MissionID)
	}
	if status.MissionStatus != "planned" {
		t.Errorf("CheckMissionState() Status = %v, want planned", status.MissionStatus)
	}
	if status.MissionIntent != "Test mission intent" {
		t.Errorf("CheckMissionState() MissionIntent = %v, want Test mission intent", status.MissionIntent)
	}
	if status.NextStep != "Run the m.apply prompt to execute this mission." {
		t.Errorf("CheckMissionState() NextStep = %v, want m.apply instruction", status.NextStep)
	}
}

func TestCheckService_CheckMissionState_ActiveMission(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	fs.MkdirAll(missionDir, 0755)

	missionContent := `---
id: test-456
type: WET
track: 2
iteration: 1
status: active
---

## INTENT
Active mission intent

## SCOPE
file2.go
`
	afero.WriteFile(fs, missionDir+"/mission.md", []byte(missionContent), 0644)

	service := NewCheckService(fs, missionDir)
	status, err := service.CheckMissionState()
	if err != nil {
		t.Fatalf("CheckMissionState() error = %v", err)
	}

	if status.MissionStatus != "active" {
		t.Errorf("CheckMissionState() Status = %v, want active", status.MissionStatus)
	}
	if status.NextStep != "Run the m.apply prompt to execute this mission." {
		t.Errorf("CheckMissionState() NextStep = %v, want m.apply instruction", status.NextStep)
	}
}

func TestCheckService_CheckMissionState_CompletedMission(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	fs.MkdirAll(missionDir, 0755)

	missionContent := `---
id: test-789
type: WET
track: 2
iteration: 1
status: completed
---

## INTENT
Completed mission intent

## SCOPE
file3.go
`
	afero.WriteFile(fs, missionDir+"/mission.md", []byte(missionContent), 0644)

	service := NewCheckService(fs, missionDir)
	status, err := service.CheckMissionState()
	if err != nil {
		t.Fatalf("CheckMissionState() error = %v", err)
	}

	if status.MissionStatus != "completed" {
		t.Errorf("CheckMissionState() Status = %v, want completed", status.MissionStatus)
	}
	if status.NextStep != "Run the m.complete prompt to finalize this mission." {
		t.Errorf("CheckMissionState() NextStep = %v, want m.complete instruction", status.NextStep)
	}
}

func TestCheckService_WithCommand_Apply_PlannedStatus(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	fs.MkdirAll(missionDir, 0755)

	missionContent := `---
id: test-100
status: planned
---

## INTENT
Test intent
`
	afero.WriteFile(fs, missionDir+"/mission.md", []byte(missionContent), 0644)

	service := NewCheckService(fs, missionDir)
	service.SetContext("apply")
	status, err := service.CheckMissionState()
	if err != nil {
		t.Fatalf("CheckMissionState() error = %v", err)
	}

	if status.NextStep != "PROCEED with m.apply execution." {
		t.Errorf("CheckMissionState() NextStep = %v, want PROCEED", status.NextStep)
	}
}

func TestCheckService_WithCommand_Apply_ActiveStatus(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	fs.MkdirAll(missionDir, 0755)

	missionContent := `---
id: test-101
status: active
---

## INTENT
Test intent
`
	afero.WriteFile(fs, missionDir+"/mission.md", []byte(missionContent), 0644)

	service := NewCheckService(fs, missionDir)
	service.SetContext("apply")
	status, err := service.CheckMissionState()
	if err != nil {
		t.Fatalf("CheckMissionState() error = %v", err)
	}

	if status.NextStep != "PROCEED with m.apply execution." {
		t.Errorf("CheckMissionState() NextStep = %v, want PROCEED", status.NextStep)
	}
}

func TestCheckService_WithCommand_Apply_CompletedStatus(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	fs.MkdirAll(missionDir, 0755)

	missionContent := `---
id: test-102
status: completed
---

## INTENT
Test intent
`
	afero.WriteFile(fs, missionDir+"/mission.md", []byte(missionContent), 0644)

	service := NewCheckService(fs, missionDir)
	service.SetContext("apply")
	status, err := service.CheckMissionState()
	if err != nil {
		t.Fatalf("CheckMissionState() error = %v", err)
	}

	if status.NextStep != "STOP. Mission must be in 'planned', 'active', or 'failed' status for m.apply." {
		t.Errorf("CheckMissionState() NextStep = %v, want STOP", status.NextStep)
	}
}

func TestCheckService_WithCommand_Apply_FailedStatus(t *testing.T) {
	// Test that failed missions can be retried with m.apply
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	fs.MkdirAll(missionDir, 0755)

	missionContent := `---
id: test-123
status: failed
---

## INTENT
Test intent
`
	afero.WriteFile(fs, missionDir+"/mission.md", []byte(missionContent), 0644)

	service := NewCheckService(fs, missionDir)
	service.SetContext("apply")
	status, err := service.CheckMissionState()
	if err != nil {
		t.Fatalf("CheckMissionState() error = %v", err)
	}

	expectedNextStep := "PROCEED with m.apply execution."
	if status.NextStep != expectedNextStep {
		t.Errorf("CheckMissionState() NextStep = %v, want %v", status.NextStep, expectedNextStep)
	}

	expectedMessage := "Mission is ready for execution or re-execution"
	if status.Message != expectedMessage {
		t.Errorf("CheckMissionState() Message = %v, want %v", status.Message, expectedMessage)
	}
}

func TestCheckService_WithCommand_Complete_ExecutedStatus(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	fs.MkdirAll(missionDir, 0755)

	missionContent := `---
id: test-103
status: executed
---

## INTENT
Test intent
`
	afero.WriteFile(fs, missionDir+"/mission.md", []byte(missionContent), 0644)

	service := NewCheckService(fs, missionDir)
	service.SetContext("complete")
	status, err := service.CheckMissionState()
	if err != nil {
		t.Fatalf("CheckMissionState() error = %v", err)
	}

	if status.NextStep != "PROCEED with m.complete execution." {
		t.Errorf("CheckMissionState() NextStep = %v, want PROCEED", status.NextStep)
	}
}

func TestCheckService_WithCommand_Complete_CompletedStatus(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	fs.MkdirAll(missionDir, 0755)

	missionContent := `---
id: test-104
status: completed
---

## INTENT
Test intent
`
	afero.WriteFile(fs, missionDir+"/mission.md", []byte(missionContent), 0644)

	service := NewCheckService(fs, missionDir)
	service.SetContext("complete")
	status, err := service.CheckMissionState()
	if err != nil {
		t.Fatalf("CheckMissionState() error = %v", err)
	}

	if status.NextStep != "PROCEED with m.complete execution." {
		t.Errorf("CheckMissionState() NextStep = %v, want PROCEED", status.NextStep)
	}
}

func TestCheckService_WithCommand_Complete_PlannedStatus(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	fs.MkdirAll(missionDir, 0755)

	missionContent := `---
id: test-105
status: planned
---

## INTENT
Test intent
`
	afero.WriteFile(fs, missionDir+"/mission.md", []byte(missionContent), 0644)

	service := NewCheckService(fs, missionDir)
	service.SetContext("complete")
	status, err := service.CheckMissionState()
	if err != nil {
		t.Fatalf("CheckMissionState() error = %v", err)
	}

	if status.NextStep != "STOP. Mission must be in 'executed' or 'completed' status for m.complete." {
		t.Errorf("CheckMissionState() NextStep = %v, want STOP", status.NextStep)
	}
}

// TestCheckService_StaleArtifactCleanup_PlanContext tests that stale artifacts are cleaned up when called from m.plan context
func TestCheckService_StaleArtifactCleanup_PlanContext(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	fs.MkdirAll(missionDir, 0755)

	// Create stale artifacts
	afero.WriteFile(fs, missionDir+"/id", []byte("old-id"), 0644)
	afero.WriteFile(fs, missionDir+"/plan.json", []byte(`{"old": "plan"}`), 0644)
	afero.WriteFile(fs, missionDir+"/execution.log", []byte("old log"), 0644)

	service := NewCheckService(fs, missionDir)
	service.SetContext("plan")
	status, err := service.CheckMissionState()
	if err != nil {
		t.Fatalf("CheckMissionState() error = %v", err)
	}

	// Verify artifacts were cleaned up
	expectedArtifacts := []string{"id", "plan.json", "execution.log"}
	if len(status.StaleArtifacts) != len(expectedArtifacts) {
		t.Errorf("Expected %d stale artifacts cleaned, got %d", len(expectedArtifacts), len(status.StaleArtifacts))
	}

	for _, artifact := range expectedArtifacts {
		found := false
		for _, cleaned := range status.StaleArtifacts {
			if cleaned == artifact {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected artifact %s to be cleaned up", artifact)
		}
	}

	// Verify plan.json and execution.log were removed (id gets recreated by IDService)
	planPath := missionDir + "/plan.json"
	if exists, _ := afero.Exists(fs, planPath); exists {
		t.Error("plan.json should have been removed but still exists")
	}

	logPath := missionDir + "/execution.log"
	if exists, _ := afero.Exists(fs, logPath); exists {
		t.Error("execution.log should have been removed but still exists")
	}

	// Verify id file exists but contains new content (not "old-id")
	idPath := missionDir + "/id"
	if exists, _ := afero.Exists(fs, idPath); !exists {
		t.Error("id file should exist (recreated by IDService)")
	} else {
		content, _ := afero.ReadFile(fs, idPath)
		if string(content) == "old-id" {
			t.Error("id file should contain new ID, not old stale content")
		}
	}
}

// TestCheckService_StaleArtifactCleanup_NonPlanContext tests that stale artifacts are NOT cleaned up when called from other contexts
func TestCheckService_StaleArtifactCleanup_NonPlanContext(t *testing.T) {
	testCases := []struct {
		name    string
		context string
	}{
		{"apply context", "apply"},
		{"complete context", "complete"},
		{"no context", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			missionDir := ".mission"
			fs.MkdirAll(missionDir, 0755)

			// Create stale artifacts
			afero.WriteFile(fs, missionDir+"/id", []byte("old-id"), 0644)
			afero.WriteFile(fs, missionDir+"/plan.json", []byte(`{"old": "plan"}`), 0644)
			afero.WriteFile(fs, missionDir+"/execution.log", []byte("old log"), 0644)

			service := NewCheckService(fs, missionDir)
			service.SetContext(tc.context)
			status, err := service.CheckMissionState()
			if err != nil {
				t.Fatalf("CheckMissionState() error = %v", err)
			}

			// Verify no artifacts were cleaned up
			if len(status.StaleArtifacts) != 0 {
				t.Errorf("Expected no stale artifacts cleaned in %s context, got %d", tc.context, len(status.StaleArtifacts))
			}

			// Verify files still exist
			artifacts := []string{"id", "plan.json", "execution.log"}
			for _, artifact := range artifacts {
				path := missionDir + "/" + artifact
				if exists, _ := afero.Exists(fs, path); !exists {
					t.Errorf("Artifact %s should not have been removed in %s context", artifact, tc.context)
				}
			}
		})
	}
}

// TestCheckService_StaleArtifactCleanup_PartialArtifacts tests cleanup when only some stale artifacts exist
func TestCheckService_StaleArtifactCleanup_PartialArtifacts(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	fs.MkdirAll(missionDir, 0755)

	// Create only some stale artifacts
	afero.WriteFile(fs, missionDir+"/id", []byte("old-id"), 0644)
	afero.WriteFile(fs, missionDir+"/plan.json", []byte(`{"old": "plan"}`), 0644)
	// execution.log intentionally not created

	service := NewCheckService(fs, missionDir)
	service.SetContext("plan")
	status, err := service.CheckMissionState()
	if err != nil {
		t.Fatalf("CheckMissionState() error = %v", err)
	}

	// Verify only existing artifacts were cleaned up
	expectedArtifacts := []string{"id", "plan.json"}
	if len(status.StaleArtifacts) != len(expectedArtifacts) {
		t.Errorf("Expected %d stale artifacts cleaned, got %d", len(expectedArtifacts), len(status.StaleArtifacts))
	}

	for _, artifact := range expectedArtifacts {
		found := false
		for _, cleaned := range status.StaleArtifacts {
			if cleaned == artifact {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected artifact %s to be cleaned up", artifact)
		}
	}

	// Verify execution.log is not in cleaned artifacts (since it didn't exist)
	for _, cleaned := range status.StaleArtifacts {
		if cleaned == "execution.log" {
			t.Error("execution.log should not be in cleaned artifacts since it didn't exist")
		}
	}
}

// Edge case: Check with corrupted mission file
func TestCheckService_CheckMissionState_CorruptedMission(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"

	err := fs.MkdirAll(missionDir, 0755)
	require.NoError(t, err)

	// Write corrupted mission file
	err = afero.WriteFile(fs, filepath.Join(missionDir, "mission.md"), []byte("invalid yaml content"), 0644)
	require.NoError(t, err)

	service := NewCheckService(fs, missionDir)
	status, err := service.CheckMissionState()
	require.Error(t, err, "Should fail with corrupted mission file")
	require.Nil(t, status)
}

// Edge case: Check with missing mission directory
func TestCheckService_CheckMissionState_MissingDirectory(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"

	service := NewCheckService(fs, missionDir)
	status, err := service.CheckMissionState()
	require.NoError(t, err, "Should handle missing directory gracefully")
	require.NotNil(t, status)
	require.False(t, status.HasActiveMission)
}
