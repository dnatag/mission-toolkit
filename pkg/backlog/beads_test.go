package backlog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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

func TestBeadsProviderList(t *testing.T) {
	tempDir := t.TempDir()

	// Create mock command runner
	mockRunner := newMockCommandRunner()

	// Mock list responses for different epic types
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-feature-1", "--json"},
		`[{"id":"bd-f1-1","title":"Implement user auth","status":"open"},{"id":"bd-f1-2","title":"Add login page","status":"closed"}]`,
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-bugfix-1", "--json"},
		`[{"id":"bd-b1-1","title":"Fix null pointer","status":"open"}]`,
		nil,
	)

	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	// Create epic cache file with pre-populated data
	epicCacheData := epicCache{
		Epics: map[string]string{
			"feature": "bd-feature-1",
			"bugfix":  "bd-bugfix-1",
		},
	}
	var data []byte
	var err error
	data, err = json.MarshalIndent(epicCacheData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal epic cache: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("Failed to write cache file: %v", err)
	}

	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: mockRunner,
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     cachePath,
	}

	// Test List with no filters
	results, err := provider.List(nil, nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 items, got %d", len(results))
	}

	// Check formatting
	expectedResults := []string{
		"- [ ] Implement user auth",
		"- [x] Add login page",
		"- [ ] Fix null pointer",
	}
	for i, expected := range expectedResults {
		if i >= len(results) {
			t.Errorf("Missing result %d: %s", i, expected)
			continue
		}
		if results[i] != expected {
			t.Errorf("Result %d mismatch:\nExpected: %s\nGot: %s", i, expected, results[i])
		}
	}
}

func TestBeadsProviderListWithFilters(t *testing.T) {
	tempDir := t.TempDir()

	// Create mock command runner
	mockRunner := newMockCommandRunner()

	// Mock list responses
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-feature-1", "--json"},
		`[{"id":"bd-f1-1","title":"Feature 1","status":"open"}]`,
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-bugfix-1", "--json"},
		`[{"id":"bd-b1-1","title":"Bugfix 1","status":"open"}]`,
		nil,
	)

	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	// Create epic cache file with pre-populated data
	epicCacheData := epicCache{
		Epics: map[string]string{
			"feature": "bd-feature-1",
			"bugfix":  "bd-bugfix-1",
		},
	}
	var data []byte
	var err error
	data, err = json.MarshalIndent(epicCacheData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal epic cache: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("Failed to write cache file: %v", err)
	}

	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: mockRunner,
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     cachePath,
	}

	// Test with include filter (only feature)
	results, err := provider.List([]string{"feature"}, nil)
	if err != nil {
		t.Fatalf("List with include failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 item with feature filter, got %d", len(results))
	}
	if results[0] != "- [ ] Feature 1" {
		t.Errorf("Expected '- [ ] Feature 1', got '%s'", results[0])
	}

	// Test with exclude filter (exclude feature)
	results, err = provider.List(nil, []string{"feature"})
	if err != nil {
		t.Fatalf("List with exclude failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 item with feature excluded, got %d", len(results))
	}
	if results[0] != "- [ ] Bugfix 1" {
		t.Errorf("Expected '- [ ] Bugfix 1', got '%s'", results[0])
	}
}

func TestBeadsProviderAdd(t *testing.T) {
	tempDir := t.TempDir()

	// Create mock command runner
	mockRunner := newMockCommandRunner()

	// Mock create response
	mockRunner.setResponse(
		[]string{"create", "New feature task", "-t", "task", "--parent", "bd-feature-1", "--json"},
		`{"id":"bd-task-1","title":"New feature task"}`,
		nil,
	)

	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	// Create epic cache file with pre-populated data
	epicCacheData := epicCache{
		Epics: map[string]string{
			"feature": "bd-feature-1",
		},
	}
	var data []byte
	var err error
	data, err = json.MarshalIndent(epicCacheData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal epic cache: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("Failed to write cache file: %v", err)
	}

	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: mockRunner,
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     cachePath,
	}

	// Test Add
	err = provider.Add("New feature task", "feature")
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}
}

