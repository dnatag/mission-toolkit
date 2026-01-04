package mission

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/spf13/afero"
)

// Archiver handles archiving mission files to completed directory
type Archiver struct {
	fs         afero.Fs
	missionDir string
}

// NewArchiver creates a new Archiver instance
func NewArchiver(fs afero.Fs, missionDir string) *Archiver {
	return &Archiver{
		fs:         fs,
		missionDir: missionDir,
	}
}

// Archive copies mission.md and execution.log to completed directory
func (a *Archiver) Archive() error {
	// Ensure completed directory exists
	completedDir := filepath.Join(a.missionDir, "completed")
	if err := a.fs.MkdirAll(completedDir, 0755); err != nil {
		return fmt.Errorf("creating completed directory: %w", err)
	}

	// Get mission ID for filename prefix
	missionPath := filepath.Join(a.missionDir, "mission.md")
	missionID, err := a.getMissionID(missionPath)
	if err != nil {
		return fmt.Errorf("getting mission ID: %w", err)
	}

	// Archive mission.md
	srcMission := filepath.Join(a.missionDir, "mission.md")
	dstMission := filepath.Join(completedDir, fmt.Sprintf("%s-mission.md", missionID))
	if err := a.copyFile(srcMission, dstMission); err != nil {
		return fmt.Errorf("archiving mission.md: %w", err)
	}

	// Archive execution.log
	srcLog := filepath.Join(a.missionDir, "execution.log")
	dstLog := filepath.Join(completedDir, fmt.Sprintf("%s-execution.log", missionID))
	if err := a.copyFile(srcLog, dstLog); err != nil {
		return fmt.Errorf("archiving execution.log: %w", err)
	}

	return nil
}

// getMissionID extracts mission ID from mission.md frontmatter
func (a *Archiver) getMissionID(missionPath string) (string, error) {
	content, err := afero.ReadFile(a.fs, missionPath)
	if err != nil {
		return "", fmt.Errorf("reading mission file: %w", err)
	}

	// Parse frontmatter to extract ID
	lines := string(content)
	if len(lines) < 10 {
		return "", fmt.Errorf("invalid mission file format")
	}

	// Simple extraction - look for "id: " line
	for i := 0; i < len(lines)-3; i++ {
		if lines[i:i+4] == "id: " {
			// Find end of line
			end := i + 4
			for end < len(lines) && lines[end] != '\n' {
				end++
			}
			return lines[i+4 : end], nil
		}
	}

	// Fallback to timestamp-based ID
	return fmt.Sprintf("archived-%d", time.Now().Unix()), nil
}

// copyFile copies a file from src to dst
func (a *Archiver) copyFile(src, dst string) error {
	content, err := afero.ReadFile(a.fs, src)
	if err != nil {
		return fmt.Errorf("reading source file %s: %w", src, err)
	}

	if err := afero.WriteFile(a.fs, dst, content, 0644); err != nil {
		return fmt.Errorf("writing destination file %s: %w", dst, err)
	}

	return nil
}
