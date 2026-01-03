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
		t.Errorf("CheckMissionState() MissionStatus = %v, want planned", status.MissionStatus)
	}
	if status.MissionIntent != "Test mission intent" {
		t.Errorf("CheckMissionState() MissionIntent = %v, want Test mission intent", status.MissionIntent)
	}
}
