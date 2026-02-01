// Package diagnosis provides functions for managing diagnosis.md files with YAML frontmatter.
// It uses the md.Document abstraction for consistent markdown handling.
package diagnosis

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dnatag/mission-toolkit/pkg/md"
	"github.com/spf13/afero"
)

// WriteDiagnosis writes a Diagnosis struct to a diagnosis.md file using pkg/md abstraction.
func WriteDiagnosis(fs afero.Fs, diagnosisPath string, diag *Diagnosis) error {
	doc := diagnosisToDocument(diag)

	content, err := doc.Write()
	if err != nil {
		return fmt.Errorf("writing document: %w", err)
	}

	if err := afero.WriteFile(fs, diagnosisPath, content, 0644); err != nil {
		return fmt.Errorf("writing diagnosis file: %w", err)
	}

	return nil
}

// CreateDiagnosis creates a new diagnosis.md file with initial structure.
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

	doc := &md.Document{
		Frontmatter: map[string]interface{}{
			"id":         diag.ID,
			"status":     diag.Status,
			"confidence": diag.Confidence,
			"created":    diag.Created,
		},
		Body: fmt.Sprintf(`## SYMPTOM
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
`, symptom),
	}

	content, err := doc.Write()
	if err != nil {
		return fmt.Errorf("writing document: %w", err)
	}

	if err := fs.MkdirAll(".mission", 0755); err != nil {
		return fmt.Errorf("creating .mission directory: %w", err)
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

	return modifyDiagnosis(fs, diagnosisPath, func(doc *md.Document) error {
		normalizedSection := normalizeSection(section)

		if isListSection(normalizedSection) {
			if err := doc.AppendSectionList(normalizedSection, []string{strings.TrimPrefix(content, "- ")}); err != nil {
				return fmt.Errorf("appending to section: %w", err)
			}
		} else {
			if err := doc.UpdateSectionContent(normalizedSection, content); err != nil {
				return fmt.Errorf("updating section: %w", err)
			}
		}
		return nil
	})
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

	return modifyDiagnosis(fs, diagnosisPath, func(doc *md.Document) error {
		normalizedSection := normalizeSection(section)
		formattedItems := formatItems(normalizedSection, items)

		if appendMode {
			if err := doc.AppendSectionList(normalizedSection, formattedItems); err != nil {
				return fmt.Errorf("appending to section: %w", err)
			}
		} else {
			if err := doc.UpdateSectionList(normalizedSection, formattedItems); err != nil {
				return fmt.Errorf("updating section: %w", err)
			}
		}
		return nil
	})
}

// diagnosisToDocument converts a Diagnosis to an md.Document.
func diagnosisToDocument(diag *Diagnosis) *md.Document {
	return &md.Document{
		Frontmatter: map[string]interface{}{
			"id":         diag.ID,
			"status":     diag.Status,
			"confidence": diag.Confidence,
			"created":    diag.Created,
		},
		Body: diag.Body,
	}
}

// modifyDiagnosis reads a diagnosis, applies a modification function to its document, and writes it back.
func modifyDiagnosis(fs afero.Fs, diagnosisPath string, modifyFn func(*md.Document) error) error {
	diag, err := ReadDiagnosis(fs, diagnosisPath)
	if err != nil {
		return fmt.Errorf("reading diagnosis: %w", err)
	}

	doc := diagnosisToDocument(diag)
	if err := modifyFn(doc); err != nil {
		return err
	}

	diag.Body = doc.Body
	return WriteDiagnosis(fs, diagnosisPath, diag)
}

// normalizeSection normalizes section names to uppercase with spaces instead of hyphens.
func normalizeSection(section string) string {
	return strings.ReplaceAll(strings.ToUpper(section), "-", " ")
}

// isListSection checks if a section is a list section.
func isListSection(normalizedSection string) bool {
	return normalizedSection == "INVESTIGATION" ||
		normalizedSection == "HYPOTHESES" ||
		normalizedSection == "AFFECTED FILES"
}

// formatItems formats items based on section type.
func formatItems(section string, items []string) []string {
	formatted := make([]string, len(items))

	switch normalizeSection(section) {
	case "INVESTIGATION":
		for i, item := range items {
			formatted[i] = "[ ] " + item
		}
	case "HYPOTHESES":
		for i, item := range items {
			formatted[i] = fmt.Sprintf("%d. %s", i+1, item)
		}
	default:
		copy(formatted, items)
	}

	return formatted
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
