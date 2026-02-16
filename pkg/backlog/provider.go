package backlog

// BacklogProvider defines the interface for backlog operations.
// Implementations must support listing, adding, completing, and cleaning up backlog items,
// as well as pattern tracking and epic decomposition.
type BacklogProvider interface {
	// List returns backlog items, optionally including completed items and filtering by type.
	// include: types to include (decomposed, refactor, future, completed)
	// exclude: types to exclude (decomposed, refactor, future, completed)
	List(include []string, exclude []string) ([]string, error)

	// Add adds a new item to the specified section.
	Add(description, itemType string) error

	// AddWithPattern adds a new item with optional pattern ID tracking.
	// For refactor items with a patternID, increments count if pattern exists.
	AddWithPattern(description, itemType, patternID string) error

	// AddMultiple adds multiple items to the specified section in a single operation.
	AddMultiple(descriptions []string, itemType string) error

	// Complete marks an item as completed and moves it to the COMPLETED section.
	Complete(itemText string) error

	// Cleanup removes completed items from the COMPLETED section.
	// If itemType is provided, only removes completed items that match that type.
	// Returns the number of items removed.
	Cleanup(itemType string) (int, error)

	// GetPatternCount returns the occurrence count for a pattern ID.
	// Returns 0 if pattern not found.
	GetPatternCount(patternID string) (int, error)

	// Decompose adds multiple sub-intents with dependency tracking.
	// Accepts a JSON string containing decomposed intent information.
	// For implementations without dependency graph support, falls back to AddMultiple.
	Decompose(jsonInput string) error
}

// Compile-time interface check: verify BacklogManager satisfies BacklogProvider
var _ BacklogProvider = (*BacklogManager)(nil)

// NewProvider creates a BacklogProvider for the given mission directory.
// Auto-detects Beads availability and returns BeadsProvider when available,
// otherwise falls back to BacklogManager.
func NewProvider(missionDir string) BacklogProvider {
	// Check if Beads is available in the mission directory
	available, err := BeadsAvailableInDir(missionDir)
	if err != nil {
		// Error indicates an issue with PATH lookup or .beads validation
		// (e.g., .beads exists but is not a directory). Fall back to BacklogManager.
		return NewManager(missionDir)
	}

	if available {
		return NewBeadsProvider(missionDir)
	}

	// Fallback to BacklogManager when Beads is not available
	return NewManager(missionDir)
}
