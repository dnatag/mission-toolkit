package mission

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

// Status represents the current state of mission artifacts
type Status struct {
	HasActiveMission bool     `json:"has_active_mission"`
	MissionStatus    string   `json:"mission_status,omitempty"`
	MissionID        string   `json:"mission_id,omitempty"`
	MissionIntent    string   `json:"mission_intent,omitempty"`
	StaleArtifacts   []string `json:"stale_artifacts_cleaned,omitempty"`
	Ready            bool     `json:"ready"`
	Message          string   `json:"message"`
	NextStep         string   `json:"next_step"`
}

// CheckService handles mission state validation using reader/writer
type CheckService struct {
	fs          afero.Fs
	missionDir  string
	missionPath string
	reader      *Reader
	idService   *IDService
	context     string
}

// NewCheckService creates a new check service
func NewCheckService(fs afero.Fs, missionDir string) *CheckService {
	return &CheckService{
		fs:          fs,
		missionDir:  missionDir,
		missionPath: filepath.Join(missionDir, "mission.md"),
		reader:      NewReader(fs),
		idService:   NewIDService(fs, missionDir),
		context:     "",
	}
}

// SetContext sets the context context for validation
func (c *CheckService) SetContext(ctx string) {
	c.context = ctx
}

// CheckMissionState validates mission state and cleans up stale artifacts
func (c *CheckService) CheckMissionState() (*Status, error) {
	status := &Status{StaleArtifacts: []string{}}

	// Check for existing mission.md
	if exists, _ := afero.Exists(c.fs, c.missionPath); exists {
		return c.handleExistingMission(status)
	}

	// Clean up stale artifacts
	if err := c.cleanupStaleArtifacts(status); err != nil {
		return nil, fmt.Errorf("failed to cleanup stale artifacts: %w", err)
	}

	// Generate new mission ID
	missionID, err := c.idService.GetOrCreateID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate mission ID: %w", err)
	}

	status.HasActiveMission = false
	status.MissionID = missionID
	status.Ready = true
	status.Message = "Ready for new mission"
	status.NextStep = "PROCEED to Step 1 (Intent Analysis)."
	return status, nil
}

// handleExistingMission processes existing mission state
func (c *CheckService) handleExistingMission(status *Status) (*Status, error) {
	mission, err := c.reader.Read(c.missionPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read mission file: %w", err)
	}

	status.HasActiveMission = true
	status.MissionStatus = mission.Status
	status.MissionID = mission.ID
	status.MissionIntent = c.extractIntent(mission.Body)
	status.Ready = false

	// Command-specific validation
	if c.context == "apply" {
		// Allow retry functionality for failed missions
		if mission.Status == "planned" || mission.Status == "active" || mission.Status == "failed" {
			status.Message = "Mission is ready for execution or re-execution"
			status.NextStep = "PROCEED with m.apply execution."
			return status, nil
		}
		status.Message = fmt.Sprintf("Mission status '%s' is not valid for m.apply", mission.Status)
		status.NextStep = "STOP. Mission must be in 'planned', 'active', or 'failed' status for m.apply."
		return status, nil
	}

	if c.context == "complete" {
		if mission.Status == "active" || mission.Status == "completed" {
			status.Message = "Mission is ready for completion or re-completion"
			status.NextStep = "PROCEED with m.complete execution."
			return status, nil
		}
		status.Message = fmt.Sprintf("Mission status '%s' is not valid for m.complete", mission.Status)
		status.NextStep = "STOP. Mission must be in 'active' or 'completed' status for m.complete."
		return status, nil
	}

	// Generic status routing (no context specified)
	switch mission.Status {
	case "clarifying":
		status.Message = "Mission is in clarifying state"
		status.NextStep = "Run the m.clarify prompt to resolve questions."
	case "planned", "active":
		status.Message = "Mission is ready for execution or re-execution"
		status.NextStep = "Run the m.apply prompt to execute this mission."
	case "completed":
		status.Message = "Mission is completed"
		status.NextStep = "Run the m.complete prompt to finalize this mission."
	default:
		status.Message = "Active mission detected - requires user decision"
		status.NextStep = "STOP. Use template libraries/displays/error-mission-exists.md to ask the user for a decision."
	}
	return status, nil
}

// extractIntent extracts the intent from mission body
func (c *CheckService) extractIntent(body string) string {
	lines := strings.Split(body, "\n")
	inIntent := false
	for _, line := range lines {
		if line == "## INTENT" {
			inIntent = true
			continue
		}
		if inIntent && line != "" && line[0] != '#' {
			return line
		}
		if inIntent && len(line) > 0 && line[0] == '#' {
			break
		}
	}
	return ""
}

// cleanupStaleArtifacts removes stale mission artifacts
func (c *CheckService) cleanupStaleArtifacts(status *Status) error {
	artifacts := []string{"id", "plan.json", "execution.log"}
	for _, artifact := range artifacts {
		path := filepath.Join(c.missionDir, artifact)
		if exists, _ := afero.Exists(c.fs, path); exists {
			if err := c.fs.Remove(path); err != nil {
				return fmt.Errorf("failed to remove %s: %w", artifact, err)
			}
			status.StaleArtifacts = append(status.StaleArtifacts, artifact)
		}
	}
	return nil
}
