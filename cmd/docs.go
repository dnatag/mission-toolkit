package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/dnatag/mission-toolkit/internal/docs"
	"github.com/spf13/cobra"
)

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate CLI documentation schema",
	Long:  `Generate a JSON schema of all CLI commands for documentation purposes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")

		if format != "json" {
			return fmt.Errorf("unsupported format: %s (only 'json' is supported)", format)
		}

		schema := docs.GenerateSchema(rootCmd)

		output, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling schema: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
	docsCmd.Flags().String("format", "json", "Output format (json)")
}
