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
	
	// Always include IDD templates
	iddFiles := []string{".idd/governance.md", ".idd/metrics.md", ".idd/backlog.md"}
	wantFiles = append(wantFiles, iddFiles...)
	
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
	case "cline":
		promptDir = ".clinerules/workflows"
	case "kiro":
		promptDir = ".kiro/prompts"
	default:
		promptDir = "prompts"
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
		{"Cline templates", "cline", "/test"},
		{"Kiro templates", "kiro", "/test"},
		{"Default templates", "default", "/test"},
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
