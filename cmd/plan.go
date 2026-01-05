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

// planInitCmd creates a new plan.json
var planInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create new plan.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		intent, _ := cmd.Flags().GetString("intent")
		missionType, _ := cmd.Flags().GetString("type")
		scope, _ := cmd.Flags().GetStringSlice("scope")
		domain, _ := cmd.Flags().GetStringSlice("domain")
		questions, _ := cmd.Flags().GetStringSlice("question")

		spec := &plan.PlanSpec{
			Intent:                 intent,
			Type:                   missionType,
			Scope:                  scope,
			Domain:                 domain,
			ClarificationQuestions: questions,
		}

		if err := spec.Write(".mission/plan.json"); err != nil {
			return fmt.Errorf("writing plan: %w", err)
		}

		fmt.Println("Plan created: .mission/plan.json")
		return nil
	},
}

// planUpdateCmd updates existing plan.json
var planUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update existing plan.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		spec, err := plan.Read(".mission/plan.json")
		if err != nil {
			return fmt.Errorf("reading plan: %w", err)
		}

		if cmd.Flags().Changed("plan") {
			spec.Plan, _ = cmd.Flags().GetStringArray("plan")
		}
		if cmd.Flags().Changed("verification") {
			spec.Verification, _ = cmd.Flags().GetString("verification")
		}

		if err := spec.Write(".mission/plan.json"); err != nil {
			return fmt.Errorf("writing plan: %w", err)
		}

		fmt.Println("Plan updated: .mission/plan.json")
		return nil
	},
}

// planAnalyzeCmd represents the plan analyze command
var planAnalyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze plan complexity and provide recommendations",
	RunE: func(cmd *cobra.Command, args []string) error {
		planFile, _ := cmd.Parent().Flags().GetString("file")
		update, _ := cmd.Flags().GetBool("update")
		analyzer := plan.NewAnalyzer(fs, getMissionID())

		result, err := analyzer.AnalyzePlan(fs, planFile)
		if err != nil {
			return fmt.Errorf("analyzing plan: %w", err)
		}

		// Auto-update track if --update flag is set
		if update {
			spec, err := plan.Read(planFile)
			if err != nil {
				return fmt.Errorf("reading plan: %w", err)
			}
			spec.Track = fmt.Sprintf("TRACK %d", result.Track)
			if err := spec.Write(planFile); err != nil {
				return fmt.Errorf("updating plan: %w", err)
			}
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
	planCmd.AddCommand(planInitCmd, planUpdateCmd, planAnalyzeCmd, planValidateCmd)

	// Add flags for init
	planInitCmd.Flags().String("intent", "", "User intent")
	planInitCmd.Flags().String("type", "", "Mission type (WET/DRY)")
	planInitCmd.Flags().StringSlice("scope", nil, "Files in scope")
	planInitCmd.Flags().StringSlice("domain", nil, "Domains")
	planInitCmd.Flags().StringSlice("question", nil, "Clarification questions")
	planInitCmd.MarkFlagRequired("intent")

	// Add flags for update
	planUpdateCmd.Flags().StringArray("plan", nil, "Plan steps")
	planUpdateCmd.Flags().String("verification", "", "Verification command")

	// Add flags for analyze and validate
	planCmd.PersistentFlags().StringP("file", "f", "", "Path to plan.json file")
	planAnalyzeCmd.Flags().Bool("update", false, "Update plan.json with track")
	planAnalyzeCmd.MarkFlagRequired("file")
	planValidateCmd.MarkFlagRequired("file")
	planValidateCmd.Flags().StringP("root", "r", ".", "Project root directory for file validation")
}
