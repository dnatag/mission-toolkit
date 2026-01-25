package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/dnatag/mission-toolkit/pkg/validation"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check [string to check]",
	Short: "Check if input is empty or whitespace",
	Long: `Validates that the input string is not empty or whitespace-only.
Useful for AI prompt validation before proceeding with execution.
Also detects diagnosis.md for debug workflow integration.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fs := afero.NewOsFs()
		result := validation.Validate(args[0], fs, ".mission")

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
