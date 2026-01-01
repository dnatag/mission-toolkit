package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dnatag/mission-toolkit/internal/mission"
	"github.com/dnatag/mission-toolkit/internal/plan"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	fs         = afero.NewOsFs()
	missionDir = ".mission"
)

// Common helper functions
func getMissionID() string {
	if missionID, err := mission.NewIDService(fs, missionDir).GetCurrentID(); err == nil {
		return missionID
	}
	return "unknown"
}

func handleError(msg string, err error) {
	fmt.Printf("Error: %s: %v\n", msg, err)
	os.Exit(1)
}

func outputJSON(data interface{}) {
	if jsonOutput, err := json.MarshalIndent(data, "", "  "); err != nil {
		handleError("formatting output", err)
	} else {
		fmt.Println(string(jsonOutput))
	}
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
		status, err := mission.NewValidationService(fs, missionDir).CheckMissionState()
		if err != nil {
			handleError("checking mission state", err)
		}
		outputJSON(status)
	},
}

// planAnalyzeCmd represents the plan analyze command
var planAnalyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze plan complexity and provide recommendations",
	Run: func(cmd *cobra.Command, args []string) {
		planFile, _ := cmd.Parent().Flags().GetString("file")
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
		planFile, _ := cmd.Parent().Flags().GetString("file")
		rootDir, _ := cmd.Flags().GetString("root")
		if rootDir == "" {
			rootDir = "."
		}

		validator := plan.NewValidatorService(fs, getMissionID(), rootDir)
		result, err := validator.ValidatePlan(planFile)
		if err != nil {
			handleError("validating plan", err)
		}

		output := validator.FormatValidationOutput(result)
		if jsonStr, err := output.FormatOutput(); err != nil {
			handleError("formatting output", err)
		} else {
			fmt.Println(jsonStr)
		}
	},
}

// planGenerateCmd represents the plan generate command
var planGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate mission.md from plan.json specification",
	Run: func(cmd *cobra.Command, args []string) {
		planFile, _ := cmd.Parent().Flags().GetString("file")
		outputFile, _ := cmd.Flags().GetString("output")
		if outputFile == "" {
			outputFile = ".mission/mission.md"
		}

		generator := plan.NewGeneratorService(fs, getMissionID())
		result, err := generator.GenerateMission(planFile, outputFile)
		if err != nil {
			handleError("generating mission", err)
		}

		output := generator.FormatGenerateOutput(result)
		if jsonStr, err := output.FormatOutput(); err != nil {
			handleError("formatting output", err)
		} else {
			fmt.Println(jsonStr)
		}
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
	planCmd.AddCommand(planCheckCmd, planAnalyzeCmd, planValidateCmd, planGenerateCmd)

	// Add flags
	planCmd.PersistentFlags().StringP("file", "f", "", "Path to plan.json file")
	planAnalyzeCmd.MarkFlagRequired("file")
	planValidateCmd.MarkFlagRequired("file")
	planGenerateCmd.MarkFlagRequired("file")
	planValidateCmd.Flags().StringP("root", "r", ".", "Project root directory for file validation")
	planGenerateCmd.Flags().StringP("output", "o", ".mission/mission.md", "Output path for generated mission.md")
}
