package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// backlogCmd represents the backlog command
var backlogCmd = &cobra.Command{
	Use:   "backlog",
	Short: "Manage the mission backlog",
	Long:  `Add items to or check the status of the mission backlog.`,
}

// backlogAddCmd adds a new item to the backlog
var backlogAddCmd = &cobra.Command{
	Use:   "add <item>",
	Short: "Add a new item to the backlog",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		item := args[0]
		backlogPath := filepath.Join(missionDir, "backlog.md")

		// Ensure mission directory exists
		if err := missionFs.MkdirAll(missionDir, 0755); err != nil {
			return fmt.Errorf("creating mission directory: %w", err)
		}

		// Read existing content or create new
		content := ""
		if exists, _ := afero.Exists(missionFs, backlogPath); exists {
			data, err := afero.ReadFile(missionFs, backlogPath)
			if err != nil {
				return fmt.Errorf("reading backlog: %w", err)
			}
			content = string(data)
		} else {
			content = "# Mission Backlog\n\n"
		}

		// Append new item
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		content += fmt.Sprintf("- [ ] %s\n", item)

		// Write back
		if err := afero.WriteFile(missionFs, backlogPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("writing backlog: %w", err)
		}

		fmt.Printf("Added to backlog: %s\n", item)
		return nil
	},
}

// backlogCheckCmd displays the current backlog
var backlogCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Display the current backlog",
	RunE: func(cmd *cobra.Command, args []string) error {
		backlogPath := filepath.Join(missionDir, "backlog.md")

		if exists, _ := afero.Exists(missionFs, backlogPath); !exists {
			fmt.Println("Backlog is empty.")
			return nil
		}

		data, err := afero.ReadFile(missionFs, backlogPath)
		if err != nil {
			return fmt.Errorf("reading backlog: %w", err)
		}

		fmt.Println(string(data))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(backlogCmd)
	backlogCmd.AddCommand(backlogAddCmd, backlogCheckCmd)
}
