package validation

import "testing"

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantValid bool
	}{
		{"empty string", "", false},
		{"whitespace only", "   ", false},
		{"placeholder", "$ARGUMENTS", false},
		{"placeholder with spaces", "  $ARGUMENTS  ", false},
		{"valid input", "Add feature", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate(tt.input)
			if result.IsValid != tt.wantValid {
				t.Errorf("Validate(%q).IsValid = %v, want %v", tt.input, result.IsValid, tt.wantValid)
			}
		})
	}
}
