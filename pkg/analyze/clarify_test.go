package analyze

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestClarifyService_ProvideTemplate(t *testing.T) {
	fs := afero.NewMemMapFs()

	// Use helper to create test logger config
	loggerConfig := CreateTestLoggerConfig(fs)

	missionContent := `---
id: test-123
status: planned
---

## INTENT
Add user authentication

## SCOPE
auth.go`

	if err := afero.WriteFile(fs, ".mission/mission.md", []byte(missionContent), 0644); err != nil {
		t.Fatal(err)
	}

	service := NewClarifyServiceWithConfig(fs, loggerConfig)
	output, err := service.ProvideTemplate()

	if err != nil {
		t.Fatalf("ProvideTemplate failed: %v", err)
	}

	// Verify output is valid JSON
	var result map[string]string
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	// Verify file was created and contains expected content
	templatePath := result["template_path"]
	content, err := afero.ReadFile(fs, templatePath)
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "## Current Intent") {
		t.Error("Template missing Current Intent section")
	}
	if !strings.Contains(contentStr, "Add user authentication") {
		t.Error("Template missing intent text")
	}
}

func TestClarifyService_ProvideTemplate_MissingMission(t *testing.T) {
	fs := afero.NewMemMapFs()

	// Use helper to create test logger config
	loggerConfig := CreateTestLoggerConfig(fs)

	service := NewClarifyServiceWithConfig(fs, loggerConfig)
	_, err := service.ProvideTemplate()

	if err == nil {
		t.Error("Expected error for missing mission.md")
	}
}
