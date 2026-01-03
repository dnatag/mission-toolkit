package mission

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
