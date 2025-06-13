package sliceset

func RemoveDuplicates[T comparable](slice []T) []T {
	var (
		seen   = make(map[T]bool)
		result = make([]T, 0, len(slice))
	)
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}
