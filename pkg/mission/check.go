package mission

import (
	"fmt"
	"path/filepath"

	"github.com/dnatag/mission-toolkit/pkg/diagnosis"
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
	*BaseService
	reader    *Reader
	idService *IDService
	context   string
}

// NewCheckService creates a new check service for the specified mission file path.
// The mission directory is derived from the path's directory component.
func NewCheckService(fs afero.Fs, path string) *CheckService {
	missionDir := filepath.Dir(path)
	base := NewBaseServiceWithPath(fs, missionDir, path)
	return &CheckService{
		BaseService: base,
		reader:      NewReader(fs, path),
		idService:   NewIDService(fs, path),
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
	if exists, _ := afero.Exists(c.FS(), c.MissionPath()); exists {
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
	status.NextStep = "PROCEED to Step 1 (Intent Analysis)"
	return status, nil
}

// handleExistingMission processes existing mission state
func (c *CheckService) handleExistingMission(status *Status) (*Status, error) {
	mission, err := c.reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read mission file: %w", err)
	}

	status.HasActiveMission = true
	status.MissionStatus = mission.Status
	status.MissionID = mission.ID
	status.MissionIntent = mission.GetIntent()
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
		if mission.Status == "executed" || mission.Status == "completed" {
			status.Message = "Mission is ready for completion or re-completion"
			status.NextStep = "PROCEED with m.complete execution."
			return status, nil
		}
		status.Message = fmt.Sprintf("Mission status '%s' is not valid for m.complete", mission.Status)
		status.NextStep = "STOP. Mission must be in 'executed' or 'completed' status for m.complete."
		return status, nil
	}

	if c.context == "debug" {
		// Check if diagnosis.md exists and validate it
		diagnosisPath := filepath.Join(c.MissionDir(), "diagnosis.md")
		exists, _ := afero.Exists(c.FS(), diagnosisPath)

		if !exists {
			status.Message = "No diagnosis.md found"
			status.NextStep = "PROCEED with m.plan execution. Create diagnosis.md first with: m diagnosis create --symptom \"...\""
			return status, nil
		}

		// Validate diagnosis file structure
		if _, err := diagnosis.ReadDiagnosis(c.FS(), diagnosisPath); err != nil {
			status.Message = fmt.Sprintf("Invalid diagnosis.md: %v", err)
			status.NextStep = "STOP. Fix diagnosis.md structure or recreate with: m diagnosis create --symptom \"...\""
			return status, nil
		}

		status.Message = "Diagnosis file exists and mission is ready for planning"
		status.NextStep = "PROCEED with m.plan execution. The diagnosis will be consumed automatically."
		return status, nil
	}

	// Generic status routing (no context specified)
	switch mission.Status {
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

// cleanupStaleArtifacts removes stale mission artifacts only when called from m.plan context.
// This prevents accidental cleanup when mission check is called from other contexts like m.apply or m.complete.
func (c *CheckService) cleanupStaleArtifacts(status *Status) error {
	// Only clean up stale artifacts when called from m.plan context
	if c.context != "plan" {
		return nil
	}

	// Define artifacts that should be cleaned up from previous missions
	artifacts := []string{"id", "plan.json", "execution.log"}

	for _, artifact := range artifacts {
		path := filepath.Join(c.MissionDir(), artifact)
		if exists, _ := afero.Exists(c.FS(), path); exists {
			if err := c.FS().Remove(path); err != nil {
				return fmt.Errorf("failed to remove stale artifact %s: %w", artifact, err)
			}
			status.StaleArtifacts = append(status.StaleArtifacts, artifact)
		}
	}

	return nil
}
