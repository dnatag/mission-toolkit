// Package md provides utilities for parsing and writing markdown documents
// with YAML frontmatter. It uses the adrg/frontmatter library for robust
// frontmatter extraction and gopkg.in/yaml.v3 for YAML serialization.
package md

import (
	"bytes"
	"fmt"

	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v3"
)

// Document represents a markdown document with YAML frontmatter and body content.
type Document struct {
	// Frontmatter contains the YAML metadata
	Frontmatter map[string]interface{}
	// Body contains the markdown content after frontmatter
	Body string
}

// Parse parses markdown content with YAML frontmatter into a Document.
// It handles documents with or without frontmatter gracefully.
//
// The function uses adrg/frontmatter to extract YAML metadata from the
// beginning of the content. If no frontmatter is present, an empty map
// is returned for Frontmatter field.
//
// Returns an error if the frontmatter YAML is malformed.
func Parse(content []byte) (*Document, error) {
	if len(content) == 0 {
		return &Document{
			Frontmatter: make(map[string]interface{}),
			Body:        "",
		}, nil
	}

	var meta map[string]interface{}

	rest, err := frontmatter.Parse(bytes.NewReader(content), &meta)
	if err != nil {
		return nil, fmt.Errorf("parsing frontmatter: %w", err)
	}

	// Initialize empty map if nil
	if meta == nil {
		meta = make(map[string]interface{})
	}

	return &Document{
		Frontmatter: meta,
		Body:        string(rest),
	}, nil
}

// Write serializes the Document back to markdown format with YAML frontmatter.
//
// The output format is:
//
//	---
//	<YAML frontmatter>
//	---
//	<body content>
//
// If Frontmatter is empty, only the body is written (no frontmatter delimiters).
// Returns the complete markdown content as bytes, or an error if YAML
// marshaling fails.
func (d *Document) Write() ([]byte, error) {
	var buf bytes.Buffer

	// Write frontmatter only if present
	if len(d.Frontmatter) > 0 {
		buf.WriteString("---\n")

		yamlData, err := yaml.Marshal(d.Frontmatter)
		if err != nil {
			return nil, fmt.Errorf("marshaling frontmatter: %w", err)
		}

		buf.Write(yamlData)
		buf.WriteString("---\n")
	}

	// Write body
	buf.WriteString(d.Body)

	return buf.Bytes(), nil
}
