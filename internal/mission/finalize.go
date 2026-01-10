package mission

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

// FinalizeService validates mission.md completeness
type FinalizeService struct {
	fs  afero.Fs
	dir string
}

// FinalizeResult represents validation result
type FinalizeResult struct {
	Valid           bool     `json:"valid"`
	MissingSections []string `json:"missing_sections,omitempty"`
	EmptySections   []string `json:"empty_sections,omitempty"`
	Message         string   `json:"message"`
}

// NewFinalizeService creates a new FinalizeService
func NewFinalizeService(fs afero.Fs, dir string) *FinalizeService {
	return &FinalizeService{
		fs:  fs,
		dir: dir,
	}
}

// Finalize validates mission.md completeness and returns JSON result
func (s *FinalizeService) Finalize() (string, error) {
	missionPath := filepath.Join(s.dir, "mission.md")

	// Read mission file
	reader := NewReader(s.fs)
	m, err := reader.Read(missionPath)
	if err != nil {
		return "", fmt.Errorf("reading mission file: %w", err)
	}

	// Validate sections
	result := s.validateSections(m)

	// Convert to JSON
	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("formatting output: %w", err)
	}

	return string(jsonOutput), nil
}

// validateSections checks that all required sections exist and are not empty
func (s *FinalizeService) validateSections(m *Mission) *FinalizeResult {
	body := m.Body
	requiredSections := []string{"INTENT", "SCOPE", "PLAN", "VERIFICATION"}

	result := &FinalizeResult{
		Valid:           true,
		MissingSections: []string{},
		EmptySections:   []string{},
	}

	for _, section := range requiredSections {
		sectionHeader := "## " + section
		if !strings.Contains(body, sectionHeader) {
			result.MissingSections = append(result.MissingSections, section)
			result.Valid = false
			continue
		}

		// Check if section is empty
		if s.isSectionEmpty(body, section) {
			result.EmptySections = append(result.EmptySections, section)
			result.Valid = false
		}
	}

	if result.Valid {
		result.Message = "Mission validated successfully"
	} else {
		result.Message = "Mission validation failed"
	}

	return result
}

// isSectionEmpty checks if a section has no content
func (s *FinalizeService) isSectionEmpty(body, sectionName string) bool {
	lines := strings.Split(body, "\n")
	inSection := false
	hasContent := false

	for _, line := range lines {
		if strings.HasPrefix(line, "## "+sectionName) {
			inSection = true
			continue
		}
		if inSection && strings.HasPrefix(line, "## ") {
			break
		}
		if inSection && strings.TrimSpace(line) != "" {
			hasContent = true
			break
		}
	}

	return !hasContent
}
