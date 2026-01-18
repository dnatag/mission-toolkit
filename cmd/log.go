package cmd

import (
	"fmt"

	"github.com/dnatag/mission-toolkit/pkg/logger"
	"github.com/dnatag/mission-toolkit/pkg/mission"
	"github.com/spf13/afero"
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
		file, _ := cmd.Flags().GetString("file")
		message := args[0]

		// Get mission ID from centralized service
		idService := mission.NewIDService(afero.NewOsFs(), ".mission/mission.md")
		missionID, err := idService.GetCurrentID()
		if err != nil {
			fmt.Printf("Warning: Could not get mission ID: %v\n", err)
			missionID = "unknown"
		}

		// Create logger config based on file flag
		config := logger.DefaultConfig()
		if file == "" {
			config.Output = logger.OutputConsole
		} else {
			config.Output = logger.OutputBoth
			config.FilePath = file
		}

		// Create logger with custom config
		log := logger.NewWithConfig(missionID, config)

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
	logCmd.Flags().StringP("file", "f", ".mission/execution.log", "Log file path (empty string for console only)")
}
