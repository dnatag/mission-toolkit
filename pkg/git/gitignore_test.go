package git

import (
	"testing"

	"github.com/spf13/afero"
)

func TestEnsureEntry(t *testing.T) {
	tests := []struct {
		name        string
		initial     string // empty string means file doesn't exist
		entry       string
		expected    string
		expectError bool
	}{
		{
			name:     "creates file if not exists",
			initial:  "",
			entry:    ".mission/",
			expected: ".mission/\n",
		},
		{
			name:     "skips if entry already exists",
			initial:  ".mission/\n",
			entry:    ".mission/",
			expected: ".mission/\n",
		},
		{
			name:     "appends entry if not present",
			initial:  "node_modules/\n",
			entry:    ".mission/",
			expected: "node_modules/\n.mission/\n",
		},
		{
			name:     "handles file without trailing newline",
			initial:  "node_modules/",
			entry:    ".mission/",
			expected: "node_modules/\n.mission/\n",
		},
		{
			name:     "skips entry with whitespace match",
			initial:  "  .mission/  \n",
			entry:    ".mission/",
			expected: "  .mission/  \n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()

			if tt.initial != "" {
				if err := afero.WriteFile(fs, ".gitignore", []byte(tt.initial), 0644); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			err := EnsureEntry(fs, ".", ".mission/")
			if (err != nil) != tt.expectError {
				t.Errorf("EnsureEntry() error = %v, expectError %v", err, tt.expectError)
				return
			}

			content, err := afero.ReadFile(fs, ".gitignore")
			if err != nil {
				t.Fatalf("failed to read .gitignore: %v", err)
			}

			if string(content) != tt.expected {
				t.Errorf("got %q, want %q", string(content), tt.expected)
			}
		})
	}
}
