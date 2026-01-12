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

//go:embed templates/intent.md
var intentTemplate string

// IntentService provides intent analysis templates
type IntentService struct {
	log *logger.Logger
}

// NewIntentService creates a new IntentService
func NewIntentService() *IntentService {
	reader := mission.NewReader(afero.NewOsFs())
	missionID, _ := reader.GetMissionID(filepath.Join(".mission", "mission.md"))
	return &IntentService{
		log: logger.New(missionID),
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

	return FormatOutput(buf.String())
}
