package logger

import (
	"testing"

	"github.com/spf13/afero"
)

func TestNew(t *testing.T) {
	missionID := "test-mission-123"

	// Use memFS for testing
	config := DefaultConfig()
	config.Fs = afero.NewMemMapFs()
	logger := NewWithConfig(missionID, config)

	if logger == nil {
		t.Error("Logger should not be nil")
	}

	if logger.missionID != missionID {
		t.Errorf("Expected mission ID %s, got %s", missionID, logger.missionID)
	}
}

func TestLogStep(t *testing.T) {
	missionID := "test-mission-123"

	// Use memFS for testing
	config := DefaultConfig()
	config.Fs = afero.NewMemMapFs()
	logger := NewWithConfig(missionID, config)

	// Test different log levels
	testCases := []struct {
		level   string
		step    string
		message string
	}{
		{"SUCCESS", "Test Step", "Test message"},
		{"ERROR", "Error Step", "Error message"},
		{"WARN", "Warning Step", "Warning message"},
		{"DEBUG", "Debug Step", "Debug message"},
		{"INFO", "Info Step", "Info message"},
	}

	for _, tc := range testCases {
		t.Run(tc.level, func(t *testing.T) {
			// This test just ensures no panic occurs
			logger.LogStep(tc.level, tc.step, tc.message)
		})
	}
}

func TestSuccessAndFailed(t *testing.T) {
	missionID := "test-mission-123"

	// Use memFS for testing
	config := DefaultConfig()
	config.Fs = afero.NewMemMapFs()
	logger := NewWithConfig(missionID, config)

	// Test convenience methods
	logger.Success("Test Step", "Success message")
	logger.Failed("Test Step", "Failure message")
}
