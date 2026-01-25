// Package validation provides input validation for mission planning.
package validation

import (
	"path/filepath"
	"strings"

	"github.com/dnatag/mission-toolkit/pkg/diagnosis"
	"github.com/spf13/afero"
)

// DiagnosisInfo contains diagnosis metadata for debug workflow integration
type DiagnosisInfo struct {
	ID            string   `json:"id"`
	RootCause     string   `json:"root_cause"`
	AffectedFiles []string `json:"affected_files"`
}

// CheckResult represents the result of input validation
type CheckResult struct {
	IsValid   bool           `json:"is_valid"`
	Message   string         `json:"message"`
	NextStep  string         `json:"next_step"`
	Diagnosis *DiagnosisInfo `json:"diagnosis,omitempty"`
}

// Validate checks if input is valid for mission planning.
// If missionDir is provided and contains a valid diagnosis.md, returns DIAGNOSIS_DETECTED.
func Validate(input string, fs afero.Fs, missionDir string) *CheckResult {
	// Check for diagnosis.md first (takes priority)
	if missionDir != "" && fs != nil {
		diagnosisPath := filepath.Join(missionDir, "diagnosis.md")
		if diag, err := diagnosis.ReadDiagnosis(fs, diagnosisPath); err == nil {
			rootCause, affectedFiles := extractDiagnosisInfo(diag.Body)
			return &CheckResult{
				IsValid:  true,
				Message:  "Diagnosis file detected",
				NextStep: "DIAGNOSIS_DETECTED",
				Diagnosis: &DiagnosisInfo{
					ID:            diag.ID,
					RootCause:     rootCause,
					AffectedFiles: affectedFiles,
				},
			}
		}
	}

	trimmed := strings.TrimSpace(input)

	if trimmed == "" {
		return &CheckResult{
			IsValid:  false,
			Message:  "Input is empty or whitespace",
			NextStep: "ASK_USER: What is your intent or goal for this task?",
		}
	}

	if trimmed == "$ARGUMENTS" {
		return &CheckResult{
			IsValid:  false,
			Message:  "Input is a placeholder - no intent provided",
			NextStep: "ASK_USER: What is your intent or goal for this task?",
		}
	}

	return &CheckResult{
		IsValid:  true,
		Message:  "Input is valid",
		NextStep: "PROCEED with execution",
	}
}

// extractDiagnosisInfo parses ROOT CAUSE and AFFECTED FILES from diagnosis body
func extractDiagnosisInfo(body string) (rootCause string, affectedFiles []string) {
	lines := strings.Split(body, "\n")
	var currentSection string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ROOT CAUSE") {
			currentSection = "root_cause"
			continue
		}
		if strings.HasPrefix(trimmed, "## AFFECTED FILES") {
			currentSection = "affected_files"
			continue
		}
		if strings.HasPrefix(trimmed, "## ") {
			currentSection = ""
			continue
		}

		if currentSection == "root_cause" && trimmed != "" && rootCause == "" {
			rootCause = trimmed
		}
		if currentSection == "affected_files" && strings.HasPrefix(trimmed, "- ") {
			affectedFiles = append(affectedFiles, strings.TrimPrefix(trimmed, "- "))
		}
	}
	return
}
