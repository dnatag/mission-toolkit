package cmd

import (
	"github.com/spf13/cobra"
)

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analysis tools for mission planning",
	Long:  `Analysis tools for breaking down intents into structured missions.`,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}
