package mission

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

// MissionStatus represents the current state of mission artifacts
type MissionStatus struct {
	HasActiveMission bool     `json:"has_active_mission"`
	MissionStatus    string   `json:"mission_status,omitempty"`
	MissionID        string   `json:"mission_id,omitempty"`
	StaleArtifacts   []string `json:"stale_artifacts_cleaned,omitempty"`
	Ready            bool     `json:"ready"`
	Message          string   `json:"message"`
}

// ValidationService handles mission state validation and cleanup
type ValidationService struct {
	fs          afero.Fs
	missionDir  string
	missionPath string
	idService   *IDService
}

// NewValidationService creates a new validation service
func NewValidationService(fs afero.Fs, missionDir string) *ValidationService {
	return &ValidationService{
		fs:          fs,
		missionDir:  missionDir,
		missionPath: filepath.Join(missionDir, "mission.md"),
		idService:   NewIDService(fs, missionDir),
	}
}

// CheckMissionState validates mission state and cleans up stale artifacts
func (v *ValidationService) CheckMissionState() (*MissionStatus, error) {
	status := &MissionStatus{StaleArtifacts: []string{}}

	// Check for existing mission.md
	if exists, _ := afero.Exists(v.fs, v.missionPath); exists {
		if data, err := afero.ReadFile(v.fs, v.missionPath); err == nil {
			content := string(data)
			status.HasActiveMission = true
			status.MissionStatus = v.extractField(content, "status: ", "unknown")
			status.MissionID = v.extractField(content, "id: ", "")
			status.Ready = false
			status.Message = "Active mission detected - requires user decision"
			return status, nil
		}
	}

	// Clean up stale artifacts
	if err := v.cleanupStaleArtifacts(status); err != nil {
		return nil, fmt.Errorf("failed to cleanup stale artifacts: %w", err)
	}

	// Generate new mission ID
	missionID, err := v.idService.GetOrCreateID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate mission ID: %w", err)
	}

	status.HasActiveMission = false
	status.MissionID = missionID
	status.Ready = true
	status.Message = "Ready for new mission"
	return status, nil
}

// cleanupStaleArtifacts removes old id, execution.log, and plan.json files
func (v *ValidationService) cleanupStaleArtifacts(status *MissionStatus) error {
	// Note: We do NOT remove "id" here because we might want to reuse it if it's valid,
	// or IDService handles it. However, the original code removed "id".
	// If we remove "id", GetOrCreateID will generate a new one.
	// We MUST remove plan.json to ensure no stale specs.
	artifacts := []string{"id", "execution.log", "plan.json"}

	for _, artifact := range artifacts {
		path := filepath.Join(v.missionDir, artifact)
		if exists, _ := afero.Exists(v.fs, path); exists {
			if err := v.fs.Remove(path); err != nil {
				return fmt.Errorf("failed to remove %s: %w", artifact, err)
			}
			status.StaleArtifacts = append(status.StaleArtifacts, artifact)
		}
	}
	return nil
}

// extractField extracts a field value from mission.md content
func (v *ValidationService) extractField(content, prefix, defaultValue string) string {
	for _, line := range strings.Split(content, "\n") {
		if trimmed := strings.TrimSpace(line); strings.HasPrefix(trimmed, prefix) {
			return strings.TrimSpace(trimmed[len(prefix):])
		}
	}
	return defaultValue
}

// ToJSON converts status to JSON string
func (s *MissionStatus) ToJSON() (string, error) {
	data, err := json.MarshalIndent(s, "", "  ")
	return string(data), err
}
