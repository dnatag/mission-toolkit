package mission

import (
	"encoding/json"
	"testing"

	"github.com/spf13/afero"
)

func TestFinalizeService_Finalize_Valid(t *testing.T) {
	fs := afero.NewMemMapFs()

	missionContent := `---
id: test-123
status: planned
track: 2
type: WET
---

## INTENT
Add user authentication

## SCOPE
auth.go
handler.go

## PLAN
- [ ] 1. Create auth handler
- [ ] 2. Add middleware

## VERIFICATION
go test ./...`

	afero.WriteFile(fs, ".mission/mission.md", []byte(missionContent), 0644)

	service := NewFinalizeService(fs, ".mission")
	output, err := service.Finalize()

	if err != nil {
		t.Fatalf("Finalize failed: %v", err)
	}

	var result FinalizeResult
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if !result.Valid {
		t.Error("Expected valid mission")
	}
	if result.Message != "Mission validated successfully" {
		t.Errorf("Expected success message, got: %s", result.Message)
	}
}

func TestFinalizeService_Finalize_MissingSection(t *testing.T) {
	fs := afero.NewMemMapFs()

	missionContent := `---
id: test-123
---

## INTENT
Add user authentication

## SCOPE
auth.go`

	afero.WriteFile(fs, ".mission/mission.md", []byte(missionContent), 0644)

	service := NewFinalizeService(fs, ".mission")
	output, err := service.Finalize()

	if err != nil {
		t.Fatalf("Finalize failed: %v", err)
	}

	var result FinalizeResult
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if result.Valid {
		t.Error("Expected invalid mission")
	}
	if len(result.MissingSections) == 0 {
		t.Error("Expected missing sections")
	}
	if !containsString(result.MissingSections, "PLAN") {
		t.Error("Expected PLAN in missing sections")
	}
}

func TestFinalizeService_Finalize_EmptySection(t *testing.T) {
	fs := afero.NewMemMapFs()

	missionContent := `---
id: test-123
---

## INTENT
Add user authentication

## SCOPE

## PLAN
- [ ] 1. Create handler

## VERIFICATION
go test ./...`

	afero.WriteFile(fs, ".mission/mission.md", []byte(missionContent), 0644)

	service := NewFinalizeService(fs, ".mission")
	output, err := service.Finalize()

	if err != nil {
		t.Fatalf("Finalize failed: %v", err)
	}

	var result FinalizeResult
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if result.Valid {
		t.Error("Expected invalid mission")
	}
	if len(result.EmptySections) == 0 {
		t.Error("Expected empty sections")
	}
	if !containsString(result.EmptySections, "SCOPE") {
		t.Error("Expected SCOPE in empty sections")
	}
}

func TestFinalizeService_Finalize_MissingFile(t *testing.T) {
	fs := afero.NewMemMapFs()

	service := NewFinalizeService(fs, ".mission")
	_, err := service.Finalize()

	if err == nil {
		t.Error("Expected error for missing mission file")
	}
}

func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func TestFinalizeService_Finalize_CleansUpTemplates(t *testing.T) {
	fs := afero.NewMemMapFs()

	missionContent := `---
id: test-123
status: planning
track: 2
type: WET
---

## INTENT
Add user authentication

## SCOPE
auth.go

## PLAN
- [ ] 1. Create handler

## VERIFICATION
go test ./...`

	afero.WriteFile(fs, ".mission/mission.md", []byte(missionContent), 0644)
	fs.Mkdir(".mission/templates", 0755)
	afero.WriteFile(fs, ".mission/templates/test.md", []byte("test"), 0644)

	service := NewFinalizeService(fs, ".mission")
	_, err := service.Finalize()

	if err != nil {
		t.Fatalf("Finalize failed: %v", err)
	}

	exists, _ := afero.DirExists(fs, ".mission/templates")
	if exists {
		t.Error("Expected templates directory to be removed")
	}
}

func TestFinalizeService_Finalize_UpdatesStatusToPlanned(t *testing.T) {
	fs := afero.NewMemMapFs()

	missionContent := `---
id: test-123
status: planning
track: 2
type: WET
---

## INTENT
Add user authentication

## SCOPE
auth.go

## PLAN
- [ ] 1. Create handler

## VERIFICATION
go test ./...`

	afero.WriteFile(fs, ".mission/mission.md", []byte(missionContent), 0644)

	service := NewFinalizeService(fs, ".mission")
	_, err := service.Finalize()

	if err != nil {
		t.Fatalf("Finalize failed: %v", err)
	}

	reader := NewReader(fs, ".mission/mission.md")
	mission, err := reader.Read()
	if err != nil {
		t.Fatalf("Failed to read mission: %v", err)
	}

	if mission.Status != "planned" {
		t.Errorf("Expected status 'planned', got: %s", mission.Status)
	}
}
