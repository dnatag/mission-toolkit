package templates

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

func TestWriteTemplates(t *testing.T) {
	tests := []struct {
		name      string
		aiType    string
		targetDir string
		wantFiles []string
	}{
		{
			name:      "Amazon Q templates",
			aiType:    "q",
			targetDir: "/test",
			wantFiles: []string{
				".idd/governance.md",
				".idd/metrics.md",
				".idd/backlog.md",
				".amazonq/prompts/idd.complete.md",
				".amazonq/prompts/idd.plan.md",
				".amazonq/prompts/idd.apply.md",
			},
		},
		{
			name:      "Claude templates",
			aiType:    "claude",
			targetDir: "/test",
			wantFiles: []string{
				".idd/governance.md",
				".idd/metrics.md",
				".idd/backlog.md",
				".claude/commands/idd.complete.md",
				".claude/commands/idd.plan.md",
				".claude/commands/idd.apply.md",
			},
		},
		{
			name:      "Gemini templates",
			aiType:    "gemini",
			targetDir: "/test",
			wantFiles: []string{
				".idd/governance.md",
				".idd/metrics.md",
				".idd/backlog.md",
				".gemini/commands/idd.complete.md",
				".gemini/commands/idd.plan.md",
				".gemini/commands/idd.apply.md",
			},
		},
		{
			name:      "Cursor templates",
			aiType:    "cursor",
			targetDir: "/test",
			wantFiles: []string{
				".idd/governance.md",
				".idd/metrics.md",
				".idd/backlog.md",
				".cursor/commands/idd.complete.md",
				".cursor/commands/idd.plan.md",
				".cursor/commands/idd.apply.md",
			},
		},
		{
			name:      "Codex templates",
			aiType:    "codex",
			targetDir: "/test",
			wantFiles: []string{
				".idd/governance.md",
				".idd/metrics.md",
				".idd/backlog.md",
				".codex/commands/idd.complete.md",
				".codex/commands/idd.plan.md",
				".codex/commands/idd.apply.md",
			},
		},
		{
			name:      "Cline templates",
			aiType:    "cline",
			targetDir: "/test",
			wantFiles: []string{
				".idd/governance.md",
				".idd/metrics.md",
				".idd/backlog.md",
				".clinerules/workflows/idd.complete.md",
				".clinerules/workflows/idd.plan.md",
				".clinerules/workflows/idd.apply.md",
			},
		},
		{
			name:      "Kiro templates",
			aiType:    "kiro",
			targetDir: "/test",
			wantFiles: []string{
				".idd/governance.md",
				".idd/metrics.md",
				".idd/backlog.md",
				".kiro/prompts/idd.complete.md",
				".kiro/prompts/idd.plan.md",
				".kiro/prompts/idd.apply.md",
			},
		},
		{
			name:      "Default templates",
			aiType:    "default",
			targetDir: "/test",
			wantFiles: []string{
				".idd/governance.md",
				".idd/metrics.md",
				".idd/backlog.md",
				"prompts/idd.complete.md",
				"prompts/idd.plan.md",
				"prompts/idd.apply.md",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()

			err := WriteTemplates(fs, tt.targetDir, tt.aiType)
			if err != nil {
				t.Fatalf("WriteTemplates() error = %v", err)
			}

			for _, file := range tt.wantFiles {
				fullPath := filepath.Join(tt.targetDir, file)
				exists, err := afero.Exists(fs, fullPath)
				if err != nil {
					t.Errorf("Error checking file %s: %v", fullPath, err)
				}
				if !exists {
					t.Errorf("Expected file %s does not exist", fullPath)
				}

				// Verify file has content
				content, err := afero.ReadFile(fs, fullPath)
				if err != nil {
					t.Errorf("Error reading file %s: %v", fullPath, err)
				}
				if len(content) == 0 {
					t.Errorf("File %s is empty", fullPath)
				}
			}
		})
	}
}
