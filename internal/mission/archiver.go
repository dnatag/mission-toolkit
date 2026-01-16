package mission

import (
	"fmt"
	"path/filepath"

	"github.com/dnatag/mission-toolkit/internal/git"
	"github.com/dnatag/mission-toolkit/internal/utils"
	"github.com/spf13/afero"
)

// Archiver handles archiving mission files to completed directory
type Archiver struct {
	fs         afero.Fs
	missionDir string
	reader     *Reader
	git        git.GitClient
}

// NewArchiver creates a new Archiver instance
func NewArchiver(fs afero.Fs, missionDir string, git git.GitClient) *Archiver {
	missionPath := filepath.Join(missionDir, "mission.md")
	return &Archiver{
		fs:         fs,
		missionDir: missionDir,
		reader:     NewReader(fs, missionPath),
		git:        git,
	}
}

// Archive copies mission artifacts to the completed directory.
// If force is true and no mission exists, this is a no-op.
// If force is false and no mission exists, returns an error.
func (a *Archiver) Archive(force bool) error {
	missionPath := filepath.Join(a.missionDir, "mission.md")

	// Check if mission file exists
	missionExists, err := afero.Exists(a.fs, missionPath)
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
	completedDir := filepath.Join(a.missionDir, "completed")
	if err := a.fs.MkdirAll(completedDir, 0755); err != nil {
		return fmt.Errorf("creating completed directory: %w", err)
	}

	missionID, err := a.reader.GetMissionID()
	if err != nil {
		return fmt.Errorf("getting mission ID: %w", err)
	}

	// Archive mission artifacts
	for _, filename := range []string{"mission.md", "execution.log"} {
		src := filepath.Join(a.missionDir, filename)
		if exists, _ := afero.Exists(a.fs, src); !exists {
			continue
		}

		dst := filepath.Join(completedDir, fmt.Sprintf("%s-%s", missionID, filename))
		if err := utils.CopyFile(a.fs, src, dst); err != nil {
			return fmt.Errorf("archiving %s: %w", filename, err)
		}
	}

	// Archive commit message
	commitMsg, err := a.git.GetCommitMessage("HEAD")
	if err != nil {
		return fmt.Errorf("archiving commit message: %w", err)
	}

	dst := filepath.Join(completedDir, fmt.Sprintf("%s-commit.msg", missionID))
	if err := afero.WriteFile(a.fs, dst, []byte(commitMsg), 0644); err != nil {
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
		filePath := filepath.Join(a.missionDir, filename)

		// Check if file exists before attempting removal
		exists, err := afero.Exists(a.fs, filePath)
		if err != nil {
			return fmt.Errorf("checking existence of %s: %w", filename, err)
		}

		if exists {
			if err := a.fs.Remove(filePath); err != nil {
				return fmt.Errorf("removing %s: %w", filename, err)
			}
		}
	}

	return nil
}
