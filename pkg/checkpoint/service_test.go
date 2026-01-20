package checkpoint

import (
	"fmt"
	"path/filepath"
	"testing"

	internalgit "github.com/dnatag/mission-toolkit/pkg/git"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func setupTestRepo(t *testing.T) (afero.Fs, *git.Repository) {
	// Use memory filesystem
	// Initialize git repo in memory
	repo, err := git.Init(memory.NewStorage(), memfs.New())
	require.NoError(t, err)

	// Create initial commit so we have a HEAD
	wt, err := repo.Worktree()
	require.NoError(t, err)

	// We need to create a file in the worktree to commit
	fs, err := wt.Filesystem.Create("README.md")
	require.NoError(t, err)
	_, err = fs.Write([]byte("# Test"))
	require.NoError(t, err)
	err = fs.Close()
	require.NoError(t, err)

	_, err = wt.Add("README.md")
	require.NoError(t, err)

	_, err = wt.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{Name: "Test", Email: "test@example.com"},
	})
	require.NoError(t, err)

	return afero.NewMemMapFs(), repo
}

func createMissionFile(t *testing.T, fs afero.Fs, missionID string, scopeFiles []string) {
	missionDir := ".mission"
	err := fs.MkdirAll(missionDir, 0755)
	require.NoError(t, err)

	scopeContent := ""
	for _, f := range scopeFiles {
		scopeContent += fmt.Sprintf("- %s\n", f)
	}

	content := fmt.Sprintf(`---
id: %s
type: WET
track: 1
iteration: 1
status: active
---

## SCOPE
%s
`, missionID, scopeContent)

	err = afero.WriteFile(fs, filepath.Join(missionDir, "mission.md"), []byte(content), 0644)
	require.NoError(t, err)
}

func TestService_Create(t *testing.T) {
	fs, repo := setupTestRepo(t)

	missionID := "test-mission"
	scopeFile := "test.txt"
	createMissionFile(t, fs, missionID, []string{scopeFile})

	// Create scope file in afero FS
	err := afero.WriteFile(fs, scopeFile, []byte("test content"), 0644)
	require.NoError(t, err)

	// Use MemGitClient
	gitClient := internalgit.NewMemGitClient(repo, fs)
	svc := NewServiceWithGit(fs, ".mission", gitClient)

	// Create checkpoint
	name, err := svc.Create(missionID)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%s-1", missionID), name)

	// Verify tag exists
	_, err = repo.Tag(name)
	require.NoError(t, err)
}

func TestService_Create_CreatesBaselineTagOnFirstCheckpoint(t *testing.T) {
	fs, repo := setupTestRepo(t)

	missionID := "test-mission"
	scopeFile := "test.txt"
	createMissionFile(t, fs, missionID, []string{scopeFile})

	err := afero.WriteFile(fs, scopeFile, []byte("test content"), 0644)
	require.NoError(t, err)

	gitClient := internalgit.NewMemGitClient(repo, fs)
	svc := NewServiceWithGit(fs, ".mission", gitClient)

	name, err := svc.Create(missionID)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%s-1", missionID), name)

	baselineTag := fmt.Sprintf("%s-baseline", missionID)
	_, err = repo.Tag(baselineTag)
	require.NoError(t, err)
}

func TestService_Create_DoesNotCreateBaselineTagOnSubsequentCheckpoints(t *testing.T) {
	fs, repo := setupTestRepo(t)

	missionID := "test-mission"
	scopeFile := "test.txt"
	createMissionFile(t, fs, missionID, []string{scopeFile})

	err := afero.WriteFile(fs, scopeFile, []byte("test content"), 0644)
	require.NoError(t, err)

	gitClient := internalgit.NewMemGitClient(repo, fs)
	svc := NewServiceWithGit(fs, ".mission", gitClient)

	_, err = svc.Create(missionID)
	require.NoError(t, err)

	err = afero.WriteFile(fs, scopeFile, []byte("updated content"), 0644)
	require.NoError(t, err)

	name, err := svc.Create(missionID)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%s-2", missionID), name)

	tags, err := repo.Tags()
	require.NoError(t, err)

	baselineCount := 0
	baselineTag := fmt.Sprintf("%s-baseline", missionID)
	err = tags.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().Short() == baselineTag {
			baselineCount++
		}
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, 1, baselineCount, "baseline tag should only be created once")
}

