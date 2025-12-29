package utils

// Section represents a Markdown section with header and content
type Section struct {
	Header  string
	Content []interface{} // Can contain strings or key-value pairs
}

// KeyValue represents a key-value pair from list items
type KeyValue struct {
	Key   string
	Value string
}
