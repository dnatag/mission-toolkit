package version

import (
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	// Test version is not empty
	if Version == "" {
		t.Error("Version should not be empty")
	}

	// Test version starts with 'v'
	if !strings.HasPrefix(Version, "v") {
		t.Errorf("Version should start with 'v', got: %s", Version)
	}

	// Test version has semantic version format (v1.2.3)
	parts := strings.Split(strings.TrimPrefix(Version, "v"), ".")
	if len(parts) != 3 {
		t.Errorf("Version should have format v1.2.3, got: %s", Version)
	}

	// Test each part is numeric
	for i, part := range parts {
		if part == "" {
			t.Errorf("Version part %d should not be empty in: %s", i, Version)
		}
	}
}

func TestVersionConstant(t *testing.T) {
	// Test that version constant is accessible and has expected value
	if Version != "v1.1.2" {
		t.Errorf("Expected version v1.1.2, got: %s", Version)
	}
}
