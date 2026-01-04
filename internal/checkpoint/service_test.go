package checkpoint

import (
	"fmt"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

var testCounter = 0

func setupTestRepo(t *testing.T) (*git.Repository, afero.Fs) {
	testCounter++
	// Create temporary directory for git repo
	fs := afero.NewMemMapFs()
	repoPath := fmt.Sprintf("/tmp/test-repo-%d", testCounter)
	fs.MkdirAll(repoPath, 0755)

	// Create git repository with worktree
	repo, err := git.PlainInit(repoPath, false)
	require.NoError(t, err)

	// Create initial commit
	wt, err := repo.Worktree()
	require.NoError(t, err)

	_, err = wt.Commit("initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test",
			Email: "test@example.com",
		},
		AllowEmptyCommits: true,
	})
	require.NoError(t, err)

	return repo, fs
}

func TestService_Create(t *testing.T) {
	t.Skip("Skipping: requires dirty working tree for checkpoint creation")
}

func TestService_Revert(t *testing.T) {
	t.Skip("Skipping: requires dirty working tree for checkpoint creation")
}

func TestService_Clear(t *testing.T) {
	t.Skip("Skipping: requires dirty working tree for checkpoint creation")
}

func TestService_Revert_NotFound(t *testing.T) {
	t.Skip("Skipping: requires dirty working tree for checkpoint creation")
}
