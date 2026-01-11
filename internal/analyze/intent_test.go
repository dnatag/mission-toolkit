package analyze

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestIntentService_ProvideTemplate(t *testing.T) {
	// Setup: Create temp directory for test
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	service := NewIntentService()
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
	content, err := os.ReadFile(templatePath)
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
	// Setup: Create temp directory for test
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	service := NewIntentService()
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
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		t.Error("Template file was not created")
	}
}
