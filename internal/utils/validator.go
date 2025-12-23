package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// ValidationReport contains the results of template validation
type ValidationReport struct {
	FilePath        string
	ParsedSections  []Section
	UnparsedContent []UnparsedElement
	IsValid         bool
	Errors          []string
}

// UnparsedElement represents content that wasn't parsed into Sections
type UnparsedElement struct {
	Type    string // "title", "paragraph", "other"
	Content string
	Line    int
}

// AllowedUnparsedContent defines content that's expected to be unparsed
var AllowedUnparsedContent = map[string]map[string]bool{
	"backlog.md": {
		"# MISSION TOOLKIT BACKLOG":       true,
		"---":                             true,
		"**Format for completed items:**": true,
		"**Format for refactoring opportunities:**": true,
		"**Format for decomposed intents:**":        true,
	},
	"metrics.md": {
		"# MISSION TOOLKIT METRICS SUMMARY":                                           true,
		"Detailed metrics with change summaries stored in completed/ with timestamps": true,
	},
}

// ValidateTemplate validates a template file using single-pass parsing
func ValidateTemplate(filePath string) (*ValidationReport, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	markdown := string(content)
	source := []byte(markdown)
	doc := goldmark.New().Parser().Parse(text.NewReader(source))

	// Single-pass extraction of both sections and unparsed content
	sections, unparsedElements := extractSectionsAndUnparsed(doc, source)
	templateType := getTemplateType(filePath)

	// Validate against allowlist
	report := &ValidationReport{
		FilePath:        filePath,
		ParsedSections:  sections,
		UnparsedContent: unparsedElements,
		IsValid:         true,
		Errors:          make([]string, 0, len(unparsedElements)),
	}

	allowedContent := AllowedUnparsedContent[templateType]
	for _, element := range unparsedElements {
		if !isContentAllowed(element.Content, allowedContent) {
			report.IsValid = false
			report.Errors = append(report.Errors,
				fmt.Sprintf("Unexpected unparsed content at line %d: %s", element.Line, element.Content))
		}
	}

	return report, nil
}

// extractSectionsAndUnparsed performs single-pass extraction of sections and unparsed content
func extractSectionsAndUnparsed(doc ast.Node, source []byte) ([]Section, []UnparsedElement) {
	var sections []Section
	var unparsed []UnparsedElement
	var current *Section

	lineNum := 1

	for child := doc.FirstChild(); child != nil; child = child.NextSibling() {
		switch n := child.(type) {
		case *ast.Heading:
			if n.Level == 2 {
				if current != nil {
					sections = append(sections, *current)
				}
				current = &Section{
					Header:  extractText(n, source),
					Content: make([]interface{}, 0, 4),
				}
			} else if n.Level == 1 {
				unparsed = append(unparsed, UnparsedElement{
					Type:    "title",
					Content: extractText(n, source),
					Line:    lineNum,
				})
			}

		case *ast.Paragraph:
			content := extractText(n, source)
			if current != nil && len(content) > 1 && content[0] == '(' && content[len(content)-1] == ')' {
				current.Content = append(current.Content, content)
			} else {
				unparsed = append(unparsed, UnparsedElement{
					Type:    "paragraph",
					Content: content,
					Line:    lineNum,
				})
			}

		case *ast.List:
			if current != nil {
				for listChild := n.FirstChild(); listChild != nil; listChild = listChild.NextSibling() {
					if listItem, ok := listChild.(*ast.ListItem); ok {
						text := extractText(listItem, source)
						if colonIndex := strings.Index(text, ":"); colonIndex != -1 {
							current.Content = append(current.Content, KeyValue{
								Key:   strings.TrimSpace(text[:colonIndex]),
								Value: strings.TrimSpace(text[colonIndex+1:]),
							})
						} else {
							current.Content = append(current.Content, text)
						}
					}
				}
			}

		case *ast.ThematicBreak:
			unparsed = append(unparsed, UnparsedElement{
				Type:    "other",
				Content: "---",
				Line:    lineNum,
			})
		}
		lineNum++
	}

	if current != nil {
		sections = append(sections, *current)
	}

	return sections, unparsed
}

// isContentAllowed checks if unparsed content is in the allowlist
func isContentAllowed(content string, allowedContent map[string]bool) bool {
	content = strings.TrimSpace(content)
	if allowedContent[content] {
		return true
	}
	// Fallback to substring matching for partial matches
	for allowed := range allowedContent {
		if strings.Contains(allowed, content) || strings.Contains(content, allowed) {
			return true
		}
	}
	return false
}

// getTemplateType extracts template type from file path
func getTemplateType(filePath string) string {
	switch {
	case strings.Contains(filePath, "backlog.md"):
		return "backlog.md"
	case strings.Contains(filePath, "metrics.md"):
		return "metrics.md"
	default:
		return "unknown"
	}
}
