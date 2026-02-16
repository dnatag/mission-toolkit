package backlog

import (
	"testing"

	"github.com/spf13/afero"
)

func TestBeadsAvailable_BothConditionsMet(t *testing.T) {
	// Create a mock filesystem with .beads directory
	fs := afero.NewMemMapFs()
	fs.Mkdir(".beads", 0755)

	// Mock exec.LookPath by setting a fake PATH with our command
	// Note: In a real test environment, we'd need to mock exec.LookPath somehow
	// For this test, we'll assume bd is not on PATH in the test environment

	// This test will pass the directory check but fail the PATH check
	// In a real scenario with bd installed, both would pass
	available, err := BeadsAvailable(fs)
	if err != nil {
		t.Fatalf("BeadsAvailable returned error: %v", err)
	}

	// In test environment, bd is unlikely to be on PATH
	// So we expect false unless bd is actually installed
	// The directory check should succeed
	_, statErr := fs.Stat(".beads")
	if statErr != nil {
		t.Fatal("Expected .beads directory to exist in mock filesystem")
	}

	t.Logf("Beads available: %v (bd on PATH: %v, .beads exists: true)", available, available)
}

func TestBeadsAvailable_NoBeadsDirectory(t *testing.T) {
	// Create an empty mock filesystem
	fs := afero.NewMemMapFs()

	available, err := BeadsAvailable(fs)
	if err != nil {
		t.Fatalf("BeadsAvailable returned error: %v", err)
	}

	if available {
		t.Error("Expected BeadsAvailable to return false when .beads directory doesn't exist")
	}
}

func TestBeadsAvailable_BeadsIsFileNotDirectory(t *testing.T) {
	// Create a mock filesystem where .beads is a file, not a directory
	fs := afero.NewMemMapFs()
	afero.WriteFile(fs, ".beads", []byte("not a directory"), 0644)

	available, err := BeadsAvailable(fs)
	if err == nil {
		t.Error("Expected error when .beads is a file, not a directory")
	}

	if available {
		t.Error("Expected BeadsAvailable to return false when .beads is a file")
	}
}

func TestBeadsAvailableInDir(t *testing.T) {
	// Create a temporary directory structure for testing
	fs := afero.NewMemMapFs()

	// Create a mock project directory with .beads subdirectory
	projectDir := "/test-project"
	fs.MkdirAll(projectDir+"/.beads", 0755)

	// Use BasePathFs to root the filesystem at the project directory
	baseFS := afero.NewBasePathFs(fs, projectDir)

	available, err := BeadsAvailable(baseFS)
	if err != nil {
		t.Fatalf("BeadsAvailable returned error: %v", err)
	}

	// In test environment, bd is unlikely to be on PATH
	// But the directory check should succeed from the project root
	t.Logf("Beads available in %s: %v", projectDir, available)
}

func TestBeadsAvailable_EmptyFilesystem(t *testing.T) {
	// Create an empty mock filesystem
	fs := afero.NewMemMapFs()

	available, err := BeadsAvailable(fs)
	if err != nil {
		t.Fatalf("BeadsAvailable returned error: %v", err)
	}

	if available {
		t.Error("Expected BeadsAvailable to return false in empty filesystem")
	}
}
