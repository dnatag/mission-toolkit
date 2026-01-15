package mission

import (
	"crypto/rand"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/afero"
)

var idPattern = regexp.MustCompile(`^\d{14}-\d{4}$`)

// IDService manages mission ID generation and persistence
type IDService struct {
	fs          afero.Fs
	missionDir  string
	idPath      string
	missionPath string
	reader      *Reader
}

// NewIDService creates a new mission ID service
func NewIDService(fs afero.Fs, missionDir string) *IDService {
	missionPath := filepath.Join(missionDir, "mission.md")
	return &IDService{
		fs:          fs,
		missionDir:  missionDir,
		idPath:      filepath.Join(missionDir, "id"),
		missionPath: missionPath,
		reader:      NewReader(fs, missionPath),
	}
}

// GetOrCreateID returns existing mission ID or creates new one (initialize once algorithm)
func (s *IDService) GetOrCreateID() (string, error) {
	// Check if ID already exists
	if id := s.readIDFile(); id != "" {
		return id, nil
	}

	// Generate and store new ID
	newID := s.generateID()
	if err := s.fs.MkdirAll(s.missionDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create mission directory: %w", err)
	}
	if err := afero.WriteFile(s.fs, s.idPath, []byte(newID), 0644); err != nil {
		return "", fmt.Errorf("failed to write mission ID: %w", err)
	}
	return newID, nil
}

// CleanupStaleID removes stale mission ID when no active mission exists
func (s *IDService) CleanupStaleID() error {
	exists, err := afero.Exists(s.fs, s.missionPath)
	if err != nil {
		return err
	}
	if !exists {
		if idExists, _ := afero.Exists(s.fs, s.idPath); idExists {
			return s.fs.Remove(s.idPath)
		}
	}
	return nil
}

// GetCurrentID returns current mission ID from active mission or stored ID
func (s *IDService) GetCurrentID() (string, error) {
	// First try to get from active mission.md
	if id := s.getIDFromMission(); id != "" {
		return id, nil
	}

	// Fallback to stored ID file
	if id := s.readIDFile(); id != "" {
		return id, nil
	}

	return "", fmt.Errorf("no active mission ID found")
}

// generateID creates new mission ID with format YYYYMMDDHHMMSS-RRRR
func (s *IDService) generateID() string {
	timestamp := time.Now().Format("20060102150405")

	// Generate 4 cryptographically secure random digits
	buf := make([]byte, 2)
	rand.Read(buf)
	random := (int(buf[0])<<8 | int(buf[1])) % 10000

	return fmt.Sprintf("%s-%04d", timestamp, random)
}

// isValidID validates mission ID format
func (s *IDService) isValidID(id string) bool {
	return idPattern.MatchString(id)
}

// readIDFile reads and validates ID from stored file
func (s *IDService) readIDFile() string {
	data, err := afero.ReadFile(s.fs, s.idPath)
	if err != nil {
		return ""
	}
	id := strings.TrimSpace(string(data))
	if id != "" && s.isValidID(id) {
		return id
	}
	return ""
}

// getIDFromMission extracts ID from mission.md file using reader
func (s *IDService) getIDFromMission() string {
	id, _ := s.reader.GetMissionID()
	return id
}
