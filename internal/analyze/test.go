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

//go:embed templates/test.md
var testTemplate string

// TestService provides test analysis templates
type TestService struct {
	fs afero.Fs
}

// NewTestService creates a new TestService
func NewTestService() *TestService {
	return &TestService{
		fs: afero.NewOsFs(),
	}
}

// NewTestServiceWithFS creates a new TestService with custom filesystem
func NewTestServiceWithFS(fs afero.Fs) *TestService {
	return &TestService{fs: fs}
}

// ProvideTemplate loads test.md template and injects current intent and scope from mission.md
func (s *TestService) ProvideTemplate() (string, error) {
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

	tmpl, err := template.New("test").Parse(testTemplate)
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

	return buf.String(), nil
}