func TestService_Create_OnlyScopeFiles(t *testing.T) {
	fs, repo := setupTestRepo(t)

	missionID := "test-scope"
	scopeFile := "scope.txt"
	otherFile := "other.txt"
	createMissionFile(t, fs, missionID, []string{scopeFile})

	// Create files in afero FS
	err := afero.WriteFile(fs, scopeFile, []byte("scope content"), 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, otherFile, []byte("other content"), 0644)
	require.NoError(t, err)

	// Use MemGitClient
	gitClient := internalgit.NewMemGitClient(repo, fs)
	svc := NewServiceWithGit(fs, ".mission", gitClient)

	// Create checkpoint
	name, err := svc.Create(missionID)
	require.NoError(t, err)

	// Verify commit content
	tagRef, err := repo.Tag(name)
	require.NoError(t, err)

	// Resolve annotated tag to commit
	tagObj, err := repo.TagObject(tagRef.Hash())
	require.NoError(t, err)

	commit, err := repo.CommitObject(tagObj.Target)
	require.NoError(t, err)

	// Check scope file is in commit
	_, err = commit.File(scopeFile)
	require.NoError(t, err)

	// Check other file is NOT in commit
	_, err = commit.File(otherFile)
	require.Error(t, err) // Should be file not found error
}

func TestService_Create_UntrackedFiles(t *testing.T) {
	fs, repo := setupTestRepo(t)

	missionID := "test-untracked"
	scopeFile := "scope.txt"
	untrackedFile := "untracked.txt"
	createMissionFile(t, fs, missionID, []string{scopeFile})

	// Create files in afero FS
	err := afero.WriteFile(fs, scopeFile, []byte("scope content"), 0644)
	require.NoError(t, err)
	err = afero.WriteFile(fs, untrackedFile, []byte("untracked content"), 0644)
	require.NoError(t, err)

	// Use MemGitClient
	gitClient := internalgit.NewMemGitClient(repo, fs)
	svc := NewServiceWithGit(fs, ".mission", gitClient)

	// Create checkpoint
	name, err := svc.Create(missionID)
	require.NoError(t, err)

	// Verify commit content
	tagRef, err := repo.Tag(name)
	require.NoError(t, err)
	tagObj, err := repo.TagObject(tagRef.Hash())
	require.NoError(t, err)
	commit, err := repo.CommitObject(tagObj.Target)
	require.NoError(t, err)

	// Check scope file is in commit
	_, err = commit.File(scopeFile)
	require.NoError(t, err)

	// Check untracked file is NOT in commit
	_, err = commit.File(untrackedFile)
	require.Error(t, err)
}

func TestService_Create_GitIgnoredFiles(t *testing.T) {
	fs, repo := setupTestRepo(t)

	missionID := "test-ignored"
	scopeFile := "ignored.log" // Assume .log is ignored
	createMissionFile(t, fs, missionID, []string{scopeFile})

	// Create .gitignore
	wt, _ := repo.Worktree()
	f, _ := wt.Filesystem.Create(".gitignore")
	_, err := f.Write([]byte("*.log\n"))
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)
	_, err = wt.Add(".gitignore")
	require.NoError(t, err)
	_, err = wt.Commit("Add gitignore", &git.CommitOptions{
		Author: &object.Signature{Name: "Test", Email: "test@example.com"},
	})
	require.NoError(t, err)

	// Create ignored file in afero FS
	err = afero.WriteFile(fs, scopeFile, []byte("ignored content"), 0644)
	require.NoError(t, err)

	// Use MemGitClient
	gitClient := internalgit.NewMemGitClient(repo, fs)
	svc := NewServiceWithGit(fs, ".mission", gitClient)

	// Create checkpoint - should succeed even if ignored because it's in scope
	name, err := svc.Create(missionID)
	require.NoError(t, err)

	// Verify commit content
	tagRef, err := repo.Tag(name)
	require.NoError(t, err)
	tagObj, err := repo.TagObject(tagRef.Hash())
	require.NoError(t, err)
	commit, err := repo.CommitObject(tagObj.Target)
	require.NoError(t, err)

	// Check ignored file IS in commit because it was in scope
	_, err = commit.File(scopeFile)
	require.NoError(t, err)
}

