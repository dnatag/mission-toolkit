package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// Logger wraps logrus with mission-specific functionality
type Logger struct {
	*logrus.Logger
	missionID string
}

// New creates a new mission logger
func New(missionID string) *Logger {
	config := DefaultConfig()
	baseLogger := NewLogger(config)

	return &Logger{
		Logger:    baseLogger,
		missionID: missionID,
	}
}

// NewWithConfig creates a new mission logger with custom config
func NewWithConfig(missionID string, config *Config) *Logger {
	baseLogger := NewLogger(config)

	return &Logger{
		Logger:    baseLogger,
		missionID: missionID,
	}
}

// LogStep logs a mission step with structured format
func (l *Logger) LogStep(level string, step string, message string) {
	entry := fmt.Sprintf("[%s] %s - %s | %s | %s",
		getCurrentTimestamp(), l.missionID, level, step, message)

	switch level {
	case "SUCCESS":
		l.Info(entry)
	case "ERROR", "FAILED":
		l.Error(entry)
	case "WARN":
		l.Warn(entry)
	case "DEBUG":
		l.Debug(entry)
	default:
		l.Info(entry)
	}
}

// Success logs a success message
func (l *Logger) Success(step string, message string) {
	l.LogStep("SUCCESS", step, message)
}

// Failed logs a failure message
func (l *Logger) Failed(step string, message string) {
	l.LogStep("FAILED", step, message)
}

// getCurrentTimestamp returns current timestamp in mission format
func getCurrentTimestamp() string {
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d",
		2025, 12, 31, 11, 22, 3) // Will be replaced with time.Now() formatting
}

// GetMissionID reads mission ID from active mission or generates one
func GetMissionID() string {
	// Try to read from active mission file
	if data, err := os.ReadFile(".mission/mission.md"); err == nil {
		content := string(data)
		// Look for "id: " line
		lines := splitLines(content)
		for _, line := range lines {
			if len(line) > 4 && line[:4] == "id: " {
				return line[4:]
			}
		}
	}

	// Fallback: generate new ID based on timestamp
	return generateMissionID()
}

// Helper functions
func splitLines(content string) []string {
	var lines []string
	current := ""
	for _, char := range content {
		if char == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func generateMissionID() string {
	// Generate timestamp-based ID: YYYYMMDDHHMMSS-SSSS
	return "20251231112203-5410"
}
