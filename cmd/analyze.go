package cmd

import (
	"fmt"

	"github.com/dnatag/mission-toolkit/internal/analyze"
	"github.com/spf13/cobra"
)

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analysis tools for mission planning",
	Long:  `Analysis tools for breaking down intents into structured missions.`,
}

// analyzeIntentCmd provides intent analysis template with user input
var analyzeIntentCmd = &cobra.Command{
	Use:   "intent <user-input>",
	Short: "Provide intent analysis template with user input",
	Long:  `Load intent.md template and inject user input for LLM analysis.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userInput := args[0]
		service := analyze.NewIntentService()
		output, err := service.ProvideTemplate(userInput)
		if err != nil {
			return fmt.Errorf("providing intent template: %w", err)
		}
		fmt.Print(output)
		return nil
	},
}

// analyzeClarifyCmd provides clarification analysis template with current intent
var analyzeClarifyCmd = &cobra.Command{
	Use:   "clarify",
	Short: "Provide clarification analysis template with current intent",
	Long:  `Load clarification.md template and inject current intent from mission.md for LLM analysis.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		service := analyze.NewClarifyService()
		output, err := service.ProvideTemplate()
		if err != nil {
			return fmt.Errorf("providing clarification template: %w", err)
		}
		fmt.Print(output)
		return nil
	},
}

// analyzeScopeCmd provides scope analysis template with current intent
var analyzeScopeCmd = &cobra.Command{
	Use:   "scope",
	Short: "Provide scope analysis template with current intent",
	Long:  `Load scope.md template and inject current intent from mission.md for LLM analysis.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		service := analyze.NewScopeService()
		output, err := service.ProvideTemplate()
		if err != nil {
			return fmt.Errorf("providing scope template: %w", err)
		}
		fmt.Print(output)
		return nil
	},
}

// analyzeTestCmd provides test analysis template with current intent and scope
var analyzeTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Provide test analysis template with current intent and scope",
	Long:  `Load test.md template and inject current intent and scope from mission.md for LLM analysis.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		service := analyze.NewTestService()
		output, err := service.ProvideTemplate()
		if err != nil {
			return fmt.Errorf("providing test template: %w", err)
		}
		fmt.Print(output)
		return nil
	},
}

// analyzeDuplicationCmd provides duplication analysis template with current intent
var analyzeDuplicationCmd = &cobra.Command{
	Use:   "duplication",
	Short: "Provide duplication analysis template with current intent",
	Long:  `Load duplication.md template and inject current intent from mission.md for LLM analysis.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		service := analyze.NewDuplicationService()
		output, err := service.ProvideTemplate()
		if err != nil {
			return fmt.Errorf("providing duplication template: %w", err)
		}
		fmt.Print(output)
		return nil
	},
}

// analyzeComplexityCmd provides complexity analysis template with current intent and scope
var analyzeComplexityCmd = &cobra.Command{
	Use:   "complexity",
	Short: "Provide complexity analysis template with current intent and scope",
	Long:  `Load complexity.md template and inject current intent and scope from mission.md for LLM analysis.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		service := analyze.NewComplexityService()
		output, err := service.ProvideTemplate()
		if err != nil {
			return fmt.Errorf("providing complexity template: %w", err)
		}
		fmt.Print(output)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.AddCommand(analyzeIntentCmd, analyzeClarifyCmd, analyzeScopeCmd, analyzeTestCmd, analyzeDuplicationCmd, analyzeComplexityCmd)
}
