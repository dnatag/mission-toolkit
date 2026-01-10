package mission

import (
	"testing"

	"github.com/spf13/afero"
)

func TestReader_Read(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		wantErr    bool
		wantID     string
		wantStatus string
		wantTrack  int
		wantType   string
	}{
		{
			name: "valid mission with frontmatter",
			content: `---
id: test-123
type: WET
track: 2
iteration: 1
status: planned
---

## INTENT
Test intent

## SCOPE
file1.go
file2.go
`,
			wantErr:    false,
			wantID:     "test-123",
			wantStatus: "planned",
			wantTrack:  2,
			wantType:   "WET",
		},
		{
			name: "missing frontmatter",
			content: `## INTENT
Test intent without frontmatter
`,
			wantErr: true,
		},
		{
			name: "invalid yaml frontmatter",
			content: `---
id: test-123
invalid yaml: [
---

## INTENT
Test
`,
			wantErr: true,
		},
		{
			name: "legacy # MISSION format",
			content: `# MISSION

type: WET
track: 2
iteration: 1
status: completed

## INTENT
Test legacy mission format

## SCOPE
file1.go
file2.go
`,
			wantErr:    false,
			wantID:     "legacy-mission",
			wantStatus: "completed",
			wantTrack:  2,
			wantType:   "WET",
		},
		{
			name: "legacy # MISSION ARCHIVE format",
			content: `# MISSION ARCHIVE

type: DRY
track: 3
iteration: 2
status: completed
completed_at: 2025-12-21T13:25:17-05:00
duration_minutes: 15

## INTENT
Test archived mission format

## SCOPE
file1.go
file2.go
`,
			wantErr:    false,
			wantID:     "legacy-mission",
			wantStatus: "completed",
			wantTrack:  3,
			wantType:   "DRY",
		},
		{
			name: "legacy # MISSION: Title format",
			content: `# MISSION: Card List Merging Function

**Track**: 2 (Standard)
**Type**: WET
**Status**: completed
**Created**: 2025-12-23T17:23:37.869-05:00
**Completed**: 2025-12-23T17:42:16.976-05:00

## INTENT
Create a function to merge two lists of cards
`,
			wantErr:    false,
			wantID:     "legacy-mission",
			wantStatus: "", // No status: field found, only **Status**
			wantTrack:  0,  // No track: field found, only **Track**
			wantType:   "", // No type: field found, only **Type**
		},
		{
			name:    "empty file",
			content: "",
			wantErr: true,
		},
		{
			name: "only frontmatter opening",
			content: `---
id: test-123
type: WET
`,
			wantErr: true,
		},
		{
			name: "frontmatter without closing",
			content: `---
id: test-123
type: WET
track: 2

## INTENT
Missing closing frontmatter
`,
			wantErr: true,
		},
		{
			name: "malformed YAML - missing quotes",
			content: `---
id: test-123
type: WET
track: not a number
status: planned
---

## INTENT
Test malformed YAML
`,
			wantErr: true,
		},
		{
			name: "frontmatter with extra dashes",
			content: `---
id: test-123
type: WET
track: 2
status: planned
---

## INTENT
Test with extra dashes in content
--- this should not break parsing ---
`,
			wantErr:    false,
			wantID:     "test-123",
			wantStatus: "planned",
			wantTrack:  2,
			wantType:   "WET",
		},
		{
			name: "legacy format with missing sections",
			content: `# MISSION

type: WET
track: 2
status: planned
`,
			wantErr:    false,
			wantID:     "legacy-mission",
			wantStatus: "planned",
			wantTrack:  2,
			wantType:   "WET",
		},
		{
			name: "legacy format with only header",
			content: `# MISSION
`,
			wantErr:    false,
			wantID:     "legacy-mission",
			wantStatus: "",
			wantTrack:  0,
			wantType:   "",
		},
		{
			name: "whitespace only file",
			content: `   
	
   `,
			wantErr: true,
		},
		{
			name: "frontmatter with missing required fields",
			content: `---
type: WET
---

## INTENT
Missing ID field
`,
			wantErr:    false,
			wantID:     "",
			wantStatus: "",
			wantTrack:  0,
			wantType:   "WET",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			path := "mission.md"
			afero.WriteFile(fs, path, []byte(tt.content), 0644)

			reader := NewReader(fs)
			mission, err := reader.Read(path)

			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if mission.ID != tt.wantID {
					t.Errorf("Read() ID = %v, want %v", mission.ID, tt.wantID)
				}
				if mission.Status != tt.wantStatus {
					t.Errorf("Read() Status = %v, want %v", mission.Status, tt.wantStatus)
				}
				if mission.Track != tt.wantTrack {
					t.Errorf("Read() Track = %v, want %v", mission.Track, tt.wantTrack)
				}
				if mission.Type != tt.wantType {
					t.Errorf("Read() Type = %v, want %v", mission.Type, tt.wantType)
				}
			}
		})
	}
}

