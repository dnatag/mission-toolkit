// Package backlog provides epic management for Beads integration.
package backlog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Epic title constants for Beads provider.
const (
	// EpicTitleFeatures represents the Features epic title.
	EpicTitleFeatures = "Features"

	// EpicTitleBugfixes represents the Bugfixes epic title.
	EpicTitleBugfixes = "Bugfixes"

	// EpicTitleDecomposedIntents represents the Decomposed Intents epic title.
	EpicTitleDecomposedIntents = "Decomposed Intents"

	// EpicTitleRefactoringOpportunities represents the Refactoring Opportunities epic title.
	EpicTitleRefactoringOpportunities = "Refactoring Opportunities"

	// EpicTitleFutureEnhancements represents the Future Enhancements epic title.
	EpicTitleFutureEnhancements = "Future Enhancements"
)

// epicCache represents the structure for storing epic ID mappings.
type epicCache struct {
	Epics map[string]string `json:"epics"` // itemType -> epicID
}

// defineEpicTypes returns the mapping of item types to epic titles.
// This encapsulates the epic type definitions used by ensureEpics.
func defineEpicTypes() map[string]string {
	return map[string]string{
		ItemTypeFeature:    EpicTitleFeatures,
		ItemTypeBugfix:     EpicTitleBugfixes,
		ItemTypeDecomposed: EpicTitleDecomposedIntents,
		ItemTypeRefactor:   EpicTitleRefactoringOpportunities,
		ItemTypeFuture:     EpicTitleFutureEnhancements,
	}
}

// ensureEpics verifies that all required epics exist in Beads.
// It loads the epic cache if it exists, otherwise creates the epics via bd CLI.
//
// This method is part of BeadsProvider and operates on the provider's epicCache
// and commandRunner fields.
func (p *BeadsProvider) ensureEpics() error {
	// Try to load existing cache
	if err := p.loadEpicCache(); err == nil {
		// Cache loaded successfully, epics already exist
		return nil
	}

	// Cache doesn't exist or is invalid, create epics
	epicTypes := defineEpicTypes()

	for itemType, epicTitle := range epicTypes {
		output, err := p.commandRunner.Run("create", epicTitle, "-t", "epic", "--json")
		if err != nil {
			return fmt.Errorf("failed to create %s epic (%s): %w", itemType, epicTitle, err)
		}

		// Parse the JSON output to extract the epic ID
		var result struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal([]byte(output), &result); err != nil {
			return fmt.Errorf("failed to parse epic creation output for %s (%s): %w", itemType, epicTitle, err)
		}

		p.epicCache.Epics[itemType] = result.ID
	}

	// Save the cache
	return p.saveEpicCache()
}

// loadEpicCache loads the epic ID mappings from the cache file.
// This method is part of BeadsProvider and operates on the provider's cachePath.
func (p *BeadsProvider) loadEpicCache() error {
	data, err := os.ReadFile(p.cachePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, p.epicCache)
}

// saveEpicCache saves the epic ID mappings to the cache file.
// This method is part of BeadsProvider and operates on the provider's cachePath.
func (p *BeadsProvider) saveEpicCache() error {
	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(p.cachePath), 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	data, err := json.MarshalIndent(p.epicCache, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal epic cache: %w", err)
	}

	return os.WriteFile(p.cachePath, data, 0644)
}
