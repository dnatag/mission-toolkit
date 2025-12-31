package plan

import (
	"testing"
)

func TestPlanSpec_Validate(t *testing.T) {
	tests := []struct {
		name    string
		spec    PlanSpec
		wantErr bool
	}{
		{
			name: "valid spec with legacy scope",
			spec: PlanSpec{
				Intent:       "Test intent",
				Scope:        []string{"file1.go", "file2.go"},
				Domain:       []string{"security"},
				Plan:         []string{"step 1", "step 2"},
				Verification: "go test",
			},
			wantErr: false,
		},
		{
			name: "valid spec with new files field",
			spec: PlanSpec{
				Intent: "Test intent",
				Files: []FileSpec{
					{Path: "file1.go", Action: FileActionModify},
					{Path: "file2.go", Action: FileActionCreate},
				},
				Domain:       []string{"security"},
				Plan:         []string{"step 1", "step 2"},
				Verification: "go test",
			},
			wantErr: false,
		},
		{
			name: "missing intent",
			spec: PlanSpec{
				Scope:        []string{"file1.go"},
				Plan:         []string{"step 1"},
				Verification: "go test",
			},
			wantErr: true,
		},
		{
			name: "empty scope and files",
			spec: PlanSpec{
				Intent:       "Test intent",
				Scope:        []string{},
				Files:        []FileSpec{},
				Plan:         []string{"step 1"},
				Verification: "go test",
			},
			wantErr: true,
		},
		{
			name: "empty plan",
			spec: PlanSpec{
				Intent:       "Test intent",
				Scope:        []string{"file1.go"},
				Plan:         []string{},
				Verification: "go test",
			},
			wantErr: true,
		},
		{
			name: "missing verification",
			spec: PlanSpec{
				Intent: "Test intent",
				Scope:  []string{"file1.go"},
				Plan:   []string{"step 1"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.spec.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PlanSpec.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlanSpec_GetScopeFiles(t *testing.T) {
	tests := []struct {
		name     string
		spec     PlanSpec
		expected []string
	}{
		{
			name: "legacy scope only",
			spec: PlanSpec{
				Scope: []string{"file1.go", "file2.go"},
			},
			expected: []string{"file1.go", "file2.go"},
		},
		{
			name: "new files only",
			spec: PlanSpec{
				Files: []FileSpec{
					{Path: "file1.go", Action: FileActionModify},
					{Path: "file2.go", Action: FileActionCreate},
				},
			},
			expected: []string{"file1.go", "file2.go"},
		},
		{
			name: "both fields with overlap",
			spec: PlanSpec{
				Scope: []string{"file1.go", "file3.go"},
				Files: []FileSpec{
					{Path: "file1.go", Action: FileActionModify},
					{Path: "file2.go", Action: FileActionCreate},
				},
			},
			expected: []string{"file1.go", "file2.go", "file3.go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.spec.GetScopeFiles()
			if len(result) != len(tt.expected) {
				t.Errorf("GetScopeFiles() returned %d files, expected %d", len(result), len(tt.expected))
				return
			}

			for i, file := range result {
				if file != tt.expected[i] {
					t.Errorf("GetScopeFiles()[%d] = %s, expected %s", i, file, tt.expected[i])
				}
			}
		})
	}
}

func TestPlanSpec_GetFileAction(t *testing.T) {
	spec := PlanSpec{
		Files: []FileSpec{
			{Path: "file1.go", Action: FileActionModify},
			{Path: "file2.go", Action: FileActionCreate},
		},
	}

	tests := []struct {
		name     string
		filePath string
		expected FileAction
	}{
		{
			name:     "existing file with modify action",
			filePath: "file1.go",
			expected: FileActionModify,
		},
		{
			name:     "existing file with create action",
			filePath: "file2.go",
			expected: FileActionCreate,
		},
		{
			name:     "non-existing file defaults to modify",
			filePath: "file3.go",
			expected: FileActionModify,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := spec.GetFileAction(tt.filePath)
			if result != tt.expected {
				t.Errorf("GetFileAction(%s) = %s, expected %s", tt.filePath, result, tt.expected)
			}
		})
	}
}
