package diagnosis

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dnatag/mission-toolkit/pkg/md"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

// Diagnosis represents a debug investigation with structured findings
type Diagnosis struct {
	ID         string    `yaml:"id"`
	Status     string    `yaml:"status"`
	Confidence string    `yaml:"confidence"`
	Created    time.Time `yaml:"created"`
	Symptom    string    `yaml:"-"`
	Body       string    `yaml:"-"`
}

// ReadDiagnosis reads and parses a diagnosis.md file using pkg/md abstraction.
// Returns an error if the file doesn't exist or has invalid format.
func ReadDiagnosis(fs afero.Fs, diagnosisPath string) (*Diagnosis, error) {
	content, err := afero.ReadFile(fs, diagnosisPath)
	if err != nil {
		return nil, fmt.Errorf("reading diagnosis file: %w", err)
	}

	doc, err := md.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("parsing diagnosis: %w", err)
	}

	// Convert frontmatter map to Diagnosis struct via YAML marshaling
	// This ensures type conversions (especially time.Time) are handled correctly
	yamlData, err := yaml.Marshal(doc.Frontmatter)
	if err != nil {
		return nil, fmt.Errorf("marshaling frontmatter: %w", err)
	}

	var diag Diagnosis
	if err := yaml.Unmarshal(yamlData, &diag); err != nil {
		return nil, fmt.Errorf("unmarshaling frontmatter: %w", err)
	}

	diag.Body = doc.Body
	return &diag, nil
}

// WriteDiagnosis writes a diagnosis struct to file using pkg/md abstraction.
// Creates the .mission directory if it doesn't exist.
func WriteDiagnosis(fs afero.Fs, diagnosisPath string, diag *Diagnosis) error {
	// Marshal diagnosis struct to YAML to get frontmatter fields
	yamlData, err := yaml.Marshal(diag)
	if err != nil {
		return fmt.Errorf("marshaling diagnosis: %w", err)
	}

	// Unmarshal into map for md.Document
	var frontmatter map[string]interface{}
	if err := yaml.Unmarshal(yamlData, &frontmatter); err != nil {
		return fmt.Errorf("unmarshaling to map: %w", err)
	}

	// Use pkg/md to write document with frontmatter
	doc := &md.Document{
		Frontmatter: frontmatter,
		Body:        diag.Body,
	}

	content, err := doc.Write()
	if err != nil {
		return fmt.Errorf("writing document: %w", err)
	}

	if err := afero.WriteFile(fs, diagnosisPath, content, 0644); err != nil {
		return fmt.Errorf("writing diagnosis file: %w", err)
	}

	return nil
}

// UpdateSection updates a specific section in the diagnosis.md file.
// If the section doesn't exist, it will be created at the end of the file.
// Section names are case-insensitive and automatically converted to uppercase.
// For list sections (INVESTIGATION, HYPOTHESES, AFFECTED FILES), content is appended rather than replaced.
func UpdateSection(fs afero.Fs, diagnosisPath string, section string, content string) error {
	if section == "" {
		return fmt.Errorf("section cannot be empty")
	}
	if content == "" {
		return fmt.Errorf("content cannot be empty")
	}

	diag, err := ReadDiagnosis(fs, diagnosisPath)
	if err != nil {
		return fmt.Errorf("reading diagnosis: %w", err)
	}

	lines := strings.Split(diag.Body, "\n")
	var result []string
	// Normalize section name: replace hyphens with spaces for matching
	normalizedSection := strings.ReplaceAll(strings.ToUpper(section), "-", " ")
	sectionHeader := "## " + normalizedSection
	foundSection := false
	upperSection := normalizedSection
	isListSection := upperSection == "INVESTIGATION" || upperSection == "HYPOTHESES" || upperSection == "AFFECTED FILES"

	for i, line := range lines {
		if strings.TrimSpace(line) == sectionHeader {
			foundSection = true
			result = append(result, line)

			if isListSection {
				// For list sections, append to existing content (skip placeholders)
				j := i + 1
				for ; j < len(lines); j++ {
					if strings.HasPrefix(strings.TrimSpace(lines[j]), "## ") {
						break
					}
					line := strings.TrimSpace(lines[j])
					// Skip placeholder lines
					if line != "- [ ] Initial investigation pending" &&
						line != "1. **[UNKNOWN]** Investigation not yet started" &&
						line != "- TBD" {
						result = append(result, lines[j])
					}
				}
				result = append(result, content, "")
				if j < len(lines) {
					result = append(result, lines[j:]...)
				}
			} else {
				// For text sections, replace content
				result = append(result, content, "")
				for j := i + 1; j < len(lines); j++ {
					if strings.HasPrefix(strings.TrimSpace(lines[j]), "## ") {
						result = append(result, lines[j:]...)
						break
					}
				}
			}
			break
		}
		result = append(result, line)
	}

	if !foundSection {
		// Add new section at end
		result = append(result, "", sectionHeader, content)
	}

	diag.Body = strings.Join(result, "\n")
	return WriteDiagnosis(fs, diagnosisPath, diag)
}

