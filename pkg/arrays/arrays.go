package arrays

// CheckEntryInArray checks if an entry is in an array, and returns true
// and the entry if it is found.
func CheckEntryInArray[T comparable](entry T, array []T) (value T, _ bool) {
	for _, e := range array {
		if e == entry {
			return e, true
		}
	}
	return value, false
}
