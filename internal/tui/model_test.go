package tui

import (
	"testing"

	"github.com/dnatag/mission-toolkit/internal/mission"
)

func TestCalculatePaneWidths(t *testing.T) {
	m := NewDashboardModel()
	m.width = 100
	m.calculatePaneWidths()

	// 80% of 100 is 80. Split 50/50 is 40 each.
	if m.leftPaneWidth != 40 {
		t.Errorf("expected leftPaneWidth 40, got %d", m.leftPaneWidth)
	}
	if m.rightPaneWidth != 40 {
		t.Errorf("expected rightPaneWidth 40, got %d", m.rightPaneWidth)
	}
}

func TestGetCurrentPageMissions(t *testing.T) {
	m := NewDashboardModel()
	m.itemsPerPage = 2

	// Empty missions
	result := m.getCurrentPageMissions()
	if result != nil {
		t.Error("expected nil for empty missions")
	}
}

func TestInit(t *testing.T) {
	model := NewDashboardModel()

	// Init returns a batch command with multiple sub-commands
	cmd := model.Init()
	if cmd == nil {
		t.Fatal("Init should return a non-nil command")
	}
}

func TestStartRefreshTicker(t *testing.T) {
	tests := []struct {
		name           string
		currentMission *mission.Mission
		wantNil        bool
	}{
		{
			name:           "nil mission returns nil",
			currentMission: nil,
			wantNil:        true,
		},
		{
			name:           "active mission returns ticker command",
			currentMission: &mission.Mission{ID: "test", Status: "active"},
			wantNil:        false,
		},
		{
			name:           "planned mission returns nil",
			currentMission: &mission.Mission{ID: "test", Status: "planned"},
			wantNil:        true,
		},
		{
			name:           "completed mission returns nil",
			currentMission: &mission.Mission{ID: "test", Status: "completed"},
			wantNil:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewDashboardModel()
			model.currentMission = tt.currentMission

			cmd := model.startRefreshTicker()

			if tt.wantNil {
				if cmd != nil {
					t.Error("expected nil command, got non-nil")
				}
			} else {
				if cmd == nil {
					t.Error("expected non-nil command for active mission")
				}
			}
		})
	}
}