func TestBeadsProviderAddInvalidType(t *testing.T) {
	tempDir := t.TempDir()

	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	// Create epic cache file with pre-populated data
	epicCacheData := epicCache{
		Epics: map[string]string{
			"feature": "bd-feature-1",
		},
	}
	var data []byte
	var err error
	data, err = json.MarshalIndent(epicCacheData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal epic cache: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("Failed to write cache file: %v", err)
	}

	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: newMockCommandRunner(),
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     cachePath,
	}

	// Test Add with invalid type
	err = provider.Add("Test task", "invalid")
	if err == nil {
		t.Fatal("Expected error for invalid type, got nil")
	}

	expectedErrMsg := "invalid type: invalid"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

func TestBeadsProviderAddWithPattern(t *testing.T) {
	tempDir := t.TempDir()

	// Create mock command runner
	mockRunner := newMockCommandRunner()

	// Mock create response
	mockRunner.setResponse(
		[]string{"create", "Refactor authentication", "-t", "task", "--parent", "bd-refactor-1", "--json"},
		`{"id":"bd-task-1","title":"Refactor authentication"}`,
		nil,
	)

	// Mock list response for finding the task
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-refactor-1", "--json"},
		`[{"id":"bd-task-1","title":"Refactor authentication","notes":""}]`,
		nil,
	)

	// Mock update response
	mockRunner.setResponse(
		[]string{"update", "bd-task-1", "--notes", "PATTERN:auth-refactor COUNT:1"},
		``,
		nil,
	)

	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	// Create epic cache file with pre-populated data
	epicCacheData := epicCache{
		Epics: map[string]string{
			"refactor": "bd-refactor-1",
			"feature":  "bd-feature-1",
		},
	}
	var data []byte
	var err error
	data, err = json.MarshalIndent(epicCacheData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal epic cache: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("Failed to write cache file: %v", err)
	}

	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: mockRunner,
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     cachePath,
	}

	// Test AddWithPattern for refactor type
	err = provider.AddWithPattern("Refactor authentication", "refactor", "auth-refactor")
	if err != nil {
		t.Fatalf("AddWithPattern failed: %v", err)
	}

	// Test AddWithPattern for non-refactor type (should not add pattern tracking)
	mockRunner.setResponse(
		[]string{"create", "New feature", "-t", "task", "--parent", "bd-feature-1", "--json"},
		`{"id":"bd-task-2","title":"New feature"}`,
		nil,
	)

	err = provider.AddWithPattern("New feature", "feature", "some-pattern")
	if err != nil {
		t.Fatalf("AddWithPattern for non-refactor failed: %v", err)
	}
}

func TestBeadsProviderAddMultiple(t *testing.T) {
	tempDir := t.TempDir()

	// Create mock command runner
	mockRunner := newMockCommandRunner()

	descriptions := []string{"Task 1", "Task 2", "Task 3"}

	// Mock create responses for each task
	for i, desc := range descriptions {
		taskID := fmt.Sprintf("bd-task-%d", i+1)
		mockRunner.setResponse(
			[]string{"create", desc, "-t", "task", "--parent", "bd-feature-1", "--json"},
			fmt.Sprintf(`{"id":"%s","title":"%s"}`, taskID, desc),
			nil,
		)
	}

	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	// Create epic cache file with pre-populated data
	epicCacheData := epicCache{
		Epics: map[string]string{
			"feature": "bd-feature-1",
		},
	}
	var data []byte
	var err error
	data, err = json.MarshalIndent(epicCacheData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal epic cache: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("Failed to write cache file: %v", err)
	}

	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: mockRunner,
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     cachePath,
	}

	// Test AddMultiple
	err = provider.AddMultiple(descriptions, "feature")
	if err != nil {
		t.Fatalf("AddMultiple failed: %v", err)
	}

	// Test with empty slice
	err = provider.AddMultiple([]string{}, "feature")
	if err != nil {
		t.Fatalf("AddMultiple with empty slice failed: %v", err)
	}

	// Test with invalid type
	err = provider.AddMultiple([]string{"Task 1"}, "invalid")
	if err == nil {
		t.Fatal("Expected error for invalid type, got nil")
	}
}

