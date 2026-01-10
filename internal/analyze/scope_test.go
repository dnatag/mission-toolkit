package analyze

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestScopeService_ProvideTemplate(t *testing.T) {
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

	service := NewScopeServiceWithFS(fs)
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
	if !strings.Contains(output, "Determine which implementation files") {
		t.Error("Output missing template content")
	}
}

func TestScopeService_ProvideTemplate_MissingMission(t *testing.T) {
	fs := afero.NewMemMapFs()
	service := NewScopeServiceWithFS(fs)
	_, err := service.ProvideTemplate()

	if err == nil {
		t.Error("Expected error for missing mission.md")
	}
}
