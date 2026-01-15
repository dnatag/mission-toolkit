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

//go:embed templates/clarification.md
var clarificationTemplate string

// ClarifyService provides clarification analysis templates
type ClarifyService struct {
	fs  afero.Fs
	log *logger.Logger
}

// NewClarifyService creates a new ClarifyService
func NewClarifyService() *ClarifyService {
	fs := afero.NewOsFs()
	return &ClarifyService{
		fs:  fs,
		log: CreateLogger(fs, nil),
	}
}

// NewClarifyServiceWithConfig creates a new ClarifyService with custom filesystem and logger config
func NewClarifyServiceWithConfig(fs afero.Fs, loggerConfig *logger.Config) *ClarifyService {
	return &ClarifyService{
		fs:  fs,
		log: CreateLogger(fs, loggerConfig),
	}
}

// ProvideTemplate loads clarification.md template and injects current intent from mission.md
func (s *ClarifyService) ProvideTemplate() (string, error) {
	s.log.LogStep(logger.LevelSuccess, "AnalyzeClarify", "Starting clarification analysis")

	missionPath := filepath.Join(".mission", "mission.md")
	reader := mission.NewReader(s.fs, missionPath)

	intent, err := reader.ReadIntent()
	if err != nil {
		return "", fmt.Errorf("reading current intent: %w", err)
	}

	tmpl, err := template.New("clarification").Parse(clarificationTemplate)
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
