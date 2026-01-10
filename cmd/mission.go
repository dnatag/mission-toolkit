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

// missionUpdateCmd updates mission status or sections
var missionUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update mission status or sections",
	RunE: func(cmd *cobra.Command, args []string) error {
		writer := mission.NewWriter(missionFs)
		missionPath := missionDir + "/mission.md"

		// Handle status update
		if cmd.Flags().Changed("status") {
			status, _ := cmd.Flags().GetString("status")
			if err := writer.UpdateStatus(missionPath, status); err != nil {
				return fmt.Errorf("updating mission status: %w", err)
			}
			fmt.Printf("Mission status updated to: %s\n", status)
			return nil
		}

		// Handle section update
		if cmd.Flags().Changed("section") {
			section, _ := cmd.Flags().GetString("section")

			// Text section update (intent, verification)
			if cmd.Flags().Changed("content") {
				content, _ := cmd.Flags().GetString("content")
				if err := writer.UpdateSection(missionPath, section, content); err != nil {
					return fmt.Errorf("updating section: %w", err)
				}
				fmt.Printf("Section '%s' updated\n", section)
				return nil
			}

			// List section update (scope, plan)
			if cmd.Flags().Changed("item") {
				items, _ := cmd.Flags().GetStringSlice("item")
				if err := writer.UpdateList(missionPath, section, items); err != nil {
					return fmt.Errorf("updating list: %w", err)
				}
				fmt.Printf("Section '%s' updated with %d items\n", section, len(items))
				return nil
			}

			return fmt.Errorf("--section requires either --content or --item")
		}

		// Handle frontmatter update
		if cmd.Flags().Changed("frontmatter") {
			frontmatter, _ := cmd.Flags().GetStringSlice("frontmatter")
			if err := writer.UpdateFrontmatter(missionPath, frontmatter); err != nil {
				return fmt.Errorf("updating frontmatter: %w", err)
			}
			fmt.Println("Frontmatter updated")
			return nil
		}

		return fmt.Errorf("no update operation specified")
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

// missionCreateCmd creates mission.md from plan.json or with just intent
var missionCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create mission.md from plan.json or with intent",
	RunE: func(cmd *cobra.Command, args []string) error {
		intent, _ := cmd.Flags().GetString("intent")
		missionType, _ := cmd.Flags().GetString("type")
		planFile, _ := cmd.Flags().GetString("file")

		// Get mission ID
		idService := mission.NewIDService(missionFs, missionDir)
		missionID, err := idService.GetOrCreateID()
		if err != nil {
			return fmt.Errorf("getting mission ID: %w", err)
		}

		writer := mission.NewWriter(missionFs)
		missionPath := missionDir + "/mission.md"

		// If intent is provided, create initial mission with just INTENT
		if intent != "" {
			if err := writer.CreateWithIntent(missionPath, missionID, intent); err != nil {
				return fmt.Errorf("creating mission with intent: %w", err)
			}
			fmt.Printf("Mission created: %s\n", missionPath)
			return nil
		}

		// Otherwise create from plan file
		if planFile == "" {
			planFile = ".mission/plan.json"
		}
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
	missionUpdateCmd.Flags().StringP("status", "s", "", "New mission status")
	missionUpdateCmd.Flags().String("section", "", "Section to update (intent, verification, scope, plan)")
	missionUpdateCmd.Flags().String("content", "", "Content for text sections")
	missionUpdateCmd.Flags().StringSlice("item", nil, "Items for list sections")
	missionUpdateCmd.Flags().StringSlice("frontmatter", nil, "Frontmatter key=value pairs")
	missionCreateCmd.Flags().String("intent", "", "Intent text for initial mission creation")
	missionCreateCmd.Flags().StringP("type", "t", "", "Mission type (clarification or final)")
	missionCreateCmd.Flags().StringP("file", "f", "", "Path to plan.json file (default: .mission/plan.json)")
}
