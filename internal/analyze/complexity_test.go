package analyze

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestComplexityService_ProvideTemplate(t *testing.T) {
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

	service := NewComplexityServiceWithFS(fs)
	output, err := service.ProvideTemplate()

	if err != nil {
		t.Fatalf("ProvideTemplate failed: %v", err)
	}

	if !strings.Contains(output, "# COMPLEXITY ANALYSIS TEMPLATE") {
		t.Error("Output missing template header")
	}
	if !strings.Contains(output, "Add user authentication") {
		t.Error("Output missing intent text")
	}
	if !strings.Contains(output, "auth.go") {
		t.Error("Output missing scope files")
	}
	if !strings.Contains(output, "Identify Domains") {
		t.Error("Output missing template content")
	}
}

func TestComplexityService_ProvideTemplate_MissingMission(t *testing.T) {
	fs := afero.NewMemMapFs()
	service := NewComplexityServiceWithFS(fs)
	_, err := service.ProvideTemplate()

	if err == nil {
		t.Error("Expected error for missing mission.md")
	}
}
