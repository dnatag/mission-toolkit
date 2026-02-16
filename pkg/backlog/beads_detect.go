package backlog

import (
	"fmt"
	"os/exec"

	"github.com/spf13/afero"
)

// BeadsAvailable checks if Beads is available for use.
// It verifies two conditions:
// 1. The 'bd' command is available on PATH (using exec.LookPath)
// 2. The .beads directory exists in the current working directory (or specified filesystem)
//
// Parameters:
//
//	fs: afero.Fs - Filesystem interface for directory operations (allows mocking for tests)
//
// Returns:
//
//	bool: true if both conditions are met, false otherwise
//	error: any error encountered during PATH lookup (directory errors return false, no error)
func BeadsAvailable(fs afero.Fs) (bool, error) {
	// Check if 'bd' command is on PATH
	_, err := exec.LookPath("bd")
	if err != nil {
		// Command not found on PATH is not an error, just means Beads is unavailable
		return false, nil
	}

	// Check if .beads directory exists
	beadsDir := ".beads"
	info, err := fs.Stat(beadsDir)
	if err != nil {
		// Directory not found is not an error, just means Beads is unavailable
		return false, nil
	}

	// Verify it's actually a directory
	if !info.IsDir() {
		return false, fmt.Errorf(".beads exists but is not a directory")
	}

	// Both checks passed
	return true, nil
}

// BeadsAvailableInDir checks if Beads is available in a specific directory.
// This is a convenience function that creates an OS filesystem rooted at the given directory.
// Useful for checking Beads availability in a specific project directory.
//
// Parameters:
//
//	dir: string - Directory to check for .beads subdirectory
//
// Returns:
//
//	bool: true if both 'bd' is on PATH and .beads directory exists in the specified directory
//	error: any error encountered during PATH lookup or directory validation
func BeadsAvailableInDir(dir string) (bool, error) {
	// Create an OS filesystem rooted at the specified directory
	baseFS := afero.NewBasePathFs(afero.NewOsFs(), dir)
	return BeadsAvailable(baseFS)
}
