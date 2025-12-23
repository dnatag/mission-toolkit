package utils

import (
	"os"
	"strings"
	"testing"
)

func TestValidateTemplate(t *testing.T) {
	tests := []struct {
		name           string
		templatePath   string
		expectValid    bool
		expectSections int
		expectUnparsed int
	}{
		{
			name:           "validate backlog template",
			templatePath:   "../../internal/templates/mission/backlog.md",
			expectValid:    true,
			expectSections: 3, // DECOMPOSED INTENTS, REFACTORING OPPORTUNITIES, FUTURE ENHANCEMENTS
			expectUnparsed: 5, // title + format instructions
		},
		{
			name:           "validate metrics template",
			templatePath:   "../../internal/templates/mission/metrics.md",
			expectValid:    true,
			expectSections: 7, // All ## sections including TECHNICAL LEARNINGS
			expectUnparsed: 2, // title + description
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report, err := ValidateTemplate(tt.templatePath)
			if err != nil {
				t.Fatalf("ValidateTemplate() error = %v", err)
			}

			if report.IsValid != tt.expectValid {
				t.Errorf("Expected IsValid = %v, got %v", tt.expectValid, report.IsValid)
				for _, err := range report.Errors {
					t.Logf("Validation error: %s", err)
				}
			}

			if len(report.ParsedSections) != tt.expectSections {
				t.Errorf("Expected %d section, got %d", tt.expectSections, len(report.ParsedSections))
				for i, section := range report.ParsedSections {
					t.Logf("Section %d: %s", i, section.Header)
				}
			}

			// Log unparsed content for debugging
			t.Logf("Unparsed content (%d items):", len(report.UnparsedContent))
			for _, element := range report.UnparsedContent {
				t.Logf("  %s (line %d): %s", element.Type, element.Line, element.Content)
			}
		})
	}
}

func TestValidateTemplateWithInvalidContent(t *testing.T) {
	// Create a temporary invalid template content
	invalidTemplate := `# MISSION TOOLKIT BACKLOG

## VALID SECTION
(This is valid)
- valid: item

This paragraph should not be here and will cause validation failure

## ANOTHER SECTION
- another item`

	// Write to temp file for testing
	tempFile := "/tmp/test_invalid_template.md"
	err := writeStringToFile(tempFile, invalidTemplate)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	report, err := ValidateTemplate(tempFile)
	if err != nil {
		t.Fatalf("ValidateTemplate() error = %v", err)
	}

	// Should be invalid due to unexpected paragraph
	if report.IsValid {
		t.Error("Expected template to be invalid, but it was marked as valid")
	}

	if len(report.Errors) == 0 {
		t.Error("Expected validation errors, but got none")
	}

	// Should still parse the valid Sections
	if len(report.ParsedSections) != 2 {
		t.Errorf("Expected 2 section, got %d", len(report.ParsedSections))
	}
}

func TestAllowedUnparsedContent(t *testing.T) {
	tests := []struct {
		content    string
		allowedMap map[string]bool
		expected   bool
	}{
		{
			content:    "# MISSION TOOLKIT BACKLOG",
			allowedMap: AllowedUnparsedContent["backlog.md"],
			expected:   true,
		},
		{
			content:    "Some random paragraph",
			allowedMap: AllowedUnparsedContent["backlog.md"],
			expected:   false,
		},
		{
			content:    "# MISSION TOOLKIT METRICS SUMMARY",
			allowedMap: AllowedUnparsedContent["metrics.md"],
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.content, func(t *testing.T) {
			result := isContentAllowed(tt.content, tt.allowedMap)
			if result != tt.expected {
				t.Errorf("isContentAllowed(%q) = %v, want %v", tt.content, result, tt.expected)
			}
		})
	}
}

