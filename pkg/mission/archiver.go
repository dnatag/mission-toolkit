package mission

import (
	"fmt"
	"path/filepath"

	"github.com/dnatag/mission-toolkit/pkg/git"
	"github.com/dnatag/mission-toolkit/pkg/utils"
	"github.com/spf13/afero"
)

// Archiver handles archiving mission files to completed directory
type Archiver struct {
	*BaseService
	reader *Reader
	git    git.GitClient
}

// NewArchiver creates a new Archiver instance for the specified mission file path.
// The mission directory is derived from the path's directory component.
func NewArchiver(fs afero.Fs, path string, git git.GitClient) *Archiver {
	missionDir := filepath.Dir(path)
	base := NewBaseServiceWithPath(fs, missionDir, path)
	return &Archiver{
		BaseService: base,
		reader:      NewReader(fs, path),
		git:         git,
	}
}

// Archive copies mission artifacts to the completed directory.
// If force is true and no mission exists, this is a no-op.
// If force is false and no mission exists, returns an error.
func (a *Archiver) Archive(force bool) error {
	missionPath := a.MissionPath()

	// Check if mission file exists
	missionExists, err := afero.Exists(a.FS(), missionPath)
	if err != nil {
		return fmt.Errorf("checking mission existence: %w", err)
	}

	// Handle no mission case based on force flag
	if !missionExists {
		if force {
			// Force flag: silently succeed (no-op)
			return nil
		}
		// No force flag: return error
		return fmt.Errorf("no current mission to archive")
	}

	// Mission exists, proceed with archiving
	completedDir := filepath.Join(a.MissionDir(), "completed")
	if err := a.FS().MkdirAll(completedDir, 0755); err != nil {
		return fmt.Errorf("creating completed directory: %w", err)
	}

	missionID, err := a.reader.GetMissionID()
	if err != nil {
		return fmt.Errorf("getting mission ID: %w", err)
	}

	// Archive mission artifacts
	for _, filename := range []string{"mission.md", "execution.log"} {
		src := filepath.Join(a.MissionDir(), filename)
		if exists, _ := afero.Exists(a.FS(), src); !exists {
			continue
		}

		dst := filepath.Join(completedDir, fmt.Sprintf("%s-%s", missionID, filename))
		if err := utils.CopyFile(a.FS(), src, dst); err != nil {
			return fmt.Errorf("archiving %s: %w", filename, err)
		}
	}

	// Archive commit message
	commitMsg, err := a.git.GetCommitMessage("HEAD")
	if err != nil {
		return fmt.Errorf("archiving commit message: %w", err)
	}

	dst := filepath.Join(completedDir, fmt.Sprintf("%s-commit.msg", missionID))
	if err := afero.WriteFile(a.FS(), dst, []byte(commitMsg), 0644); err != nil {
		return fmt.Errorf("writing commit message: %w", err)
	}

	return nil
}

// CleanupObsoleteFiles removes obsolete mission files after successful archive.
// This function safely removes temporary mission files that are no longer needed
// after the mission has been archived to the completed directory.
func (a *Archiver) CleanupObsoleteFiles() error {
	// Files to clean up after successful archive
	filesToClean := []string{"execution.log", "mission.md", "id", "plan.json"}

	for _, filename := range filesToClean {
		filePath := filepath.Join(a.MissionDir(), filename)

		// Check if file exists before attempting removal
		exists, err := afero.Exists(a.FS(), filePath)
		if err != nil {
			return fmt.Errorf("checking existence of %s: %w", filename, err)
		}

		if exists {
			if err := a.FS().Remove(filePath); err != nil {
				return fmt.Errorf("removing %s: %w", filename, err)
			}
		}
	}

	return nil
}
