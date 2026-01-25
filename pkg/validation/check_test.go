package validation

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantValid bool
	}{
		{"empty string", "", false},
		{"whitespace only", "   ", false},
		{"placeholder", "$ARGUMENTS", false},
		{"placeholder with spaces", "  $ARGUMENTS  ", false},
		{"valid input", "Add feature", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate(tt.input, nil, "")
			if result.IsValid != tt.wantValid {
				t.Errorf("Validate(%q).IsValid = %v, want %v", tt.input, result.IsValid, tt.wantValid)
			}
		})
	}
}

func TestValidate_DiagnosisDetected(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	fs.MkdirAll(missionDir, 0755)

	diagnosisContent := `---
id: DIAG-20260124-120000
status: confirmed
confidence: high
---

## SYMPTOM
Test symptom

## ROOT CAUSE
Session middleware skipped for /api/login endpoint

## AFFECTED FILES
- internal/auth/session.go
- internal/middleware/auth.go
`
	afero.WriteFile(fs, filepath.Join(missionDir, "diagnosis.md"), []byte(diagnosisContent), 0644)

	// Test with empty input - should return DIAGNOSIS_DETECTED
	result := Validate("", fs, missionDir)
	if result.NextStep != "DIAGNOSIS_DETECTED" {
		t.Errorf("Expected DIAGNOSIS_DETECTED, got %s", result.NextStep)
	}
	if result.Diagnosis == nil {
		t.Fatal("Expected Diagnosis to be populated")
	}
	if result.Diagnosis.ID != "DIAG-20260124-120000" {
		t.Errorf("Expected diagnosis ID DIAG-20260124-120000, got %s", result.Diagnosis.ID)
	}

	// Test with non-empty input - should still return DIAGNOSIS_DETECTED (diagnosis takes priority)
	result = Validate("fix the bug", fs, missionDir)
	if result.NextStep != "DIAGNOSIS_DETECTED" {
		t.Errorf("Expected DIAGNOSIS_DETECTED even with input, got %s", result.NextStep)
	}
}

func TestValidate_NoDiagnosis(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := ".mission"
	fs.MkdirAll(missionDir, 0755)

	// No diagnosis.md - should return normal behavior
	result := Validate("Add feature", fs, missionDir)
	if result.NextStep != "PROCEED with execution" {
		t.Errorf("Expected PROCEED, got %s", result.NextStep)
	}
	if result.Diagnosis != nil {
		t.Error("Expected Diagnosis to be nil when no diagnosis.md exists")
	}
}