func TestService_Restore(t *testing.T) {
	fs, repo := setupTestRepo(t)

	missionID := "test-restore"
	scopeFile := "restore.txt"
	createMissionFile(t, fs, missionID, []string{scopeFile})

	// Initial state v1
	err := afero.WriteFile(fs, scopeFile, []byte("v1"), 0644)
	require.NoError(t, err)

	// Use MemGitClient
	gitClient := internalgit.NewMemGitClient(repo, fs)
	svc := NewServiceWithGit(fs, ".mission", gitClient)

	// Create checkpoint v1
	name1, err := svc.Create(missionID)
	require.NoError(t, err)

	// Modify to v2
	err = afero.WriteFile(fs, scopeFile, []byte("v2"), 0644)
	require.NoError(t, err)

	// Restore v1
	err = svc.Restore(name1)
	require.NoError(t, err)

	// Verify content in afero fs is v1
	content, err := afero.ReadFile(fs, scopeFile)
	require.NoError(t, err)
	require.Equal(t, "v1", string(content))
}

func TestService_Restore_UntrackedFiles(t *testing.T) {
	fs, repo := setupTestRepo(t)

	missionID := "test-restore-untracked"
	scopeFile := "restore.txt"
	untrackedFile := "untracked.txt"
	createMissionFile(t, fs, missionID, []string{scopeFile})

	// Initial state
	err := afero.WriteFile(fs, scopeFile, []byte("v1"), 0644)
	require.NoError(t, err)

	// Use MemGitClient
	gitClient := internalgit.NewMemGitClient(repo, fs)
	svc := NewServiceWithGit(fs, ".mission", gitClient)

	// Create checkpoint
	name1, err := svc.Create(missionID)
	require.NoError(t, err)

	// Create untracked file
	err = afero.WriteFile(fs, untrackedFile, []byte("untracked"), 0644)
	require.NoError(t, err)

	// Restore
	err = svc.Restore(name1)
	require.NoError(t, err)

	// Verify untracked file still exists and is untouched
	content, err := afero.ReadFile(fs, untrackedFile)
	require.NoError(t, err)
	require.Equal(t, "untracked", string(content))
}

func TestService_Clear(t *testing.T) {
	fs, repo := setupTestRepo(t)

	missionID := "test-clear"
	scopeFile := "clear.txt"
	createMissionFile(t, fs, missionID, []string{scopeFile})

	err := afero.WriteFile(fs, scopeFile, []byte("content"), 0644)
	require.NoError(t, err)

	// Use MemGitClient
	gitClient := internalgit.NewMemGitClient(repo, fs)
	svc := NewServiceWithGit(fs, ".mission", gitClient)

	// Create checkpoints
	_, err = svc.Create(missionID)
	require.NoError(t, err)

	err = afero.WriteFile(fs, scopeFile, []byte("content2"), 0644)
	require.NoError(t, err)
	_, err = svc.Create(missionID)
	require.NoError(t, err)

	// Clear
	count, err := svc.Clear(missionID)
	require.NoError(t, err)
	require.Equal(t, 3, count) // 2 checkpoints + 1 baseline tag

	// Verify tags gone
	tags, _ := gitClient.ListTags(missionID)
	require.Empty(t, tags)
}

