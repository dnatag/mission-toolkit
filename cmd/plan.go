package cmd

import (
	"fmt"
	"os"

	"github.com/dnatag/mission-toolkit/internal/mission"
	"github.com/dnatag/mission-toolkit/internal/plan"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// Common helper functions
func getMissionID() string {
	idService := mission.NewIDService(afero.NewOsFs(), ".mission")
	if missionID, err := idService.GetCurrentID(); err == nil {
		return missionID
	}
	return "unknown"
}

func handleError(msg string, err error) {
	fmt.Printf("Error: %s: %v\n", msg, err)
	os.Exit(1)
}

func requireFlag(cmd *cobra.Command, flag, errMsg string) string {
	value, _ := cmd.Flags().GetString(flag)
	if value == "" {
		fmt.Println(errMsg)
		os.Exit(1)
	}
	return value
}

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
	Run: func(cmd *cobra.Command, args []string) {
		status, err := mission.NewValidationService(afero.NewOsFs(), ".mission").CheckMissionState()
		if err != nil {
			handleError("checking mission state", err)
		}

		if jsonOutput, err := status.ToJSON(); err != nil {
			handleError("formatting output", err)
		} else {
			fmt.Println(jsonOutput)
		}
	},
}

// planAnalyzeCmd represents the plan analyze command
var planAnalyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze plan complexity and provide recommendations",
	Run: func(cmd *cobra.Command, args []string) {
		planFile := requireFlag(cmd, "file", "Error: --file flag is required")
		fs := afero.NewOsFs()
		analyzer := plan.NewAnalyzer(fs, getMissionID())

		result, err := analyzer.AnalyzePlan(fs, planFile)
		if err != nil {
			handleError("analyzing plan", err)
		}

		if jsonOutput, err := plan.FormatResult(result); err != nil {
			handleError("formatting output", err)
		} else {
			fmt.Println(jsonOutput)
		}
	},
}

// planValidateCmd represents the plan validate command
var planValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate plan with comprehensive safety checks",
	Run: func(cmd *cobra.Command, args []string) {
		planFile := requireFlag(cmd, "file", "Error: --file flag is required")
		rootDir, _ := cmd.Flags().GetString("root")
		if rootDir == "" {
			rootDir = "."
		}

		fs := afero.NewOsFs()
		validator := plan.NewValidatorService(fs, getMissionID(), rootDir)

		result, err := validator.ValidatePlan(planFile)
		if err != nil {
			handleError("validating plan", err)
		}

		if jsonOutput, err := result.ToJSON(); err != nil {
			handleError("formatting output", err)
		} else {
			fmt.Println(jsonOutput)
		}
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
	planCmd.AddCommand(planCheckCmd, planAnalyzeCmd, planValidateCmd)

	// Add flags
	planAnalyzeCmd.Flags().StringP("file", "f", "", "Path to plan.json file to analyze")
	planValidateCmd.Flags().StringP("file", "f", "", "Path to plan.json file to validate")
	planValidateCmd.Flags().StringP("root", "r", ".", "Project root directory for file validation")
}
