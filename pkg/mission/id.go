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
	*BaseService
	idPath string
	reader *Reader
}

// NewIDService creates a new mission ID service for the specified mission file path.
// The mission directory and ID file path are derived from the path's directory component.
func NewIDService(fs afero.Fs, path string) *IDService {
	base := NewBaseServiceWithPath(fs, "", path)
	missionDir := filepath.Dir(path)
	return &IDService{
		BaseService: base,
		idPath:      filepath.Join(missionDir, "id"),
		reader:      NewReader(fs, path),
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
	if err := s.FS().MkdirAll(s.MissionDir(), 0755); err != nil {
		return "", fmt.Errorf("failed to create mission directory: %w", err)
	}
	if err := afero.WriteFile(s.FS(), s.idPath, []byte(newID), 0644); err != nil {
		return "", fmt.Errorf("failed to write mission ID: %w", err)
	}
	return newID, nil
}

// CleanupStaleID removes stale mission ID when no active mission exists
func (s *IDService) CleanupStaleID() error {
	exists, err := afero.Exists(s.FS(), s.MissionPath())
	if err != nil {
		return err
	}
	if !exists {
		if idExists, _ := afero.Exists(s.FS(), s.idPath); idExists {
			return s.FS().Remove(s.idPath)
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
	data, err := afero.ReadFile(s.FS(), s.idPath)
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
