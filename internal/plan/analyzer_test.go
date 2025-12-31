package plan

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestAnalyzer_AnalyzePlan(t *testing.T) {
	fs := afero.NewMemMapFs()
	analyzer := NewAnalyzer(fs, "test-mission")

	planJSON := `{"intent":"Test intent","scope":["file1.go","file2.go","file1_test.go"],"domain":["security"],"plan":["step 1","step 2"],"verification":"go test"}`

	if err := afero.WriteFile(fs, "/tmp/plan.json", []byte(planJSON), 0644); err != nil {
		t.Fatalf("Failed to create test plan: %v", err)
	}

	result, err := analyzer.AnalyzePlan(fs, "/tmp/plan.json")
	if err != nil {
		t.Fatalf("AnalyzePlan() error = %v", err)
	}

	if result.Track != 3 || result.Recommendation != "proceed" {
		t.Errorf("Expected track=3, recommendation=proceed, got track=%d, recommendation=%s",
			result.Track, result.Recommendation)
	}
}

func TestAnalyzer_DetectFileActions(t *testing.T) {
	fs := afero.NewMemMapFs()
	analyzer := NewAnalyzer(fs, "test-mission")

	// Create an existing file
	if err := afero.WriteFile(fs, "existing.go", []byte("package main"), 0644); err != nil {
		t.Fatalf("Failed to create existing file: %v", err)
	}

	filePaths := []string{"existing.go", "new.go"}
	result := analyzer.detectFileActions(fs, filePaths)

	if len(result) != 2 {
		t.Fatalf("Expected 2 file specs, got %d", len(result))
	}

	// Check existing file is marked for modification
	if result[0].Path != "existing.go" || result[0].Action != FileActionModify {
		t.Errorf("Expected existing.go to be marked for modify, got %s:%s", result[0].Path, result[0].Action)
	}

	// Check new file is marked for creation
	if result[1].Path != "new.go" || result[1].Action != FileActionCreate {
		t.Errorf("Expected new.go to be marked for create, got %s:%s", result[1].Path, result[1].Action)
	}
}

func TestFormatResult(t *testing.T) {
	result := &ComplexityResult{
		Track: 2, Confidence: "High", Reasoning: "2 files = Track 2", Recommendation: "proceed",
	}

	json, err := FormatResult(result)
	if err != nil {
		t.Fatalf("FormatResult() error = %v", err)
	}

	if json == "" || !strings.Contains(json, `"track": 2`) {
		t.Errorf("Expected valid JSON with track field, got: %s", json)
	}
}
