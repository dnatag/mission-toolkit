package analyze

import (
	_ "embed"

	"github.com/dnatag/mission-toolkit/pkg/logger"
	"github.com/spf13/afero"
)

//go:embed templates/intent.md
var intentTemplate string

// IntentService provides intent analysis templates
type IntentService struct {
	*BaseService
}

// NewIntentService creates a new IntentService
func NewIntentService() *IntentService {
	return &IntentService{
		BaseService: NewBaseService(),
	}
}

// NewIntentServiceWithConfig creates a new IntentService with custom filesystem and logger config
func NewIntentServiceWithConfig(fs afero.Fs, loggerConfig *logger.Config) *IntentService {
	return &IntentService{
		BaseService: NewBaseServiceWithConfig(fs, loggerConfig),
	}
}

// ProvideTemplate loads intent.md template and injects user input
func (s *IntentService) ProvideTemplate(userInput string) (string, error) {
	s.Log().LogStep(logger.LevelSuccess, "AnalyzeIntent", "Starting intent analysis")

	output, err := s.ExecuteTemplate("intent", intentTemplate, map[string]string{"UserInput": userInput})
	if err != nil {
		return "", err
	}

	return s.FormatOutput(output)
}
