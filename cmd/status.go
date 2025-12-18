/*
Copyright © 2025 Yi Xie dnatag@gmail.com

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/dnatag/idd/internal/tui"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display current and completed mission status with interactive TUI",
	Long: `Display the current mission status and browse completed missions using an 
interactive Terminal User Interface (TUI). Shows mission details, progress, 
and provides clear guidance on next steps.

Use Tab to switch between current mission and completed missions history.
Use arrow keys to navigate through completed missions.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := tui.RunStatusTUI(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running status TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
