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
				Markets: markets,
			},
			expectErr: false,
		},
		{
			name: "invalid ticker string does not match ticker ID",
			marketMap: types.MarketMap{
				Markets: map[string]types.Market{
					"invalid": {
						Ticker:    ethusdt.Ticker,
						Paths:     ethusdt.Paths,
						Providers: ethusdt.Providers,
					},
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
				Markets: map[string]types.Market{
					constants.BITCOIN_USD.String(): {
						Providers: types.Providers{
							Providers: []types.ProviderConfig{
								{
									Name:           coinbase.Name,
									OffChainTicker: "BTC-USD",
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
			name: "provider includes a ticker that is supported by the main set - no paths",
			cfg: types.MarketMap{
				Markets: map[string]types.Market{
					constants.BITCOIN_USD.String(): {
						Ticker: constants.BITCOIN_USD,
						Providers: types.Providers{
							Providers: []types.ProviderConfig{
								{
									Name:           coinbase.Name,
									OffChainTicker: "BTC-USD",
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
			name: "includes a ticker that is not supported",
			cfg: types.MarketMap{
				Markets: map[string]types.Market{
					constants.BITCOIN_USD.String(): {},
				},
				AggregationType: types.AggregationType_INDEX_PRICE_AGGREGATION,
			},
			err: true,
		},
		{
			name: "paths includes a path that has no operations",
			cfg: types.MarketMap{
				Markets: map[string]types.Market{
					constants.BITCOIN_USD.String(): {
						Ticker: constants.BITCOIN_USD,
						Paths: types.Paths{
							Paths: []types.Path{
								{
									Operations: []types.Operation{},
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
			name: "paths includes a path that has too many operations",
			cfg: types.MarketMap{
				Markets: map[string]types.Market{
					constants.BITCOIN_USD.String(): {
						Ticker: constants.BITCOIN_USD,
						Paths: types.Paths{
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
				},
				AggregationType: types.AggregationType_INDEX_PRICE_AGGREGATION,
			},
			err: true,
		},
		{
			name: "operation includes a ticker that is not supported",
			cfg: types.MarketMap{
				Markets: map[string]types.Market{
					constants.BITCOIN_USD.String(): {
						Ticker: constants.BITCOIN_USD,
						Paths: types.Paths{
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
						Providers: types.Providers{},
					},
				},

				AggregationType: types.AggregationType_INDEX_PRICE_AGGREGATION,
			},
			err: true,
		},
		{
			name: "operation includes a provider that does not support the ticker",
			cfg: types.MarketMap{
				Markets: map[string]types.Market{
					constants.BITCOIN_USD.String(): {
						Ticker: constants.BITCOIN_USD,
						Paths: types.Paths{
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
				},
				AggregationType: types.AggregationType_INDEX_PRICE_AGGREGATION,
			},
			err: true,
		},
		{
			name: "provider does not support a ticker included in an operation",
			cfg: types.MarketMap{
				Markets: map[string]types.Market{
					constants.BITCOIN_USD.String(): {
						Ticker: constants.BITCOIN_USD,
						Providers: types.Providers{
							Providers: []types.ProviderConfig{},
						},
						Paths: types.Paths{
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
				},
				AggregationType: types.AggregationType_INDEX_PRICE_AGGREGATION,
			},
			err: true,
		},
		{
			name: "valid single path",
			cfg: types.MarketMap{
				Markets: map[string]types.Market{
					constants.BITCOIN_USD.String(): {
						Ticker: constants.BITCOIN_USD,
						Providers: types.Providers{
							Providers: []types.ProviderConfig{
								{
									Name:           coinbase.Name,
									OffChainTicker: "BTC-USD",
								},
							},
						},
						Paths: types.Paths{
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
				},
				AggregationType: types.AggregationType_INDEX_PRICE_AGGREGATION,
			},
			err: false,
		},
		{
			name: "path includes a index ticker that is not supported",
			cfg: types.MarketMap{
				Markets: map[string]types.Market{
					constants.BITCOIN_USDT.String(): {
						Ticker: constants.BITCOIN_USDT,
						Providers: types.Providers{
							Providers: []types.ProviderConfig{
								{
									Name:           coinbase.Name,
									OffChainTicker: "BTC-USDT",
								},
							},
						},
					},

					constants.BITCOIN_USD.String(): {
						Ticker: constants.BITCOIN_USD,
						Paths: types.Paths{
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
				},
				AggregationType: types.AggregationType_INDEX_PRICE_AGGREGATION,
			},
			err: true,
		},
		{
			name: "second operation is not an index price provider",
			cfg: types.MarketMap{
				Markets: map[string]types.Market{
					constants.BITCOIN_USDT.String(): {
						Ticker: constants.BITCOIN_USD,
						Providers: types.Providers{
							Providers: []types.ProviderConfig{
								{
									Name:           coinbase.Name,
									OffChainTicker: "BTC-USDT",
								},
							},
						},
						Paths: types.Paths{
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
					constants.BITCOIN_USD.String(): {},
					constants.USDT_USD.String():    {},
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
