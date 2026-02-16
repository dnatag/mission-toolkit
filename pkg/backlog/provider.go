package backlog

// BacklogProvider defines the interface for backlog operations.
// Implementations must support listing, adding, completing, and cleaning up backlog items,
// as well as epic decomposition. Pattern-specific methods like AddWithPattern and
// GetPatternCount are excluded as they are implementation-specific (e.g., markdown-based).
type BacklogProvider interface {
	// List returns backlog items, optionally including completed items and filtering by type.
	// include: types to include (decomposed, refactor, future, completed)
	// exclude: types to exclude (decomposed, refactor, future, completed)
	List(include []string, exclude []string) ([]string, error)

	// Add adds a new item to the specified section.
	Add(description, itemType string) error

	// AddMultiple adds multiple items to the specified section in a single operation.
	AddMultiple(descriptions []string, itemType string) error

	// Complete marks an item as completed and moves it to the COMPLETED section.
	Complete(itemText string) error

	// Cleanup removes completed items from the COMPLETED section.
	// If itemType is provided, only removes completed items that match that type.
	// Returns the number of items removed.
	Cleanup(itemType string) (int, error)

	// Decompose adds multiple sub-intents with dependency tracking.
	// Accepts a JSON string containing decomposed intent information.
	// For implementations without dependency graph support, falls back to AddMultiple.
	Decompose(jsonInput string) error
}

// Compile-time interface check: verify BacklogManager satisfies BacklogProvider
var _ BacklogProvider = (*BacklogManager)(nil)
