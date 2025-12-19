/*
Copyright © 2025 Yi Xie dnatag@gmail.com

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/dnatag/mission-toolkit/internal/templates"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var aiType string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Mission Toolkit project with templates for specified AI type",
	Long: `Initialize Mission Toolkit project structure with templates
for the specified AI assistant type. Creates .mission/ directory with governance files
and AI-specific prompt templates.

Supported AI types: q, claude, gemini, cursor, codex, cline, kiro`,
	Run: func(cmd *cobra.Command, args []string) {
		// Validate AI type
		if err := templates.ValidateAIType(aiType); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Get current working directory
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
			os.Exit(1)
		}

		// Use real filesystem
		fs := afero.NewOsFs()

		// Write templates
		if err := templates.WriteTemplates(fs, cwd, aiType); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing templates: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Mission Toolkit project initialized successfully for AI type: %s\n", aiType)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Add --ai flag
	initCmd.Flags().StringVar(&aiType, "ai", "", "AI assistant type (required). Supported: q, claude, gemini, cursor, codex, cline, kiro")
	initCmd.MarkFlagRequired("ai")
}
