package plan

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/afero"
)

// Analyzer handles plan analysis operations
type Analyzer struct {
	complexity *ComplexityEngine
}

// NewAnalyzer creates a new plan analyzer
func NewAnalyzer(fs afero.Fs, missionID string) *Analyzer {
	return &Analyzer{
		complexity: NewComplexityEngine(fs, missionID),
	}
}

// AnalyzePlan performs comprehensive analysis of a plan specification
func (a *Analyzer) AnalyzePlan(fs afero.Fs, planFile string) (*ComplexityResult, error) {
	spec, err := LoadPlanSpec(fs, planFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load plan spec: %w", err)
	}

	// Populate file action annotations if not already present
	if len(spec.Files) == 0 && len(spec.Scope) > 0 {
		spec.Files = a.detectFileActions(fs, spec.Scope)
	}

	return a.complexity.AnalyzeComplexity(spec)
}

// AnalyzePlanFromSpec performs analysis on an existing PlanSpec
func (a *Analyzer) AnalyzePlanFromSpec(spec *PlanSpec) (*ComplexityResult, error) {
	return a.complexity.AnalyzeComplexity(spec)
}

// detectFileActions analyzes file paths to determine if they should be modified or created
func (a *Analyzer) detectFileActions(fs afero.Fs, filePaths []string) []FileSpec {
	files := make([]FileSpec, 0, len(filePaths))

	for _, filePath := range filePaths {
		action := FileActionCreate
		if exists, _ := afero.Exists(fs, filePath); exists {
			action = FileActionModify
		}

		files = append(files, FileSpec{
			Path:   filePath,
			Action: action,
		})
	}

	return files
}

// FormatResult converts analysis result to JSON
func FormatResult(result *ComplexityResult) (string, error) {
	data, err := json.MarshalIndent(result, "", "  ")
	return string(data), err
}
