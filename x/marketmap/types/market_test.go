package types_test

import (
	"testing"

	"github.com/skip-mev/slinky/x/marketmap/types"
	"github.com/stretchr/testify/require"
)

func TestMarketConfig(t *testing.T) {
	testCases := []struct {
		name   string
		market types.MarketConfig
		expErr bool
	}{
		{
			name: "valid market",
			market: types.MarketConfig{
				Name: "binance",
				TickerConfigs: map[string]types.TickerConfig{
					1: {
						Ticker: types.Ticker{
							Id:               1,
							Base:             "BITCOIN",
							Quote:            "USDT",
							Decimals:         8,
							MinProviderCount: 1,
						},
						OffChainTicker: "BTC/USDT",
					},
				},
			},
			expErr: false,
		},
		{
			name: "empty name",
			market: types.MarketConfig{
				Name: "",
				TickerConfigs: map[uint64]types.TickerConfig{
					1: {
						Ticker: types.Ticker{
							Id:               1,
							Base:             "BITCOIN",
							Quote:            "USDT",
							Decimals:         8,
							MinProviderCount: 1,
						},
						OffChainTicker: "BTC/USDT",
					},
				},
			},
			expErr: true,
		},
		{
			name: "empty ticker configs",
			market: types.MarketConfig{
				Name:          "binance",
				TickerConfigs: map[uint64]types.TickerConfig{},
			},
			expErr: true,
		},
		{
			name: "invalid ticker config",
			market: types.MarketConfig{
				Name: "binance",
				TickerConfigs: map[uint64]types.TickerConfig{
					1: {
						Ticker:         types.Ticker{},
						OffChainTicker: "BTC/USDT",
					},
				},
			},
			expErr: true,
		},
		{
			name: "invalid id",
			market: types.MarketConfig{
				Name: "binance",
				TickerConfigs: map[uint64]types.TickerConfig{
					2: {
						Ticker: types.Ticker{
							Id:               1,
							Base:             "BITCOIN",
							Quote:            "USDT",
							Decimals:         8,
							MinProviderCount: 1,
						},
						OffChainTicker: "BTC/USDT",
					},
				},
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		err := tc.market.ValidateBasic()
		if tc.expErr {
			require.Error(t, err)
			return
		}

		require.NoError(t, err)
	}
}
