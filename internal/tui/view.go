package tui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dnatag/mission-toolkit/internal/mission"
)

// View renders the dashboard UI
func (m DashboardModel) View() string {
	var sections []string

	sections = append(sections, titleStyle.Render("Mission Dashboard"))
	sections = append(sections, "")

	if m.currentMission != nil {
		sections = append(sections, m.renderCurrentMission())
	} else {
		sections = append(sections, m.renderNoMission())
	}

	sections = append(sections, "")
	sections = append(sections, strings.Repeat("─", 60))
	sections = append(sections, "")

	if m.selectedMission != nil {
		sections = append(sections, m.renderDashboardView())
	} else {
		sections = append(sections, fmt.Sprintf("Completed Missions (%d)", m.totalCount))
		sections = append(sections, "")
		sections = append(sections, m.renderCompletedMissions())
	}

	sections = append(sections, "")
	if m.selectedMission != nil {
		if m.selectedMission.Status == "completed" {
			sections = append(sections, helpStyle.Render("Tab: switch panes (mission|log|commit) • Esc: back to list • q: quit"))
		} else {
			sections = append(sections, helpStyle.Render("Tab: switch panes (mission|log) • Esc: back to list • q: quit"))
		}
	} else {
		sections = append(sections, helpStyle.Render("↑/↓: navigate • ←/→: prev/next page • Enter: view details • q: quit"))
	}

	return strings.Join(sections, "\n")
}

// renderDashboardView renders the appropriate pane layout
func (m DashboardModel) renderDashboardView() string {
	// Always use two-pane layout, right pane content depends on currentPane
	return m.renderTwoPaneLayout()
}

// renderTwoPaneLayout renders a two-pane layout
func (m DashboardModel) renderTwoPaneLayout() string {
	leftPane := m.renderMissionPane()
	var rightPane string

	switch m.currentPane {
	case CommitPane:
		if m.selectedMission.Status == "completed" {
			rightPane = m.renderCommitPane()
		} else {
			// Fallback to log if not completed
			rightPane = m.renderExecutionLogPane()
		}
	case ExecutionLogPane:
		rightPane = m.renderExecutionLogPane()
	default:
		// Default to execution log
		rightPane = m.renderExecutionLogPane()
	}

	// Apply width constraints
	if m.leftPaneWidth > 0 {
		leftPane = lipgloss.NewStyle().Width(m.leftPaneWidth).Render(leftPane)
	}
	if m.rightPaneWidth > 0 {
		rightPane = lipgloss.NewStyle().Width(m.rightPaneWidth).Render(rightPane)
	}

	if m.currentPane == MissionPane {
		leftPane = activePaneStyle.Render(leftPane)
		rightPane = paneStyle.Render(rightPane)
	} else {
		leftPane = paneStyle.Render(leftPane)
		rightPane = activePaneStyle.Render(rightPane)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
}

// renderMissionPane renders the mission details pane
func (m DashboardModel) renderMissionPane() string {
	return m.renderMissionDetails(m.selectedMission)
}

// renderExecutionLogPane renders the execution log pane with recent entries
func (m DashboardModel) renderExecutionLogPane() string {
	if !m.executionLogLoaded {
		return "⏳ Loading execution log..."
	}

	if m.executionLog == "" || m.executionLog == "No execution log available" {
		return "📝 No execution log available"
	}

	lines := strings.Split(m.executionLog, "\n")

	var formattedLines []string
	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			formattedLines = append(formattedLines, line)
		}
	}

	const maxLines = 10
	if len(formattedLines) > maxLines {
		formattedLines = formattedLines[len(formattedLines)-maxLines:]
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

// renderCurrentMission renders the current mission box
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

// renderNoMission renders the no mission box
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

// renderCompletedMissions renders the list of completed missions
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

// renderMissionDetails renders detailed mission information
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
