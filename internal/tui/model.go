package tui

import (
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/dnatag/mission-toolkit/internal/mission"
)

// Pane represents different dashboard panes
type Pane int

const (
	MissionPane Pane = iota
	ExecutionLogPane
	CommitPane
)

// DashboardModel represents the comprehensive dashboard state
type DashboardModel struct {
	currentMission    *mission.Mission
	completedMissions []*mission.Mission
	selectedIndex     int
	selectedMission   *mission.Mission
	scrollOffset      int
	viewportHeight    int
	currentPage       int
	itemsPerPage      int
	searchMode        bool
	searchQuery       string
	filteredMissions  []*mission.Mission
	width             int
	height            int

	// Dashboard-specific state
	currentPane        Pane
	executionLog       string
	commitMessage      string
	executionLogLoaded bool
	commitLoaded       bool
	refreshTicker      *time.Ticker

	// Lazy loading state
	totalCount        int
	loadedCount       int
	loading           bool
	loadError         error
	prefetchedPages   map[int][]*mission.Mission
	prefetching       bool
	pendingPageChange int

	// Split-pane width configuration
	leftPaneWidth   int
	middlePaneWidth int
	rightPaneWidth  int
}

// NewDashboardModel creates a new dashboard model with default settings
func NewDashboardModel() DashboardModel {
	return DashboardModel{
		selectedIndex:     0,
		currentPage:       0,
		itemsPerPage:      5,
		prefetchedPages:   make(map[int][]*mission.Mission),
		pendingPageChange: -1,
		currentPane:       MissionPane,
	}
}

// Message types
type refreshTickMsg struct {
	time time.Time
}

type executionLogMsg struct {
	content string
	err     error
}

type commitMsg struct {
	content string
	err     error
}

type currentMissionMsg struct {
	mission *mission.Mission
	err     error
}

type initialMissionsMsg struct {
	missions    []*mission.Mission
	totalCount  int
	loadedCount int
	offset      int
	err         error
}

// Init initializes the dashboard model
func (m DashboardModel) Init() tea.Cmd {
	return tea.Batch(
		loadCurrentMission,
		loadInitialMissions,
		m.startRefreshTicker(),
	)
}

// startRefreshTicker creates a ticker for live refresh of active mission logs
func (m DashboardModel) startRefreshTicker() tea.Cmd {
	if m.currentMission != nil && m.currentMission.Status == "active" {
		return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
			return refreshTickMsg{time: t}
		})
	}
	return nil
}

// getCurrentPageMissions returns missions for the current page
func (m DashboardModel) getCurrentPageMissions() []*mission.Mission {
	missions := m.completedMissions
	totalItems := len(missions)
	if totalItems == 0 {
		return nil
	}

	start := m.currentPage * m.itemsPerPage
	if start >= totalItems {
		return nil
	}

	end := min(totalItems, start+m.itemsPerPage)
	return missions[start:end]
}

// calculatePaneWidths sets pane widths based on terminal width with default ratios
// Default: 40% left, 30% middle, 30% right
func (m *DashboardModel) calculatePaneWidths() {
	if m.width == 0 {
		return
	}

	m.leftPaneWidth = int(float64(m.width) * 0.4)
	m.middlePaneWidth = int(float64(m.width) * 0.3)
	m.rightPaneWidth = m.width - m.leftPaneWidth - m.middlePaneWidth
}
