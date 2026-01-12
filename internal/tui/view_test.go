package tui

import (
	"strings"
	"testing"

	"github.com/dnatag/mission-toolkit/internal/mission"
)

func TestRenderNoMission(t *testing.T) {
	m := NewDashboardModel()
	output := m.renderNoMission()

	if !strings.Contains(output, "NO ACTIVE MISSION") {
		t.Error("expected output to contain 'NO ACTIVE MISSION'")
	}
}

func TestRenderCompletedMissions_Empty(t *testing.T) {
	m := NewDashboardModel()
	output := m.renderCompletedMissions()

	if output != "No completed missions found" {
		t.Errorf("expected 'No completed missions found', got '%s'", output)
	}
}

func TestRenderMissionDetails(t *testing.T) {
	m := NewDashboardModel()
	testMission := &mission.Mission{
		ID:     "test-123",
		Status: "completed",
		Type:   "WET",
		Track:  2,
		Body:   "## INTENT\nTest intent",
	}

	output := m.renderMissionDetails(testMission)
	if !strings.Contains(output, "COMPLETED") {
		t.Error("expected output to contain 'COMPLETED'")
	}
	if !strings.Contains(output, "Test intent") {
		t.Error("expected output to contain 'Test intent'")
	}
}

func TestRenderDashboardView(t *testing.T) {
	m := NewDashboardModel()
	m.width = 100 // Set width to trigger pane calculation
	m.calculatePaneWidths()

	// Case 1: Active Mission (2 panes)
	m.selectedMission = &mission.Mission{ID: "active", Status: "active", Body: "## INTENT\nActive"}
	m.currentPane = MissionPane

	output := m.renderDashboardView()
	// Should contain mission content
	if !strings.Contains(output, "Active") {
		t.Error("expected active mission view to contain body content")
	}
	// Should NOT contain commit pane placeholder
	if strings.Contains(output, "No commit info") {
		t.Error("expected active mission view NOT to contain commit info")
	}

	// Case 2: Completed Mission (Commit Pane)
	m.selectedMission = &mission.Mission{ID: "completed", Status: "completed", Body: "## INTENT\nDone"}
	m.currentPane = CommitPane
	m.commitLoaded = true
	m.commitMessage = "test commit"

	output = m.renderDashboardView()
	if !strings.Contains(output, "test commit") {
		t.Error("expected completed mission view to contain commit message in CommitPane")
	}
}

func TestRenderTwoPaneLayout(t *testing.T) {
	m := NewDashboardModel()
	m.width = 100
	m.calculatePaneWidths()

	m.selectedMission = &mission.Mission{ID: "test", Status: "active", Body: "## INTENT\nBody"}
	m.currentPane = MissionPane
	m.executionLogLoaded = true
	m.executionLog = "LogContent"

	// MissionPane active
	output := m.renderTwoPaneLayout()
	if !strings.Contains(output, "Body") {
		t.Error("expected renderTwoPaneLayout to contain mission body")
	}
	if !strings.Contains(output, "LogContent") {
		t.Error("expected renderTwoPaneLayout to contain log content")
	}
}
