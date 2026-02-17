package backlog

// Item type constants for backlog items.
const (
	ItemTypeFeature    = "feature"
	ItemTypeBugfix     = "bugfix"
	ItemTypeDecomposed = "decomposed"
	ItemTypeRefactor   = "refactor"
	ItemTypeFuture     = "future"
)

// Markdown list format constants.
const (
	ListItemOpenFormat   = "- [ ] "
	ListItemClosedFormat = "- [x] "
)

// ValidTypes is the set of valid backlog item types.
var ValidTypes = map[string]bool{
	ItemTypeFeature:    true,
	ItemTypeBugfix:     true,
	ItemTypeDecomposed: true,
	ItemTypeRefactor:   true,
	ItemTypeFuture:     true,
}

// IsValidType checks if the given type is a valid backlog item type.
func IsValidType(itemType string) bool {
	return ValidTypes[itemType]
}
