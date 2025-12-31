package plan

import (
	"testing"
)

func TestPlanSpec_Validate(t *testing.T) {
	tests := []struct {
		name    string
		spec    PlanSpec
		wantErr bool
	}{
		{
			name: "valid spec",
			spec: PlanSpec{
				Intent:       "Test intent",
				Scope:        []string{"file1.go", "file2.go"},
				Domain:       []string{"security"},
				Plan:         []string{"step 1", "step 2"},
				Verification: "go test",
			},
			wantErr: false,
		},
		{
			name: "missing intent",
			spec: PlanSpec{
				Scope:        []string{"file1.go"},
				Plan:         []string{"step 1"},
				Verification: "go test",
			},
			wantErr: true,
		},
		{
			name: "empty scope",
			spec: PlanSpec{
				Intent:       "Test intent",
				Scope:        []string{},
				Plan:         []string{"step 1"},
				Verification: "go test",
			},
			wantErr: true,
		},
		{
			name: "empty plan",
			spec: PlanSpec{
				Intent:       "Test intent",
				Scope:        []string{"file1.go"},
				Plan:         []string{},
				Verification: "go test",
			},
			wantErr: true,
		},
		{
			name: "missing verification",
			spec: PlanSpec{
				Intent: "Test intent",
				Scope:  []string{"file1.go"},
				Plan:   []string{"step 1"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.spec.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PlanSpec.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
