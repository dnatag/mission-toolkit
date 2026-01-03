package cmd

import (
	"encoding/json"
	"fmt"

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

func init() {
	rootCmd.AddCommand(missionCmd)
	missionCmd.AddCommand(missionCheckCmd, missionUpdateCmd, missionIDCmd)

	// Add flags
	missionUpdateCmd.Flags().StringP("status", "s", "", "New mission status (required)")
	missionUpdateCmd.MarkFlagRequired("status")
}
