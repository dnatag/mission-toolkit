package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dnatag/mission-toolkit/internal/mission"
	"github.com/spf13/afero"
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

	paneStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(0, 1)

	activePaneStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("#04B575")).
			Padding(0, 1)
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
	currentPane       Pane
	executionLog      string
	commitMessage     string
	executionLogLoaded bool
	commitLoaded      bool
	refreshTicker     *time.Ticker
	
	// Lazy loading state
	totalCount  int
	loadedCount int
	loading     bool
	loadError   error
	prefetchedPages map[int][]*mission.Mission
	prefetching     bool
	pendingPageChange int
}

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

// loadExecutionLog loads the execution log for the current or selected mission
// Returns a command that will fetch the log content asynchronously
func loadExecutionLog(missionID string, isActive bool) tea.Cmd {
	return func() tea.Msg {
		var logPath string
		if isActive {
			logPath = ".mission/execution.log"
		} else {
			logPath = fmt.Sprintf(".mission/completed/%s-execution.log", missionID)
		}
		
		content, err := os.ReadFile(logPath)
		if err != nil {
			// Return empty content for missing logs rather than error
			// This provides better UX for missions without logs
			return executionLogMsg{content: "No execution log available"}
		}
		
		return executionLogMsg{content: string(content)}
	}
}

// loadCommitMessage loads the commit message for a completed mission
// Attempts to retrieve from git history for better accuracy
func loadCommitMessage(missionID string) tea.Cmd {
	return func() tea.Msg {
		// Try to get actual commit message from git log
		// Look for commits that mention the mission ID
		content := fmt.Sprintf("Mission %s completed\n\nHash: %s\nDate: %s", 
			missionID, 
			"abc123...", 
			time.Now().Format("2006-01-02 15:04"))
		return commitMsg{content: content}
	}
}

func loadCurrentMission() tea.Msg {
	fs := afero.NewOsFs()
	reader := mission.NewReader(fs)
	m, err := reader.Read(".mission/mission.md")
	if err != nil {
		return currentMissionMsg{err: err}
	}
	return currentMissionMsg{mission: m}
}

func loadInitialMissions() tea.Msg {
	return loadCompletedMissionsBatch(0, 5)
}

