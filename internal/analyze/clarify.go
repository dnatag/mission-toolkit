package analyze

import (
	"bytes"
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/afero"
)

//go:embed templates/clarification.md
var clarificationTemplate string

// ClarifyService provides clarification analysis templates
type ClarifyService struct {
	fs afero.Fs
}

// NewClarifyService creates a new ClarifyService
func NewClarifyService() *ClarifyService {
	return &ClarifyService{
		fs: afero.NewOsFs(),
	}
}

// NewClarifyServiceWithFS creates a new ClarifyService with custom filesystem
func NewClarifyServiceWithFS(fs afero.Fs) *ClarifyService {
	return &ClarifyService{fs: fs}
}

// ProvideTemplate loads clarification.md template and injects current intent from mission.md
func (s *ClarifyService) ProvideTemplate() (string, error) {
	intent, err := s.readCurrentIntent()
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

	return buf.String(), nil
}

// readCurrentIntent reads the INTENT section from .mission/mission.md
func (s *ClarifyService) readCurrentIntent() (string, error) {
	missionPath := filepath.Join(".mission", "mission.md")
	content, err := afero.ReadFile(s.fs, missionPath)
	if err != nil {
		return "", fmt.Errorf("reading mission.md: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	inIntentSection := false
	var intentLines []string

	for _, line := range lines {
		if strings.HasPrefix(line, "## INTENT") {
			inIntentSection = true
			continue
		}
		// Stop at next section
		if inIntentSection && strings.HasPrefix(line, "## ") {
			break
		}
		// Collect non-empty lines in INTENT section
		if inIntentSection && strings.TrimSpace(line) != "" {
			intentLines = append(intentLines, line)
		}
	}

	if len(intentLines) == 0 {
		return "", fmt.Errorf("no intent found in mission.md")
	}

	return strings.TrimSpace(strings.Join(intentLines, " ")), nil
}
