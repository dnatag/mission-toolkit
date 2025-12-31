package cmd

import (
	"fmt"

	"github.com/dnatag/mission-toolkit/internal/mission"
	"github.com/dnatag/mission-toolkit/internal/plan"
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

// planAnalyzeCmd represents the plan analyze command
var planAnalyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze plan complexity and provide recommendations",
	Long:  `Analyze plan.json files to calculate complexity track based on file counts and domain multipliers.`,
	Run: func(cmd *cobra.Command, args []string) {
		planFile, _ := cmd.Flags().GetString("file")
		if planFile == "" {
			fmt.Println("Error: --file flag is required")
			return
		}

		// Get mission ID for logging
		idService := mission.NewIDService(afero.NewOsFs(), ".mission")
		missionID, err := idService.GetCurrentID()
		if err != nil {
			missionID = "unknown"
		}

		// Create analyzer
		fs := afero.NewOsFs()
		analyzer := plan.NewAnalyzer(fs, missionID)

		// Analyze plan
		result, err := analyzer.AnalyzePlan(fs, planFile)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Format and output result
		jsonOutput, err := plan.FormatResult(result)
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
	planCmd.AddCommand(planAnalyzeCmd)

	// Add flags for analyze command
	planAnalyzeCmd.Flags().StringP("file", "f", "", "Path to plan.json file to analyze")
}
