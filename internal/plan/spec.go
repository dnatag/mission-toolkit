package plan

import "fmt"

// PlanSpec represents the structure for mission planning data
type PlanSpec struct {
	Intent       string   `json:"intent"`
	Scope        []string `json:"scope"`
	Domain       []string `json:"domain"`
	Plan         []string `json:"plan"`
	Verification string   `json:"verification"`
}

// Validate checks if the PlanSpec has required fields
func (p *PlanSpec) Validate() error {
	if p.Intent == "" {
		return fmt.Errorf("intent is required")
	}
	if len(p.Scope) == 0 {
		return fmt.Errorf("scope is required")
	}
	if len(p.Plan) == 0 {
		return fmt.Errorf("plan is required")
	}
	if p.Verification == "" {
		return fmt.Errorf("verification is required")
	}
	return nil
}
