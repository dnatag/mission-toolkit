package git

import (
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRepo creates an in-memory git repository for testing
// Following the pattern from internal/checkpoint/service_test.go
func setupTestRepo(t *testing.T) (afero.Fs, *git.Repository) {
	t.Helper()

	// Create in-memory git repository
	repo, err := git.Init(memory.NewStorage(), memfs.New())
	require.NoError(t, err, "Failed to create in-memory git repository")

	// Create initial commit so we have a HEAD
	wt, err := repo.Worktree()
	require.NoError(t, err, "Failed to get worktree")

	// Create a file in the worktree to commit
	fs, err := wt.Filesystem.Create("README.md")
	require.NoError(t, err, "Failed to create README.md")
	_, err = fs.Write([]byte("# Test Repository"))
	require.NoError(t, err, "Failed to write README.md")
	err = fs.Close()
	require.NoError(t, err, "Failed to close README.md")

	_, err = wt.Add("README.md")
	require.NoError(t, err, "Failed to add README.md")

	_, err = wt.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{Name: "Test User", Email: "test@example.com"},
	})
	require.NoError(t, err, "Failed to create initial commit")

	return afero.NewMemMapFs(), repo
}

func TestMemGitClient_Add(t *testing.T) {
	tests := []struct {
		name        string
		setupFiles  map[string]string // files to create in afero fs
		addFiles    []string          // files to add via git client
		wantErr     bool
		errContains string
	}{
		{
			name: "add single existing file",
			setupFiles: map[string]string{
				"test.txt": "content",
			},
			addFiles: []string{"test.txt"},
			wantErr:  false,
		},
		{
			name: "add multiple existing files",
			setupFiles: map[string]string{
				"file1.txt": "content1",
				"file2.txt": "content2",
				"file3.txt": "content3",
			},
			addFiles: []string{"file1.txt", "file2.txt", "file3.txt"},
			wantErr:  false,
		},
		{
			name:       "add non-existent file",
			setupFiles: map[string]string{},
			addFiles:   []string{"nonexistent.txt"},
			wantErr:    false, // Should not error, just skip
		},
		{
			name: "add file that was deleted",
			setupFiles: map[string]string{
				"delete_me.txt": "will be deleted",
			},
			addFiles: []string{"delete_me.txt"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs, repo := setupTestRepo(t)

			// Create setup files in afero fs
			for path, content := range tt.setupFiles {
				err := afero.WriteFile(fs, path, []byte(content), 0644)
				require.NoError(t, err)
			}

			client := NewMemGitClient(repo, fs)
			err := client.Add(tt.addFiles)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMemGitClient_Commit(t *testing.T) {
	tests := []struct {
		name          string
		message       string
		setupAddFile  bool // whether to add a file before committing
		wantErr       bool
		errIs         error
		errContains   string
		wantHashEmpty bool
	}{
		{
			name:          "commit with changes",
			message:       "Test commit",
			setupAddFile:  true,
			wantErr:       false,
			wantHashEmpty: false,
		},
		{
			name:          "commit without changes returns ErrNoChanges",
			message:       "No changes commit",
			setupAddFile:  false,
			wantErr:       true,
			errIs:         ErrNoChanges,
			wantHashEmpty: true,
		},
		{
			name:          "empty commit message",
			message:       "",
			setupAddFile:  true,
			wantErr:       false,
			wantHashEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs, repo := setupTestRepo(t)
			client := NewMemGitClient(repo, fs)

			if tt.setupAddFile {
				// Create and add a file to have changes
				err := afero.WriteFile(fs, "newfile.txt", []byte("new content"), 0644)
				require.NoError(t, err)
				err = client.Add([]string{"newfile.txt"})
				require.NoError(t, err)
			}

			hash, err := client.Commit(tt.message)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errIs != nil {
					assert.ErrorIs(t, err, tt.errIs)
				}
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				if tt.wantHashEmpty {
					assert.Empty(t, hash)
				} else {
					assert.NotEmpty(t, hash, "commit hash should not be empty")
					// Verify hash is a valid git hash (40 hex chars for SHA-1, or short form)
					_, err := repo.CommitObject(plumbing.NewHash(hash))
					assert.NoError(t, err, "commit hash should be valid")
				}
			}
		})
	}
}

func TestMemGitClient_CreateTag(t *testing.T) {
	tests := []struct {
		name        string
		tagName     string
		commitHash  string
		wantErr     bool
		errContains string
	}{
		{
			name:       "create tag on HEAD commit",
			tagName:    "v1.0.0",
			commitHash: "HEAD",
			wantErr:    false,
		},
		{
			name:       "create tag with specific commit hash",
			tagName:    "test-tag",
			commitHash: "", // Will use HEAD
			wantErr:    false,
		},
		{
			name:        "create tag with invalid commit hash",
			tagName:     "invalid-tag",
			commitHash:  "invalid-hash-123",
			wantErr:     true,
			errContains: "object not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs, repo := setupTestRepo(t)
			client := NewMemGitClient(repo, fs)

			// Get HEAD hash if not specified
			commitHash := tt.commitHash
			if commitHash == "" || commitHash == "HEAD" {
				head, err := repo.Head()
				require.NoError(t, err)
				commitHash = head.Hash().String()
			}

			err := client.CreateTag(tt.tagName, commitHash)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				// Verify tag exists
				tagRef, err := repo.Tag(tt.tagName)
				require.NoError(t, err, "tag should exist after creation")
				assert.NotNil(t, tagRef)
			}
		})
	}
}