func TestBeadsProviderGetPatternCount(t *testing.T) {
	tempDir := t.TempDir()

	// Create mock command runner
	mockRunner := newMockCommandRunner()

	// Mock list response with pattern notes
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-refactor-1", "--json"},
		`[
			{"id":"bd-r1","title":"Refactor 1","notes":"PATTERN:auth COUNT:1"},
			{"id":"bd-r2","title":"Refactor 2","notes":"PATTERN:auth COUNT:2"},
			{"id":"bd-r3","title":"Refactor 3","notes":"PATTERN:db COUNT:1"}
		]`,
		nil,
	)

	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	// Create epic cache file with pre-populated data
	epicCacheData := epicCache{
		Epics: map[string]string{
			"refactor": "bd-refactor-1",
		},
	}
	var data []byte
	var err error
	data, err = json.MarshalIndent(epicCacheData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal epic cache: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("Failed to write cache file: %v", err)
	}

	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: mockRunner,
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     cachePath,
	}

	// Test GetPatternCount for existing pattern
	count, err := provider.GetPatternCount("auth")
	if err != nil {
		t.Fatalf("GetPatternCount failed: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected pattern count 2 for 'auth', got %d", count)
	}

	// Test GetPatternCount for pattern with single occurrence
	count, err = provider.GetPatternCount("db")
	if err != nil {
		t.Fatalf("GetPatternCount for 'db' failed: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected pattern count 1 for 'db', got %d", count)
	}

	// Test GetPatternCount for non-existent pattern
	count, err = provider.GetPatternCount("nonexistent")
	if err != nil {
		t.Fatalf("GetPatternCount for non-existent failed: %v", err)
	}

	if count != 0 {
		t.Errorf("Expected pattern count 0 for non-existent, got %d", count)
	}
}

func TestBeadsProviderComplete(t *testing.T) {
	tempDir := t.TempDir()

	// Create mock command runner
	mockRunner := newMockCommandRunner()

	// Mock list responses for each epic type
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-feature-1", "--json"},
		`[{"id":"bd-task-1","title":"Feature task 1"},{"id":"bd-task-2","title":"Task to complete"}]`,
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-bugfix-1", "--json"},
		`[{"id":"bd-task-3","title":"Bugfix task"}]`,
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-decomposed-1", "--json"},
		`[]`,
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-refactor-1", "--json"},
		`[]`,
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-future-1", "--json"},
		`[]`,
		nil,
	)

	// Mock close response
	mockRunner.setResponse(
		[]string{"close", "bd-task-2", "--reason", "Completed"},
		``,
		nil,
	)

	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	// Create epic cache file with pre-populated data
	epicCacheData := epicCache{
		Epics: map[string]string{
			"feature":    "bd-feature-1",
			"bugfix":     "bd-bugfix-1",
			"decomposed": "bd-decomposed-1",
			"refactor":   "bd-refactor-1",
			"future":     "bd-future-1",
		},
	}
	data, err := json.MarshalIndent(epicCacheData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal epic cache: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("Failed to write cache file: %v", err)
	}

	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: mockRunner,
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     cachePath,
	}

	// Test Complete
	err = provider.Complete("Task to complete")
	if err != nil {
		t.Fatalf("Complete failed: %v", err)
	}
}

