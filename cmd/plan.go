package cmd

import (
	"fmt"

	"github.com/dnatag/mission-toolkit/internal/mission"
	"github.com/dnatag/mission-toolkit/internal/plan"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var fs = afero.NewOsFs()

// getMissionID returns current mission ID or "unknown"
func getMissionID() string {
	id, _ := mission.NewIDService(fs, missionDir).GetCurrentID()
	if id == "" {
		return "unknown"
	}
	return id
}

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Mission planning tools",
	Long:  `Mission planning tools for analysis and validation.`,
}

// planAnalyzeCmd represents the plan analyze command
var planAnalyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze plan complexity and provide recommendations",
	RunE: func(cmd *cobra.Command, args []string) error {
		planFile, _ := cmd.Parent().Flags().GetString("file")
		analyzer := plan.NewAnalyzer(fs, getMissionID())

		result, err := analyzer.AnalyzePlan(fs, planFile)
		if err != nil {
			return fmt.Errorf("analyzing plan: %w", err)
		}

		jsonOutput, err := plan.FormatResult(result)
		if err != nil {
			return fmt.Errorf("formatting output: %w", err)
		}

		fmt.Println(jsonOutput)
		return nil
	},
}

// planValidateCmd represents the plan validate command
var planValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate plan with comprehensive safety checks",
	RunE: func(cmd *cobra.Command, args []string) error {
		planFile, _ := cmd.Parent().Flags().GetString("file")
		rootDir, _ := cmd.Flags().GetString("root")
		if rootDir == "" {
			rootDir = "."
		}

		validator := plan.NewValidatorService(fs, getMissionID(), rootDir)
		result, err := validator.ValidatePlan(planFile)
		if err != nil {
			return fmt.Errorf("validating plan: %w", err)
		}

		output := validator.FormatValidationOutput(result)
		jsonStr, err := output.FormatOutput()
		if err != nil {
			return fmt.Errorf("formatting output: %w", err)
		}

		fmt.Println(jsonStr)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
	planCmd.AddCommand(planAnalyzeCmd, planValidateCmd)

	// Add flags
	planCmd.PersistentFlags().StringP("file", "f", "", "Path to plan.json file")
	planAnalyzeCmd.MarkFlagRequired("file")
	planValidateCmd.MarkFlagRequired("file")
	planValidateCmd.Flags().StringP("root", "r", ".", "Project root directory for file validation")
}
