// Package backlog provides backlog management functionality with support for
// multiple backend implementations.
package backlog

// Item type constants for backlog items.
// These types are used across both BacklogManager and BeadsProvider implementations.
const (
	// ItemTypeFeature represents a new feature or enhancement.
	ItemTypeFeature = "feature"

	// ItemTypeBugfix represents a bug fix or correction.
	ItemTypeBugfix = "bugfix"

	// ItemTypeDecomposed represents a decomposed sub-intent.
	ItemTypeDecomposed = "decomposed"

	// ItemTypeRefactor represents a refactoring opportunity.
	ItemTypeRefactor = "refactor"

	// ItemTypeFuture represents a future enhancement idea.
	ItemTypeFuture = "future"
)

// Markdown list format constants for matching BacklogManager output format.
const (
	// ListItemOpenFormat represents an unchecked/open list item.
	ListItemOpenFormat = "- [ ] "

	// ListItemClosedFormat represents a checked/closed list item.
	ListItemClosedFormat = "- [x] "
)
