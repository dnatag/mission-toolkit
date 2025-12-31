package plan

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dnatag/mission-toolkit/internal/logger"
	"github.com/spf13/afero"
)

// ComplexityResult represents the result of complexity analysis
type ComplexityResult struct {
	Track          int      `json:"track"`
	Confidence     string   `json:"confidence"`
	Reasoning      string   `json:"reasoning"`
	Recommendation string   `json:"recommendation"`
	Warnings       []string `json:"warnings,omitempty"`
	TestGaps       []string `json:"test_gaps,omitempty"`
}

// ComplexityEngine calculates mission complexity
type ComplexityEngine struct {
	log *logger.Logger
}

var (
	complexDomains = map[string]bool{
		"security": true, "performance": true, "complex-algo": true, "high-risk": true,
		"cross-cutting": true, "authentication": true, "authorization": true,
		"cryptography": true, "real-time": true, "compliance": true,
	}
	recommendations = [5]string{"", "atomic-edit", "proceed", "proceed", "decompose"}
	programmingExts = map[string]bool{
		".go": true, ".java": true, ".kt": true, ".cs": true, ".rs": true,
		".py": true, ".ts": true, ".js": true, ".c": true, ".cpp": true,
		".cc": true, ".cxx": true, ".h": true, ".hpp": true,
	}
	nonSourceExts = map[string]bool{
		".md": true, ".txt": true, ".json": true, ".yaml": true, ".yml": true,
		".xml": true, ".html": true, ".css": true, ".png": true, ".jpg": true,
	}
)

// NewComplexityEngine creates a new complexity analysis engine
func NewComplexityEngine(fs afero.Fs, missionID string) *ComplexityEngine {
	return &ComplexityEngine{log: logger.New(missionID)}
}

// AnalyzeComplexity calculates track complexity for a plan spec
func (e *ComplexityEngine) AnalyzeComplexity(spec *PlanSpec) (*ComplexityResult, error) {
	implFiles := countImplementationFiles(spec.Scope)
	baseTrack := calculateBaseTrack(implFiles)
	multipliers := calculateDomainMultipliers(spec.Domain)
	finalTrack := calculateFinalTrack(baseTrack, multipliers)
	testGaps := detectTestGaps(spec.Scope)

	e.log.LogStep("INFO", "complexity-result",
		fmt.Sprintf("Track %d (files:%d, base:%d, mult:%d)", finalTrack, implFiles, baseTrack, multipliers))

	return &ComplexityResult{
		Track:          finalTrack,
		Confidence:     calculateConfidence(implFiles, spec.Domain),
		Reasoning:      generateReasoning(baseTrack, multipliers, implFiles, spec.Domain),
		Recommendation: generateRecommendation(finalTrack),
		TestGaps:       testGaps,
	}, nil
}

// countImplementationFiles counts non-test files
func countImplementationFiles(scope []string) int {
	count := 0
	for _, file := range scope {
		if !strings.Contains(file, "_test.") && !strings.HasSuffix(file, ".md") {
			count++
		}
	}
	return count
}

// calculateBaseTrack determines base complexity track
func calculateBaseTrack(fileCount int) int {
	if fileCount == 0 {
		return 1
	}
	if fileCount <= 5 {
		return 2
	}
	if fileCount <= 9 {
		return 3
	}
	return 4
}

// calculateDomainMultipliers applies domain-based complexity increases
func calculateDomainMultipliers(domains []string) int {
	multipliers := 0
	for _, domain := range domains {
		if complexDomains[strings.ToLower(domain)] {
			multipliers++
		}
	}
	return multipliers
}

// calculateFinalTrack applies multipliers with max cap
func calculateFinalTrack(baseTrack, multipliers int) int {
	final := baseTrack + multipliers
	if final > 3 {
		return 3
	}
	return final
}

// calculateConfidence determines confidence level
func calculateConfidence(fileCount int, domains []string) string {
	if fileCount > 0 && len(domains) > 0 {
		return "High"
	}
	if fileCount > 0 || len(domains) > 0 {
		return "Medium"
	}
	return "Low"
}

// generateReasoning creates human-readable reasoning
func generateReasoning(baseTrack, multipliers, fileCount int, domains []string) string {
	base := fmt.Sprintf("%d files = Track %d", fileCount, baseTrack)
	if multipliers > 0 {
		return fmt.Sprintf("%s + %d domain multipliers (%v) = Track %d",
			base, multipliers, domains, baseTrack+multipliers)
	}
	return base
}

// generateRecommendation provides action recommendation
func generateRecommendation(track int) string {
	if track > 0 && track < len(recommendations) {
		return recommendations[track]
	}
	return "review"
}

