package templates

import (
	"embed"
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

//go:embed idd/*.md
var iddTemplates embed.FS

//go:embed prompts/*.md
var promptTemplates embed.FS

// SupportedAITypes lists all supported AI types
var SupportedAITypes = []string{"q", "claude", "gemini", "cursor", "codex", "cline", "kiro"}

// ValidateAIType checks if the provided AI type is supported
func ValidateAIType(aiType string) error {
	for _, supported := range SupportedAITypes {
		if aiType == supported {
			return nil
		}
	}
	return fmt.Errorf("unsupported AI type '%s'. Supported types: %v", aiType, SupportedAITypes)
}

// WriteTemplates writes embedded templates to the specified filesystem
func WriteTemplates(fs afero.Fs, targetDir string, aiType string) error {
	// Write IDD templates to .idd directory
	iddDir := filepath.Join(targetDir, ".idd")
	if err := fs.MkdirAll(iddDir, 0755); err != nil {
		return err
	}

	iddFiles := []string{"governance.md", "metrics.md", "backlog.md"}
	for _, file := range iddFiles {
		content, err := iddTemplates.ReadFile("idd/" + file)
		if err != nil {
			return err
		}
		if err := afero.WriteFile(fs, filepath.Join(iddDir, file), content, 0644); err != nil {
			return err
		}
	}

	// Write prompt templates based on AI type
	var promptDir string
	switch aiType {
	case "q":
		promptDir = filepath.Join(targetDir, ".amazonq", "prompts")
	case "claude":
		promptDir = filepath.Join(targetDir, ".claude", "commands")
	case "gemini":
		promptDir = filepath.Join(targetDir, ".gemini", "commands")
	case "cursor":
		promptDir = filepath.Join(targetDir, ".cursor", "commands")
	case "codex":
		promptDir = filepath.Join(targetDir, ".codex", "commands")
	case "cline":
		promptDir = filepath.Join(targetDir, ".clinerules", "workflows")
	case "kiro":
		promptDir = filepath.Join(targetDir, ".kiro", "prompts")
	default:
		promptDir = filepath.Join(targetDir, "prompts")
	}

	if err := fs.MkdirAll(promptDir, 0755); err != nil {
		return err
	}

	promptFiles := []string{"idd.complete.md", "idd.plan.md", "idd.apply.md", "idd.clarify.md"}
	for _, file := range promptFiles {
		content, err := promptTemplates.ReadFile("prompts/" + file)
		if err != nil {
			return err
		}
		if err := afero.WriteFile(fs, filepath.Join(promptDir, file), content, 0644); err != nil {
			return err
		}
	}

	return nil
}
