// Package validation provides input validation for mission planning.
package validation

import "strings"

// CheckResult represents the result of input validation
type CheckResult struct {
	IsValid  bool   `json:"is_valid"`
	Message  string `json:"message"`
	NextStep string `json:"next_step"`
}

// Validate checks if input is valid for mission planning
func Validate(input string) *CheckResult {
	trimmed := strings.TrimSpace(input)

	if trimmed == "" {
		return &CheckResult{
			IsValid:  false,
			Message:  "Input is empty or whitespace",
			NextStep: "ASK_USER: What is your intent or goal for this task?",
		}
	}

	if trimmed == "$ARGUMENTS" {
		return &CheckResult{
			IsValid:  false,
			Message:  "Input is a placeholder - no intent provided",
			NextStep: "ASK_USER: What is your intent or goal for this task?",
		}
	}

	return &CheckResult{
		IsValid:  true,
		Message:  "Input is valid",
		NextStep: "PROCEED with execution",
	}
}
