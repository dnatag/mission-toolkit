package backlog

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

func TestNewProvider_BeadsAvailable(t *testing.T) {
	// Create a temporary directory with .beads subdirectory
	tempDir, err := os.MkdirTemp("", "provider-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create .beads directory to simulate Beads availability
	beadsDir := filepath.Join(tempDir, ".beads")
	if err := os.Mkdir(beadsDir, 0755); err != nil {
		t.Fatalf("Failed to create .beads directory: %v", err)
	}

	// Call NewProvider
	provider := NewProvider(tempDir)

	// Verify provider is not nil
	if provider == nil {
		t.Fatal("Expected non-nil provider")
	}

	// Note: We cannot verify the exact type (BeadsProvider vs BacklogManager)
	// without exported type assertions, and the actual type depends on whether
	// 'bd' is on PATH. The factory pattern abstracts this away.
	// We verify the provider implements the interface by calling a method.
	_, err = provider.List(nil, nil)
	if err != nil {
		// Error is acceptable (no backlog.md), but provider should not be nil
		t.Logf("List() returned error (expected for new directory): %v", err)
	}
}

func TestNewProvider_BeadsNotAvailable(t *testing.T) {
	// Create a temporary directory WITHOUT .beads subdirectory
	tempDir, err := os.MkdirTemp("", "provider-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Call NewProvider
	provider := NewProvider(tempDir)

	// Verify provider is not nil
	if provider == nil {
		t.Fatal("Expected non-nil provider")
	}

	// Verify provider implements the interface
	_, err = provider.List(nil, nil)
	if err != nil {
		// Error is acceptable (no backlog.md), but provider should not be nil
		t.Logf("List() returned error (expected for new directory): %v", err)
	}
}

func TestNewProvider_EmptyDirectory(t *testing.T) {
	// Create an empty temporary directory
	tempDir, err := os.MkdirTemp("", "provider-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Call NewProvider with empty directory
	provider := NewProvider(tempDir)

	// Verify provider is not nil (should return BacklogManager as fallback)
	if provider == nil {
		t.Fatal("Expected non-nil provider")
	}
}

func TestNewProvider_InvalidPath(t *testing.T) {
	// Call NewProvider with non-existent path
	// The factory should still return a valid provider (BacklogManager)
	provider := NewProvider("/nonexistent/path")

	// Verify provider is not nil
	if provider == nil {
		t.Fatal("Expected non-nil provider even for invalid path")
	}
}

func TestBeadsAvailableInDir_Integration(t *testing.T) {
	// Integration test with afero filesystem
	fs := afero.NewMemMapFs()

	// Create a mock project directory with .beads subdirectory
	projectDir := "/test-project"
	fs.MkdirAll(filepath.Join(projectDir, ".beads"), 0755)

	// Use BasePathFs to root the filesystem at the project directory
	baseFS := afero.NewBasePathFs(fs, projectDir)

	// Test the underlying BeadsAvailable function
	available, err := BeadsAvailable(baseFS)
	if err != nil {
		t.Fatalf("BeadsAvailable returned error: %v", err)
	}

	// In test environment, bd is unlikely to be on PATH
	// But the directory check should succeed
	t.Logf("Beads available in %s: %v", projectDir, available)
}

func TestNewProvider_BeadsDirectoryIsFile(t *testing.T) {
	// Create a temporary directory where .beads is a file, not a directory
	tempDir, err := os.MkdirTemp("", "provider-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create .beads as a file instead of a directory
	beadsFile := filepath.Join(tempDir, ".beads")
	if err := os.WriteFile(beadsFile, []byte("not a directory"), 0644); err != nil {
		t.Fatalf("Failed to create .beads file: %v", err)
	}

	// Call NewProvider - should handle the error gracefully and return BacklogManager
	provider := NewProvider(tempDir)

	// Verify provider is not nil (should fallback to BacklogManager)
	if provider == nil {
		t.Fatal("Expected non-nil provider (fallback to BacklogManager)")
	}
}
