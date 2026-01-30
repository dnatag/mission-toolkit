package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/dnatag/mission-toolkit/pkg/docs"
	"github.com/spf13/cobra"
)

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate CLI documentation schema",
	Long:  `Generate a JSON schema of all CLI commands for documentation purposes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")

		switch format {
		case "json":
			schema := docs.GenerateSchema(rootCmd)
			output, err := json.MarshalIndent(schema, "", "  ")
			if err != nil {
				return fmt.Errorf("marshaling schema: %w", err)
			}
			fmt.Println(string(output))
		case "condensed":
			output := docs.GenerateCondensedMarkdown(rootCmd)
			fmt.Println(output)
		default:
			return fmt.Errorf("unsupported format: %s (supported: json, condensed)", format)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
	docsCmd.Flags().String("format", "json", "Output format (json, condensed)")
}
