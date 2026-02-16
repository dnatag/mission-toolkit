package backlog

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// mockCommandRunner implements CommandRunner for testing.
type mockCommandRunner struct {
	outputs map[string]string // args -> output
	errors  map[string]error  // args -> error
}

// newMockCommandRunner creates a new mockCommandRunner.
func newMockCommandRunner() *mockCommandRunner {
	return &mockCommandRunner{
		outputs: make(map[string]string),
		errors:  make(map[string]error),
	}
}

// setResponse sets the mock response for a given command.
func (m *mockCommandRunner) setResponse(args []string, output string, err error) {
	key := argsToString(args)
	m.outputs[key] = output
	m.errors[key] = err
}

// Run executes the mock command.
func (m *mockCommandRunner) Run(args ...string) (string, error) {
	key := argsToString(args)
	output, ok := m.outputs[key]
	if !ok {
		return "", fmt.Errorf("unexpected command: %s", key)
	}
	return output, m.errors[key]
}

// argsToString converts args to a string key for map lookup.
func argsToString(args []string) string {
	key := ""
	for i, arg := range args {
		if i > 0 {
			key += " "
		}
		key += arg
	}
	return key
}

func TestNewBeadsProvider(t *testing.T) {
	projectRoot := "/test/project"
	provider := NewBeadsProvider(projectRoot)

	if provider == nil {
		t.Fatal("NewBeadsProvider returned nil")
	}

	if provider.projectRoot != projectRoot {
		t.Errorf("expected projectRoot %s, got %s", projectRoot, provider.projectRoot)
	}

	if provider.commandRunner == nil {
		t.Error("expected commandRunner to be set")
	}

	if provider.epicCache == nil {
		t.Error("expected epicCache to be initialized")
	}

	if provider.epicCache.Epics == nil {
		t.Error("expected epicCache.Epics to be initialized")
	}

	expectedCachePath := filepath.Join(projectRoot, ".mission", "beads-epics.json")
	if provider.cachePath != expectedCachePath {
		t.Errorf("expected cachePath %s, got %s", expectedCachePath, provider.cachePath)
	}
}

func TestEnsureEpics_CreatesEpicsOnFirstCall(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create mock command runner that returns epic IDs
	mockRunner := newMockCommandRunner()
	mockRunner.setResponse(
		[]string{"create", "Features", "-t", "epic", "--json"},
		`{"id": "bd-feature-1"}`,
		nil,
	)
	mockRunner.setResponse(
		[]string{"create", "Bugfixes", "-t", "epic", "--json"},
		`{"id": "bd-bugfix-1"}`,
		nil,
	)
	mockRunner.setResponse(
		[]string{"create", "Decomposed Intents", "-t", "epic", "--json"},
		`{"id": "bd-decomposed-1"}`,
		nil,
	)
	mockRunner.setResponse(
		[]string{"create", "Refactoring Opportunities", "-t", "epic", "--json"},
		`{"id": "bd-refactor-1"}`,
		nil,
	)
	mockRunner.setResponse(
		[]string{"create", "Future Enhancements", "-t", "epic", "--json"},
		`{"id": "bd-future-1"}`,
		nil,
	)

	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: mockRunner,
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     filepath.Join(tempDir, ".mission", "beads-epics.json"),
	}

	// Ensure epics
	err := provider.ensureEpics()
	if err != nil {
		t.Fatalf("ensureEpics failed: %v", err)
	}

	// Verify all epic IDs were stored
	expectedEpics := map[string]string{
		"feature":    "bd-feature-1",
		"bugfix":     "bd-bugfix-1",
		"decomposed": "bd-decomposed-1",
		"refactor":   "bd-refactor-1",
		"future":     "bd-future-1",
	}

	for itemType, expectedID := range expectedEpics {
		actualID, ok := provider.epicCache.Epics[itemType]
		if !ok {
			t.Errorf("epic cache missing entry for %s", itemType)
			continue
		}
		if actualID != expectedID {
			t.Errorf("expected epic ID %s for %s, got %s", expectedID, itemType, actualID)
		}
	}

	// Verify cache file was created
	data, err := os.ReadFile(provider.cachePath)
	if err != nil {
		t.Fatalf("failed to read cache file: %v", err)
	}

	var cache epicCache
	if err := json.Unmarshal(data, &cache); err != nil {
		t.Fatalf("failed to parse cache file: %v", err)
	}

	for itemType, expectedID := range expectedEpics {
		actualID, ok := cache.Epics[itemType]
		if !ok {
			t.Errorf("cache file missing entry for %s", itemType)
			continue
		}
		if actualID != expectedID {
			t.Errorf("expected cache file epic ID %s for %s, got %s", expectedID, itemType, actualID)
		}
	}
}

