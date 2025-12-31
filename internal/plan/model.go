package plan

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

// LoadPlanSpec reads and parses a plan.json file
func LoadPlanSpec(fs afero.Fs, filename string) (*PlanSpec, error) {
	data, err := afero.ReadFile(fs, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read plan file: %w", err)
	}

	var spec PlanSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse plan JSON: %w", err)
	}

	return &spec, nil
}

// SavePlanSpec writes a PlanSpec to a plan.json file
func SavePlanSpec(fs afero.Fs, spec *PlanSpec, filename string) error {
	if err := spec.Validate(); err != nil {
		return fmt.Errorf("invalid plan spec: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := fs.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal plan JSON: %w", err)
	}

	if err := afero.WriteFile(fs, filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write plan file: %w", err)
	}

	return nil
}
