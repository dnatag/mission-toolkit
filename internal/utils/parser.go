package utils

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// ParseSections parses Markdown Sections
func ParseSections(markdown string) []Section {
	source := []byte(markdown)
	doc := goldmark.New().Parser().Parse(text.NewReader(source))

	var sections []Section
	var current *Section

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
			}

		case *ast.Paragraph:
			if current != nil {
				if text := extractText(n, source); len(text) > 1 && text[0] == '(' && text[len(text)-1] == ')' {
					current.Content = append(current.Content, text)
				}
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
		}
	}

	if current != nil {
		sections = append(sections, *current)
	}

	return sections
}

// extractText extracts text content from an AST node recursively
func extractText(node ast.Node, source []byte) string {
	var buf bytes.Buffer

	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if text, ok := n.(*ast.Text); ok {
				buf.Write(text.Segment.Value(source))
			}
		}
		return ast.WalkContinue, nil
	})

	return strings.TrimSpace(buf.String())
}