func TestBeadsProviderCompleteNotFound(t *testing.T) {
	tempDir := t.TempDir()

	// Create mock command runner
	mockRunner := newMockCommandRunner()

	// Mock list responses returning no matching task
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-feature-1", "--json"},
		`[{"id":"bd-task-1","title":"Other task"}]`,
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-bugfix-1", "--json"},
		`[]`,
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-decomposed-1", "--json"},
		`[]`,
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-refactor-1", "--json"},
		`[]`,
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-future-1", "--json"},
		`[]`,
		nil,
	)

	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	// Create epic cache file with pre-populated data
	epicCacheData := epicCache{
		Epics: map[string]string{
			"feature":    "bd-feature-1",
			"bugfix":     "bd-bugfix-1",
			"decomposed": "bd-decomposed-1",
			"refactor":   "bd-refactor-1",
			"future":     "bd-future-1",
		},
	}
	data, err := json.MarshalIndent(epicCacheData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal epic cache: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("Failed to write cache file: %v", err)
	}

	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: mockRunner,
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     cachePath,
	}

	// Test Complete with non-existent item
	err = provider.Complete("Non-existent task")
	if err == nil {
		t.Fatal("Expected error for non-existent task, got nil")
	}
	expectedErr := "item not found: Non-existent task"
	if err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestBeadsProviderCleanupAll(t *testing.T) {
	tempDir := t.TempDir()

	// Create mock command runner
	mockRunner := newMockCommandRunner()

	// Mock list responses with mixed open and closed items
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-feature-1", "--json"},
		`[{"id":"bd-task-1","title":"Open feature","status":"open"},{"id":"bd-task-2","title":"Closed feature","status":"closed"}]`,
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-bugfix-1", "--json"},
		`[{"id":"bd-task-3","title":"Closed bugfix","status":"closed"}]`,
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-decomposed-1", "--json"},
		`[{"id":"bd-task-4","title":"Open decomposed","status":"open"}]`,
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-refactor-1", "--json"},
		`[]`,
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-future-1", "--json"},
		`[]`,
		nil,
	)

	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	// Create epic cache file with pre-populated data
	epicCacheData := epicCache{
		Epics: map[string]string{
			"feature":    "bd-feature-1",
			"bugfix":     "bd-bugfix-1",
			"decomposed": "bd-decomposed-1",
			"refactor":   "bd-refactor-1",
			"future":     "bd-future-1",
		},
	}
	data, err := json.MarshalIndent(epicCacheData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal epic cache: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("Failed to write cache file: %v", err)
	}

	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: mockRunner,
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     cachePath,
	}

	// Test Cleanup with no itemType (all types)
	count, err := provider.Cleanup("")
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	// Should count 2 closed items (1 feature + 1 bugfix)
	if count != 2 {
		t.Errorf("Expected 2 closed items, got %d", count)
	}
}

func TestBeadsProviderCleanupByType(t *testing.T) {
	tempDir := t.TempDir()

	// Create mock command runner
	mockRunner := newMockCommandRunner()

	// Mock list response for feature epic only
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-feature-1", "--json"},
		`[{"id":"bd-task-1","title":"Open feature","status":"open"},{"id":"bd-task-2","title":"Closed feature 1","status":"closed"},{"id":"bd-task-3","title":"Closed feature 2","status":"closed"}]`,
		nil,
	)

	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	// Create epic cache file with pre-populated data
	epicCacheData := epicCache{
		Epics: map[string]string{
			"feature":    "bd-feature-1",
			"bugfix":     "bd-bugfix-1",
			"decomposed": "bd-decomposed-1",
			"refactor":   "bd-refactor-1",
			"future":     "bd-future-1",
		},
	}
	data, err := json.MarshalIndent(epicCacheData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal epic cache: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("Failed to write cache file: %v", err)
	}

	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: mockRunner,
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     cachePath,
	}

	// Test Cleanup with itemType specified
	count, err := provider.Cleanup("feature")
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	// Should count 2 closed items in feature epic
	if count != 2 {
		t.Errorf("Expected 2 closed feature items, got %d", count)
	}
}

