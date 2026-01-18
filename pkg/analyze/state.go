package analyze

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

// PlanState represents the analysis state stored in plan.json during mission planning.
// It tracks the progression from original user intent through refinement, complexity
// analysis, and final plan generation.
type PlanState struct {
	OriginalIntent string   `json:"original_intent"`
	RefinedIntent  string   `json:"refined_intent,omitempty"`
	MissionType    string   `json:"mission_type,omitempty"` // WET or DRY
	Track          int      `json:"track,omitempty"`        // 1-4 complexity track
	Scope          []string `json:"scope,omitempty"`
	Domains        []string `json:"domains,omitempty"`
	PlanSteps      []string `json:"plan_steps,omitempty"`
	Verification   string   `json:"verification,omitempty"`
}

// LoadState reads and parses a plan.json file from the filesystem.
// Returns an error if the file doesn't exist or contains invalid JSON.
func LoadState(fs afero.Fs, path string) (*PlanState, error) {
	data, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, fmt.Errorf("reading plan state: %w", err)
	}

	var state PlanState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("parsing plan state: %w", err)
	}

	return &state, nil
}

// SaveState writes a PlanState to plan.json with proper formatting.
// Creates the parent directory if it doesn't exist.
func SaveState(fs afero.Fs, state *PlanState, path string) error {
	dir := filepath.Dir(path)
	if err := fs.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling plan state: %w", err)
	}

	if err := afero.WriteFile(fs, path, data, 0644); err != nil {
		return fmt.Errorf("writing plan state: %w", err)
	}

	return nil
}
