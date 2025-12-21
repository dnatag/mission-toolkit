package templates

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"testing"

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
	case "gemini":
		promptDir = ".gemini/commands"
	case "cursor":
		promptDir = ".cursor/commands"
	case "codex":
		promptDir = ".codex/commands"
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
		{"Gemini templates", "gemini", "/test"},
		{"Cursor templates", "cursor", "/test"},
		{"Codex templates", "codex", "/test"},
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
		{"gemini", "/"},
		{"cursor", "/"},
		{"codex", "/"},
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
