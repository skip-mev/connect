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
		{
			name: "single provider with single ticker",
			markets: map[string]types.MarketConfig{
				"binance": {
					Name: "binance",
					TickerConfigs: map[string]types.TickerConfig{
						"BITCOIN/USDT": {
							Ticker:         btcusdt,
							OffChainTicker: "BTC/USDT",
						},
					},
				},
			},
			tickers: map[string]types.PathsConfig{
				"BITCOIN/USDT": {
					Ticker: btcusdt,
					Paths: []types.Path{
						{
							Operations: []types.Operation{
								{
									Ticker: btcusdt,
								},
							},
						},
					},
				},
			},
			expErr: false,
		},
		{
			name: "single provider with an unsupported ticker",
			markets: map[string]types.MarketConfig{
				"binance": {
					Name: "binance",
					TickerConfigs: map[string]types.TickerConfig{
						"BITCOIN/USDT": {
							Ticker:         btcusdt,
							OffChainTicker: "BTC/USDT",
						},
					},
				},
			},
			tickers: map[string]types.PathsConfig{
				"ETHEREUM/USDT": {
					Ticker: ethusdt,
					Paths: []types.Path{
						{
							Operations: []types.Operation{
								{
									Ticker: ethusdt,
								},
							},
						},
					},
				},
			},
			expErr: true,
		},
		{
			name: "single bad provider market config",
			markets: map[string]types.MarketConfig{
				"binance": {},
			},
			tickers: map[string]types.PathsConfig{},
			expErr:  true,
		},
		{
			name: "single bad provider with mismatching provider name",
			markets: map[string]types.MarketConfig{
				"binance": {
					Name: "coinbase",
					TickerConfigs: map[string]types.TickerConfig{
						"BITCOIN/USDT": {
							Ticker:         btcusdt,
							OffChainTicker: "BTC/USDT",
						},
					},
				},
			},
			tickers: map[string]types.PathsConfig{},
			expErr:  true,
		},
		{
			name: "1 good provider and a bad ticker config",
			markets: map[string]types.MarketConfig{
				"binance": {
					Name: "binance",
					TickerConfigs: map[string]types.TickerConfig{
						"BITCOIN/USDT": {
							Ticker:         btcusdt,
							OffChainTicker: "BTC/USDT",
						},
					},
				},
			},
			tickers: map[string]types.PathsConfig{
				"BITCOIN/USDT": {
					Ticker: types.Ticker{},
					Paths:  []types.Path{},
				},
			},
			expErr: true,
		},
		{
			name: "1 good provider and a mismatching ticker config",
			markets: map[string]types.MarketConfig{
				"binance": {
					Name: "binance",
					TickerConfigs: map[string]types.TickerConfig{
						"BITCOIN/USDT": {
							Ticker:         btcusdt,
							OffChainTicker: "BTC/USDT",
						},
					},
				},
			},
			tickers: map[string]types.PathsConfig{
				"ETHEREUM/USDT": {
					Ticker: btcusdt,
					Paths: []types.Path{
						{
							Operations: []types.Operation{
								{
									Ticker: btcusdt,
								},
							},
						},
					},
				},
			},
			expErr: true,
		},
		{
			name: "1 good provider but no support for conversion path",
			markets: map[string]types.MarketConfig{
				"binance": {
					Name: "binance",
					TickerConfigs: map[string]types.TickerConfig{
						"BITCOIN/USDT": {
							Ticker:         btcusdt,
							OffChainTicker: "BTC/USDT",
						},
					},
				},
			},
			tickers: map[string]types.PathsConfig{
				"BITCOIN/USD": {
					Ticker: types.Ticker{
						Base:             "BITCOIN",
						Quote:            "USD",
						Decimals:         8,
						MinProviderCount: 1,
					},
					Paths: []types.Path{
						{
							Operations: []types.Operation{
								{
									Ticker: types.Ticker{
										Base:             "BITCOIN",
										Quote:            "USDT",
										Decimals:         8,
										MinProviderCount: 1,
									},
								},
								{
									Ticker: types.Ticker{
										Base:             "USDT",
										Quote:            "USD",
										Decimals:         8,
										MinProviderCount: 1,
									},
								},
							},
						},
					},
				},
			},
			expErr: true,
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
