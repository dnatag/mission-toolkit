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
