package backlog

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// patternRegex matches [PATTERN:id][COUNT:n] format
var patternRegex = regexp.MustCompile(`\[PATTERN:([^\]]+)\]\[COUNT:(\d+)\]`)

// GetPatternCount returns the occurrence count for a pattern ID.
// Returns 0 if pattern not found.
func (m *BacklogManager) GetPatternCount(patternID string) (int, error) {
	if err := m.ensureBacklogExists(); err != nil {
		return 0, err
	}

	body, _, err := m.readBacklogWithMetadata()
	if err != nil {
		return 0, err
	}

	for _, line := range strings.Split(body, "\n") {
		matches := patternRegex.FindStringSubmatch(line)
		if len(matches) == 3 && matches[1] == patternID {
			count, _ := strconv.Atoi(matches[2])
			return count, nil
		}
	}
	return 0, nil
}

// incrementPatternCount increments the count for an existing pattern ID.
func (m *BacklogManager) incrementPatternCount(patternID string) error {
	body, _, err := m.readBacklogWithMetadata()
	if err != nil {
		return err
	}

	lines := strings.Split(body, "\n")
	for i, line := range lines {
		matches := patternRegex.FindStringSubmatch(line)
		if len(matches) == 3 && matches[1] == patternID {
			count, _ := strconv.Atoi(matches[2])
			newCount := count + 1
			lines[i] = patternRegex.ReplaceAllString(line, fmt.Sprintf("[PATTERN:%s][COUNT:%d]", patternID, newCount))
			action := fmt.Sprintf("Incremented pattern %s count to %d", patternID, newCount)
			return m.writeBacklogWithMetadata(strings.Join(lines, "\n"), action)
		}
	}
	return fmt.Errorf("pattern not found: %s", patternID)
}
