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

//go:embed libraries/**/*.md libraries/**/*.sh
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

	missionFiles := []string{"governance.md", "metrics.md", "backlog.md"}
	preserveSet := map[string]bool{"metrics.md": true, "backlog.md": true}

	for _, file := range missionFiles {
		content, err := missionTemplates.ReadFile("mission/" + file)
		if err != nil {
			return err
		}

		contentStr := strings.ReplaceAll(string(content), "/m.", prefix+"m.")
		filePath := filepath.Join(missionDir, file)

		if preserveSet[file] {
			if existingContent, err := afero.ReadFile(fs, filePath); err == nil {
				contentStr = preserveUserSections(contentStr, string(existingContent))
			}
		}

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

	promptFiles := []string{"m.complete.md", "m.plan.md", "m.apply.md", "m.clarify.md"}
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

// preserveUserSections preserves all user content sections generically
func preserveUserSections(templateContent, existingContent string) string {
	templateSections := extractAllSections(templateContent)
	existingSections := extractAllSections(existingContent)

	result := templateContent
	for heading, existingSection := range existingSections {
		if _, exists := templateSections[heading]; exists && strings.TrimSpace(existingSection) != "" {
			result = replaceSection(result, heading, existingSection)
		}
	}

	return result
}

// extractAllSections extracts all sections from markdown content
func extractAllSections(content string) map[string]string {
	sections := make(map[string]string)
	lines := strings.Split(content, "\n")

	var currentHeading string
	var start int

	for i, line := range lines {
		if trimmed := strings.TrimSpace(line); strings.HasPrefix(trimmed, "## ") {
			if currentHeading != "" {
				sections[currentHeading] = strings.Join(lines[start:i], "\n")
			}
			currentHeading = strings.TrimPrefix(trimmed, "## ")
			start = i + 1
		}
	}

	if currentHeading != "" {
		sections[currentHeading] = strings.Join(lines[start:], "\n")
	}

	return sections
}

// replaceSection replaces content under a specific heading
func replaceSection(content, heading, newContent string) string {
	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines)+10)
	inSection := false

	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "## ") {
			if inSection {
				result = append(result, strings.Split(newContent, "\n")...)
				inSection = false
			}
			if strings.Contains(line, heading) {
				result = append(result, line)
				inSection = true
				continue
			}
		}

		if !inSection {
			result = append(result, line)
		}
	}

	if inSection {
		result = append(result, strings.Split(newContent, "\n")...)
	}

	return strings.Join(result, "\n")
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
