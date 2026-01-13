package analyze

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/dnatag/mission-toolkit/internal/logger"
	"github.com/dnatag/mission-toolkit/internal/mission"
	"github.com/spf13/afero"
)

// CreateTestLoggerConfig creates a console-only logger config for tests
func CreateTestLoggerConfig(fs afero.Fs) *logger.Config {
	return &logger.Config{
		Level:    logger.DefaultConfig().Level,
		Format:   logger.DefaultConfig().Format,
		Output:   logger.OutputConsole,
		FilePath: "/dev/null", // Use /dev/null to prevent file creation
		Fs:       fs,          // Use same in-memory filesystem
	}
}

// CreateLogger creates a logger with optional config, reading mission ID from filesystem
func CreateLogger(fs afero.Fs, loggerConfig *logger.Config) *logger.Logger {
	reader := mission.NewReader(fs)
	missionID, _ := reader.GetMissionID(filepath.Join(".mission", "mission.md"))

	if loggerConfig != nil {
		return logger.NewWithConfig(missionID, loggerConfig)
	}
	return logger.New(missionID)
}

// FormatOutput writes template to .mission/templates/ and returns JSON with path
func FormatOutput(templateContent string) (string, error) {
	return FormatOutputWithFS(afero.NewOsFs(), templateContent)
}

// FormatOutputWithFS writes template using provided filesystem
func FormatOutputWithFS(fs afero.Fs, templateContent string) (string, error) {
	// Write to .mission/templates directory
	templateDir := filepath.Join(".mission", "templates")
	if err := fs.MkdirAll(templateDir, 0755); err != nil {
		return "", fmt.Errorf("creating template dir: %w", err)
	}

	timestamp := time.Now().Format("20060102-150405")
	fileName := fmt.Sprintf("analysis-%s.md", timestamp)
	filePath := filepath.Join(templateDir, fileName)

	if err := afero.WriteFile(fs, filePath, []byte(templateContent), 0644); err != nil {
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
