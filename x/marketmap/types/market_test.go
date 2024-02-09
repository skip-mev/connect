package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/x/marketmap/types"
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
					"BITCOIN/USDT": {
						Ticker: types.Ticker{
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
				TickerConfigs: map[string]types.TickerConfig{
					"BITCOIN/USDT": {
						Ticker: types.Ticker{
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
				TickerConfigs: map[string]types.TickerConfig{},
			},
			expErr: true,
		},
		{
			name: "invalid ticker config",
			market: types.MarketConfig{
				Name: "binance",
				TickerConfigs: map[string]types.TickerConfig{
					"BITCOIN/USDT": {
						Ticker:         types.Ticker{},
						OffChainTicker: "BTC/USDT",
					},
				},
			},
			expErr: true,
		},
		{
			name: "missing ticker config key",
			market: types.MarketConfig{
				Name: "binance",
				TickerConfigs: map[string]types.TickerConfig{
					"BITCOIN/USDC": {
						Ticker: types.Ticker{
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
			name: "duplicate on-chain ticker",
			market: types.MarketConfig{
				Name: "binance",
				TickerConfigs: map[string]types.TickerConfig{
					"BITCOIN/USDT": {
						Ticker: types.Ticker{
							Base:             "BITCOIN",
							Quote:            "USDT",
							Decimals:         8,
							MinProviderCount: 1,
						},
						OffChainTicker: "BTC/USDT",
					},
					"BITCOIN/USDC": {
						Ticker: types.Ticker{
							Base:             "BITCOIN",
							Quote:            "USDT",
							Decimals:         8,
							MinProviderCount: 1,
						},
						OffChainTicker: "BTC/USDC",
					},
				},
			},
			expErr: true,
		},
		{
			name: "duplicate off-chain ticker",
			market: types.MarketConfig{
				Name: "binance",
				TickerConfigs: map[string]types.TickerConfig{
					"BITCOIN/USDT": {
						Ticker: types.Ticker{
							Base:             "BITCOIN",
							Quote:            "USDT",
							Decimals:         8,
							MinProviderCount: 1,
						},
						OffChainTicker: "BTC/USDT",
					},
					"BITCOIN/USDC": {
						Ticker: types.Ticker{
							Base:             "BITCOIN",
							Quote:            "USDC",
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
