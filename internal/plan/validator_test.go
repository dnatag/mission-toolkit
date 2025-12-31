package plan

import (
	"encoding/json"
	"testing"

	"github.com/spf13/afero"
)

func TestValidatorService_ValidatePlan_ValidPlan(t *testing.T) {
	fs := afero.NewMemMapFs()
	validator := NewValidatorService(fs, "test-mission", "/project")

	// Create a valid plan
	validPlan := PlanSpec{
		Intent:       "Add user authentication",
		Scope:        []string{"auth.go", "auth_test.go"},
		Domain:       []string{"security"},
		Plan:         []string{"Create auth service", "Add tests"},
		Verification: "go test ./auth",
	}

	planJSON, _ := json.Marshal(validPlan)
	afero.WriteFile(fs, "/project/plan.json", planJSON, 0644)

	// Create existing file
	afero.WriteFile(fs, "/project/auth.go", []byte("package auth"), 0644)

	result, err := validator.ValidatePlan("/project/plan.json")
	if err != nil {
		t.Fatalf("ValidatePlan() error = %v", err)
	}

	if !result.Valid {
		t.Errorf("Expected valid plan, got invalid. Errors: %v", result.Errors)
	}

	if len(result.Errors) > 0 {
		t.Errorf("Expected no errors, got: %v", result.Errors)
	}
}

func TestValidatorService_ValidatePlan_MissingFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	validator := NewValidatorService(fs, "test-mission", "/project")

	result, err := validator.ValidatePlan("/project/nonexistent.json")
	if err != nil {
		t.Fatalf("ValidatePlan() error = %v", err)
	}

	if result.Valid {
		t.Error("Expected invalid plan for missing file")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected errors for missing file")
	}
}

func TestValidatorService_ValidatePlan_InvalidJSON(t *testing.T) {
	fs := afero.NewMemMapFs()
	validator := NewValidatorService(fs, "test-mission", "/project")

	// Create invalid JSON
	afero.WriteFile(fs, "/project/invalid.json", []byte("invalid json"), 0644)

	result, err := validator.ValidatePlan("/project/invalid.json")
	if err != nil {
		t.Fatalf("ValidatePlan() error = %v", err)
	}

	if result.Valid {
		t.Error("Expected invalid plan for malformed JSON")
	}

	if len(result.FormatIssues) == 0 {
		t.Error("Expected format issues for invalid JSON")
	}
}

func TestValidatorService_ValidatePlan_MissingRequiredFields(t *testing.T) {
	fs := afero.NewMemMapFs()
	validator := NewValidatorService(fs, "test-mission", "/project")

	testCases := []struct {
		name           string
		plan           PlanSpec
		expectedErrors int
	}{
		{
			name: "missing intent",
			plan: PlanSpec{
				Scope: []string{"file.go"},
				Plan:  []string{"step 1"},
			},
			expectedErrors: 1,
		},
		{
			name: "missing scope",
			plan: PlanSpec{
				Intent: "test intent",
				Plan:   []string{"step 1"},
			},
			expectedErrors: 1,
		},
		{
			name: "missing plan",
			plan: PlanSpec{
				Intent: "test intent",
				Scope:  []string{"file.go"},
			},
			expectedErrors: 1,
		},
		{
			name:           "all missing",
			plan:           PlanSpec{},
			expectedErrors: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			planJSON, _ := json.Marshal(tc.plan)
			planPath := "/project/" + tc.name + ".json"
			afero.WriteFile(fs, planPath, planJSON, 0644)

			result, err := validator.ValidatePlan(planPath)
			if err != nil {
				t.Fatalf("ValidatePlan() error = %v", err)
			}

			if result.Valid {
				t.Error("Expected invalid plan for missing required fields")
			}

			if len(result.Errors) != tc.expectedErrors {
				t.Errorf("Expected %d errors, got %d: %v", tc.expectedErrors, len(result.Errors), result.Errors)
			}
		})
	}
}

