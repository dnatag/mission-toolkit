package analyze

import (
	"bytes"
	_ "embed"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/dnatag/mission-toolkit/internal/mission"
	"github.com/spf13/afero"
)

//go:embed templates/duplication.md
var duplicationTemplate string

// DuplicationService provides duplication analysis templates
type DuplicationService struct {
	fs afero.Fs
}

// NewDuplicationService creates a new DuplicationService
func NewDuplicationService() *DuplicationService {
	return &DuplicationService{
		fs: afero.NewOsFs(),
	}
}

// NewDuplicationServiceWithFS creates a new DuplicationService with custom filesystem
func NewDuplicationServiceWithFS(fs afero.Fs) *DuplicationService {
	return &DuplicationService{fs: fs}
}

// ProvideTemplate loads duplication.md template and injects current intent from mission.md
func (s *DuplicationService) ProvideTemplate() (string, error) {
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
