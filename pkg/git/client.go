package git

import "errors"

var ErrNoChanges = errors.New("no changes to commit")

// GitClient defines the interface for git operations
type GitClient interface {
	Add(files []string) error
	Commit(message string) (string, error)
	// CommitNoVerify creates a commit without running git hooks (pre-commit, commit-msg, etc.).
	// This is used for internal checkpoint commits to avoid hook interference.
	CommitNoVerify(message string) (string, error)
	CreateTag(name string, commitHash string) error
	Restore(checkpointName string, files []string) error
	ListTags(prefix string) ([]string, error)
	DeleteTag(name string) error
	GetTagCommit(tagName string) (string, error)
	SoftReset(commitHash string) error
	GetCommitMessage(commitHash string) (string, error)
	IsTracked(path string) (bool, error)
	GetCommitParent(commitHash string) (string, error)
	GetUnstagedFiles() ([]string, error)
	// GetUntrackedFiles returns a list of files that exist in the working directory
	// but are not tracked by git (status "??"). These files need manual cleanup.
	GetUntrackedFiles() ([]string, error)
}