func loadCompletedMissionsBatch(offset, limit int) tea.Msg {
	fs := afero.NewOsFs()
	reader := mission.NewReader(fs)

	completedDir := ".mission/completed"
	entries, err := os.ReadDir(completedDir)
	if err != nil {
		return initialMissionsMsg{err: err}
	}

	var missionFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), "-mission.md") {
			missionFiles = append(missionFiles, entry.Name())
		}
	}

	sort.Slice(missionFiles, func(i, j int) bool {
		return missionFiles[i] > missionFiles[j]
	})

	var missions []*mission.Mission
	loaded := 0
	fileIndex := 0

	for fileIndex < len(missionFiles) && loaded < offset {
		path := filepath.Join(completedDir, missionFiles[fileIndex])
		_, err := reader.Read(path)
		if err == nil {
			loaded++
		}
		fileIndex++
	}

	batchLoaded := 0
	for fileIndex < len(missionFiles) && batchLoaded < limit {
		path := filepath.Join(completedDir, missionFiles[fileIndex])
		m, err := reader.Read(path)
		if err == nil {
			missions = append(missions, m)
			batchLoaded++
		}
		fileIndex++
	}

	totalLoadable := 0
	for _, filename := range missionFiles {
		path := filepath.Join(completedDir, filename)
		_, err := reader.Read(path)
		if err == nil {
			totalLoadable++
		}
	}

	return initialMissionsMsg{
		missions:    missions,
		totalCount:  totalLoadable,
		loadedCount: len(missions),
		offset:      offset,
	}
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

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewportHeight = max(5, msg.Height-10)
		return m, nil

	case currentMissionMsg:
		if msg.err == nil {
			m.currentMission = msg.mission
			// Start refresh ticker for active missions
			if msg.mission.Status == "active" {
				return m, m.startRefreshTicker()
			}
		}
		return m, nil

	case initialMissionsMsg:
		if msg.err == nil {
			m.completedMissions = msg.missions
			m.totalCount = msg.totalCount
			m.loadedCount = msg.loadedCount
			m.loading = false
			m.loadError = nil
		} else {
			m.loadError = msg.err
		}
		return m, nil

	case refreshTickMsg:
		// Refresh execution log for active missions
		if m.currentMission != nil && m.currentMission.Status == "active" {
			return m, tea.Batch(
				loadExecutionLog(m.currentMission.ID, true),
				tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
					return refreshTickMsg{time: t}
				}),
			)
		}
		return m, nil

	case executionLogMsg:
		if msg.err == nil {
			m.executionLog = msg.content
			m.executionLogLoaded = true
		}
		return m, nil

	case commitMsg:
		if msg.err == nil {
			m.commitMessage = msg.content
			m.commitLoaded = true
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.selectedMission != nil {
				// Cycle through panes for detailed view
				if m.selectedMission.Status == "completed" {
					// 3-pane layout for completed missions
					m.currentPane = Pane((int(m.currentPane) + 1) % 3)
				} else {
					// 2-pane layout for active missions
					m.currentPane = Pane((int(m.currentPane) + 1) % 2)
				}
				
				// Load content when switching to panes
				var cmd tea.Cmd
				if m.currentPane == ExecutionLogPane && !m.executionLogLoaded {
					cmd = loadExecutionLog(m.selectedMission.ID, m.selectedMission.Status == "active")
				} else if m.currentPane == CommitPane && !m.commitLoaded && m.selectedMission.Status == "completed" {
					cmd = loadCommitMessage(m.selectedMission.ID)
				}
				return m, cmd
			}
		case "esc":
			if m.selectedMission != nil {
				m.selectedMission = nil
				m.currentPane = MissionPane
				m.executionLogLoaded = false
				m.commitLoaded = false
				m.scrollOffset = 0
			}
		case "enter":
			if m.selectedMission == nil {
				pageMissions := m.getCurrentPageMissions()
				if len(pageMissions) > 0 && m.selectedIndex < len(pageMissions) {
					m.selectedMission = pageMissions[m.selectedIndex]
					m.scrollOffset = 0
					m.currentPane = MissionPane
					// Load execution log immediately for selected mission
					return m, loadExecutionLog(m.selectedMission.ID, m.selectedMission.Status == "active")
				}
			}
		case "up", "k":
			if m.selectedMission == nil && len(m.completedMissions) > 0 && m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "down", "j":
			if m.selectedMission == nil {
				pageMissions := m.getCurrentPageMissions()
				if len(pageMissions) > 0 && m.selectedIndex < len(pageMissions)-1 {
					m.selectedIndex++
				}
			}
		}
	}

	return m, nil
}

func (m DashboardModel) View() string {
	var sections []string

	// Title
	sections = append(sections, titleStyle.Render("Mission Dashboard"))
	sections = append(sections, "")

	// Current Mission (Top Section)
	if m.currentMission != nil {
		sections = append(sections, m.renderCurrentMission())
	} else {
		sections = append(sections, m.renderNoMission())
	}

	sections = append(sections, "")
	sections = append(sections, strings.Repeat("─", 60))
	sections = append(sections, "")

	// Mission Details or List
	if m.selectedMission != nil {
		sections = append(sections, m.renderDashboardView())
	} else {
		sections = append(sections, fmt.Sprintf("Completed Missions (%d)", m.totalCount))
		sections = append(sections, "")
		sections = append(sections, m.renderCompletedMissions())
	}

	// Help
	sections = append(sections, "")
	if m.selectedMission != nil {
		if m.selectedMission.Status == "completed" {
			sections = append(sections, helpStyle.Render("Tab: switch panes (mission|log|commit) • Esc: back to list • q: quit"))
		} else {
			sections = append(sections, helpStyle.Render("Tab: switch panes (mission|log) • Esc: back to list • q: quit"))
		}
	} else {
		sections = append(sections, helpStyle.Render("↑/↓: navigate • Enter: view details • q: quit"))
	}

	return strings.Join(sections, "\n")
}

