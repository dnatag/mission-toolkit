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

func TestTestService_readCurrentIntent(t *testing.T) {
	fs := afero.NewMemMapFs()

	missionContent := `---
id: test-123
---

## INTENT
Fix authentication bug in login handler

## SCOPE
handler.go`

	if err := afero.WriteFile(fs, ".mission/mission.md", []byte(missionContent), 0644); err != nil {
		t.Fatal(err)
	}

	service := NewTestServiceWithFS(fs)
	intent, err := service.readCurrentIntent()

	if err != nil {
		t.Fatalf("readCurrentIntent failed: %v", err)
	}

	expected := "Fix authentication bug in login handler"
	if intent != expected {
		t.Errorf("Expected intent %q, got %q", expected, intent)
	}
}

func TestTestService_readCurrentScope(t *testing.T) {
	fs := afero.NewMemMapFs()

	missionContent := `---
id: test-123
---

## INTENT
Add user authentication

## SCOPE
auth.go
handler.go
middleware.go`

	if err := afero.WriteFile(fs, ".mission/mission.md", []byte(missionContent), 0644); err != nil {
		t.Fatal(err)
	}

	service := NewTestServiceWithFS(fs)
	scope, err := service.readCurrentScope()

	if err != nil {
		t.Fatalf("readCurrentScope failed: %v", err)
	}

	if !strings.Contains(scope, "auth.go") {
		t.Error("Scope missing auth.go")
	}
	if !strings.Contains(scope, "handler.go") {
		t.Error("Scope missing handler.go")
	}
	if !strings.Contains(scope, "middleware.go") {
		t.Error("Scope missing middleware.go")
	}
}
