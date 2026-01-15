package analyze

import (
	_ "embed"
	"fmt"
	"path/filepath"

	"github.com/dnatag/mission-toolkit/internal/logger"
	"github.com/dnatag/mission-toolkit/internal/mission"
	"github.com/spf13/afero"
)

//go:embed templates/complexity.md
var complexityTemplate string

// ComplexityService provides complexity analysis templates
type ComplexityService struct {
	*BaseService
}

// NewComplexityService creates a new ComplexityService
func NewComplexityService() *ComplexityService {
	return &ComplexityService{
		BaseService: NewBaseService(),
	}
}

// NewComplexityServiceWithConfig creates a new ComplexityService with custom filesystem and logger config
func NewComplexityServiceWithConfig(fs afero.Fs, loggerConfig *logger.Config) *ComplexityService {
	return &ComplexityService{
		BaseService: NewBaseServiceWithConfig(fs, loggerConfig),
	}
}

// ProvideTemplate loads complexity.md template and injects current intent and scope from mission.md
func (s *ComplexityService) ProvideTemplate() (string, error) {
	s.Log().LogStep(logger.LevelSuccess, "AnalyzeComplexity", "Starting complexity analysis")

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

	output, err := s.ExecuteTemplate("complexity", complexityTemplate, map[string]string{
		"CurrentIntent": intent,
		"CurrentScope":  scope,
	})
	if err != nil {
		return "", err
	}

	return s.FormatOutput(output)
}
