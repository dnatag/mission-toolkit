package templates

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

//go:embed mission/*.md
var missionTemplates embed.FS

//go:embed prompts/*.md
var promptTemplates embed.FS

//go:embed libraries/**/*.md
var libraryTemplates embed.FS

// SupportedAITypes lists all supported AI types
var SupportedAITypes = []string{"q", "claude", "kiro", "opencode"}

// getSlashPrefix returns the appropriate slash command prefix for the AI type
func getSlashPrefix(aiType string) string {
	switch aiType {
	case "q", "kiro":
		return "@"
	default:
		return "/"
	}
}

// ValidateAIType checks if the provided AI type is supported
func ValidateAIType(aiType string) error {
	for _, supported := range SupportedAITypes {
		if supported == aiType {
			return nil
		}
	}
	return fmt.Errorf("unsupported AI type '%s'. Supported types: %v", aiType, SupportedAITypes)
}

// WriteTemplates writes embedded templates to the specified filesystem
func WriteTemplates(fs afero.Fs, targetDir string, aiType string) error {
	// Validate AI type first
	if err := ValidateAIType(aiType); err != nil {
		return err
	}

	prefix := getSlashPrefix(aiType)

	// Write Mission Toolkit templates to .mission directory
	missionDir := filepath.Join(targetDir, ".mission")
	if err := fs.MkdirAll(missionDir, 0755); err != nil {
		return err
	}

	missionFiles := []string{"governance.md", "backlog.md"}

	for _, file := range missionFiles {
		filePath := filepath.Join(missionDir, file)

		// If the file is backlog.md and it already exists, skip it.
		if file == "backlog.md" {
			if exists, _ := afero.Exists(fs, filePath); exists {
				continue
			}
		}

		content, err := missionTemplates.ReadFile("mission/" + file)
		if err != nil {
			return err
		}

		contentStr := strings.ReplaceAll(string(content), "/m.", prefix+"m.")

		if err := afero.WriteFile(fs, filePath, []byte(contentStr), 0644); err != nil {
			return err
		}
	}

	// AI-specific directory mapping
	aiDirs := map[string]string{
		"q":        ".amazonq/prompts",
		"claude":   ".claude/commands",
		"kiro":     ".kiro/prompts",
		"opencode": ".opencode/command",
	}

	promptDir := filepath.Join(targetDir, aiDirs[aiType])

	if err := fs.MkdirAll(promptDir, 0755); err != nil {
		return err
	}

	promptFiles := []string{"m.complete.md", "m.plan.md", "m.apply.md"}
	for _, file := range promptFiles {
		content, err := promptTemplates.ReadFile("prompts/" + file)
		if err != nil {
			return err
		}

		contentStr := strings.ReplaceAll(string(content), "/m.", prefix+"m.")
		if err := afero.WriteFile(fs, filepath.Join(promptDir, file), []byte(contentStr), 0644); err != nil {
			return err
		}
	}

	return nil
}

// WriteLibraryTemplates writes embedded library templates to .mission/libraries
func WriteLibraryTemplates(fs afero.Fs, targetDir string, aiType string) error {
	prefix := getSlashPrefix(aiType)

	libraryDir := filepath.Join(targetDir, ".mission", "libraries")
	if err := fs.MkdirAll(libraryDir, 0755); err != nil {
		return err
	}

	// Walk through all library template files using embed.FS
	return walkEmbedFS(libraryTemplates, "libraries", func(path string, isDir bool) error {
		if isDir || path == "libraries" {
			return nil
		}

		content, err := libraryTemplates.ReadFile(path)
		if err != nil {
			return err
		}

		relPath := strings.TrimPrefix(path, "libraries/")
		targetPath := filepath.Join(libraryDir, relPath)

		if err := fs.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		contentStr := strings.ReplaceAll(string(content), "/m.", prefix+"m.")

		// Set executable permissions for .sh files
		fileMode := os.FileMode(0644)
		if strings.HasSuffix(relPath, ".sh") {
			fileMode = 0755
		}

		return afero.WriteFile(fs, targetPath, []byte(contentStr), fileMode)
	})
}

// walkEmbedFS walks through an embedded filesystem
func walkEmbedFS(fsys embed.FS, root string, fn func(path string, isDir bool) error) error {
	entries, err := fsys.ReadDir(root)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(root, entry.Name())
		if err := fn(path, entry.IsDir()); err != nil {
			return err
		}

		if entry.IsDir() {
			if err := walkEmbedFS(fsys, path, fn); err != nil {
				return err
			}
		}
	}

	return nil
}