func TestService_Consolidate(t *testing.T) {
	fs, repo := setupTestRepo(t)

	missionID := "test-consolidate"
	scopeFile1 := "file1.txt"
	scopeFile2 := "file2.txt"
	createMissionFile(t, fs, missionID, []string{scopeFile1, scopeFile2})

	// Use MemGitClient
	gitClient := internalgit.NewMemGitClient(repo, fs)
	svc := NewServiceWithGit(fs, ".mission", gitClient)

	// --- Checkpoint 1 ---
	err := afero.WriteFile(fs, scopeFile1, []byte("v1"), 0644)
	require.NoError(t, err)
	_, err = svc.Create(missionID)
	require.NoError(t, err)

	// --- Checkpoint 2 ---
	err = afero.WriteFile(fs, scopeFile2, []byte("v1"), 0644)
	require.NoError(t, err)
	_, err = svc.Create(missionID)
	require.NoError(t, err)

	// --- Final change (no checkpoint) ---
	err = afero.WriteFile(fs, scopeFile1, []byte("v2"), 0644)
	require.NoError(t, err)

	// Consolidate
	commitMsg := "Final commit"
	result, err := svc.Consolidate(missionID, commitMsg)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify final commit
	commit, err := repo.CommitObject(plumbing.NewHash(result.CommitHash))
	require.NoError(t, err)
	require.Equal(t, commitMsg, commit.Message)

	// Verify file contents in final commit
	f1, err := commit.File(scopeFile1)
	require.NoError(t, err)
	content1, _ := f1.Contents()
	require.Equal(t, "v2", content1)

	f2, err := commit.File(scopeFile2)
	require.NoError(t, err)
	content2, _ := f2.Contents()
	require.Equal(t, "v1", content2)

	// Verify checkpoints are cleared
	tags, err := gitClient.ListTags(missionID)
	require.NoError(t, err)
	require.Empty(t, tags)
}

func TestService_Consolidate_NoChanges(t *testing.T) {
	fs, repo := setupTestRepo(t)

	missionID := "test-consolidate-no-changes"
	scopeFile := "file.txt"
	createMissionFile(t, fs, missionID, []string{scopeFile})

	// Use MemGitClient
	gitClient := internalgit.NewMemGitClient(repo, fs)
	svc := NewServiceWithGit(fs, ".mission", gitClient)

	// Create file and commit it initially (so it's tracked)
	err := afero.WriteFile(fs, scopeFile, []byte("initial"), 0644)
	require.NoError(t, err)
	wt, _ := repo.Worktree()
	f, _ := wt.Filesystem.Create(scopeFile)
	f.Write([]byte("initial"))
	f.Close()
	wt.Add(scopeFile)
	wt.Commit("Initial state", &git.CommitOptions{
		Author: &object.Signature{Name: "Test", Email: "test@example.com"},
	})

	// Try to consolidate without any changes
	_, err = svc.Consolidate(missionID, "Final commit")
	require.ErrorContains(t, err, "creating final commit")
}

func TestService_Consolidate_WithUntrackedFile(t *testing.T) {
	fs, repo := setupTestRepo(t)

	missionID := "test-consolidate-untracked"
	scopeFile := "scope.txt"
	untrackedFile := "untracked.txt"
	createMissionFile(t, fs, missionID, []string{scopeFile})

	// Use MemGitClient
	gitClient := internalgit.NewMemGitClient(repo, fs)
	svc := NewServiceWithGit(fs, ".mission", gitClient)

	// Modify scope file
	err := afero.WriteFile(fs, scopeFile, []byte("v1"), 0644)
	require.NoError(t, err)

	// Create untracked file in worktree fs
	wt, _ := repo.Worktree()
	f, _ := wt.Filesystem.Create(untrackedFile)
	f.Write([]byte("untracked"))
	f.Close()

	// Consolidate
	result, err := svc.Consolidate(missionID, "Final commit")
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify final commit
	commit, err := repo.CommitObject(plumbing.NewHash(result.CommitHash))
	require.NoError(t, err)

	// Check scope file is in commit
	_, err = commit.File(scopeFile)
	require.NoError(t, err)

	// Check untracked file is NOT in commit
	_, err = commit.File(untrackedFile)
	require.Error(t, err)

	// Verify unstaged files includes the untracked file
	require.Contains(t, result.UnstagedFiles, untrackedFile)
}

