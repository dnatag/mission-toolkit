package cmd

import (
	"fmt"

	"github.com/dnatag/mission-toolkit/internal/checkpoint"
	"github.com/dnatag/mission-toolkit/internal/mission"
	"github.com/spf13/cobra"
)

// checkpointCmd represents the checkpoint command
var checkpointCmd = &cobra.Command{
	Use:   "checkpoint",
	Short: "Checkpoint management for mission execution",
	Long:  `Create, revert, and clear checkpoints during mission execution.`,
}

// checkpointCreateCmd creates a new checkpoint
var checkpointCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a checkpoint of current working directory state",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get mission ID
		idService := mission.NewIDService(missionFs, missionDir)
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
		return nil
	},
}

// checkpointRevertCmd reverts to a checkpoint
var checkpointRevertCmd = &cobra.Command{
	Use:   "revert <checkpoint>",
	Short: "Revert working directory to specified checkpoint",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		checkpointName := args[0]

		// Handle --all flag
		if cmd.Flags().Changed("all") {
			// Get mission ID
			idService := mission.NewIDService(missionFs, missionDir)
			missionID, err := idService.GetCurrentID()
			if err != nil {
				return fmt.Errorf("getting mission ID: %w", err)
			}

			// Create checkpoint service
			svc, err := checkpoint.NewService(missionFs, missionDir)
			if err != nil {
				return fmt.Errorf("initializing checkpoint service: %w", err)
			}

			// Clear all checkpoints (which reverts all changes)
			count, err := svc.Clear(missionID)
			if err != nil {
				return fmt.Errorf("reverting all changes: %w", err)
			}

			fmt.Printf("Reverted all changes, cleared %d checkpoint(s)\n", count)
			return nil
		}

		// Create checkpoint service
		svc, err := checkpoint.NewService(missionFs, missionDir)
		if err != nil {
			return fmt.Errorf("initializing checkpoint service: %w", err)
		}

		// Revert to checkpoint
		if err := svc.Revert(checkpointName); err != nil {
			return fmt.Errorf("reverting to checkpoint: %w", err)
		}

		fmt.Printf("Reverted to checkpoint: %s\n", checkpointName)
		return nil
	},
}

// checkpointClearCmd clears all checkpoints for current mission
var checkpointClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all checkpoints for current mission",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get mission ID
		idService := mission.NewIDService(missionFs, missionDir)
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

func init() {
	rootCmd.AddCommand(checkpointCmd)
	checkpointCmd.AddCommand(checkpointCreateCmd, checkpointRevertCmd, checkpointClearCmd)

	// Add flags
	checkpointRevertCmd.Flags().Bool("all", false, "Revert all mission changes")
}
