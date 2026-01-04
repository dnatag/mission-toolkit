package cmd

import (
	"path/filepath"
	"testing"

	"github.com/dnatag/mission-toolkit/internal/mission"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

// MockGitClient for testing
type MockGitClient struct {
	commitMessage string
}

func (m *MockGitClient) Add(files []string) error {
	return nil
}

func (m *MockGitClient) Commit(message string) (string, error) {
	return "mock-hash", nil
}

func (m *MockGitClient) CreateTag(name string, commitHash string) error {
	return nil
}

func (m *MockGitClient) Restore(checkpointName string, files []string) error {
	return nil
}

func (m *MockGitClient) ListTags(prefix string) ([]string, error) {
	return []string{}, nil
}

func (m *MockGitClient) DeleteTag(name string) error {
	return nil
}

func (m *MockGitClient) GetTagCommit(tagName string) (string, error) {
	return "mock-commit-hash", nil
}

func (m *MockGitClient) SoftReset(commitHash string) error {
	return nil
}

func (m *MockGitClient) GetCommitMessage(commitHash string) (string, error) {
	return m.commitMessage, nil
}

func TestMissionArchiveCmd_WithCleanup(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	
	// Create mission files
	missionID := "test-mission-789"
	missionContent := `---
id: ` + missionID + `
---
Body
`
	err := fs.MkdirAll(missionDir, 0755)
	require.NoError(t, err)
	
	// Create files that should be cleaned up
	filesToCreate := []string{"mission.md", "execution.log", "id", "plan.json"}
	for _, filename := range filesToCreate {
		content := "content"
		if filename == "mission.md" {
			content = missionContent
		}
		err = afero.WriteFile(fs, filepath.Join(missionDir, filename), []byte(content), 0644)
		require.NoError(t, err)
	}

	// Mock GitClient
	mockGit := &MockGitClient{
		commitMessage: "feat: test commit",
	}

	// Test archive with cleanup
	archiver := mission.NewArchiver(fs, missionDir, mockGit)
	
	// Archive
	err = archiver.Archive()
	require.NoError(t, err)
	
	// Cleanup
	err = archiver.CleanupObsoleteFiles()
	require.NoError(t, err)

	// Verify files were archived
	completedDir := filepath.Join(missionDir, "completed")
	archivedMissionPath := filepath.Join(completedDir, missionID+"-mission.md")
	exists, err := afero.Exists(fs, archivedMissionPath)
	require.NoError(t, err)
	require.True(t, exists, "mission.md should be archived")

	// Verify obsolete files were cleaned up
	for _, filename := range filesToCreate {
		exists, err := afero.Exists(fs, filepath.Join(missionDir, filename))
		require.NoError(t, err)
		require.False(t, exists, "%s should be cleaned up", filename)
	}
}
