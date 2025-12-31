package cmd

import (
	"fmt"

	"github.com/dnatag/mission-toolkit/internal/logger"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log [message]",
	Short: "Log messages to mission execution log",
	Long:  `Log messages to the mission execution log with structured formatting for both CLI and AI components.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		level, _ := cmd.Flags().GetString("level")
		step, _ := cmd.Flags().GetString("step")
		missionIDFlag, _ := cmd.Flags().GetString("mission-id")
		message := args[0]

		// Get mission ID - use flag if provided, otherwise auto-detect
		var missionID string
		if missionIDFlag != "" {
			missionID = missionIDFlag
		} else {
			missionID = logger.GetMissionID()
		}

		// Create logger
		log := logger.New(missionID)

		// Log the message
		log.LogStep(level, step, message)

		fmt.Printf("Logged: [%s] %s\n", level, message)
	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	// Add flags
	logCmd.Flags().StringP("level", "l", "INFO", "Log level (DEBUG, INFO, WARN, ERROR, SUCCESS)")
	logCmd.Flags().StringP("step", "s", "General", "Mission step name")
	logCmd.Flags().StringP("mission-id", "m", "", "Mission ID (auto-detected if not provided)")
}
