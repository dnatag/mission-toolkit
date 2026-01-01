package plan

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/afero"
)

// ServiceBase provides common functionality for plan services
type ServiceBase struct {
	fs        afero.Fs
	missionID string
}

// NewServiceBase creates a new service base with common dependencies
func NewServiceBase(fs afero.Fs, missionID string) ServiceBase {
	return ServiceBase{
		fs:        fs,
		missionID: missionID,
	}
}

// JSONMarshaler interface for types that can be marshaled to JSON
type JSONMarshaler interface {
	ToJSON() (string, error)
}

// MarshalToJSON provides common JSON marshaling functionality
func MarshalToJSON(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	return string(data), nil
}

// ErrorResult represents a common error result structure
type ErrorResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// NewErrorResult creates a standardized error result
func NewErrorResult(format string, args ...interface{}) *ErrorResult {
	return &ErrorResult{
		Success: false,
		Message: fmt.Sprintf(format, args...),
	}
}

// ToJSON converts ErrorResult to JSON string
func (r *ErrorResult) ToJSON() (string, error) {
	return MarshalToJSON(r)
}
