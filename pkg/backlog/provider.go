package backlog

// BacklogProvider defines the interface for backlog operations.
type BacklogProvider interface {
	List(include []string, exclude []string) ([]string, error)
	Add(description, itemType string) error
	AddWithPattern(description, itemType, patternID string) error
	AddMultiple(descriptions []string, itemType string) error
	Complete(itemText string) error
	Cleanup(itemType string) (int, error)
	GetPatternCount(patternID string) (int, error)
	Decompose(jsonInput string) error
}
