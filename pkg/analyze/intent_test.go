package analyze

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestIntentService_ProvideTemplate(t *testing.T) {
	fs := afero.NewMemMapFs()

	// Use helper to create test logger config
	loggerConfig := CreateTestLoggerConfig(fs)

	service := NewIntentServiceWithConfig(fs, loggerConfig)
	output, err := service.ProvideTemplate("add auth")

	if err != nil {
		t.Fatalf("ProvideTemplate failed: %v", err)
	}

	// Verify output is valid JSON
	var result map[string]string
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	// Verify JSON structure
	templatePath, ok := result["template_path"]
	if !ok {
		t.Error("Output missing template_path field")
	}

	// Verify file was created and contains expected content
	content, err := afero.ReadFile(fs, templatePath)
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "## User Input") {
		t.Error("Template missing User Input section")
	}
	if !strings.Contains(contentStr, "add auth") {
		t.Error("Template missing user input text")
	}
}

func TestIntentService_ProvideTemplate_EmptyInput(t *testing.T) {
	fs := afero.NewMemMapFs()

	// Use helper to create test logger config
	loggerConfig := CreateTestLoggerConfig(fs)

	service := NewIntentServiceWithConfig(fs, loggerConfig)
	output, err := service.ProvideTemplate("")

	if err != nil {
		t.Fatalf("ProvideTemplate failed: %v", err)
	}

	// Verify output is valid JSON
	var result map[string]string
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	// Verify file exists
	templatePath := result["template_path"]
	if exists, _ := afero.Exists(fs, templatePath); !exists {
		t.Error("Template file was not created")
	}
}
