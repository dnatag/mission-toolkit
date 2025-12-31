package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Level != logrus.InfoLevel {
		t.Errorf("Expected InfoLevel, got %v", config.Level)
	}

	if config.Format != "text" {
		t.Errorf("Expected text format, got %s", config.Format)
	}

	if config.Output != "console" {
		t.Errorf("Expected console output, got %s", config.Output)
	}
}

func TestNewLogger(t *testing.T) {
	config := DefaultConfig()
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
	}

	logger := NewLogger(config)

	if logger.GetLevel() != logrus.DebugLevel {
		t.Errorf("Expected DebugLevel, got %v", logger.GetLevel())
	}
}
