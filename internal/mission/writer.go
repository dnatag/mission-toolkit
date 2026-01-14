package mission

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dnatag/mission-toolkit/internal/logger"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// Writer handles writing mission files and updating status.
type Writer struct {
	fs   afero.Fs
	path string
}

// NewWriter creates a new mission writer rooted at missionDir.
// The mission file path is always <missionDir>/mission.md.
func NewWriter(fs afero.Fs, missionDir string) *Writer {
	return &Writer{
		fs:   fs,
		path: filepath.Join(missionDir, "mission.md"),
	}
}

// NewWriterWithPath creates a new mission writer for an explicit mission file path.
func NewWriterWithPath(fs afero.Fs, path string) *Writer {
	return &Writer{
		fs:   fs,
		path: path,
	}
}

// Write writes a Mission struct to the writer's mission file with YAML frontmatter.
func (w *Writer) Write(mission *Mission) error {
	content, err := w.format(mission)
	if err != nil {
		return err
	}

	return afero.WriteFile(w.fs, w.path, []byte(content), 0644)
}

// UpdateStatus updates the status field in the mission file while preserving the body.
func (w *Writer) UpdateStatus(newStatus string) error {
	mission, err := NewReader(w.fs).Read(w.path)
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
func (w *Writer) UpdateSection(section string, content string) error {
	mission, err := NewReader(w.fs).Read(w.path)
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
		// Add new section at end
		result = append(result, "", sectionHeader)
		result = append(result, content)
	}

	mission.Body = strings.Join(result, "\n")
	return w.Write(mission)
}

// UpdateList updates a list section (scope, plan) with optional append mode.
// This method preserves the structure of the mission file by properly handling
// section boundaries and ensuring subsequent sections remain intact.
func (w *Writer) UpdateList(section string, items []string, appendMode bool) error {
	mission, err := NewReader(w.fs).Read(w.path)
	if err != nil {
		return fmt.Errorf("reading mission: %w", err)
	}

	lines := strings.Split(mission.Body, "\n")
	var result []string
	sectionHeader := "## " + strings.ToUpper(section)
	foundSection := false
	lineIndex := 0

	// Process each line, looking for the target section
	for lineIndex < len(lines) {
		currentLine := lines[lineIndex]

		if strings.TrimSpace(currentLine) == sectionHeader {
			// Found the target section header
			result = append(result, currentLine)

			if appendMode {
				// Preserve existing items and add new ones
				existingItems := w.extractExistingItems(lines, lineIndex+1)
				result = append(result, existingItems...)
			}

			// Add new items with appropriate formatting
			w.addFormattedItems(&result, section, items)

			foundSection = true
			// Skip existing content in this section and continue processing
			lineIndex = w.skipSectionContent(lines, lineIndex+1)
			continue
		}

		// Copy non-target lines as-is
		result = append(result, currentLine)
		lineIndex++
	}

	if !foundSection {
		// Add new section at end if it doesn't exist
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
	mission, err := NewReader(w.fs).Read(w.path)
	if err != nil {
		return fmt.Errorf("reading mission: %w", err)
	}

	// Handle logging if message is provided
	if message != "" {
		if status == "" {
			status = "INFO"
		}
		log := logger.New(mission.ID)
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

		if inPlan && strings.HasPrefix(trimmed, "- [ ]") {
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
	mission, err := NewReader(w.fs).Read(w.path)
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
			mission.Domains = value
		}
	}

	return w.Write(mission)
}

// format converts a Mission struct to markdown with YAML frontmatter.
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

	if mission.Domains != "" {
		frontmatter["domains"] = mission.Domains
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
