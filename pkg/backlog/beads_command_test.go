// Package backlog provides tests for command execution abstraction.
package backlog

import (
	"fmt"
	"testing"
)

// mockCommandRunner implements CommandRunner for testing.
type mockCommandRunner struct {
	outputs map[string]string // args -> output
	errors  map[string]error  // args -> error
}

// newMockCommandRunner creates a new mockCommandRunner.
func newMockCommandRunner() *mockCommandRunner {
	return &mockCommandRunner{
		outputs: make(map[string]string),
		errors:  make(map[string]error),
	}
}

// setResponse sets the mock response for a given command.
func (m *mockCommandRunner) setResponse(args []string, output string, err error) {
	key := argsToString(args)
	m.outputs[key] = output
	m.errors[key] = err
}

// Run executes the mock command.
func (m *mockCommandRunner) Run(args ...string) (string, error) {
	key := argsToString(args)
	output, ok := m.outputs[key]
	if !ok {
		return "", fmt.Errorf("unexpected command: %s", key)
	}
	return output, m.errors[key]
}

// argsToString converts args to a string key for map lookup.
func argsToString(args []string) string {
	key := ""
	for i, arg := range args {
		if i > 0 {
			key += " "
		}
		key += arg
	}
	return key
}

func TestNewBDCommandRunner(t *testing.T) {
	runner := NewBDCommandRunner("/tmp")
	if runner == nil {
		t.Fatal("NewBDCommandRunner returned nil")
	}

	bdRunner, ok := runner.(*bdCommandRunner)
	if !ok {
		t.Fatal("NewBDCommandRunner did not return *bdCommandRunner")
	}

	if bdRunner.workDir != "/tmp" {
		t.Errorf("Expected workDir /tmp, got %s", bdRunner.workDir)
	}
}

func TestBDCommandRunnerImplementsCommandRunner(t *testing.T) {
	var _ CommandRunner = (*bdCommandRunner)(nil)
	// This test ensures bdCommandRunner implements CommandRunner at compile time
}
