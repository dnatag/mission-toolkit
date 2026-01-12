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

//go:embed templates/duplication.md
var duplicationTemplate string

// DuplicationService provides duplication analysis templates
type DuplicationService struct {
	fs  afero.Fs
	log *logger.Logger
}

// NewDuplicationService creates a new DuplicationService
func NewDuplicationService() *DuplicationService {
	fs := afero.NewOsFs()
	reader := mission.NewReader(fs)
	missionID, _ := reader.GetMissionID(filepath.Join(".mission", "mission.md"))
	return &DuplicationService{
		fs:  fs,
		log: logger.New(missionID),
	}
}

// NewDuplicationServiceWithFS creates a new DuplicationService with custom filesystem
func NewDuplicationServiceWithFS(fs afero.Fs) *DuplicationService {
	reader := mission.NewReader(fs)
	missionID, _ := reader.GetMissionID(filepath.Join(".mission", "mission.md"))
	return &DuplicationService{
		fs:  fs,
		log: logger.New(missionID),
	}
}

// ProvideTemplate loads duplication.md template and injects current intent from mission.md
func (s *DuplicationService) ProvideTemplate() (string, error) {
	s.log.LogStep(logger.LevelSuccess, "AnalyzeDuplication", "Starting duplication analysis")

	reader := mission.NewReader(s.fs)
	missionPath := filepath.Join(".mission", "mission.md")

	intent, err := reader.ReadIntent(missionPath)
	if err != nil {
		return "", fmt.Errorf("reading current intent: %w", err)
	}

	tmpl, err := template.New("duplication").Parse(duplicationTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	data := map[string]string{"CurrentIntent": intent}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return FormatOutput(buf.String())
}
