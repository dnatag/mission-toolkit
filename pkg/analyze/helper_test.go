package analyze

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dnatag/mission-toolkit/pkg/logger"
	"github.com/spf13/afero"
)

func TestCreateTestLoggerConfig(t *testing.T) {
	fs := afero.NewMemMapFs()
	config := CreateTestLoggerConfig(fs)

	if config.Output != logger.OutputConsole {
		t.Errorf("Expected OutputConsole, got %s", config.Output)
	}

	if config.FilePath != "/dev/null" {
		t.Errorf("Expected /dev/null, got %s", config.FilePath)
	}

	if config.Fs != fs {
		t.Error("Expected same filesystem instance")
	}
}

func TestCreateLogger(t *testing.T) {
	fs := afero.NewMemMapFs()

	// Test with console config (don't use nil to avoid real filesystem)
	consoleConfig := CreateTestLoggerConfig(fs)
	log1 := CreateLogger(fs, consoleConfig)
	if log1 == nil {
		t.Error("Expected logger, got nil")
	}

	// Test with custom config
	config := CreateTestLoggerConfig(fs)
	log2 := CreateLogger(fs, config)
	if log2 == nil {
		t.Error("Expected logger, got nil")
	}
}

func TestFormatOutputWithFS(t *testing.T) {
	fs := afero.NewMemMapFs()
	templateContent := "# Test Template\nThis is test content"

	output, err := FormatOutputWithFS(fs, templateContent)
	if err != nil {
		t.Fatalf("FormatOutputWithFS failed: %v", err)
	}

	// Verify JSON structure
	var result map[string]string
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	// Verify required fields
	if _, ok := result["instruction"]; !ok {
		t.Error("Output missing instruction field")
	}

	templatePath, ok := result["template_path"]
	if !ok {
		t.Error("Output missing template_path field")
	}

	// Verify file was created
	exists, err := afero.Exists(fs, templatePath)
	if err != nil {
		t.Fatalf("Error checking file existence: %v", err)
	}
	if !exists {
		t.Error("Template file was not created")
	}

	// Verify file content
	content, err := afero.ReadFile(fs, templatePath)
	if err != nil {
		t.Fatalf("Error reading template file: %v", err)
	}
	if string(content) != templateContent {
		t.Errorf("Expected %q, got %q", templateContent, string(content))
	}
}

func TestFormatOutput(t *testing.T) {
	// Setup: Create temp directory for test
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	templateContent := "# Test Template\nThis is test content"

	output, err := FormatOutput(templateContent)
	if err != nil {
		t.Fatalf("FormatOutput failed: %v", err)
	}

	// Verify JSON structure
	var result map[string]string
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	// Verify JSON fields
	templatePath, ok := result["template_path"]
	if !ok {
		t.Error("Output missing template_path field")
	}

	instruction, ok := result["instruction"]
	if !ok {
		t.Error("Output missing instruction field")
	}

	// Verify instruction content
	if !strings.Contains(instruction, "Use file read tool") {
		t.Error("Instruction missing expected text")
	}

	// Verify file was created
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		t.Errorf("Template file was not created at %s", templatePath)
	}

	// Verify file content
	content, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}

	if string(content) != templateContent {
		t.Errorf("File content mismatch.\nExpected: %s\nGot: %s", templateContent, string(content))
	}

	// Verify file path structure
	if !strings.Contains(templatePath, ".mission/templates/") {
		t.Errorf("Template path should contain .mission/templates/, got: %s", templatePath)
	}

	if !strings.HasSuffix(templatePath, ".md") {
		t.Errorf("Template path should end with .md, got: %s", templatePath)
	}
}

func TestFormatOutput_CreatesDirectory(t *testing.T) {
	// Setup: Create temp directory for test
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	// Ensure .mission/templates doesn't exist
	templateDir := filepath.Join(".mission", "templates")
	if _, err := os.Stat(templateDir); !os.IsNotExist(err) {
		t.Fatal("Template directory should not exist before test")
	}

	_, err := FormatOutput("test content")
	if err != nil {
		t.Fatalf("FormatOutput failed: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		t.Error("Template directory was not created")
	}
}

func TestFormatOutput_EmptyContent(t *testing.T) {
	// Setup: Create temp directory for test
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	output, err := FormatOutput("")
	if err != nil {
		t.Fatalf("FormatOutput failed with empty content: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	// Verify file was created even with empty content
	templatePath := result["template_path"]
	content, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}

	if len(content) != 0 {
		t.Errorf("Expected empty file, got: %s", string(content))
	}
}
