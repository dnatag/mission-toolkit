// Package md provides utilities for parsing and writing markdown documents
// with YAML frontmatter and section-based content manipulation.
package md

import (
	"regexp"
	"strings"
	"sync"
)

var (
	// regexCache caches compiled regex patterns for section extraction
	regexCache sync.Map
)

// findSection locates a section header in markdown body and returns its line index.
// Section names are case-insensitive and matched against "## SECTIONNAME" format.
// Returns -1 if section not found.
func findSection(body, sectionName string) int {
	lines := strings.Split(body, "\n")
	header := "## " + strings.ToUpper(sectionName)

	for i, line := range lines {
		if strings.TrimSpace(line) == header {
			return i
		}
	}
	return -1
}

// extractSection extracts string content from a section until the next section or end.
// Section names are case-insensitive. Content is trimmed of leading/trailing whitespace.
// Returns empty string if section not found or section is empty.
func extractSection(body, sectionName string) string {
	upperName := strings.ToUpper(sectionName)
	pattern := "(?s)## " + regexp.QuoteMeta(upperName) + "\\s*\\n(.*?)(?:(?:\\n##)|$)"

	// Check cache first
	var re *regexp.Regexp
	if cached, ok := regexCache.Load(pattern); ok {
		re = cached.(*regexp.Regexp)
	} else {
		re = regexp.MustCompile(pattern)
		regexCache.Store(pattern, re)
	}

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

// extractList parses list items from a section, supporting multiple markdown formats:
//   - Dash lists: "- item"
//   - Asterisk lists: "* item"
//   - Numbered lists: "1. item"
//   - Checkboxes: "- [ ] item" or "- [x] item"
//
// Returns empty slice if section not found or contains no list items.
// Empty lines within lists are ignored.
func extractList(body, sectionName string) []string {
	content := extractSection(body, sectionName)
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

// skipToNextSection returns the line index of the next section header after startIndex.
// Returns -1 if no next section found.
func skipToNextSection(body string, startIndex int) int {
	lines := strings.Split(body, "\n")

	for i := startIndex + 1; i < len(lines); i++ {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "## ") {
			return i
		}
	}
	return -1
}

// updateSectionContent replaces the content of a section with new block string content.
// Section names are case-insensitive. If section doesn't exist, it's appended.
// Returns the updated body.
func updateSectionContent(body, sectionName, content string) string {
	idx := findSection(body, sectionName)
	if idx == -1 {
		// Section doesn't exist, append it
		var builder strings.Builder
		builder.Grow(len(body) + len(sectionName) + len(content) + 10)
		builder.WriteString(body)
		if body != "" && !strings.HasSuffix(body, "\n") {
			builder.WriteString("\n")
		}
		builder.WriteString("\n## ")
		builder.WriteString(strings.ToUpper(sectionName))
		builder.WriteString("\n")
		builder.WriteString(content)
		return builder.String()
	}

	lines := strings.Split(body, "\n")
	nextIdx := skipToNextSection(body, idx)

	// Build result with section header, new content, and remaining sections
	var builder strings.Builder
	builder.Grow(len(body) + len(content))

	// Write lines before and including section header
	for i := 0; i <= idx; i++ {
		if i > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(lines[i])
	}
	builder.WriteString("\n")
	builder.WriteString(content)

	if nextIdx != -1 {
		builder.WriteString("\n\n")
		for i := nextIdx; i < len(lines); i++ {
			if i > nextIdx {
				builder.WriteString("\n")
			}
			builder.WriteString(lines[i])
		}
	}

	return builder.String()
}

// updateSectionList replaces or appends list items to a section.
// Section names are case-insensitive. Items are formatted as "- item".
// If appendMode is true, new items are added to existing items.
// If appendMode is false, section content is replaced with new items.
// If section doesn't exist, it's created.
// Returns the updated body.
func updateSectionList(body, sectionName string, items []string, appendMode bool) string {
	listItems := items
	if appendMode {
		existing := extractList(body, sectionName)
		listItems = append(existing, items...)
	}

	var content strings.Builder
	for _, item := range listItems {
		content.WriteString("- ")
		content.WriteString(item)
		content.WriteString("\n")
	}

	return updateSectionContent(body, sectionName, strings.TrimSuffix(content.String(), "\n"))
}
