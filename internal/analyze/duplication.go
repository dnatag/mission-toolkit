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
	intent, err := s.readCurrentIntent()
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

	return buf.String(), nil
}

// readCurrentIntent reads the INTENT section from .mission/mission.md
func (s *DuplicationService) readCurrentIntent() (string, error) {
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
		if inIntentSection && strings.HasPrefix(line, "## ") {
			break
		}
		if inIntentSection && strings.TrimSpace(line) != "" {
			intentLines = append(intentLines, line)
		}
	}

	if len(intentLines) == 0 {
		return "", fmt.Errorf("no intent found in mission.md")
	}

	return strings.TrimSpace(strings.Join(intentLines, " ")), nil
}
