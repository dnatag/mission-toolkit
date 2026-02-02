package backlog

import (
	"fmt"
	"strings"
)

// getSectionType extracts the type from a section header
func (m *BacklogManager) getSectionType(section string) string {
	switch section {
	case "## DECOMPOSED INTENTS":
		return "decomposed"
	case "## REFACTORING OPPORTUNITIES":
		return "refactor"
	case "## FUTURE ENHANCEMENTS":
		return "future"
	case "## FEATURES":
		return "feature"
	case "## BUGFIXES":
		return "bugfix"
	default:
		return ""
	}
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
	case "feature":
		return "## FEATURES"
	case "bugfix":
		return "## BUGFIXES"
	default:
		return ""
	}
}

// isInSection checks if the current section matches the item type
func (m *BacklogManager) isInSection(sectionHeader, itemType string) bool {
	expectedHeader := m.getSectionHeader(itemType)
	return sectionHeader == expectedHeader
}

// findAndModifySection finds a section by header and applies a modifier function to insert items.
// Returns the modified lines or an error if the section is not found.
func (m *BacklogManager) findAndModifySection(lines []string, sectionHeader string, modifier func() []string) ([]string, error) {
	result := make([]string, 0, len(lines)+10)

	for i, line := range lines {
		result = append(result, line)

		if strings.TrimSpace(line) == sectionHeader {
			// Find the end of this section
			j := i + 1
			for j < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[j]), "## ") && strings.TrimSpace(lines[j]) != "" {
				result = append(result, lines[j])
				j++
			}
			// Apply modifier to get items to insert
			newItems := modifier()
			result = append(result, newItems...)
			// Add remaining lines
			result = append(result, lines[j:]...)
			return result, nil
		}
	}

	return nil, fmt.Errorf("section %s not found in backlog", sectionHeader)
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
