// Package md provides utilities for parsing and writing markdown documents
// with YAML frontmatter. It uses the adrg/frontmatter library for robust
// frontmatter extraction and gopkg.in/yaml.v3 for YAML serialization.
package md

import (
	"bytes"
	"fmt"
	"strings"

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

// GetSection retrieves the content of a section by name.
// Section names are case-insensitive.
// Returns empty string and nil error if section doesn't exist.
func (d *Document) GetSection(name string) (string, error) {
	if err := validateSectionName(name); err != nil {
		return "", err
	}
	return extractSection(d.Body, name), nil
}

// GetList retrieves list items from a section by name.
// Section names are case-insensitive.
// Returns empty slice and nil error if section doesn't exist or contains no lists.
func (d *Document) GetList(name string) ([]string, error) {
	if err := validateSectionName(name); err != nil {
		return nil, err
	}
	return extractList(d.Body, name), nil
}

// UpdateSectionContent replaces the content of a section.
// Section names are case-insensitive.
// If section doesn't exist, it will be created.
func (d *Document) UpdateSectionContent(name, content string) error {
	if err := validateSectionName(name); err != nil {
		return err
	}
	d.Body = updateSectionContent(d.Body, name, content)
	return nil
}

// UpdateSectionList replaces list items in a section.
// Section names are case-insensitive. Items are formatted as "- item".
// If section doesn't exist, it will be created.
func (d *Document) UpdateSectionList(name string, items []string) error {
	if err := validateSectionName(name); err != nil {
		return err
	}
	d.Body = updateSectionList(d.Body, name, items, false)
	return nil
}

// AppendSectionList appends list items to a section.
// Section names are case-insensitive. Items are formatted as "- item".
// If section doesn't exist, it will be created with the items.
func (d *Document) AppendSectionList(name string, items []string) error {
	if err := validateSectionName(name); err != nil {
		return err
	}
	d.Body = updateSectionList(d.Body, name, items, true)
	return nil
}

// HasSection checks if a section exists in the document.
// Section names are case-insensitive.
func (d *Document) HasSection(name string) bool {
	return findSection(d.Body, name) != -1
}

// ListSections returns all section names found in the document.
// Section names are returned in the order they appear.
func (d *Document) ListSections() []string {
	sections := []string{}
	lines := strings.Split(d.Body, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			section := strings.TrimPrefix(trimmed, "## ")
			sections = append(sections, section)
		}
	}
	return sections
}

// validateSectionName checks if a section name is valid.
// Returns error if name is empty, contains invalid characters,
// or could cause issues in markdown parsing.
func validateSectionName(name string) error {
	if name == "" {
		return fmt.Errorf("section name cannot be empty")
	}

	// Check for characters that could break markdown or cause issues
	if strings.ContainsAny(name, "\n\r\t") {
		return fmt.Errorf("section name contains invalid characters")
	}

	// Limit length to prevent abuse
	if len(name) > 100 {
		return fmt.Errorf("section name too long (max 100 characters)")
	}

	return nil
}
