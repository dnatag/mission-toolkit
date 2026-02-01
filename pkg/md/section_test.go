package md

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocument_HasSection(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		sectionName string
		want        bool
	}{
		{
			name: "section exists",
			body: `## INTENT
Content here

## SCOPE
More content`,
			sectionName: "INTENT",
			want:        true,
		},
		{
			name: "section in middle",
			body: `## INTENT
Content

## SCOPE
Files here

## PLAN
Steps`,
			sectionName: "SCOPE",
			want:        true,
		},
		{
			name: "section not found",
			body: `## INTENT
Content`,
			sectionName: "MISSING",
			want:        false,
		},
		{
			name: "case insensitive match",
			body: `## INTENT
Content`,
			sectionName: "intent",
			want:        true,
		},
		{
			name:        "empty body",
			body:        "",
			sectionName: "INTENT",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &Document{Body: tt.body}
			got := doc.HasSection(tt.sectionName)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDocument_GetSection(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		sectionName string
		want        string
	}{
		{
			name: "extract simple content",
			body: `## INTENT
Add user authentication

## SCOPE
auth.go`,
			sectionName: "INTENT",
			want:        "Add user authentication",
		},
		{
			name: "extract multiline content",
			body: `## INTENT
Line 1
Line 2
Line 3

## SCOPE
Files`,
			sectionName: "INTENT",
			want:        "Line 1\nLine 2\nLine 3",
		},
		{
			name: "section at end",
			body: `## INTENT
Content

## VERIFICATION
go test ./...`,
			sectionName: "VERIFICATION",
			want:        "go test ./...",
		},
		{
			name: "section not found",
			body: `## INTENT
Content`,
			sectionName: "MISSING",
			want:        "",
		},
		{
			name: "empty section",
			body: `## INTENT

## SCOPE
Files`,
			sectionName: "INTENT",
			want:        "",
		},
		{
			name:        "empty body",
			body:        "",
			sectionName: "INTENT",
			want:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &Document{Body: tt.body}
			got, err := doc.GetSection(tt.sectionName)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDocument_GetList(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		sectionName string
		want        []string
	}{
		{
			name: "dash list",
			body: `## SCOPE
- file1.go
- file2.go
- file3.go`,
			sectionName: "SCOPE",
			want:        []string{"file1.go", "file2.go", "file3.go"},
		},
		{
			name: "asterisk list",
			body: `## SCOPE
* file1.go
* file2.go`,
			sectionName: "SCOPE",
			want:        []string{"file1.go", "file2.go"},
		},
		{
			name: "numbered list",
			body: `## PLAN
1. First step
2. Second step
3. Third step`,
			sectionName: "PLAN",
			want:        []string{"First step", "Second step", "Third step"},
		},
		{
			name: "checkbox list",
			body: `## PLAN
- [ ] Todo item
- [x] Done item
- [ ] Another todo`,
			sectionName: "PLAN",
			want:        []string{"Todo item", "Done item", "Another todo"},
		},
		{
			name: "mixed list formats",
			body: `## SCOPE
- file1.go
* file2.go
- [ ] file3.go`,
			sectionName: "SCOPE",
			want:        []string{"file1.go", "file2.go", "file3.go"},
		},
		{
			name: "list with empty lines",
			body: `## SCOPE
- file1.go

- file2.go

- file3.go`,
			sectionName: "SCOPE",
			want:        []string{"file1.go", "file2.go", "file3.go"},
		},
		{
			name: "section not found",
			body: `## INTENT
Content`,
			sectionName: "SCOPE",
			want:        []string{},
		},
		{
			name: "empty section",
			body: `## SCOPE

## PLAN
Steps`,
			sectionName: "SCOPE",
			want:        []string{},
		},
		{
			name: "no list items",
			body: `## SCOPE
Just plain text
No list markers`,
			sectionName: "SCOPE",
			want:        []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &Document{Body: tt.body}
			got, err := doc.GetList(tt.sectionName)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDocument_ListSections(t *testing.T) {
	tests := []struct {
		name string
		body string
		want []string
	}{
		{
			name: "multiple sections",
			body: `## INTENT
Content

## SCOPE
Files

## PLAN
Steps`,
			want: []string{"INTENT", "SCOPE", "PLAN"},
		},
		{
			name: "single section",
			body: `## INTENT
Content`,
			want: []string{"INTENT"},
		},
		{
			name: "no sections",
			body: `Just plain text
No sections here`,
			want: []string{},
		},
		{
			name: "empty body",
			body: "",
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &Document{Body: tt.body}
			got := doc.ListSections()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDocument_UpdateSectionContent(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		sectionName string
		content     string
		want        string
	}{
		{
			name: "replace existing section",
			body: `## INTENT
Old content

## SCOPE
Files`,
			sectionName: "INTENT",
			content:     "New content",
			want: `## INTENT
New content

## SCOPE
Files`,
		},
		{
			name: "replace middle section",
			body: `## INTENT
Content

## SCOPE
Old files

## PLAN
Steps`,
			sectionName: "SCOPE",
			content:     "New files",
			want: `## INTENT
Content

## SCOPE
New files

## PLAN
Steps`,
		},
		{
			name: "replace last section",
			body: `## INTENT
Content

## SCOPE
Old files`,
			sectionName: "SCOPE",
			content:     "New files",
			want: `## INTENT
Content

## SCOPE
New files`,
		},
		{
			name:        "add new section to empty body",
			body:        "",
			sectionName: "INTENT",
			content:     "New content",
			want: `
## INTENT
New content`,
		},
		{
			name: "add new section to existing body",
			body: `## INTENT
Content`,
			sectionName: "SCOPE",
			content:     "Files",
			want: `## INTENT
Content

## SCOPE
Files`,
		},
		{
			name: "case insensitive section name",
			body: `## INTENT
Old content`,
			sectionName: "intent",
			content:     "New content",
			want: `## INTENT
New content`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &Document{Body: tt.body}
			err := doc.UpdateSectionContent(tt.sectionName, tt.content)
			require.NoError(t, err)
			assert.Equal(t, tt.want, doc.Body)
		})
	}
}

func TestDocument_UpdateSectionList(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		sectionName string
		items       []string
		want        string
	}{
		{
			name: "replace list items",
			body: `## SCOPE
- file1.go
- file2.go

## PLAN
Steps`,
			sectionName: "SCOPE",
			items:       []string{"file3.go", "file4.go"},
			want: `## SCOPE
- file3.go
- file4.go

## PLAN
Steps`,
		},
		{
			name:        "create new list section",
			body:        "## INTENT\nContent",
			sectionName: "SCOPE",
			items:       []string{"file1.go", "file2.go"},
			want: `## INTENT
Content

## SCOPE
- file1.go
- file2.go`,
		},
		{
			name:        "empty items list",
			body:        "## SCOPE\n- file1.go",
			sectionName: "SCOPE",
			items:       []string{},
			want:        "## SCOPE\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &Document{Body: tt.body}
			err := doc.UpdateSectionList(tt.sectionName, tt.items)
			require.NoError(t, err)
			assert.Equal(t, tt.want, doc.Body)
		})
	}
}

func TestDocument_AppendSectionList(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		sectionName string
		items       []string
		want        string
	}{
		{
			name: "append list items",
			body: `## SCOPE
- file1.go
- file2.go

## PLAN
Steps`,
			sectionName: "SCOPE",
			items:       []string{"file3.go"},
			want: `## SCOPE
- file1.go
- file2.go
- file3.go

## PLAN
Steps`,
		},
		{
			name: "append to empty section",
			body: `## SCOPE

## PLAN
Steps`,
			sectionName: "SCOPE",
			items:       []string{"file1.go"},
			want: `## SCOPE
- file1.go

## PLAN
Steps`,
		},
		{
			name:        "append to new section",
			body:        "## INTENT\nContent",
			sectionName: "SCOPE",
			items:       []string{"file1.go", "file2.go"},
			want: `## INTENT
Content

## SCOPE
- file1.go
- file2.go`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &Document{Body: tt.body}
			err := doc.AppendSectionList(tt.sectionName, tt.items)
			require.NoError(t, err)
			assert.Equal(t, tt.want, doc.Body)
		})
	}
}

func TestDocument_ValidateSectionName(t *testing.T) {
	tests := []struct {
		name        string
		sectionName string
		wantErr     bool
	}{
		{
			name:        "valid name",
			sectionName: "INTENT",
			wantErr:     false,
		},
		{
			name:        "valid with spaces",
			sectionName: "MY SECTION",
			wantErr:     false,
		},
		{
			name:        "empty name",
			sectionName: "",
			wantErr:     true,
		},
		{
			name:        "contains newline",
			sectionName: "INTENT\nSCOPE",
			wantErr:     true,
		},
		{
			name:        "contains tab",
			sectionName: "INTENT\tSCOPE",
			wantErr:     true,
		},
		{
			name:        "too long",
			sectionName: strings.Repeat("A", 101),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &Document{Body: "## TEST\nContent"}
			_, err := doc.GetSection(tt.sectionName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
