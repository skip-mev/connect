package types_test

import (
	"testing"

	"github.com/skip-mev/slinky/x/marketmap/types"
	"github.com/stretchr/testify/require"
)

func TestAggregateMarketConfig(t *testing.T) {
	testCases := []struct {
		name    string
		markets map[string]types.MarketConfig
		tickers map[string]types.PathsConfig
		expErr  bool
	}{
		{
			name:    "empty config",
			markets: map[string]types.MarketConfig{},
			tickers: map[string]types.PathsConfig{},
			expErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := types.NewAggregateMarketConfig(tc.markets, tc.tickers)
			if tc.expErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
