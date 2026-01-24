package diagnosis

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestCreateDiagnosis(t *testing.T) {
	tests := []struct {
		name    string
		symptom string
		wantErr bool
	}{
		{
			name:    "creates diagnosis with symptom",
			symptom: "API returns 500 error on login",
			wantErr: false,
		},
		{
			name:    "rejects empty symptom",
			symptom: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			path := ".mission/diagnosis.md"

			err := CreateDiagnosis(fs, path, tt.symptom)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify file exists
			exists, err := afero.Exists(fs, path)
			require.NoError(t, err)
			require.True(t, exists)

			// Verify content structure
			content, err := afero.ReadFile(fs, path)
			require.NoError(t, err)

			contentStr := string(content)
			require.True(t, strings.HasPrefix(contentStr, "---\n"), "should start with frontmatter")
			require.Contains(t, contentStr, "id: DIAG-")
			require.Contains(t, contentStr, "status: investigating")
			require.Contains(t, contentStr, "confidence: low")
			require.Contains(t, contentStr, "## SYMPTOM")
			require.Contains(t, contentStr, tt.symptom)
			require.Contains(t, contentStr, "## INVESTIGATION")
			require.Contains(t, contentStr, "## HYPOTHESES")
			require.Contains(t, contentStr, "## ROOT CAUSE")
			require.Contains(t, contentStr, "## AFFECTED FILES")
			require.Contains(t, contentStr, "## RECOMMENDED FIX")
			require.Contains(t, contentStr, "## REPRODUCTION")
		})
	}
}

func TestDiagnosisExists(t *testing.T) {
	tests := []struct {
		name       string
		setupFile  bool
		wantExists bool
		wantErr    bool
	}{
		{
			name:       "file exists",
			setupFile:  true,
			wantExists: true,
			wantErr:    false,
		},
		{
			name:       "file does not exist",
			setupFile:  false,
			wantExists: false,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			path := ".mission/diagnosis.md"

			if tt.setupFile {
				fs.MkdirAll(".mission", 0755)
				afero.WriteFile(fs, path, []byte("test content"), 0644)
			}

			exists, err := DiagnosisExists(fs, path)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantExists, exists)
		})
	}
}

