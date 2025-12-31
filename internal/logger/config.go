package logger

import (
	"io"
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

const (
	FormatText    = "text"
	FormatJSON    = "json"
	OutputConsole = "console"
	OutputFile    = "file"
	OutputBoth    = "both"
)

// DefaultConfig returns default logger configuration
func DefaultConfig() *Config {
	return &Config{
		Level:    logrus.InfoLevel,
		Format:   FormatText,
		Output:   OutputConsole,
		FilePath: ".mission/execution.log",
	}
}

// NewLogger creates a new logger with the given configuration
func NewLogger(config *Config) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(config.Level)

	// Set formatter
	if config.Format == FormatJSON {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}

	// Set output
	logger.SetOutput(getOutput(config))
	return logger
}

// getOutput returns the appropriate output writer based on config
func getOutput(config *Config) io.Writer {
	switch config.Output {
	case OutputFile:
		if file := openLogFile(config.FilePath); file != nil {
			return file
		}
		fallthrough
	case OutputBoth:
		// TODO: Implement multi-writer for both console and file
		fallthrough
	default:
		return os.Stdout
	}
}

// openLogFile opens log file for writing, returns nil on error
func openLogFile(filePath string) *os.File {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil
	}
	return file
}
