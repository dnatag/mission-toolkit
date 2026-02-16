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
	"strings"
)

// List format constants for matching BacklogManager output format.
const (
	ListItemOpenFormat   = "- [ ] "
	ListItemClosedFormat = "- [x] "
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
// It queries all relevant epics and formats the output to match the existing
// BacklogManager format with "- [ ] description" for open items and "- [x] description"
// for completed items.
func (p *BeadsProvider) List(include []string, exclude []string) ([]string, error) {
	// Ensure epics are loaded
	if err := p.ensureEpics(); err != nil {
		return nil, fmt.Errorf("failed to ensure epics: %w", err)
	}

	// Determine which epic types to query
	epicTypes := []string{
		EpicTypeFeature,
		EpicTypeBugfix,
		EpicTypeDecomposed,
		EpicTypeRefactor,
		EpicTypeFuture,
	}

	// Apply include/exclude filters
	queryTypes := p.filterEpicTypes(epicTypes, include, exclude)

	var results []string
	for _, itemType := range queryTypes {
		epicID, ok := p.epicCache.Epics[itemType]
		if !ok {
			continue // Skip if epic doesn't exist
		}

		// Query children of this epic
		output, err := p.commandRunner.Run("list", "--parent", epicID, "--json")
		if err != nil {
			return nil, fmt.Errorf("failed to list %s items: %w", itemType, err)
		}

		// Parse JSON output and format results
		items, err := p.parseListOutput(output)
		if err != nil {
			return nil, fmt.Errorf("failed to parse list output for %s: %w", itemType, err)
		}

		results = append(results, items...)
	}

	return results, nil
}

// filterEpicTypes filters epic types based on include/exclude parameters.
// Uses O(n) lookup with maps for better performance with larger filter lists.
func (p *BeadsProvider) filterEpicTypes(types []string, include []string, exclude []string) []string {
	if len(include) == 0 && len(exclude) == 0 {
		return types
	}

	// Convert exclude/include slices to maps for O(1) lookup
	excludeMap := make(map[string]bool, len(exclude))
	for _, ex := range exclude {
		excludeMap[ex] = true
	}

	includeMap := make(map[string]bool, len(include))
	for _, inc := range include {
		includeMap[inc] = true
	}

	filtered := make([]string, 0, len(types))
	for _, t := range types {
		// Check exclude first
		if excludeMap[t] {
			continue
		}

		// If include is specified, only include matching types
		if len(include) > 0 && !includeMap[t] {
			continue
		}

		filtered = append(filtered, t)
	}

	return filtered
}

// parseListOutput parses JSON output from bd list command and formats items.
// Expected format is an array of items with fields: id, title, status.
func (p *BeadsProvider) parseListOutput(output string) ([]string, error) {
	var items []struct {
		ID     string `json:"id"`
		Title  string `json:"title"`
		Status string `json:"status"`
	}

	if err := json.Unmarshal([]byte(output), &items); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	results := make([]string, 0, len(items))
	for _, item := range items {
		// Format: "- [ ] description" or "- [x] description"
		// Status "open" -> "- [ ]", status "closed" -> "- [x]"
		var prefix string
		switch item.Status {
		case "closed":
			prefix = ListItemClosedFormat
		default:
			prefix = ListItemOpenFormat
		}

		results = append(results, prefix+item.Title)
	}

	return results, nil
}

// Add adds a new item to the specified section.
// It creates a new task under the appropriate epic as a child item.
func (p *BeadsProvider) Add(description, itemType string) error {
	// Ensure epics are loaded
	if err := p.ensureEpics(); err != nil {
		return fmt.Errorf("failed to ensure epics: %w", err)
	}

	// Validate item type
	validTypes := map[string]bool{
		EpicTypeFeature:    true,
		EpicTypeBugfix:     true,
		EpicTypeDecomposed: true,
		EpicTypeRefactor:   true,
		EpicTypeFuture:     true,
	}

	if !validTypes[itemType] {
		return fmt.Errorf("invalid type: %s. Valid types: feature, bugfix, decomposed, refactor, future", itemType)
	}

	// Look up epic ID for this item type
	epicID, ok := p.epicCache.Epics[itemType]
	if !ok {
		return fmt.Errorf("epic not found for type: %s", itemType)
	}

	// Create task under the epic
	output, err := p.commandRunner.Run("create", description, "-t", "task", "--parent", epicID, "--json")
	if err != nil {
		return fmt.Errorf("failed to create %s item '%s': %w", itemType, description, err)
	}

	// Parse JSON output to confirm creation
	var result struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return fmt.Errorf("failed to parse create output: %w", err)
	}

	if result.ID == "" {
		return fmt.Errorf("create returned empty ID for %s item '%s'", itemType, description)
	}

	return nil
}

