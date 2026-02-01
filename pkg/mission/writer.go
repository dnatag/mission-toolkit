package mission

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dnatag/mission-toolkit/pkg/logger"
	"github.com/dnatag/mission-toolkit/pkg/md"
	"github.com/spf13/afero"
)

// Writer handles writing mission files and updating status.
type Writer struct {
	*BaseService
	loggerConfig *logger.Config // Logger configuration for execution logging
}

// NewWriter creates a new mission writer for the specified mission file path.
// The mission directory is derived from the path's directory component.
// Logger is configured with default settings (both console and file output).
func NewWriter(fs afero.Fs, path string) *Writer {
	missionDir := filepath.Dir(path)
	return &Writer{
		BaseService:  NewBaseServiceWithPath(fs, missionDir, path),
		loggerConfig: logger.DefaultConfig(),
	}
}

// Write writes a Mission struct to the writer's mission file with YAML frontmatter.
func (w *Writer) Write(mission *Mission) error {
	content, err := w.format(mission)
	if err != nil {
		return err
	}

	return afero.WriteFile(w.FS(), w.MissionPath(), []byte(content), 0644)
}

// UpdateStatus updates the status field in the mission file while preserving the body.
func (w *Writer) UpdateStatus(newStatus string) error {
	mission, err := NewReader(w.FS(), w.MissionPath()).Read()
	if err != nil {
		return fmt.Errorf("failed to read mission: %w", err)
	}
	mission.Status = newStatus
	return w.Write(mission)
}

// CreateWithIntent creates initial mission.md with just INTENT section.
func (w *Writer) CreateWithIntent(missionID string, intent string) error {
	mission := &Mission{
		ID:        missionID,
		Iteration: 1,
		Status:    "planning",
		Body:      fmt.Sprintf("## INTENT\n%s\n", intent),
	}
	return w.Write(mission)
}

// normalizePlanContent converts all list formats to checkbox format.
func normalizePlanContent(content string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		indent := line[:len(line)-len(trimmed)]
		var normalized string

		// Already checkbox format
		if strings.HasPrefix(trimmed, "- [ ] ") || strings.HasPrefix(trimmed, "- [x] ") {
			continue
		}
		// Dash with number: "- 1. item" → "- [ ] item"
		if matched, _ := regexp.MatchString(`^-\s+\d+\.\s`, trimmed); matched {
			parts := strings.SplitN(trimmed, ". ", 2)
			if len(parts) == 2 {
				normalized = indent + "- [ ] " + parts[1]
			}
		} else if matched, _ := regexp.MatchString(`^\d+\.\s`, trimmed); matched {
			// Pure numbered: "1. item" → "- [ ] item"
			parts := strings.SplitN(trimmed, ". ", 2)
			if len(parts) == 2 {
				normalized = indent + "- [ ] " + parts[1]
			}
		} else if strings.HasPrefix(trimmed, "- ") {
			// Plain dash: "- item" → "- [ ] item"
			normalized = indent + "- [ ] " + strings.TrimPrefix(trimmed, "- ")
		} else if strings.HasPrefix(trimmed, "* ") {
			// Asterisk: "* item" → "- [ ] item"
			normalized = indent + "- [ ] " + strings.TrimPrefix(trimmed, "* ")
		}

		if normalized != "" {
			lines[i] = normalized
		}
	}
	return strings.Join(lines, "\n")
}

// UpdateSection updates a text section (intent, verification).
func (w *Writer) UpdateSection(section string, content string) error {
	// Allow empty section name for backward compatibility (creates "## " section)
	if section == "" {
		mission, err := NewReader(w.FS(), w.MissionPath()).Read()
		if err != nil {
			return fmt.Errorf("reading mission: %w", err)
		}
		mission.Body = mission.Body + "\n## \n" + content + "\n"
		return w.Write(mission)
	}

	// Normalize PLAN section to checkbox format
	if strings.ToLower(section) == "plan" {
		content = normalizePlanContent(content)
	}

	doc, err := w.parseDocument()
	if err != nil {
		return err
	}

	if err := doc.UpdateSectionContent(strings.ToUpper(section), content); err != nil {
		return fmt.Errorf("updating section %q: %w", section, err)
	}

	return w.writeDocument(doc)
}

