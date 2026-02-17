// Package backlog provides tests for epic management.
package beads

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	backlog "github.com/dnatag/mission-toolkit/pkg/backlog"
)

func TestOrderedEpicTypes(t *testing.T) {
	epicTypes := OrderedEpicTypes()

	expected := []EpicTypeDef{
		{backlog.ItemTypeFeature, EpicTitleFeatures},
		{backlog.ItemTypeBugfix, EpicTitleBugfixes},
		{backlog.ItemTypeDecomposed, EpicTitleDecomposedIntents},
		{backlog.ItemTypeRefactor, EpicTitleRefactoringOpportunities},
		{backlog.ItemTypeFuture, EpicTitleFutureEnhancements},
	}

	if len(epicTypes) != len(expected) {
		t.Fatalf("Expected %d epic types, got %d", len(expected), len(epicTypes))
	}

	for i, et := range epicTypes {
		if et.ItemType != expected[i].ItemType || et.EpicTitle != expected[i].EpicTitle {
			t.Errorf("Index %d: expected {%s, %s}, got {%s, %s}", i, expected[i].ItemType, expected[i].EpicTitle, et.ItemType, et.EpicTitle)
		}
	}
}

func TestLoadEpicCache(t *testing.T) {
	tempDir := t.TempDir()
	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	// Create test epic cache data
	EpicCacheData := EpicCache{
		Epics: map[string]string{
			backlog.ItemTypeFeature: "bd-feature-1",
			backlog.ItemTypeBugfix:  "bd-bugfix-1",
		},
	}
	data, err := json.MarshalIndent(EpicCacheData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal epic cache: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		t.Fatalf("Failed to write cache file: %v", err)
	}

	provider := &Provider{
		projectRoot: tempDir,
		epicCache:   &EpicCache{Epics: make(map[string]string)},
		cachePath:   cachePath,
	}

	err = provider.loadEpicCache()
	if err != nil {
		t.Fatalf("loadEpicCache failed: %v", err)
	}

	if len(provider.epicCache.Epics) != 2 {
		t.Errorf("Expected 2 epics, got %d", len(provider.epicCache.Epics))
	}

	if provider.epicCache.Epics[backlog.ItemTypeFeature] != "bd-feature-1" {
		t.Errorf("Expected feature epic ID bd-feature-1, got %s", provider.epicCache.Epics[backlog.ItemTypeFeature])
	}
}

func TestLoadEpicCacheMissingFile(t *testing.T) {
	tempDir := t.TempDir()
	cachePath := filepath.Join(tempDir, ".mission", "beads-epics.json")

	provider := &Provider{
		projectRoot: tempDir,
		epicCache:   &EpicCache{Epics: make(map[string]string)},
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

	provider := &Provider{
		projectRoot: tempDir,
		epicCache: &EpicCache{
			Epics: map[string]string{
				backlog.ItemTypeFeature: "bd-feature-1",
				backlog.ItemTypeBugfix:  "bd-bugfix-1",
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
	var loadedCache EpicCache
	if err := json.Unmarshal(data, &loadedCache); err != nil {
		t.Fatalf("Failed to unmarshal cache: %v", err)
	}

	if len(loadedCache.Epics) != 2 {
		t.Errorf("Expected 2 epics in saved cache, got %d", len(loadedCache.Epics))
	}

	if loadedCache.Epics[backlog.ItemTypeFeature] != "bd-feature-1" {
		t.Errorf("Expected feature epic ID bd-feature-1, got %s", loadedCache.Epics[backlog.ItemTypeFeature])
	}
}
