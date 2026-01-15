package analyze

import (
	_ "embed"
	"fmt"
	"path/filepath"

	"github.com/dnatag/mission-toolkit/internal/logger"
	"github.com/dnatag/mission-toolkit/internal/mission"
	"github.com/spf13/afero"
)

//go:embed templates/duplication.md
var duplicationTemplate string

// DuplicationService provides duplication analysis templates
type DuplicationService struct {
	*BaseService
}

// NewDuplicationService creates a new DuplicationService
func NewDuplicationService() *DuplicationService {
	return &DuplicationService{
		BaseService: NewBaseService(),
	}
}

// NewDuplicationServiceWithConfig creates a new DuplicationService with custom filesystem and logger config
func NewDuplicationServiceWithConfig(fs afero.Fs, loggerConfig *logger.Config) *DuplicationService {
	return &DuplicationService{
		BaseService: NewBaseServiceWithConfig(fs, loggerConfig),
	}
}

// ProvideTemplate loads duplication.md template and injects current intent from mission.md
func (s *DuplicationService) ProvideTemplate() (string, error) {
	s.Log().LogStep(logger.LevelSuccess, "AnalyzeDuplication", "Starting duplication analysis")

	missionPath := filepath.Join(".mission", "mission.md")
	reader := mission.NewReader(s.FS(), missionPath)

	intent, err := reader.ReadIntent()
	if err != nil {
		return "", fmt.Errorf("reading current intent: %w", err)
	}

	output, err := s.ExecuteTemplate("duplication", duplicationTemplate, map[string]string{"CurrentIntent": intent})
	if err != nil {
		return "", err
	}

	return s.FormatOutput(output)
}
