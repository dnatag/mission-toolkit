package mission

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArchiver(t *testing.T) {
	tests := []struct {
		name          string
		setupFiles    map[string]string
		expectedFiles []string
		expectedError string
	}{
		{
			name: "successful archive",
			setupFiles: map[string]string{
				".mission/mission.md": `---
id: 20260103220303-9114
iteration: 1
status: active
track: 2
type: WET
---

## INTENT
Test mission

## SCOPE
test.go

## PLAN
- [ ] Test step

## VERIFICATION
echo "test"`,
				".mission/execution.log": "2026-01-03 22:03:03 [INFO] Test log entry",
			},
			expectedFiles: []string{
				".mission/completed/20260103220303-9114-mission.md",
				".mission/completed/20260103220303-9114-execution.log",
			},
		},
		{
			name: "missing mission file",
			setupFiles: map[string]string{
				".mission/execution.log": "test log",
			},
			expectedError: "getting mission ID",
		},
		{
			name: "missing execution log",
			setupFiles: map[string]string{
				".mission/mission.md": `---
id: test-123
---

## INTENT
Test`,
			},
			expectedError: "archiving execution.log",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup filesystem
			fs := afero.NewMemMapFs()
			missionDir := ".mission"

			// Create test files
			for path, content := range tt.setupFiles {
				dir := filepath.Dir(path)
				require.NoError(t, fs.MkdirAll(dir, 0755))
				require.NoError(t, afero.WriteFile(fs, path, []byte(content), 0644))
			}

			// Create archiver and run
			archiver := NewArchiver(fs, missionDir)
			err := archiver.Archive()

			// Check error expectation
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)

			// Check expected files exist
			for _, expectedFile := range tt.expectedFiles {
				exists, err := afero.Exists(fs, expectedFile)
				require.NoError(t, err)
				assert.True(t, exists, "Expected file %s to exist", expectedFile)

				// Verify content was copied
				content, err := afero.ReadFile(fs, expectedFile)
				require.NoError(t, err)
				assert.NotEmpty(t, content)
			}
		})
	}
}

func TestArchiver_getMissionID(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectedID  string
		expectError bool
	}{
		{
			name: "valid mission with ID",
			content: `---
id: 20260103220303-9114
status: active
---

## INTENT
Test`,
			expectedID: "20260103220303-9114",
		},
		{
			name: "mission without ID",
			content: `---
status: active
---

## INTENT
Test`,
			expectError: false, // Should fallback to timestamp
		},
		{
			name:        "invalid file",
			content:     "short",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			missionPath := ".mission/mission.md"

			require.NoError(t, fs.MkdirAll(".mission", 0755))
			require.NoError(t, afero.WriteFile(fs, missionPath, []byte(tt.content), 0644))

			archiver := NewArchiver(fs, ".mission")
			id, err := archiver.getMissionID(missionPath)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.expectedID != "" {
				assert.Equal(t, tt.expectedID, id)
			} else {
				assert.NotEmpty(t, id)
				assert.Contains(t, id, "archived-")
			}
		})
	}
}
