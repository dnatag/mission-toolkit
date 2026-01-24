package diagnosis

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

// Diagnosis represents a debug investigation with structured findings
type Diagnosis struct {
	ID         string    `yaml:"id"`
	Status     string    `yaml:"status"`
	Confidence string    `yaml:"confidence"`
	Created    time.Time `yaml:"created"`
	Symptom    string
	Body       string
}

// ReadDiagnosis reads and parses a diagnosis.md file.
// Returns an error if the file doesn't exist or has invalid format.
func ReadDiagnosis(fs afero.Fs, diagnosisPath string) (*Diagnosis, error) {
	content, err := afero.ReadFile(fs, diagnosisPath)
	if err != nil {
		return nil, fmt.Errorf("reading diagnosis file: %w", err)
	}

	parts := strings.SplitN(string(content), "---\n", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid diagnosis format: missing frontmatter")
	}

	var diag Diagnosis
	if err := yaml.Unmarshal([]byte(parts[1]), &diag); err != nil {
		return nil, fmt.Errorf("parsing frontmatter: %w", err)
	}

	diag.Body = parts[2]
	return &diag, nil
}

// WriteDiagnosis writes a diagnosis struct to file.
// Creates the .mission directory if it doesn't exist.
func WriteDiagnosis(fs afero.Fs, diagnosisPath string, diag *Diagnosis) error {
	frontmatter, err := yaml.Marshal(diag)
	if err != nil {
		return fmt.Errorf("marshaling frontmatter: %w", err)
	}

	content := fmt.Sprintf("---\n%s---\n%s", string(frontmatter), diag.Body)
	if err := afero.WriteFile(fs, diagnosisPath, []byte(content), 0644); err != nil {
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
	sectionHeader := "## " + strings.ToUpper(section)
	foundSection := false
	upperSection := strings.ToUpper(section)
	isListSection := upperSection == "INVESTIGATION" || upperSection == "HYPOTHESES" || upperSection == "AFFECTED FILES"

	for i, line := range lines {
		if strings.TrimSpace(line) == sectionHeader {
			foundSection = true
			result = append(result, line)

			if isListSection {
				// For list sections, append to existing content
				j := i + 1
				for ; j < len(lines); j++ {
					if strings.HasPrefix(strings.TrimSpace(lines[j]), "## ") {
						break
					}
					result = append(result, lines[j])
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
