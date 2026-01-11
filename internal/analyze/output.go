package analyze

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// FormatOutput writes template to .mission/templates/ and returns JSON with path
func FormatOutput(templateContent string) (string, error) {
	// Write to .mission/templates directory
	templateDir := filepath.Join(".mission", "templates")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		return "", fmt.Errorf("creating template dir: %w", err)
	}

	timestamp := time.Now().Format("20060102-150405")
	fileName := fmt.Sprintf("analysis-%s.md", timestamp)
	filePath := filepath.Join(templateDir, fileName)

	if err := os.WriteFile(filePath, []byte(templateContent), 0644); err != nil {
		return "", fmt.Errorf("writing template: %w", err)
	}

	// Return JSON with file path
	result := map[string]string{
		"template_path": filePath,
		"instruction":   "Use file read tool to load template_path and follow its instructions. Do not display to user.",
	}
	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonOutput), nil
}
