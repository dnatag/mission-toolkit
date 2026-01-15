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

//go:embed templates/decompose.md
var decomposeTemplate string

// DecomposeService provides decomposition analysis templates for Track 4 epics
type DecomposeService struct {
	fs  afero.Fs
	log *logger.Logger
}

// NewDecomposeService creates a new DecomposeService
func NewDecomposeService() *DecomposeService {
	fs := afero.NewOsFs()
	return &DecomposeService{
		fs:  fs,
		log: CreateLogger(fs, nil),
	}
}

// NewDecomposeServiceWithConfig creates a new DecomposeService with custom filesystem and logger config
func NewDecomposeServiceWithConfig(fs afero.Fs, loggerConfig *logger.Config) *DecomposeService {
	return &DecomposeService{
		fs:  fs,
		log: CreateLogger(fs, loggerConfig),
	}
}

// ProvideTemplate loads decompose.md template and injects current intent and scope from mission.md
func (s *DecomposeService) ProvideTemplate() (string, error) {
	s.log.LogStep(logger.LevelSuccess, "AnalyzeDecompose", "Starting decomposition analysis")

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

	tmpl, err := template.New("decompose").Parse(decomposeTemplate)
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
