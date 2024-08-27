package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/providers/apis/defi/types"
)

func TestBlockAgeChecker_IsHeightValid(t *testing.T) {
	tests := []struct {
		name       string
		lastHeight uint64
		waitTime   time.Duration
		maxAge     time.Duration
		newHeight  uint64
		isValid    bool
	}{
		{
			name:       "valid 0s no timeout",
			lastHeight: 0,
			waitTime:   0,
			maxAge:     10 * time.Minute,
			newHeight:  0,
			isValid:    true,
		},
		{
			name:       "valid new height no timeout",
			lastHeight: 0,
			waitTime:   0,
			maxAge:     10 * time.Minute,
			newHeight:  0,
			isValid:    true,
		},
		{
			name:       "invalid 0s due to timeout",
			lastHeight: 0,
			waitTime:   10 * time.Millisecond,
			maxAge:     1 * time.Millisecond,
			newHeight:  0,
			isValid:    false,
		},
		{
			name:       "valid timeout but block height increase",
			lastHeight: 0,
			waitTime:   10 * time.Millisecond,
			maxAge:     1 * time.Millisecond,
			newHeight:  1,
			isValid:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bc := types.NewBlockAgeChecker(tt.maxAge)

			got := bc.IsHeightValid(tt.lastHeight)
			require.True(t, got)
			time.Sleep(tt.waitTime)

			got = bc.IsHeightValid(tt.newHeight)
			require.Equal(t, tt.isValid, got)
		})
	}
}
