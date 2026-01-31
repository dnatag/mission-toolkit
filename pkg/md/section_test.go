package md

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindSection(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		sectionName string
		want        int
	}{
		{
			name: "section exists",
			body: `## INTENT
Content here

## SCOPE
More content`,
			sectionName: "INTENT",
			want:        0,
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
			want:        3,
		},
		{
			name: "section not found",
			body: `## INTENT
Content`,
			sectionName: "MISSING",
			want:        -1,
		},
		{
			name: "case insensitive match",
			body: `## INTENT
Content`,
			sectionName: "intent",
			want:        0,
		},
		{
			name:        "empty body",
			body:        "",
			sectionName: "INTENT",
			want:        -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindSection(tt.body, tt.sectionName)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExtractSection(t *testing.T) {
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
			got := ExtractSection(tt.body, tt.sectionName)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExtractList(t *testing.T) {
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
			got := ExtractList(tt.body, tt.sectionName)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSkipToNextSection(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		startIndex int
		want       int
	}{
		{
			name: "next section exists",
			body: `## INTENT
Content

## SCOPE
Files`,
			startIndex: 0,
			want:       3,
		},
		{
			name: "multiple sections",
			body: `## INTENT
Content

## SCOPE
Files

## PLAN
Steps`,
			startIndex: 3,
			want:       6,
		},
		{
			name: "no next section",
			body: `## INTENT
Content
More content`,
			startIndex: 0,
			want:       -1,
		},
		{
			name: "start at last section",
			body: `## INTENT
Content

## SCOPE
Files`,
			startIndex: 3,
			want:       -1,
		},
		{
			name:       "empty body",
			body:       "",
			startIndex: 0,
			want:       -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SkipToNextSection(tt.body, tt.startIndex)
			assert.Equal(t, tt.want, got)
		})
	}
}
