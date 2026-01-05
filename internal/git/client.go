package git

import "errors"

var ErrNoChanges = errors.New("no changes to commit")

// GitClient defines the interface for git operations
type GitClient interface {
	Add(files []string) error
	Commit(message string) (string, error)
	CreateTag(name string, commitHash string) error
	Restore(checkpointName string, files []string) error
	ListTags(prefix string) ([]string, error)
	DeleteTag(name string) error
	GetTagCommit(tagName string) (string, error)
	SoftReset(commitHash string) error
	GetCommitMessage(commitHash string) (string, error)
	IsTracked(path string) (bool, error)
	GetCommitParent(commitHash string) (string, error)
}
