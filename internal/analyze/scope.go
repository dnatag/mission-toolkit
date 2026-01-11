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

//go:embed templates/scope.md
var scopeTemplate string

// ScopeService provides scope analysis templates
type ScopeService struct {
	fs afero.Fs
}

// NewScopeService creates a new ScopeService
func NewScopeService() *ScopeService {
	return &ScopeService{
		fs: afero.NewOsFs(),
	}
}

// NewScopeServiceWithFS creates a new ScopeService with custom filesystem
func NewScopeServiceWithFS(fs afero.Fs) *ScopeService {
	return &ScopeService{fs: fs}
}

// ProvideTemplate loads scope.md template and injects current intent from mission.md
func (s *ScopeService) ProvideTemplate() (string, error) {
	reader := mission.NewReader(s.fs)
	missionPath := filepath.Join(".mission", "mission.md")

	intent, err := reader.ReadIntent(missionPath)
	if err != nil {
		return "", fmt.Errorf("reading current intent: %w", err)
	}

	tmpl, err := template.New("scope").Parse(scopeTemplate)
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
