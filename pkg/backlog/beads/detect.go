package beads

import (
	"fmt"
	"os/exec"

	"github.com/spf13/afero"
)

// Available checks if Beads is available for use.
// Verifies: bd on PATH + .beads directory exists in the given filesystem.
func Available(fs afero.Fs) (bool, error) {
	_, err := exec.LookPath("bd")
	if err != nil {
		return false, nil
	}

	info, err := fs.Stat(".beads")
	if err != nil {
		return false, nil
	}

	if !info.IsDir() {
		return false, fmt.Errorf(".beads exists but is not a directory")
	}

	return true, nil
}

// AvailableInDir checks if Beads is available in a specific directory.
func AvailableInDir(dir string) (bool, error) {
	baseFS := afero.NewBasePathFs(afero.NewOsFs(), dir)
	return Available(baseFS)
}
