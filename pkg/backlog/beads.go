// Package backlog provides backlog management functionality with support for
// multiple backend implementations. The BeadsProvider implementation uses the
// Beads (bd) CLI for dependency-aware task management.
package backlog

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Epic type constants for Beads provider.
const (
	EpicTypeFeature    = "feature"
	EpicTypeBugfix     = "bugfix"
	EpicTypeDecomposed = "decomposed"
	EpicTypeRefactor   = "refactor"
	EpicTypeFuture     = "future"
)

// Epic title constants for Beads provider.
const (
	EpicTitleFeatures                 = "Features"
	EpicTitleBugfixes                 = "Bugfixes"
	EpicTitleDecomposedIntents        = "Decomposed Intents"
	EpicTitleRefactoringOpportunities = "Refactoring Opportunities"
	EpicTitleFutureEnhancements       = "Future Enhancements"
)

// CommandRunner defines the interface for executing external commands.
// This abstraction enables testability by allowing mock implementations.
type CommandRunner interface {
	Run(args ...string) (string, error)
}

// bdCommandRunner implements CommandRunner for the bd CLI.
type bdCommandRunner struct {
	workDir string
}

// NewBDCommandRunner creates a new bdCommandRunner for the given working directory.
func NewBDCommandRunner(workDir string) CommandRunner {
	return &bdCommandRunner{workDir: workDir}
}

// Run executes a bd command with the provided arguments.
// It returns the combined stdout/stderr output and any error encountered.
func (r *bdCommandRunner) Run(args ...string) (string, error) {
	cmd := exec.Command("bd", args...)
	cmd.Dir = r.workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("bd %v failed: %w", args, err)
	}
	return string(output), nil
}

// epicCache represents the structure for storing epic ID mappings.
type epicCache struct {
	Epics map[string]string `json:"epics"` // itemType -> epicID
}

// BeadsProvider implements BacklogProvider using the bd CLI.
// It manages backlog items through Beads' dependency-aware task system.
type BeadsProvider struct {
	projectRoot   string
	commandRunner CommandRunner
	epicCache     *epicCache
	cachePath     string
}

// NewBeadsProvider creates a new BeadsProvider for the given project directory.
// It initializes the command runner and epic cache, and ensures that the required
// epics exist in Beads.
func NewBeadsProvider(projectRoot string) *BeadsProvider {
	cachePath := filepath.Join(projectRoot, ".mission", "beads-epics.json")
	return &BeadsProvider{
		projectRoot:   projectRoot,
		commandRunner: NewBDCommandRunner(projectRoot),
		epicCache:     &epicCache{Epics: make(map[string]string)},
		cachePath:     cachePath,
	}
}

// defineEpicTypes returns the mapping of item types to epic titles.
// This encapsulates the epic type definitions used by ensureEpics.
func defineEpicTypes() map[string]string {
	return map[string]string{
		EpicTypeFeature:    EpicTitleFeatures,
		EpicTypeBugfix:     EpicTitleBugfixes,
		EpicTypeDecomposed: EpicTitleDecomposedIntents,
		EpicTypeRefactor:   EpicTitleRefactoringOpportunities,
		EpicTypeFuture:     EpicTitleFutureEnhancements,
	}
}

// ensureEpics verifies that all required epics exist in Beads.
// It loads the epic cache if it exists, otherwise creates the epics via bd CLI.
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
func (p *BeadsProvider) loadEpicCache() error {
	data, err := os.ReadFile(p.cachePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, p.epicCache)
}

// saveEpicCache saves the epic ID mappings to the cache file.
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

// Compile-time interface check: verify BeadsProvider satisfies BacklogProvider
// Note: These are placeholder methods that will be implemented in subsequent tasks.
var _ BacklogProvider = (*BeadsProvider)(nil)

// List returns backlog items from Beads.
// TODO: Implement in Task 5
func (p *BeadsProvider) List(include []string, exclude []string) ([]string, error) {
	return nil, fmt.Errorf("List: not yet implemented")
}

// Add adds a new item to the specified section.
// TODO: Implement in Task 5
func (p *BeadsProvider) Add(description, itemType string) error {
	return fmt.Errorf("Add: not yet implemented")
}

// AddWithPattern adds a new item with pattern tracking for refactor items.
// TODO: Implement in Task 5
func (p *BeadsProvider) AddWithPattern(description, itemType, patternID string) error {
	return fmt.Errorf("AddWithPattern: not yet implemented")
}

// AddMultiple adds multiple items to the specified section.
// TODO: Implement in Task 5
func (p *BeadsProvider) AddMultiple(descriptions []string, itemType string) error {
	return fmt.Errorf("AddMultiple: not yet implemented")
}

// Complete marks an item as completed.
// TODO: Implement in Task 6
func (p *BeadsProvider) Complete(itemText string) error {
	return fmt.Errorf("Complete: not yet implemented")
}

// Cleanup removes completed items.
// TODO: Implement in Task 6
func (p *BeadsProvider) Cleanup(itemType string) (int, error) {
	return 0, fmt.Errorf("Cleanup: not yet implemented")
}

// GetPatternCount returns the occurrence count for a pattern ID.
// TODO: Implement in Task 6
func (p *BeadsProvider) GetPatternCount(patternID string) (int, error) {
	return 0, fmt.Errorf("GetPatternCount: not yet implemented")
}

// Decompose adds multiple sub-intents with dependency tracking.
// TODO: Implement in Task 7
func (p *BeadsProvider) Decompose(jsonInput string) error {
	return fmt.Errorf("Decompose: not yet implemented")
}
