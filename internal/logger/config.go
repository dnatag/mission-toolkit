package logger

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// Config holds logger configuration
type Config struct {
	Level    logrus.Level
	Format   string // "text" or "json"
	Output   string // "console", "file", or "both"
	FilePath string
}

// DefaultConfig returns default logger configuration
func DefaultConfig() *Config {
	return &Config{
		Level:    logrus.InfoLevel,
		Format:   "text",
		Output:   "console",
		FilePath: ".mission/execution.log",
	}
}

// NewLogger creates a new logger with the given configuration
func NewLogger(config *Config) *logrus.Logger {
	logger := logrus.New()

	// Set log level
	logger.SetLevel(config.Level)

	// Set formatter
	if config.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// Set output
	switch config.Output {
	case "file":
		if err := ensureDir(config.FilePath); err == nil {
			if file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
				logger.SetOutput(file)
			}
		}
	case "both":
		// For now, default to console. Multi-writer can be added later
		logger.SetOutput(os.Stdout)
	default:
		logger.SetOutput(os.Stdout)
	}

	return logger
}

// ensureDir creates directory if it doesn't exist
func ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	return os.MkdirAll(dir, 0755)
}
