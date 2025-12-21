package cmd

import (
	"fmt"

	"github.com/dnatag/mission-toolkit/internal/templates"
	"github.com/dnatag/mission-toolkit/internal/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display the current version of the Mission Toolkit CLI and embedded templates.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("CLI: %s\n", version.CLIVersion)
		fmt.Printf("Templates: %s\n", templates.TemplateVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
