package tui

import (
	"github.com/charmbracelet/bubbletea"
)

var styles = NewStyles()

// RunDashboardTUI starts the TUI for dashboard display
func RunDashboardTUI() error {
	p := tea.NewProgram(NewDashboardModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
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
