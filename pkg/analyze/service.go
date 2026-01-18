package analyze

import (
	"fmt"

	"github.com/spf13/afero"
)

// Service orchestrates the analysis workflow for mission planning.
// It manages plan state transitions and coordinates between different
// analysis phases (intent, scope, complexity, etc.).
type Service struct {
	fs        afero.Fs
	missionID string
}

// NewService creates a new analysis service instance.
func NewService(fs afero.Fs, missionID string) *Service {
	return &Service{
		fs:        fs,
		missionID: missionID,
	}
}

// InitializePlan creates a new plan.json with the given user intent.
// This is the first step in the analysis workflow.
func (s *Service) InitializePlan(intent string) error {
	state := &PlanState{
		OriginalIntent: intent,
	}

	if err := SaveState(s.fs, state, ".mission/plan.json"); err != nil {
		return fmt.Errorf("initializing plan: %w", err)
	}

	return nil
}

// GetPlanState retrieves the current plan state from disk.
// Returns an error if plan.json doesn't exist or is invalid.
func (s *Service) GetPlanState() (*PlanState, error) {
	return LoadState(s.fs, ".mission/plan.json")
}

// UpdatePlanState persists changes to the plan state.
// Used by analysis steps to incrementally build up the plan.
func (s *Service) UpdatePlanState(state *PlanState) error {
	return SaveState(s.fs, state, ".mission/plan.json")
}
