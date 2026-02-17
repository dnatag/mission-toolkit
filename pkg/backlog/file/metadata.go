package file

import (
	"fmt"
	"os"
	"time"

	"github.com/dnatag/mission-toolkit/pkg/md"
	"gopkg.in/yaml.v3"
)

// BacklogMetadata represents the frontmatter metadata for backlog.md.
type BacklogMetadata struct {
	LastUpdated time.Time `yaml:"last_updated"`
	LastAction  string    `yaml:"last_action"`
}

func (m *Manager) readBacklogWithMetadata() (string, *BacklogMetadata, error) {
	content, err := os.ReadFile(m.backlogPath)
	if err != nil {
		return "", nil, fmt.Errorf("reading backlog file: %w", err)
	}

	doc, err := md.Parse(content)
	if err != nil {
		return "", nil, fmt.Errorf("parsing backlog: %w", err)
	}

	var metadata BacklogMetadata
	if len(doc.Frontmatter) > 0 {
		yamlData, err := yaml.Marshal(doc.Frontmatter)
		if err != nil {
			return "", nil, fmt.Errorf("marshaling frontmatter: %w", err)
		}
		if err := yaml.Unmarshal(yamlData, &metadata); err != nil {
			return "", nil, fmt.Errorf("unmarshaling frontmatter: %w", err)
		}
	}

	return doc.Body, &metadata, nil
}

func (m *Manager) writeBacklogWithMetadata(content string, action string) error {
	if err := os.MkdirAll(m.missionDir, 0755); err != nil {
		return fmt.Errorf("creating mission directory: %w", err)
	}

	frontmatter := map[string]interface{}{
		"last_updated": time.Now(),
		"last_action":  action,
	}

	doc := &md.Document{
		Frontmatter: frontmatter,
		Body:        content,
	}

	data, err := doc.Write()
	if err != nil {
		return fmt.Errorf("writing document: %w", err)
	}

	if err := os.WriteFile(m.backlogPath, data, 0644); err != nil {
		return fmt.Errorf("writing backlog file: %w", err)
	}
	return nil
}