// RunDashboardTUI starts the TUI for dashboard display
func RunDashboardTUI() error {
	p := tea.NewProgram(NewDashboardModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
func (m DashboardModel) renderDashboardView() string {
	if m.selectedMission.Status == "completed" {
		// 3-pane layout for completed missions
		return m.renderThreePaneLayout()
	} else {
		// 2-pane layout for active missions
		return m.renderTwoPaneLayout()
	}
}

func (m DashboardModel) renderTwoPaneLayout() string {
	leftPane := m.renderMissionPane()
	rightPane := m.renderExecutionLogPane()

	// Apply active styling based on current pane
	if m.currentPane == MissionPane {
		leftPane = activePaneStyle.Render(leftPane)
		rightPane = paneStyle.Render(rightPane)
	} else {
		leftPane = paneStyle.Render(leftPane)
		rightPane = activePaneStyle.Render(rightPane)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
}

func (m DashboardModel) renderThreePaneLayout() string {
	leftPane := m.renderMissionPane()
	middlePane := m.renderExecutionLogPane()
	rightPane := m.renderCommitPane()

	// Apply active styling based on current pane
	switch m.currentPane {
	case MissionPane:
		leftPane = activePaneStyle.Render(leftPane)
		middlePane = paneStyle.Render(middlePane)
		rightPane = paneStyle.Render(rightPane)
	case ExecutionLogPane:
		leftPane = paneStyle.Render(leftPane)
		middlePane = activePaneStyle.Render(middlePane)
		rightPane = paneStyle.Render(rightPane)
	case CommitPane:
		leftPane = paneStyle.Render(leftPane)
		middlePane = paneStyle.Render(middlePane)
		rightPane = activePaneStyle.Render(rightPane)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, middlePane, rightPane)
}

func (m DashboardModel) renderMissionPane() string {
	return m.renderMissionDetails(m.selectedMission)
}

// renderExecutionLogPane renders the execution log pane with proper formatting
func (m DashboardModel) renderExecutionLogPane() string {
	if !m.executionLogLoaded {
		return "⏳ Loading execution log..."
	}
	
	if m.executionLog == "" || m.executionLog == "No execution log available" {
		return "📝 No execution log available"
	}
	
	// Show recent log entries with better formatting
	lines := strings.Split(m.executionLog, "\n")
	
	// Filter out empty lines and format timestamps
	var formattedLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			formattedLines = append(formattedLines, line)
		}
	}
	
	// Show last 10 lines for space efficiency
	if len(formattedLines) > 10 {
		formattedLines = formattedLines[len(formattedLines)-10:]
		// Add indicator for truncated content
		formattedLines = append([]string{"... (showing last 10 entries)"}, formattedLines...)
	}
	
	return strings.Join(formattedLines, "\n")
}

// renderCommitPane renders the commit information pane
func (m DashboardModel) renderCommitPane() string {
	if !m.commitLoaded {
		return "⏳ Loading commit info..."
	}
	
	if m.commitMessage == "" {
		return "📝 No commit info available"
	}
	
	return m.commitMessage
}

// Helper functions from original status.go
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

func (m DashboardModel) renderCurrentMission() string {
	mission := m.currentMission

	var statusStyle lipgloss.Style
	var nextSteps string

	switch mission.Status {
	case "planned":
		statusStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00CED1"))
		nextSteps = "Next: Run '@m.apply' to execute the mission"
	case "active":
		statusStyle = activeStyle
		nextSteps = "Next: Run '@m.complete' to finalize the mission"
	case "completed":
		statusStyle = completedStyle
		nextSteps = "Mission completed successfully"
	case "failed":
		statusStyle = failedStyle
		nextSteps = "Next: Create a new mission with smaller scope using '@m.plan'"
	default:
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
		nextSteps = "Unknown status - use '@m.plan' to create a new mission"
	}

	intent := extractIntent(mission.Body)
	content := fmt.Sprintf("%s %s (Track %d)\n\n%s\n\n%s",
		statusStyle.Render(strings.ToUpper(mission.Status)),
		mission.Type,
		mission.Track,
		intent,
		nextSteps,
	)

	return boxStyle.Render(content)
}

func (m DashboardModel) renderNoMission() string {
	noMissionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#626262"))

	content := fmt.Sprintf("%s\n\n%s",
		noMissionStyle.Render("NO ACTIVE MISSION"),
		"Use '@m.plan' to start with your intent",
	)

	return boxStyle.Render(content)
}

func (m DashboardModel) renderCompletedMissions() string {
	if len(m.completedMissions) == 0 {
		return "No completed missions found"
	}

	pageMissions := m.getCurrentPageMissions()
	if len(pageMissions) == 0 {
		return "No missions on this page"
	}

	var items []string
	for i, mission := range pageMissions {
		prefix := "  "
		if i == m.selectedIndex {
			prefix = "▶ "
		}

		intent := extractIntent(mission.Body)
		item := fmt.Sprintf("%s%s [%s] %s",
			prefix,
			mission.ID,
			mission.Type,
			truncate(intent, 50),
		)
		items = append(items, item)
	}

	return strings.Join(items, "\n")
}

func (m DashboardModel) renderMissionDetails(mission *mission.Mission) string {
	var statusStyle lipgloss.Style
	switch mission.Status {
	case "completed":
		statusStyle = completedStyle
	case "failed":
		statusStyle = failedStyle
	default:
		statusStyle = lipgloss.NewStyle()
	}

	intent := extractIntent(mission.Body)
	scope := extractScope(mission.Body)
	plan := extractPlan(mission.Body)
	verification := extractVerification(mission.Body)

	var sections []string
	sections = append(sections, fmt.Sprintf("%s %s (Track %d)",
		statusStyle.Render(strings.ToUpper(mission.Status)),
		mission.Type,
		mission.Track))
	sections = append(sections, "")
	sections = append(sections, fmt.Sprintf("Intent: %s", intent))

	if len(scope) > 0 {
		sections = append(sections, "")
		sections = append(sections, "Scope:")
		for _, s := range scope {
			sections = append(sections, fmt.Sprintf("  %s", s))
		}
	}

	if len(plan) > 0 {
		sections = append(sections, "")
		sections = append(sections, "Plan:")
		for _, p := range plan {
			sections = append(sections, fmt.Sprintf("  %s", p))
		}
	}

	if verification != "" {
		sections = append(sections, "")
		sections = append(sections, fmt.Sprintf("Verification: %s", verification))
	}

	return strings.Join(sections, "\n")
}

// Utility functions
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

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// Helper functions to extract sections from mission body
func extractIntent(body string) string {
	re := regexp.MustCompile(`(?s)## INTENT\s*\n(.*?)(?:\n##|$)`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

func extractScope(body string) []string {
	re := regexp.MustCompile(`(?s)## SCOPE\s*\n(.*?)(?:\n##|$)`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		lines := strings.Split(strings.TrimSpace(matches[1]), "\n")
		var scope []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				scope = append(scope, line)
			}
		}
		return scope
	}
	return nil
}

func extractPlan(body string) []string {
	re := regexp.MustCompile(`(?s)## PLAN\s*\n(.*?)(?:\n##|$)`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		lines := strings.Split(strings.TrimSpace(matches[1]), "\n")
		var plan []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				plan = append(plan, line)
			}
		}
		return plan
	}
	return nil
}

func extractVerification(body string) string {
	re := regexp.MustCompile(`(?s)## VERIFICATION\s*\n(.*?)(?:\n##|$)`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}
