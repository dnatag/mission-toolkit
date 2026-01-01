package plan

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestGeneratorService_GenerateMission(t *testing.T) {
	tests := []struct {
		name        string
		planSpec    PlanSpec
		missionID   string
		expectError bool
		checkFields []string
	}{
		{
			name: "basic WET mission generation",
			planSpec: PlanSpec{
				Intent:       "Add user authentication",
				Scope:        []string{"auth.go", "auth_test.go"},
				Plan:         []string{"Create auth service", "Add tests"},
				Verification: "go test ./auth",
			},
			missionID:   "test-123",
			expectError: false,
			checkFields: []string{"id: test-123", "type: WET", "status: planned"},
		},
		{
			name: "mission with Files field",
			planSpec: PlanSpec{
				Intent: "Refactor authentication system",
				Files: []FileSpec{
					{Path: "auth/service.go", Action: FileActionModify},
					{Path: "auth/handler.go", Action: FileActionModify},
					{Path: "auth/middleware.go", Action: FileActionCreate},
				},
				Plan:         []string{"Extract common interface", "Implement factory pattern", "Update tests"},
				Verification: "go test ./auth/... && go build",
			},
			missionID:   "refactor-456",
			expectError: false,
			checkFields: []string{"id: refactor-456", "type: WET", "Refactor authentication system"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			generator := NewGeneratorService(fs, tt.missionID)

			// Create plan.json
			planData, _ := json.Marshal(tt.planSpec)
			afero.WriteFile(fs, "plan.json", planData, 0644)

			// Generate mission
			result, err := generator.GenerateMission("plan.json", "mission.md")

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !result.Success {
				t.Errorf("Expected success=true, got %v", result.Success)
			}

			// Check generated content
			content, err := afero.ReadFile(fs, "mission.md")
			if err != nil {
				t.Fatalf("Failed to read generated mission: %v", err)
			}

			contentStr := string(content)
			for _, field := range tt.checkFields {
				if !strings.Contains(contentStr, field) {
					t.Errorf("Generated mission missing field: %s", field)
				}
			}

			// Verify structure
			if !strings.Contains(contentStr, "# MISSION") {
				t.Error("Missing mission header")
			}
			if !strings.Contains(contentStr, "## INTENT") {
				t.Error("Missing intent section")
			}
			if !strings.Contains(contentStr, "## SCOPE") {
				t.Error("Missing scope section")
			}
			if !strings.Contains(contentStr, "## PLAN") {
				t.Error("Missing plan section")
			}
			if !strings.Contains(contentStr, "## VERIFICATION") {
				t.Error("Missing verification section")
			}
		})
	}
}

func TestGeneratorService_GenerateMission_FileErrors(t *testing.T) {
	fs := afero.NewMemMapFs()
	generator := NewGeneratorService(fs, "test-id")

	// Test missing plan file
	result, err := generator.GenerateMission("nonexistent.json", "mission.md")
	if err == nil {
		t.Error("Expected error for missing plan file")
	}
	if result.Success {
		t.Error("Expected success=false for missing plan file")
	}

	// Test invalid JSON
	afero.WriteFile(fs, "invalid.json", []byte("invalid json"), 0644)
	result, err = generator.GenerateMission("invalid.json", "mission.md")
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
	if result.Success {
		t.Error("Expected success=false for invalid JSON")
	}
}

func TestGenerateResult_ToJSON(t *testing.T) {
	result := &GenerateResult{
		Success:     true,
		Message:     "Test message",
		OutputFile:  "test.md",
		PlanFile:    "plan.json",
		MissionType: "WET",
		Track:       2,
	}

	jsonStr, err := result.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert to JSON: %v", err)
	}

	// Parse back to verify
	var parsed GenerateResult
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if parsed.Success != result.Success {
		t.Errorf("JSON conversion failed for Success field")
	}
	if parsed.Message != result.Message {
		t.Errorf("JSON conversion failed for Message field")
	}
	if parsed.Track != result.Track {
		t.Errorf("JSON conversion failed for Track field")
	}
}