func TestService_Consolidate_WithFileDeletion(t *testing.T) {
	fs, repo := setupTestRepo(t)

	missionID := "test-consolidate-deletion"
	scopeFile := "file.txt"
	createMissionFile(t, fs, missionID, []string{scopeFile})

	// Use MemGitClient
	gitClient := internalgit.NewMemGitClient(repo, fs)
	svc := NewServiceWithGit(fs, ".mission", gitClient)

	// Create file and commit it initially
	err := afero.WriteFile(fs, scopeFile, []byte("initial"), 0644)
	require.NoError(t, err)
	wt, _ := repo.Worktree()
	f, _ := wt.Filesystem.Create(scopeFile)
	f.Write([]byte("initial"))
	f.Close()
	wt.Add(scopeFile)
	wt.Commit("Initial state", &git.CommitOptions{
		Author: &object.Signature{Name: "Test", Email: "test@example.com"},
	})

	// Delete file
	err = fs.Remove(scopeFile)
	require.NoError(t, err)

	// Consolidate
	result, err := svc.Consolidate(missionID, "Final commit")
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify final commit
	commit, err := repo.CommitObject(plumbing.NewHash(result.CommitHash))
	require.NoError(t, err)

	// Check file is NOT in commit
	_, err = commit.File(scopeFile)
	require.Error(t, err)
}

// Edge case tests for checkpoint service

// TestService_Create_InvalidMissionIDFormat verifies that Create fails gracefully
// when the mission ID is empty or invalid.
func TestService_Create_InvalidMissionIDFormat(t *testing.T) {
	fs, repo := setupTestRepo(t)

	missionID := ""
	createMissionFile(t, fs, missionID, []string{"test.txt"})

	// Create test file in git worktree
	wt, _ := repo.Worktree()
	f, _ := wt.Filesystem.Create("test.txt")
	f.Write([]byte("content"))
	f.Close()

	gitClient := internalgit.NewMemGitClient(repo, fs)
	svc := NewServiceWithGit(fs, ".mission", gitClient)

	// Test with empty mission ID - should fail
	_, err := svc.Create(missionID)
	require.Error(t, err)
}

// TestService_Create_MultipleCheckpoints verifies that multiple checkpoints
// can be created for the same mission with incrementing numbers.
func TestService_Create_MultipleCheckpoints(t *testing.T) {
	fs, repo := setupTestRepo(t)

	missionID := "test-mission"
	scopeFile := "test.txt"
	createMissionFile(t, fs, missionID, []string{scopeFile})

	err := afero.WriteFile(fs, scopeFile, []byte("initial"), 0644)
	require.NoError(t, err)

	gitClient := internalgit.NewMemGitClient(repo, fs)
	svc := NewServiceWithGit(fs, ".mission", gitClient)

	// Create first checkpoint successfully
	name1, err := svc.Create(missionID)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%s-1", missionID), name1)

	// Verify tag exists
	_, err = repo.Tag(name1)
	require.NoError(t, err)
}

// TestService_Restore_NonExistentCheckpoint verifies that Restore fails with
// an appropriate error when attempting to restore a checkpoint that doesn't exist.
func TestService_Restore_NonExistentCheckpoint(t *testing.T) {
	fs, repo := setupTestRepo(t)

	missionID := "test-mission"
	createMissionFile(t, fs, missionID, []string{"test.txt"})

	gitClient := internalgit.NewMemGitClient(repo, fs)
	svc := NewServiceWithGit(fs, ".mission", gitClient)

	// Try to restore non-existent checkpoint
	err := svc.Restore("non-existent-checkpoint")
	require.Error(t, err)
	require.Contains(t, err.Error(), "tag not found")
}

func TestService_Restore_GitOperationsFail(t *testing.T) {
	fs, repo := setupTestRepo(t)

	missionID := "test-mission"
	scopeFile := "test.txt"
	createMissionFile(t, fs, missionID, []string{scopeFile})

	err := afero.WriteFile(fs, scopeFile, []byte("initial"), 0644)
	require.NoError(t, err)

	gitClient := internalgit.NewMemGitClient(repo, fs)
	svc := NewServiceWithGit(fs, ".mission", gitClient)

	// Create checkpoint
	name, err := svc.Create(missionID)
	require.NoError(t, err)

	// Modify file
	err = afero.WriteFile(fs, scopeFile, []byte("modified"), 0644)
	require.NoError(t, err)

	// Restore should work
	err = svc.Restore(name)
	require.NoError(t, err)

	// Verify file was restored
	content, err := afero.ReadFile(fs, scopeFile)
	require.NoError(t, err)
	require.Equal(t, "initial", string(content))
}

