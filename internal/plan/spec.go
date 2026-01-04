package plan

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/afero"
)

var fs = afero.NewOsFs()

// FileAction represents the action to be performed on a file
type FileAction string

const (
	FileActionModify FileAction = "modify"
	FileActionCreate FileAction = "create"
)

// FileSpec represents a file with its associated action
type FileSpec struct {
	Path   string     `json:"path"`
	Action FileAction `json:"action"`
}

// PlanSpec represents the structure for mission planning data
type PlanSpec struct {
	Intent                 string     `json:"intent"`
	Type                   string     `json:"type,omitempty"`  // WET or DRY
	Scope                  []string   `json:"scope,omitempty"` // Legacy field for backward compatibility
	Files                  []FileSpec `json:"files,omitempty"` // New field with action annotations
	Domain                 []string   `json:"domain,omitempty"`
	Track                  string     `json:"track,omitempty"` // TRACK 1-4
	Plan                   []string   `json:"plan,omitempty"`
	Verification           string     `json:"verification,omitempty"`
	ClarificationQuestions []string   `json:"clarification_questions,omitempty"` // For clarification missions
}

// GetScopeFiles returns all file paths from both legacy scope and new files fields
func (p *PlanSpec) GetScopeFiles() []string {
	var files []string

	// Add files from new Files field
	for _, file := range p.Files {
		files = append(files, file.Path)
	}

	// Add files from legacy Scope field for backward compatibility
	for _, file := range p.Scope {
		// Only add if not already in Files
		found := false
		for _, existing := range files {
			if existing == file {
				found = true
				break
			}
		}
		if !found {
			files = append(files, file)
		}
	}

	return files
}

// GetFileAction returns the action for a specific file path
func (p *PlanSpec) GetFileAction(filePath string) FileAction {
	for _, file := range p.Files {
		if file.Path == filePath {
			return file.Action
		}
	}
	// Default to modify for backward compatibility
	return FileActionModify
}

// Validate checks if the PlanSpec has required fields
func (p *PlanSpec) Validate() error {
	if p.Intent == "" {
		return fmt.Errorf("intent is required")
	}

	// Check if we have files in either scope or files field
	allFiles := p.GetScopeFiles()
	if len(allFiles) == 0 {
		return fmt.Errorf("scope or files is required")
	}

	if len(p.Plan) == 0 {
		return fmt.Errorf("plan is required")
	}
	if p.Verification == "" {
		return fmt.Errorf("verification is required")
	}
	return nil
}

// Read reads a PlanSpec from a JSON file
func Read(path string) (*PlanSpec, error) {
	data, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, fmt.Errorf("reading plan file: %w", err)
	}

	var spec PlanSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("parsing plan file: %w", err)
	}

	return &spec, nil
}

// Write writes a PlanSpec to a JSON file
func (p *PlanSpec) Write(path string) error {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling plan: %w", err)
	}

	if err := afero.WriteFile(fs, path, data, 0644); err != nil {
		return fmt.Errorf("writing plan file: %w", err)
	}

	return nil
}
