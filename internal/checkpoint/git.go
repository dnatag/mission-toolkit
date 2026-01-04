package checkpoint

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

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
}

// CmdGitClient implements GitClient using the git CLI
type CmdGitClient struct {
	workDir string
}

// NewCmdGitClient creates a new CmdGitClient
func NewCmdGitClient(workDir string) *CmdGitClient {
	return &CmdGitClient{workDir: workDir}
}

func (c *CmdGitClient) run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = c.workDir
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func (c *CmdGitClient) Add(files []string) error {
	if len(files) == 0 {
		return nil
	}
	_, err := c.run(append([]string{"add"}, files...)...)
	return err
}

func (c *CmdGitClient) Commit(message string) (string, error) {
	output, err := c.run("commit", "-m", message)
	if err != nil {
		if strings.Contains(output, "nothing to commit") {
			return "", ErrNoChanges
		}
		return "", fmt.Errorf("git commit failed: %s", output)
	}
	return c.run("rev-parse", "HEAD")
}

func (c *CmdGitClient) CreateTag(name string, commitHash string) error {
	if commitHash == "" {
		_, err := c.run("tag", name)
		return err
	}
	_, err := c.run("tag", name, commitHash)
	return err
}

func (c *CmdGitClient) Restore(checkpointName string, files []string) error {
	if len(files) == 0 {
		return nil
	}

	if _, err := c.run("rev-parse", checkpointName); err != nil {
		return fmt.Errorf("checkpoint not found: %s", checkpointName)
	}

	_, err := c.run(append([]string{"checkout", checkpointName, "--"}, files...)...)
	return err
}

func (c *CmdGitClient) ListTags(prefix string) ([]string, error) {
	out, err := c.run("tag", "-l", prefix+"*")
	if err != nil || out == "" {
		return nil, err
	}
	return strings.Split(out, "\n"), nil
}

func (c *CmdGitClient) DeleteTag(name string) error {
	_, err := c.run("tag", "-d", name)
	return err
}

func (c *CmdGitClient) GetTagCommit(tagName string) (string, error) {
	return c.run("rev-parse", tagName+"^{commit}")
}

func (c *CmdGitClient) SoftReset(commitHash string) error {
	_, err := c.run("reset", "--soft", commitHash)
	return err
}
