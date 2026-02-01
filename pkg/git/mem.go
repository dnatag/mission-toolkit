package git

import (
	"io"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/afero"
)

// MemGitClient implements GitClient using go-git in memory
type MemGitClient struct {
	repo *git.Repository
	fs   afero.Fs // The filesystem used by the service (afero)
}

// NewMemGitClient creates a new MemGitClient
func NewMemGitClient(repo *git.Repository, fs afero.Fs) *MemGitClient {
	return &MemGitClient{repo: repo, fs: fs}
}

func (c *MemGitClient) Add(files []string) error {
	wt, err := c.repo.Worktree()
	if err != nil {
		return err
	}
	for _, file := range files {
		// Check if file exists in afero fs
		exists, _ := afero.Exists(c.fs, file)
		if !exists {
			// If file doesn't exist, it might be a deletion
			// We need to remove it from the index
			_, err := wt.Remove(file)
			if err != nil {
				// Ignore if file is not in worktree
			}
			continue
		}
		// We need to ensure the file exists in the go-git worktree filesystem
		// Since we are mocking, we copy from afero fs to go-git fs
		content, err := afero.ReadFile(c.fs, file)
		if err != nil {
			return err
		}

		f, err := wt.Filesystem.Create(file)
		if err != nil {
			return err
		}
		_, err = f.Write(content)
		f.Close()
		if err != nil {
			return err
		}

		_, err = wt.Add(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *MemGitClient) Commit(message string) (string, error) {
	wt, err := c.repo.Worktree()
	if err != nil {
		return "", err
	}

	// Check for changes
	status, err := wt.Status()
	if err != nil {
		return "", err
	}
	if status.IsClean() {
		return "", ErrNoChanges
	}

	hash, err := wt.Commit(message, &git.CommitOptions{
		Author: &object.Signature{Name: "Test", Email: "test@example.com"},
	})
	if err != nil {
		return "", err
	}
	return hash.String(), nil
}

func (c *MemGitClient) CreateTag(name string, commitHash string) error {
	hash := plumbing.NewHash(commitHash)
	_, err := c.repo.CreateTag(name, hash, &git.CreateTagOptions{
		Message: name,
		Tagger:  &object.Signature{Name: "Mission Toolkit", Email: "mission@toolkit.local"},
	})
	return err
}

func (c *MemGitClient) Restore(checkpointName string, files []string) error {
	tagRef, err := c.repo.Tag(checkpointName)
	if err != nil {
		return err
	}

	commit, err := c.repo.CommitObject(tagRef.Hash())
	if err != nil {
		// Try annotated tag
		tagObj, errTag := c.repo.TagObject(tagRef.Hash())
		if errTag == nil {
			commit, err = c.repo.CommitObject(tagObj.Target)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	tree, err := commit.Tree()
	if err != nil {
		return err
	}

	wt, err := c.repo.Worktree()
	if err != nil {
		return err
	}

	for _, path := range files {
		file, err := tree.File(path)
		if err != nil {
			// File not in checkpoint, remove from fs
			c.fs.Remove(path)
			wt.Filesystem.Remove(path)
			continue
		}

		reader, err := file.Reader()
		if err != nil {
			return err
		}
		content, err := io.ReadAll(reader)
		reader.Close()
		if err != nil {
			return err
		}

		// Update afero fs
		if err := afero.WriteFile(c.fs, path, content, 0644); err != nil {
			return err
		}

		// Update go-git fs
		f, err := wt.Filesystem.Create(path)
		if err != nil {
			return err
		}
		f.Write(content)
		f.Close()
	}
	return nil
}

func (c *MemGitClient) ListTags(prefix string) ([]string, error) {
	var tags []string
	iter, err := c.repo.Tags()
	if err != nil {
		return nil, err
	}
	err = iter.ForEach(func(ref *plumbing.Reference) error {
		if strings.HasPrefix(ref.Name().Short(), prefix) {
			tags = append(tags, ref.Name().Short())
		}
		return nil
	})
	return tags, err
}

func (c *MemGitClient) DeleteTag(name string) error {
	return c.repo.DeleteTag(name)
}

func (c *MemGitClient) GetTagCommit(tagName string) (string, error) {
	tagRef, err := c.repo.Tag(tagName)
	if err == nil {
		if tagObj, err := c.repo.TagObject(tagRef.Hash()); err == nil {
			return tagObj.Target.String(), nil
		}
		return tagRef.Hash().String(), nil
	}

	// Fall back to interpreting tagName as commit-ish (e.g., prefix + "^{commit}")
	hash := plumbing.NewHash(tagName)
	if hash.IsZero() {
		return "", err
	}
	return hash.String(), nil
}

func (c *MemGitClient) SoftReset(commitHash string) error {
	wt, err := c.repo.Worktree()
	if err != nil {
		return err
	}
	hash := plumbing.NewHash(commitHash)
	return wt.Reset(&git.ResetOptions{
		Commit: hash,
		Mode:   git.SoftReset,
	})
}

func (c *MemGitClient) GetCommitMessage(commitHash string) (string, error) {
	var hash plumbing.Hash
	if commitHash == "HEAD" {
		ref, err := c.repo.Head()
		if err != nil {
			return "", err
		}
		hash = ref.Hash()
	} else {
		hash = plumbing.NewHash(commitHash)
	}

	commit, err := c.repo.CommitObject(hash)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(commit.Message), nil
}

func (c *MemGitClient) IsTracked(path string) (bool, error) {
	// Check the index directly to see if the file is tracked
	idx, err := c.repo.Storer.Index()
	if err != nil {
		return false, err
	}

	_, err = idx.Entry(path)
	if err == object.ErrEntryNotFound || (err != nil && strings.Contains(err.Error(), "entry not found")) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *MemGitClient) GetCommitParent(commitHash string) (string, error) {
	hash := plumbing.NewHash(commitHash)
	commit, err := c.repo.CommitObject(hash)
	if err != nil {
		return "", err
	}
	if commit.NumParents() == 0 {
		return "", nil // No parent (initial commit)
	}
	parent, err := commit.Parent(0)
	if err != nil {
		return "", err
	}
	return parent.Hash.String(), nil
}

func (c *MemGitClient) GetUnstagedFiles() ([]string, error) {
	wt, err := c.repo.Worktree()
	if err != nil {
		return nil, err
	}
	status, err := wt.Status()
	if err != nil {
		return nil, err
	}
	var files []string
	for path, s := range status {
		// Worktree status: untracked or modified but not staged
		if s.Worktree == git.Untracked || s.Worktree == git.Modified {
			files = append(files, path)
		}
	}
	return files, nil
}

// GetUntrackedFiles returns files that exist in the working directory but are not tracked by git.
func (c *MemGitClient) GetUntrackedFiles() ([]string, error) {
	wt, err := c.repo.Worktree()
	if err != nil {
		return nil, err
	}
	status, err := wt.Status()
	if err != nil {
		return nil, err
	}
	var files []string
	for path, s := range status {
		// Only untracked files
		if s.Worktree == git.Untracked {
			files = append(files, path)
		}
	}
	return files, nil
}
