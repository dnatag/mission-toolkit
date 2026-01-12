package tui

import (
	"testing"

	"github.com/charmbracelet/bubbletea"
	"github.com/dnatag/mission-toolkit/internal/mission"
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
	// 80% of 120 is 96. Split 50/50 is 48.
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

func TestUpdate_KeyTab(t *testing.T) {
	m := NewDashboardModel()
	// Need a selected mission to switch panes
	m.selectedMission = &mission.Mission{ID: "test", Status: "active"}

	// Initial state is MissionPane
	if m.currentPane != MissionPane {
		t.Errorf("expected initial pane MissionPane, got %v", m.currentPane)
	}

	// Press Tab -> ExecutionLogPane
	msg := tea.KeyMsg{Type: tea.KeyTab}
	updated, _ := m.Update(msg)
	model := updated.(DashboardModel)

	if model.currentPane != ExecutionLogPane {
		t.Errorf("expected pane ExecutionLogPane, got %v", model.currentPane)
	}

	// Press Tab -> MissionPane (cycle for active)
	updated, _ = model.Update(msg)
	model = updated.(DashboardModel)

	if model.currentPane != MissionPane {
		t.Errorf("expected pane MissionPane, got %v", model.currentPane)
	}
}

func TestUpdate_KeyEnter_Esc(t *testing.T) {
	m := NewDashboardModel()
	// Setup completed missions
	m.completedMissions = []*mission.Mission{
		{ID: "m1", Status: "completed"},
		{ID: "m2", Status: "completed"},
	}
	m.itemsPerPage = 5

	// Select first item (index 0)
	msgEnter := tea.KeyMsg{Type: tea.KeyEnter}
	updated, cmd := m.Update(msgEnter)
	model := updated.(DashboardModel)

	if model.selectedMission == nil {
		t.Error("expected selectedMission to be set")
	} else if model.selectedMission.ID != "m1" {
		t.Errorf("expected selectedMission ID 'm1', got '%s'", model.selectedMission.ID)
	}

	// Verify no command is returned (lazy loading)
	if cmd != nil {
		t.Error("expected no command on Enter (lazy loading)")
	}

	// Verify execution log is not loaded yet
	if model.executionLogLoaded {
		t.Error("expected executionLogLoaded to be false (lazy loading)")
	}

	// Press Esc to deselect
	msgEsc := tea.KeyMsg{Type: tea.KeyEsc}
	updated, _ = model.Update(msgEsc)
	model = updated.(DashboardModel)

	if model.selectedMission != nil {
		t.Error("expected selectedMission to be nil after Esc")
	}

	// Verify flags are reset
	if model.executionLogLoaded {
		t.Error("expected executionLogLoaded to be false after Esc")
	}
	if model.commitLoaded {
		t.Error("expected commitLoaded to be false after Esc")
	}
}

func TestUpdate_KeyNavigation(t *testing.T) {
	m := NewDashboardModel()
	m.completedMissions = []*mission.Mission{
		{ID: "m1"}, {ID: "m2"}, {ID: "m3"},
	}
	m.itemsPerPage = 5

	// Initial index 0
	if m.selectedIndex != 0 {
		t.Errorf("expected initial index 0, got %d", m.selectedIndex)
	}

	// Down -> 1
	msgDown := tea.KeyMsg{Type: tea.KeyDown}
	updated, _ := m.Update(msgDown)
	model := updated.(DashboardModel)
	if model.selectedIndex != 1 {
		t.Errorf("expected index 1, got %d", model.selectedIndex)
	}

	// Down -> 2
	updated, _ = model.Update(msgDown)
	model = updated.(DashboardModel)
	if model.selectedIndex != 2 {
		t.Errorf("expected index 2, got %d", model.selectedIndex)
	}

	// Down -> 2 (bound check)
	updated, _ = model.Update(msgDown)
	model = updated.(DashboardModel)
	if model.selectedIndex != 2 {
		t.Errorf("expected index 2 (bound), got %d", model.selectedIndex)
	}

	// Up -> 1
	msgUp := tea.KeyMsg{Type: tea.KeyUp}
	updated, _ = model.Update(msgUp)
	model = updated.(DashboardModel)
	if model.selectedIndex != 1 {
		t.Errorf("expected index 1, got %d", model.selectedIndex)
	}
}

func TestUpdate_KeyPagination(t *testing.T) {
	m := NewDashboardModel()
	// Setup 12 missions with 5 per page = 3 pages
	m.completedMissions = []*mission.Mission{
		{ID: "m1"}, {ID: "m2"}, {ID: "m3"}, {ID: "m4"}, {ID: "m5"},
		{ID: "m6"}, {ID: "m7"}, {ID: "m8"}, {ID: "m9"}, {ID: "m10"},
		{ID: "m11"}, {ID: "m12"},
	}
	m.itemsPerPage = 5

	// Initial page 0
	if m.currentPage != 0 {
		t.Errorf("expected initial page 0, got %d", m.currentPage)
	}

	// Right -> page 1
	msgRight := tea.KeyMsg{Type: tea.KeyRight}
	updated, _ := m.Update(msgRight)
	model := updated.(DashboardModel)
	if model.currentPage != 1 {
		t.Errorf("expected page 1, got %d", model.currentPage)
	}
	if model.selectedIndex != 0 {
		t.Errorf("expected selectedIndex reset to 0, got %d", model.selectedIndex)
	}

	// Right -> page 2
	updated, _ = model.Update(msgRight)
	model = updated.(DashboardModel)
	if model.currentPage != 2 {
		t.Errorf("expected page 2, got %d", model.currentPage)
	}

	// Right -> page 2 (bound check, can't go beyond last page)
	updated, _ = model.Update(msgRight)
	model = updated.(DashboardModel)
	if model.currentPage != 2 {
		t.Errorf("expected page 2 (bound), got %d", model.currentPage)
	}

	// Left -> page 1
	msgLeft := tea.KeyMsg{Type: tea.KeyLeft}
	updated, _ = model.Update(msgLeft)
	model = updated.(DashboardModel)
	if model.currentPage != 1 {
		t.Errorf("expected page 1, got %d", model.currentPage)
	}
	if model.selectedIndex != 0 {
		t.Errorf("expected selectedIndex reset to 0, got %d", model.selectedIndex)
	}

	// Left -> page 0
	updated, _ = model.Update(msgLeft)
	model = updated.(DashboardModel)
	if model.currentPage != 0 {
		t.Errorf("expected page 0, got %d", model.currentPage)
	}

	// Left -> page 0 (bound check, can't go below 0)
	updated, _ = model.Update(msgLeft)
	model = updated.(DashboardModel)
	if model.currentPage != 0 {
		t.Errorf("expected page 0 (bound), got %d", model.currentPage)
	}
}
