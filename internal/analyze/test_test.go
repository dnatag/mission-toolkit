package analyze

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestTestService_ProvideTemplate(t *testing.T) {
	fs := afero.NewMemMapFs()

	missionContent := `---
id: test-123
status: planned
---

## INTENT
Add user authentication

## SCOPE
auth.go
handler.go`

	if err := afero.WriteFile(fs, ".mission/mission.md", []byte(missionContent), 0644); err != nil {
		t.Fatal(err)
	}

	service := NewTestServiceWithFS(fs)
	output, err := service.ProvideTemplate()

	if err != nil {
		t.Fatalf("ProvideTemplate failed: %v", err)
	}

	if !strings.Contains(output, "## Current Mission Context") {
		t.Error("Output missing Current Mission Context section")
	}
	if !strings.Contains(output, "Add user authentication") {
		t.Error("Output missing intent text")
	}
	if !strings.Contains(output, "auth.go") {
		t.Error("Output missing scope files")
	}
	if !strings.Contains(output, "## Test Analysis Instructions") {
		t.Error("Output missing template content")
	}
}

func TestTestService_ProvideTemplate_MissingMission(t *testing.T) {
	fs := afero.NewMemMapFs()
	service := NewTestServiceWithFS(fs)
	_, err := service.ProvideTemplate()

	if err == nil {
		t.Error("Expected error for missing mission.md")
	}
}
