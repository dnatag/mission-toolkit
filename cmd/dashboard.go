/*
Copyright © 2025 Yi Xie dnatag@gmail.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/dnatag/mission-toolkit/internal/tui"
	"github.com/spf13/cobra"
)

// dashboardCmd represents the dashboard command
var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Display comprehensive mission dashboard with execution logs in interactive TUI",
	Long: `Display a comprehensive mission dashboard with execution logs using an 
interactive Terminal User Interface (TUI). Shows mission details, execution progress,
and commit history with adaptive layout for active vs completed missions.

Features:
- Split-pane view: mission.md | execution.log | commit.msg (for completed)
- Live refresh for active mission execution logs
- Lazy loading for completed mission logs and commit messages
- Keyboard navigation between panes (Tab/Shift+Tab)

Use ↑/↓ to navigate missions, Enter to view details, / to search.
Use ←/→ for pagination, Tab to switch panes, Esc to go back, q to quit.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := tui.RunDashboardTUI(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running dashboard TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}
