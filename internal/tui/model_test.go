package tui

import (
	"testing"
)

func TestCalculatePaneWidths(t *testing.T) {
	m := NewDashboardModel()
	m.width = 100
	m.calculatePaneWidths()

	if m.leftPaneWidth != 40 {
		t.Errorf("expected leftPaneWidth 40, got %d", m.leftPaneWidth)
	}
	if m.middlePaneWidth != 30 {
		t.Errorf("expected middlePaneWidth 30, got %d", m.middlePaneWidth)
	}
	if m.rightPaneWidth != 30 {
		t.Errorf("expected rightPaneWidth 30, got %d", m.rightPaneWidth)
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
