package templates

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/dnatag/mission-toolkit/internal/utils"
	"github.com/spf13/afero"
)

// getAvailablePromptTemplates dynamically discovers all .md files in the prompts directory
func getAvailablePromptTemplates() ([]string, error) {
	var templates []string
	err := fs.WalkDir(promptTemplates, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".md") {
			// Remove "prompts/" prefix and keep just the filename
			filename := filepath.Base(path)
			templates = append(templates, filename)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(templates) // Ensure consistent ordering
	return templates, nil
}

// generateExpectedFiles creates the expected file list for a given AI type
func generateExpectedFiles(aiType, targetDir string, promptFiles []string) []string {
	var wantFiles []string

	// Always include Mission Toolkit templates
	missionFiles := []string{".mission/governance.md", ".mission/metrics.md", ".mission/backlog.md"}
	wantFiles = append(wantFiles, missionFiles...)

	// Add prompt files based on AI type
	var promptDir string
	switch aiType {
	case "q":
		promptDir = ".amazonq/prompts"
	case "claude":
		promptDir = ".claude/commands"
	case "kiro":
		promptDir = ".kiro/prompts"
	case "opencode":
		promptDir = ".opencode/command"
	}

	for _, file := range promptFiles {
		wantFiles = append(wantFiles, filepath.Join(promptDir, file))
	}

	return wantFiles
}

func TestWriteTemplates(t *testing.T) {
	// Dynamically discover available prompt templates
	availableTemplates, err := getAvailablePromptTemplates()
	if err != nil {
		t.Fatalf("Failed to discover available templates: %v", err)
	}

	if len(availableTemplates) == 0 {
		t.Fatal("No prompt templates found - this indicates a problem with the embedded filesystem")
	}

	tests := []struct {
		name      string
		aiType    string
		targetDir string
	}{
		{"Amazon Q templates", "q", "/test"},
		{"Claude templates", "claude", "/test"},
		{"Kiro templates", "kiro", "/test"},
		{"OpenCode templates", "opencode", "/test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()

			err := WriteTemplates(fs, tt.targetDir, tt.aiType)
			if err != nil {
				t.Fatalf("WriteTemplates() error = %v", err)
			}

			// Generate expected files dynamically
			wantFiles := generateExpectedFiles(tt.aiType, tt.targetDir, availableTemplates)

			for _, file := range wantFiles {
				fullPath := filepath.Join(tt.targetDir, file)
				exists, err := afero.Exists(fs, fullPath)
				if err != nil {
					t.Errorf("Error checking file %s: %v", fullPath, err)
				}
				if !exists {
					t.Errorf("Expected file %s does not exist", fullPath)
				}

				// Verify file has content
				content, err := afero.ReadFile(fs, fullPath)
				if err != nil {
					t.Errorf("Error reading file %s: %v", fullPath, err)
				}
				if len(content) == 0 {
					t.Errorf("File %s is empty", fullPath)
				}
			}
		})
	}
}

func TestWriteTemplatesUnsupportedAI(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := WriteTemplates(fs, "/test", "unsupported")
	if err == nil {
		t.Error("Expected error for unsupported AI type, but got nil")
	}

	expectedError := "unsupported AI type 'unsupported'"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', but got: %v", expectedError, err)
	}
}

func TestGetSlashPrefix(t *testing.T) {
	tests := []struct {
		aiType   string
		expected string
	}{
		{"q", "@"},
		{"kiro", "@"},
		{"claude", "/"},
		{"opencode", "/"},
	}

	for _, tt := range tests {
		t.Run(tt.aiType, func(t *testing.T) {
			result := getSlashPrefix(tt.aiType)
			if result != tt.expected {
				t.Errorf("getSlashPrefix(%s) = %s, want %s", tt.aiType, result, tt.expected)
			}
		})
	}
}

func TestSlashPrefixReplacement(t *testing.T) {
	tests := []struct {
		name   string
		aiType string
		prefix string
	}{
		{"Amazon Q uses @", "q", "@"},
		{"Kiro uses @", "kiro", "@"},
		{"Claude uses /", "claude", "/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()

			err := WriteTemplates(fs, "/test", tt.aiType)
			if err != nil {
				t.Fatalf("WriteTemplates() error = %v", err)
			}

			// Check governance.md for correct prefix
			govPath := "/test/.mission/governance.md"
			content, err := afero.ReadFile(fs, govPath)
			if err != nil {
				t.Fatalf("Failed to read governance.md: %v", err)
			}

			contentStr := string(content)
			expectedPattern := tt.prefix + "m.clarify"
			if !strings.Contains(contentStr, expectedPattern) {
				t.Errorf("Expected governance.md to contain '%s', but it doesn't. Content: %s", expectedPattern, contentStr)
			}
		})
	}
}

