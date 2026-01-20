package cmd

import (
	"fmt"

	"github.com/dnatag/mission-toolkit/pkg/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display the current version of the Mission Toolkit CLI and embedded templates.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", version.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
