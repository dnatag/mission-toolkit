package analyze

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"
)

//go:embed templates/intent.md
var intentTemplate string

// IntentService provides intent analysis templates
type IntentService struct{}

// NewIntentService creates a new IntentService
func NewIntentService() *IntentService {
	return &IntentService{}
}

// ProvideTemplate loads intent.md template and injects user input
func (s *IntentService) ProvideTemplate(userInput string) (string, error) {
	tmpl, err := template.New("intent").Parse(intentTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	data := map[string]string{"UserInput": userInput}
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}
