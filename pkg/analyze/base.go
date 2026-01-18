package analyze

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/dnatag/mission-toolkit/pkg/logger"
	"github.com/spf13/afero"
)

// BaseService provides the foundation for analysis services.
// It contains common fields and methods shared across all analysis services.
type BaseService struct {
	fs  afero.Fs
	log *logger.Logger
}

// NewBaseService creates a new BaseService with OS filesystem
func NewBaseService() *BaseService {
	fs := afero.NewOsFs()
	return &BaseService{
		fs:  fs,
		log: CreateLogger(fs, nil),
	}
}

// NewBaseServiceWithConfig creates a new BaseService with custom filesystem and logger config
func NewBaseServiceWithConfig(fs afero.Fs, loggerConfig *logger.Config) *BaseService {
	return &BaseService{
		fs:  fs,
		log: CreateLogger(fs, loggerConfig),
	}
}

// FS returns the filesystem
func (s *BaseService) FS() afero.Fs {
	return s.fs
}

// Log returns the logger
func (s *BaseService) Log() *logger.Logger {
	return s.log
}

// ExecuteTemplate parses and executes a template with the given data
func (s *BaseService) ExecuteTemplate(name, templateContent string, data map[string]string) (string, error) {
	// Validate template name to prevent potential issues
	if strings.TrimSpace(name) == "" {
		return "", fmt.Errorf("template name cannot be empty")
	}

	// Validate template content
	if strings.TrimSpace(templateContent) == "" {
		return "", fmt.Errorf("template content cannot be empty")
	}

	tmpl, err := template.New(name).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("parsing template %q: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template %q: %w", name, err)
	}

	return buf.String(), nil
}

// FormatOutput writes template content to .mission/templates/ and returns JSON with path
func (s *BaseService) FormatOutput(templateContent string) (string, error) {
	return FormatOutputWithFS(s.fs, templateContent)
}