func TestUpdateSection(t *testing.T) {
	tests := []struct {
		name           string
		initialContent string
		section        string
		newContent     string
		wantContains   string
		wantErr        bool
	}{
		{
			name: "updates existing section",
			initialContent: `---
id: DIAG-123
status: investigating
confidence: low
created: 2026-01-24T10:00:00Z
symptom: test
---

## SYMPTOM
test

## ROOT CAUSE
To be determined
`,
			section:      "ROOT CAUSE",
			newContent:   "Fixed issue",
			wantContains: "Fixed issue",
			wantErr:      false,
		},
		{
			name: "adds new section if not exists",
			initialContent: `---
id: DIAG-123
status: investigating
confidence: low
created: 2026-01-24T10:00:00Z
symptom: test
---

## SYMPTOM
test
`,
			section:      "ROOT CAUSE",
			newContent:   "New finding",
			wantContains: "New finding",
			wantErr:      false,
		},
		{
			name: "preserves other sections",
			initialContent: `---
id: DIAG-123
status: investigating
confidence: low
created: 2026-01-24T10:00:00Z
symptom: test
---

## SYMPTOM
test

## ROOT CAUSE
old

## AFFECTED FILES
- file.go
`,
			section:      "ROOT CAUSE",
			newContent:   "updated",
			wantContains: "## AFFECTED FILES",
			wantErr:      false,
		},
		{
			name: "rejects empty section",
			initialContent: `---
id: DIAG-123
status: investigating
confidence: low
created: 2026-01-24T10:00:00Z
symptom: test
---

## SYMPTOM
test
`,
			section:    "",
			newContent: "content",
			wantErr:    true,
		},
		{
			name: "rejects empty content",
			initialContent: `---
id: DIAG-123
status: investigating
confidence: low
created: 2026-01-24T10:00:00Z
symptom: test
---

## SYMPTOM
test
`,
			section:    "ROOT CAUSE",
			newContent: "",
			wantErr:    true,
		},
		{
			name: "appends to INVESTIGATION section",
			initialContent: `---
id: DIAG-123
status: investigating
confidence: low
created: 2026-01-24T10:00:00Z
symptom: test
---

## SYMPTOM
test

## INVESTIGATION
- [x] Checked file.go

## ROOT CAUSE
TBD
`,
			section:      "INVESTIGATION",
			newContent:   "- [ ] Review logs",
			wantContains: "- [x] Checked file.go",
			wantErr:      false,
		},
		{
			name: "appends to HYPOTHESES section",
			initialContent: `---
id: DIAG-123
status: investigating
confidence: low
created: 2026-01-24T10:00:00Z
symptom: test
---

## SYMPTOM
test

## HYPOTHESES
1. **[HIGH]** First hypothesis

## ROOT CAUSE
TBD
`,
			section:      "HYPOTHESES",
			newContent:   "2. **[LOW]** Second hypothesis",
			wantContains: "1. **[HIGH]** First hypothesis",
			wantErr:      false,
		},
		{
			name: "appends to AFFECTED FILES section",
			initialContent: `---
id: DIAG-123
status: investigating
confidence: low
created: 2026-01-24T10:00:00Z
symptom: test
---

## SYMPTOM
test

## AFFECTED FILES
- file1.go

## ROOT CAUSE
TBD
`,
			section:      "AFFECTED FILES",
			newContent:   "- file2.go",
			wantContains: "- file1.go",
			wantErr:      false,
		},
		{
			name: "replaces non-list section",
			initialContent: `---
id: DIAG-123
status: investigating
confidence: low
created: 2026-01-24T10:00:00Z
symptom: test
---

## SYMPTOM
test

## ROOT CAUSE
Old cause

## AFFECTED FILES
- file.go
`,
			section:      "ROOT CAUSE",
			newContent:   "New cause",
			wantContains: "New cause",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			path := ".mission/diagnosis.md"

			// Create initial file
			fs.MkdirAll(".mission", 0755)
			afero.WriteFile(fs, path, []byte(tt.initialContent), 0644)

			err := UpdateSection(fs, path, tt.section, tt.newContent)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify updated content
			content, err := afero.ReadFile(fs, path)
			require.NoError(t, err)
			require.Contains(t, string(content), tt.wantContains)

			// For append tests, verify new content is also present
			if strings.Contains(tt.name, "appends") {
				require.Contains(t, string(content), tt.newContent)
			}

			// For replace tests, verify old content is gone
			if tt.name == "replaces non-list section" {
				require.NotContains(t, string(content), "Old cause")
			}
		})
	}
}

func TestUpdateFrontmatter(t *testing.T) {
	tests := []struct {
		name           string
		initialContent string
		status         string
		confidence     string
		wantErr        bool
		wantStatus     string
		wantConfidence string
	}{
		{
			name: "updates status only",
			initialContent: `---
id: DIAG-123
status: investigating
confidence: low
created: 2026-01-24T10:00:00Z
symptom: test
---

## SYMPTOM
test
`,
			status:         "confirmed",
			confidence:     "",
			wantErr:        false,
			wantStatus:     "confirmed",
			wantConfidence: "low",
		},
		{
			name: "updates confidence only",
			initialContent: `---
id: DIAG-123
status: investigating
confidence: low
created: 2026-01-24T10:00:00Z
symptom: test
---

## SYMPTOM
test
`,
			status:         "",
			confidence:     "high",
			wantErr:        false,
			wantStatus:     "investigating",
			wantConfidence: "high",
		},
		{
			name: "updates both status and confidence",
			initialContent: `---
id: DIAG-123
status: investigating
confidence: low
created: 2026-01-24T10:00:00Z
symptom: test
---

## SYMPTOM
test
`,
			status:         "confirmed",
			confidence:     "high",
			wantErr:        false,
			wantStatus:     "confirmed",
			wantConfidence: "high",
		},
		{
			name: "rejects invalid status",
			initialContent: `---
id: DIAG-123
status: investigating
confidence: low
created: 2026-01-24T10:00:00Z
symptom: test
---

## SYMPTOM
test
`,
			status:     "invalid",
			confidence: "",
			wantErr:    true,
		},
		{
			name: "rejects invalid confidence",
			initialContent: `---
id: DIAG-123
status: investigating
confidence: low
created: 2026-01-24T10:00:00Z
symptom: test
---

## SYMPTOM
test
`,
			status:     "",
			confidence: "invalid",
			wantErr:    true,
		},
		{
			name: "rejects empty status and confidence",
			initialContent: `---
id: DIAG-123
status: investigating
confidence: low
created: 2026-01-24T10:00:00Z
symptom: test
---

## SYMPTOM
test
`,
			status:     "",
			confidence: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			path := ".mission/diagnosis.md"

			// Create initial file
			fs.MkdirAll(".mission", 0755)
			afero.WriteFile(fs, path, []byte(tt.initialContent), 0644)

			err := UpdateFrontmatter(fs, path, tt.status, tt.confidence)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify updated frontmatter
			diag, err := ReadDiagnosis(fs, path)
			require.NoError(t, err)
			require.Equal(t, tt.wantStatus, diag.Status)
			require.Equal(t, tt.wantConfidence, diag.Confidence)
		})
	}
}