// AddWithPattern adds a new item with pattern tracking for refactor items.
// For refactor items with a patternID, it tracks the occurrence count in the task notes.
func (p *BeadsProvider) AddWithPattern(description, itemType, patternID string) error {
	// First, add the item normally
	if err := p.Add(description, itemType); err != nil {
		return err
	}

	// Only add pattern tracking for refactor items
	if itemType != EpicTypeRefactor || patternID == "" {
		return nil
	}

	// Get the current pattern count to determine the new count
	currentCount, err := p.GetPatternCount(patternID)
	if err != nil {
		return fmt.Errorf("failed to get pattern count: %w", err)
	}
	newCount := currentCount + 1

	// Find the most recently created task (the one we just added)
	// We need to get the task ID from the most recent create operation
	// For now, we'll query the refactor epic to find the task by description
	epicID, ok := p.epicCache.Epics[itemType]
	if !ok {
		return fmt.Errorf("epic not found for type: %s", itemType)
	}

	output, err := p.commandRunner.Run("list", "--parent", epicID, "--json")
	if err != nil {
		return fmt.Errorf("failed to list %s items: %w", itemType, err)
	}

	var items []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}
	if err := json.Unmarshal([]byte(output), &items); err != nil {
		return fmt.Errorf("failed to parse list output: %w", err)
	}

	// Find the task with matching description (most recently created)
	var taskID string
	for i := len(items) - 1; i >= 0; i-- {
		if items[i].Title == description {
			taskID = items[i].ID
			break
		}
	}

	if taskID == "" {
		return fmt.Errorf("failed to find newly created task for '%s'", description)
	}

	// Update task notes with pattern tracking
	notes := fmt.Sprintf("PATTERN:%s COUNT:%d", patternID, newCount)
	_, err = p.commandRunner.Run("update", taskID, "--notes", notes)
	if err != nil {
		return fmt.Errorf("failed to update task notes: %w", err)
	}

	return nil
}

// AddMultiple adds multiple items to the specified section.
// It creates each description as a separate task under the appropriate epic.
func (p *BeadsProvider) AddMultiple(descriptions []string, itemType string) error {
	if len(descriptions) == 0 {
		return nil
	}

	// Ensure epics are loaded
	if err := p.ensureEpics(); err != nil {
		return fmt.Errorf("failed to ensure epics: %w", err)
	}

	// Validate item type
	validTypes := map[string]bool{
		EpicTypeFeature:    true,
		EpicTypeBugfix:     true,
		EpicTypeDecomposed: true,
		EpicTypeRefactor:   true,
		EpicTypeFuture:     true,
	}

	if !validTypes[itemType] {
		return fmt.Errorf("invalid type: %s. Valid types: feature, bugfix, decomposed, refactor, future", itemType)
	}

	// Look up epic ID for this item type
	epicID, ok := p.epicCache.Epics[itemType]
	if !ok {
		return fmt.Errorf("epic not found for type: %s", itemType)
	}

	// Create each item
	for _, description := range descriptions {
		output, err := p.commandRunner.Run("create", description, "-t", "task", "--parent", epicID, "--json")
		if err != nil {
			return fmt.Errorf("failed to create %s item '%s': %w", itemType, description, err)
		}

		// Parse JSON output to confirm creation
		var result struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal([]byte(output), &result); err != nil {
			return fmt.Errorf("failed to parse create output: %w", err)
		}

		if result.ID == "" {
			return fmt.Errorf("create returned empty ID for %s item '%s'", itemType, description)
		}
	}

	return nil
}

