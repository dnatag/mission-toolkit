package tui

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/dnatag/mission-toolkit/internal/mission"
)

const (
	fixedPaneWidth  = 60
	fixedPaneHeight = 12
)

// View renders the dashboard UI
func (m DashboardModel) View() string {
	var sections []string

	sections = append(sections, styles.Title.Render("Mission Dashboard"))
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
			sections = append(sections, styles.Help.Render("Tab: switch panes (mission|log|commit) • ↑↓←→: scroll content • Esc: back to list • q: quit"))
		} else {
			sections = append(sections, styles.Help.Render("Tab: switch panes (mission|log) • ↑↓←→: scroll content • Esc: back to list • q: quit"))
		}
	} else {
		sections = append(sections, styles.Help.Render("↑/↓: navigate • ←/→: prev/next page • Enter: view details • q: quit"))
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

	// Apply scrollable content with proper pane identification
	leftPane = m.applyScrollableContent(leftPane, true)    // Left pane
	rightPane = m.applyScrollableContent(rightPane, false) // Right pane

	if m.currentPane == MissionPane {
		leftPane = styles.ActivePane.Render(leftPane)
		rightPane = styles.Pane.Render(rightPane)
	} else {
		leftPane = styles.Pane.Render(leftPane)
		rightPane = styles.ActivePane.Render(rightPane)
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
		statusStyle = styles.Active
		nextSteps = "Next: Run '@m.complete' to finalize the mission"
	case "completed":
		statusStyle = styles.Completed
		nextSteps = "Mission completed successfully"
	case "failed":
		statusStyle = styles.Failed
		nextSteps = "Next: Create a new mission with smaller scope using '@m.plan'"
	default:
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
		nextSteps = "Unknown status - use '@m.plan' to create a new mission"
	}

	content := fmt.Sprintf("%s %s (Track %d)\n\n%s\n\n%s",
		statusStyle.Render(strings.ToUpper(mission.Status)),
		mission.Type,
		mission.Track,
		mission.GetIntent(),
		nextSteps,
	)

	return styles.Box.Render(content)
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

	return styles.Box.Render(content)
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

		item := fmt.Sprintf("%s%s [%s] %s",
			prefix,
			mission.ID,
			mission.Type,
			truncate(mission.GetIntent(), 50),
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
		statusStyle = styles.Completed
	case "failed":
		statusStyle = styles.Failed
	default:
		statusStyle = lipgloss.NewStyle()
	}

	intent := mission.GetIntent()
	scope := mission.GetScope()
	plan := mission.GetPlan()
	verification := mission.GetVerification()

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

// applyFixedDimensions applies fixed width and height to pane content with scrolling support.
// Content is clipped based on scroll position and scrollbars are added if needed.
func (m DashboardModel) applyFixedDimensions(content string) string {
	// For backward compatibility, use the original fixed dimensions without scrolling
	lines := strings.Split(content, "\n")

	// Ensure exactly fixedPaneHeight lines
	if len(lines) > fixedPaneHeight {
		lines = lines[:fixedPaneHeight]
	}
	for len(lines) < fixedPaneHeight {
		lines = append(lines, "")
	}

	// Ensure each line is exactly fixedPaneWidth characters
	for i, line := range lines {
		lineLen := len(line)
		if lineLen > fixedPaneWidth {
			// Truncate with ellipsis
			lines[i] = line[:fixedPaneWidth-3] + "..."
		} else if lineLen < fixedPaneWidth {
			// Pad with spaces
			lines[i] = line + strings.Repeat(" ", fixedPaneWidth-lineLen)
		}
	}

	return strings.Join(lines, "\n")
}

// clipContentToViewport extracts the visible portion of content based on scroll position.
// It handles both horizontal and vertical clipping, padding the result to fill the viewport.
// Returns the clipped lines padded to the visible dimensions.
func clipContentToViewport(lines []string, scrollX, scrollY, visibleWidth, visibleHeight int) []string {
	var visibleLines []string
	
	// Extract visible lines based on vertical scroll
	for i := scrollY; i < scrollY+visibleHeight && i < len(lines); i++ {
		if i < 0 {
			continue
		}
		
		line := lines[i]
		
		// Apply horizontal clipping
		if scrollX >= len(line) {
			visibleLines = append(visibleLines, "")
			continue
		}
		
		endX := scrollX + visibleWidth
		if endX > len(line) {
			endX = len(line)
		}
		visibleLines = append(visibleLines, line[scrollX:endX])
	}

	// Pad to visible height
	for len(visibleLines) < visibleHeight {
		visibleLines = append(visibleLines, "")
	}

	// Pad each line to visible width
	for i, line := range visibleLines {
		lineWidth := displayWidth(line)
		if lineWidth < visibleWidth {
			visibleLines[i] = line + strings.Repeat(" ", visibleWidth-lineWidth)
		}
	}

	return visibleLines
}

// renderVerticalScrollbar adds vertical scrollbar indicators to each line.
// The scrollbar shows navigation arrows and a position indicator when content exceeds viewport.
// Returns the lines with scrollbar characters appended.
func renderVerticalScrollbar(lines []string, scrollY, totalLines, visibleHeight int) []string {
	needsScrollbar := totalLines > visibleHeight
	result := make([]string, len(lines))
	
	for i := range lines {
		var scrollbarChar string
		
		if !needsScrollbar {
			scrollbarChar = " "
		} else if i == 0 && scrollY > 0 {
			scrollbarChar = "▲" // Up arrow - can scroll up
		} else if i == len(lines)-1 && scrollY+visibleHeight < totalLines {
			scrollbarChar = "▼" // Down arrow - can scroll down
		} else if totalLines > visibleHeight {
			// Show position indicator
			scrollPosition := float64(scrollY) / float64(totalLines-visibleHeight)
			indicatorPos := int(scrollPosition*float64(visibleHeight-2)) + 1
			if i == indicatorPos {
				scrollbarChar = "█" // Position indicator
			} else {
				scrollbarChar = "│" // Scrollbar track
			}
		} else {
			scrollbarChar = " "
		}
		
		result[i] = lines[i] + scrollbarChar
	}
	
	return result
}

// renderHorizontalScrollbar creates the horizontal scrollbar line.
// Shows navigation arrows and position indicator when content width exceeds viewport.
// Returns a string representing the scrollbar with position indicators.
func renderHorizontalScrollbar(scrollX, maxWidth, visibleWidth, totalWidth int) string {
	needsScrollbar := maxWidth > visibleWidth
	
	if !needsScrollbar {
		return strings.Repeat(" ", totalWidth)
	}
	
	var scrollbar strings.Builder
	scrollbar.Grow(totalWidth)
	
	for i := 0; i < visibleWidth; i++ {
		var char string
		
		if i == 0 && scrollX > 0 {
			char = "◀" // Left arrow - can scroll left
		} else if i == visibleWidth-1 && scrollX+visibleWidth < maxWidth {
			char = "▶" // Right arrow - can scroll right
		} else {
			// Show position indicator
			scrollPosition := float64(scrollX) / float64(maxWidth-visibleWidth)
			indicatorPos := int(scrollPosition*float64(visibleWidth-2)) + 1
			if i == indicatorPos {
				char = "█" // Position indicator
			} else {
				char = "─" // Scrollbar track
			}
		}
		
		scrollbar.WriteString(char)
	}
	
	scrollbar.WriteString(" ") // Corner space for vertical scrollbar
	return scrollbar.String()
}

// applyScrollableContent handles content clipping and scrollbar rendering for panes.
// It supports both horizontal and vertical scrolling with visual indicators.
// Parameters:
//   - content: The raw content to be displayed
//   - isLeftPane: true for left pane (mission), false for right pane (log/commit)
//
// Returns: Formatted content with scrollbars and proper clipping
func (m *DashboardModel) applyScrollableContent(content string, isLeftPane bool) string {
	lines := strings.Split(content, "\n")

	// Calculate content dimensions
	maxWidth := 0
	for _, line := range lines {
		if displayWidth(line) > maxWidth {
			maxWidth = displayWidth(line)
		}
	}

	// Update max dimensions for scrolling
	if isLeftPane {
		m.leftPaneMaxWidth = maxWidth
		m.leftPaneMaxHeight = len(lines)
	} else {
		m.rightPaneMaxWidth = maxWidth
		m.rightPaneMaxHeight = len(lines)
	}

	// Get scroll position
	scrollX := m.leftPaneScrollX
	scrollY := m.leftPaneScrollY
	if !isLeftPane {
		scrollX = m.rightPaneScrollX
		scrollY = m.rightPaneScrollY
	}

	// Calculate visible area (reserve space for scrollbars)
	visibleWidth := fixedPaneWidth - 1   // Reserve 1 char for vertical scrollbar
	visibleHeight := fixedPaneHeight - 1 // Reserve 1 line for horizontal scrollbar

	// Clip content based on scroll position
	visibleLines := clipContentToViewport(lines, scrollX, scrollY, visibleWidth, visibleHeight)

	// Add vertical scrollbar with position indicator
	visibleLines = renderVerticalScrollbar(visibleLines, scrollY, len(lines), visibleHeight)

	// Add horizontal scrollbar with position indicator
	horizontalScrollbar := renderHorizontalScrollbar(scrollX, maxWidth, visibleWidth, fixedPaneWidth)
	visibleLines = append(visibleLines, horizontalScrollbar)

	return strings.Join(visibleLines, "\n")
}

// displayWidth calculates the visual width of a string, accounting for emojis.
// Emojis typically render as 2 character widths in most terminals, which can
// cause alignment issues if not accounted for in layout calculations.
func displayWidth(s string) int {
	// Comprehensive emoji detection covering major Unicode emoji blocks
	emojiPattern := regexp.MustCompile(`[\x{1F600}-\x{1F64F}]|[\x{1F300}-\x{1F5FF}]|[\x{1F680}-\x{1F6FF}]|[\x{1F1E0}-\x{1F1FF}]|[\x{2600}-\x{26FF}]|[\x{2700}-\x{27BF}]|[\x{1F900}-\x{1F9FF}]|[\x{1F018}-\x{1F270}]`)

	// Count regular characters (runes, not bytes)
	width := utf8.RuneCountInString(s)

	// Add extra width for emojis (they typically display as 2 characters wide)
	emojiMatches := emojiPattern.FindAllString(s, -1)
	width += len(emojiMatches)

	return width
}
