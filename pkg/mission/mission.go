package mission

import (
	"regexp"
	"strings"
)

// Mission represents a mission with YAML frontmatter metadata and markdown body
type Mission struct {
	// Frontmatter fields
	ID            string `yaml:"id"`
	Type          string `yaml:"type"`
	Domains       string `yaml:"domains,omitempty"`
	Track         int    `yaml:"track"`
	Iteration     int    `yaml:"iteration"`
	Status        string `yaml:"status"`
	ParentMission string `yaml:"parent_mission,omitempty"`

	// Markdown body (everything after frontmatter)
	Body string
}

// GetScope extracts the list of files from the SCOPE section of the mission body
func (m *Mission) GetScope() []string {
	var scope []string
	lines := strings.Split(m.Body, "\n")
	inScope := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "## ") {
			inScope = strings.EqualFold(trimmed, "## SCOPE")
			continue
		}

		if inScope && trimmed != "" {
			// Remove list markers (- or *)
			clean := strings.TrimLeft(trimmed, "-* ")
			if clean != "" {
				scope = append(scope, clean)
			}
		}
	}
	return scope
}

// GetIntent extracts the intent from mission body
func (m *Mission) GetIntent() string {
	return extractSection(m.Body, "INTENT")
}

// GetPlan extracts the plan steps from mission body
func (m *Mission) GetPlan() []string {
	content := extractSection(m.Body, "PLAN")
	if content == "" {
		return nil
	}

	lines := strings.Split(content, "\n")
	var plan []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			plan = append(plan, trimmed)
		}
	}
	return plan
}

// GetVerification extracts the verification command from mission body
func (m *Mission) GetVerification() string {
	return extractSection(m.Body, "VERIFICATION")
}

// extractSection extracts a section from mission body
func extractSection(body, sectionName string) string {
	re := regexp.MustCompile("(?s)## " + sectionName + "\\s*\\n(.*?)(?:\\n##|$)")
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}
