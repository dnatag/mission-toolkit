package mission

import (
	"fmt"
	"path/filepath"
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

// UpdateSection updates a text section (intent, verification).
// UpdateSection updates a text section (intent, verification).
// Ensures exactly one empty line between sections for consistent formatting.
func (w *Writer) UpdateSection(section string, content string) error {
	mission, err := NewReader(w.FS(), w.MissionPath()).Read()
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
			result = append(result, line, content, "")

			// Skip old content until next section or end
			for j := i + 1; j < len(lines); j++ {
				if strings.HasPrefix(strings.TrimSpace(lines[j]), "## ") {
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
		result = append(result, "", sectionHeader, content)
	}

	mission.Body = strings.Join(result, "\n")
	return w.Write(mission)
}

// UpdateList updates a list section (scope, plan) with optional append mode.
// Ensures exactly one empty line between sections for consistent formatting.
func (w *Writer) UpdateList(section string, items []string, appendMode bool) error {
	mission, err := NewReader(w.FS(), w.MissionPath()).Read()
	if err != nil {
		return fmt.Errorf("reading mission: %w", err)
	}

	lines := strings.Split(mission.Body, "\n")
	var result []string
	sectionHeader := "## " + strings.ToUpper(section)
	foundSection := false
	lineIndex := 0

	for lineIndex < len(lines) {
		currentLine := lines[lineIndex]

		if strings.TrimSpace(currentLine) == sectionHeader {
			result = append(result, currentLine)

			if appendMode {
				existingItems := w.extractExistingItems(lines, lineIndex+1)
				result = append(result, existingItems...)
			}

			w.addFormattedItems(&result, section, items)

			foundSection = true
			nextSectionIndex := w.skipSectionContent(lines, lineIndex+1)
			if nextSectionIndex < len(lines) {
				result = append(result, "")
			}
			lineIndex = nextSectionIndex
			continue
		}

		result = append(result, currentLine)
		lineIndex++
	}

	if !foundSection {
		result = append(result, "", sectionHeader)
		w.addFormattedItems(&result, section, items)
	}

	mission.Body = strings.Join(result, "\n")
	return w.Write(mission)
}

// extractExistingItems collects existing items from a section.
func (w *Writer) extractExistingItems(lines []string, startIndex int) []string {
	var existingItems []string
	for j := startIndex; j < len(lines); j++ {
		if strings.HasPrefix(strings.TrimSpace(lines[j]), "## ") {
			break
		}
		if strings.TrimSpace(lines[j]) != "" {
			existingItems = append(existingItems, lines[j])
		}
	}
	return existingItems
}

// addFormattedItems adds items with proper formatting based on section type.
func (w *Writer) addFormattedItems(result *[]string, section string, items []string) {
	if section == "plan" {
		for _, item := range items {
			*result = append(*result, "- [ ] "+item)
		}
	} else {
		*result = append(*result, items...)
	}
}

// skipSectionContent skips content until the next section and returns the index.
// This ensures proper section boundary handling when updating list sections.
func (w *Writer) skipSectionContent(lines []string, startIndex int) int {
	for j := startIndex; j < len(lines); j++ {
		if strings.HasPrefix(strings.TrimSpace(lines[j]), "## ") {
			// Found next section header, return its index
			return j
		}
	}
	// No more sections found, return end of file
	return len(lines)
}

// MarkPlanStepComplete marks a specific plan step as completed and optionally logs a message.
func (w *Writer) MarkPlanStepComplete(step int, status, message string) error {
	mission, err := NewReader(w.FS(), w.MissionPath()).Read()
	if err != nil {
		return fmt.Errorf("reading mission: %w", err)
	}

	// Handle logging if message is provided
	if message != "" {
		if status == "" {
			status = "INFO"
		}
		var log *logger.Logger
		if w.loggerConfig != nil {
			log = logger.NewWithConfig(mission.ID, w.loggerConfig)
		} else {
			log = logger.New(mission.ID)
		}
		log.LogStep(status, fmt.Sprintf("Plan Step %d", step), message)
	}

	lines := strings.Split(mission.Body, "\n")
	var result []string
	inPlan := false
	planStepCount := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "## PLAN" {
			inPlan = true
			result = append(result, line)
			continue
		} else if strings.HasPrefix(trimmed, "## ") {
			inPlan = false
		}

		if inPlan && strings.HasPrefix(trimmed, "- [") {
			planStepCount++
			if planStepCount == step {
				line = strings.Replace(line, "- [ ]", "- [x]", 1)
			}
		}

		result = append(result, line)
	}

	if step > planStepCount {
		return fmt.Errorf("step %d not found (total steps: %d)", step, planStepCount)
	}

	mission.Body = strings.Join(result, "\n")
	return w.Write(mission)
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
