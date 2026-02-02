package cmd

import (
	"fmt"
	"strings"

	"github.com/dnatag/mission-toolkit/pkg/checkpoint"
	"github.com/dnatag/mission-toolkit/pkg/mission"
	"github.com/spf13/cobra"
)

// checkpointCmd represents the checkpoint command
var checkpointCmd = &cobra.Command{
	Use:   "checkpoint",
	Short: "Checkpoint management for mission execution",
	Long:  `Create, restore, and clear checkpoints during mission execution.`,
}

// checkpointCreateCmd creates a new checkpoint
var checkpointCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a checkpoint of current working directory state",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get mission ID
		idService := mission.NewIDService(missionFs, missionPath)
		missionID, err := idService.GetCurrentID()
		if err != nil {
			return fmt.Errorf("getting mission ID: %w", err)
		}

		// Create checkpoint service
		svc, err := checkpoint.NewService(missionFs, missionDir)
		if err != nil {
			return fmt.Errorf("initializing checkpoint service: %w", err)
		}

		// Create checkpoint
		checkpointName, err := svc.Create(missionID)
		if err != nil {
			return fmt.Errorf("creating checkpoint: %w", err)
		}

		fmt.Printf("Checkpoint created: %s\n", checkpointName)

		// Display baseline tag info if this is the first checkpoint
		if strings.HasSuffix(checkpointName, "-1") {
			baselineTag := missionID + "-baseline"
			fmt.Printf("\n📌 Baseline tag created: %s\n", baselineTag)
			fmt.Printf("   View cumulative changes: git diff %s\n", baselineTag)
		}

		return nil
	},
}

// checkpointRestoreCmd restores to a checkpoint
var checkpointRestoreCmd = &cobra.Command{
	Use:   "restore <checkpoint>",
	Short: "Restore working directory to specified checkpoint",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		checkpointName := args[0]

		// Handle --all flag
		if cmd.Flags().Changed("all") {
			// Get mission ID
			idService := mission.NewIDService(missionFs, missionPath)
			missionID, err := idService.GetCurrentID()
			if err != nil {
				return fmt.Errorf("getting mission ID: %w", err)
			}

			// Create checkpoint service
			svc, err := checkpoint.NewService(missionFs, missionDir)
			if err != nil {
				return fmt.Errorf("initializing checkpoint service: %w", err)
			}

			// Restore all changes and clear checkpoints
			count, untrackedFiles, err := svc.RestoreAll(missionID)
			if err != nil {
				return fmt.Errorf("reverting all changes: %w", err)
			}

			fmt.Printf("Restored all changes, cleared %d checkpoint(s)\n", count)

			// Warn about untracked files if any exist
			if len(untrackedFiles) > 0 {
				fmt.Println("\n⚠️  Warning: Untracked files detected (created during mission):")
				for _, file := range untrackedFiles {
					fmt.Printf("  - %s\n", file)
				}
				fmt.Println("\nTo remove these files, run:")
				fmt.Println("  git clean -fd")
			}

			return nil
		}

		// Create checkpoint service
		svc, err := checkpoint.NewService(missionFs, missionDir)
		if err != nil {
			return fmt.Errorf("initializing checkpoint service: %w", err)
		}

		// Restore to checkpoint
		if err := svc.Restore(checkpointName); err != nil {
			return fmt.Errorf("restoring to checkpoint: %w", err)
		}

		fmt.Printf("Restored to checkpoint: %s\n", checkpointName)
		return nil
	},
}

// checkpointClearCmd clears all checkpoints for current mission
var checkpointClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all checkpoints for current mission",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get mission ID
		idService := mission.NewIDService(missionFs, missionPath)
		missionID, err := idService.GetCurrentID()
		if err != nil {
			return fmt.Errorf("getting mission ID: %w", err)
		}

		// Create checkpoint service
		svc, err := checkpoint.NewService(missionFs, missionDir)
		if err != nil {
			return fmt.Errorf("initializing checkpoint service: %w", err)
		}

		// Clear checkpoints
		count, err := svc.Clear(missionID)
		if err != nil {
			return fmt.Errorf("clearing checkpoints: %w", err)
		}

		fmt.Printf("Cleared %d checkpoint(s) for mission %s\n", count, missionID)
		return nil
	},
}

// checkpointCommitCmd creates the final commit for the mission
var checkpointCommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Create final commit for the mission and clear checkpoints",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get mission ID
		idService := mission.NewIDService(missionFs, missionPath)
		missionID, err := idService.GetCurrentID()
		if err != nil {
			return fmt.Errorf("getting mission ID: %w", err)
		}

		// Get commit message from flag
		commitMsg, _ := cmd.Flags().GetString("message")
		if commitMsg == "" {
			return fmt.Errorf("commit message cannot be empty")
		}

		// Create checkpoint service
		svc, err := checkpoint.NewService(missionFs, missionDir)
		if err != nil {
			return fmt.Errorf("initializing checkpoint service: %w", err)
		}

		// Consolidate and commit
		result, err := svc.Consolidate(missionID, commitMsg)
		if err != nil {
			return fmt.Errorf("consolidating commit: %w", err)
		}

		fmt.Printf("Final commit created: %s\n", result.CommitHash)

		if len(result.UnstagedFiles) > 0 {
			fmt.Printf("\n⚠️  UNSTAGED FILES DETECTED:\n")
			for _, f := range result.UnstagedFiles {
				fmt.Printf("   - %s\n", f)
			}
			fmt.Printf("\nStatus: Recorded for display in completion template.\n")
			fmt.Printf("No action needed - user will handle after mission completion.\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkpointCmd)
	checkpointCmd.AddCommand(checkpointCreateCmd, checkpointRestoreCmd, checkpointClearCmd, checkpointCommitCmd)

	// Add flags
	checkpointRestoreCmd.Flags().Bool("all", false, "Restore all mission changes")
	checkpointCommitCmd.Flags().StringP("message", "m", "", "Commit message for the final commit")
	checkpointCommitCmd.MarkFlagRequired("message")
}
