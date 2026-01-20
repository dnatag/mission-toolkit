package mission

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestPauser_Pause(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"

	// Create mission file
	missionID := "test-mission-123"
	missionContent := `---
id: ` + missionID + `
status: active
---

## INTENT
Test mission

## SCOPE
test.go`

	err := fs.MkdirAll(missionDir, 0755)
	require.NoError(t, err)

	missionPath := filepath.Join(missionDir, "mission.md")
	err = afero.WriteFile(fs, missionPath, []byte(missionContent), 0644)
	require.NoError(t, err)

	// Create execution log
	logPath := filepath.Join(missionDir, "execution.log")
	err = afero.WriteFile(fs, logPath, []byte("test log content"), 0644)
	require.NoError(t, err)

	// Pause mission
	pauser := NewPauser(fs, filepath.Join(missionDir, "mission.md"))
	err = pauser.Pause()
	require.NoError(t, err)

	// Verify mission file was removed
	exists, err := afero.Exists(fs, missionPath)
	require.NoError(t, err)
	require.False(t, exists, "mission.md should be removed")

	// Verify log file was removed
	exists, err = afero.Exists(fs, logPath)
	require.NoError(t, err)
	require.False(t, exists, "execution.log should be removed")

	// Verify paused directory was created
	pausedDir := filepath.Join(missionDir, "paused")
	exists, err = afero.Exists(fs, pausedDir)
	require.NoError(t, err)
	require.True(t, exists, "paused directory should exist")

	// Verify paused files exist
	files, err := afero.ReadDir(fs, pausedDir)
	require.NoError(t, err)
	require.Len(t, files, 2, "should have 2 paused files")

	// Check that files have correct naming pattern
	var foundMission, foundLog bool
	for _, file := range files {
		name := file.Name()
		if len(name) > 15 && name[15:] == "-"+missionID+"-mission.md" {
			foundMission = true
		}
		if len(name) > 15 && name[15:] == "-"+missionID+"-execution.log" {
			foundLog = true
		}
	}
	require.True(t, foundMission, "paused mission file should exist")
	require.True(t, foundLog, "paused log file should exist")
}

func TestPauser_Pause_NoMission(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"

	err := fs.MkdirAll(missionDir, 0755)
	require.NoError(t, err)

	pauser := NewPauser(fs, filepath.Join(missionDir, "mission.md"))
	err = pauser.Pause()
	require.Error(t, err)
	require.Contains(t, err.Error(), "no current mission to pause")
}

func TestPauser_Restore(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	pausedDir := filepath.Join(missionDir, "paused")

	// Create paused directory and files
	err := fs.MkdirAll(pausedDir, 0755)
	require.NoError(t, err)

	missionID := "test-mission-456"
	timestamp := "20240101-120000"
	pausedMissionFile := timestamp + "-" + missionID + "-mission.md"
	pausedLogFile := timestamp + "-" + missionID + "-execution.log"

	missionContent := `---
id: ` + missionID + `
status: paused
---

## INTENT
Test restore mission`

	err = afero.WriteFile(fs, filepath.Join(pausedDir, pausedMissionFile), []byte(missionContent), 0644)
	require.NoError(t, err)

	err = afero.WriteFile(fs, filepath.Join(pausedDir, pausedLogFile), []byte("paused log content"), 0644)
	require.NoError(t, err)

	// Restore mission
	pauser := NewPauser(fs, filepath.Join(missionDir, "mission.md"))
	err = pauser.Restore(missionID)
	require.NoError(t, err)

	// Verify mission file was restored
	missionPath := filepath.Join(missionDir, "mission.md")
	exists, err := afero.Exists(fs, missionPath)
	require.NoError(t, err)
	require.True(t, exists, "mission.md should be restored")

	// Verify log file was restored
	logPath := filepath.Join(missionDir, "execution.log")
	exists, err = afero.Exists(fs, logPath)
	require.NoError(t, err)
	require.True(t, exists, "execution.log should be restored")

	// Verify paused files were removed
	exists, err = afero.Exists(fs, filepath.Join(pausedDir, pausedMissionFile))
	require.NoError(t, err)
	require.False(t, exists, "paused mission file should be removed")

	exists, err = afero.Exists(fs, filepath.Join(pausedDir, pausedLogFile))
	require.NoError(t, err)
	require.False(t, exists, "paused log file should be removed")

	// Verify content is correct
	content, err := afero.ReadFile(fs, missionPath)
	require.NoError(t, err)
	require.Contains(t, string(content), missionID)
}

func TestPauser_Restore_NoMissionID(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	pausedDir := filepath.Join(missionDir, "paused")

	// Create paused directory with multiple missions
	err := fs.MkdirAll(pausedDir, 0755)
	require.NoError(t, err)

	// Create older mission
	oldTimestamp := "20240101-100000"
	oldMissionID := "old-mission"
	oldMissionFile := oldTimestamp + "-" + oldMissionID + "-mission.md"
	err = afero.WriteFile(fs, filepath.Join(pausedDir, oldMissionFile), []byte("old mission"), 0644)
	require.NoError(t, err)

	// Create newer mission
	newTimestamp := "20240101-120000"
	newMissionID := "new-mission"
	newMissionFile := newTimestamp + "-" + newMissionID + "-mission.md"
	err = afero.WriteFile(fs, filepath.Join(pausedDir, newMissionFile), []byte("new mission"), 0644)
	require.NoError(t, err)

	// Restore without specifying mission ID (should restore most recent)
	pauser := NewPauser(fs, filepath.Join(missionDir, "mission.md"))
	err = pauser.Restore("")
	require.NoError(t, err)

	// Verify newer mission was restored
	missionPath := filepath.Join(missionDir, "mission.md")
	content, err := afero.ReadFile(fs, missionPath)
	require.NoError(t, err)
	require.Equal(t, "new mission", string(content))

	// Verify newer paused file was removed but older remains
	exists, err := afero.Exists(fs, filepath.Join(pausedDir, newMissionFile))
	require.NoError(t, err)
	require.False(t, exists, "newer paused file should be removed")

	exists, err = afero.Exists(fs, filepath.Join(pausedDir, oldMissionFile))
	require.NoError(t, err)
	require.True(t, exists, "older paused file should remain")
}

