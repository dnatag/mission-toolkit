package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/dnatag/mission-toolkit/pkg/backlog"
	"github.com/dnatag/mission-toolkit/pkg/backlog/beads"
	"github.com/dnatag/mission-toolkit/pkg/backlog/file"
	"github.com/spf13/cobra"
)

// newBacklogProvider creates the appropriate BacklogProvider based on Beads availability.
func newBacklogProvider() backlog.BacklogProvider {
	projectRoot := filepath.Dir(missionDir)
	if projectRoot == "." {
		projectRoot = ""
	}

	available, err := beads.AvailableInDir(projectRoot)
	if err == nil && available {
		return beads.NewProvider(projectRoot)
	}

	return file.NewManager(missionDir)
}

var backlogCmd = &cobra.Command{
	Use:   "backlog",
	Short: "Manage mission backlog items",
	Long:  `Manage mission backlog items including decomposed intents, refactoring opportunities, and future enhancements.`,
}

var backlogListCmd = &cobra.Command{
	Use:   "list",
	Short: "List backlog items",
	RunE: func(cmd *cobra.Command, args []string) error {
		include, _ := cmd.Flags().GetStringArray("include")
		exclude, _ := cmd.Flags().GetStringArray("exclude")

		if len(include) > 0 && len(exclude) > 0 {
			return fmt.Errorf("--include and --exclude are mutually exclusive")
		}

		provider := newBacklogProvider()
		items, err := provider.List(include, exclude)
		if err != nil {
			return fmt.Errorf("listing backlog: %w", err)
		}

		for _, item := range items {
			fmt.Println(item)
		}
		return nil
	},
}

var backlogAddCmd = &cobra.Command{
	Use:   "add [description...]",
	Short: "Add one or more backlog items",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		itemType, _ := cmd.Flags().GetString("type")
		patternID, _ := cmd.Flags().GetString("pattern-id")

		provider := newBacklogProvider()

		if len(args) == 1 {
			if err := provider.AddWithPattern(args[0], itemType, patternID); err != nil {
				return fmt.Errorf("adding backlog item: %w", err)
			}
			if patternID != "" {
				count, _ := provider.GetPatternCount(patternID)
				fmt.Printf("Added backlog item (pattern: %s, count: %d): %s\n", patternID, count, args[0])
			} else {
				fmt.Printf("Added backlog item: %s\n", args[0])
			}
		} else {
			if err := provider.AddMultiple(args, itemType); err != nil {
				return fmt.Errorf("adding backlog items: %w", err)
			}
			fmt.Printf("Added %d backlog items\n", len(args))
		}
		return nil
	},
}

var backlogCompleteCmd = &cobra.Command{
	Use:   "complete",
	Short: "Mark a backlog item as complete",
	RunE: func(cmd *cobra.Command, args []string) error {
		item, _ := cmd.Flags().GetString("item")
		if item == "" {
			return fmt.Errorf("--item flag is required")
		}

		provider := newBacklogProvider()
		if err := provider.Complete(item); err != nil {
			return fmt.Errorf("completing backlog item: %w", err)
		}

		fmt.Printf("Completed backlog item: %s\n", item)
		return nil
	},
}

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

		provider := newBacklogProvider()
		count, err := provider.Cleanup(itemType)
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

var backlogDecomposeCmd = &cobra.Command{
	Use:   "decompose",
	Short: "Decompose an epic into sub-intents with dependency tracking",
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonInput, _ := cmd.Flags().GetString("json")
		if jsonInput == "" {
			return fmt.Errorf("--json flag is required")
		}

		provider := newBacklogProvider()
		if err := provider.Decompose(jsonInput); err != nil {
			return fmt.Errorf("decomposing backlog: %w", err)
		}

		fmt.Println("Decomposed epic into sub-intents")
		return nil
	},
}

var backlogBeadsCmd = &cobra.Command{
	Use:   "beads",
	Short: "Beads integration commands",
}

var backlogBeadsAvailableCmd = &cobra.Command{
	Use:   "available",
	Short: "Check if Beads (bd) is available for backlog management",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectRoot := filepath.Dir(missionDir)
		if projectRoot == "." {
			projectRoot = ""
		}
		available, err := beads.AvailableInDir(projectRoot)
		if err != nil {
			return fmt.Errorf("checking beads availability: %w", err)
		}
		if available {
			fmt.Println("true")
		} else {
			fmt.Println("false")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(backlogCmd)
	backlogCmd.AddCommand(backlogListCmd, backlogAddCmd, backlogCompleteCmd, backlogCleanupCmd, backlogDecomposeCmd, backlogBeadsCmd)
	backlogBeadsCmd.AddCommand(backlogBeadsAvailableCmd)

	backlogListCmd.Flags().StringArray("include", []string{}, "Include only these types (decomposed, refactor, future, completed)")
	backlogListCmd.Flags().StringArray("exclude", []string{}, "Exclude these types (decomposed, refactor, future, completed)")
	backlogAddCmd.Flags().String("type", "", "Item type (decomposed, refactor, future)")
	backlogAddCmd.MarkFlagRequired("type")
	backlogAddCmd.Flags().String("pattern-id", "", "Pattern ID for Rule-of-Three tracking (refactor type only)")
	backlogCompleteCmd.Flags().String("item", "", "Exact text of the item to complete")
	backlogCompleteCmd.MarkFlagRequired("item")
	backlogCleanupCmd.Flags().String("type", "", "Filter by item type (decomposed, refactor, future)")
	backlogDecomposeCmd.Flags().String("json", "", "JSON input from m analyze decompose output")
	backlogDecomposeCmd.MarkFlagRequired("json")
}
