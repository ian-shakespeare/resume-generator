package utils

// Check if some element meets the condition
func Some[T any](arr []T, cond func(T) bool) bool {
	for i := 0; i < len(arr); i += 1 {
		if cond(arr[i]) {
			return true
		}
	}
	return false
}

// Find the index of an element that meets the condition.
// Returns `-1` if no element matches the condition.
func Find[T any](arr []T, cond func(T) bool) int {
	for i := 0; i < len(arr); i += 1 {
		if cond(arr[i]) {
			return i
		}
	}
	return -1
}
