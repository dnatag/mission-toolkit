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
