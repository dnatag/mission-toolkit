package tui

import (
	"strings"
	"testing"

	"github.com/dnatag/mission-toolkit/internal/mission"
)

func TestNewDashboardModel(t *testing.T) {
	model := NewDashboardModel()

	if model.currentPane != MissionPane {
		t.Errorf("Expected initial pane to be MissionPane, got %v", model.currentPane)
	}

	if model.itemsPerPage != 5 {
		t.Errorf("Expected itemsPerPage to be 5, got %d", model.itemsPerPage)
	}
}

func TestDashboardModel_PaneNavigation(t *testing.T) {
	model := NewDashboardModel()

	// Set up a completed mission for testing 3-pane layout
	model.selectedMission = &mission.Mission{
		ID:     "test-mission",
		Status: "completed",
		Type:   "WET",
		Track:  2,
		Body:   "## INTENT\nTest mission",
	}

	// Test pane cycling for completed mission (3 panes)
	if model.currentPane != MissionPane {
		t.Errorf("Expected initial pane to be MissionPane")
	}

	// Simulate tab press
	model.currentPane = Pane((int(model.currentPane) + 1) % 3)
	if model.currentPane != ExecutionLogPane {
		t.Errorf("Expected pane to be ExecutionLogPane after tab")
	}

	model.currentPane = Pane((int(model.currentPane) + 1) % 3)
	if model.currentPane != CommitPane {
		t.Errorf("Expected pane to be CommitPane after second tab")
	}

	model.currentPane = Pane((int(model.currentPane) + 1) % 3)
	if model.currentPane != MissionPane {
		t.Errorf("Expected pane to cycle back to MissionPane")
	}
}

func TestDashboardModel_ActiveMissionLayout(t *testing.T) {
	model := NewDashboardModel()

	// Set up an active mission for testing 2-pane layout
	model.selectedMission = &mission.Mission{
		ID:     "active-mission",
		Status: "active",
		Type:   "WET",
		Track:  2,
		Body:   "## INTENT\nActive test mission",
	}

	// Test pane cycling for active mission (2 panes)
	model.currentPane = Pane((int(model.currentPane) + 1) % 2)
	if model.currentPane != ExecutionLogPane {
		t.Errorf("Expected pane to be ExecutionLogPane for active mission")
	}

	model.currentPane = Pane((int(model.currentPane) + 1) % 2)
	if model.currentPane != MissionPane {
		t.Errorf("Expected pane to cycle back to MissionPane for active mission")
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"this is a very long string", 10, "this is..."},
		{"exactly10c", 10, "exactly10c"},
		{"", 5, ""},
		{"abc", 3, "abc"},
	}

	for _, test := range tests {
		result := truncate(test.input, test.maxLen)
		if result != test.expected {
			t.Errorf("truncate(%q, %d) = %q, expected %q",
				test.input, test.maxLen, result, test.expected)
		}
	}
}

func TestRenderExecutionLogPane(t *testing.T) {
	model := NewDashboardModel()

	// Test loading state
	result := model.renderExecutionLogPane()
	if !strings.Contains(result, "Loading") {
		t.Errorf("Expected loading message, got: %s", result)
	}

	// Test loaded state with content
	model.executionLogLoaded = true
	model.executionLog = "line1\nline2\nline3"
	result = model.renderExecutionLogPane()
	if !strings.Contains(result, "line1") {
		t.Errorf("Expected log content, got: %s", result)
	}

	// Test empty log
	model.executionLog = ""
	result = model.renderExecutionLogPane()
	if !strings.Contains(result, "No execution log") {
		t.Errorf("Expected no log message, got: %s", result)
	}
}

func TestRenderCommitPane(t *testing.T) {
	model := NewDashboardModel()

	// Test loading state
	result := model.renderCommitPane()
	if !strings.Contains(result, "Loading") {
		t.Errorf("Expected loading message, got: %s", result)
	}

	// Test loaded state with content
	model.commitLoaded = true
	model.commitMessage = "feat: add dashboard\n\nHash: abc123"
	result = model.renderCommitPane()
	if !strings.Contains(result, "feat: add dashboard") {
		t.Errorf("Expected commit content, got: %s", result)
	}
}
