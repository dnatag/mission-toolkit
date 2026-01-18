package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Level != logrus.InfoLevel {
		t.Errorf("Expected InfoLevel, got %v", config.Level)
	}

	if config.Format != "text" {
		t.Errorf("Expected text format, got %s", config.Format)
	}

	if config.Output != "both" {
		t.Errorf("Expected both output, got %s", config.Output)
	}
}

func TestNewLogger(t *testing.T) {
	config := DefaultConfig()
	config.Fs = afero.NewMemMapFs() // Use memFS for testing
	logger := NewLogger(config)

	if logger == nil {
		t.Error("Logger should not be nil")
	}

	if logger.GetLevel() != logrus.InfoLevel {
		t.Errorf("Expected InfoLevel, got %v", logger.GetLevel())
	}
}

func TestNewLoggerWithJSONFormat(t *testing.T) {
	config := &Config{
		Level:  logrus.DebugLevel,
		Format: "json",
		Output: "console",
		Fs:     afero.NewMemMapFs(), // Use memFS for testing
	}

	logger := NewLogger(config)

	if logger.GetLevel() != logrus.DebugLevel {
		t.Errorf("Expected DebugLevel, got %v", logger.GetLevel())
	}
}

func TestGetOutputWithMemFS(t *testing.T) {
	fs := afero.NewMemMapFs()

	// Test console output
	config := &Config{
		Level:    logrus.InfoLevel,
		Format:   FormatText,
		Output:   OutputConsole,
		FilePath: "/test/output.log",
		Fs:       fs,
	}

	output := getOutput(config)
	if output == nil {
		t.Error("Console output should not be nil")
	}

	// Test file output with memFS
	config.Output = OutputFile
	output = getOutput(config)
	if output == nil {
		t.Error("File output should not be nil")
	}

	// Verify file was created in memFS
	exists, err := afero.Exists(fs, "/test/output.log")
	if err != nil {
		t.Errorf("Error checking file existence: %v", err)
	}
	if !exists {
		t.Error("Log file should have been created in memFS")
	}

	// Test both output with memFS
	config.Output = OutputBoth
	output = getOutput(config)
	if output == nil {
		t.Error("Both output should not be nil")
	}
}

func TestConfigConstants(t *testing.T) {
	// Test format constants
	if FormatText != "text" {
		t.Errorf("Expected FormatText to be 'text', got %s", FormatText)
	}

	if FormatJSON != "json" {
		t.Errorf("Expected FormatJSON to be 'json', got %s", FormatJSON)
	}

	// Test output constants
	if OutputConsole != "console" {
		t.Errorf("Expected OutputConsole to be 'console', got %s", OutputConsole)
	}

	if OutputFile != "file" {
		t.Errorf("Expected OutputFile to be 'file', got %s", OutputFile)
	}

	if OutputBoth != "both" {
		t.Errorf("Expected OutputBoth to be 'both', got %s", OutputBoth)
	}
}
