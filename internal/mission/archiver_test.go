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
	err = archiver.Archive()
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
	err = archiver.Archive()
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
