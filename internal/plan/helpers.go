package plan

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/afero"
)

// OutputResponse standardizes CLI output with next_step guidance
type OutputResponse struct {
	Data     interface{} `json:",inline"`
	NextStep string      `json:"next_step"`
}

// NewOutputResponse creates a standardized output response
func NewOutputResponse(data interface{}, nextStep string) OutputResponse {
	return OutputResponse{
		Data:     data,
		NextStep: nextStep,
	}
}

// FormatOutput converts OutputResponse to JSON string
func (o OutputResponse) FormatOutput() (string, error) {
	// Merge the data fields with next_step
	if dataMap, ok := o.Data.(map[string]interface{}); ok {
		dataMap["next_step"] = o.NextStep
		return ToJSON(dataMap)
	}

	// For structs, create a new map
	dataBytes, err := json.Marshal(o.Data)
	if err != nil {
		return "", err
	}

	var dataMap map[string]interface{}
	if err := json.Unmarshal(dataBytes, &dataMap); err != nil {
		return "", err
	}

	dataMap["next_step"] = o.NextStep
	return ToJSON(dataMap)
}

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
