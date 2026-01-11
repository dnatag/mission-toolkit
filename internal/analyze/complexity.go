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

//go:embed templates/complexity.md
var complexityTemplate string

// ComplexityService provides complexity analysis templates
type ComplexityService struct {
	fs afero.Fs
}

// NewComplexityService creates a new ComplexityService
func NewComplexityService() *ComplexityService {
	return &ComplexityService{
		fs: afero.NewOsFs(),
	}
}

// NewComplexityServiceWithFS creates a new ComplexityService with custom filesystem
func NewComplexityServiceWithFS(fs afero.Fs) *ComplexityService {
	return &ComplexityService{fs: fs}
}

// ProvideTemplate loads complexity.md template and injects current intent and scope from mission.md
func (s *ComplexityService) ProvideTemplate() (string, error) {
	reader := mission.NewReader(s.fs)
	missionPath := filepath.Join(".mission", "mission.md")

	intent, err := reader.ReadIntent(missionPath)
	if err != nil {
		return "", fmt.Errorf("reading current intent: %w", err)
	}

	scope, err := reader.ReadScope(missionPath)
	if err != nil {
		return "", fmt.Errorf("reading current scope: %w", err)
	}

	tmpl, err := template.New("complexity").Parse(complexityTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	data := map[string]string{
		"CurrentIntent": intent,
		"CurrentScope":  scope,
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return FormatOutput(buf.String())
}
