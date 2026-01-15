package analyze

import (
	_ "embed"
	"fmt"
	"path/filepath"

	"github.com/dnatag/mission-toolkit/internal/logger"
	"github.com/dnatag/mission-toolkit/internal/mission"
	"github.com/spf13/afero"
)

//go:embed templates/clarification.md
var clarificationTemplate string

// ClarifyService provides clarification analysis templates
type ClarifyService struct {
	*BaseService
}

// NewClarifyService creates a new ClarifyService
func NewClarifyService() *ClarifyService {
	return &ClarifyService{
		BaseService: NewBaseService(),
	}
}

// NewClarifyServiceWithConfig creates a new ClarifyService with custom filesystem and logger config
func NewClarifyServiceWithConfig(fs afero.Fs, loggerConfig *logger.Config) *ClarifyService {
	return &ClarifyService{
		BaseService: NewBaseServiceWithConfig(fs, loggerConfig),
	}
}

// ProvideTemplate loads clarification.md template and injects current intent from mission.md
func (s *ClarifyService) ProvideTemplate() (string, error) {
	s.Log().LogStep(logger.LevelSuccess, "AnalyzeClarify", "Starting clarification analysis")

	missionPath := filepath.Join(".mission", "mission.md")
	reader := mission.NewReader(s.FS(), missionPath)

	intent, err := reader.ReadIntent()
	if err != nil {
		return "", fmt.Errorf("reading current intent: %w", err)
	}

	output, err := s.ExecuteTemplate("clarification", clarificationTemplate, map[string]string{"CurrentIntent": intent})
	if err != nil {
		return "", err
	}

	return s.FormatOutput(output)
}
