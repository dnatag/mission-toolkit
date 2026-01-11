package tui

import (
	"testing"

	"github.com/charmbracelet/bubbletea"
)

func TestUpdate_WindowSize(t *testing.T) {
	m := NewDashboardModel()
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}

	updated, _ := m.Update(msg)
	model := updated.(DashboardModel)

	if model.width != 120 {
		t.Errorf("expected width 120, got %d", model.width)
	}
	if model.height != 40 {
		t.Errorf("expected height 40, got %d", model.height)
	}
	if model.leftPaneWidth != 48 {
		t.Errorf("expected leftPaneWidth 48, got %d", model.leftPaneWidth)
	}
}

func TestUpdate_KeyQuit(t *testing.T) {
	m := NewDashboardModel()
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}

	_, cmd := m.Update(msg)
	if cmd == nil {
		t.Error("expected quit command")
	}
}

func TestUpdate_CurrentMissionMsg(t *testing.T) {
	m := NewDashboardModel()

	// Test with nil mission and no error
	msg := currentMissionMsg{
		mission: nil,
		err:     nil,
	}

	updated, _ := m.Update(msg)
	model := updated.(DashboardModel)

	if model.currentMission != nil {
		t.Error("expected currentMission to be nil")
	}
}

func TestUpdate_ExecutionLogMsg(t *testing.T) {
	m := NewDashboardModel()
	msg := executionLogMsg{
		content: "test log content",
		err:     nil,
	}

	updated, _ := m.Update(msg)
	model := updated.(DashboardModel)

	if model.executionLog != "test log content" {
		t.Errorf("expected executionLog 'test log content', got '%s'", model.executionLog)
	}
	if !model.executionLogLoaded {
		t.Error("expected executionLogLoaded to be true")
	}
}