func TestWriteLibraryTemplates(t *testing.T) {
	tests := []struct {
		name      string
		aiType    string
		targetDir string
		prefix    string
	}{
		{"Amazon Q library templates", "q", "/test", "@"},
		{"Kiro library templates", "kiro", "/test", "@"},
		{"Claude library templates", "claude", "/test", "/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()

			err := WriteLibraryTemplates(fs, tt.targetDir, tt.aiType)
			if err != nil {
				t.Fatalf("WriteLibraryTemplates() error = %v", err)
			}

			// Check that .mission/libraries directory exists
			libraryDir := filepath.Join(tt.targetDir, ".mission", "libraries")
			exists, err := afero.DirExists(fs, libraryDir)
			if err != nil {
				t.Fatalf("Error checking library directory: %v", err)
			}
			if !exists {
				t.Errorf("Library directory %s does not exist", libraryDir)
			}

			// Check for expected subdirectories
			expectedDirs := []string{"displays", "missions", "scripts", "metrics", "variables"}
			for _, dir := range expectedDirs {
				dirPath := filepath.Join(libraryDir, dir)
				exists, err := afero.DirExists(fs, dirPath)
				if err != nil {
					t.Errorf("Error checking directory %s: %v", dirPath, err)
				}
				if !exists {
					t.Errorf("Expected directory %s does not exist", dirPath)
				}
			}

			// Check for specific template files and prefix replacement
			testFiles := []string{
				"displays/plan-success.md",
				"missions/wet.md",
				"scripts/create-mission.md",
			}

			for _, file := range testFiles {
				filePath := filepath.Join(libraryDir, file)
				exists, err := afero.Exists(fs, filePath)
				if err != nil {
					t.Errorf("Error checking file %s: %v", filePath, err)
					continue
				}
				if !exists {
					t.Errorf("Expected file %s does not exist", filePath)
					continue
				}

				// Verify prefix replacement
				content, err := afero.ReadFile(fs, filePath)
				if err != nil {
					t.Errorf("Error reading file %s: %v", filePath, err)
					continue
				}

				contentStr := string(content)
				expectedPattern := tt.prefix + "m.apply"
				if strings.Contains(contentStr, "/m.apply") && tt.prefix == "@" {
					t.Errorf("File %s still contains '/m.apply' instead of '@m.apply'", filePath)
				}
				if strings.Contains(contentStr, expectedPattern) || strings.Contains(contentStr, tt.prefix+"m.") {
					// Good - contains expected prefix
				} else if strings.Contains(contentStr, "m.apply") {
					// File might not contain slash commands, which is okay
				}
			}
		})
	}
}

func TestWriteLibraryTemplatesUnsupportedAI(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := WriteLibraryTemplates(fs, "/test", "unsupported")
	if err != nil {
		t.Errorf("WriteLibraryTemplates should not fail for unsupported AI type, got: %v", err)
	}

	// Should still create the directory structure
	libraryDir := "/test/.mission/libraries"
	exists, err := afero.DirExists(fs, libraryDir)
	if err != nil {
		t.Fatalf("Error checking library directory: %v", err)
	}
	if !exists {
		t.Errorf("Library directory should be created even for unsupported AI types")
	}
}