// Complete marks an item as completed.
// It searches all epic types for a task with matching title and marks it as closed.
func (p *BeadsProvider) Complete(itemText string) error {
	// Ensure epics are loaded
	if err := p.ensureEpics(); err != nil {
		return fmt.Errorf("failed to ensure epics: %w", err)
	}

	// Search all epic types for the matching task
	epicTypes := []string{
		EpicTypeFeature,
		EpicTypeBugfix,
		EpicTypeDecomposed,
		EpicTypeRefactor,
		EpicTypeFuture,
	}

	for _, itemType := range epicTypes {
		epicID, ok := p.epicCache.Epics[itemType]
		if !ok {
			continue
		}

		// Query children of this epic
		output, err := p.commandRunner.Run("list", "--parent", epicID, "--json")
		if err != nil {
			return fmt.Errorf("failed to list %s items: %w", itemType, err)
		}

		var items []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		}
		if err := json.Unmarshal([]byte(output), &items); err != nil {
			return fmt.Errorf("failed to parse list output for %s: %w", itemType, err)
		}

		// Find the task with matching title
		for _, item := range items {
			if item.Title == itemText {
				// Mark the task as closed
				_, err := p.commandRunner.Run("close", item.ID, "--reason", "Completed")
				if err != nil {
					return fmt.Errorf("failed to close task '%s': %w", itemText, err)
				}
				return nil
			}
		}
	}

	return fmt.Errorf("item not found: %s", itemText)
}

// Cleanup removes completed items.
// For Beads, this counts and returns the number of closed items across the specified epic types.
// Actual deletion is not performed since closed items are already archived in Beads.
func (p *BeadsProvider) Cleanup(itemType string) (int, error) {
	// Ensure epics are loaded
	if err := p.ensureEpics(); err != nil {
		return 0, fmt.Errorf("failed to ensure epics: %w", err)
	}

	// Determine which epic types to query
	var queryTypes []string
	if itemType != "" {
		// Validate item type
		validTypes := map[string]bool{
			EpicTypeFeature:    true,
			EpicTypeBugfix:     true,
			EpicTypeDecomposed: true,
			EpicTypeRefactor:   true,
			EpicTypeFuture:     true,
		}
		if !validTypes[itemType] {
			return 0, fmt.Errorf("invalid type: %s. Valid types: feature, bugfix, decomposed, refactor, future", itemType)
		}
		queryTypes = []string{itemType}
	} else {
		// Query all epic types
		queryTypes = []string{
			EpicTypeFeature,
			EpicTypeBugfix,
			EpicTypeDecomposed,
			EpicTypeRefactor,
			EpicTypeFuture,
		}
	}

	// Count closed items across the specified epic types
	closedCount := 0
	for _, itemType := range queryTypes {
		epicID, ok := p.epicCache.Epics[itemType]
		if !ok {
			continue
		}

		// Query children of this epic
		output, err := p.commandRunner.Run("list", "--parent", epicID, "--json")
		if err != nil {
			return 0, fmt.Errorf("failed to list %s items: %w", itemType, err)
		}

		var items []struct {
			ID     string `json:"id"`
			Title  string `json:"title"`
			Status string `json:"status"`
		}
		if err := json.Unmarshal([]byte(output), &items); err != nil {
			return 0, fmt.Errorf("failed to parse list output for %s: %w", itemType, err)
		}

		// Count closed items
		for _, item := range items {
			if item.Status == "closed" {
				closedCount++
			}
		}
	}

	return closedCount, nil
}

// GetPatternCount returns the occurrence count for a pattern ID.
// It queries the refactor epic and counts items with matching pattern notes.
func (p *BeadsProvider) GetPatternCount(patternID string) (int, error) {
	// Ensure epics are loaded
	if err := p.ensureEpics(); err != nil {
		return 0, fmt.Errorf("failed to ensure epics: %w", err)
	}

	// Get the refactor epic ID
	epicID, ok := p.epicCache.Epics[EpicTypeRefactor]
	if !ok {
		return 0, nil // No refactor epic means no patterns
	}

	// Query all items in the refactor epic
	output, err := p.commandRunner.Run("list", "--parent", epicID, "--json")
	if err != nil {
		return 0, fmt.Errorf("failed to list refactor items: %w", err)
	}

	var items []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Notes string `json:"notes"`
	}
	if err := json.Unmarshal([]byte(output), &items); err != nil {
		return 0, fmt.Errorf("failed to parse list output: %w", err)
	}

	// Count items with matching pattern in notes
	maxCount := 0
	patternPrefix := fmt.Sprintf("PATTERN:%s", patternID)
	for _, item := range items {
		if strings.Contains(item.Notes, patternPrefix) {
			// Parse COUNT: from notes (format: "PATTERN:xxx COUNT:N")
			parts := strings.Split(item.Notes, " ")
			for _, part := range parts {
				if strings.HasPrefix(part, "COUNT:") {
					countStr := strings.TrimPrefix(part, "COUNT:")
					var count int
					if _, err := fmt.Sscanf(countStr, "%d", &count); err == nil {
						if count > maxCount {
							maxCount = count
						}
					}
				}
			}
		}
	}

	return maxCount, nil
}

