package cmd

import (
	"fmt"

	"github.com/dnatag/mission-toolkit/internal/mission"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Mission planning tools",
	Long:  `Mission planning tools for analysis, validation, and generation.`,
}

// planCheckCmd represents the plan check command
var planCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check mission state and cleanup stale artifacts",
	Long:  `Check mission state, generate mission ID, and cleanup stale artifacts.`,
	Run: func(cmd *cobra.Command, args []string) {
		idService := mission.NewIDService(afero.NewOsFs(), ".mission")

		// Cleanup stale ID first
		if err := idService.CleanupStaleID(); err != nil {
			fmt.Printf("Warning: Failed to cleanup stale ID: %v\n", err)
		}

		// Get or create mission ID
		missionID, err := idService.GetOrCreateID()
		if err != nil {
			fmt.Printf("Error: Failed to get mission ID: %v\n", err)
			return
		}

		fmt.Printf("Mission ID: %s\n", missionID)
		fmt.Printf("Status: Ready for planning\n")
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
	planCmd.AddCommand(planCheckCmd)
}
