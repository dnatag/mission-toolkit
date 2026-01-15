package analyze

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestBaseService_NewBaseService(t *testing.T) {
	service := NewBaseService()

	if service == nil {
		t.Fatal("NewBaseService returned nil")
	}

	if service.FS() == nil {
		t.Error("Expected fs to be initialized")
	}

	if service.Log() == nil {
		t.Error("Expected log to be initialized")
	}
}

func TestBaseService_NewBaseServiceWithConfig(t *testing.T) {
	fs := afero.NewMemMapFs()
	loggerConfig := CreateTestLoggerConfig(fs)

	service := NewBaseServiceWithConfig(fs, loggerConfig)

	if service == nil {
		t.Fatal("NewBaseServiceWithConfig returned nil")
	}

	if service.FS() != fs {
		t.Error("Expected fs to match provided filesystem")
	}

	if service.Log() == nil {
		t.Error("Expected log to be initialized")
	}
}

func TestBaseService_FS(t *testing.T) {
	fs := afero.NewMemMapFs()
	service := NewBaseServiceWithConfig(fs, CreateTestLoggerConfig(fs))

	if service.FS() != fs {
		t.Error("FS() did not return the expected filesystem")
	}
}

func TestBaseService_Log(t *testing.T) {
	fs := afero.NewMemMapFs()
	service := NewBaseServiceWithConfig(fs, CreateTestLoggerConfig(fs))

	if service.Log() == nil {
		t.Error("Log() returned nil")
	}
}

func TestBaseService_ExecuteTemplate(t *testing.T) {
	tests := []struct {
		name            string
		templateName    string
		templateContent string
		data            map[string]string
		wantOutput      string
		wantErr         bool
		errContains     string
	}{
		{
			name:            "simple substitution",
			templateName:    "test",
			templateContent: "Hello, {{.Name}}!",
			data:            map[string]string{"Name": "World"},
			wantOutput:      "Hello, World!",
			wantErr:         false,
		},
		{
			name:            "multiple substitutions",
			templateName:    "multi",
			templateContent: "{{.Greeting}}, {{.Name}}! Your status is {{.Status}}.",
			data:            map[string]string{"Greeting": "Hello", "Name": "Alice", "Status": "active"},
			wantOutput:      "Hello, Alice! Your status is active.",
			wantErr:         false,
		},
		{
			name:            "missing variable",
			templateName:    "missing",
			templateContent: "Hello, {{.Name}}!",
			data:            map[string]string{},
			wantOutput:      "Hello, <no value>!",
			wantErr:         false,
		},
		{
			name:            "empty template name",
			templateName:    "",
			templateContent: "Hello, {{.Name}}!",
			data:            map[string]string{"Name": "World"},
			wantErr:         true,
			errContains:     "template name cannot be empty",
		},
		{
			name:            "whitespace only template name",
			templateName:    "   ",
			templateContent: "Hello, {{.Name}}!",
			data:            map[string]string{"Name": "World"},
			wantErr:         true,
			errContains:     "template name cannot be empty",
		},
		{
			name:            "empty template content",
			templateName:    "test",
			templateContent: "",
			data:            map[string]string{"Name": "World"},
			wantErr:         true,
			errContains:     "template content cannot be empty",
		},
		{
			name:            "whitespace only template content",
			templateName:    "test",
			templateContent: "   ",
			data:            map[string]string{"Name": "World"},
			wantErr:         true,
			errContains:     "template content cannot be empty",
		},
		{
			name:            "invalid template syntax",
			templateName:    "invalid",
			templateContent: "Hello, {{.Name}",
			data:            map[string]string{"Name": "World"},
			wantErr:         true,
			errContains:     "parsing template",
		},
		{
			name:            "nil data",
			templateName:    "nil-data",
			templateContent: "Static content without variables",
			data:            nil,
			wantOutput:      "Static content without variables",
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			service := NewBaseServiceWithConfig(fs, CreateTestLoggerConfig(fs))

			got, err := service.ExecuteTemplate(tt.templateName, tt.templateContent, tt.data)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ExecuteTemplate() expected error containing %q, got nil", tt.errContains)
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ExecuteTemplate() error = %q, want error containing %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("ExecuteTemplate() unexpected error: %v", err)
			}

			if got != tt.wantOutput {
				t.Errorf("ExecuteTemplate() = %q, want %q", got, tt.wantOutput)
			}
		})
	}
}

func TestBaseService_FormatOutput(t *testing.T) {
	tests := []struct {
		name            string
		templateContent string
		wantErr         bool
		validateOutput  func(string) bool
	}{
		{
			name:            "valid template",
			templateContent: "# Test Template\n\nContent here.",
			wantErr:         false,
			validateOutput: func(output string) bool {
				return len(output) > 0 && strings.Contains(output, "template_path")
			},
		},
		{
			name:            "empty template",
			templateContent: "",
			wantErr:         false,
			validateOutput: func(output string) bool {
				return len(output) > 0
			},
		},
		{
			name:            "markdown template",
			templateContent: "# Header\n\n* Item 1\n* Item 2",
			wantErr:         false,
			validateOutput: func(output string) bool {
				return len(output) > 0 && strings.Contains(output, "template_path")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			service := NewBaseServiceWithConfig(fs, CreateTestLoggerConfig(fs))

			output, err := service.FormatOutput(tt.templateContent)

			if (err != nil) != tt.wantErr {
				t.Errorf("FormatOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.validateOutput(output) {
				t.Errorf("FormatOutput() output validation failed for: %s", output)
			}
		})
	}
}
