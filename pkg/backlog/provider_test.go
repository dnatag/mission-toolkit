package backlog

import "testing"

func TestIsValidType(t *testing.T) {
	valid := []string{"feature", "bugfix", "decomposed", "refactor", "future"}
	for _, v := range valid {
		if !IsValidType(v) {
			t.Errorf("expected %q to be valid", v)
		}
	}
	if IsValidType("invalid") {
		t.Error("expected 'invalid' to be invalid")
	}
}
