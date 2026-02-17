package backlog

import (
	"fmt"
	"strings"
)

// Reconciler detects and syncs rogue edits between backlog.md and Beads.
// It compares items in both sources and identifies discrepancies.
type Reconciler struct {
	fileProvider  BacklogProvider // For reading backlog.md
	beadsProvider BacklogProvider // For reading/writing Beads
}

// NewReconciler creates a new Reconciler with the given providers.
func NewReconciler(fileProvider, beadsProvider BacklogProvider) *Reconciler {
	return &Reconciler{
		fileProvider:  fileProvider,
		beadsProvider: beadsProvider,
	}
}

// RogueEdit represents an item that exists in backlog.md but not in Beads.
type RogueEdit struct {
	Description  string   // The item description
	Type         string   // The item type (feature, bugfix, decomposed, refactor, future)
	Status       string   // The item status ("open" or "completed")
	Dependencies []string // Dependencies extracted from description (if any)
}

// DetectRogueEdits compares items between backlog.md and Beads,
// returning items that exist in backlog.md but not in Beads.
//
// The method:
// 1. Lists all items from both providers
// 2. Normalizes descriptions for comparison
// 3. Identifies items present in backlog.md but missing from Beads
// 4. Returns the list of rogue edits with their inferred types
func (r *Reconciler) DetectRogueEdits() ([]RogueEdit, error) {
	// Get all items from backlog.md (including completed for full sync)
	fileItems, err := r.fileProvider.List(nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list backlog.md items: %w", err)
	}

	// Get all items from Beads (including completed for full comparison)
	beadsItems, err := r.beadsProvider.List(nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list Beads items: %w", err)
	}

	// Create a map of Beads items for quick lookup
	beadsSet := make(map[string]bool, len(beadsItems))
	for _, item := range beadsItems {
		normalized := normalizeDescription(item)
		beadsSet[normalized] = true
	}

	// Find rogue edits (items in backlog.md but not in Beads)
	var rogueEdits []RogueEdit
	for _, item := range fileItems {
		normalized := normalizeDescription(item)

		// Skip if the item exists in Beads
		if beadsSet[normalized] {
			continue
		}

		// Parse the item to extract type, status, and description
		edit := parseBacklogItem(item)
		if edit != nil {
			rogueEdits = append(rogueEdits, *edit)
		}
	}

	return rogueEdits, nil
}

// ImportToBeads imports the given rogue edits into Beads.
// Each edit is added to the appropriate epic based on its type.
//
// The method preserves:
// - Item type (feature, bugfix, decomposed, refactor, future)
// - Dependencies (extracted from description)
// - Completion status (completed items are closed after creation)
//
// Returns the number of items successfully imported and any error encountered.
func (r *Reconciler) ImportToBeads(edits []RogueEdit) (int, error) {
	// Get current Beads items for deduplication
	beadsItems, err := r.beadsProvider.List(nil, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to list Beads items for dedup: %w", err)
	}
	beadsSet := make(map[string]bool, len(beadsItems))
	for _, item := range beadsItems {
		beadsSet[normalizeDescription(item)] = true
	}

	imported := 0

	for _, edit := range edits {
		// Skip if already in Beads (dedup guard)
		if beadsSet[normalizeDescription(edit.Description)] {
			continue
		}

		// Build description with dependencies if present
		description := edit.Description
		if len(edit.Dependencies) > 0 {
			description += fmt.Sprintf(" (depends on: %s)", strings.Join(edit.Dependencies, ", "))
		}

		if err := r.beadsProvider.Add(description, edit.Type); err != nil {
			return imported, fmt.Errorf("failed to import '%s': %w", edit.Description, err)
		}

		if edit.Status == "completed" {
			r.beadsProvider.Complete(edit.Description) //nolint:errcheck // Non-critical: item already created
		}

		imported++
	}

	return imported, nil
}

// normalizeDescription removes list markers and normalizes whitespace for comparison.
func normalizeDescription(item string) string {
	// Trim leading/trailing spaces first
	item = strings.TrimSpace(item)

	// Remove list markers: "- [ ] " or "- [x] "
	item = strings.TrimPrefix(item, ListItemOpenFormat)
	item = strings.TrimPrefix(item, ListItemClosedFormat)

	// Remove pattern markers: "[PATTERN:xxx][COUNT:nn] "
	item = removePatternMarker(item)

	// Remove completion timestamp: " (Completed: YYYY-MM-DD)"
	item = removeCompletionTimestamp(item)

	// Normalize whitespace
	item = strings.TrimSpace(item)
	item = strings.Join(strings.Fields(item), " ")

	return item
}

