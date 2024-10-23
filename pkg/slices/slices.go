package slices

// Chunk chunks a slice into batches of chunkSize.
// example: {1,2,3,4,5}, chunkSize = 2 -> {1,2}, {3,4}, {5}
func Chunk[T any](input []T, chunkSize int) [][]T {
	if len(input) <= chunkSize {
		return [][]T{input}
	}
	var chunks [][]T
	for i := 0; i < len(input); i += chunkSize {
		end := i + chunkSize

		if end > len(input) {
			end = len(input)
		}

		chunks = append(chunks, input[i:end])
	}
	return chunks
}
