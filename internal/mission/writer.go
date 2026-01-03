package mission

import (
	"fmt"
	"strings"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// PlanSpec represents the structure for mission planning data (minimal subset needed for mission creation)
type PlanSpec struct {
	Intent       string   `json:"intent"`
	Type         string   `json:"type"` // WET or DRY
	Scope        []string `json:"scope"`
	Plan         []string `json:"plan"`
	Verification string   `json:"verification"`
}

// Writer handles writing mission files and updating status
type Writer struct {
	fs afero.Fs
}

// NewWriter creates a new mission writer
func NewWriter(fs afero.Fs) *Writer {
	return &Writer{fs: fs}
}

// Write writes a Mission struct to a file with YAML frontmatter
func (w *Writer) Write(path string, mission *Mission) error {
	content, err := w.format(mission)
	if err != nil {
		return err
	}

	return afero.WriteFile(w.fs, path, []byte(content), 0644)
}

// UpdateStatus updates the status field in a mission file while preserving the body
func (w *Writer) UpdateStatus(path string, newStatus string) error {
	reader := NewReader(w.fs)
	mission, err := reader.Read(path)
	if err != nil {
		return fmt.Errorf("failed to read mission: %w", err)
	}

	mission.Status = newStatus
	return w.Write(path, mission)
}

// CreateFromPlan creates a mission.md file from a PlanSpec
func (w *Writer) CreateFromPlan(path string, missionID string, track int, spec *PlanSpec) error {
	// Build mission body
	var body strings.Builder

	// INTENT section
	body.WriteString("## INTENT\n")
	body.WriteString(spec.Intent)
	body.WriteString("\n\n")

	// SCOPE section
	body.WriteString("## SCOPE\n")
	for _, file := range spec.Scope {
		body.WriteString(file)
		body.WriteString("\n")
	}
	body.WriteString("\n")

	// PLAN section
	body.WriteString("## PLAN\n")
	for _, step := range spec.Plan {
		body.WriteString("- [ ] ")
		body.WriteString(step)
		body.WriteString("\n")
	}
	body.WriteString("\n")

	// VERIFICATION section
	body.WriteString("## VERIFICATION\n")
	body.WriteString(spec.Verification)
	body.WriteString("\n")

	// Create mission struct
	missionType := "WET"
	if spec.Type == "DRY" {
		missionType = "DRY"
	}

	mission := &Mission{
		ID:        missionID,
		Type:      missionType,
		Track:     track,
		Iteration: 1,
		Status:    "planned",
		Body:      body.String(),
	}

	return w.Write(path, mission)
}

// format converts a Mission struct to markdown with YAML frontmatter
func (w *Writer) format(mission *Mission) (string, error) {
	// Create frontmatter map
	frontmatter := map[string]interface{}{
		"id":        mission.ID,
		"type":      mission.Type,
		"track":     mission.Track,
		"iteration": mission.Iteration,
		"status":    mission.Status,
	}

	if mission.ParentMission != "" {
		frontmatter["parent_mission"] = mission.ParentMission
	}

	// Marshal to YAML
	yamlData, err := yaml.Marshal(frontmatter)
	if err != nil {
		return "", fmt.Errorf("failed to marshal frontmatter: %w", err)
	}

	// Build final content
	var sb strings.Builder
	sb.WriteString("---\n")
	sb.Write(yamlData)
	sb.WriteString("---\n\n")
	sb.WriteString(mission.Body)

	return sb.String(), nil
}