func TestValidatorService_ValidatePlan_SecurityIssues(t *testing.T) {
	fs := afero.NewMemMapFs()
	validator := NewValidatorService(fs, "test-mission", "/project")

	testCases := []struct {
		name           string
		plan           PlanSpec
		expectedIssues int
	}{
		{
			name: "path traversal in scope",
			plan: PlanSpec{
				Intent: "test",
				Scope:  []string{"../../../etc/passwd"},
				Plan:   []string{"step 1"},
			},
			expectedIssues: 1,
		},
		{
			name: "suspicious intent",
			plan: PlanSpec{
				Intent: "rm -rf / everything",
				Scope:  []string{"file.go"},
				Plan:   []string{"step 1"},
			},
			expectedIssues: 1,
		},
		{
			name: "suspicious plan step",
			plan: PlanSpec{
				Intent: "test",
				Scope:  []string{"file.go"},
				Plan:   []string{"sudo rm -rf /"},
			},
			expectedIssues: 2, // "sudo" and "rm -rf"
		},
		{
			name: "dangerous verification",
			plan: PlanSpec{
				Intent:       "test",
				Scope:        []string{"file.go"},
				Plan:         []string{"step 1"},
				Verification: "rm -rf test_data",
			},
			expectedIssues: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			planJSON, _ := json.Marshal(tc.plan)
			planPath := "/project/" + tc.name + ".json"
			afero.WriteFile(fs, planPath, planJSON, 0644)

			result, err := validator.ValidatePlan(planPath)
			if err != nil {
				t.Fatalf("ValidatePlan() error = %v", err)
			}

			if len(result.SecurityIssues) != tc.expectedIssues {
				t.Errorf("Expected %d security issues, got %d: %v",
					tc.expectedIssues, len(result.SecurityIssues), result.SecurityIssues)
			}
		})
	}
}

func TestValidatorService_ValidatePlan_FileAccess(t *testing.T) {
	fs := afero.NewMemMapFs()
	validator := NewValidatorService(fs, "test-mission", "/project")

	// Create plan with mix of existing and new files
	plan := PlanSpec{
		Intent: "test file access",
		Scope:  []string{"existing.go", "new.go", "subdir/new.go"},
		Plan:   []string{"step 1"},
	}

	planJSON, _ := json.Marshal(plan)
	afero.WriteFile(fs, "/project/plan.json", planJSON, 0644)

	// Create existing file
	afero.WriteFile(fs, "/project/existing.go", []byte("package main"), 0644)

	result, err := validator.ValidatePlan("/project/plan.json")
	if err != nil {
		t.Fatalf("ValidatePlan() error = %v", err)
	}

	// Should have file validation messages for all files
	if len(result.FileValidation) != 3 {
		t.Errorf("Expected 3 file validation messages, got %d: %v",
			len(result.FileValidation), result.FileValidation)
	}
}

func TestValidatorService_ValidatePlan_UnknownDomain(t *testing.T) {
	fs := afero.NewMemMapFs()
	validator := NewValidatorService(fs, "test-mission", "/project")

	plan := PlanSpec{
		Intent: "test",
		Scope:  []string{"file.go"},
		Domain: []string{"security", "unknown-domain", "performance"},
		Plan:   []string{"step 1"},
	}

	planJSON, _ := json.Marshal(plan)
	afero.WriteFile(fs, "/project/plan.json", planJSON, 0644)

	result, err := validator.ValidatePlan("/project/plan.json")
	if err != nil {
		t.Fatalf("ValidatePlan() error = %v", err)
	}

	// Should have warning for unknown domain
	if len(result.Warnings) == 0 {
		t.Error("Expected warning for unknown domain")
	}

	foundUnknownDomainWarning := false
	for _, warning := range result.Warnings {
		if contains(warning, "unknown-domain") {
			foundUnknownDomainWarning = true
			break
		}
	}

	if !foundUnknownDomainWarning {
		t.Errorf("Expected warning about unknown-domain, got warnings: %v", result.Warnings)
	}
}

func TestValidationResult_ToJSON(t *testing.T) {
	result := &ValidationResult{
		Valid:          false,
		Errors:         []string{"error 1", "error 2"},
		Warnings:       []string{"warning 1"},
		SecurityIssues: []string{"security issue 1"},
	}

	jsonStr, err := result.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() error = %v", err)
	}

	// Verify it's valid JSON by unmarshaling
	var parsed ValidationResult
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		t.Fatalf("Generated JSON is invalid: %v", err)
	}

	if parsed.Valid != result.Valid {
		t.Errorf("JSON parsing changed Valid field: expected %t, got %t", result.Valid, parsed.Valid)
	}

	if len(parsed.Errors) != len(result.Errors) {
		t.Errorf("JSON parsing changed Errors count: expected %d, got %d", len(result.Errors), len(parsed.Errors))
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsAt(s, substr))))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
