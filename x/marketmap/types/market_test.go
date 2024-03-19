package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/x/marketmap/types"
)

func TestMarketMapValidateBasic(t *testing.T) {
	testCases := []struct {
		name      string
		marketMap types.MarketMap
		expectErr bool
	}{
		{
			name:      "valid empty",
			marketMap: types.MarketMap{},
			expectErr: false,
		},
		{
			name: "valid map",
			marketMap: types.MarketMap{
				Tickers: map[string]types.Ticker{
					ethusdt.String(): ethusdt,
					btcusdt.String(): btcusdt,
					usdcusd.String(): usdcusd,
				},
				Paths: map[string]types.Paths{
					ethusdt.String(): ethusdtPaths,
					btcusdt.String(): btcusdtPaths,
					usdcusd.String(): usdcusdPaths,
				},
				Providers: map[string]types.Providers{
					ethusdt.String(): ethusdtProviders,
					btcusdt.String(): btcusdtProviders,
					usdcusd.String(): usdcusdProviders,
				},
			},
			expectErr: false,
		},
		{
			name: "invalid mismatch ticker",
			marketMap: types.MarketMap{
				Tickers: map[string]types.Ticker{
					ethusdt.String(): ethusdt,
					btcusdt.String(): btcusdt,
					usdcusd.String(): usdcusd,
				},
				Paths: map[string]types.Paths{
					ethusdt.String(): ethusdtPaths,
					btcusdt.String(): btcusdtPaths,
					usdcusd.String(): usdcusdPaths,
				},
				Providers: map[string]types.Providers{
					usdtusd.String(): usdtusdProviders,
					btcusdt.String(): btcusdtProviders,
					usdcusd.String(): usdcusdProviders,
				},
			},
			expectErr: true,
		},
		{
			name: "invalid ticker does not exist for a given provider",
			marketMap: types.MarketMap{
				Tickers: map[string]types.Ticker{
					btcusdt.String(): btcusdt,
					usdcusd.String(): usdcusd,
				},
				Paths: map[string]types.Paths{
					btcusdt.String(): btcusdtPaths,
					usdcusd.String(): usdcusdPaths,
				},
				Providers: map[string]types.Providers{
					ethusdt.String(): ethusdtProviders,
					btcusdt.String(): btcusdtProviders,
					usdcusd.String(): usdcusdProviders,
				},
			},
			expectErr: true,
		},
		{
			name: "invalid ticker string does not match ticker ID",
			marketMap: types.MarketMap{
				Tickers: map[string]types.Ticker{
					"invalid":        ethusdt,
					btcusdt.String(): btcusdt,
					usdcusd.String(): usdcusd,
				},
				Paths: map[string]types.Paths{
					ethusdt.String(): ethusdtPaths,
					btcusdt.String(): btcusdtPaths,
					usdcusd.String(): usdcusdPaths,
				},
				Providers: map[string]types.Providers{
					ethusdt.String(): ethusdtProviders,
					btcusdt.String(): btcusdtProviders,
					usdcusd.String(): usdcusdProviders,
				},
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.marketMap.ValidateBasic()
			if tc.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestValidateMarketMap(t *testing.T) {
	testCases := []struct {
		name string
		cfg  types.MarketMap
		err  bool
	}{
		{
			name: "empty market map",
			cfg: types.MarketMap{
				AggregationType: types.AggregationType_INDEX_PRICE_AGGREGATION,
			},
			err: false,
		},
		{
			name: "provider includes a ticker that is not supported by the main set",
			cfg: types.MarketMap{
				Providers: map[string]types.Providers{
					constants.BITCOIN_USD.String(): {
						Providers: []types.ProviderConfig{
							{
								Name:           coinbase.Name,
								OffChainTicker: "BTC-USD",
							},
						},
					},
				},
				AggregationType: types.AggregationType_INDEX_PRICE_AGGREGATION,
			},
			err: true,
		},
		{
			name: "provider includes a ticker that is supported by the main set - no paths",
			cfg: types.MarketMap{
				Tickers: map[string]types.Ticker{
					constants.BITCOIN_USD.String(): constants.BITCOIN_USD,
				},
				Providers: map[string]types.Providers{
					constants.BITCOIN_USD.String(): {
						Providers: []types.ProviderConfig{
							{
								Name:           coinbase.Name,
								OffChainTicker: "BTC-USD",
							},
						},
					},
				},
				AggregationType: types.AggregationType_INDEX_PRICE_AGGREGATION,
			},
			err: false,
		},
		{
			name: "path includes a ticker that is not supported",
			cfg: types.MarketMap{
				Paths: map[string]types.Paths{
					constants.BITCOIN_USD.String(): {},
				},
				AggregationType: types.AggregationType_INDEX_PRICE_AGGREGATION,
			},
			err: true,
		},
		{
			name: "paths includes a path that has no operations",
			cfg: types.MarketMap{
				Tickers: map[string]types.Ticker{
					constants.BITCOIN_USD.String(): constants.BITCOIN_USD,
				},
				Paths: map[string]types.Paths{
					constants.BITCOIN_USD.String(): {
						Paths: []types.Path{
							{
								Operations: []types.Operation{},
							},
						},
					},
				},
				AggregationType: types.AggregationType_INDEX_PRICE_AGGREGATION,
			},
			err: true,
		},
		{
			name: "paths includes a path that has too many operations",
			cfg: types.MarketMap{
				Tickers: map[string]types.Ticker{
					constants.BITCOIN_USD.String(): constants.BITCOIN_USD,
				},
				Paths: map[string]types.Paths{
					constants.BITCOIN_USD.String(): {
						Paths: []types.Path{
							{
								Operations: []types.Operation{
									{
										CurrencyPair: constants.BITCOIN_USD.CurrencyPair,
										Provider:     coinbase.Name,
									},
									{
										CurrencyPair: constants.BITCOIN_USD.CurrencyPair,
										Provider:     coinbase.Name,
									},
									{
										CurrencyPair: constants.BITCOIN_USD.CurrencyPair,
										Provider:     coinbase.Name,
									},
								},
							},
						},
					},
				},
				AggregationType: types.AggregationType_INDEX_PRICE_AGGREGATION,
			},
			err: true,
		},
		{
			name: "operation includes a ticker that is not supported",
			cfg: types.MarketMap{
				Tickers: map[string]types.Ticker{
					constants.BITCOIN_USD.String(): constants.BITCOIN_USD,
				},
				Paths: map[string]types.Paths{
					constants.BITCOIN_USD.String(): {
						Paths: []types.Path{
							{
								Operations: []types.Operation{
									{
										CurrencyPair: constants.BITCOIN_USDT.CurrencyPair,
										Provider:     coinbase.Name,
									},
								},
							},
						},
					},
				},
				AggregationType: types.AggregationType_INDEX_PRICE_AGGREGATION,
			},
			err: true,
		},
		{
			name: "operation includes a provider that does not support the ticker",
			cfg: types.MarketMap{
				Tickers: map[string]types.Ticker{
					constants.BITCOIN_USD.String(): constants.BITCOIN_USD,
				},
				Paths: map[string]types.Paths{
					constants.BITCOIN_USD.String(): {
						Paths: []types.Path{
							{
								Operations: []types.Operation{
									{
										CurrencyPair: constants.BITCOIN_USD.CurrencyPair,
										Provider:     coinbase.Name,
									},
								},
							},
						},
					},
				},
				AggregationType: types.AggregationType_INDEX_PRICE_AGGREGATION,
			},
			err: true,
		},
		{
			name: "provider does not support a ticker included in an operation",
			cfg: types.MarketMap{
				Tickers: map[string]types.Ticker{
					constants.BITCOIN_USD.String(): constants.BITCOIN_USD,
				},
				Providers: map[string]types.Providers{
					constants.BITCOIN_USD.String(): {
						Providers: []types.ProviderConfig{},
					},
				},
				Paths: map[string]types.Paths{
					constants.BITCOIN_USD.String(): {
						Paths: []types.Path{
							{
								Operations: []types.Operation{
									{
										CurrencyPair: constants.BITCOIN_USD.CurrencyPair,
										Provider:     coinbase.Name,
									},
								},
							},
						},
					},
				},
				AggregationType: types.AggregationType_INDEX_PRICE_AGGREGATION,
			},
			err: true,
		},
		{
			name: "valid single path",
			cfg: types.MarketMap{
				Tickers: map[string]types.Ticker{
					constants.BITCOIN_USD.String(): constants.BITCOIN_USD,
				},
				Providers: map[string]types.Providers{
					constants.BITCOIN_USD.String(): {
						Providers: []types.ProviderConfig{
							{
								Name:           coinbase.Name,
								OffChainTicker: "BTC-USD",
							},
						},
					},
				},
				Paths: map[string]types.Paths{
					constants.BITCOIN_USD.String(): {
						Paths: []types.Path{
							{
								Operations: []types.Operation{
									{
										CurrencyPair: constants.BITCOIN_USD.CurrencyPair,
										Provider:     coinbase.Name,
									},
								},
							},
						},
					},
				},
				AggregationType: types.AggregationType_INDEX_PRICE_AGGREGATION,
			},
			err: false,
		},
		{
			name: "path includes a index ticker that is not supported",
			cfg: types.MarketMap{
				Tickers: map[string]types.Ticker{
					constants.BITCOIN_USDT.String(): constants.BITCOIN_USDT,
					constants.BITCOIN_USD.String():  constants.BITCOIN_USD,
				},
				Providers: map[string]types.Providers{
					constants.BITCOIN_USDT.String(): {
						Providers: []types.ProviderConfig{
							{
								Name:           coinbase.Name,
								OffChainTicker: "BTC-USDT",
							},
						},
					},
				},
				Paths: map[string]types.Paths{
					constants.BITCOIN_USD.String(): {
						Paths: []types.Path{
							{
								Operations: []types.Operation{
									{
										CurrencyPair: constants.BITCOIN_USDT.CurrencyPair,
										Provider:     coinbase.Name,
									},
									{
										CurrencyPair: constants.USDT_USD.CurrencyPair,
										Provider:     types.IndexPrice,
									},
								},
							},
						},
					},
				},
				AggregationType: types.AggregationType_INDEX_PRICE_AGGREGATION,
			},
			err: true,
		},
		{
			name: "second operation is not an index price provider",
			cfg: types.MarketMap{
				Tickers: map[string]types.Ticker{
					constants.BITCOIN_USDT.String(): constants.BITCOIN_USDT,
					constants.BITCOIN_USD.String():  constants.BITCOIN_USD,
					constants.USDT_USD.String():     constants.USDT_USD,
				},
				Providers: map[string]types.Providers{
					constants.BITCOIN_USDT.String(): {
						Providers: []types.ProviderConfig{
							{
								Name:           coinbase.Name,
								OffChainTicker: "BTC-USDT",
							},
						},
					},
				},
				Paths: map[string]types.Paths{
					constants.BITCOIN_USD.String(): {
						Paths: []types.Path{
							{
								Operations: []types.Operation{
									{
										CurrencyPair: constants.BITCOIN_USDT.CurrencyPair,
										Provider:     coinbase.Name,
									},
									{
										CurrencyPair: constants.USDT_USD.CurrencyPair,
										Provider:     coinbase.Name,
									},
								},
							},
						},
					},
				},
				AggregationType: types.AggregationType_INDEX_PRICE_AGGREGATION,
			},
			err: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cfg.ValidateBasic()
			if tc.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMarketMapEqual(t *testing.T) {
	cases := []struct {
		name      string
		marketMap types.MarketMap
		other     types.MarketMap
		expect    bool
	}{
		{
			name:      "empty market map",
			marketMap: types.MarketMap{},
			other:     types.MarketMap{},
			expect:    true,
		},
		{
			name: "same market map",
			marketMap: types.MarketMap{
				Tickers: map[string]types.Ticker{
					ethusdt.String(): ethusdt,
				},
				Paths: map[string]types.Paths{
					ethusdt.String(): ethusdtPaths,
				},
				Providers: map[string]types.Providers{
					ethusdt.String(): ethusdtProviders,
				},
			},
			other: types.MarketMap{
				Tickers: map[string]types.Ticker{
					ethusdt.String(): ethusdt,
				},
				Paths: map[string]types.Paths{
					ethusdt.String(): ethusdtPaths,
				},
				Providers: map[string]types.Providers{
					ethusdt.String(): ethusdtProviders,
				},
			},
			expect: true,
		},
		{
			name: "different tickers",
			marketMap: types.MarketMap{
				Tickers: map[string]types.Ticker{
					ethusdt.String(): ethusdt,
				},
				Paths: map[string]types.Paths{
					ethusdt.String(): ethusdtPaths,
				},
				Providers: map[string]types.Providers{
					ethusdt.String(): ethusdtProviders,
				},
			},
			other: types.MarketMap{
				Tickers: map[string]types.Ticker{
					btcusdt.String(): btcusdt,
				},
				Paths: map[string]types.Paths{
					btcusdt.String(): btcusdtPaths,
				},
				Providers: map[string]types.Providers{
					btcusdt.String(): btcusdtProviders,
				},
			},
			expect: false,
		},
		{
			name: "different paths",
			marketMap: types.MarketMap{
				Tickers: map[string]types.Ticker{
					ethusdt.String(): ethusdt,
				},
				Paths: map[string]types.Paths{
					ethusdt.String(): ethusdtPaths,
				},
				Providers: map[string]types.Providers{
					ethusdt.String(): ethusdtProviders,
				},
			},
			other: types.MarketMap{
				Tickers: map[string]types.Ticker{
					ethusdt.String(): ethusdt,
				},
				Paths: map[string]types.Paths{
					btcusdt.String(): btcusdtPaths,
				},
				Providers: map[string]types.Providers{
					btcusdt.String(): btcusdtProviders,
				},
			},
			expect: false,
		},
		{
			name: "different providers",
			marketMap: types.MarketMap{
				Tickers: map[string]types.Ticker{
					ethusdt.String(): ethusdt,
				},
				Paths: map[string]types.Paths{
					ethusdt.String(): ethusdtPaths,
				},
				Providers: map[string]types.Providers{
					ethusdt.String(): ethusdtProviders,
				},
			},
			other: types.MarketMap{
				Tickers: map[string]types.Ticker{
					ethusdt.String(): ethusdt,
				},
				Paths: map[string]types.Paths{
					ethusdt.String(): ethusdtPaths,
				},
				Providers: map[string]types.Providers{
					btcusdt.String(): btcusdtProviders,
				},
			},
			expect: false,
		},

		{
			name: "different aggregation type",
			marketMap: types.MarketMap{
				Tickers: map[string]types.Ticker{
					ethusdt.String(): ethusdt,
				},
				Paths: map[string]types.Paths{
					ethusdt.String(): ethusdtPaths,
				},
				Providers: map[string]types.Providers{
					ethusdt.String(): ethusdtProviders,
				},
				AggregationType: types.AggregationType_INDEX_PRICE_AGGREGATION,
			},
			other: types.MarketMap{
				Tickers: map[string]types.Ticker{
					ethusdt.String(): ethusdt,
				},
				Paths: map[string]types.Paths{
					ethusdt.String(): ethusdtPaths,
				},
				Providers: map[string]types.Providers{
					ethusdt.String(): ethusdtProviders,
				},
			},
			expect: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expect, tc.marketMap.Equal(tc.other))
		})
	}
}
