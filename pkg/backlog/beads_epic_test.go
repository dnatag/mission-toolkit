// Package backlog provides tests for epic management.
package backlog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDefineEpicTypes(t *testing.T) {
	epicTypes := defineEpicTypes()

	expectedTypes := []string{ItemTypeFeature, ItemTypeBugfix, ItemTypeDecomposed, ItemTypeRefactor, ItemTypeFuture}
	expectedTitles := []string{EpicTitleFeatures, EpicTitleBugfixes, EpicTitleDecomposedIntents, EpicTitleRefactoringOpportunities, EpicTitleFutureEnhancements}

	if len(epicTypes) != len(expectedTypes) {
		t.Errorf("Expected %d epic types, got %d", len(expectedTypes), len(epicTypes))
	}

	for i, itemType := range expectedTypes {
		title, ok := epicTypes[itemType]
		if !ok {
			t.Errorf("Missing epic type: %s", itemType)
			continue
		}
		if title != expectedTitles[i] {
			t.Errorf("Epic type %s: expected title '%s', got '%s'", itemType, expectedTitles[i], title)
		}
	}
}

func TestLoadEpicCache(t *testing.T) {
	tempDir := t.TempDir()
	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	// Create test epic cache data
	epicCacheData := epicCache{
		Epics: map[string]string{
			ItemTypeFeature: "bd-feature-1",
			ItemTypeBugfix:  "bd-bugfix-1",
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
		projectRoot: tempDir,
		epicCache:   &epicCache{Epics: make(map[string]string)},
		cachePath:   cachePath,
	}

	err = provider.loadEpicCache()
	if err != nil {
		t.Fatalf("loadEpicCache failed: %v", err)
	}

	if len(provider.epicCache.Epics) != 2 {
		t.Errorf("Expected 2 epics, got %d", len(provider.epicCache.Epics))
	}

	if provider.epicCache.Epics[ItemTypeFeature] != "bd-feature-1" {
		t.Errorf("Expected feature epic ID bd-feature-1, got %s", provider.epicCache.Epics[ItemTypeFeature])
	}
}

func TestLoadEpicCacheMissingFile(t *testing.T) {
	tempDir := t.TempDir()
	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	provider := &BeadsProvider{
		projectRoot: tempDir,
		epicCache:   &epicCache{Epics: make(map[string]string)},
		cachePath:   cachePath,
	}

	err := provider.loadEpicCache()
	if err == nil {
		t.Fatal("Expected error for missing cache file, got nil")
	}
}

func TestSaveEpicCache(t *testing.T) {
	tempDir := t.TempDir()
	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	provider := &BeadsProvider{
		projectRoot: tempDir,
		epicCache: &epicCache{
			Epics: map[string]string{
				ItemTypeFeature: "bd-feature-1",
				ItemTypeBugfix:  "bd-bugfix-1",
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
		t.Fatalf("Failed to read cache file: %v", err)
	}

	// Verify content
	var loadedCache epicCache
	if err := json.Unmarshal(data, &loadedCache); err != nil {
		t.Fatalf("Failed to unmarshal cache: %v", err)
	}

	if len(loadedCache.Epics) != 2 {
		t.Errorf("Expected 2 epics in saved cache, got %d", len(loadedCache.Epics))
	}

	if loadedCache.Epics[ItemTypeFeature] != "bd-feature-1" {
		t.Errorf("Expected feature epic ID bd-feature-1, got %s", loadedCache.Epics[ItemTypeFeature])
	}
}
