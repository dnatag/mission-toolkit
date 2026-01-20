package analyze

import (
	_ "embed"
	"fmt"
	"path/filepath"

	"github.com/dnatag/mission-toolkit/pkg/logger"
	"github.com/dnatag/mission-toolkit/pkg/mission"
	"github.com/spf13/afero"
)

//go:embed templates/test.md
var testTemplate string

// TestService provides test analysis templates
type TestService struct {
	*BaseService
}

// NewTestService creates a new TestService
func NewTestService() *TestService {
	return &TestService{
		BaseService: NewBaseService(),
	}
}

// NewTestServiceWithConfig creates a new TestService with custom filesystem and logger config
func NewTestServiceWithConfig(fs afero.Fs, loggerConfig *logger.Config) *TestService {
	return &TestService{
		BaseService: NewBaseServiceWithConfig(fs, loggerConfig),
	}
}

// ProvideTemplate loads test.md template and injects current intent and scope from mission.md
func (s *TestService) ProvideTemplate() (string, error) {
	s.Log().LogStep(logger.LevelSuccess, "AnalyzeTest", "Starting test analysis")

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

	output, err := s.ExecuteTemplate("test", testTemplate, map[string]string{
		"CurrentIntent": intent,
		"CurrentScope":  scope,
	})
	if err != nil {
		return "", err
	}

	return s.FormatOutput(output)
}
