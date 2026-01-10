package mission

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dnatag/mission-toolkit/internal/plan"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

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
	mission, err := NewReader(w.fs).Read(path)
	if err != nil {
		return fmt.Errorf("failed to read mission: %w", err)
	}
	mission.Status = newStatus
	return w.Write(path, mission)
}

// CreateWithIntent creates initial mission.md with just INTENT section
func (w *Writer) CreateWithIntent(path string, missionID string, intent string) error {
	mission := &Mission{
		ID:        missionID,
		Iteration: 1,
		Status:    "planning",
		Body:      fmt.Sprintf("## INTENT\n%s\n", intent),
	}
	return w.Write(path, mission)
}

// UpdateSection updates a text section (intent, verification)
func (w *Writer) UpdateSection(path string, section string, content string) error {
	mission, err := NewReader(w.fs).Read(path)
	if err != nil {
		return fmt.Errorf("reading mission: %w", err)
	}

	lines := strings.Split(mission.Body, "\n")
	var result []string
	sectionHeader := "## " + strings.ToUpper(section)
	foundSection := false

	for i, line := range lines {
		if strings.TrimSpace(line) == sectionHeader {
			foundSection = true
			result = append(result, line)
			result = append(result, content)

			// Skip old content until next section or end
			for j := i + 1; j < len(lines); j++ {
				if strings.HasPrefix(strings.TrimSpace(lines[j]), "## ") {
					result = append(result, "")
					result = append(result, lines[j:]...)
					break
				}
				if j == len(lines)-1 {
					// Reached end without finding next section
					break
				}
			}
			break
		}
		result = append(result, line)
	}

	if !foundSection {
		return fmt.Errorf("section %s not found", section)
	}

	mission.Body = strings.Join(result, "\n")
	return w.Write(path, mission)
}

// UpdateList updates a list section (scope, plan)
func (w *Writer) UpdateList(path string, section string, items []string) error {
	mission, err := NewReader(w.fs).Read(path)
	if err != nil {
		return fmt.Errorf("reading mission: %w", err)
	}

	lines := strings.Split(mission.Body, "\n")
	var result []string
	sectionHeader := "## " + strings.ToUpper(section)
	foundSection := false

	for i, line := range lines {
		if strings.TrimSpace(line) == sectionHeader {
			result = append(result, line)
			// Add items
			if section == "plan" {
				for _, item := range items {
					result = append(result, "- [ ] "+item)
				}
			} else {
				result = append(result, items...)
			}
			foundSection = true
			// Skip old content until next section
			for j := i + 1; j < len(lines); j++ {
				if strings.HasPrefix(strings.TrimSpace(lines[j]), "## ") {
					result = append(result, "")
					result = append(result, lines[j:]...)
					break
				}
			}
			break
		}
		result = append(result, line)
	}

	if !foundSection {
		// Add new section at end
		result = append(result, "", sectionHeader)
		if section == "plan" {
			for _, item := range items {
				result = append(result, "- [ ] "+item)
			}
		} else {
			result = append(result, items...)
		}
	}

	mission.Body = strings.Join(result, "\n")
	return w.Write(path, mission)
}

// UpdateFrontmatter updates frontmatter fields
func (w *Writer) UpdateFrontmatter(path string, pairs []string) error {
	mission, err := NewReader(w.fs).Read(path)
	if err != nil {
		return fmt.Errorf("reading mission: %w", err)
	}

	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid frontmatter pair: %s", pair)
		}
		key, value := parts[0], parts[1]

		switch key {
		case "track":
			var track int
			fmt.Sscanf(value, "%d", &track)
			mission.Track = track
		case "type":
			mission.Type = value
		case "domains":
			// Store as type for now (would need Mission struct update for domains)
			mission.Type = value
		}
	}

	return w.Write(path, mission)
}

// CreateFromPlanFile creates a mission.md file from a plan.json file
func (w *Writer) CreateFromPlanFile(planPath string, missionPath string, missionID string, missionType string) error {
	data, err := afero.ReadFile(w.fs, planPath)
	if err != nil {
		return fmt.Errorf("reading plan file: %w", err)
	}

	var spec plan.PlanSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return fmt.Errorf("parsing plan file: %w", err)
	}

	if missionType == "clarification" {
		return w.CreateClarification(missionPath, missionID, &spec)
	}

	track := w.calculateTrack(&spec)
	return w.CreateFromPlan(missionPath, missionID, track, &spec)
}

// calculateTrack determines the track based on plan spec
func (w *Writer) calculateTrack(spec *plan.PlanSpec) int {
	if spec.Track != "" {
		var track int
		fmt.Sscanf(spec.Track, "TRACK %d", &track)
		return track
	}
	fileCount := len(spec.GetScopeFiles())
	if fileCount == 0 {
		return 1
	}
	if fileCount > 5 {
		return 3
	}
	return 2
}

// CreateClarification creates a clarification mission.md
func (w *Writer) CreateClarification(path string, missionID string, spec *plan.PlanSpec) error {
	mission := &Mission{
		ID:        missionID,
		Type:      "CLARIFICATION",
		Track:     2,
		Iteration: 1,
		Status:    "planned",
		Body:      w.buildClarificationBody(spec),
	}

	return w.Write(path, mission)
}

// buildClarificationBody constructs clarification mission body
func (w *Writer) buildClarificationBody(spec *plan.PlanSpec) string {
	var body strings.Builder

	body.WriteString("## INTENT\n")
	body.WriteString(spec.Intent)
	body.WriteString("\n\n## CLARIFICATION QUESTIONS\n")
	for i, question := range spec.ClarificationQuestions {
		body.WriteString(fmt.Sprintf("%d. %s\n", i+1, question))
	}
	body.WriteString("\n## INSTRUCTIONS\n")
	body.WriteString("Please answer the questions above to refine the intent.\n")
	body.WriteString("Once answered, run `/m.plan` again with the refined intent.\n")

	return body.String()
}

// CreateFromPlan creates a mission.md file from a PlanSpec
func (w *Writer) CreateFromPlan(path string, missionID string, track int, spec *plan.PlanSpec) error {
	mission := &Mission{
		ID:        missionID,
		Type:      spec.Type,
		Track:     track,
		Iteration: 1,
		Status:    "planned",
		Body:      w.buildBody(spec),
	}

	return w.Write(path, mission)
}

// buildBody constructs mission body from PlanSpec
func (w *Writer) buildBody(spec *plan.PlanSpec) string {
	var body strings.Builder

	body.WriteString("## INTENT\n")
	body.WriteString(spec.Intent)
	body.WriteString("\n\n## SCOPE\n")
	body.WriteString(strings.Join(spec.Scope, "\n"))
	body.WriteString("\n\n## PLAN\n")
	for _, step := range spec.Plan {
		body.WriteString("- [ ] ")
		body.WriteString(step)
		body.WriteString("\n")
	}
	body.WriteString("\n## VERIFICATION\n")
	body.WriteString(spec.Verification)
	body.WriteString("\n")

	return body.String()
}

// format converts a Mission struct to markdown with YAML frontmatter
func (w *Writer) format(mission *Mission) (string, error) {
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

	yamlData, err := yaml.Marshal(frontmatter)
	if err != nil {
		return "", fmt.Errorf("failed to marshal frontmatter: %w", err)
	}

	var sb strings.Builder
	sb.WriteString("---\n")
	sb.Write(yamlData)
	sb.WriteString("---\n\n")
	sb.WriteString(mission.Body)

	return sb.String(), nil
}
