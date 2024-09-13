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

// Check if array contains element
func Contains[T comparable](arr []T, value T) bool {
	return Some(arr, func(existing T) bool {
		return existing == value
	})
}

func Map[In, Out any](arr []In, transformFunc func(elem In, index int) Out) []Out {
	out := []Out{}
	for i := 0; i < len(arr); i += 1 {
		out = append(out, transformFunc(arr[i], i))
	}
	return out
}

func Filter[T any](arr []T, cond func(elem T) bool) []T {
	out := []T{}
	for i := 0; i < len(arr); i += 1 {
		if cond(arr[i]) {
			out = append(out, arr[i])
		}
	}
	return out
}