// CreateDiagnosis generates a new diagnosis.md file with the given symptom.
// It creates a structured investigation template following the debug workflow design.
func CreateDiagnosis(fs afero.Fs, diagnosisPath string, symptom string) error {
	if symptom == "" {
		return fmt.Errorf("symptom cannot be empty")
	}

	diag := Diagnosis{
		ID:         fmt.Sprintf("DIAG-%s", time.Now().Format("20060102-150405")),
		Status:     "investigating",
		Confidence: "low",
		Created:    time.Now(),
		Symptom:    symptom,
	}

	frontmatter, err := yaml.Marshal(diag)
	if err != nil {
		return fmt.Errorf("marshaling frontmatter: %w", err)
	}

	content := fmt.Sprintf(`---
%s---

## SYMPTOM
%s

## INVESTIGATION
- [ ] Initial investigation pending

## HYPOTHESES
1. **[UNKNOWN]** Investigation not yet started

## ROOT CAUSE
To be determined

## AFFECTED FILES
- TBD

## RECOMMENDED FIX
To be determined after investigation

## REPRODUCTION
TBD
`, string(frontmatter), symptom)

	if err := fs.MkdirAll(".mission", 0755); err != nil {
		return fmt.Errorf("creating .mission directory: %w", err)
	}

	if err := afero.WriteFile(fs, diagnosisPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing diagnosis file: %w", err)
	}

	return nil
}

