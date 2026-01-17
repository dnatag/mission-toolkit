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
		include, _ := cmd.Flags().GetStringArray("include")
		exclude, _ := cmd.Flags().GetStringArray("exclude")

		// Validate mutual exclusivity
		if len(include) > 0 && len(exclude) > 0 {
			return fmt.Errorf("--include and --exclude are mutually exclusive")
		}

		manager := backlog.NewManager(missionDir)
		items, err := manager.List(include, exclude)
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
		patternID, _ := cmd.Flags().GetString("pattern-id")

		manager := backlog.NewManager(missionDir)

		if len(args) == 1 {
			if err := manager.AddWithPattern(args[0], itemType, patternID); err != nil {
				return fmt.Errorf("adding backlog item: %w", err)
			}
			if patternID != "" {
				count, _ := manager.GetPatternCount(patternID)
				fmt.Printf("Added backlog item (pattern: %s, count: %d): %s\n", patternID, count, args[0])
			} else {
				fmt.Printf("Added backlog item: %s\n", args[0])
			}
		} else {
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

// backlogResolveCmd is deprecated - pattern count tracking replaces RESOLVED workflow
var backlogResolveCmd = &cobra.Command{
	Use:        "resolve",
	Short:      "DEPRECATED: Use pattern-id tracking instead",
	Deprecated: "Pattern count tracking (--pattern-id) replaces the RESOLVED workflow",
	Hidden:     true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("resolve command is deprecated: use --pattern-id with 'backlog add' for Rule-of-Three tracking")
	},
}

func init() {
	rootCmd.AddCommand(backlogCmd)
	backlogCmd.AddCommand(backlogListCmd, backlogAddCmd, backlogCompleteCmd, backlogCleanupCmd)

	// Add flags
	backlogListCmd.Flags().StringArray("include", []string{}, "Include only these types (decomposed, refactor, future, completed)")
	backlogListCmd.Flags().StringArray("exclude", []string{}, "Exclude these types (decomposed, refactor, future, completed)")
	backlogAddCmd.Flags().String("type", "", "Item type (decomposed, refactor, future)")
	backlogAddCmd.MarkFlagRequired("type")
	backlogAddCmd.Flags().String("pattern-id", "", "Pattern ID for Rule-of-Three tracking (refactor type only)")
	backlogCompleteCmd.Flags().String("item", "", "Exact text of the item to complete")
	backlogCompleteCmd.MarkFlagRequired("item")
	backlogCleanupCmd.Flags().String("type", "", "Filter by item type (decomposed, refactor, future)")
}
