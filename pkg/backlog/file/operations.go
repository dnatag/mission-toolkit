package file

import (
	"fmt"
	"strings"
)

func (m *Manager) getSectionType(section string) string {
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

func (m *Manager) getSectionHeader(itemType string) string {
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

func (m *Manager) isInSection(sectionHeader, itemType string) bool {
	expectedHeader := m.getSectionHeader(itemType)
	return sectionHeader == expectedHeader
}

func (m *Manager) findAndModifySection(lines []string, sectionHeader string, modifier func() []string) ([]string, error) {
	result := make([]string, 0, len(lines)+10)

	for i, line := range lines {
		result = append(result, line)

		if strings.TrimSpace(line) == sectionHeader {
			j := i + 1
			for j < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[j]), "## ") && strings.TrimSpace(lines[j]) != "" {
				result = append(result, lines[j])
				j++
			}
			newItems := modifier()
			result = append(result, newItems...)
			result = append(result, lines[j:]...)
			return result, nil
		}
	}

	return nil, fmt.Errorf("section %s not found in backlog", sectionHeader)
}

func (m *Manager) matchesItemType(item, itemType string) bool {
	switch itemType {
	case "decomposed":
		return strings.Contains(item, "(from Epic:")
	case "refactor":
		lowerItem := strings.ToLower(item)
		return strings.Contains(lowerItem, "refactor") || strings.Contains(lowerItem, "extract")
	case "future":
		return false
	default:
		return false
	}
}
