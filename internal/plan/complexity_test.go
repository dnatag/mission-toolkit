package plan

import (
	"testing"

	"github.com/spf13/afero"
)

func TestComplexityEngine_AnalyzeComplexity(t *testing.T) {
	fs := afero.NewMemMapFs()
	engine := NewComplexityEngine(fs, "test-mission")

	spec := &PlanSpec{
		Intent:       "Test intent",
		Scope:        []string{"file1.go", "file2.go", "file1_test.go"},
		Domain:       []string{"security"},
		Plan:         []string{"step 1", "step 2"},
		Verification: "go test",
	}

	result, err := engine.AnalyzeComplexity(spec)
	if err != nil {
		t.Fatalf("AnalyzeComplexity() error = %v", err)
	}

	if result.Track != 3 {
		t.Errorf("Expected track 3, got %d", result.Track)
	}

	if result.Confidence != "High" {
		t.Errorf("Expected confidence 'High', got '%s'", result.Confidence)
	}

	if result.Recommendation != "proceed" {
		t.Errorf("Expected recommendation 'proceed', got '%s'", result.Recommendation)
	}
}

func TestComplexityEngine_CountImplementationFiles(t *testing.T) {

	testCases := []struct {
		name     string
		scope    []string
		expected int
	}{
		{
			name:     "mixed files",
			scope:    []string{"file1.go", "file2.go", "file1_test.go", "README.md"},
			expected: 2,
		},
		{
			name:     "only test files",
			scope:    []string{"file1_test.go", "file2_test.go"},
			expected: 0,
		},
		{
			name:     "no files",
			scope:    []string{},
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := countImplementationFiles(tc.scope)
			if result != tc.expected {
				t.Errorf("Expected %d, got %d", tc.expected, result)
			}
		})
	}
}

func TestComplexityEngine_CalculateBaseTrack(t *testing.T) {

	testCases := []struct {
		fileCount int
		expected  int
	}{
		{0, 1},  // Atomic
		{3, 2},  // Standard
		{7, 3},  // Robust
		{12, 4}, // Epic
	}

	for _, tc := range testCases {
		result := calculateBaseTrack(tc.fileCount)
		if result != tc.expected {
			t.Errorf("Files %d: expected track %d, got %d", tc.fileCount, tc.expected, result)
		}
	}
}

func TestComplexityEngine_CalculateDomainMultipliers(t *testing.T) {

	testCases := []struct {
		name     string
		domains  []string
		expected int
	}{
		{
			name:     "security domain",
			domains:  []string{"security"},
			expected: 1,
		},
		{
			name:     "multiple domains",
			domains:  []string{"security", "performance", "unknown"},
			expected: 2,
		},
		{
			name:     "no domains",
			domains:  []string{},
			expected: 0,
		},
		{
			name:     "unknown domains",
			domains:  []string{"unknown", "invalid"},
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := calculateDomainMultipliers(tc.domains)
			if result != tc.expected {
				t.Errorf("Expected %d, got %d", tc.expected, result)
			}
		})
	}
}

func TestComplexityEngine_CalculateFinalTrack(t *testing.T) {

	testCases := []struct {
		baseTrack   int
		multipliers int
		expected    int
	}{
		{2, 1, 3},
		{3, 2, 3}, // Capped at 3
		{1, 0, 1},
		{4, 1, 4}, // Track 4 preserved (files > 9)
	}

	for _, tc := range testCases {
		result := calculateFinalTrack(tc.baseTrack, tc.multipliers)
		if result != tc.expected {
			t.Errorf("Base %d + multipliers %d: expected %d, got %d",
				tc.baseTrack, tc.multipliers, tc.expected, result)
		}
	}
}
