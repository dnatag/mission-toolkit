package mission

import (
	"testing"

	"github.com/spf13/afero"
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

func TestCheckService_WithCommand_Complete_ActiveStatus(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	fs.MkdirAll(missionDir, 0755)

	missionContent := `---
id: test-103
status: active
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

	if status.NextStep != "STOP. Mission must be in 'active' or 'completed' status for m.complete." {
		t.Errorf("CheckMissionState() NextStep = %v, want STOP", status.NextStep)
	}
}