// DiagnosisExists checks if a diagnosis.md file exists
func DiagnosisExists(fs afero.Fs, diagnosisPath string) (bool, error) {
	_, err := fs.Stat(diagnosisPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// UpdateList updates a list section (investigation, hypotheses, affected files) with optional append mode.
// Section names are case-insensitive and hyphens are normalized to spaces (e.g., "affected-files" matches "AFFECTED FILES").
// In replace mode (default), existing content is removed and replaced with new items.
// In append mode, new items are added after existing items.
func UpdateList(fs afero.Fs, diagnosisPath string, section string, items []string, appendMode bool) error {
	if section == "" {
		return fmt.Errorf("section cannot be empty")
	}
	if len(items) == 0 {
		return fmt.Errorf("items cannot be empty")
	}

	diag, err := ReadDiagnosis(fs, diagnosisPath)
	if err != nil {
		return fmt.Errorf("reading diagnosis: %w", err)
	}

	lines := strings.Split(diag.Body, "\n")
	var result []string
	// Normalize section name: replace hyphens with spaces for matching
	normalizedSection := strings.ReplaceAll(strings.ToUpper(section), "-", " ")
	sectionHeader := "## " + normalizedSection
	foundSection := false
	lineIndex := 0

	for lineIndex < len(lines) {
		currentLine := lines[lineIndex]

		if strings.TrimSpace(currentLine) == sectionHeader {
			result = append(result, currentLine)

			if appendMode {
				existingItems := extractExistingItems(lines, lineIndex+1)
				result = append(result, existingItems...)
			}

			addFormattedItems(&result, normalizedSection, items)

			foundSection = true
			nextSectionIndex := skipSectionContent(lines, lineIndex+1)
			if nextSectionIndex < len(lines) {
				result = append(result, "")
			}
			lineIndex = nextSectionIndex
			continue
		}

		result = append(result, currentLine)
		lineIndex++
	}

	if !foundSection {
		result = append(result, "", sectionHeader)
		addFormattedItems(&result, normalizedSection, items)
	}

	diag.Body = strings.Join(result, "\n")
	return WriteDiagnosis(fs, diagnosisPath, diag)
}

// extractExistingItems collects existing items from a section.
func extractExistingItems(lines []string, startIndex int) []string {
	var existingItems []string
	for j := startIndex; j < len(lines); j++ {
		if strings.HasPrefix(strings.TrimSpace(lines[j]), "## ") {
			break
		}
		line := strings.TrimSpace(lines[j])
		// Skip placeholder lines
		if line != "" && line != "- [ ] Initial investigation pending" &&
			line != "1. **[UNKNOWN]** Investigation not yet started" &&
			line != "- TBD" {
			existingItems = append(existingItems, lines[j])
		}
	}
	return existingItems
}

// skipSectionContent skips content until the next section and returns the index.
func skipSectionContent(lines []string, startIndex int) int {
	for j := startIndex; j < len(lines); j++ {
		if strings.HasPrefix(strings.TrimSpace(lines[j]), "## ") {
			return j
		}
	}
	return len(lines)
}

// addFormattedItems adds items with proper formatting based on section type.
func addFormattedItems(result *[]string, section string, items []string) {
	// Normalize section name for comparison
	normalizedSection := strings.ReplaceAll(strings.ToUpper(section), "-", " ")

	if normalizedSection == "INVESTIGATION" {
		for _, item := range items {
			*result = append(*result, "- [ ] "+item)
		}
	} else if normalizedSection == "HYPOTHESES" {
		for i, item := range items {
			*result = append(*result, fmt.Sprintf("%d. %s", i+1, item))
		}
	} else {
		for _, item := range items {
			*result = append(*result, "- "+item)
		}
	}
}

// UpdateFrontmatter updates status and/or confidence fields in diagnosis frontmatter.
// Either status or confidence can be empty to skip updating that field.
// Valid status values: investigating, confirmed, inconclusive
// Valid confidence values: low, medium, high
func UpdateFrontmatter(fs afero.Fs, diagnosisPath string, status string, confidence string) error {
	if status == "" && confidence == "" {
		return fmt.Errorf("at least one of status or confidence must be provided")
	}

	diag, err := ReadDiagnosis(fs, diagnosisPath)
	if err != nil {
		return fmt.Errorf("reading diagnosis: %w", err)
	}

	if status != "" {
		validStatuses := map[string]bool{
			"investigating": true,
			"confirmed":     true,
			"inconclusive":  true,
		}
		if !validStatuses[status] {
			return fmt.Errorf("invalid status: %s (must be investigating, confirmed, or inconclusive)", status)
		}
		diag.Status = status
	}

	if confidence != "" {
		validConfidence := map[string]bool{
			"low":    true,
			"medium": true,
			"high":   true,
		}
		if !validConfidence[confidence] {
			return fmt.Errorf("invalid confidence: %s (must be low, medium, or high)", confidence)
		}
		diag.Confidence = confidence
	}

	return WriteDiagnosis(fs, diagnosisPath, diag)
}

// Finalize validates diagnosis.md completeness and returns JSON result.
// Checks for required sections (SYMPTOM, INVESTIGATION, HYPOTHESES, ROOT CAUSE,
// AFFECTED FILES, RECOMMENDED FIX) and validates frontmatter fields.
// REPRODUCTION is optional. RECOMMENDED FIX is only required if status is "confirmed".
// Returns JSON with "valid" boolean and optional "missing_sections" array.
func Finalize(fs afero.Fs, diagnosisPath string) (string, error) {
	diag, err := ReadDiagnosis(fs, diagnosisPath)
	if err != nil {
		return "", fmt.Errorf("reading diagnosis file: %w", err)
	}

	var missingSections []string
	requiredSections := []string{"SYMPTOM", "INVESTIGATION", "HYPOTHESES", "ROOT CAUSE", "AFFECTED FILES"}

	// RECOMMENDED FIX is only required if status is confirmed
	if diag.Status == "confirmed" {
		requiredSections = append(requiredSections, "RECOMMENDED FIX")
	}

	for _, section := range requiredSections {
		sectionHeader := "## " + section
		if !strings.Contains(diag.Body, sectionHeader) {
			missingSections = append(missingSections, section)
		}
	}

	result := map[string]interface{}{
		"valid":   len(missingSections) == 0,
		"message": "Diagnosis validated successfully",
	}

	if len(missingSections) > 0 {
		result["valid"] = false
		result["missing_sections"] = missingSections
		result["message"] = fmt.Sprintf("Missing required sections: %s", strings.Join(missingSections, ", "))
	}

	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("formatting output: %w", err)
	}

	return string(jsonOutput), nil
}
