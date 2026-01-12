package tui

import (
	"testing"
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
