package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

func TestProviderMarketMap(t *testing.T) {
	testCases := []struct {
		name         string
		providerName string
		configs      types.TickerToProviderConfig
		expErr       bool
	}{
		{
			name:         "empty configs",
			providerName: "test",
			configs:      types.TickerToProviderConfig{},
			expErr:       true,
		},
		{
			name:         "empty provider name",
			providerName: "",
			configs: types.TickerToProviderConfig{
				constants.BITCOIN_USD: {
					Name:           "test",
					OffChainTicker: "BTC-USD",
				},
			},
			expErr: true,
		},
		{
			name:         "empty off-chain ticker",
			providerName: "test",
			configs: types.TickerToProviderConfig{
				constants.BITCOIN_USD: {
					Name:           "test",
					OffChainTicker: "",
				},
			},
			expErr: true,
		},
		{
			name:         "invalid ticker",
			providerName: "test",
			configs: types.TickerToProviderConfig{
				mmtypes.NewTicker("BTC", "USD", 8, 0): {
					Name:           "test",
					OffChainTicker: "BTC-USD",
				},
			},
			expErr: true,
		},
		{
			name:         "valid configs",
			providerName: "test",
			configs: types.TickerToProviderConfig{
				constants.BITCOIN_USD: {
					Name:           "test",
					OffChainTicker: "BTC-USD",
				},
				constants.BITCOIN_USDC: {
					Name:           "test",
					OffChainTicker: "BTC-USDC",
				},
			},
			expErr: false,
		},
		{
			name:         "mismatch provider name and config",
			providerName: "test",
			configs: types.TickerToProviderConfig{
				constants.BITCOIN_USD: {
					Name:           "invalid",
					OffChainTicker: "BTC-USD",
				},
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := types.NewProviderMarketMap(tc.providerName, tc.configs)
			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
