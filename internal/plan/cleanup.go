package plan

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

// CleanupStalePlanFiles removes stale plan.json files from the mission directory
func CleanupStalePlanFiles(fs afero.Fs, missionDir string) error {
	planFile := filepath.Join(missionDir, "plan.json")

	// Check if plan.json exists
	exists, err := afero.Exists(fs, planFile)
	if err != nil {
		return fmt.Errorf("failed to check plan file: %w", err)
	}
	if !exists {
		// File doesn't exist, nothing to clean
		return nil
	}

	// Remove the stale plan.json file
	if err := fs.Remove(planFile); err != nil {
		return fmt.Errorf("failed to remove stale plan file: %w", err)
	}

	return nil
}

// HasStalePlanFile checks if a stale plan.json exists in the mission directory
func HasStalePlanFile(fs afero.Fs, missionDir string) bool {
	planFile := filepath.Join(missionDir, "plan.json")
	exists, _ := afero.Exists(fs, planFile)
	return exists
}
