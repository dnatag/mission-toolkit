// Package backlog provides command execution abstraction for Beads integration.
package backlog

import (
	"fmt"
	"os/exec"
)

// CommandRunner defines the interface for executing external commands.
// This abstraction enables testability by allowing mock implementations.
type CommandRunner interface {
	Run(args ...string) (string, error)
}

// bdCommandRunner implements CommandRunner for the bd CLI.
type bdCommandRunner struct {
	workDir string
}

// NewBDCommandRunner creates a new bdCommandRunner for the given working directory.
// The workDir specifies the directory where the bd command will be executed.
func NewBDCommandRunner(workDir string) CommandRunner {
	return &bdCommandRunner{workDir: workDir}
}

// Run executes a bd command with the provided arguments.
// It returns the combined stdout/stderr output and any error encountered.
// The command is executed in the configured working directory.
func (r *bdCommandRunner) Run(args ...string) (string, error) {
	cmd := exec.Command("bd", args...)
	cmd.Dir = r.workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("bd %v failed: %w", args, err)
	}
	return string(output), nil
}
