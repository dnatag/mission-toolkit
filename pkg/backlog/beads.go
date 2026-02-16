// Package backlog provides backlog management functionality with support for
// multiple backend implementations. The BeadsProvider implementation uses the
// Beads (bd) CLI for dependency-aware task management.
package backlog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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
		ItemTypeFeature,
		ItemTypeBugfix,
		ItemTypeDecomposed,
		ItemTypeRefactor,
		ItemTypeFuture,
	}

	// Apply include/exclude filters
	queryTypes := filterStringSlice(epicTypes, include, exclude)

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
		ItemTypeFeature:    true,
		ItemTypeBugfix:     true,
		ItemTypeDecomposed: true,
		ItemTypeRefactor:   true,
		ItemTypeFuture:     true,
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

	// Generate snapshot after successful operation
	p.generateSnapshot()

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
	if itemType != ItemTypeRefactor || patternID == "" {
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

	// Generate snapshot after successful operation
	p.generateSnapshot()

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
		ItemTypeFeature:    true,
		ItemTypeBugfix:     true,
		ItemTypeDecomposed: true,
		ItemTypeRefactor:   true,
		ItemTypeFuture:     true,
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

	// Generate snapshot after successful operation
	p.generateSnapshot()

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
		ItemTypeFeature,
		ItemTypeBugfix,
		ItemTypeDecomposed,
		ItemTypeRefactor,
		ItemTypeFuture,
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
				// Generate snapshot after successful operation
				p.generateSnapshot()
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
			ItemTypeFeature:    true,
			ItemTypeBugfix:     true,
			ItemTypeDecomposed: true,
			ItemTypeRefactor:   true,
			ItemTypeFuture:     true,
		}
		if !validTypes[itemType] {
			return 0, fmt.Errorf("invalid type: %s. Valid types: feature, bugfix, decomposed, refactor, future", itemType)
		}
		queryTypes = []string{itemType}
	} else {
		// Query all epic types
		queryTypes = []string{
			ItemTypeFeature,
			ItemTypeBugfix,
			ItemTypeDecomposed,
			ItemTypeRefactor,
			ItemTypeFuture,
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
	epicID, ok := p.epicCache.Epics[ItemTypeRefactor]
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
	epicID, ok := p.epicCache.Epics[ItemTypeDecomposed]
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

	// Generate snapshot after successful operation
	p.generateSnapshot()

	return nil
}

// generateSnapshot writes the current Beads state to .mission/backlog.md.
// This creates a canonical representation of the backlog that can be used
// for reconciliation and manual inspection. Errors are logged but don't fail operations.
func (p *BeadsProvider) generateSnapshot() {
	// Build snapshot content
	var content strings.Builder

	// Add frontmatter
	content.WriteString("---\n")
	content.WriteString(fmt.Sprintf("last_updated: %s\n", time.Now().Format(time.RFC3339)))
	content.WriteString("source: beads\n")
	content.WriteString("---\n\n")

	// Add title
	content.WriteString("# Backlog\n\n")

	// Define section order and types
	sections := []struct {
		title string
		desc  string
		typ   string
	}{
		{"FEATURES", "User-defined feature requests and enhancements", ItemTypeFeature},
		{"BUGFIXES", "Bug reports and issues to be fixed", ItemTypeBugfix},
		{"DECOMPOSED INTENTS", "Sub-intents from Track 4 Epic requests that need separate missions", ItemTypeDecomposed},
		{"REFACTORING OPPORTUNITIES", "Detected duplication patterns that need DRY missions", ItemTypeRefactor},
		{"FUTURE ENHANCEMENTS", "Ideas and improvements for later consideration", ItemTypeFuture},
		{"COMPLETED", "History of completed backlog items", "completed"},
	}

	// Build each section
	for _, section := range sections {
		content.WriteString(fmt.Sprintf("## %s\n", section.title))
		content.WriteString(fmt.Sprintf("(%s)\n", section.desc))

		// Get items for this section
		var items []string
		var err error

		if section.typ == "completed" {
			// Get completed items
			items, err = p.List([]string{"completed"}, nil)
		} else {
			// Get open items for this type
			items, err = p.List([]string{section.typ}, []string{"completed"})
		}

		if err != nil {
			// Log error but continue with empty section
			fmt.Fprintf(os.Stderr, "Warning: failed to list %s items for snapshot: %v\n", section.typ, err)
			items = []string{}
		}

		// Add items to section
		for _, item := range items {
			content.WriteString(item + "\n")
		}

		content.WriteString("\n")
	}

	// Write to file
	backlogPath := filepath.Join(p.projectRoot, ".mission", "backlog.md")
	if err := os.MkdirAll(filepath.Dir(backlogPath), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to create .mission directory for snapshot: %v\n", err)
		return
	}

	if err := os.WriteFile(backlogPath, []byte(content.String()), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to write snapshot to %s: %v\n", backlogPath, err)
		return
	}
}
