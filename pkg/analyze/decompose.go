package analyze

import (
	_ "embed"
	"fmt"
	"path/filepath"

	"github.com/dnatag/mission-toolkit/pkg/logger"
	"github.com/dnatag/mission-toolkit/pkg/mission"
	"github.com/spf13/afero"
)

//go:embed templates/decompose.md
var decomposeTemplate string

// DecomposeService provides decomposition analysis templates for Track 4 epics
type DecomposeService struct {
	*BaseService
}

// NewDecomposeService creates a new DecomposeService
func NewDecomposeService() *DecomposeService {
	return &DecomposeService{
		BaseService: NewBaseService(),
	}
}

// NewDecomposeServiceWithConfig creates a new DecomposeService with custom filesystem and logger config
func NewDecomposeServiceWithConfig(fs afero.Fs, loggerConfig *logger.Config) *DecomposeService {
	return &DecomposeService{
		BaseService: NewBaseServiceWithConfig(fs, loggerConfig),
	}
}

// ProvideTemplate loads decompose.md template and injects current intent and scope from mission.md
func (s *DecomposeService) ProvideTemplate() (string, error) {
	s.Log().LogStep(logger.LevelSuccess, "AnalyzeDecompose", "Starting decomposition analysis")

	missionPath := filepath.Join(".mission", "mission.md")
	reader := mission.NewReader(s.FS(), missionPath)

	intent, err := reader.ReadIntent()
	if err != nil {
		return "", fmt.Errorf("reading current intent: %w", err)
	}

	scope, err := reader.ReadScope()
	if err != nil {
		return "", fmt.Errorf("reading current scope: %w", err)
	}

	output, err := s.ExecuteTemplate("decompose", decomposeTemplate, map[string]string{
		"CurrentIntent": intent,
		"CurrentScope":  scope,
	})
	if err != nil {
		return "", err
	}

	return s.FormatOutput(output)
}
