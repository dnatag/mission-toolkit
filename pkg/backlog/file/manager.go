package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Manager handles backlog file operations using markdown.
type Manager struct {
	missionDir   string
	backlogPath  string
	patternRegex *regexp.Regexp
}

// NewManager creates a new Manager for the given mission directory.
func NewManager(missionDir string) *Manager {
	return &Manager{
		missionDir:   missionDir,
		backlogPath:  filepath.Join(missionDir, "backlog.md"),
		patternRegex: regexp.MustCompile(`\[PATTERN:([^\]]+)\]\[COUNT:(\d+)\]`),
	}
}

// List returns backlog items, optionally including completed items and filtering by type.
func (m *Manager) List(include []string, exclude []string) ([]string, error) {
	if err := m.validateFilters(include, exclude); err != nil {
		return nil, err
	}

	if err := m.ensureBacklogExists(); err != nil {
		return nil, err
	}

	file, err := os.Open(m.backlogPath)
	if err != nil {
		return nil, fmt.Errorf("opening backlog file: %w", err)
	}
	defer file.Close()

	return m.scanBacklogItems(file, include, exclude)
}

func (m *Manager) validateFilters(include, exclude []string) error {
	for _, t := range include {
		if t != "completed" {
			if err := m.validateType(t); err != nil {
				return err
			}
		}
	}
	for _, t := range exclude {
		if t != "completed" {
			if err := m.validateType(t); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *Manager) scanBacklogItems(file *os.File, include, exclude []string) ([]string, error) {
	var items []string
	scanner := bufio.NewScanner(file)
	inCompletedSection := false
	currentSection := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "## COMPLETED" {
			inCompletedSection = true
			currentSection = ""
			continue
		}

		if strings.HasPrefix(line, "## ") {
			inCompletedSection = false
			currentSection = line
			continue
		}

		if strings.HasPrefix(line, "- [ ]") || strings.HasPrefix(line, "- [x]") {
			if m.shouldIncludeItem(line, inCompletedSection, currentSection, include, exclude) {
				items = append(items, line)
			}
		}
	}

	return items, scanner.Err()
}

func (m *Manager) shouldIncludeItem(line string, inCompletedSection bool, currentSection string, include, exclude []string) bool {
	if inCompletedSection {
		return m.shouldIncludeCompleted(include, exclude)
	}
	return m.shouldIncludeTyped(currentSection, include, exclude)
}

func (m *Manager) shouldIncludeCompleted(include, exclude []string) bool {
	if len(exclude) > 0 && contains(exclude, "completed") {
		return false
	}
	if len(include) > 0 {
		return contains(include, "completed")
	}
	return false
}

func (m *Manager) shouldIncludeTyped(currentSection string, include, exclude []string) bool {
	itemType := m.getSectionType(currentSection)

	if contains(exclude, itemType) {
		return false
	}

	if len(include) > 0 {
		return contains(include, itemType) || contains(include, "completed")
	}
	return true
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func (m *Manager) validateType(itemType string) error {
	validTypes := map[string]bool{
		"decomposed": true,
		"refactor":   true,
		"future":     true,
		"feature":    true,
		"bugfix":     true,
	}

	if !validTypes[itemType] {
		return fmt.Errorf("invalid type: %s. Valid types: decomposed, refactor, future, feature, bugfix", itemType)
	}
	return nil
}

// Add adds a new item to the specified section.
func (m *Manager) Add(description, itemType string) error {
	return m.AddWithPattern(description, itemType, "")
}

// AddWithPattern adds a new item with optional pattern ID tracking.
func (m *Manager) AddWithPattern(description, itemType, patternID string) error {
	if err := m.validateType(itemType); err != nil {
		return err
	}

	if err := m.ensureBacklogExists(); err != nil {
		return err
	}

	patternID = strings.TrimSpace(patternID)
	if patternID != "" && itemType != "refactor" {
		patternID = ""
	}

	if patternID != "" {
		count, err := m.GetPatternCount(patternID)
		if err != nil {
			return err
		}
		if count > 0 {
			return m.incrementPatternCount(patternID)
		}
	}

	body, _, err := m.readBacklogWithMetadata()
	if err != nil {
		return err
	}

	sectionHeader := m.getSectionHeader(itemType)
	lines := strings.Split(body, "\n")

	result, err := m.findAndModifySection(lines, sectionHeader, func() []string {
		if patternID != "" {
			return []string{fmt.Sprintf("- [ ] [PATTERN:%s][COUNT:2] %s", patternID, description)}
		}
		return []string{fmt.Sprintf("- [ ] %s", description)}
	})
	if err != nil {
		return err
	}

	action := fmt.Sprintf("Added %s item: %s", itemType, description)
	return m.writeBacklogWithMetadata(strings.Join(result, "\n"), action)
}

// AddMultiple adds multiple items to the specified section in a single operation.
func (m *Manager) AddMultiple(descriptions []string, itemType string) error {
	if err := m.validateType(itemType); err != nil {
		return err
	}

	if err := m.ensureBacklogExists(); err != nil {
		return err
	}

	body, _, err := m.readBacklogWithMetadata()
	if err != nil {
		return err
	}

	sectionHeader := m.getSectionHeader(itemType)
	lines := strings.Split(body, "\n")

	result, err := m.findAndModifySection(lines, sectionHeader, func() []string {
		items := make([]string, len(descriptions))
		for i, desc := range descriptions {
			items[i] = fmt.Sprintf("- [ ] %s", desc)
		}
		return items
	})
	if err != nil {
		return err
	}

	action := fmt.Sprintf("Added %d %s items", len(descriptions), itemType)
	return m.writeBacklogWithMetadata(strings.Join(result, "\n"), action)
}

// Complete marks an item as completed and moves it to the COMPLETED section.
func (m *Manager) Complete(itemText string) error {
	if err := m.ensureBacklogExists(); err != nil {
		return err
	}

	body, _, err := m.readBacklogWithMetadata()
	if err != nil {
		return err
	}

	lines := strings.Split(body, "\n")
	result := make([]string, 0, len(lines))
	var completedItem string
	itemFound := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- [ ]") && strings.Contains(trimmed, itemText) {
			timestamp := time.Now().Format("2006-01-02")
			itemDesc := strings.TrimPrefix(trimmed, "- [ ] ")
			completedItem = fmt.Sprintf("- [x] %s (Completed: %s)", itemDesc, timestamp)
			itemFound = true
			continue
		}
		result = append(result, line)
	}

	if !itemFound {
		return fmt.Errorf("item not found: %s", itemText)
	}

	return m.addToCompletedSection(result, completedItem, itemText)
}

