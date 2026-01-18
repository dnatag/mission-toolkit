package mission

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestArchiver_Archive(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	completedDir := filepath.Join(missionDir, "completed")

	// Create dummy mission files
	missionID := "test-mission-123"
	missionContent := `---
id: ` + missionID + `
---
Body
`
	err := fs.MkdirAll(missionDir, 0755)
	require.NoError(t, err)
	err = afero.WriteFile(fs, filepath.Join(missionDir, "mission.md"), []byte(missionContent), 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, filepath.Join(missionDir, "execution.log"), []byte("log content"), 0644)
	require.NoError(t, err)

	// Mock GitClient
	mockGit := &MockGitClient{
		commitMessage: "feat: test commit",
	}

	// Archive
	archiver := NewArchiver(fs, missionDir, mockGit)
	err = archiver.Archive(false)
	require.NoError(t, err)

	// Verify files were copied
	archivedMissionPath := filepath.Join(completedDir, missionID+"-mission.md")
	exists, err := afero.Exists(fs, archivedMissionPath)
	require.NoError(t, err)
	require.True(t, exists, "mission.md should be archived")

	archivedLogPath := filepath.Join(completedDir, missionID+"-execution.log")
	exists, err = afero.Exists(fs, archivedLogPath)
	require.NoError(t, err)
	require.True(t, exists, "execution.log should be archived")

	// Verify commit message was archived
	archivedCommitPath := filepath.Join(completedDir, missionID+"-commit.msg")
	exists, err = afero.Exists(fs, archivedCommitPath)
	require.NoError(t, err)
	require.True(t, exists, "commit.msg should be archived")
	content, err := afero.ReadFile(fs, archivedCommitPath)
	require.NoError(t, err)
	require.Equal(t, "feat: test commit", string(content))
}

func TestArchiver_Archive_MissingOptionalFiles(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	completedDir := filepath.Join(missionDir, "completed")

	// Create only essential mission file
	missionID := "test-mission-456"
	missionContent := `---
id: ` + missionID + `
---
Body
`
	err := fs.MkdirAll(missionDir, 0755)
	require.NoError(t, err)
	err = afero.WriteFile(fs, filepath.Join(missionDir, "mission.md"), []byte(missionContent), 0644)
	require.NoError(t, err)

	// Mock GitClient
	mockGit := &MockGitClient{
		commitMessage: "feat: test commit",
	}

	// Archive
	archiver := NewArchiver(fs, missionDir, mockGit)
	err = archiver.Archive(false)
	require.NoError(t, err)

	// Verify essential file was copied
	archivedMissionPath := filepath.Join(completedDir, missionID+"-mission.md")
	exists, err := afero.Exists(fs, archivedMissionPath)
	require.NoError(t, err)
	require.True(t, exists, "mission.md should be archived")

	// Verify optional files were not created
	archivedLogPath := filepath.Join(completedDir, missionID+"-execution.log")
	exists, err = afero.Exists(fs, archivedLogPath)
	require.NoError(t, err)
	require.False(t, exists, "execution.log should not be archived if it doesn't exist")
}

func TestArchiver_CleanupObsoleteFiles(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"

	// Create obsolete files that should be cleaned up
	err := fs.MkdirAll(missionDir, 0755)
	require.NoError(t, err)

	filesToCreate := []string{"execution.log", "mission.md", "id", "plan.json"}
	for _, filename := range filesToCreate {
		err = afero.WriteFile(fs, filepath.Join(missionDir, filename), []byte("content"), 0644)
		require.NoError(t, err)
	}

	// Mock GitClient
	mockGit := &MockGitClient{}

	// Execute cleanup
	archiver := NewArchiver(fs, missionDir, mockGit)
	err = archiver.CleanupObsoleteFiles()
	require.NoError(t, err)

	// Verify all obsolete files were removed
	for _, filename := range filesToCreate {
		exists, err := afero.Exists(fs, filepath.Join(missionDir, filename))
		require.NoError(t, err)
		require.False(t, exists, "%s should be removed after cleanup", filename)
	}
}

func TestArchiver_CleanupObsoleteFiles_MissingFiles(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"

	// Create directory but no files
	err := fs.MkdirAll(missionDir, 0755)
	require.NoError(t, err)

	// Mock GitClient
	mockGit := &MockGitClient{}

	// Cleanup should not fail if files don't exist
	archiver := NewArchiver(fs, missionDir, mockGit)
	err = archiver.CleanupObsoleteFiles()
	require.NoError(t, err)
}

func TestArchiver_Archive_ForceWithNoMission(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"

	// Create directory but no mission file
	err := fs.MkdirAll(missionDir, 0755)
	require.NoError(t, err)

	// Mock GitClient
	mockGit := &MockGitClient{}

	// Archive with force=true should succeed as no-op when no mission exists
	archiver := NewArchiver(fs, missionDir, mockGit)
	err = archiver.Archive(true)
	require.NoError(t, err, "Archive with force=true should succeed when no mission exists")
}

func TestArchiver_Archive_NoForceWithNoMission(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"

	// Create directory but no mission file
	err := fs.MkdirAll(missionDir, 0755)
	require.NoError(t, err)

	// Mock GitClient
	mockGit := &MockGitClient{}

	// Archive with force=false should return error when no mission exists
	archiver := NewArchiver(fs, missionDir, mockGit)
	err = archiver.Archive(false)
	require.Error(t, err, "Archive with force=false should fail when no mission exists")
	require.Contains(t, err.Error(), "no current mission to archive")
}

// Edge case: Archive with corrupted mission file
func TestArchiver_Archive_CorruptedMissionFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"

	err := fs.MkdirAll(missionDir, 0755)
	require.NoError(t, err)

	// Write invalid YAML frontmatter
	err = afero.WriteFile(fs, filepath.Join(missionDir, "mission.md"), []byte("---\ninvalid: yaml: content:\n---\nBody"), 0644)
	require.NoError(t, err)

	mockGit := &MockGitClient{}
	archiver := NewArchiver(fs, missionDir, mockGit)

	err = archiver.Archive(false)
	require.Error(t, err, "Should fail with corrupted mission file")
}

// Edge case: Archive with read-only completed directory
func TestArchiver_Archive_ReadOnlyDestination(t *testing.T) {
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

	// Create read-only completed directory
	completedDir := filepath.Join(missionDir, "completed")
	err = fs.MkdirAll(completedDir, 0444)
	require.NoError(t, err)

	mockGit := &MockGitClient{}
	archiver := NewArchiver(fs, missionDir, mockGit)

	err = archiver.Archive(false)
	// Should handle permission errors gracefully
	if err != nil {
		require.Contains(t, err.Error(), "permission denied")
	}
}