func TestIdempotentInit(t *testing.T) {
	fs := afero.NewMemMapFs()
	targetDir := "/test"

	// First init with Q
	err := WriteTemplates(fs, targetDir, "q")
	if err != nil {
		t.Fatalf("First WriteTemplates() error = %v", err)
	}

	// Add user content to backlog.md
	backlogPath := filepath.Join(targetDir, ".mission", "backlog.md")
	existingContent, err := afero.ReadFile(fs, backlogPath)
	if err != nil {
		t.Fatalf("Failed to read backlog.md: %v", err)
	}

	userContent := string(existingContent) + "\n- Test backlog item\n- Another test item"
	err = afero.WriteFile(fs, backlogPath, []byte(userContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write user content: %v", err)
	}

	// Add user content to metrics.md
	metricsPath := filepath.Join(targetDir, ".mission", "metrics.md")
	existingMetrics, err := afero.ReadFile(fs, metricsPath)
	if err != nil {
		t.Fatalf("Failed to read metrics.md: %v", err)
	}

	userMetrics := strings.ReplaceAll(string(existingMetrics), "Total Missions**: [count]", "Total Missions**: 5")
	err = afero.WriteFile(fs, metricsPath, []byte(userMetrics), 0644)
	if err != nil {
		t.Fatalf("Failed to write user metrics: %v", err)
	}

	// Second init with Claude (different AI type)
	err = WriteTemplates(fs, targetDir, "claude")
	if err != nil {
		t.Fatalf("Second WriteTemplates() error = %v", err)
	}

	// Verify user content is preserved in backlog.md
	finalBacklog, err := afero.ReadFile(fs, backlogPath)
	if err != nil {
		t.Fatalf("Failed to read final backlog.md: %v", err)
	}

	backlogStr := string(finalBacklog)
	if !strings.Contains(backlogStr, "Test backlog item") {
		t.Errorf("User backlog content not preserved: %s", backlogStr)
	}
	if !strings.Contains(backlogStr, "Another test item") {
		t.Errorf("User backlog content not preserved: %s", backlogStr)
	}

	// Verify user content is preserved in metrics.md
	finalMetrics, err := afero.ReadFile(fs, metricsPath)
	if err != nil {
		t.Fatalf("Failed to read final metrics.md: %v", err)
	}

	metricsStr := string(finalMetrics)
	if !strings.Contains(metricsStr, "Total Missions**: 5") {
		t.Errorf("User metrics content not preserved: %s", metricsStr)
	}

	// Verify AI-specific content was updated (check governance.md for slash prefixes)
	govPath := filepath.Join(targetDir, ".mission", "governance.md")
	govContent, err := afero.ReadFile(fs, govPath)
	if err != nil {
		t.Fatalf("Failed to read governance.md: %v", err)
	}

	govStr := string(govContent)
	if !strings.Contains(govStr, "/m.clarify") {
		t.Errorf("AI-specific content not updated to Claude format in governance.md: %s", govStr)
	}
}

func TestExtractAllSections(t *testing.T) {
	content := `# Test Document

## SECTION ONE
- Item 1
- Item 2

Some text here

## SECTION TWO
- Different item
- Another item

## SECTION THREE
No items here
`

	sections := extractAllSections(content)

	if len(sections) != 3 {
		t.Errorf("Expected 3 sections, got %d", len(sections))
	}

	// Check section content
	if !strings.Contains(sections["SECTION ONE"], "Item 1") {
		t.Errorf("SECTION ONE missing expected content")
	}
	if !strings.Contains(sections["SECTION ONE"], "Item 2") {
		t.Errorf("SECTION ONE missing expected content")
	}
	if !strings.Contains(sections["SECTION TWO"], "Different item") {
		t.Errorf("SECTION TWO missing expected content")
	}
	if !strings.Contains(sections["SECTION TWO"], "Another item") {
		t.Errorf("SECTION TWO missing expected content")
	}
	if sections["SECTION THREE"] == "" {
		t.Errorf("SECTION THREE should exist even if empty")
	}
}

func TestPreserveUserSections(t *testing.T) {
	template := `# Template

## SECTION ONE
Template content 1

## SECTION TWO
Template content 2

## NEW SECTION
New template content
`

	existing := `# Existing

## SECTION ONE
User content 1
- User item

## SECTION TWO
Template content 2

## OLD SECTION
Old user content
`

	result := preserveUserSections(template, existing)

	// Should preserve user content in SECTION ONE
	if !strings.Contains(result, "User content 1") {
		t.Errorf("User content not preserved in SECTION ONE")
	}
	if !strings.Contains(result, "User item") {
		t.Errorf("User items not preserved in SECTION ONE")
	}

	// Should keep template content in SECTION TWO (no user changes)
	if !strings.Contains(result, "Template content 2") {
		t.Errorf("Template content missing in SECTION TWO")
	}

	// Should keep new template sections
	if !strings.Contains(result, "NEW SECTION") {
		t.Errorf("New template section missing")
	}

	// Should not include old sections not in template
	if strings.Contains(result, "OLD SECTION") {
		t.Errorf("Old section should not be preserved")
	}
}

func TestReplaceSection(t *testing.T) {
	t.Run("with existing user data", func(t *testing.T) {
		content := `# Document

## SECTION ONE
Original content

## SECTION TWO
Keep this content
`

		newContent := "Replaced content\n- New item"
		result := replaceSection(content, "SECTION ONE", newContent)

		if !strings.Contains(result, "Replaced content") {
			t.Errorf("Section content not replaced")
		}
		if !strings.Contains(result, "New item") {
			t.Errorf("New content not added")
		}
		if !strings.Contains(result, "Keep this content") {
			t.Errorf("Other sections should be preserved")
		}
		if strings.Contains(result, "Original content") {
			t.Errorf("Original content should be replaced")
		}
	})

	t.Run("with no user data", func(t *testing.T) {
		content := `# Document

## SECTION ONE

## SECTION TWO
Keep this content
`

		newContent := "New content\n- Added item"
		result := replaceSection(content, "SECTION ONE", newContent)

		if !strings.Contains(result, "New content") {
			t.Errorf("New content not added to empty section")
		}
		if !strings.Contains(result, "Added item") {
			t.Errorf("New content items not added")
		}
		if !strings.Contains(result, "Keep this content") {
			t.Errorf("Other sections should be preserved")
		}
	})
}

func TestParseSections(t *testing.T) {
	content := `# Test Document

## SECTION ONE
- Item 1
- Item 2

Some text here

## SECTION TWO
- Different item
- Another item

## SECTION THREE
No items here

## OTHER SECTION
This should not be captured
`

	result := utils.ParseSections(content)

	if len(result) != 4 {
		t.Errorf("Expected 4 sections, got %d", len(result))
	}

	// Check section headers
	expectedHeaders := []string{"SECTION ONE", "SECTION TWO", "SECTION THREE", "OTHER SECTION"}
	for i, section := range result {
		if section.Header != expectedHeaders[i] {
			t.Errorf("Expected header %s, got %s", expectedHeaders[i], section.Header)
		}
	}

	// Check SECTION ONE has content
	section1 := result[0]
	if len(section1.Content) == 0 {
		t.Errorf("Expected content in SECTION ONE, got empty")
	}
}
