package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/dnatag/mission-toolkit/internal/plan"
	"github.com/spf13/cobra"
)

// CheckResult represents the result of input validation
type CheckResult struct {
	IsValid  bool   `json:"is_valid"`
	Message  string `json:"message"`
	NextStep string `json:"next_step"`
}

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check [string to check]",
	Short: "Check if input is empty or whitespace",
	Long: `Validates that the input string is not empty or whitespace-only.
Useful for AI prompt validation before proceeding with execution.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]
		result := &CheckResult{}

		if plan.IsEmptyOrWhitespace(input) {
			result.IsValid = false
			result.Message = "Input is empty or whitespace"
			result.NextStep = "ASK_USER: What is your intent or goal for this task?"
		} else {
			result.IsValid = true
			result.Message = "Input is valid"
			result.NextStep = "PROCEED with execution"
		}

		jsonOutput, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("formatting output: %w", err)
		}
		fmt.Println(string(jsonOutput))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
