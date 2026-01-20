package logger

import (
	"github.com/sirupsen/logrus"
)

// Logger wraps logrus with mission-specific functionality
type Logger struct {
	*logrus.Logger
	missionID string
}

const (
	LevelSuccess = "SUCCESS"
	LevelError   = "ERROR"
	LevelFailed  = "FAILED"
	LevelWarn    = "WARN"
	LevelDebug   = "DEBUG"
)

// New creates a new mission logger
func New(missionID string) *Logger {
	return NewWithConfig(missionID, DefaultConfig())
}

// NewWithConfig creates a new mission logger with custom config
func NewWithConfig(missionID string, config *Config) *Logger {
	return &Logger{
		Logger:    NewLogger(config),
		missionID: missionID,
	}
}

// LogStep logs a mission step with structured format
func (l *Logger) LogStep(level, step, message string) {
	fields := logrus.Fields{
		"mission_id": l.missionID,
		"step":       step,
		"level":      level,
	}

	switch level {
	case LevelSuccess:
		l.WithFields(fields).Info(message)
	case LevelError, LevelFailed:
		l.WithFields(fields).Error(message)
	case LevelWarn:
		l.WithFields(fields).Warn(message)
	case LevelDebug:
		l.WithFields(fields).Debug(message)
	default:
		l.WithFields(fields).Info(message)
	}
}

// Success logs a success message
func (l *Logger) Success(step, message string) {
	l.LogStep(LevelSuccess, step, message)
}

// Failed logs a failure message
func (l *Logger) Failed(step, message string) {
	l.LogStep(LevelFailed, step, message)
}