func TestBeadsProviderCleanupInvalidType(t *testing.T) {
	tempDir := t.TempDir()

	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	// Create epic cache file with pre-populated data
	epicCacheData := epicCache{
		Epics: map[string]string{
			"feature": "bd-feature-1",
		},
	}
	data, err := json.MarshalIndent(epicCacheData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal epic cache: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("Failed to write cache file: %v", err)
	}

	mockRunner := newMockCommandRunner()

	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: mockRunner,
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     cachePath,
	}

	// Test Cleanup with invalid itemType
	count, err := provider.Cleanup("invalid")
	if err == nil {
		t.Fatal("Expected error for invalid type, got nil")
	}
	if count != 0 {
		t.Errorf("Expected count 0 for invalid type, got %d", count)
	}
}

func TestBeadsProviderDecompose(t *testing.T) {
	tempDir := t.TempDir()

	// Create mock command runner
	mockRunner := newMockCommandRunner()

	// Mock create responses for sub-intents
	mockRunner.setResponse(
		[]string{"create", "Implement authentication", "-t", "task", "--parent", "bd-decomposed-1", "--json"},
		`{"id":"bd-task-1"}`,
		nil,
	)
	mockRunner.setResponse(
		[]string{"create", "Add user registration", "-t", "task", "--parent", "bd-decomposed-1", "--json"},
		`{"id":"bd-task-2"}`,
		nil,
	)
	mockRunner.setResponse(
		[]string{"create", "Add login form", "-t", "task", "--parent", "bd-decomposed-1", "--json"},
		`{"id":"bd-task-3"}`,
		nil,
	)

	// Mock add-dep responses (task 2 depends on task 1)
	mockRunner.setResponse(
		[]string{"add-dep", "bd-task-2", "bd-task-1"},
		``,
		nil,
	)
	// Mock add-dep responses (task 3 depends on task 2)
	mockRunner.setResponse(
		[]string{"add-dep", "bd-task-3", "bd-task-2"},
		``,
		nil,
	)

	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	// Create epic cache file with pre-populated data
	epicCacheData := epicCache{
		Epics: map[string]string{
			"decomposed": "bd-decomposed-1",
		},
	}
	data, err := json.MarshalIndent(epicCacheData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal epic cache: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("Failed to write cache file: %v", err)
	}

	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: mockRunner,
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     cachePath,
	}

	// Test Decompose with valid JSON and dependencies
	jsonInput := `{
		"action": "decompose",
		"sub_intents": [
			{
				"intent": "Implement authentication",
				"rationale": "Core authentication system",
				"estimated_files": 3,
				"dependencies": []
			},
			{
				"intent": "Add user registration",
				"rationale": "User registration flow",
				"estimated_files": 2,
				"dependencies": ["Implement authentication"]
			},
			{
				"intent": "Add login form",
				"rationale": "Login UI",
				"estimated_files": 1,
				"dependencies": ["Add user registration"]
			}
		],
		"decomposition_rationale": "Breaking down authentication feature"
	}`

	err = provider.Decompose(jsonInput)
	if err != nil {
		t.Fatalf("Decompose failed: %v", err)
	}
}

func TestBeadsProviderDecomposeInvalidJSON(t *testing.T) {
	tempDir := t.TempDir()

	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	// Create epic cache file with pre-populated data
	epicCacheData := epicCache{
		Epics: map[string]string{
			"decomposed": "bd-decomposed-1",
		},
	}
	data, err := json.MarshalIndent(epicCacheData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal epic cache: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("Failed to write cache file: %v", err)
	}

	mockRunner := newMockCommandRunner()

	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: mockRunner,
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     cachePath,
	}

	// Test Decompose with invalid JSON
	err = provider.Decompose("invalid json")
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}
	if !strings.Contains(err.Error(), "parsing decompose JSON") {
		t.Errorf("Expected JSON parsing error, got: %v", err)
	}
}

