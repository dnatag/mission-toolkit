package analyze

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/dnatag/mission-toolkit/internal/logger"
	"github.com/spf13/afero"
)

func TestTestService_ProvideTemplate(t *testing.T) {
	fs := afero.NewMemMapFs()

	// Use console-only logger config to prevent execution.log creation
	loggerConfig := &logger.Config{
		Level:    logger.DefaultConfig().Level,
		Format:   logger.DefaultConfig().Format,
		Output:   logger.OutputConsole,
		FilePath: "", // Empty path to prevent file creation
		Fs:       fs, // Use same in-memory filesystem
	}

	missionContent := `---
id: test-123
status: planned
---

## INTENT
Add user authentication

## SCOPE
auth.go
handler.go`

	if err := afero.WriteFile(fs, ".mission/mission.md", []byte(missionContent), 0644); err != nil {
		t.Fatal(err)
	}

	service := NewTestServiceWithConfig(fs, loggerConfig)
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
	if !strings.Contains(contentStr, "## Current Mission Context") {
		t.Error("Template missing Current Mission Context section")
	}
	if !strings.Contains(contentStr, "Add user authentication") {
		t.Error("Template missing intent text")
	}
	if !strings.Contains(contentStr, "auth.go") {
		t.Error("Template missing scope files")
	}
}

func TestTestService_ProvideTemplate_MissingMission(t *testing.T) {
	fs := afero.NewMemMapFs()

	// Use console-only logger config to prevent execution.log creation
	loggerConfig := &logger.Config{
		Level:    logger.DefaultConfig().Level,
		Format:   logger.DefaultConfig().Format,
		Output:   logger.OutputConsole,
		FilePath: "", // Empty path to prevent file creation
		Fs:       fs, // Use same in-memory filesystem
	}

	service := NewTestServiceWithConfig(fs, loggerConfig)
	_, err := service.ProvideTemplate()

	if err == nil {
		t.Error("Expected error for missing mission.md")
	}
}
