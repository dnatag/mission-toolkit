// Package md provides utilities for parsing and writing markdown documents
// with YAML frontmatter and section-based content manipulation.
package md

import (
	"regexp"
	"strings"
)

// FindSection locates a section header in markdown body and returns its line index.
// Section names are case-insensitive and matched against "## SECTIONNAME" format.
// Returns -1 if section not found.
//
// Example:
//
//	body := "## INTENT\nContent\n## SCOPE\nFiles"
//	idx := FindSection(body, "scope") // returns 2
func FindSection(body, sectionName string) int {
	lines := strings.Split(body, "\n")
	header := "## " + strings.ToUpper(sectionName)

	for i, line := range lines {
		if strings.TrimSpace(line) == header {
			return i
		}
	}
	return -1
}

// ExtractSection extracts string content from a section until the next section or end.
// Section names are case-insensitive. Content is trimmed of leading/trailing whitespace.
// Returns empty string if section not found or section is empty.
//
// Example:
//
//	body := "## INTENT\nAdd feature\n## SCOPE\nFiles"
//	content := ExtractSection(body, "intent") // returns "Add feature"
func ExtractSection(body, sectionName string) string {
	pattern := "(?s)## " + regexp.QuoteMeta(strings.ToUpper(sectionName)) + "\\s*\\n(.*?)(?:(?:\\n##)|$)"
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(body)

	if len(matches) > 1 {
		content := strings.TrimSpace(matches[1])
		// Empty section followed immediately by another section
		if strings.HasPrefix(content, "##") {
			return ""
		}
		return content
	}
	return ""
}

// ExtractList parses list items from a section, supporting multiple markdown formats:
//   - Dash lists: "- item"
//   - Asterisk lists: "* item"
//   - Numbered lists: "1. item"
//   - Checkboxes: "- [ ] item" or "- [x] item"
//
// Returns empty slice if section not found or contains no list items.
// Empty lines within lists are ignored.
//
// Example:
//
//	body := "## SCOPE\n- file1.go\n- file2.go"
//	items := ExtractList(body, "scope") // returns ["file1.go", "file2.go"]
func ExtractList(body, sectionName string) []string {
	content := ExtractSection(body, sectionName)
	items := []string{}

	if content == "" {
		return items
	}

	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// Check checkbox patterns first (more specific than plain dash)
		if strings.HasPrefix(trimmed, "- [ ] ") {
			items = append(items, strings.TrimSpace(trimmed[6:]))
		} else if strings.HasPrefix(trimmed, "- [x] ") {
			items = append(items, strings.TrimSpace(trimmed[6:]))
		} else if strings.HasPrefix(trimmed, "- ") {
			items = append(items, strings.TrimSpace(trimmed[2:]))
		} else if strings.HasPrefix(trimmed, "* ") {
			items = append(items, strings.TrimSpace(trimmed[2:]))
		} else if matched, _ := regexp.MatchString(`^\d+\.\s`, trimmed); matched {
			// Numbered list: "1. item"
			parts := strings.SplitN(trimmed, ". ", 2)
			if len(parts) == 2 {
				items = append(items, strings.TrimSpace(parts[1]))
			}
		}
	}

	return items
}

// SkipToNextSection returns the line index of the next section header after startIndex.
// Useful for iterating through sections or skipping section content.
// Returns -1 if no next section found.
//
// Example:
//
//	body := "## INTENT\nContent\n## SCOPE\nFiles"
//	idx := FindSection(body, "intent")        // returns 0
//	next := SkipToNextSection(body, idx)      // returns 2
func SkipToNextSection(body string, startIndex int) int {
	lines := strings.Split(body, "\n")

	for i := startIndex + 1; i < len(lines); i++ {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "## ") {
			return i
		}
	}
	return -1
}
