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

// ToJSON converts any struct to JSON string
func ToJSON(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	return string(data), nil
}
