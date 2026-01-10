package analyze

import (
	"strings"
	"testing"
)

func TestIntentService_ProvideTemplate(t *testing.T) {
	service := NewIntentService()
	output, err := service.ProvideTemplate("add auth")

	if err != nil {
		t.Fatalf("ProvideTemplate failed: %v", err)
	}

	// Verify output contains all expected sections
	if !strings.Contains(output, "## User Input") {
		t.Error("Output missing User Input section")
	}
	if !strings.Contains(output, "add auth") {
		t.Error("Output missing user input text")
	}
	if !strings.Contains(output, "Distill raw user input") {
		t.Error("Output missing template content")
	}
}

func TestIntentService_ProvideTemplate_EmptyInput(t *testing.T) {
	service := NewIntentService()
	output, err := service.ProvideTemplate("")

	if err != nil {
		t.Fatalf("ProvideTemplate failed: %v", err)
	}

	if !strings.Contains(output, "## User Input") {
		t.Error("Output missing User Input section")
	}
}
