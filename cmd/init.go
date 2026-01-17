/*
Copyright © 2025 Yi Xie dnatag@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dnatag/mission-toolkit/internal/docs"
	"github.com/dnatag/mission-toolkit/internal/git"
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

If a Git repository is not found, it will be initialized automatically.`,
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

		// Write library templates
		if err := templates.WriteLibraryTemplates(fs, cwd, aiType); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing library templates: %v\n", err)
			os.Exit(1)
		}

		// Generate cli-reference.md from Cobra commands
		cliRefPath := filepath.Join(cwd, ".mission", "libraries", "cli-reference.md")
		cliRefContent := docs.GenerateMarkdown(rootCmd)
		if err := afero.WriteFile(fs, cliRefPath, []byte(cliRefContent), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing cli-reference.md: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Mission Toolkit project initialized successfully for AI type: %s\n", aiType)

		// Add .mission/ to .gitignore
		if err := git.EnsureEntry(fs, cwd, ".mission/"); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating .gitignore: %v\n", err)
			os.Exit(1)
		}

		// Check for Git repository and initialize if not found
		if _, err := os.Stat(".git"); os.IsNotExist(err) {
			fmt.Println("No Git repository found. Initializing a new one...")
			gitCmd := exec.Command("git", "init")
			if output, err := gitCmd.CombinedOutput(); err != nil {
				fmt.Fprintf(os.Stderr, "Error initializing Git repository: %s\n", output)
				os.Exit(1)
			}
			fmt.Println("Git repository initialized successfully.")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Add --ai flag
	initCmd.Flags().StringVar(&aiType, "ai", "", "AI assistant type (required). Supported: q, claude, kiro, opencode")
	initCmd.MarkFlagRequired("ai")
}
