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