func TestBeadsProviderDecomposeEmptySubIntents(t *testing.T) {
	tempDir := t.TempDir()

	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	// Create epic cache file with pre-populated data
	epicCacheData := epicCache{
		Epics: map[string]string{
			"decomposed": "bd-decomposed-1",
		},
	}
	data, err := json.MarshalIndent(epicCacheData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal epic cache: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("Failed to write cache file: %v", err)
	}

	mockRunner := newMockCommandRunner()

	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: mockRunner,
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     cachePath,
	}

	// Test Decompose with empty sub_intents
	jsonInput := `{
		"action": "decompose",
		"sub_intents": [],
		"decomposition_rationale": "Test"
	}`

	err = provider.Decompose(jsonInput)
	if err == nil {
		t.Fatal("Expected error for empty sub_intents, got nil")
	}
	if !strings.Contains(err.Error(), "no sub-intents found") {
		t.Errorf("Expected 'no sub-intents found' error, got: %v", err)
	}
}

func TestBeadsProviderDecomposeMissingDependency(t *testing.T) {
	tempDir := t.TempDir()

	// Create mock command runner
	mockRunner := newMockCommandRunner()

	// Mock create responses
	mockRunner.setResponse(
		[]string{"create", "Task A", "-t", "task", "--parent", "bd-decomposed-1", "--json"},
		`{"id":"bd-task-1"}`,
		nil,
	)
	mockRunner.setResponse(
		[]string{"create", "Task B", "-t", "task", "--parent", "bd-decomposed-1", "--json"},
		`{"id":"bd-task-2"}`,
		nil,
	)

	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	// Create epic cache file with pre-populated data
	epicCacheData := epicCache{
		Epics: map[string]string{
			"decomposed": "bd-decomposed-1",
		},
	}
	data, err := json.MarshalIndent(epicCacheData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal epic cache: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("Failed to write cache file: %v", err)
	}

	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: mockRunner,
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     cachePath,
	}

	// Test Decompose with missing dependency
	jsonInput := `{
		"action": "decompose",
		"sub_intents": [
			{
				"intent": "Task A",
				"rationale": "First task",
				"estimated_files": 1,
				"dependencies": []
			},
			{
				"intent": "Task B",
				"rationale": "Second task",
				"estimated_files": 1,
				"dependencies": ["Non-existent Task"]
			}
		],
		"decomposition_rationale": "Test"
	}`

	err = provider.Decompose(jsonInput)
	if err == nil {
		t.Fatal("Expected error for missing dependency, got nil")
	}
	if !strings.Contains(err.Error(), "dependency task not found") {
		t.Errorf("Expected 'dependency task not found' error, got: %v", err)
	}
}

func TestBeadsProviderDecomposeCLICreateFailure(t *testing.T) {
	tempDir := t.TempDir()

	// Create mock command runner
	mockRunner := newMockCommandRunner()

	// Mock create failure
	mockRunner.setResponse(
		[]string{"create", "Task A", "-t", "task", "--parent", "bd-decomposed-1", "--json"},
		``,
		fmt.Errorf("bd CLI failed"),
	)

	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	// Create epic cache file with pre-populated data
	epicCacheData := epicCache{
		Epics: map[string]string{
			"decomposed": "bd-decomposed-1",
		},
	}
	data, err := json.MarshalIndent(epicCacheData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal epic cache: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("Failed to write cache file: %v", err)
	}

	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: mockRunner,
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     cachePath,
	}

	// Test Decompose with CLI create failure
	jsonInput := `{
		"action": "decompose",
		"sub_intents": [
			{
				"intent": "Task A",
				"rationale": "First task",
				"estimated_files": 1,
				"dependencies": []
			}
		],
		"decomposition_rationale": "Test"
	}`

	err = provider.Decompose(jsonInput)
	if err == nil {
		t.Fatal("Expected error for CLI create failure, got nil")
	}
	if !strings.Contains(err.Error(), "failed to create task") {
		t.Errorf("Expected 'failed to create task' error, got: %v", err)
	}
}

