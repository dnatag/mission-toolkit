package templates

import (
	"embed"
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

//go:embed mission/*.md
var missionTemplates embed.FS

//go:embed prompts/*.md
var promptTemplates embed.FS

// SupportedAITypes lists all supported AI types
var SupportedAITypes = []string{"q", "claude", "gemini", "cursor", "codex", "cline", "kiro", "opencode"}

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
	// Write Mission Toolkit templates to .mission directory
	missionDir := filepath.Join(targetDir, ".mission")
	if err := fs.MkdirAll(missionDir, 0755); err != nil {
		return err
	}

	missionFiles := []string{"governance.md", "metrics.md", "backlog.md"}
	for _, file := range missionFiles {
		content, err := missionTemplates.ReadFile("mission/" + file)
		if err != nil {
			return err
		}
		if err := afero.WriteFile(fs, filepath.Join(missionDir, file), content, 0644); err != nil {
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
	case "opencode":
		promptDir = filepath.Join(targetDir, ".opencode", "command")
	default:
		return fmt.Errorf("unsupported AI type '%s'. Supported types: %v", aiType, SupportedAITypes)
	}

	if err := fs.MkdirAll(promptDir, 0755); err != nil {
		return err
	}

	promptFiles := []string{"m.complete.md", "m.plan.md", "m.apply.md", "m.clarify.md"}
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
