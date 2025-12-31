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
		validationService := mission.NewValidationService(afero.NewOsFs(), ".mission")

		status, err := validationService.CheckMissionState()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Output JSON for AI consumption
		jsonOutput, err := status.ToJSON()
		if err != nil {
			fmt.Printf("Error formatting output: %v\n", err)
			return
		}

		fmt.Println(jsonOutput)
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
	planCmd.AddCommand(planCheckCmd)
}
