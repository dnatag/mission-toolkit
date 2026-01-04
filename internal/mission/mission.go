package mission

import "strings"

// Mission represents a mission with YAML frontmatter metadata and markdown body
type Mission struct {
	// Frontmatter fields
	ID            string `yaml:"id"`
	Type          string `yaml:"type"`
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
