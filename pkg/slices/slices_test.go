package slices_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/pkg/slices"
)

func TestChunkSlice(t *testing.T) {
	testCases := []struct {
		name      string
		input     []string
		chunkSize int
		expected  [][]string
	}{
		{
			name:      "Empty slice",
			input:     []string{},
			chunkSize: 3,
			expected:  [][]string{{}},
		},
		{
			name:      "Slice smaller than chunk size",
			input:     []string{"a", "b"},
			chunkSize: 3,
			expected:  [][]string{{"a", "b"}},
		},
		{
			name:      "Slice equal to chunk size",
			input:     []string{"a", "b", "c"},
			chunkSize: 3,
			expected:  [][]string{{"a", "b", "c"}},
		},
		{
			name:      "Slice larger than chunk size",
			input:     []string{"a", "b", "c", "d", "e"},
			chunkSize: 2,
			expected:  [][]string{{"a", "b"}, {"c", "d"}, {"e"}},
		},
		{
			name:      "Chunk size of 1",
			input:     []string{"a", "b", "c"},
			chunkSize: 1,
			expected:  [][]string{{"a"}, {"b"}, {"c"}},
		},
		{
			name:      "Large slice with uneven chunks",
			input:     []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
			chunkSize: 3,
			expected:  [][]string{{"1", "2", "3"}, {"4", "5", "6"}, {"7", "8", "9"}, {"10"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := slices.Chunk(tc.input, tc.chunkSize)
			require.Equal(t, tc.expected, result)
		})
	}
}