// Decompose adds multiple sub-intents with dependency tracking.
// It parses JSON input containing sub-intents with dependencies and creates
// them in Beads with proper dependency wiring using the bd CLI's dependency graph.
//
// The JSON input format matches the BacklogProvider interface:
//
//	{
//	  "action": "decompose",
//	  "sub_intents": [
//	    {
//	      "intent": "Task description",
//	      "rationale": "Why this task is needed",
//	      "estimated_files": 3,
//	      "dependencies": ["Other task description"]
//	    }
//	  ],
//	  "decomposition_rationale": "Why we decomposed this"
//	}
//
// The implementation:
// 1. Creates all tasks in the decomposed epic first
// 2. Maps intent descriptions to Beads task IDs
// 3. Wires dependencies using Beads' native dependency graph
//
// This two-phase approach ensures all tasks exist before wiring dependencies,
// avoiding ordering constraints in the input.
func (p *BeadsProvider) Decompose(jsonInput string) error {
	// Parse the decompose JSON input
	var decompose struct {
		Action     string `json:"action"`
		SubIntents []struct {
			Intent         string   `json:"intent"`
			Rationale      string   `json:"rationale"`
			EstimatedFiles int      `json:"estimated_files"`
			Dependencies   []string `json:"dependencies"`
		} `json:"sub_intents"`
		DecompositionRationale string `json:"decomposition_rationale"`
	}

	if err := json.Unmarshal([]byte(jsonInput), &decompose); err != nil {
		return fmt.Errorf("parsing decompose JSON: %w", err)
	}

	if len(decompose.SubIntents) == 0 {
		return fmt.Errorf("no sub-intents found in decompose input")
	}

	// Ensure epics are loaded
	if err := p.ensureEpics(); err != nil {
		return fmt.Errorf("failed to ensure epics: %w", err)
	}

	// Get the decomposed epic ID
	epicID, ok := p.epicCache.Epics[EpicTypeDecomposed]
	if !ok {
		return fmt.Errorf("decomposed epic not found in cache")
	}

	// Phase 1: Create all tasks and build the intent-to-task-ID mapping.
	// We create all tasks first to avoid ordering issues with dependencies.
	intentToTaskID := make(map[string]string, len(decompose.SubIntents))
	for _, subIntent := range decompose.SubIntents {
		// Create task using bd CLI: bd create "<intent>" -t task --parent <epicID> --json
		output, err := p.commandRunner.Run("create", subIntent.Intent, "-t", "task", "--parent", epicID, "--json")
		if err != nil {
			return fmt.Errorf("failed to create task for sub-intent '%s': %w", subIntent.Intent, err)
		}

		// Parse JSON output to extract the created task ID
		var result struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal([]byte(output), &result); err != nil {
			return fmt.Errorf("failed to parse create output for '%s': %w", subIntent.Intent, err)
		}

		if result.ID == "" {
			return fmt.Errorf("bd create returned empty ID for sub-intent '%s'", subIntent.Intent)
		}

		// Store the mapping from intent description to task ID for dependency wiring
		intentToTaskID[subIntent.Intent] = result.ID
	}

	// Phase 2: Wire dependencies using Beads' native dependency graph.
	// Now that all tasks exist, we can safely add dependencies between them.
	for _, subIntent := range decompose.SubIntents {
		taskID := intentToTaskID[subIntent.Intent]

		// Skip tasks with no dependencies
		if len(subIntent.Dependencies) == 0 {
			continue
		}

		// For each dependency reference, find the corresponding task ID and add the dependency
		for _, depIntent := range subIntent.Dependencies {
			depTaskID, ok := intentToTaskID[depIntent]
			if !ok {
				return fmt.Errorf("dependency task not found: '%s' referenced by '%s'", depIntent, subIntent.Intent)
			}

			// Add dependency using bd CLI: bd add-dep <taskID> <depTaskID>
			_, err := p.commandRunner.Run("add-dep", taskID, depTaskID)
			if err != nil {
				return fmt.Errorf("failed to add dependency from '%s' to '%s': %w", subIntent.Intent, depIntent, err)
			}
		}
	}

	return nil
}
