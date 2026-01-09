package backlog

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// BacklogManager handles backlog file operations
type BacklogManager struct {
	missionDir  string
	backlogPath string
}

// NewManager creates a new BacklogManager
func NewManager(missionDir string) *BacklogManager {
	return &BacklogManager{
		missionDir:  missionDir,
		backlogPath: filepath.Join(missionDir, "backlog.md"),
	}
}

// List returns backlog items, optionally including completed items
func (m *BacklogManager) List(includeCompleted bool) ([]string, error) {
	if err := m.ensureBacklogExists(); err != nil {
		return nil, err
	}

	file, err := os.Open(m.backlogPath)
	if err != nil {
		return nil, fmt.Errorf("opening backlog file: %w", err)
	}
	defer file.Close()

	var items []string
	scanner := bufio.NewScanner(file)
	inCompletedSection := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "## COMPLETED" {
			inCompletedSection = true
			continue
		}

		if strings.HasPrefix(line, "## ") {
			inCompletedSection = false
			continue
		}

		if strings.HasPrefix(line, "- [ ]") || strings.HasPrefix(line, "- [x]") {
			if strings.HasPrefix(line, "- [x]") && !includeCompleted {
				continue
			}
			if inCompletedSection && !includeCompleted {
				continue
			}
			items = append(items, line)
		}
	}

	return items, scanner.Err()
}

// Add adds a new item to the specified section
func (m *BacklogManager) Add(description, itemType string) error {
	if err := m.validateType(itemType); err != nil {
		return err
	}

	if err := m.ensureBacklogExists(); err != nil {
		return err
	}

	content, err := m.readBacklogContent()
	if err != nil {
		return err
	}

	sectionHeader := m.getSectionHeader(itemType)
	newItem := fmt.Sprintf("- [ ] %s", description)

	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines)+1)

	for i, line := range lines {
		result = append(result, line)

		if strings.TrimSpace(line) == sectionHeader {
			// Find the end of this section
			j := i + 1
			for j < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[j]), "## ") && strings.TrimSpace(lines[j]) != "" {
				result = append(result, lines[j])
				j++
			}
			// Add the new item
			result = append(result, newItem)
			// Add remaining lines
			result = append(result, lines[j:]...)
			return m.writeBacklogContent(strings.Join(result, "\n"))
		}
	}

	return fmt.Errorf("section %s not found in backlog", sectionHeader)
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

	content, err := m.readBacklogContent()
	if err != nil {
		return err
	}

	sectionHeader := m.getSectionHeader(itemType)

	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines)+len(descriptions))

	for i, line := range lines {
		result = append(result, line)

		if strings.TrimSpace(line) == sectionHeader {
			// Find the end of this section
			j := i + 1
			for j < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[j]), "## ") && strings.TrimSpace(lines[j]) != "" {
				result = append(result, lines[j])
				j++
			}
			// Add all new items
			for _, desc := range descriptions {
				result = append(result, fmt.Sprintf("- [ ] %s", desc))
			}
			// Add remaining lines
			result = append(result, lines[j:]...)
			return m.writeBacklogContent(strings.Join(result, "\n"))
		}
	}

	return fmt.Errorf("section %s not found in backlog", sectionHeader)
}

// Complete marks an item as completed and moves it to the COMPLETED section
func (m *BacklogManager) Complete(itemText string) error {
	if err := m.ensureBacklogExists(); err != nil {
		return err
	}

	content, err := m.readBacklogContent()
	if err != nil {
		return err
	}

	lines := strings.Split(content, "\n")
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
	return m.addToCompletedSection(result, completedItem)
}

// addToCompletedSection adds a completed item to the COMPLETED section
func (m *BacklogManager) addToCompletedSection(lines []string, completedItem string) error {
	result := make([]string, 0, len(lines)+1)

	for i, line := range lines {
		result = append(result, line)
		if strings.TrimSpace(line) == "## COMPLETED" {
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
			return m.writeBacklogContent(strings.Join(result, "\n"))
		}
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

## DECOMPOSED INTENTS
*This section lists atomic tasks that were broken down from a larger user-defined epic.*

## REFACTORING OPPORTUNITIES
*This section lists technical debt and refactoring opportunities identified by the AI during planning or execution.*

## FUTURE ENHANCEMENTS
*This section is for user-defined ideas and future feature requests.*

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
	if err := os.MkdirAll(m.missionDir, 0755); err != nil {
		return fmt.Errorf("creating mission directory: %w", err)
	}

	if err := os.WriteFile(m.backlogPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing backlog file: %w", err)
	}
	return nil
}

// validateType validates the item type
func (m *BacklogManager) validateType(itemType string) error {
	validTypes := map[string]bool{
		"decomposed": true,
		"refactor":   true,
		"future":     true,
	}

	if !validTypes[itemType] {
		return fmt.Errorf("invalid type: %s. Valid types: decomposed, refactor, future", itemType)
	}
	return nil
}

// getSectionHeader returns the markdown header for the given type
func (m *BacklogManager) getSectionHeader(itemType string) string {
	switch itemType {
	case "decomposed":
		return "## DECOMPOSED INTENTS"
	case "refactor":
		return "## REFACTORING OPPORTUNITIES"
	case "future":
		return "## FUTURE ENHANCEMENTS"
	default:
		return ""
	}
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

	content, err := m.readBacklogContent()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(content, "\n")
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
		if err := m.writeBacklogContent(strings.Join(result, "\n")); err != nil {
			return 0, err
		}
	}

	return removedCount, nil
}

// matchesItemType checks if a completed item matches the specified type.
// Decomposed items contain "(from Epic:" marker.
// This is a heuristic based on how items are typically added to the backlog.
func (m *BacklogManager) matchesItemType(item, itemType string) bool {
	switch itemType {
	case "decomposed":
		// Decomposed items typically have "(from Epic:" marker
		return strings.Contains(item, "(from Epic:")
	case "refactor":
		// Refactor items typically have "Refactor" or "Extract" in the description
		lowerItem := strings.ToLower(item)
		return strings.Contains(lowerItem, "refactor") || strings.Contains(lowerItem, "extract")
	case "future":
		// Future items don't have specific markers, so we can't reliably identify them
		// This will effectively not match any items unless explicitly marked
		return false
	default:
		return false
	}
}
