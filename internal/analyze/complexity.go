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
	intent, err := s.readCurrentIntent()
	if err != nil {
		return "", fmt.Errorf("reading current intent: %w", err)
	}

	scope, err := s.readCurrentScope()
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

	return buf.String(), nil
}

// readCurrentIntent reads the INTENT section from .mission/mission.md
func (s *ComplexityService) readCurrentIntent() (string, error) {
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

// readCurrentScope reads the SCOPE section from .mission/mission.md
func (s *ComplexityService) readCurrentScope() (string, error) {
	missionPath := filepath.Join(".mission", "mission.md")
	content, err := afero.ReadFile(s.fs, missionPath)
	if err != nil {
		return "", fmt.Errorf("reading mission.md: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	inScopeSection := false
	var scopeLines []string

	for _, line := range lines {
		if strings.HasPrefix(line, "## SCOPE") {
			inScopeSection = true
			continue
		}
		if inScopeSection && strings.HasPrefix(line, "## ") {
			break
		}
		if inScopeSection && strings.TrimSpace(line) != "" {
			scopeLines = append(scopeLines, line)
		}
	}

	if len(scopeLines) == 0 {
		return "", fmt.Errorf("no scope found in mission.md")
	}

	return strings.Join(scopeLines, "\n"), nil
}
