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

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.AddCommand(analyzeIntentCmd)
}
