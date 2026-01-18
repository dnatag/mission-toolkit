package analyze

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/dnatag/mission-toolkit/pkg/logger"
	"github.com/spf13/afero"
)

func TestDecomposeService_ProvideTemplate(t *testing.T) {
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
Build complete payment processing system

## SCOPE
payments/handler.go
payments/service.go
payments/model.go`

	if err := afero.WriteFile(fs, ".mission/mission.md", []byte(missionContent), 0644); err != nil {
		t.Fatal(err)
	}

	service := NewDecomposeServiceWithConfig(fs, loggerConfig)
	output, err := service.ProvideTemplate()

	if err != nil {
		t.Fatalf("ProvideTemplate failed: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	templatePath := result["template_path"]
	content, err := afero.ReadFile(fs, templatePath)
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "# DECOMPOSE ANALYSIS TEMPLATE") {
		t.Error("Template missing header")
	}
	if !strings.Contains(contentStr, "Build complete payment processing system") {
		t.Error("Template missing intent text")
	}
	if !strings.Contains(contentStr, "payments/handler.go") {
		t.Error("Template missing scope text")
	}
}

func TestDecomposeService_ProvideTemplate_MissingMission(t *testing.T) {
	fs := afero.NewMemMapFs()

	// Use console-only logger config to prevent execution.log creation
	loggerConfig := &logger.Config{
		Level:    logger.DefaultConfig().Level,
		Format:   logger.DefaultConfig().Format,
		Output:   logger.OutputConsole,
		FilePath: "", // Empty path to prevent file creation
		Fs:       fs, // Use same in-memory filesystem
	}

	service := NewDecomposeServiceWithConfig(fs, loggerConfig)
	_, err := service.ProvideTemplate()

	if err == nil {
		t.Error("Expected error for missing mission.md")
	}
}
