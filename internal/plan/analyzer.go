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
	return a.complexity.AnalyzeComplexity(spec)
}

// AnalyzePlanFromSpec performs analysis on an existing PlanSpec
func (a *Analyzer) AnalyzePlanFromSpec(spec *PlanSpec) (*ComplexityResult, error) {
	return a.complexity.AnalyzeComplexity(spec)
}

// FormatResult converts analysis result to JSON
func FormatResult(result *ComplexityResult) (string, error) {
	data, err := json.MarshalIndent(result, "", "  ")
	return string(data), err
}
