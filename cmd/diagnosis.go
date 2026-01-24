package cmd

import (
	"fmt"

	"github.com/dnatag/mission-toolkit/pkg/diagnosis"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	diagnosisFs   = afero.NewOsFs()
	diagnosisPath = ".mission/diagnosis.md"
)

// diagnosisCmd represents the diagnosis command
var diagnosisCmd = &cobra.Command{
	Use:   "diagnosis",
	Short: "Diagnosis management commands",
	Long:  `Commands for managing debug investigations and diagnosis files.`,
}

// diagnosisCreateCmd creates a new diagnosis.md file
var diagnosisCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new diagnosis.md file with symptom",
	RunE: func(cmd *cobra.Command, args []string) error {
		symptom, _ := cmd.Flags().GetString("symptom")
		if symptom == "" {
			return fmt.Errorf("--symptom is required")
		}

		if err := diagnosis.CreateDiagnosis(diagnosisFs, diagnosisPath, symptom); err != nil {
			return fmt.Errorf("creating diagnosis: %w", err)
		}

		fmt.Printf("Diagnosis created: %s\n", diagnosisPath)
		return nil
	},
}

// diagnosisUpdateCmd updates a section in diagnosis.md
var diagnosisUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a section or frontmatter in diagnosis.md",
	RunE: func(cmd *cobra.Command, args []string) error {
		section, _ := cmd.Flags().GetString("section")
		content, _ := cmd.Flags().GetString("content")
		status, _ := cmd.Flags().GetString("status")
		confidence, _ := cmd.Flags().GetString("confidence")

		// Handle frontmatter update
		if status != "" || confidence != "" {
			if err := diagnosis.UpdateFrontmatter(diagnosisFs, diagnosisPath, status, confidence); err != nil {
				return fmt.Errorf("updating diagnosis frontmatter: %w", err)
			}
			fmt.Println("Diagnosis frontmatter updated")
			return nil
		}

		// Handle section update
		if section == "" {
			return fmt.Errorf("--section is required (or use --status/--confidence for frontmatter)")
		}
		if content == "" {
			return fmt.Errorf("--content is required")
		}

		if err := diagnosis.UpdateSection(diagnosisFs, diagnosisPath, section, content); err != nil {
			return fmt.Errorf("updating diagnosis: %w", err)
		}

		fmt.Printf("Diagnosis section '%s' updated\n", section)
		return nil
	},
}

// diagnosisFinalizeCmd validates and displays diagnosis.md
var diagnosisFinalizeCmd = &cobra.Command{
	Use:   "finalize",
	Short: "Validate and display diagnosis.md for review",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := diagnosis.Finalize(diagnosisFs, diagnosisPath)
		if err != nil {
			return fmt.Errorf("finalizing diagnosis: %w", err)
		}

		fmt.Print(result)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(diagnosisCmd)
	diagnosisCmd.AddCommand(diagnosisCreateCmd, diagnosisUpdateCmd, diagnosisFinalizeCmd)

	diagnosisCreateCmd.Flags().String("symptom", "", "Symptom description for the diagnosis")
	diagnosisCreateCmd.MarkFlagRequired("symptom")

	diagnosisUpdateCmd.Flags().String("section", "", "Section to update (e.g., 'ROOT CAUSE', 'INVESTIGATION')")
	diagnosisUpdateCmd.Flags().String("content", "", "New content for the section")
	diagnosisUpdateCmd.Flags().String("status", "", "Update diagnosis status (investigating, confirmed, inconclusive)")
	diagnosisUpdateCmd.Flags().String("confidence", "", "Update confidence level (low, medium, high)")
}
