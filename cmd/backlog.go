package cmd

import (
	"fmt"

	"github.com/dnatag/mission-toolkit/internal/backlog"
	"github.com/spf13/cobra"
)

// backlogCmd represents the backlog command
var backlogCmd = &cobra.Command{
	Use:   "backlog",
	Short: "Manage mission backlog items",
	Long:  `Manage mission backlog items including decomposed intents, refactoring opportunities, and future enhancements.`,
}

// backlogListCmd lists backlog items
var backlogListCmd = &cobra.Command{
	Use:   "list",
	Short: "List backlog items",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		
		manager := backlog.NewManager(missionDir)
		items, err := manager.List(all)
		if err != nil {
			return fmt.Errorf("listing backlog: %w", err)
		}

		for _, item := range items {
			fmt.Println(item)
		}
		return nil
	},
}

// backlogAddCmd adds a new backlog item
var backlogAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new backlog item",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("item description is required")
		}
		
		itemType, _ := cmd.Flags().GetString("type")
		description := args[0]
		
		manager := backlog.NewManager(missionDir)
		if err := manager.Add(description, itemType); err != nil {
			return fmt.Errorf("adding backlog item: %w", err)
		}

		fmt.Printf("Added backlog item: %s\n", description)
		return nil
	},
}

// backlogCompleteCmd marks a backlog item as complete
var backlogCompleteCmd = &cobra.Command{
	Use:   "complete",
	Short: "Mark a backlog item as complete",
	RunE: func(cmd *cobra.Command, args []string) error {
		item, _ := cmd.Flags().GetString("item")
		if item == "" {
			return fmt.Errorf("--item flag is required")
		}
		
		manager := backlog.NewManager(missionDir)
		if err := manager.Complete(item); err != nil {
			return fmt.Errorf("completing backlog item: %w", err)
		}

		fmt.Printf("Completed backlog item: %s\n", item)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(backlogCmd)
	backlogCmd.AddCommand(backlogListCmd, backlogAddCmd, backlogCompleteCmd)

	// Add flags
	backlogListCmd.Flags().Bool("all", false, "Include completed items")
	backlogAddCmd.Flags().String("type", "", "Item type (decomposed, refactor, future)")
	backlogAddCmd.MarkFlagRequired("type")
	backlogCompleteCmd.Flags().String("item", "", "Exact text of the item to complete")
	backlogCompleteCmd.MarkFlagRequired("item")
}
