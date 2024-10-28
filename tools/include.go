package tools

// Include slice包含item
func Include[T comparable](slice []T, item T) bool {
	for i := range slice {
		if slice[i] == item {
			return true
		}
	}
	return false
}

// Includes slice包含全部items
func Includes[T comparable](slice []T, items ...T) bool {
	m := make(map[T]struct{}, len(slice))
	for i := range slice {
		m[slice[i]] = struct{}{}
	}
	for i := range items {
		if _, ok := m[items[i]]; !ok {
			return false
		}
	}
	return true
}

// IncludesAny slice包含任意一个items
func IncludesAny[T comparable](slice []T, items ...T) bool {
	m := make(map[T]struct{}, len(slice))
	for i := range slice {
		m[slice[i]] = struct{}{}
	}
	for i := range items {
		if _, ok := m[items[i]]; ok {
			return true
		}
	}
	return false
}
