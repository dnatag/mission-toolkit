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

//go:embed templates/scope.md
var scopeTemplate string

// ScopeService provides scope analysis templates
type ScopeService struct {
	fs  afero.Fs
	log *logger.Logger
}

// NewScopeService creates a new ScopeService
func NewScopeService() *ScopeService {
	fs := afero.NewOsFs()
	return &ScopeService{
		fs:  fs,
		log: CreateLogger(fs, nil),
	}
}

// NewScopeServiceWithConfig creates a new ScopeService with custom filesystem and logger config
func NewScopeServiceWithConfig(fs afero.Fs, loggerConfig *logger.Config) *ScopeService {
	return &ScopeService{
		fs:  fs,
		log: CreateLogger(fs, loggerConfig),
	}
}

// ProvideTemplate loads scope.md template and injects current intent from mission.md
func (s *ScopeService) ProvideTemplate() (string, error) {
	s.log.LogStep(logger.LevelSuccess, "AnalyzeScope", "Starting scope analysis")

	missionPath := filepath.Join(".mission", "mission.md")
	reader := mission.NewReader(s.fs, missionPath)

	intent, err := reader.ReadIntent()
	if err != nil {
		return "", fmt.Errorf("reading current intent: %w", err)
	}

	tmpl, err := template.New("scope").Parse(scopeTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	data := map[string]string{"CurrentIntent": intent}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return FormatOutputWithFS(s.fs, buf.String())
}