func TestBeadsProviderDecomposeAddDepFailure(t *testing.T) {
	tempDir := t.TempDir()

	// Create mock command runner
	mockRunner := newMockCommandRunner()

	// Mock create responses
	mockRunner.setResponse(
		[]string{"create", "Task A", "-t", "task", "--parent", "bd-decomposed-1", "--json"},
		`{"id":"bd-task-1"}`,
		nil,
	)
	mockRunner.setResponse(
		[]string{"create", "Task B", "-t", "task", "--parent", "bd-decomposed-1", "--json"},
		`{"id":"bd-task-2"}`,
		nil,
	)

	// Mock add-dep failure
	mockRunner.setResponse(
		[]string{"add-dep", "bd-task-2", "bd-task-1"},
		``,
		fmt.Errorf("bd add-dep failed"),
	)

	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	// Create epic cache file with pre-populated data
	epicCacheData := epicCache{
		Epics: map[string]string{
			"decomposed": "bd-decomposed-1",
		},
	}
	data, err := json.MarshalIndent(epicCacheData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal epic cache: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("Failed to write cache file: %v", err)
	}

	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: mockRunner,
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     cachePath,
	}

	// Test Decompose with add-dep failure
	jsonInput := `{
		"action": "decompose",
		"sub_intents": [
			{
				"intent": "Task A",
				"rationale": "First task",
				"estimated_files": 1,
				"dependencies": []
			},
			{
				"intent": "Task B",
				"rationale": "Second task",
				"estimated_files": 1,
				"dependencies": ["Task A"]
			}
		],
		"decomposition_rationale": "Test"
	}`

	err = provider.Decompose(jsonInput)
	if err == nil {
		t.Fatal("Expected error for add-dep failure, got nil")
	}
	if !strings.Contains(err.Error(), "failed to add dependency") {
		t.Errorf("Expected 'failed to add dependency' error, got: %v", err)
	}
}

func TestGenerateSnapshot_Success(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create mock command runner with epic and item responses
	mockRunner := newMockCommandRunner()

	// Setup epic creation responses
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

	// Setup list responses for each epic type
	featureItems := []map[string]interface{}{
		{"id": "task-1", "title": "Add new feature", "status": "open"},
		{"id": "task-2", "title": "Completed feature", "status": "closed"},
	}
	featureJSON, _ := json.Marshal(featureItems)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-feature-1", "--json"},
		string(featureJSON),
		nil,
	)

	bugfixItems := []map[string]interface{}{
		{"id": "task-3", "title": "Fix login bug", "status": "open"},
	}
	bugfixJSON, _ := json.Marshal(bugfixItems)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-bugfix-1", "--json"},
		string(bugfixJSON),
		nil,
	)

	// Setup empty responses for other types
	emptyJSON, _ := json.Marshal([]interface{}{})
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-decomposed-1", "--json"},
		string(emptyJSON),
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-refactor-1", "--json"},
		string(emptyJSON),
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-future-1", "--json"},
		string(emptyJSON),
		nil,
	)

	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: mockRunner,
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     filepath.Join(tempDir, ".mission", "beads-epics.json"),
	}

	// Ensure epics are loaded
	if err := provider.ensureEpics(); err != nil {
		t.Fatalf("ensureEpics failed: %v", err)
	}

	// Generate snapshot
	provider.generateSnapshot()

	// Verify snapshot file was created
	backlogPath := filepath.Join(tempDir, ".mission", "backlog.md")
	if _, err := os.Stat(backlogPath); os.IsNotExist(err) {
		t.Fatal("Snapshot file was not created")
	}

	// Read and verify snapshot content
	content, err := os.ReadFile(backlogPath)
	if err != nil {
		t.Fatalf("Failed to read snapshot: %v", err)
	}

	snapshot := string(content)

	// Verify frontmatter
	if !strings.Contains(snapshot, "---") {
		t.Error("Snapshot missing frontmatter")
	}
	if !strings.Contains(snapshot, "source: beads") {
		t.Error("Snapshot missing source: beads in frontmatter")
	}

	// Verify title
	if !strings.Contains(snapshot, "# Backlog") {
		t.Error("Snapshot missing title")
	}

	// Verify sections
	expectedSections := []string{
		"## FEATURES",
		"## BUGFIXES",
		"## DECOMPOSED INTENTS",
		"## REFACTORING OPPORTUNITIES",
		"## FUTURE ENHANCEMENTS",
		"## COMPLETED",
	}
	for _, section := range expectedSections {
		if !strings.Contains(snapshot, section) {
			t.Errorf("Snapshot missing section: %s", section)
		}
	}

	// Verify items
	if !strings.Contains(snapshot, "Add new feature") {
		t.Error("Snapshot missing feature item")
	}
	if !strings.Contains(snapshot, "Fix login bug") {
		t.Error("Snapshot missing bugfix item")
	}
}