// detectTestGaps identifies missing test files using language-specific patterns
func detectTestGaps(scope []string) []string {
	testFiles := make(map[string]bool, len(scope)/2)
	sourceFiles := make([]string, 0, len(scope))

	for _, file := range scope {
		if isTestFile(file) {
			if sourceFile := getSourceFileFromTest(file); sourceFile != "" {
				testFiles[sourceFile] = true
			}
		} else if isSourceFile(file) {
			sourceFiles = append(sourceFiles, file)
		}
	}

	var gaps []string
	for _, sourceFile := range sourceFiles {
		if !testFiles[sourceFile] {
			if testFile := getTestFileForSource(sourceFile); testFile != "" {
				gaps = append(gaps, testFile)
			}
		}
	}
	return gaps
}

// isTestFile checks if a file is a test file using language-specific patterns
func isTestFile(file string) bool {
	ext := strings.ToLower(filepath.Ext(file))
	filename := strings.ToLower(filepath.Base(file))

	switch ext {
	case ".go":
		return strings.Contains(file, "_test.go")
	case ".java":
		return strings.Contains(file, "Test.java")
	case ".cs":
		return strings.Contains(file, "Test.cs")
	case ".rs":
		return strings.Contains(file, "tests/") || strings.HasSuffix(file, "_test.rs")
	case ".py":
		return strings.HasPrefix(filename, "test_") || strings.HasSuffix(file, "_test.py")
	case ".js", ".ts":
		return strings.Contains(file, ".test.") || strings.Contains(file, ".spec.")
	case ".c", ".cpp", ".cc", ".cxx":
		return strings.Contains(file, "_test.")
	default:
		return strings.Contains(filename, "test")
	}
}

// isSourceFile checks if a file is a source file
func isSourceFile(file string) bool {
	if isTestFile(file) || strings.HasSuffix(file, ".md") {
		return false
	}
	ext := strings.ToLower(filepath.Ext(file))
	if programmingExts[ext] {
		return true
	}
	// For unknown extensions, consider it a source file if it has an extension
	// and isn't a common non-source extension
	if ext != "" {
		return !nonSourceExts[ext]
	}
	return false
}

// getSourceFileFromTest maps a test file back to its source file
func getSourceFileFromTest(testFile string) string {
	ext := strings.ToLower(filepath.Ext(testFile))
	switch ext {
	case ".go":
		return strings.Replace(testFile, "_test.go", ".go", 1)
	case ".java":
		return strings.Replace(testFile, "Test.java", ".java", 1)
	case ".cs":
		return strings.Replace(testFile, "Test.cs", ".cs", 1)
	case ".rs":
		if strings.Contains(testFile, "tests/") {
			base := filepath.Base(testFile)
			if strings.HasSuffix(base, "_test.rs") {
				return strings.Replace(base, "_test.rs", ".rs", 1)
			}
		}
		return strings.Replace(testFile, "_test.rs", ".rs", 1)
	case ".py":
		if strings.HasPrefix(filepath.Base(testFile), "test_") {
			return strings.Replace(testFile, "test_", "", 1)
		}
		return strings.Replace(testFile, "_test.py", ".py", 1)
	case ".js", ".ts":
		if strings.Contains(testFile, ".test.") {
			return strings.Replace(testFile, ".test.", ".", 1)
		}
		return strings.Replace(testFile, ".spec.", ".", 1)
	case ".c", ".cpp", ".cc", ".cxx":
		return strings.Replace(testFile, "_test.", ".", 1)
	}
	return ""
}

// getTestFileForSource generates the expected test file name for a source file
func getTestFileForSource(sourceFile string) string {
	ext := strings.ToLower(filepath.Ext(sourceFile))
	switch ext {
	case ".go":
		return strings.Replace(sourceFile, ".go", "_test.go", 1)
	case ".java":
		return strings.Replace(sourceFile, ".java", "Test.java", 1)
	case ".cs":
		return strings.Replace(sourceFile, ".cs", "Test.cs", 1)
	case ".rs":
		base := filepath.Base(sourceFile)
		return filepath.Join("tests", strings.Replace(base, ".rs", "_test.rs", 1))
	case ".py":
		base := filepath.Base(sourceFile)
		dir := filepath.Dir(sourceFile)
		return filepath.Join(dir, "test_"+base)
	case ".js":
		return strings.Replace(sourceFile, ".js", ".test.js", 1)
	case ".ts":
		return strings.Replace(sourceFile, ".ts", ".test.ts", 1)
	case ".c":
		return strings.Replace(sourceFile, ".c", "_test.c", 1)
	case ".cpp", ".cc", ".cxx":
		return strings.Replace(sourceFile, ext, "_test"+ext, 1)
	default:
		if ext != "" {
			return strings.TrimSuffix(sourceFile, ext) + "_test" + ext
		}
	}
	return ""
}