func TestMemGitClient_ListTags(t *testing.T) {
	tests := []struct {
		name      string
		setupTags []string // tags to create before testing
		prefix    string   // prefix to filter by
		wantCount int      // expected number of tags
		wantTags  []string // expected tag names
	}{
		{
			name:      "list all tags",
			setupTags: []string{"v1.0.0", "v2.0.0", "feature-1"},
			prefix:    "",
			wantCount: 3,
			wantTags:  []string{"v1.0.0", "v2.0.0", "feature-1"},
		},
		{
			name:      "list tags with prefix",
			setupTags: []string{"v1.0.0", "v2.0.0", "feature-1"},
			prefix:    "v",
			wantCount: 2,
			wantTags:  []string{"v1.0.0", "v2.0.0"},
		},
		{
			name:      "list tags with no matches",
			setupTags: []string{"v1.0.0", "v2.0.0"},
			prefix:    "x",
			wantCount: 0,
			wantTags:  []string{},
		},
		{
			name:      "list tags when none exist",
			setupTags: []string{},
			prefix:    "",
			wantCount: 0,
			wantTags:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs, repo := setupTestRepo(t)
			client := NewMemGitClient(repo, fs)

			// Create setup tags
			head, _ := repo.Head()
			headHash := head.Hash().String()
			for _, tag := range tt.setupTags {
				err := client.CreateTag(tag, headHash)
				require.NoError(t, err)
			}

			tags, err := client.ListTags(tt.prefix)
			require.NoError(t, err)
			assert.Len(t, tags, tt.wantCount)

			if len(tt.wantTags) > 0 {
				for _, wantTag := range tt.wantTags {
					assert.Contains(t, tags, wantTag)
				}
			}
		})
	}
}

func TestMemGitClient_DeleteTag(t *testing.T) {
	fs, repo := setupTestRepo(t)
	client := NewMemGitClient(repo, fs)

	// Create a tag first
	head, _ := repo.Head()
	headHash := head.Hash().String()
	err := client.CreateTag("test-tag", headHash)
	require.NoError(t, err)

	// Verify tag exists
	_, err = repo.Tag("test-tag")
	require.NoError(t, err)

	// Delete the tag
	err = client.DeleteTag("test-tag")
	require.NoError(t, err)

	// Verify tag is gone
	_, err = repo.Tag("test-tag")
	assert.Error(t, err, "tag should not exist after deletion")
}