func TestEnsureEpics_LoadsCacheOnSubsequentCalls(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create an existing cache file
	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")
	existingCache := epicCache{
		Epics: map[string]string{
			"feature":    "bd-feature-existing",
			"bugfix":     "bd-bugfix-existing",
			"decomposed": "bd-decomposed-existing",
			"refactor":   "bd-refactor-existing",
			"future":     "bd-future-existing",
		},
	}

	data, err := json.MarshalIndent(existingCache, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal cache: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		t.Fatalf("failed to create cache directory: %v", err)
	}

	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("failed to write cache file: %v", err)
	}

	// Create provider with mock runner (should not be called)
	mockRunner := newMockCommandRunner()
	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: mockRunner,
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     cachePath,
	}

	// Ensure epics should load from cache, not create new ones
	err = provider.ensureEpics()
	if err != nil {
		t.Fatalf("ensureEpics failed: %v", err)
	}

	// Verify existing epic IDs were loaded
	for itemType, expectedID := range existingCache.Epics {
		actualID, ok := provider.epicCache.Epics[itemType]
		if !ok {
			t.Errorf("epic cache missing entry for %s", itemType)
			continue
		}
		if actualID != expectedID {
			t.Errorf("expected epic ID %s for %s, got %s", expectedID, itemType, actualID)
		}
	}
}

func TestSaveEpicCache(t *testing.T) {
	tempDir := t.TempDir()
	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	provider := &BeadsProvider{
		projectRoot: tempDir,
		epicCache: &epicCache{
			Epics: map[string]string{
				"feature": "bd-test-1",
				"bugfix":  "bd-test-2",
			},
		},
		cachePath: cachePath,
	}

	err := provider.saveEpicCache()
	if err != nil {
		t.Fatalf("saveEpicCache failed: %v", err)
	}

	// Verify file was created
	data, err := os.ReadFile(cachePath)
	if err != nil {
		t.Fatalf("failed to read cache file: %v", err)
	}

	var cache epicCache
	if err := json.Unmarshal(data, &cache); err != nil {
		t.Fatalf("failed to parse cache file: %v", err)
	}

	if cache.Epics["feature"] != "bd-test-1" {
		t.Errorf("expected feature epic ID bd-test-1, got %s", cache.Epics["feature"])
	}

	if cache.Epics["bugfix"] != "bd-test-2" {
		t.Errorf("expected bugfix epic ID bd-test-2, got %s", cache.Epics["bugfix"])
	}
}

func TestLoadEpicCache(t *testing.T) {
	tempDir := t.TempDir()
	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	// Create a cache file
	testCache := epicCache{
		Epics: map[string]string{
			"feature": "bd-load-test-1",
			"bugfix":  "bd-load-test-2",
		},
	}

	data, err := json.MarshalIndent(testCache, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal cache: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		t.Fatalf("failed to create cache directory: %v", err)
	}

	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("failed to write cache file: %v", err)
	}

	provider := &BeadsProvider{
		projectRoot: tempDir,
		epicCache:   &epicCache{Epics: make(map[string]string)},
		cachePath:   cachePath,
	}

	err = provider.loadEpicCache()
	if err != nil {
		t.Fatalf("loadEpicCache failed: %v", err)
	}

	if provider.epicCache.Epics["feature"] != "bd-load-test-1" {
		t.Errorf("expected feature epic ID bd-load-test-1, got %s", provider.epicCache.Epics["feature"])
	}

	if provider.epicCache.Epics["bugfix"] != "bd-load-test-2" {
		t.Errorf("expected bugfix epic ID bd-load-test-2, got %s", provider.epicCache.Epics["bugfix"])
	}
}

func TestBdCommandRunner_Run(t *testing.T) {
	// This test requires the bd CLI to be installed
	// Skip if bd is not available
	if _, err := exec.LookPath("bd"); err != nil {
		t.Skip("bd CLI not installed, skipping integration test")
	}

	tempDir := t.TempDir()
	runner := NewBDCommandRunner(tempDir)

	// Test with a simple command like 'bd --help'
	output, err := runner.Run("--help")
	if err != nil {
		t.Logf("bd --help failed (expected if bd not initialized): %v", err)
	}

	if output == "" {
		t.Log("bd --help produced no output")
	}

	t.Logf("bd --help output: %s", output)
}
