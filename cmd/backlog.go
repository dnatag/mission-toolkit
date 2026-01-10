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
		itemType, _ := cmd.Flags().GetString("type")

		manager := backlog.NewManager(missionDir)
		items, err := manager.List(all, itemType)
		if err != nil {
			return fmt.Errorf("listing backlog: %w", err)
		}

		for _, item := range items {
			fmt.Println(item)
		}
		return nil
	},
}

// backlogAddCmd adds backlog items (single or multiple)
var backlogAddCmd = &cobra.Command{
	Use:   "add [description...]",
	Short: "Add one or more backlog items",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		itemType, _ := cmd.Flags().GetString("type")

		manager := backlog.NewManager(missionDir)

		if len(args) == 1 {
			// Single item - use existing Add method
			if err := manager.Add(args[0], itemType); err != nil {
				return fmt.Errorf("adding backlog item: %w", err)
			}
			fmt.Printf("Added backlog item: %s\n", args[0])
		} else {
			// Multiple items - use AddMultiple
			if err := manager.AddMultiple(args, itemType); err != nil {
				return fmt.Errorf("adding backlog items: %w", err)
			}
			fmt.Printf("Added %d backlog items\n", len(args))
		}
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

// backlogCleanupCmd removes completed items from the backlog
var backlogCleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Remove completed items from the backlog",
	Long: `Remove completed items from the COMPLETED section of the backlog.

By default, removes all completed items. Use --type to filter by item type.

Examples:
  m backlog cleanup                    # Remove all completed items
  m backlog cleanup --type decomposed  # Remove only completed decomposed epic items`,
	RunE: func(cmd *cobra.Command, args []string) error {
		itemType, _ := cmd.Flags().GetString("type")

		manager := backlog.NewManager(missionDir)
		count, err := manager.Cleanup(itemType)
		if err != nil {
			return fmt.Errorf("cleaning up backlog: %w", err)
		}

		if count == 0 {
			fmt.Println("No completed items to remove")
		} else {
			fmt.Printf("Removed %d completed item(s)\n", count)
		}
		return nil
	},
}

// backlogResolveCmd marks a refactor opportunity as resolved via DRY conversion
var backlogResolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Mark a refactor opportunity as resolved via DRY conversion",
	Long: `Mark a refactor opportunity as resolved via DRY conversion.

This marks the item in-place with [RESOLVED] prefix and timestamp, allowing future
duplication detection to recognize this pattern has been addressed.

Example:
  m backlog resolve --item "Refactor email validation in handlers"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		item, _ := cmd.Flags().GetString("item")
		if item == "" {
			return fmt.Errorf("--item flag is required")
		}

		manager := backlog.NewManager(missionDir)
		if err := manager.Resolve(item); err != nil {
			return fmt.Errorf("marking item as resolved: %w", err)
		}

		fmt.Printf("Marked as resolved: %s\n", item)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(backlogCmd)
	backlogCmd.AddCommand(backlogListCmd, backlogAddCmd, backlogCompleteCmd, backlogResolveCmd, backlogCleanupCmd)

	// Add flags
	backlogListCmd.Flags().Bool("all", false, "Include completed items")
	backlogListCmd.Flags().String("type", "", "Filter by item type (decomposed, refactor, future)")
	backlogAddCmd.Flags().String("type", "", "Item type (decomposed, refactor, future)")
	backlogAddCmd.MarkFlagRequired("type")
	backlogCompleteCmd.Flags().String("item", "", "Exact text of the item to complete")
	backlogCompleteCmd.MarkFlagRequired("item")
	backlogResolveCmd.Flags().String("item", "", "Exact text of the refactor item to mark as resolved")
	backlogResolveCmd.MarkFlagRequired("item")
	backlogCleanupCmd.Flags().String("type", "", "Filter by item type (decomposed, refactor, future)")
}
