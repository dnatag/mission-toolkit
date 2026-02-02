package git

import (
	"fmt"
	"os/exec"
	"strings"
)

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
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}

func (c *CmdGitClient) Add(files []string) error {
	if len(files) == 0 {
		return nil
	}
	args := append([]string{"add"}, files...)
	_, err := c.run(args...)
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

	// Get the commit hash
	out, err := c.run("rev-parse", "HEAD")
	return strings.TrimSpace(out), err
}

func (c *CmdGitClient) CreateTag(name string, commitHash string) error {
	args := []string{"tag", name}
	if commitHash != "" {
		args = append(args, strings.TrimSpace(commitHash))
	}
	_, err := c.run(args...)
	return err
}

func (c *CmdGitClient) Restore(checkpointName string, files []string) error {
	if len(files) == 0 {
		return nil
	}

	// Check if tag exists
	if _, err := c.run("rev-parse", checkpointName); err != nil {
		return fmt.Errorf("checkpoint not found: %s", checkpointName)
	}

	// Checkout files from the checkpoint
	args := []string{"checkout", checkpointName, "--"}
	args = append(args, files...)
	_, err := c.run(args...)
	return err
}

func (c *CmdGitClient) ListTags(prefix string) ([]string, error) {
	out, err := c.run("tag", "-l", prefix+"*")
	if err != nil {
		return nil, err
	}
	var tags []string
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if line != "" {
			tags = append(tags, line)
		}
	}
	return tags, nil
}

func (c *CmdGitClient) DeleteTag(name string) error {
	_, err := c.run("tag", "-d", name)
	return err
}

func (c *CmdGitClient) GetTagCommit(tagName string) (string, error) {
	ref := tagName
	if !strings.Contains(tagName, "^{") {
		ref = tagName + "^{commit}"
	}

	out, err := c.run("rev-parse", ref)
	return strings.TrimSpace(out), err
}

func (c *CmdGitClient) SoftReset(commitHash string) error {
	_, err := c.run("reset", "--soft", commitHash)
	return err
}

func (c *CmdGitClient) GetCommitMessage(commitHash string) (string, error) {
	out, err := c.run("log", "-1", "--pretty=%B", commitHash)
	return strings.TrimSpace(out), err
}

func (c *CmdGitClient) IsTracked(path string) (bool, error) {
	// git ls-files --error-unmatch <file> returns 0 if tracked, 1 if not
	_, err := c.run("ls-files", "--error-unmatch", path)
	if err != nil {
		// If exit code is 1, it's not tracked.
		// We treat this as false, not an error.
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c *CmdGitClient) GetCommitParent(commitHash string) (string, error) {
	out, err := c.run("rev-parse", commitHash+"^")
	return strings.TrimSpace(out), err
}

func (c *CmdGitClient) GetUnstagedFiles() ([]string, error) {
	out, err := c.run("status", "--porcelain")
	if err != nil {
		return nil, err
	}
	var files []string
	for _, line := range strings.Split(out, "\n") {
		if len(line) < 3 {
			continue
		}
		// Porcelain format: XY filename (X=index, Y=worktree)
		// Unstaged: Y is not space and X is space or ?
		x, y := line[0], line[1]
		if y != ' ' && (x == ' ' || x == '?') {
			files = append(files, strings.TrimSpace(line[3:]))
		}
	}
	return files, nil
}

// GetUntrackedFiles returns files that exist in the working directory but are not tracked by git.
// These files have status "??" in git status --porcelain output.
func (c *CmdGitClient) GetUntrackedFiles() ([]string, error) {
	out, err := c.run("status", "--porcelain")
	if err != nil {
		return nil, err
	}
	var files []string
	for _, line := range strings.Split(out, "\n") {
		if len(line) < 3 {
			continue
		}
		// Porcelain format: ?? filename means untracked
		if line[0] == '?' && line[1] == '?' {
			files = append(files, strings.TrimSpace(line[3:]))
		}
	}
	return files, nil
}