func TestMemGitClient_GetTagCommit(t *testing.T) {
	tests := []struct {
		name       string
		tagName    string
		setupTag   bool
		wantErr    bool
		wantCommit bool
	}{
		{
			name:       "get commit from existing tag",
			tagName:    "v1.0.0",
			setupTag:   true,
			wantErr:    false,
			wantCommit: true,
		},
		{
			name:       "get commit from non-existent tag",
			tagName:    "nonexistent",
			setupTag:   false,
			wantErr:    true,
			wantCommit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs, repo := setupTestRepo(t)
			client := NewMemGitClient(repo, fs)

			if tt.setupTag {
				head, _ := repo.Head()
				headHash := head.Hash().String()
				err := client.CreateTag(tt.tagName, headHash)
				require.NoError(t, err)
			}

			commitHash, err := client.GetTagCommit(tt.tagName)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.wantCommit {
					assert.NotEmpty(t, commitHash)
					// Verify it's a valid commit
					_, err := repo.CommitObject(plumbing.NewHash(commitHash))
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestMemGitClient_Restore(t *testing.T) {
	tests := []struct {
		name         string
		setupFiles   map[string]string // initial files
		modifyFiles  map[string]string // modifications before restore
		tagName      string
		restoreFiles []string // files to restore
		wantErr      bool
		verifyFiles  map[string]string // expected file contents after restore
	}{
		{
			name: "restore single file to checkpoint state",
			setupFiles: map[string]string{
				"test.txt": "original",
			},
			modifyFiles: map[string]string{
				"test.txt": "modified",
			},
			tagName:      "checkpoint-1",
			restoreFiles: []string{"test.txt"},
			wantErr:      false,
			verifyFiles: map[string]string{
				"test.txt": "original",
			},
		},
		{
			name: "restore multiple files to checkpoint state",
			setupFiles: map[string]string{
				"file1.txt": "v1",
				"file2.txt": "v2",
			},
			modifyFiles: map[string]string{
				"file1.txt": "v1-modified",
				"file2.txt": "v2-modified",
			},
			tagName:      "checkpoint-1",
			restoreFiles: []string{"file1.txt", "file2.txt"},
			wantErr:      false,
			verifyFiles: map[string]string{
				"file1.txt": "v1",
				"file2.txt": "v2",
			},
		},
		{
			name: "restore file that was deleted in checkpoint",
			setupFiles: map[string]string{
				"existing.txt": "content",
			},
			modifyFiles: map[string]string{
				"deleted.txt": "will be deleted",
			},
			tagName:      "checkpoint-1",
			restoreFiles: []string{"deleted.txt"},
			wantErr:      false,
			verifyFiles: map[string]string{
				"existing.txt": "content",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs, repo := setupTestRepo(t)
			client := NewMemGitClient(repo, fs)

			// Create initial files and add to git
			for path, content := range tt.setupFiles {
				err := afero.WriteFile(fs, path, []byte(content), 0644)
				require.NoError(t, err)
				err = client.Add([]string{path})
				require.NoError(t, err)
			}

			// Create checkpoint
			_, err := client.Commit("Initial state")
			require.NoError(t, err)
			head, _ := repo.Head()
			err = client.CreateTag(tt.tagName, head.Hash().String())
			require.NoError(t, err)

			// Modify files
			for path, content := range tt.modifyFiles {
				err := afero.WriteFile(fs, path, []byte(content), 0644)
				require.NoError(t, err)
			}

			// Restore
			err = client.Restore(tt.tagName, tt.restoreFiles)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				// Verify file contents
				for path, expectedContent := range tt.verifyFiles {
					content, err := afero.ReadFile(fs, path)
					if expectedContent == "" {
						// File should not exist
						assert.Error(t, err)
					} else {
						require.NoError(t, err)
						assert.Equal(t, expectedContent, string(content))
					}
				}
			}
		})
	}
}

func TestMemGitClient_SoftReset(t *testing.T) {
	fs, repo := setupTestRepo(t)
	client := NewMemGitClient(repo, fs)

	// Create two commits
	err := afero.WriteFile(fs, "file1.txt", []byte("content1"), 0644)
	require.NoError(t, err)
	err = client.Add([]string{"file1.txt"})
	require.NoError(t, err)
	commit1, err := client.Commit("First commit")
	require.NoError(t, err)

	err = afero.WriteFile(fs, "file2.txt", []byte("content2"), 0644)
	require.NoError(t, err)
	err = client.Add([]string{"file2.txt"})
	require.NoError(t, err)
	_, err = client.Commit("Second commit")
	require.NoError(t, err)

	// Reset to first commit
	err = client.SoftReset(commit1)
	require.NoError(t, err)

	// Verify HEAD points to first commit
	head, err := repo.Head()
	require.NoError(t, err)
	assert.Equal(t, commit1, head.Hash().String())
}

func TestMemGitClient_GetCommitMessage(t *testing.T) {
	tests := []struct {
		name         string
		commitHash   string
		commitMsg    string
		wantErr      bool
		wantContains string
	}{
		{
			name:         "get message for HEAD",
			commitHash:   "HEAD",
			commitMsg:    "Test commit message",
			wantErr:      false,
			wantContains: "Test commit message",
		},
		{
			name:         "get message for specific commit",
			commitHash:   "", // Will be filled with actual hash
			commitMsg:    "Specific commit",
			wantErr:      false,
			wantContains: "Specific commit",
		},
		{
			name:         "get message for invalid commit",
			commitHash:   "invalid-hash",
			commitMsg:    "",
			wantErr:      true,
			wantContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs, repo := setupTestRepo(t)
			client := NewMemGitClient(repo, fs)

			// Create a commit with specific message
			err := afero.WriteFile(fs, "test.txt", []byte("test"), 0644)
			require.NoError(t, err)
			err = client.Add([]string{"test.txt"})
			require.NoError(t, err)
			hash, err := client.Commit(tt.commitMsg)
			require.NoError(t, err)

			// Use the actual hash if commitHash is empty
			commitHash := tt.commitHash
			if commitHash == "" {
				commitHash = hash
			}

			message, err := client.GetCommitMessage(commitHash)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, message, tt.wantContains)
			}
		})
	}
}

