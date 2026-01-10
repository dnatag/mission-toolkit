package analyze

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestClarifyService_ProvideTemplate(t *testing.T) {
	fs := afero.NewMemMapFs()

	missionContent := `---
id: test-123
status: planned
---

## INTENT
Add user authentication

## SCOPE
auth.go`

	if err := afero.WriteFile(fs, ".mission/mission.md", []byte(missionContent), 0644); err != nil {
		t.Fatal(err)
	}

	service := NewClarifyServiceWithFS(fs)
	output, err := service.ProvideTemplate()

	if err != nil {
		t.Fatalf("ProvideTemplate failed: %v", err)
	}

	if !strings.Contains(output, "## Current Intent") {
		t.Error("Output missing Current Intent section")
	}
	if !strings.Contains(output, "Add user authentication") {
		t.Error("Output missing intent text")
	}
	if !strings.Contains(output, "Scan user intent") {
		t.Error("Output missing template content")
	}
}

func TestClarifyService_ProvideTemplate_MissingMission(t *testing.T) {
	fs := afero.NewMemMapFs()
	service := NewClarifyServiceWithFS(fs)
	_, err := service.ProvideTemplate()

	if err == nil {
		t.Error("Expected error for missing mission.md")
	}
}
