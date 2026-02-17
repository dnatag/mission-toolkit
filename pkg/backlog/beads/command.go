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
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("bd %v failed: %w", args, err)
	}
	return string(output), nil
}

// extractJSON strips non-JSON preamble (warnings, etc.) from bd output.
func extractJSON(output string) string {
	for i, c := range output {
		if c == '{' || c == '[' {
			if i > 0 {
				return output[i:]
			}
			break
		}
	}
	return output
}
