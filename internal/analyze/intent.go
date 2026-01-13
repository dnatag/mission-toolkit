package analyze

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/dnatag/mission-toolkit/internal/logger"
	"github.com/spf13/afero"
)

//go:embed templates/intent.md
var intentTemplate string

// IntentService provides intent analysis templates
type IntentService struct {
	fs  afero.Fs
	log *logger.Logger
}

// NewIntentService creates a new IntentService
func NewIntentService() *IntentService {
	fs := afero.NewOsFs()
	return &IntentService{
		fs:  fs,
		log: CreateLogger(fs, nil),
	}
}

// NewIntentServiceWithConfig creates a new IntentService with custom filesystem and logger config
func NewIntentServiceWithConfig(fs afero.Fs, loggerConfig *logger.Config) *IntentService {
	return &IntentService{
		fs:  fs,
		log: CreateLogger(fs, loggerConfig),
	}
}

// ProvideTemplate loads intent.md template and injects user input
func (s *IntentService) ProvideTemplate(userInput string) (string, error) {
	s.log.LogStep(logger.LevelSuccess, "AnalyzeIntent", "Starting intent analysis")

	tmpl, err := template.New("intent").Parse(intentTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	data := map[string]string{"UserInput": userInput}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return FormatOutputWithFS(s.fs, buf.String())
}
