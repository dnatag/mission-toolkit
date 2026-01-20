package tui

import "github.com/charmbracelet/lipgloss"

// Styles contains all lipgloss style definitions for the TUI dashboard.
type Styles struct {
	Title      lipgloss.Style
	Active     lipgloss.Style
	Completed  lipgloss.Style
	Failed     lipgloss.Style
	Help       lipgloss.Style
	Box        lipgloss.Style
	Pane       lipgloss.Style
	ActivePane lipgloss.Style
}

// NewStyles creates and returns a new Styles instance with default styling.
func NewStyles() *Styles {
	return &Styles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1),
		Active: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575")),
		Completed: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")),
		Failed: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF5F87")),
		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")),
		Box: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2),
		Pane: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(0, 1),
		ActivePane: lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("#04B575")).
			Padding(0, 1),
	}
}
