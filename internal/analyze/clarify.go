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

//go:embed templates/clarification.md
var clarificationTemplate string

// ClarifyService provides clarification analysis templates
type ClarifyService struct {
	fs  afero.Fs
	log *logger.Logger
}

// NewClarifyService creates a new ClarifyService
func NewClarifyService() *ClarifyService {
	fs := afero.NewOsFs()
	reader := mission.NewReader(fs)
	missionID, _ := reader.GetMissionID(filepath.Join(".mission", "mission.md"))
	return &ClarifyService{
		fs:  fs,
		log: logger.New(missionID),
	}
}

// NewClarifyServiceWithFS creates a new ClarifyService with custom filesystem
func NewClarifyServiceWithFS(fs afero.Fs) *ClarifyService {
	reader := mission.NewReader(fs)
	missionID, _ := reader.GetMissionID(filepath.Join(".mission", "mission.md"))
	return &ClarifyService{
		fs:  fs,
		log: logger.New(missionID),
	}
}

// ProvideTemplate loads clarification.md template and injects current intent from mission.md
func (s *ClarifyService) ProvideTemplate() (string, error) {
	s.log.LogStep(logger.LevelSuccess, "AnalyzeClarify", "Starting clarification analysis")

	reader := mission.NewReader(s.fs)
	missionPath := filepath.Join(".mission", "mission.md")

	intent, err := reader.ReadIntent(missionPath)
	if err != nil {
		return "", fmt.Errorf("reading current intent: %w", err)
	}

	tmpl, err := template.New("clarification").Parse(clarificationTemplate)
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