func TestGenerateSnapshot_ListError(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create mock command runner
	mockRunner := newMockCommandRunner()

	// Setup epic creation responses
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

	// Setup list response that returns an error
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-feature-1", "--json"},
		"",
		fmt.Errorf("list error"),
	)

	// Setup empty responses for other types
	emptyJSON, _ := json.Marshal([]interface{}{})
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-bugfix-1", "--json"},
		string(emptyJSON),
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-decomposed-1", "--json"},
		string(emptyJSON),
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-refactor-1", "--json"},
		string(emptyJSON),
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-future-1", "--json"},
		string(emptyJSON),
		nil,
	)

	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: mockRunner,
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     filepath.Join(tempDir, ".mission", "beads-epics.json"),
	}

	// Ensure epics are loaded
	if err := provider.ensureEpics(); err != nil {
		t.Fatalf("ensureEpics failed: %v", err)
	}

	// Generate snapshot should not fail even with list error
	provider.generateSnapshot()

	// Verify snapshot file was still created
	backlogPath := filepath.Join(tempDir, ".mission", "backlog.md")
	if _, err := os.Stat(backlogPath); os.IsNotExist(err) {
		t.Fatal("Snapshot file should have been created even with list error")
	}
}

func TestGenerateSnapshot_WriteError(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create a read-only backlog.md file to cause a write error
	backlogPath := filepath.Join(tempDir, ".mission", "backlog.md")
	if err := os.MkdirAll(filepath.Dir(backlogPath), 0755); err != nil {
		t.Fatalf("Failed to create .mission directory: %v", err)
	}
	if err := os.WriteFile(backlogPath, []byte("existing content"), 0444); err != nil {
		t.Fatalf("Failed to create read-only file: %v", err)
	}

	// Create mock command runner
	mockRunner := newMockCommandRunner()

	// Setup epic creation responses
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

	// Setup empty responses
	emptyJSON, _ := json.Marshal([]interface{}{})
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-feature-1", "--json"},
		string(emptyJSON),
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-bugfix-1", "--json"},
		string(emptyJSON),
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-decomposed-1", "--json"},
		string(emptyJSON),
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-refactor-1", "--json"},
		string(emptyJSON),
		nil,
	)
	mockRunner.setResponse(
		[]string{"list", "--parent", "bd-future-1", "--json"},
		string(emptyJSON),
		nil,
	)

	provider := &BeadsProvider{
		projectRoot:   tempDir,
		commandRunner: mockRunner,
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     filepath.Join(tempDir, ".mission", "beads-epics.json"),
	}

	// Ensure epics are loaded
	if err := provider.ensureEpics(); err != nil {
		t.Fatalf("ensureEpics failed: %v", err)
	}

	// Generate snapshot should not panic with write error
	provider.generateSnapshot()
}
