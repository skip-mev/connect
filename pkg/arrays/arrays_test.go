package arrays_test

import (
	"testing"

	"github.com/skip-mev/connect/v2/pkg/arrays"
)

func TestCheckEntryInArray(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		entry int
		array []int
		want  bool
	}{
		{
			name:  "entry in array",
			entry: 1,
			array: []int{1, 2, 3},
			want:  true,
		},
		{
			name:  "entry not in array",
			entry: 4,
			array: []int{1, 2, 3},
			want:  false,
		},
		{
			name:  "empty array",
			entry: 1,
			array: []int{},
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if _, got := arrays.CheckEntryInArray(tt.entry, tt.array); got != tt.want {
				t.Errorf("CheckEntryInArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
