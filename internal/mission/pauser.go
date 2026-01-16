package mission

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/dnatag/mission-toolkit/internal/utils"
	"github.com/spf13/afero"
)

// Pauser handles pausing and restoring missions
type Pauser struct {
	fs         afero.Fs
	missionDir string
	reader     *Reader
}

// NewPauser creates a new Pauser instance
func NewPauser(fs afero.Fs, missionDir string) *Pauser {
	missionPath := filepath.Join(missionDir, "mission.md")
	return &Pauser{
		fs:         fs,
		missionDir: missionDir,
		reader:     NewReader(fs, missionPath),
	}
}

// Pause moves the current mission to .mission/paused/ with timestamp.
// The paused mission is saved with format: TIMESTAMP-MISSIONID-mission.md
// along with its execution log if it exists.
func (p *Pauser) Pause() error {
	missionPath := filepath.Join(p.missionDir, "mission.md")

	// Check if mission exists
	exists, err := afero.Exists(p.fs, missionPath)
	if err != nil {
		return fmt.Errorf("checking mission existence: %w", err)
	}

	if !exists {
		return fmt.Errorf("no current mission to pause")
	}

	// Read mission to get ID for naming
	mission, err := p.reader.Read()
	if err != nil {
		return fmt.Errorf("reading mission: %w", err)
	}

	// Create paused directory if it doesn't exist
	pausedDir := filepath.Join(p.missionDir, "paused")
	if err := p.fs.MkdirAll(pausedDir, 0755); err != nil {
		return fmt.Errorf("creating paused directory: %w", err)
	}

	// Generate timestamp for paused mission
	timestamp := time.Now().Format("20060102-150405")
	pausedFilename := fmt.Sprintf("%s-%s-mission.md", timestamp, mission.ID)
	pausedPath := filepath.Join(pausedDir, pausedFilename)

	// Copy mission file to paused directory
	if err := utils.CopyFile(p.fs, missionPath, pausedPath); err != nil {
		return fmt.Errorf("copying mission to paused directory: %w", err)
	}

	// Copy execution log if it exists
	logPath := filepath.Join(p.missionDir, "execution.log")
	if exists, _ := afero.Exists(p.fs, logPath); exists {
		pausedLogFilename := fmt.Sprintf("%s-%s-execution.log", timestamp, mission.ID)
		pausedLogPath := filepath.Join(pausedDir, pausedLogFilename)
		if err := utils.CopyFile(p.fs, logPath, pausedLogPath); err != nil {
			return fmt.Errorf("copying execution log: %w", err)
		}
	}

	// Remove current mission files
	if err := p.fs.Remove(missionPath); err != nil {
		return fmt.Errorf("removing current mission: %w", err)
	}

	if exists, _ := afero.Exists(p.fs, logPath); exists {
		if err := p.fs.Remove(logPath); err != nil {
			return fmt.Errorf("removing execution log: %w", err)
		}
	}

	return nil
}

// Restore moves a paused mission back to active state.
// If missionID is empty, restores the most recently paused mission.
// If missionID is provided, restores the specific mission with that ID.
func (p *Pauser) Restore(missionID string) error {
	pausedDir := filepath.Join(p.missionDir, "paused")

	// Check if paused directory exists
	exists, err := afero.Exists(p.fs, pausedDir)
	if err != nil {
		return fmt.Errorf("checking paused directory: %w", err)
	}

	if !exists {
		return fmt.Errorf("no paused missions found")
	}

	// List paused missions
	files, err := afero.ReadDir(p.fs, pausedDir)
	if err != nil {
		return fmt.Errorf("reading paused directory: %w", err)
	}

	var missionFile string
	var logFile string

	if missionID == "" {
		// Find the most recent paused mission by timestamp
		var latestTime time.Time
		for _, file := range files {
			if filepath.Ext(file.Name()) == ".md" {
				// Extract timestamp from filename (YYYYMMDD-HHMMSS format)
				if len(file.Name()) >= 15 {
					timeStr := file.Name()[:15] // YYYYMMDD-HHMMSS
					if t, err := time.Parse("20060102-150405", timeStr); err == nil {
						if t.After(latestTime) {
							latestTime = t
							missionFile = file.Name()
							// Look for corresponding log file
							logName := timeStr + file.Name()[15:]
							logName = logName[:len(logName)-10] + "-execution.log" // Replace -mission.md with -execution.log
							logFile = logName
						}
					}
				}
			}
		}
	} else {
		// Find specific mission by ID
		for _, file := range files {
			if filepath.Ext(file.Name()) == ".md" && len(file.Name()) > 16 {
				// Extract mission ID from filename: TIMESTAMP-MISSIONID-mission.md
				parts := file.Name()[16:] // Remove timestamp part
				if len(parts) > 11 && parts[len(parts)-11:] == "-mission.md" {
					extractedID := parts[:len(parts)-11] // Remove -mission.md suffix
					if extractedID == missionID {
						missionFile = file.Name()
						// Look for corresponding log file
						logFile = file.Name()[:len(file.Name())-11] + "-execution.log"
						break
					}
				}
			}
		}
	}

	if missionFile == "" {
		if missionID == "" {
			return fmt.Errorf("no paused missions found")
		}
		return fmt.Errorf("paused mission with ID %s not found", missionID)
	}

	// Check if current mission exists
	currentMissionPath := filepath.Join(p.missionDir, "mission.md")
	if exists, _ := afero.Exists(p.fs, currentMissionPath); exists {
		return fmt.Errorf("current mission exists, pause it first before restoring")
	}

	// Restore mission file
	pausedMissionPath := filepath.Join(pausedDir, missionFile)
	if err := utils.CopyFile(p.fs, pausedMissionPath, currentMissionPath); err != nil {
		return fmt.Errorf("restoring mission file: %w", err)
	}

	// Restore log file if it exists
	pausedLogPath := filepath.Join(pausedDir, logFile)
	if exists, _ := afero.Exists(p.fs, pausedLogPath); exists {
		currentLogPath := filepath.Join(p.missionDir, "execution.log")
		if err := utils.CopyFile(p.fs, pausedLogPath, currentLogPath); err != nil {
			return fmt.Errorf("restoring execution log: %w", err)
		}
	}

	// Remove paused files
	if err := p.fs.Remove(pausedMissionPath); err != nil {
		return fmt.Errorf("removing paused mission file: %w", err)
	}

	if exists, _ := afero.Exists(p.fs, pausedLogPath); exists {
		if err := p.fs.Remove(pausedLogPath); err != nil {
			return fmt.Errorf("removing paused log file: %w", err)
		}
	}

	return nil
}