func TestMemGitClient_IsTracked(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		addFile  bool
		wantBool bool
	}{
		{
			name:     "tracked file",
			path:     "tracked.txt",
			addFile:  true,
			wantBool: true,
		},
		{
			name:     "untracked file",
			path:     "untracked.txt",
			addFile:  false,
			wantBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs, repo := setupTestRepo(t)
			client := NewMemGitClient(repo, fs)

			if tt.addFile {
				err := afero.WriteFile(fs, tt.path, []byte("content"), 0644)
				require.NoError(t, err)
				err = client.Add([]string{tt.path})
				require.NoError(t, err)
			}

			tracked, err := client.IsTracked(tt.path)
			require.NoError(t, err)
			assert.Equal(t, tt.wantBool, tracked)
		})
	}
}

func TestMemGitClient_GetCommitParent(t *testing.T) {
	fs, repo := setupTestRepo(t)
	client := NewMemGitClient(repo, fs)

	// Get the initial commit from setupTestRepo as our baseline
	head, _ := repo.Head()
	initialCommit := head.Hash().String()

	// Create first new commit
	err := afero.WriteFile(fs, "file1.txt", []byte("content1"), 0644)
	require.NoError(t, err)
	err = client.Add([]string{"file1.txt"})
	require.NoError(t, err)
	commit1, err := client.Commit("First commit")
	require.NoError(t, err)

	// Create second commit
	err = afero.WriteFile(fs, "file2.txt", []byte("content2"), 0644)
	require.NoError(t, err)
	err = client.Add([]string{"file2.txt"})
	require.NoError(t, err)
	commit2, err := client.Commit("Second commit")
	require.NoError(t, err)

	// Test: Get parent of second commit (should be first commit)
	parent, err := client.GetCommitParent(commit2)
	require.NoError(t, err)
	assert.Equal(t, commit1, parent)

	// Test: Get parent of first commit (should be the initial commit from setup)
	parent, err = client.GetCommitParent(commit1)
	require.NoError(t, err)
	assert.Equal(t, initialCommit, parent)

	// Test: Get parent of the very first commit (should be empty)
	parent, err = client.GetCommitParent(initialCommit)
	require.NoError(t, err)
	assert.Empty(t, parent)
}

func TestMemGitClient_GetUnstagedFiles(t *testing.T) {
	fs, repo := setupTestRepo(t)
	client := NewMemGitClient(repo, fs)

	// Initially no unstaged files
	unstaged, err := client.GetUnstagedFiles()
	require.NoError(t, err)
	assert.Empty(t, unstaged)

	// Create untracked file in worktree fs
	wt, _ := repo.Worktree()
	f, _ := wt.Filesystem.Create("untracked.txt")
	f.Write([]byte("untracked"))
	f.Close()

	unstaged, err = client.GetUnstagedFiles()
	require.NoError(t, err)
	assert.Contains(t, unstaged, "untracked.txt")

	// Stage and commit the file
	wt.Add("untracked.txt")
	wt.Commit("Add file", &git.CommitOptions{
		Author: &object.Signature{Name: "Test", Email: "test@example.com"},
	})

	// Now no unstaged files
	unstaged, err = client.GetUnstagedFiles()
	require.NoError(t, err)
	assert.Empty(t, unstaged)
}
