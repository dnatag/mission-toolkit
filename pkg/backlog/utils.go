package backlog

// FilterStringSlice filters a slice based on include/exclude parameters.
func FilterStringSlice(items []string, include []string, exclude []string) []string {
	if len(include) == 0 && len(exclude) == 0 {
		return items
	}

	excludeMap := make(map[string]bool, len(exclude))
	for _, ex := range exclude {
		excludeMap[ex] = true
	}

	includeMap := make(map[string]bool, len(include))
	for _, inc := range include {
		includeMap[inc] = true
	}

	filtered := make([]string, 0, len(items))
	for _, item := range items {
		if excludeMap[item] {
			continue
		}
		if len(include) > 0 && !includeMap[item] {
			continue
		}
		filtered = append(filtered, item)
	}

	return filtered
}

// ContainsString checks if a string slice contains a specific value.
func ContainsString(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}

// RemoveString returns a new slice with all occurrences of val removed.
func RemoveString(slice []string, val string) []string {
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if s != val {
			result = append(result, s)
		}
	}
	return result
}