func TestService_Clear_NoCheckpointsToClean(t *testing.T) {
	fs, repo := setupTestRepo(t)
	gitClient := internalgit.NewMemGitClient(repo, fs)
	svc := NewServiceWithGit(fs, ".mission", gitClient)

	missionID := "test-mission"

	// Clear when no checkpoints exist
	count, err := svc.Clear(missionID)
	require.NoError(t, err)
	require.Equal(t, 0, count)
}

func TestService_Clear_RemovesAllCheckpoints(t *testing.T) {
	fs, repo := setupTestRepo(t)

	missionID := "test-mission"
	scopeFile := "test.txt"
	createMissionFile(t, fs, missionID, []string{scopeFile})

	err := afero.WriteFile(fs, scopeFile, []byte("content"), 0644)
	require.NoError(t, err)

	gitClient := internalgit.NewMemGitClient(repo, fs)
	svc := NewServiceWithGit(fs, ".mission", gitClient)

	// Create checkpoint (creates baseline tag)
	name, err := svc.Create(missionID)
	require.NoError(t, err)

	// Verify baseline tag exists
	baselineTag := fmt.Sprintf("%s-baseline", missionID)
	_, err = repo.Tag(baselineTag)
	require.NoError(t, err)

	// Clear should work and remove checkpoint (count includes baseline)
	count, err := svc.Clear(missionID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, count, 1)

	// Verify checkpoint tag is removed
	_, err = repo.Tag(name)
	require.Error(t, err)
}

func TestService_Consolidate_EmptyScope(t *testing.T) {
	fs, repo := setupTestRepo(t)

	missionID := "test-mission"
	createMissionFile(t, fs, missionID, []string{}) // Empty scope

	gitClient := internalgit.NewMemGitClient(repo, fs)
	svc := NewServiceWithGit(fs, ".mission", gitClient)

	// Consolidate with empty scope should fail
	_, err := svc.Consolidate(missionID, "Empty commit")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no files in mission scope")
}

func TestService_Consolidate_UnstagedFiles(t *testing.T) {
	fs, repo := setupTestRepo(t)

	missionID := "test-mission"
	scopeFile := "test.txt"
	unstagedFile := "unstaged.txt"
	createMissionFile(t, fs, missionID, []string{scopeFile})

	gitClient := internalgit.NewMemGitClient(repo, fs)
	svc := NewServiceWithGit(fs, ".mission", gitClient)

	// Create scope file in both filesystems
	err := afero.WriteFile(fs, scopeFile, []byte("content"), 0644)
	require.NoError(t, err)

	wt, _ := repo.Worktree()
	f, _ := wt.Filesystem.Create(scopeFile)
	f.Write([]byte("content"))
	f.Close()

	// Create unstaged file (not in scope) in git worktree only
	f2, _ := wt.Filesystem.Create(unstagedFile)
	f2.Write([]byte("unstaged"))
	f2.Close()

	// Consolidate
	result, err := svc.Consolidate(missionID, "Commit with unstaged")
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify unstaged file is reported
	require.Contains(t, result.UnstagedFiles, unstagedFile)
}

func TestService_Consolidate_CommitCreationFails(t *testing.T) {
	fs, repo := setupTestRepo(t)

	missionID := "test-mission"
	createMissionFile(t, fs, missionID, []string{"test.txt"})

	gitClient := internalgit.NewMemGitClient(repo, fs)
	svc := NewServiceWithGit(fs, ".mission", gitClient)

	// Try to consolidate with empty commit message (should fail)
	_, err := svc.Consolidate(missionID, "")
	require.Error(t, err)
}