func (m *Manager) addToCompletedSection(lines []string, completedItem string, itemText string) error {
	result := make([]string, 0, len(lines)+1)

	for i, line := range lines {
		result = append(result, line)
		if strings.TrimSpace(line) == "## COMPLETED" {
			j := i + 1
			for j < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[j]), "## ") && strings.TrimSpace(lines[j]) != "" {
				result = append(result, lines[j])
				j++
			}
			result = append(result, completedItem)
			result = append(result, lines[j:]...)
			action := fmt.Sprintf("Completed item: %s", itemText)
			return m.writeBacklogWithMetadata(strings.Join(result, "\n"), action)
		}
	}

	// COMPLETED section not found, create it
	result = append(result, "", "## COMPLETED", "(History of completed backlog items)", completedItem)
	action := fmt.Sprintf("Completed item: %s", itemText)
	return m.writeBacklogWithMetadata(strings.Join(result, "\n"), action)
}

func (m *Manager) ensureBacklogExists() error {
	if _, err := os.Stat(m.backlogPath); os.IsNotExist(err) {
		return m.createBacklogFile()
	}
	return nil
}

func (m *Manager) createBacklogFile() error {
	template := `# Mission Backlog

## FEATURES
*User-defined feature requests and enhancements.*

## BUGFIXES
*Bug reports and issues to be fixed.*

## DECOMPOSED INTENTS
*Atomic tasks broken down from larger epics.*

## REFACTORING OPPORTUNITIES
*Technical debt and refactoring opportunities identified during development.*

## FUTURE ENHANCEMENTS
*Ideas and future feature requests for later consideration.*

## COMPLETED
*History of completed backlog items.*
`
	return m.writeBacklogContent(template)
}

func (m *Manager) readBacklogContent() (string, error) {
	content, err := os.ReadFile(m.backlogPath)
	if err != nil {
		return "", fmt.Errorf("reading backlog file: %w", err)
	}
	return string(content), nil
}

func (m *Manager) writeBacklogContent(content string) error {
	return m.writeBacklogWithMetadata(content, "")
}

// Cleanup removes completed items from the COMPLETED section.
func (m *Manager) Cleanup(itemType string) (int, error) {
	if itemType != "" {
		if err := m.validateType(itemType); err != nil {
			return 0, err
		}
	}

	if err := m.ensureBacklogExists(); err != nil {
		return 0, err
	}

	body, _, err := m.readBacklogWithMetadata()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(body, "\n")
	result := make([]string, 0, len(lines))
	inCompletedSection := false
	removedCount := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "## COMPLETED" {
			inCompletedSection = true
			result = append(result, line)
			continue
		}

		if strings.HasPrefix(trimmed, "## ") {
			inCompletedSection = false
			result = append(result, line)
			continue
		}

		if inCompletedSection && strings.HasPrefix(trimmed, "- [x]") {
			if itemType == "" {
				removedCount++
				continue
			}
			if m.matchesItemType(trimmed, itemType) {
				removedCount++
				continue
			}
		}

		result = append(result, line)
	}

	if removedCount > 0 {
		action := fmt.Sprintf("Cleaned up %d completed items", removedCount)
		if itemType != "" {
			action = fmt.Sprintf("Cleaned up %d completed %s items", removedCount, itemType)
		}
		if err := m.writeBacklogWithMetadata(strings.Join(result, "\n"), action); err != nil {
			return 0, err
		}
	}

	return removedCount, nil
}

// Decompose adds multiple sub-intents with dependency tracking.
// Falls back to AddMultiple with dependency hints in descriptions.
func (m *Manager) Decompose(jsonInput string) error {
	var decompose struct {
		Action     string `json:"action"`
		SubIntents []struct {
			Intent         string   `json:"intent"`
			Rationale      string   `json:"rationale"`
			EstimatedFiles int      `json:"estimated_files"`
			Dependencies   []string `json:"dependencies"`
		} `json:"sub_intents"`
		DecompositionRationale string `json:"decomposition_rationale"`
	}

	if err := json.Unmarshal([]byte(jsonInput), &decompose); err != nil {
		return fmt.Errorf("parsing decompose JSON: %w", err)
	}

	if len(decompose.SubIntents) == 0 {
		return fmt.Errorf("no sub-intents found in decompose input")
	}

	descriptions := make([]string, len(decompose.SubIntents))
	for i, subIntent := range decompose.SubIntents {
		description := subIntent.Intent
		if len(subIntent.Dependencies) > 0 {
			description += fmt.Sprintf(" (depends on: %s)", strings.Join(subIntent.Dependencies, ", "))
		}
		descriptions[i] = description
	}

	return m.AddMultiple(descriptions, "decomposed")
}