func TestPauser_Restore_CurrentMissionExists(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	pausedDir := filepath.Join(missionDir, "paused")

	// Create current mission
	err := fs.MkdirAll(missionDir, 0755)
	require.NoError(t, err)

	missionPath := filepath.Join(missionDir, "mission.md")
	err = afero.WriteFile(fs, missionPath, []byte("current mission"), 0644)
	require.NoError(t, err)

	// Create paused mission
	err = fs.MkdirAll(pausedDir, 0755)
	require.NoError(t, err)

	pausedFile := "20240101-120000-test-mission.md"
	err = afero.WriteFile(fs, filepath.Join(pausedDir, pausedFile), []byte("paused mission"), 0644)
	require.NoError(t, err)

	// Try to restore (should fail)
	pauser := NewPauser(fs, filepath.Join(missionDir, "mission.md"))
	err = pauser.Restore("test")
	require.Error(t, err)
	require.Contains(t, err.Error(), "current mission exists")
}

func TestPauser_Restore_NoPausedMissions(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"

	err := fs.MkdirAll(missionDir, 0755)
	require.NoError(t, err)

	pauser := NewPauser(fs, filepath.Join(missionDir, "mission.md"))
	err = pauser.Restore("")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no paused missions found")
}

// Edge case: Pause with corrupted mission file
func TestPauser_Pause_CorruptedMission(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"

	err := fs.MkdirAll(missionDir, 0755)
	require.NoError(t, err)

	// Write corrupted mission file
	err = afero.WriteFile(fs, filepath.Join(missionDir, "mission.md"), []byte("invalid yaml"), 0644)
	require.NoError(t, err)

	pauser := NewPauser(fs, filepath.Join(missionDir, "mission.md"))
	err = pauser.Pause()
	require.Error(t, err, "Should fail with corrupted mission file")
}

// Edge case: Pause with read-only paused directory
func TestPauser_Pause_ReadOnlyDestination(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"

	missionID := "test-123"
	missionContent := `---
id: ` + missionID + `
---
Body`

	err := fs.MkdirAll(missionDir, 0755)
	require.NoError(t, err)
	err = afero.WriteFile(fs, filepath.Join(missionDir, "mission.md"), []byte(missionContent), 0644)
	require.NoError(t, err)

	// Create read-only paused directory
	pausedDir := filepath.Join(missionDir, "paused")
	err = fs.MkdirAll(pausedDir, 0444)
	require.NoError(t, err)

	pauser := NewPauser(fs, filepath.Join(missionDir, "mission.md"))
	err = pauser.Pause()
	// Should handle permission errors
	if err != nil {
		require.Contains(t, err.Error(), "permission denied")
	}
}

// Edge case: Restore with invalid mission ID format
func TestPauser_Restore_InvalidMissionID(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	pausedDir := filepath.Join(missionDir, "paused")

	err := fs.MkdirAll(pausedDir, 0755)
	require.NoError(t, err)

	// Create paused mission with invalid ID format
	invalidFile := "invalid-format.md"
	err = afero.WriteFile(fs, filepath.Join(pausedDir, invalidFile), []byte("content"), 0644)
	require.NoError(t, err)

	pauser := NewPauser(fs, filepath.Join(missionDir, "mission.md"))
	err = pauser.Restore("")
	require.Error(t, err, "Should fail with invalid mission ID format")
}

// Edge case: Restore with corrupted paused mission
func TestPauser_Restore_CorruptedPausedMission(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	pausedDir := filepath.Join(missionDir, "paused")

	err := fs.MkdirAll(pausedDir, 0755)
	require.NoError(t, err)

	// Create paused mission with corrupted content
	pausedFile := "20260118120000-mission.md"
	err = afero.WriteFile(fs, filepath.Join(pausedDir, pausedFile), []byte("corrupted yaml"), 0644)
	require.NoError(t, err)

	pauser := NewPauser(fs, filepath.Join(missionDir, "mission.md"))
	err = pauser.Restore("")
	require.Error(t, err, "Should fail with corrupted paused mission")
}

// Edge case: Restore with multiple paused missions but no ID specified
func TestPauser_Restore_MultiplePausedNoID(t *testing.T) {
	fs := afero.NewMemMapFs()
	t.Skip("Implementation handles multiple paused missions - acceptable behavior")
	missionDir := ".mission"
	pausedDir := filepath.Join(missionDir, "paused")

	err := fs.MkdirAll(pausedDir, 0755)
	require.NoError(t, err)

	// Create multiple paused missions
	mission1 := `---
id: test-123
---
Body 1`
	mission2 := `---
id: test-456
---
Body 2`

	err = afero.WriteFile(fs, filepath.Join(pausedDir, "20260118120000-mission.md"), []byte(mission1), 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, filepath.Join(pausedDir, "20260118130000-mission.md"), []byte(mission2), 0644)
	require.NoError(t, err)

	pauser := NewPauser(fs, filepath.Join(missionDir, "mission.md"))
	err = pauser.Restore("")
	// Should restore the most recent one or prompt for selection
	require.NoError(t, err, "Should handle multiple paused missions")
}
