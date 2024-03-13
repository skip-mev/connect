package oracle_test

import (
	"testing"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math/oracle"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	"github.com/stretchr/testify/require"
)

func TestGetTickerFromOperation(t *testing.T) {
	t.Run("has ticker included in the market config", func(t *testing.T) {
		m, err := oracle.NewMedianAggregator(logger, marketmap)
		require.NoError(t, err)

		operation := mmtypes.Operation{
			CurrencyPair: BTC_USD.CurrencyPair,
		}
		ticker, err := m.GetTickerFromOperation(operation)
		require.NoError(t, err)
		require.Equal(t, BTC_USD, ticker)
	})

	t.Run("has ticker not included in the market config", func(t *testing.T) {
		m, err := oracle.NewMedianAggregator(logger, marketmap)
		require.NoError(t, err)

		operation := mmtypes.Operation{
			CurrencyPair: constants.MOG_USD.CurrencyPair,
		}
		ticker, err := m.GetTickerFromOperation(operation)
		require.Error(t, err)
		require.Empty(t, ticker)
	})
}

func TestGetProviderPrice(t *testing.T) {
	t.Run("does not have a ticker in the config", func(t *testing.T) {
		m, err := oracle.NewMedianAggregator(logger, marketmap)
		require.NoError(t, err)

		operation := mmtypes.Operation{
			CurrencyPair: constants.MOG_USD.CurrencyPair,
		}
		_, err = m.GetProviderPrice(operation)
		require.Error(t, err)
	})

	t.Run("has no provider prices or index prices", func(t *testing.T) {
		m, err := oracle.NewMedianAggregator(logger, marketmap)
		require.NoError(t, err)

		// Attempt to retrieve the provider.
		operation := mmtypes.Operation{
			CurrencyPair: BTC_USD.CurrencyPair,
			Provider:     coinbase.Name,
		}
		_, err = m.GetProviderPrice(operation)
		require.Error(t, err)

		// Attempt to retrieve the index price.
		operation = mmtypes.Operation{
			CurrencyPair: BTC_USD.CurrencyPair,
			Provider:     oracle.IndexPrice,
		}
		_, err = m.GetProviderPrice(operation)
		require.Error(t, err)
	})

	t.Run("has provider prices but no index prices", func(t *testing.T) {
		m, err := oracle.NewMedianAggregator(logger, marketmap)
		require.NoError(t, err)

		// Set the provider price.
		prices := types.TickerPrices{
			BTC_USD: createPrice(100, BTC_USD.Decimals),
		}
		m.PriceAggregator.SetProviderData(coinbase.Name, prices)

		// Attempt to retrieve the provider.
		operation := mmtypes.Operation{
			CurrencyPair: BTC_USD.CurrencyPair,
			Provider:     coinbase.Name,
		}
		price, err := m.GetProviderPrice(operation)
		require.NoError(t, err)
		require.Equal(t, createPrice(100, oracle.ScaledDecimals), price)

		// Attempt to retrieve the index price.
		operation = mmtypes.Operation{
			CurrencyPair: BTC_USD.CurrencyPair,
			Provider:     oracle.IndexPrice,
		}
		_, err = m.GetProviderPrice(operation)
		require.Error(t, err)
	})

	t.Run("has provider prices and index prices", func(t *testing.T) {
		m, err := oracle.NewMedianAggregator(logger, marketmap)
		require.NoError(t, err)

		// Set the provider price.
		prices := types.TickerPrices{
			BTC_USD: createPrice(100, BTC_USD.Decimals),
		}
		m.PriceAggregator.SetProviderData(coinbase.Name, prices)

		// Set the index price.
		m.PriceAggregator.SetAggregatedData(prices)

		// Attempt to retrieve the provider.
		operation := mmtypes.Operation{
			CurrencyPair: BTC_USD.CurrencyPair,
			Provider:     coinbase.Name,
		}
		price, err := m.GetProviderPrice(operation)
		require.NoError(t, err)
		require.Equal(t, createPrice(100, oracle.ScaledDecimals), price)

		// Attempt to retrieve the index price.
		operation = mmtypes.Operation{
			CurrencyPair: BTC_USD.CurrencyPair,
			Provider:     oracle.IndexPrice,
		}
		price, err = m.GetProviderPrice(operation)
		require.NoError(t, err)
		require.Equal(t, createPrice(100, oracle.ScaledDecimals), price)
	})

	t.Run("has provider prices and can correctly scale up", func(t *testing.T) {
		m, err := oracle.NewMedianAggregator(logger, marketmap)
		require.NoError(t, err)

		// Set the provider price.
		prices := types.TickerPrices{
			BTC_USD: createPrice(40_000, BTC_USD.Decimals),
		}
		m.PriceAggregator.SetProviderData(coinbase.Name, prices)

		// Attempt to retrieve the provider.
		operation := mmtypes.Operation{
			CurrencyPair: BTC_USD.CurrencyPair,
			Provider:     coinbase.Name,
		}
		price, err := m.GetProviderPrice(operation)
		require.NoError(t, err)
		require.Equal(t, createPrice(40_000, oracle.ScaledDecimals), price)
	})

	t.Run("has provider prices and can correctly invert", func(t *testing.T) {
		m, err := oracle.NewMedianAggregator(logger, marketmap)
		require.NoError(t, err)

		// Set the provider price.
		prices := types.TickerPrices{
			BTC_USD: createPrice(40_000, BTC_USD.Decimals),
		}
		m.PriceAggregator.SetProviderData(coinbase.Name, prices)

		// Attempt to retrieve the provider.
		operation := mmtypes.Operation{
			CurrencyPair: BTC_USD.CurrencyPair,
			Provider:     coinbase.Name,
			Invert:       true,
		}
		price, err := m.GetProviderPrice(operation)
		require.NoError(t, err)
		expectedPrice := createPrice(0.000025, oracle.ScaledDecimals)
		verifyPrice(t, expectedPrice, price)
	})
}

