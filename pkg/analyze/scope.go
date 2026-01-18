package analyze

import (
	_ "embed"
	"fmt"
	"path/filepath"

	"github.com/dnatag/mission-toolkit/pkg/logger"
	"github.com/dnatag/mission-toolkit/pkg/mission"
	"github.com/spf13/afero"
)

//go:embed templates/scope.md
var scopeTemplate string

// ScopeService provides scope analysis templates
type ScopeService struct {
	*BaseService
}

// NewScopeService creates a new ScopeService
func NewScopeService() *ScopeService {
	return &ScopeService{
		BaseService: NewBaseService(),
	}
}

// NewScopeServiceWithConfig creates a new ScopeService with custom filesystem and logger config
func NewScopeServiceWithConfig(fs afero.Fs, loggerConfig *logger.Config) *ScopeService {
	return &ScopeService{
		BaseService: NewBaseServiceWithConfig(fs, loggerConfig),
	}
}

// ProvideTemplate loads scope.md template and injects current intent from mission.md
func (s *ScopeService) ProvideTemplate() (string, error) {
	s.Log().LogStep(logger.LevelSuccess, "AnalyzeScope", "Starting scope analysis")

	missionPath := filepath.Join(".mission", "mission.md")
	reader := mission.NewReader(s.FS(), missionPath)

	intent, err := reader.ReadIntent()
	if err != nil {
		return "", fmt.Errorf("reading current intent: %w", err)
	}

	output, err := s.ExecuteTemplate("scope", scopeTemplate, map[string]string{"CurrentIntent": intent})
	if err != nil {
		return "", err
	}

	return s.FormatOutput(output)
}
