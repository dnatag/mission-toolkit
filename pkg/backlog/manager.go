package backlog

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// BacklogManager handles backlog file operations
type BacklogManager struct {
	missionDir   string
	backlogPath  string
	patternRegex *regexp.Regexp
}

// NewManager creates a new BacklogManager
func NewManager(missionDir string) *BacklogManager {
	return &BacklogManager{
		missionDir:   missionDir,
		backlogPath:  filepath.Join(missionDir, "backlog.md"),
		patternRegex: regexp.MustCompile(`\[PATTERN:([^\]]+)\]\[COUNT:(\d+)\]`),
	}
}

// List returns backlog items, optionally including completed items and filtering by type
func (m *BacklogManager) List(include []string, exclude []string) ([]string, error) {
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

// validateFilters validates include and exclude type filters
func (m *BacklogManager) validateFilters(include, exclude []string) error {
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

// scanBacklogItems scans the backlog file and returns filtered items
func (m *BacklogManager) scanBacklogItems(file *os.File, include, exclude []string) ([]string, error) {
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

// shouldIncludeItem determines if an item should be included based on filters
func (m *BacklogManager) shouldIncludeItem(line string, inCompletedSection bool, currentSection string, include, exclude []string) bool {
	if inCompletedSection {
		return m.shouldIncludeCompleted(include, exclude)
	}
	return m.shouldIncludeTyped(currentSection, include, exclude)
}

// shouldIncludeCompleted checks if completed items should be included
func (m *BacklogManager) shouldIncludeCompleted(include, exclude []string) bool {
	if len(exclude) > 0 && contains(exclude, "completed") {
		return false
	}
	if len(include) > 0 {
		return contains(include, "completed")
	}
	return false
}

// shouldIncludeTyped checks if typed items should be included
func (m *BacklogManager) shouldIncludeTyped(currentSection string, include, exclude []string) bool {
	itemType := m.getSectionType(currentSection)

	if contains(exclude, itemType) {
		return false
	}

	if len(include) > 0 {
		return contains(include, itemType) || contains(include, "completed")
	}
	return true
}

// contains checks if a string slice contains a value
func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// validateType validates the item type
func (m *BacklogManager) validateType(itemType string) error {
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
// If patternID is provided for refactor type, it tracks occurrence count.
func (m *BacklogManager) Add(description, itemType string) error {
	return m.AddWithPattern(description, itemType, "")
}

// AddWithPattern adds a new item with optional pattern ID tracking.
// For refactor items with a patternID, increments count if pattern exists.
func (m *BacklogManager) AddWithPattern(description, itemType, patternID string) error {
	if err := m.validateType(itemType); err != nil {
		return err
	}

	if err := m.ensureBacklogExists(); err != nil {
		return err
	}

	// Pattern ID only applies to refactor type
	patternID = strings.TrimSpace(patternID)
	if patternID != "" && itemType != "refactor" {
		patternID = "" // Ignore pattern ID for non-refactor types
	}

	// If pattern ID provided, check for existing pattern and increment
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
			// Start count at 2 since detecting a pattern means duplication already exists (2+ instances)
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
// This is more efficient than calling Add multiple times when adding multiple items.
func (m *BacklogManager) AddMultiple(descriptions []string, itemType string) error {
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

// Complete marks an item as completed and moves it to the COMPLETED section
func (m *BacklogManager) Complete(itemText string) error {
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

	// Find and remove the item from its current section
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- [ ]") && strings.Contains(trimmed, itemText) {
			// Create completed item with timestamp
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

	// Add to COMPLETED section
	return m.addToCompletedSection(result, completedItem, itemText)
}

// addToCompletedSection adds a completed item to the COMPLETED section
func (m *BacklogManager) addToCompletedSection(lines []string, completedItem string, itemText string) error {
	result := make([]string, 0, len(lines)+1)
	completedSectionFound := false

	for i, line := range lines {
		result = append(result, line)
		if strings.TrimSpace(line) == "## COMPLETED" {
			completedSectionFound = true
			// Find the end of the completed section
			j := i + 1
			for j < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[j]), "## ") && strings.TrimSpace(lines[j]) != "" {
				result = append(result, lines[j])
				j++
			}
			// Add the completed item
			result = append(result, completedItem)
			// Add remaining lines
			result = append(result, lines[j:]...)
			action := fmt.Sprintf("Completed item: %s", itemText)
			return m.writeBacklogWithMetadata(strings.Join(result, "\n"), action)
		}
	}

	// If COMPLETED section not found, create it at the end
	if !completedSectionFound {
		result = append(result, "", "## COMPLETED", "(History of completed backlog items)", completedItem)
		action := fmt.Sprintf("Completed item: %s", itemText)
		return m.writeBacklogWithMetadata(strings.Join(result, "\n"), action)
	}

	return fmt.Errorf("COMPLETED section not found in backlog")
}

// ensureBacklogExists creates the backlog file if it doesn't exist
func (m *BacklogManager) ensureBacklogExists() error {
	if _, err := os.Stat(m.backlogPath); os.IsNotExist(err) {
		return m.createBacklogFile()
	}
	return nil
}

// createBacklogFile creates a new backlog file with the standard structure
func (m *BacklogManager) createBacklogFile() error {
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

// readBacklogContent reads the entire backlog file content
func (m *BacklogManager) readBacklogContent() (string, error) {
	content, err := os.ReadFile(m.backlogPath)
	if err != nil {
		return "", fmt.Errorf("reading backlog file: %w", err)
	}
	return string(content), nil
}

// writeBacklogContent writes content to the backlog file
func (m *BacklogManager) writeBacklogContent(content string) error {
	return m.writeBacklogWithMetadata(content, "")
}

// Cleanup removes completed items from the COMPLETED section of the backlog.
// If itemType is provided, only removes completed items that match that type.
// If itemType is empty, removes all completed items.
// Returns the number of items removed.
//
// Type matching is heuristic-based:
//   - "decomposed": matches items containing "(from Epic:" marker
//   - "refactor": matches items containing "refactor" or "extract" (case-insensitive)
//   - "future": cannot be reliably identified (returns 0 matches)
func (m *BacklogManager) Cleanup(itemType string) (int, error) {
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

		// Check if this is a completed item in the COMPLETED section
		if inCompletedSection && strings.HasPrefix(trimmed, "- [x]") {
			if itemType == "" {
				// Remove all completed items
				removedCount++
				continue
			}

			// Filter by type using markers
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
