package mission

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestReadMissionFile(t *testing.T) {
	// Create a temporary mission file content
	missionContent := `# MISSION

type: WET
track: 2
iteration: 1
status: active
completed_at: 2025-12-18T10:20:00.000-05:00

## INTENT
Test mission for parsing

## SCOPE
file1.go
file2.go

## PLAN
- [ ] Step 1
- [x] Step 2 completed
- [ ] Step 3

## VERIFICATION
go test ./...`

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "test-mission-*.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(missionContent); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	// Test reading the mission file
	mission, err := readMissionFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read mission file: %v", err)
	}

	// Verify parsed content
	if mission.Type != "WET" {
		t.Errorf("Expected Type 'WET', got '%s'", mission.Type)
	}
	if mission.Track != "2" {
		t.Errorf("Expected Track '2', got '%s'", mission.Track)
	}
	if mission.Status != "active" {
		t.Errorf("Expected Status 'active', got '%s'", mission.Status)
	}
	if mission.Intent != "Test mission for parsing" {
		t.Errorf("Expected Intent 'Test mission for parsing', got '%s'", mission.Intent)
	}
	if len(mission.Scope) != 2 {
		t.Errorf("Expected 2 scope items, got %d", len(mission.Scope))
	}
	if len(mission.Plan) != 3 {
		t.Errorf("Expected 3 plan items, got %d", len(mission.Plan))
	}
	if mission.Verification != "go test ./..." {
		t.Errorf("Expected Verification 'go test ./...', got '%s'", mission.Verification)
	}
	if mission.CompletedAt == nil {
		t.Error("Expected CompletedAt to be parsed")
	} else {
		expectedTime, _ := time.Parse(time.RFC3339, "2025-12-18T10:20:00.000-05:00")
		if !mission.CompletedAt.Equal(expectedTime) {
			t.Errorf("Expected CompletedAt %v, got %v", expectedTime, *mission.CompletedAt)
		}
	}
}

func TestReadCompletedMissions(t *testing.T) {
	// Create temporary directory structure
	tmpDir, err := os.MkdirTemp("", "test-mission-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	// Create .mission/completed directory
	completedDir := ".mission/completed"
	if err := os.MkdirAll(completedDir, 0755); err != nil {
		t.Fatalf("Failed to create completed dir: %v", err)
	}

	// Create test mission files
	mission1Content := `# MISSION
type: WET
track: 2
status: completed
completed_at: 2025-12-18T10:15:00.000-05:00

## INTENT
First test mission`

	mission2Content := `# MISSION
type: DRY
track: 3
status: completed
completed_at: 2025-12-18T10:20:00.000-05:00

## INTENT
Second test mission`

	// Write mission files
	mission1Path := filepath.Join(completedDir, "2025-12-18-10-15-mission.md")
	mission2Path := filepath.Join(completedDir, "2025-12-18-10-20-mission.md")
	
	if err := os.WriteFile(mission1Path, []byte(mission1Content), 0644); err != nil {
		t.Fatalf("Failed to write mission1: %v", err)
	}
	if err := os.WriteFile(mission2Path, []byte(mission2Content), 0644); err != nil {
		t.Fatalf("Failed to write mission2: %v", err)
	}

	// Test reading completed missions
	missions, err := ReadCompletedMissions()
	if err != nil {
		t.Fatalf("Failed to read completed missions: %v", err)
	}

	if len(missions) != 2 {
		t.Errorf("Expected 2 missions, got %d", len(missions))
	}

	// Verify sorting (newest first)
	if len(missions) >= 2 {
		if missions[0].Intent != "Second test mission" {
			t.Errorf("Expected newest mission first, got '%s'", missions[0].Intent)
		}
		if missions[1].Intent != "First test mission" {
			t.Errorf("Expected oldest mission second, got '%s'", missions[1].Intent)
		}
	}
}

func TestReadCurrentMission(t *testing.T) {
	// Create temporary directory structure
	tmpDir, err := os.MkdirTemp("", "test-mission-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	// Create .mission directory
	if err := os.MkdirAll(".mission", 0755); err != nil {
		t.Fatalf("Failed to create .mission dir: %v", err)
	}

	// Create current mission file
	missionContent := `# MISSION
type: WET
track: 2
status: active

## INTENT
Current test mission`

	missionPath := ".mission/mission.md"
	if err := os.WriteFile(missionPath, []byte(missionContent), 0644); err != nil {
		t.Fatalf("Failed to write current mission: %v", err)
	}

	// Test reading current mission
	mission, err := ReadCurrentMission()
	if err != nil {
		t.Fatalf("Failed to read current mission: %v", err)
	}

	if mission.Type != "WET" {
		t.Errorf("Expected Type 'WET', got '%s'", mission.Type)
	}
	if mission.Status != "active" {
		t.Errorf("Expected Status 'active', got '%s'", mission.Status)
	}
	if mission.Intent != "Current test mission" {
		t.Errorf("Expected Intent 'Current test mission', got '%s'", mission.Intent)
	}
}
