package plan

import (
	"encoding/json"
	"testing"
)

func TestNewOutputResponse(t *testing.T) {
	data := map[string]interface{}{
		"success": true,
		"message": "test message",
	}
	nextStep := "PROCEED to next step"

	output := NewOutputResponse(data, nextStep)

	if output.NextStep != nextStep {
		t.Errorf("Expected NextStep %s, got %s", nextStep, output.NextStep)
	}

	if output.Data == nil {
		t.Error("Expected Data to be set")
	}
}

func TestOutputResponse_FormatOutput_WithMap(t *testing.T) {
	data := map[string]interface{}{
		"valid":  true,
		"errors": []string{},
	}
	nextStep := "PROCEED to validation"

	output := NewOutputResponse(data, nextStep)
	jsonStr, err := output.FormatOutput()

	if err != nil {
		t.Fatalf("FormatOutput failed: %v", err)
	}

	// Parse back to verify next_step was added
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if result["next_step"] != nextStep {
		t.Errorf("Expected next_step %s, got %v", nextStep, result["next_step"])
	}

	if result["valid"] != true {
		t.Errorf("Expected valid true, got %v", result["valid"])
	}
}

func TestOutputResponse_FormatOutput_WithStruct(t *testing.T) {
	type TestStruct struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	data := TestStruct{
		Success: true,
		Message: "test message",
	}
	nextStep := "PROCEED to next step"

	output := NewOutputResponse(data, nextStep)
	jsonStr, err := output.FormatOutput()

	if err != nil {
		t.Fatalf("FormatOutput failed: %v", err)
	}

	// Parse back to verify next_step was added
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if result["next_step"] != nextStep {
		t.Errorf("Expected next_step %s, got %v", nextStep, result["next_step"])
	}

	if result["success"] != true {
		t.Errorf("Expected success true, got %v", result["success"])
	}

	if result["message"] != "test message" {
		t.Errorf("Expected message 'test message', got %v", result["message"])
	}
}
