package types_test

import (
	"testing"

	"github.com/skip-mev/slinky/x/marketmap/types"
	"github.com/stretchr/testify/require"
)

func TestOperation(t *testing.T) {
	t.Run("valid operation", func(t *testing.T) {
		ticker := types.Ticker{
			Base:             "BITCOIN",
			Quote:            "USDT",
			Decimals:         8,
			MinProviderCount: 1,
		}

		_, err := types.NewOperation(ticker, false)
		require.NoError(t, err)
	})

	t.Run("invalid operation", func(t *testing.T) {
		ticker := types.Ticker{}
		_, err := types.NewOperation(ticker, false)
		require.Error(t, err)
	})
}
