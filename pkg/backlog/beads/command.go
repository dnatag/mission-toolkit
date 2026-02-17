package beads

import (
	"fmt"
	"os/exec"
)

// CommandRunner defines the interface for executing external commands.
type CommandRunner interface {
	Run(args ...string) (string, error)
}

type bdCommandRunner struct {
	workDir string
}

// NewBDCommandRunner creates a new CommandRunner for the bd CLI.
func NewBDCommandRunner(workDir string) CommandRunner {
	return &bdCommandRunner{workDir: workDir}
}

func (r *bdCommandRunner) Run(args ...string) (string, error) {
	cmd := exec.Command("bd", args...)
	cmd.Dir = r.workDir
	output, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return string(ee.Stderr), fmt.Errorf("bd %v failed:\n%s", args, ee.Stderr)
		}
		return "", fmt.Errorf("bd %v failed: %w", args, err)
	}
	return string(output), nil
}
