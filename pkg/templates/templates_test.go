package templates

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
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
func generateExpectedFiles(aiType string, promptFiles []string) []string {
	var wantFiles []string

	// Always include Mission Toolkit templates
	missionFiles := []string{".mission/governance.md", ".mission/backlog.md"}
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
			wantFiles := generateExpectedFiles(tt.aiType, availableTemplates)

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
			expectedPattern := tt.prefix + "m.apply"
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
			expectedDirs := []string{"displays"}
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

func TestWriteLibraryTemplatesShellFiles(t *testing.T) {
	fs := afero.NewMemMapFs()
	targetDir := "/test"

	err := WriteLibraryTemplates(fs, targetDir, "q")
	if err != nil {
		t.Fatalf("WriteLibraryTemplates() error = %v", err)
	}

	// Walk through the library directory to find .sh files
	libraryDir := filepath.Join(targetDir, ".mission", "libraries")

	var shFiles []string
	err = afero.Walk(fs, libraryDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".sh") {
			shFiles = append(shFiles, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Error walking library directory: %v", err)
	}

	// Verify .sh files have executable permissions
	for _, shFile := range shFiles {
		info, err := fs.Stat(shFile)
		if err != nil {
			t.Errorf("Error getting file info for %s: %v", shFile, err)
			continue
		}

		mode := info.Mode()
		if mode.Perm() != 0755 {
			t.Errorf("Shell file %s has permissions %o, expected 0755", shFile, mode.Perm())
		}
	}

	// Verify .md files still have regular permissions
	var mdFiles []string
	err = afero.Walk(fs, libraryDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".md") {
			mdFiles = append(mdFiles, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Error walking library directory for .md files: %v", err)
	}

	for _, mdFile := range mdFiles {
		info, err := fs.Stat(mdFile)
		if err != nil {
			t.Errorf("Error getting file info for %s: %v", mdFile, err)
			continue
		}

		mode := info.Mode()
		if mode.Perm() != 0644 {
			t.Errorf("Markdown file %s has permissions %o, expected 0644", mdFile, mode.Perm())
		}
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

	// First init
	err := WriteTemplates(fs, targetDir, "q")
	require.NoError(t, err)

	// Add user content to backlog.md
	backlogPath := filepath.Join(targetDir, ".mission", "backlog.md")
	userContent := "## FUTURE ENHANCEMENTS\n- My custom backlog item"
	err = afero.WriteFile(fs, backlogPath, []byte(userContent), 0644)
	require.NoError(t, err)

	// Second init
	err = WriteTemplates(fs, targetDir, "claude")
	require.NoError(t, err)

	// Verify user content is preserved in backlog.md
	finalBacklog, err := afero.ReadFile(fs, backlogPath)
	require.NoError(t, err)
	require.Equal(t, userContent, string(finalBacklog), "backlog.md should not be overwritten")

	// Verify other files were updated (e.g., governance.md)
	govPath := filepath.Join(targetDir, ".mission", "governance.md")
	govContent, err := afero.ReadFile(fs, govPath)
	require.NoError(t, err)
	require.Contains(t, string(govContent), "/m.", "governance.md should be updated with new prefix")
}
