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
parent_mission: 2025-12-17-15-30-mission.md

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
	if mission.ParentMission != "2025-12-17-15-30-mission.md" {
		t.Errorf("Expected ParentMission '2025-12-17-15-30-mission.md', got '%s'", mission.ParentMission)
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

	// Test new timestamp format
	mission3Content := `# MISSION
id: 20251229174130-1221
type: WET
track: 2
status: completed
completed_at: 2025-12-29T17:49:26.790-05:00

## INTENT
New timestamp format mission`

	// Write mission files
	mission1Path := filepath.Join(completedDir, "2025-12-18-10-15-mission.md")
	mission2Path := filepath.Join(completedDir, "2025-12-18-10-20-mission.md")
	mission3Path := filepath.Join(completedDir, "20251229174130-1221-mission.md")

	if err := os.WriteFile(mission1Path, []byte(mission1Content), 0644); err != nil {
		t.Fatalf("Failed to write mission1: %v", err)
	}
	if err := os.WriteFile(mission2Path, []byte(mission2Content), 0644); err != nil {
		t.Fatalf("Failed to write mission2: %v", err)
	}
	if err := os.WriteFile(mission3Path, []byte(mission3Content), 0644); err != nil {
		t.Fatalf("Failed to write mission3: %v", err)
	}

	// Test reading completed missions
	missions, err := ReadCompletedMissions()
	if err != nil {
		t.Fatalf("Failed to read completed missions: %v", err)
	}

	if len(missions) != 3 {
		t.Errorf("Expected 3 missions, got %d", len(missions))
	}

	// Verify new timestamp format is loaded
	found := false
	for _, mission := range missions {
		if mission.Intent == "New timestamp format mission" {
			found = true
			break
		}
	}
	if !found {
		t.Error("New timestamp format mission not found in loaded missions")
	}

	// Verify sorting (newest first)
	if len(missions) >= 3 {
		if missions[0].Intent != "New timestamp format mission" {
			t.Errorf("Expected newest mission first, got '%s'", missions[0].Intent)
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

func TestGetEffectiveTime(t *testing.T) {
	// Test with CompletedAt available
	completedTime, _ := time.Parse(time.RFC3339, "2025-12-18T10:20:00.000-05:00")
	mission1 := &Mission{
		CompletedAt: &completedTime,
		FilePath:    ".mission/completed/2025-12-18-10-20-mission.md",
	}

	effectiveTime := getEffectiveTime(mission1)
	if effectiveTime == nil || !effectiveTime.Equal(completedTime) {
		t.Errorf("Expected CompletedAt time, got %v", effectiveTime)
	}

	// Test with old filename format fallback (YYYY-MM-DD-HH-MM)
	mission2 := &Mission{
		CompletedAt: nil,
		FilePath:    ".mission/completed/2025-12-20-22-31-mission.md",
	}

	effectiveTime = getEffectiveTime(mission2)
	expectedTime, _ := time.Parse("2006-01-02-15-04", "2025-12-20-22-31")
	if effectiveTime == nil || !effectiveTime.Equal(expectedTime) {
		t.Errorf("Expected old format filename time %v, got %v", expectedTime, effectiveTime)
	}

	// Test with new filename format fallback (YYYYMMDDHHMMSS-SSSS)
	mission3 := &Mission{
		CompletedAt: nil,
		FilePath:    ".mission/completed/20251220223145-1234-mission.md",
	}

	effectiveTime = getEffectiveTime(mission3)
	expectedTime, _ = time.Parse("20060102150405", "20251220223145")
	if effectiveTime == nil || !effectiveTime.Equal(expectedTime) {
		t.Errorf("Expected new format filename time %v, got %v", expectedTime, effectiveTime)
	}

	// Test with invalid filename
	mission4 := &Mission{
		CompletedAt: nil,
		FilePath:    ".mission/completed/invalid-filename.md",
	}

	effectiveTime = getEffectiveTime(mission4)
	if effectiveTime != nil {
		t.Errorf("Expected nil for invalid filename, got %v", effectiveTime)
	}
}

func TestSortingWithMixedTimestamps(t *testing.T) {
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

	// Create missions with mixed timestamp scenarios
	// Mission with CompletedAt (newest)
	mission1Content := `# MISSION
type: WET
track: 2
status: completed
completed_at: 2025-12-25T10:00:00.000-05:00

## INTENT
Mission with CompletedAt`

	// Mission with old filename format (middle)
	mission2Content := `# MISSION
type: WET
track: 2
status: completed

## INTENT
Mission with old filename format`

	// Mission with new filename format (second newest)
	mission3Content := `# MISSION
type: WET
track: 2
status: completed

## INTENT
Mission with new filename format`

	// Mission with invalid filename format (should go to end)
	mission4Content := `# MISSION
type: WET
track: 2
status: completed

## INTENT
Mission with invalid filename`

	// Write mission files
	mission1Path := filepath.Join(completedDir, "2025-12-25-10-00-mission.md")
	mission2Path := filepath.Join(completedDir, "2025-12-20-22-31-mission.md")
	mission3Path := filepath.Join(completedDir, "20251224143045-2847-mission.md")
	mission4Path := filepath.Join(completedDir, "invalid-format-mission.md")

	if err := os.WriteFile(mission1Path, []byte(mission1Content), 0644); err != nil {
		t.Fatalf("Failed to write mission1: %v", err)
	}
	if err := os.WriteFile(mission2Path, []byte(mission2Content), 0644); err != nil {
		t.Fatalf("Failed to write mission2: %v", err)
	}
	if err := os.WriteFile(mission3Path, []byte(mission3Content), 0644); err != nil {
		t.Fatalf("Failed to write mission3: %v", err)
	}
	if err := os.WriteFile(mission4Path, []byte(mission4Content), 0644); err != nil {
		t.Fatalf("Failed to write mission4: %v", err)
	}

	// Test reading and sorting
	missions, err := ReadCompletedMissions()
	if err != nil {
		t.Fatalf("Failed to read completed missions: %v", err)
	}

	if len(missions) != 4 {
		t.Errorf("Expected 4 missions, got %d", len(missions))
	}

	// Verify sorting order: CompletedAt first, new format second, old format third, invalid last
	if len(missions) >= 4 {
		if missions[0].Intent != "Mission with CompletedAt" {
			t.Errorf("Expected CompletedAt mission first, got '%s'", missions[0].Intent)
		}
		if missions[1].Intent != "Mission with new filename format" {
			t.Errorf("Expected new format mission second, got '%s'", missions[1].Intent)
		}
		if missions[2].Intent != "Mission with old filename format" {
			t.Errorf("Expected old format mission third, got '%s'", missions[2].Intent)
		}
		if missions[3].Intent != "Mission with invalid filename" {
			t.Errorf("Expected invalid filename mission last, got '%s'", missions[3].Intent)
		}
	}
}
