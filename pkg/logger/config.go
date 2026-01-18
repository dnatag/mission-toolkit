package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// Config holds logger configuration
type Config struct {
	Level    logrus.Level
	Format   string // "text" or "json"
	Output   string // "console", "file", or "both"
	FilePath string
	Fs       afero.Fs // Filesystem interface for testing
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
		Output:   OutputBoth,
		FilePath: ".mission/execution.log",
		Fs:       afero.NewOsFs(), // Use real filesystem by default
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
		if file := openLogFile(config.Fs, config.FilePath); file != nil {
			return file
		}
		return os.Stdout
	case OutputBoth:
		if file := openLogFile(config.Fs, config.FilePath); file != nil {
			return io.MultiWriter(os.Stdout, file)
		}
		return os.Stdout
	default:
		return os.Stdout
	}
}

// openLogFile opens log file for writing using the provided filesystem, returns nil on error
func openLogFile(fs afero.Fs, filePath string) afero.File {
	// Return nil if filePath is empty or invalid to prevent file creation
	if filePath == "" || filePath == "/dev/null" {
		return nil
	}
	if err := fs.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil
	}
	file, err := fs.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil
	}
	return file
}
