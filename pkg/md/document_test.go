package md

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantMeta map[string]interface{}
		wantBody string
		wantErr  bool
	}{
		{
			name: "valid frontmatter and body",
			content: `---
title: Test Document
author: John Doe
---

# Hello World

This is the body.`,
			wantMeta: map[string]interface{}{
				"title":  "Test Document",
				"author": "John Doe",
			},
			wantBody: "\n# Hello World\n\nThis is the body.",
			wantErr:  false,
		},
		{
			name: "no frontmatter",
			content: `# Hello World

This is just body content.`,
			wantMeta: map[string]interface{}{},
			wantBody: "# Hello World\n\nThis is just body content.",
			wantErr:  false,
		},
		{
			name:     "empty content",
			content:  "",
			wantMeta: map[string]interface{}{},
			wantBody: "",
			wantErr:  false,
		},
		{
			name: "frontmatter only",
			content: `---
title: Test
---`,
			wantMeta: map[string]interface{}{
				"title": "Test",
			},
			wantBody: "",
			wantErr:  false,
		},
		{
			name: "complex frontmatter",
			content: `---
id: 123
tags:
  - go
  - markdown
metadata:
  created: 2026-01-31
  status: active
---

Body content here.`,
			wantMeta: map[string]interface{}{
				"id": 123,
				"tags": []interface{}{
					"go",
					"markdown",
				},
				"metadata": map[interface{}]interface{}{
					"created": "2026-01-31",
					"status":  "active",
				},
			},
			wantBody: "\nBody content here.",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse([]byte(tt.content))

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, doc)
			assert.Equal(t, tt.wantBody, doc.Body)
			assert.Equal(t, tt.wantMeta, doc.Frontmatter)
		})
	}
}

func TestDocument_Write(t *testing.T) {
	tests := []struct {
		name    string
		doc     *Document
		want    string
		wantErr bool
	}{
		{
			name: "with frontmatter and body",
			doc: &Document{
				Frontmatter: map[string]interface{}{
					"title":  "Test",
					"author": "John",
				},
				Body: "\n# Content\n\nBody text.",
			},
			want: `---
author: John
title: Test
---

# Content

Body text.`,
			wantErr: false,
		},
		{
			name: "body only",
			doc: &Document{
				Frontmatter: map[string]interface{}{},
				Body:        "# Just Body\n\nNo frontmatter.",
			},
			want:    "# Just Body\n\nNo frontmatter.",
			wantErr: false,
		},
		{
			name: "empty document",
			doc: &Document{
				Frontmatter: map[string]interface{}{},
				Body:        "",
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "frontmatter only",
			doc: &Document{
				Frontmatter: map[string]interface{}{
					"status": "active",
				},
				Body: "",
			},
			want: `---
status: active
---
`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.doc.Write()

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestParseWrite_RoundTrip(t *testing.T) {
	original := `---
id: test-123
status: active
tags:
  - markdown
  - go
---

# Test Document

This is a test.`

	// Parse
	doc, err := Parse([]byte(original))
	require.NoError(t, err)
	require.NotNil(t, doc)

	// Write
	output, err := doc.Write()
	require.NoError(t, err)

	// Parse again
	doc2, err := Parse(output)
	require.NoError(t, err)
	require.NotNil(t, doc2)

	// Compare
	assert.Equal(t, doc.Body, doc2.Body)
	assert.Equal(t, doc.Frontmatter, doc2.Frontmatter)
}