func TestFinalize(t *testing.T) {
	tests := []struct {
		name             string
		diagnosisContent string
		wantValid        bool
		wantMissing      []string
	}{
		{
			name: "valid diagnosis with investigating status",
			diagnosisContent: `---
id: DIAG-123
status: investigating
confidence: low
created: 2026-01-24T10:00:00Z
symptom: test
---

## SYMPTOM
test

## INVESTIGATION
- [ ] Check logs

## HYPOTHESES
1. **[HIGH]** First hypothesis

## ROOT CAUSE
TBD

## AFFECTED FILES
- file.go
`,
			wantValid:   true,
			wantMissing: nil,
		},
		{
			name: "valid diagnosis with confirmed status and recommended fix",
			diagnosisContent: `---
id: DIAG-123
status: confirmed
confidence: high
created: 2026-01-24T10:00:00Z
symptom: test
---

## SYMPTOM
test

## INVESTIGATION
- [x] Checked logs

## HYPOTHESES
1. **[HIGH]** First hypothesis

## ROOT CAUSE
Session not initialized

## AFFECTED FILES
- file.go

## RECOMMENDED FIX
Initialize session before use
`,
			wantValid:   true,
			wantMissing: nil,
		},
		{
			name: "invalid - confirmed status missing recommended fix",
			diagnosisContent: `---
id: DIAG-123
status: confirmed
confidence: high
created: 2026-01-24T10:00:00Z
symptom: test
---

## SYMPTOM
test

## INVESTIGATION
- [x] Checked logs

## HYPOTHESES
1. **[HIGH]** First hypothesis

## ROOT CAUSE
Session not initialized

## AFFECTED FILES
- file.go
`,
			wantValid:   false,
			wantMissing: []string{"RECOMMENDED FIX"},
		},
		{
			name: "missing required sections",
			diagnosisContent: `---
id: DIAG-123
status: investigating
confidence: low
created: 2026-01-24T10:00:00Z
symptom: test
---

## SYMPTOM
test

## INVESTIGATION
- [ ] Check logs
`,
			wantValid:   false,
			wantMissing: []string{"HYPOTHESES", "ROOT CAUSE", "AFFECTED FILES"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			path := ".mission/diagnosis.md"

			// Create diagnosis file
			fs.MkdirAll(".mission", 0755)
			afero.WriteFile(fs, path, []byte(tt.diagnosisContent), 0644)

			result, err := Finalize(fs, path)
			require.NoError(t, err)

			// Parse JSON result
			var resultMap map[string]interface{}
			err = json.Unmarshal([]byte(result), &resultMap)
			require.NoError(t, err)

			require.Equal(t, tt.wantValid, resultMap["valid"])

			if !tt.wantValid {
				missingSections := resultMap["missing_sections"].([]interface{})
				require.Len(t, missingSections, len(tt.wantMissing))
			}
		})
	}
}
