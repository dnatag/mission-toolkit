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
		{4, 1, 3}, // Capped at 3
	}

	for _, tc := range testCases {
		result := calculateFinalTrack(tc.baseTrack, tc.multipliers)
		if result != tc.expected {
			t.Errorf("Base %d + multipliers %d: expected %d, got %d",
				tc.baseTrack, tc.multipliers, tc.expected, result)
		}
	}
}

func TestComplexityEngine_DetectTestGaps(t *testing.T) {
	testCases := []struct {
		name     string
		scope    []string
		expected []string
	}{
		// Go language tests
		{
			name:     "Go - missing tests",
			scope:    []string{"file1.go", "file2.go"},
			expected: []string{"file1_test.go", "file2_test.go"},
		},
		{
			name:     "Go - has tests",
			scope:    []string{"file1.go", "file1_test.go"},
			expected: []string{},
		},
		{
			name:     "Go - mixed files",
			scope:    []string{"file1.go", "file2.go", "file1_test.go", "README.md"},
			expected: []string{"file2_test.go"},
		},

		// Java language tests
		{
			name:     "Java - missing tests",
			scope:    []string{"User.java", "Service.java"},
			expected: []string{"UserTest.java", "ServiceTest.java"},
		},
		{
			name:     "Java - has tests",
			scope:    []string{"User.java", "UserTest.java"},
			expected: []string{},
		},

		// Python language tests
		{
			name:     "Python - missing tests",
			scope:    []string{"user.py", "service.py"},
			expected: []string{"test_user.py", "test_service.py"},
		},
		{
			name:     "Python - has tests",
			scope:    []string{"user.py", "test_user.py"},
			expected: []string{},
		},
		{
			name:     "Python - suffix style tests",
			scope:    []string{"user.py", "user_test.py"},
			expected: []string{},
		},

		// TypeScript language tests
		{
			name:     "TypeScript - missing tests",
			scope:    []string{"user.ts", "service.ts"},
			expected: []string{"user.test.ts", "service.test.ts"},
		},
		{
			name:     "TypeScript - has test files",
			scope:    []string{"user.ts", "user.test.ts"},
			expected: []string{},
		},
		{
			name:     "TypeScript - spec files",
			scope:    []string{"user.ts", "user.spec.ts"},
			expected: []string{},
		},

		// JavaScript language tests
		{
			name:     "JavaScript - missing tests",
			scope:    []string{"user.js", "service.js"},
			expected: []string{"user.test.js", "service.test.js"},
		},

		// C# language tests
		{
			name:     "C# - missing tests",
			scope:    []string{"User.cs", "Service.cs"},
			expected: []string{"UserTest.cs", "ServiceTest.cs"},
		},
		{
			name:     "C# - has tests",
			scope:    []string{"User.cs", "UserTest.cs"},
			expected: []string{},
		},

		// Rust language tests
		{
			name:     "Rust - missing tests",
			scope:    []string{"user.rs", "service.rs"},
			expected: []string{"tests/user_test.rs", "tests/service_test.rs"},
		},
		{
			name:     "Rust - has tests",
			scope:    []string{"user.rs", "tests/user_test.rs"},
			expected: []string{},
		},

		// C/C++ language tests
		{
			name:     "C - missing tests",
			scope:    []string{"user.c", "service.c"},
			expected: []string{"user_test.c", "service_test.c"},
		},
		{
			name:     "C++ - missing tests",
			scope:    []string{"user.cpp", "service.cpp"},
			expected: []string{"user_test.cpp", "service_test.cpp"},
		},

		// Unknown language fallback
		{
			name:     "Unknown language - keyword fallback",
			scope:    []string{"user.xyz", "service.xyz"},
			expected: []string{"user_test.xyz", "service_test.xyz"},
		},

		// Mixed languages
		{
			name:     "Mixed languages",
			scope:    []string{"user.go", "Service.java", "api.py", "user_test.go"},
			expected: []string{"ServiceTest.java", "test_api.py"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := detectTestGaps(tc.scope)
			if len(result) != len(tc.expected) {
				t.Errorf("Expected %d gaps, got %d. Expected: %v, Got: %v",
					len(tc.expected), len(result), tc.expected, result)
				return
			}

			for i, gap := range result {
				if gap != tc.expected[i] {
					t.Errorf("Expected gap '%s', got '%s'", tc.expected[i], gap)
				}
			}
		})
	}
}