// UpdateList updates a list section (scope, plan) with optional append mode.
func (w *Writer) UpdateList(section string, items []string, appendMode bool) error {
	doc, err := w.parseDocument()
	if err != nil {
		return err
	}

	sectionName := strings.ToUpper(section)

	// Format items based on section type
	// Plan sections use checkboxes: "- [ ] item"
	// Other sections use plain lists: "- item"
	formattedItems := items
	if section == "plan" {
		formattedItems = make([]string, len(items))
		for i, item := range items {
			formattedItems[i] = "[ ] " + item
		}
	}

	var updateErr error
	if appendMode {
		updateErr = doc.AppendSectionList(sectionName, formattedItems)
	} else {
		updateErr = doc.UpdateSectionList(sectionName, formattedItems)
	}

	if updateErr != nil {
		return fmt.Errorf("updating %q section: %w", section, updateErr)
	}

	return w.writeDocument(doc)
}

// MarkPlanStepComplete marks a specific plan step as completed and optionally logs a message.
func (w *Writer) MarkPlanStepComplete(step int, status, message string) error {
	mission, err := NewReader(w.FS(), w.MissionPath()).Read()
	if err != nil {
		return fmt.Errorf("reading mission: %w", err)
	}

	// Handle logging if message is provided
	if message != "" {
		w.logPlanStep(mission.ID, step, status, message)
	}

	doc, err := w.parseDocument()
	if err != nil {
		return err
	}

	// Get the raw section content to preserve checkbox states
	planContent, err := doc.GetSection("PLAN")
	if err != nil {
		return fmt.Errorf("getting plan section: %w", err)
	}

	updatedContent, err := w.markStepComplete(planContent, step)
	if err != nil {
		return err
	}

	if err := doc.UpdateSectionContent("PLAN", updatedContent); err != nil {
		return fmt.Errorf("updating plan section: %w", err)
	}

	return w.writeDocument(doc)
}

// parseDocument reads and parses the mission document.
func (w *Writer) parseDocument() (*md.Document, error) {
	data, err := afero.ReadFile(w.FS(), w.MissionPath())
	if err != nil {
		return nil, fmt.Errorf("reading mission file: %w", err)
	}

	doc, err := md.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parsing mission: %w", err)
	}

	return doc, nil
}

// writeDocument writes the document back to the mission file.
func (w *Writer) writeDocument(doc *md.Document) error {
	output, err := doc.Write()
	if err != nil {
		return fmt.Errorf("writing document: %w", err)
	}

	return afero.WriteFile(w.FS(), w.MissionPath(), output, 0644)
}

// logPlanStep logs a plan step completion message.
func (w *Writer) logPlanStep(missionID string, step int, status, message string) {
	if status == "" {
		status = "INFO"
	}
	var log *logger.Logger
	if w.loggerConfig != nil {
		log = logger.NewWithConfig(missionID, w.loggerConfig)
	} else {
		log = logger.New(missionID)
	}
	log.LogStep(status, fmt.Sprintf("Plan Step %d", step), message)
}

// markStepComplete marks a specific step as complete in the plan content.
// Expects plan content to be in checkbox format (normalized by UpdateSection/UpdateList).
func (w *Writer) markStepComplete(planContent string, step int) (string, error) {
	lines := strings.Split(planContent, "\n")
	planStepCount := 0

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- [ ] ") || strings.HasPrefix(trimmed, "- [x] ") {
			planStepCount++
			if planStepCount == step {
				lines[i] = strings.Replace(line, "- [ ]", "- [x]", 1)
			}
		}
	}

	if step < 1 || step > planStepCount {
		return "", fmt.Errorf("step %d not found (total steps: %d)", step, planStepCount)
	}

	return strings.Join(lines, "\n"), nil
}

// UpdateFrontmatter updates frontmatter fields.
func (w *Writer) UpdateFrontmatter(pairs []string) error {
	mission, err := NewReader(w.FS(), w.MissionPath()).Read()
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
			mission.Domains = append(mission.Domains, value)
		}
	}

	return w.Write(mission)
}

// format converts a Mission struct to markdown with YAML frontmatter.
// format formats a Mission struct into markdown with YAML frontmatter using pkg/md abstraction.
func (w *Writer) format(mission *Mission) (string, error) {
	// Build frontmatter map with required fields
	frontmatter := map[string]interface{}{
		"id":        mission.ID,
		"type":      mission.Type,
		"track":     mission.Track,
		"iteration": mission.Iteration,
		"status":    mission.Status,
	}

	// Add optional fields if present
	if mission.ParentMission != "" {
		frontmatter["parent_mission"] = mission.ParentMission
	}

	if len(mission.Domains) > 0 {
		frontmatter["domains"] = mission.Domains
	}

	// Use pkg/md to write document with frontmatter
	doc := &md.Document{
		Frontmatter: frontmatter,
		Body:        mission.Body,
	}

	data, err := doc.Write()
	if err != nil {
		return "", fmt.Errorf("writing document: %w", err)
	}

	return string(data), nil
}