// parseBacklogItem extracts type, status, and description from a backlog item.
// Returns nil if the item cannot be parsed.
func parseBacklogItem(item string) *RogueEdit {
	edit := &RogueEdit{}

	// Determine status from checkbox
	if strings.HasPrefix(item, ListItemClosedFormat) {
		edit.Status = "completed"
		item = strings.TrimPrefix(item, ListItemClosedFormat)
	} else if strings.HasPrefix(item, ListItemOpenFormat) {
		edit.Status = "open"
		item = strings.TrimPrefix(item, ListItemOpenFormat)
	} else {
		// Not a valid backlog item
		return nil
	}

	// Remove pattern markers if present
	item = removePatternMarker(item)

	// Remove completion timestamp if present
	item = removeCompletionTimestamp(item)

	// Extract dependencies from description
	edit.Description, edit.Dependencies = extractDependencies(item)

	// Infer type from description content (heuristic approach)
	edit.Type = inferTypeFromDescription(edit.Description)

	return edit
}

// removePatternMarker removes pattern tracking markers from a description.
// Pattern format: "[PATTERN:xxx][COUNT:nn] description"
// Returns the description with all pattern markers removed.
func removePatternMarker(item string) string {
	// Look for pattern marker pattern
	if !strings.Contains(item, "[PATTERN:") {
		return item
	}

	// Find the closing "] " after COUNT marker
	// Example: "[PATTERN:helper-extract][COUNT:2] Extract helper function"
	// We need to find the second "] " (after COUNT)
	countIdx := strings.Index(item, "][COUNT:")
	if countIdx == -1 {
		// Only PATTERN marker, no COUNT
		if idx := strings.Index(item, "] "); idx != -1 {
			return item[idx+2:]
		}
		return item
	}

	// Find the closing bracket after COUNT
	countEnd := strings.Index(item[countIdx:], "] ")
	if countEnd == -1 {
		return item
	}

	// Return everything after the pattern marker
	return item[countIdx+countEnd+2:]
}

// removeCompletionTimestamp removes completion timestamp from a description.
func removeCompletionTimestamp(item string) string {
	// Timestamp format: " (Completed: YYYY-MM-DD)"
	if strings.Contains(item, " (Completed:") {
		idx := strings.Index(item, " (Completed:")
		if idx != -1 {
			item = item[:idx]
		}
	}
	return item
}

// extractDependencies parses dependency information from a description.
// Returns the clean description and any dependencies found.
func extractDependencies(item string) (string, []string) {
	// Look for dependency marker: " (depends on: X, Y, Z)"
	if strings.Contains(item, " (depends on:") {
		idx := strings.Index(item, " (depends on:")
		if idx != -1 {
			description := strings.TrimSpace(item[:idx])
			depPart := item[idx+13:] // Skip " (depends on: "

			// Find the closing parenthesis
			endIdx := strings.Index(depPart, ")")
			if endIdx != -1 {
				depList := depPart[:endIdx]
				dependencies := strings.Split(depList, ",")
				for i, dep := range dependencies {
					dependencies[i] = strings.TrimSpace(dep)
				}
				return description, dependencies
			}
		}
	}
	return strings.TrimSpace(item), nil
}

// inferTypeFromDescription attempts to determine item type from description content.
// This is a heuristic approach based on common patterns.
func inferTypeFromDescription(description string) string {
	desc := strings.ToLower(description)

	// Check for decomposed items (epic breakdowns)
	if strings.Contains(desc, "(from epic:") || strings.Contains(desc, "sub-intent") {
		return ItemTypeDecomposed
	}

	// Check for refactor items
	if strings.Contains(desc, "refactor") || strings.Contains(desc, "extract") ||
		strings.Contains(desc, "restructure") || strings.Contains(desc, "consolidate") {
		return ItemTypeRefactor
	}

	// Check for bugfix items
	if strings.Contains(desc, "fix") || strings.Contains(desc, "bug") ||
		strings.Contains(desc, "error") || strings.Contains(desc, "crash") ||
		strings.Contains(desc, "issue") {
		return ItemTypeBugfix
	}

	// Check for feature items (default for new capabilities)
	if strings.Contains(desc, "add") || strings.Contains(desc, "implement") ||
		strings.Contains(desc, "create") || strings.Contains(desc, "support") ||
		strings.Contains(desc, "feature") {
		return ItemTypeFeature
	}

	// Default to feature for unclear items
	return ItemTypeFeature
}