func TestDirectChildIteration(t *testing.T) {
	// Test that direct child iteration works more reliably than ast.Walk
	templateContent := `# MISSION TOOLKIT TEST

## SECTION ONE
- item: value

Some paragraph content

## SECTION TWO  
- another: item

---

Format instructions after separator`

	tempFile := "/tmp/test_direct_iteration.md"
	err := writeStringToFile(tempFile, templateContent)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	report, err := ValidateTemplate(tempFile)
	if err != nil {
		t.Fatalf("ValidateTemplate() error = %v", err)
	}

	// Should parse 2 section (## sections)
	if len(report.ParsedSections) != 2 {
		t.Errorf("Expected 2 section, got %d", len(report.ParsedSections))
	}

	// Should identify unparsed elements including format instructions after separator
	foundTitle := false
	foundParagraph := false
	foundSeparator := false

	for _, element := range report.UnparsedContent {
		if element.Type == "title" && strings.Contains(element.Content, "MISSION TOOLKIT TEST") {
			foundTitle = true
		}
		if element.Type == "paragraph" && strings.Contains(element.Content, "Some paragraph content") {
			foundParagraph = true
		}
		if element.Type == "other" && element.Content == "---" {
			foundSeparator = true
		}
	}

	if !foundTitle {
		t.Error("Expected to find title in unparsed content")
	}
	if !foundParagraph {
		t.Error("Expected to find paragraph in unparsed content")
	}
	if !foundSeparator {
		t.Error("Expected to find separator in unparsed content")
	}
}

func TestFlexibleAllowlistApproach(t *testing.T) {
	// Test that allowlist approach provides flexibility for legitimate template elements
	tests := []struct {
		name        string
		content     string
		allowedMap  map[string]bool
		expected    bool
		description string
	}{
		{
			name:        "exact match",
			content:     "# MISSION TOOLKIT BACKLOG",
			allowedMap:  map[string]bool{"# MISSION TOOLKIT BACKLOG": true},
			expected:    true,
			description: "Exact string match should work",
		},
		{
			name:        "partial match in allowed",
			content:     "BACKLOG",
			allowedMap:  map[string]bool{"# MISSION TOOLKIT BACKLOG": true},
			expected:    true,
			description: "Content contained in allowed should work",
		},
		{
			name:        "partial match in content",
			content:     "# MISSION TOOLKIT BACKLOG EXTENDED",
			allowedMap:  map[string]bool{"# MISSION TOOLKIT BACKLOG": true},
			expected:    true,
			description: "Allowed contained in content should work",
		},
		{
			name:        "no match",
			content:     "Random content",
			allowedMap:  map[string]bool{"# MISSION TOOLKIT BACKLOG": true},
			expected:    false,
			description: "Unrelated content should not match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isContentAllowed(tt.content, tt.allowedMap)
			if result != tt.expected {
				t.Errorf("%s: isContentAllowed(%q) = %v, want %v",
					tt.description, tt.content, result, tt.expected)
			}
		})
	}
}

// Helper function to write string to file for testing
func writeStringToFile(filePath, content string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}

func TestGetTemplateType(t *testing.T) {
	tests := []struct {
		filePath string
		expected string
	}{
		{"internal/templates/mission/backlog.md", "backlog.md"},
		{"internal/templates/mission/metrics.md", "metrics.md"},
		{"/some/other/path/file.md", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.filePath, func(t *testing.T) {
			result := getTemplateType(tt.filePath)
			if result != tt.expected {
				t.Errorf("getTemplateType(%q) = %q, want %q", tt.filePath, result, tt.expected)
			}
		})
	}
}

func TestTemplateValidationBalance(t *testing.T) {
	// Test balancing strict parsing with expected unparsed content
	strictTemplate := `## VALID SECTION
- item: value`

	balancedTemplate := `# TITLE (expected unparsed)

## VALID SECTION  
- item: value

*Description text (expected unparsed)*

---

Format instructions (expected unparsed)`

	tests := []struct {
		name     string
		content  string
		expected bool
		reason   string
	}{
		{
			name:     "strict parsing only",
			content:  strictTemplate,
			expected: true,
			reason:   "Only parsed content should be valid",
		},
		{
			name:     "balanced with expected unparsed",
			content:  balancedTemplate,
			expected: false, // Will be invalid without proper allowlist
			reason:   "Mixed content requires allowlist for unparsed elements",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempFile := "/tmp/test_balance_" + tt.name + ".md"
			err := writeStringToFile(tempFile, tt.content)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}

			report, err := ValidateTemplate(tempFile)
			if err != nil {
				t.Fatalf("ValidateTemplate() error = %v", err)
			}

			if report.IsValid != tt.expected {
				t.Logf("Reason: %s", tt.reason)
				t.Errorf("Expected IsValid = %v, got %v", tt.expected, report.IsValid)
				for _, validationErr := range report.Errors {
					t.Logf("Validation error: %s", validationErr)
				}
			}
		})
	}
}