func TestReader_ParseEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name: "frontmatter with nested YAML",
			content: `---
id: test-123
metadata:
  nested: value
  list:
    - item1
    - item2
type: WET
---

## INTENT
Test nested YAML
`,
			wantErr: false,
		},
		{
			name: "legacy format with colon in content",
			content: `# MISSION

type: WET
track: 2

## INTENT
Test with colon: in content
This should not be parsed as metadata: value
`,
			wantErr: false,
		},
		{
			name: "frontmatter with unicode characters",
			content: `---
id: test-123-ñ
type: WET
track: 2
status: planned
---

## INTENT
Test unicode: ñáéíóú
`,
			wantErr: false,
		},
		{
			name: "very large frontmatter",
			content: `---
id: test-123
type: WET
track: 2
status: planned
description: "large content"
---

## INTENT
Test large frontmatter
`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			path := "mission.md"
			afero.WriteFile(fs, path, []byte(tt.content), 0644)

			reader := NewReader(fs)
			_, err := reader.Read(path)

			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReader_ReadNonexistentFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	reader := NewReader(fs)

	_, err := reader.Read("nonexistent.md")
	if err == nil {
		t.Error("Read() expected error for nonexistent file, got nil")
	}
}

func TestReader_GetMissionID(t *testing.T) {
	fs := afero.NewMemMapFs()
	content := `---
id: test-mission-123
type: WET
track: 2
iteration: 1
status: planned
---

## INTENT
Test content
`
	path := "mission.md"
	afero.WriteFile(fs, path, []byte(content), 0644)

	reader := NewReader(fs)
	id, err := reader.GetMissionID(path)
	if err != nil {
		t.Fatalf("GetMissionID() error = %v", err)
	}

	if id != "test-mission-123" {
		t.Errorf("GetMissionID() = %v, want test-mission-123", id)
	}
}

func TestReader_GetMissionStatus(t *testing.T) {
	fs := afero.NewMemMapFs()
	content := `---
id: test-123
type: WET
track: 2
iteration: 1
status: active
---

## INTENT
Test
`
	path := "mission.md"
	afero.WriteFile(fs, path, []byte(content), 0644)

	reader := NewReader(fs)
	status, err := reader.GetMissionStatus(path)
	if err != nil {
		t.Fatalf("GetMissionStatus() error = %v", err)
	}

	if status != "active" {
		t.Errorf("GetMissionStatus() = %v, want active", status)
	}
}

func TestReader_GetMissionIDNoFrontmatter(t *testing.T) {
	fs := afero.NewMemMapFs()
	content := `## INTENT
No frontmatter
`
	path := "mission.md"
	afero.WriteFile(fs, path, []byte(content), 0644)

	reader := NewReader(fs)
	_, err := reader.GetMissionID(path)
	if err == nil {
		t.Error("GetMissionID() expected error for missing frontmatter, got nil")
	}
}

func TestReader_LegacyFormatVariants(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		wantID     string
		wantStatus string
		wantType   string
	}{
		{
			name: "legacy with ID field",
			content: `# MISSION

id: custom-legacy-123
type: DRY
track: 3
status: active

## INTENT
Test with custom ID
`,
			wantID:     "custom-legacy-123",
			wantStatus: "active",
			wantType:   "DRY",
		},
		{
			name: "legacy with extra metadata",
			content: `# MISSION ARCHIVE

id: archived-123
type: WET
track: 2
status: completed
completed_at: 2025-12-21T13:25:17-05:00
duration_minutes: 15
parent_mission: parent-123

## INTENT
Test with extra metadata
`,
			wantID:     "archived-123",
			wantStatus: "completed",
			wantType:   "WET",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			path := "mission.md"
			afero.WriteFile(fs, path, []byte(tt.content), 0644)

			reader := NewReader(fs)
			mission, err := reader.Read(path)
			if err != nil {
				t.Fatalf("Read() error = %v", err)
			}

			if mission.ID != tt.wantID {
				t.Errorf("Read() ID = %v, want %v", mission.ID, tt.wantID)
			}
			if mission.Status != tt.wantStatus {
				t.Errorf("Read() Status = %v, want %v", mission.Status, tt.wantStatus)
			}
			if mission.Type != tt.wantType {
				t.Errorf("Read() Type = %v, want %v", mission.Type, tt.wantType)
			}
		})
	}
}

func TestReader_ReadIntent(t *testing.T) {
	fs := afero.NewMemMapFs()
	reader := NewReader(fs)

	content := `---
id: test-123
status: planned
---

## INTENT
Add user authentication

## SCOPE
auth.go`

	afero.WriteFile(fs, "mission.md", []byte(content), 0644)

	intent, err := reader.ReadIntent("mission.md")
	if err != nil {
		t.Fatalf("ReadIntent failed: %v", err)
	}

	expected := "Add user authentication"
	if intent != expected {
		t.Errorf("Expected intent %q, got %q", expected, intent)
	}
}

func TestReader_ReadScope(t *testing.T) {
	fs := afero.NewMemMapFs()
	reader := NewReader(fs)

	content := `---
id: test-123
status: planned
---

## INTENT
Add user authentication

## SCOPE
auth.go
handler.go
middleware.go`

	afero.WriteFile(fs, "mission.md", []byte(content), 0644)

	scope, err := reader.ReadScope("mission.md")
	if err != nil {
		t.Fatalf("ReadScope failed: %v", err)
	}

	if !contains(scope, "auth.go") {
		t.Error("Scope missing auth.go")
	}
	if !contains(scope, "handler.go") {
		t.Error("Scope missing handler.go")
	}
	if !contains(scope, "middleware.go") {
		t.Error("Scope missing middleware.go")
	}
}

func TestReader_ReadIntent_MissingSection(t *testing.T) {
	fs := afero.NewMemMapFs()
	reader := NewReader(fs)

	content := `---
id: test-123
---

## SCOPE
auth.go`

	afero.WriteFile(fs, "mission.md", []byte(content), 0644)

	_, err := reader.ReadIntent("mission.md")
	if err == nil {
		t.Error("Expected error for missing INTENT section")
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) >= len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
