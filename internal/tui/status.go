package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dnatag/idd/internal/mission"
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	activeStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#04B575"))

	completedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262"))

	failedStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF5F87"))

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262"))

	boxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 2)
)

type Model struct {
	currentMission    *mission.Mission
	completedMissions []*mission.Mission
	selectedIndex     int
	showCompleted     bool
	width             int
	height            int
}

func NewModel() Model {
	return Model{
		selectedIndex: 0,
		showCompleted: false,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		loadCurrentMission,
		loadCompletedMissions,
	)
}

func loadCurrentMission() tea.Msg {
	mission, err := mission.ReadCurrentMission()
	if err != nil {
		return currentMissionMsg{err: err}
	}
	return currentMissionMsg{mission: mission}
}

func loadCompletedMissions() tea.Msg {
	missions, err := mission.ReadCompletedMissions()
	if err != nil {
		return completedMissionsMsg{err: err}
	}
	return completedMissionsMsg{missions: missions}
}

type currentMissionMsg struct {
	mission *mission.Mission
	err     error
}

type completedMissionsMsg struct {
	missions []*mission.Mission
	err      error
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case currentMissionMsg:
		if msg.err == nil {
			m.currentMission = msg.mission
		}
		return m, nil

	case completedMissionsMsg:
		if msg.err == nil {
			m.completedMissions = msg.missions
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			m.showCompleted = !m.showCompleted
			m.selectedIndex = 0
		case "up", "k":
			if m.showCompleted && len(m.completedMissions) > 0 {
				m.selectedIndex = max(0, m.selectedIndex-1)
			}
		case "down", "j":
			if m.showCompleted && len(m.completedMissions) > 0 {
				m.selectedIndex = min(len(m.completedMissions)-1, m.selectedIndex+1)
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	var sections []string

	// Title
	sections = append(sections, titleStyle.Render("IDD Mission Status"))
	sections = append(sections, "")

	// Current Mission
	if m.currentMission != nil {
		sections = append(sections, m.renderCurrentMission())
	} else {
		sections = append(sections, "No active mission found")
	}

	sections = append(sections, "")

	// Tab indicator
	tabIndicator := "Current Mission"
	if m.showCompleted {
		tabIndicator = fmt.Sprintf("Completed Missions (%d)", len(m.completedMissions))
	}
	sections = append(sections, fmt.Sprintf("View: %s", tabIndicator))
	sections = append(sections, "")

	// Completed missions (if tab is active)
	if m.showCompleted {
		sections = append(sections, m.renderCompletedMissions())
	}

	// Help
	sections = append(sections, "")
	sections = append(sections, helpStyle.Render("Tab: switch view • ↑/↓: navigate • q: quit"))

	return strings.Join(sections, "\n")
}

func (m Model) renderCurrentMission() string {
	mission := m.currentMission
	
	var statusStyle lipgloss.Style
	var nextSteps string

	switch mission.Status {
	case "active":
		statusStyle = activeStyle
		nextSteps = "Next: Run 'idd apply' to execute the mission"
	case "completed":
		statusStyle = completedStyle
		nextSteps = "Mission completed successfully"
	case "failed":
		statusStyle = failedStyle
		nextSteps = "Next: Create a new mission with smaller scope"
	default:
		statusStyle = lipgloss.NewStyle()
		nextSteps = "Unknown status"
	}

	content := fmt.Sprintf("%s %s (Track %s)\n\n%s\n\n%s",
		statusStyle.Render(strings.ToUpper(mission.Status)),
		mission.Type,
		mission.Track,
		mission.Intent,
		nextSteps,
	)

	return boxStyle.Render(content)
}

func (m Model) renderCompletedMissions() string {
	if len(m.completedMissions) == 0 {
		return "No completed missions found"
	}

	var items []string
	for i, mission := range m.completedMissions {
		prefix := "  "
		if i == m.selectedIndex {
			prefix = "▶ "
		}

		timeStr := "Unknown"
		if mission.CompletedAt != nil {
			timeStr = mission.CompletedAt.Format("2006-01-02 15:04")
		}

		item := fmt.Sprintf("%s%s [%s] %s",
			prefix,
			timeStr,
			mission.Type,
			truncate(mission.Intent, 50),
		)
		items = append(items, item)
	}

	return strings.Join(items, "\n")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// RunStatusTUI starts the TUI for status display
func RunStatusTUI() error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
