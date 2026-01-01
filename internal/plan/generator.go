package plan

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/afero"
)

// GeneratorService handles mission.md generation from plan.json
type GeneratorService struct {
	ServiceBase
}

// GenerateResult represents the result of mission generation
type GenerateResult struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	OutputFile  string `json:"output_file"`
	PlanFile    string `json:"plan_file"`
	MissionType string `json:"mission_type"`
	Track       int    `json:"track"`
}

// NewGeneratorService creates a new generator service
func NewGeneratorService(fs afero.Fs, missionID string) *GeneratorService {
	return &GeneratorService{
		ServiceBase: NewServiceBase(fs, missionID),
	}
}

// GenerateMission creates mission.md from plan.json specification
func (g *GeneratorService) GenerateMission(planFile, outputFile string) (*GenerateResult, error) {
	planData, err := afero.ReadFile(g.fs, planFile)
	if err != nil {
		return g.errorResult("Failed to read plan file: %v", err), err
	}

	var planSpec PlanSpec
	if err := json.Unmarshal(planData, &planSpec); err != nil {
		return g.errorResult("Failed to parse plan JSON: %v", err), err
	}

	complexity, err := NewAnalyzer(g.fs, g.missionID).AnalyzePlanFromSpec(&planSpec)
	if err != nil {
		return g.errorResult("Failed to analyze plan complexity: %v", err), err
	}

	missionContent := g.generateMissionContent(&planSpec, complexity)
	if err := afero.WriteFile(g.fs, outputFile, []byte(missionContent), 0644); err != nil {
		return g.errorResult("Failed to write mission file: %v", err), err
	}

	return &GenerateResult{
		Success: true, Message: "Mission file generated successfully",
		OutputFile: outputFile, PlanFile: planFile,
		MissionType: "WET", Track: complexity.Track,
	}, nil
}

// errorResult creates an error result
func (g *GeneratorService) errorResult(format string, err error) *GenerateResult {
	return &GenerateResult{Success: false, Message: fmt.Sprintf(format, err)}
}

// generateMissionContent creates the mission.md content from PlanSpec
func (g *GeneratorService) generateMissionContent(spec *PlanSpec, complexity *ComplexityResult) string {
	var content strings.Builder

	// Header
	content.WriteString(fmt.Sprintf("# MISSION\n\nid: %s\ntype: WET\ntrack: %d\niteration: 1\nstatus: planned\n\n", g.missionID, complexity.Track))

	// Intent
	content.WriteString(fmt.Sprintf("## INTENT\n%s\n\n", spec.Intent))

	// Scope
	content.WriteString("## SCOPE\n")
	for _, file := range spec.Scope {
		content.WriteString(file + "\n")
	}

	// Plan
	content.WriteString("\n## PLAN\n")
	for _, step := range spec.Plan {
		content.WriteString(fmt.Sprintf("- [ ] %s\n", step))
	}

	// Verification and Instructions
	content.WriteString(fmt.Sprintf("\n## VERIFICATION\n%s\n\n## EXECUTION INSTRUCTIONS\n⚠️  DO NOT EXECUTE THIS MISSION NOW\n- This is PLANNING PHASE only\n- Run: `@m.apply` to execute this mission (requires user approval)\n- Run: `@m.complete` to archive when finished", spec.Verification))

	return content.String()
}

// ToJSON converts GenerateResult to JSON string
func (r *GenerateResult) ToJSON() (string, error) {
	return MarshalToJSON(r)
}
