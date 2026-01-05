package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/dnatag/mission-toolkit/internal/git"
	"github.com/dnatag/mission-toolkit/internal/mission"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	missionFs  = afero.NewOsFs()
	missionDir = ".mission"
)

// missionCmd represents the mission command
var missionCmd = &cobra.Command{
	Use:   "mission",
	Short: "Mission management commands",
	Long:  `Commands for managing mission state, status, and IDs.`,
}

// missionCheckCmd checks mission state
var missionCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check mission state and validate artifacts",
	RunE: func(cmd *cobra.Command, args []string) error {
		checkService := mission.NewCheckService(missionFs, missionDir)

		// Set command context if provided
		context, _ := cmd.Flags().GetString("context")
		if context != "" {
			checkService.SetContext(context)
		}

		status, err := checkService.CheckMissionState()
		if err != nil {
			return fmt.Errorf("checking mission state: %w", err)
		}

		jsonOutput, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			return fmt.Errorf("formatting output: %w", err)
		}
		fmt.Println(string(jsonOutput))
		return nil
	},
}

// missionUpdateCmd updates mission status
var missionUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update mission status",
	RunE: func(cmd *cobra.Command, args []string) error {
		status, _ := cmd.Flags().GetString("status")
		if status == "" {
			return fmt.Errorf("--status flag is required")
		}

		writer := mission.NewWriter(missionFs)
		if err := writer.UpdateStatus(missionDir+"/mission.md", status); err != nil {
			return fmt.Errorf("updating mission status: %w", err)
		}

		fmt.Printf("Mission status updated to: %s\n", status)
		return nil
	},
}

// missionIDCmd gets or creates mission ID
var missionIDCmd = &cobra.Command{
	Use:   "id",
	Short: "Get or create mission ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		idService := mission.NewIDService(missionFs, missionDir)

		id, err := idService.GetCurrentID()
		if err != nil {
			id, err = idService.GetOrCreateID()
			if err != nil {
				return fmt.Errorf("generating mission ID: %w", err)
			}
		}

		fmt.Println(id)
		return nil
	},
}

// missionCreateCmd creates mission.md from plan.json
var missionCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create mission.md from plan.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		missionType, _ := cmd.Flags().GetString("type")
		planFile, _ := cmd.Flags().GetString("file")
		if planFile == "" {
			planFile = ".mission/plan.json"
		}

		// Get mission ID
		idService := mission.NewIDService(missionFs, missionDir)
		missionID, err := idService.GetOrCreateID()
		if err != nil {
			return fmt.Errorf("getting mission ID: %w", err)
		}

		// Create mission using Writer
		writer := mission.NewWriter(missionFs)
		missionPath := missionDir + "/mission.md"

		if err := writer.CreateFromPlanFile(planFile, missionPath, missionID, missionType); err != nil {
			return fmt.Errorf("creating mission file: %w", err)
		}

		fmt.Printf("Mission created: %s\n", missionPath)
		return nil
	},
}

// missionArchiveCmd archives mission.md and execution.log to completed directory
// and cleans up obsolete files
var missionArchiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive mission files to completed directory and clean up obsolete files",
	RunE: func(cmd *cobra.Command, args []string) error {
		gitClient := git.NewCmdGitClient(".")
		archiver := mission.NewArchiver(missionFs, missionDir, gitClient)

		if err := archiver.Archive(); err != nil {
			return fmt.Errorf("archiving mission: %w", err)
		}

		// Clean up obsolete files after successful archive
		if err := archiver.CleanupObsoleteFiles(); err != nil {
			return fmt.Errorf("cleaning up obsolete files: %w", err)
		}

		fmt.Println("Mission archived successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(missionCmd)
	missionCmd.AddCommand(missionCheckCmd, missionUpdateCmd, missionIDCmd, missionCreateCmd, missionArchiveCmd)

	// Add flags
	missionCheckCmd.Flags().StringP("context", "c", "", "Context for validation (apply or complete)")
	missionUpdateCmd.Flags().StringP("status", "s", "", "New mission status (required)")
	missionUpdateCmd.MarkFlagRequired("status")
	missionCreateCmd.Flags().StringP("type", "t", "", "Mission type (clarification or final)")
	missionCreateCmd.Flags().StringP("file", "f", "", "Path to plan.json file (default: .mission/plan.json)")
	missionCreateCmd.MarkFlagRequired("type")
}
