package plan

import (
	"fmt"
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
	NextStep       string   `json:"next_step"`
	Warnings       []string `json:"warnings,omitempty"`
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
)

// NewComplexityEngine creates a new complexity analysis engine
func NewComplexityEngine(fs afero.Fs, missionID string) *ComplexityEngine {
	// Create logger config with same filesystem as the engine
	config := logger.DefaultConfig()
	config.Fs = fs

	return &ComplexityEngine{log: logger.NewWithConfig(missionID, config)}
}

// AnalyzeComplexity calculates track complexity for a plan spec
func (e *ComplexityEngine) AnalyzeComplexity(spec *PlanSpec) (*ComplexityResult, error) {
	allFiles := spec.GetScopeFiles()
	implFiles := countImplementationFiles(allFiles)
	baseTrack := calculateBaseTrack(implFiles)
	multipliers := calculateDomainMultipliers(spec.Domain)
	finalTrack := calculateFinalTrack(baseTrack, multipliers)

	// Warnings slice is now for future, non-test-gap warnings.
	var warnings []string

	e.log.LogStep("INFO", "complexity-result",
		fmt.Sprintf("Track %d (files:%d, base:%d, mult:%d)", finalTrack, implFiles, baseTrack, multipliers))

	return &ComplexityResult{
		Track:          finalTrack,
		Confidence:     calculateConfidence(implFiles, spec.Domain),
		Reasoning:      generateReasoning(baseTrack, multipliers, implFiles, spec.Domain),
		Recommendation: generateRecommendation(finalTrack),
		NextStep:       generateNextStep(finalTrack),
		Warnings:       warnings,
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
	// If base track is already 4 (files > 9), keep it 4 regardless of multipliers
	if baseTrack >= 4 {
		return 4
	}

	final := baseTrack + multipliers
	// Cap at 3 for normal flows, unless base was already 4
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

// generateNextStep provides explicit instructions for the AI
func generateNextStep(track int) string {
	switch track {
	case 1:
		return "STOP. Use template libraries/displays/plan-atomic.md to provide a direct code suggestion."
	case 2, 3:
		return "PROCEED to Step 4 (Validation)."
	case 4:
		return "STOP. Decompose this intent into 3-5 sub-intents in .mission/backlog.md and use template libraries/displays/plan-epic.md."
	default:
		return "Review the analysis and decide whether to proceed or decompose."
	}
}
