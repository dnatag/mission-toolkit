package analyze

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestScopeService_ProvideTemplate(t *testing.T) {
	// Setup: Create temp directory for test
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	fs := afero.NewMemMapFs()

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

	service := NewScopeServiceWithFS(fs)
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
	content, err := os.ReadFile(templatePath)
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

func TestScopeService_ProvideTemplate_MissingMission(t *testing.T) {
	fs := afero.NewMemMapFs()
	service := NewScopeServiceWithFS(fs)
	_, err := service.ProvideTemplate()

	if err == nil {
		t.Error("Expected error for missing mission.md")
	}
}
