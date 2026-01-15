package analyze

import (
	"bytes"
	_ "embed"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/dnatag/mission-toolkit/internal/logger"
	"github.com/dnatag/mission-toolkit/internal/mission"
	"github.com/spf13/afero"
)

//go:embed templates/complexity.md
var complexityTemplate string

// ComplexityService provides complexity analysis templates
type ComplexityService struct {
	fs  afero.Fs
	log *logger.Logger
}

// NewComplexityService creates a new ComplexityService
func NewComplexityService() *ComplexityService {
	fs := afero.NewOsFs()
	return &ComplexityService{
		fs:  fs,
		log: CreateLogger(fs, nil),
	}
}

// NewComplexityServiceWithConfig creates a new ComplexityService with custom filesystem and logger config
func NewComplexityServiceWithConfig(fs afero.Fs, loggerConfig *logger.Config) *ComplexityService {
	return &ComplexityService{
		fs:  fs,
		log: CreateLogger(fs, loggerConfig),
	}
}

// ProvideTemplate loads complexity.md template and injects current intent and scope from mission.md
func (s *ComplexityService) ProvideTemplate() (string, error) {
	s.log.LogStep(logger.LevelSuccess, "AnalyzeComplexity", "Starting complexity analysis")

	missionPath := filepath.Join(".mission", "mission.md")
	reader := mission.NewReader(s.fs, missionPath)

	intent, err := reader.ReadIntent()
	if err != nil {
		return "", fmt.Errorf("reading current intent: %w", err)
	}

	scope, err := reader.ReadScope()
	if err != nil {
		return "", fmt.Errorf("reading current scope: %w", err)
	}

	tmpl, err := template.New("complexity").Parse(complexityTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	data := map[string]string{
		"CurrentIntent": intent,
		"CurrentScope":  scope,
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return FormatOutputWithFS(s.fs, buf.String())
}