func TestValidateMarketMap(t *testing.T) {
	testCases := []struct {
		name string
		cfg  mmtypes.MarketMap
		err  bool
	}{
		{
			name: "empty market map",
			cfg:  mmtypes.MarketMap{},
			err:  false,
		},
		{
			name: "provider includes a ticker that is not supported by the main set",
			cfg: mmtypes.MarketMap{
				Providers: map[string]mmtypes.Providers{
					BTC_USD.String(): {
						Providers: []mmtypes.ProviderConfig{
							{
								Name:           coinbase.Name,
								OffChainTicker: "BTC-USD",
							},
						},
					},
				},
			},
			err: true,
		},
		{
			name: "provider includes a ticker that is supported by the main set - no paths",
			cfg: mmtypes.MarketMap{
				Tickers: map[string]mmtypes.Ticker{
					BTC_USD.String(): BTC_USD,
				},
				Providers: map[string]mmtypes.Providers{
					BTC_USD.String(): {
						Providers: []mmtypes.ProviderConfig{
							{
								Name:           coinbase.Name,
								OffChainTicker: "BTC-USD",
							},
						},
					},
				},
			},
			err: false,
		},
		{
			name: "path includes a ticker that is not supported",
			cfg: mmtypes.MarketMap{
				Paths: map[string]mmtypes.Paths{
					BTC_USD.String(): {},
				},
			},
			err: true,
		},
		{
			name: "paths includes a path that has no operations",
			cfg: mmtypes.MarketMap{
				Tickers: map[string]mmtypes.Ticker{
					BTC_USD.String(): BTC_USD,
				},
				Paths: map[string]mmtypes.Paths{
					BTC_USD.String(): {
						Paths: []mmtypes.Path{
							{
								Operations: []mmtypes.Operation{},
							},
						},
					},
				},
			},
			err: true,
		},
		{
			name: "paths includes a path that has too many operations",
			cfg: mmtypes.MarketMap{
				Tickers: map[string]mmtypes.Ticker{
					BTC_USD.String(): BTC_USD,
				},
				Paths: map[string]mmtypes.Paths{
					BTC_USD.String(): {
						Paths: []mmtypes.Path{
							{
								Operations: []mmtypes.Operation{
									{
										CurrencyPair: BTC_USD.CurrencyPair,
										Provider:     coinbase.Name,
									},
									{
										CurrencyPair: BTC_USD.CurrencyPair,
										Provider:     coinbase.Name,
									},
									{
										CurrencyPair: BTC_USD.CurrencyPair,
										Provider:     coinbase.Name,
									},
								},
							},
						},
					},
				},
			},
			err: true,
		},
		{
			name: "operation includes a ticker that is not supported",
			cfg: mmtypes.MarketMap{
				Tickers: map[string]mmtypes.Ticker{
					BTC_USD.String(): BTC_USD,
				},
				Paths: map[string]mmtypes.Paths{
					BTC_USD.String(): {
						Paths: []mmtypes.Path{
							{
								Operations: []mmtypes.Operation{
									{
										CurrencyPair: BTC_USDT.CurrencyPair,
										Provider:     coinbase.Name,
									},
								},
							},
						},
					},
				},
			},
			err: true,
		},
		{
			name: "operation includes a provider that does not support the ticker",
			cfg: mmtypes.MarketMap{
				Tickers: map[string]mmtypes.Ticker{
					BTC_USD.String(): BTC_USD,
				},
				Paths: map[string]mmtypes.Paths{
					BTC_USD.String(): {
						Paths: []mmtypes.Path{
							{
								Operations: []mmtypes.Operation{
									{
										CurrencyPair: BTC_USD.CurrencyPair,
										Provider:     coinbase.Name,
									},
								},
							},
						},
					},
				},
			},
			err: true,
		},
		{
			name: "provider does not support a ticker included in an operation",
			cfg: mmtypes.MarketMap{
				Tickers: map[string]mmtypes.Ticker{
					BTC_USD.String(): BTC_USD,
				},
				Providers: map[string]mmtypes.Providers{
					BTC_USD.String(): {
						Providers: []mmtypes.ProviderConfig{},
					},
				},
				Paths: map[string]mmtypes.Paths{
					BTC_USD.String(): {
						Paths: []mmtypes.Path{
							{
								Operations: []mmtypes.Operation{
									{
										CurrencyPair: BTC_USD.CurrencyPair,
										Provider:     coinbase.Name,
									},
								},
							},
						},
					},
				},
			},
			err: true,
		},
		{
			name: "valid single path",
			cfg: mmtypes.MarketMap{
				Tickers: map[string]mmtypes.Ticker{
					BTC_USD.String(): BTC_USD,
				},
				Providers: map[string]mmtypes.Providers{
					BTC_USD.String(): {
						Providers: []mmtypes.ProviderConfig{
							{
								Name:           coinbase.Name,
								OffChainTicker: "BTC-USD",
							},
						},
					},
				},
				Paths: map[string]mmtypes.Paths{
					BTC_USD.String(): {
						Paths: []mmtypes.Path{
							{
								Operations: []mmtypes.Operation{
									{
										CurrencyPair: BTC_USD.CurrencyPair,
										Provider:     coinbase.Name,
									},
								},
							},
						},
					},
				},
			},
			err: false,
		},
		{
			name: "path includes a index ticker that is not supported",
			cfg: mmtypes.MarketMap{
				Tickers: map[string]mmtypes.Ticker{
					BTC_USDT.String(): BTC_USDT,
					BTC_USD.String():  BTC_USD,
				},
				Providers: map[string]mmtypes.Providers{
					BTC_USDT.String(): {
						Providers: []mmtypes.ProviderConfig{
							{
								Name:           coinbase.Name,
								OffChainTicker: "BTC-USDT",
							},
						},
					},
				},
				Paths: map[string]mmtypes.Paths{
					BTC_USD.String(): {
						Paths: []mmtypes.Path{
							{
								Operations: []mmtypes.Operation{
									{
										CurrencyPair: BTC_USDT.CurrencyPair,
										Provider:     coinbase.Name,
									},
									{
										CurrencyPair: USDT_USD.CurrencyPair,
										Provider:     oracle.IndexPrice,
									},
								},
							},
						},
					},
				},
			},
			err: true,
		},
		{
			name: "second operation is not an index price provider",
			cfg: mmtypes.MarketMap{
				Tickers: map[string]mmtypes.Ticker{
					BTC_USDT.String(): BTC_USDT,
					BTC_USD.String():  BTC_USD,
					USDT_USD.String(): USDT_USD,
				},
				Providers: map[string]mmtypes.Providers{
					BTC_USDT.String(): {
						Providers: []mmtypes.ProviderConfig{
							{
								Name:           coinbase.Name,
								OffChainTicker: "BTC-USDT",
							},
						},
					},
				},
				Paths: map[string]mmtypes.Paths{
					BTC_USD.String(): {
						Paths: []mmtypes.Path{
							{
								Operations: []mmtypes.Operation{
									{
										CurrencyPair: BTC_USDT.CurrencyPair,
										Provider:     coinbase.Name,
									},
									{
										CurrencyPair: USDT_USD.CurrencyPair,
										Provider:     coinbase.Name,
									},
								},
							},
						},
					},
				},
			},
			err: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := oracle.ValidateMarketMap(tc.cfg)
			if tc.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
