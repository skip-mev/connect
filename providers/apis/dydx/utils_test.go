package dydx_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/constants"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	coinbaseapi "github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/apis/dydx"
	dydxtypes "github.com/skip-mev/slinky/providers/apis/dydx/types"
	coinbasews "github.com/skip-mev/slinky/providers/websockets/coinbase"
	"github.com/skip-mev/slinky/providers/websockets/kucoin"
	"github.com/skip-mev/slinky/providers/websockets/okx"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

func TestConvertMarketParamsToMarketMap(t *testing.T) {
	testCases := []struct {
		name     string
		params   dydxtypes.QueryAllMarketParamsResponse
		expected mmtypes.GetMarketMapResponse
		err      bool
	}{
		{
			name:   "empty market params",
			params: dydxtypes.QueryAllMarketParamsResponse{},
			expected: mmtypes.GetMarketMapResponse{
				MarketMap: mmtypes.MarketMap{
					Markets:         make(map[string]mmtypes.Market),
					AggregationType: mmtypes.AggregationType_INDEX_PRICE_AGGREGATION,
				},
			},
			err: false,
		},
		{
			name: "single market param",
			params: dydxtypes.QueryAllMarketParamsResponse{
				MarketParams: []dydxtypes.MarketParam{
					{
						Pair:               "BTC-USD", // Taken from dYdX mainnet
						Exponent:           -5,
						MinExchanges:       3,
						ExchangeConfigJson: "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"BTCUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Bybit\",\"ticker\":\"BTCUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"BTC-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"btcusdt\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"XXBTZUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"BTC-USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Mexc\",\"ticker\":\"BTC_USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"BTC-USDT\",\"adjustByMarket\":\"USDT-USD\"}]}",
					},
					{
						Pair:               "ETH-USD", // Taken from dYdX mainnet
						MinExchanges:       3,
						Exponent:           -6,
						ExchangeConfigJson: "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"ETHUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Bybit\",\"ticker\":\"ETHUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"ETH-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"ethusdt\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"XETHZUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"ETH-USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Mexc\",\"ticker\":\"ETH_USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"ETH-USDT\",\"adjustByMarket\":\"USDT-USD\"}]}",
					},
					{
						Pair:               "USDT-USD", // Taken from dYdX mainnet
						MinExchanges:       3,
						Exponent:           -9,
						ExchangeConfigJson: "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"USDCUSDT\",\"invert\":true},{\"exchangeName\":\"Bybit\",\"ticker\":\"USDCUSDT\",\"invert\":true},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"USDT-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"ethusdt\",\"adjustByMarket\":\"ETH-USD\",\"invert\":true},{\"exchangeName\":\"Kraken\",\"ticker\":\"USDTZUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"BTC-USDT\",\"adjustByMarket\":\"BTC-USD\",\"invert\":true},{\"exchangeName\":\"Okx\",\"ticker\":\"USDC-USDT\",\"invert\":true}]}",
					},
				},
			},
			expected: convertedResponse,
			err:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := dydx.ConvertMarketParamsToMarketMap(tc.params)
			if tc.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				for _, market := range tc.expected.MarketMap.Markets {
					require.Equal(t, market, resp.MarketMap.Markets[market.Ticker.String()])
				}

				require.Equal(t, tc.expected.MarketMap.Markets, resp.MarketMap.Markets)
			}
		})
	}
}

func TestCreateCurrencyPairFromMarket(t *testing.T) {
	t.Run("good ticker", func(t *testing.T) {
		pair := "BTC-USD"
		cp, err := dydx.CreateCurrencyPairFromPair(pair)
		require.NoError(t, err)
		require.Equal(t, cp.Base, "BTC")
		require.Equal(t, cp.Quote, "USD")
	})

	t.Run("bad ticker", func(t *testing.T) {
		pair := "BTCUSD"
		_, err := dydx.CreateCurrencyPairFromPair(pair)
		require.Error(t, err)
	})

	t.Run("lower casing still corrects", func(t *testing.T) {
		pair := "btc-usd"
		cp, err := dydx.CreateCurrencyPairFromPair(pair)
		require.NoError(t, err)
		require.Equal(t, cp.Base, "BTC")
		require.Equal(t, cp.Quote, "USD")
	})
}

func TestCreateTickerFromMarket(t *testing.T) {
	testCases := []struct {
		name     string
		market   dydxtypes.MarketParam
		expected mmtypes.Ticker
		err      bool
	}{
		{
			name: "valid market",
			market: dydxtypes.MarketParam{
				Pair:         "BTC-USD",
				MinExchanges: 3,
				Exponent:     -8,
			},
			expected: mmtypes.Ticker{
				CurrencyPair:     slinkytypes.NewCurrencyPair("BTC", "USD"),
				Decimals:         8,
				MinProviderCount: 3,
			},
			err: false,
		},
		{
			name: "invalid market",
			market: dydxtypes.MarketParam{
				Pair:         "BTCUSD",
				MinExchanges: 3,
				Exponent:     -8,
			},
			expected: mmtypes.Ticker{},
			err:      true,
		},
		{
			name: "invalid number of exchanges",
			market: dydxtypes.MarketParam{
				Pair:         "BTC-USD",
				MinExchanges: 0,
				Exponent:     -8,
			},
			expected: mmtypes.Ticker{},
			err:      true,
		},
		{
			name: "invalid exponent",
			market: dydxtypes.MarketParam{
				Pair:         "BTC-USD",
				MinExchanges: 3,
				Exponent:     0,
			},
			expected: mmtypes.Ticker{},
			err:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ticker, err := dydx.CreateTickerFromMarket(tc.market)
			if tc.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, ticker)
			}
		})
	}
}

func TestConvertExchangeConfigJSON(t *testing.T) {
	testCases := []struct {
		name              string
		ticker            mmtypes.Ticker
		config            dydxtypes.ExchangeConfigJson
		expectedPaths     mmtypes.Paths
		expectedProviders mmtypes.Providers
		expectedErr       bool
	}{
		{
			name: "handles duplicate configs",
			ticker: mmtypes.Ticker{
				CurrencyPair:     slinkytypes.NewCurrencyPair("BTC", "USD"),
				Decimals:         8,
				MinProviderCount: 3,
			},
			config: dydxtypes.ExchangeConfigJson{
				Exchanges: []dydxtypes.ExchangeMarketConfigJson{
					{
						ExchangeName: "CoinbasePro",
						Ticker:       "BTC-USD",
					},
					{
						ExchangeName: "CoinbasePro",
						Ticker:       "BTC-USD",
					},
				},
			},
			expectedPaths: mmtypes.Paths{
				Paths: []mmtypes.Path{
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     coinbaseapi.Name,
								CurrencyPair: constants.BITCOIN_USD.CurrencyPair,
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     coinbasews.Name,
								CurrencyPair: constants.BITCOIN_USD.CurrencyPair,
								Invert:       false,
							},
						},
					},
				},
			},
			expectedProviders: mmtypes.Providers{
				Providers: []mmtypes.ProviderConfig{
					{
						Name:           coinbaseapi.Name,
						OffChainTicker: "BTC-USD",
					},
					{
						Name:           coinbasews.Name,
						OffChainTicker: "BTC-USD",
					},
				},
			},
			expectedErr: false,
		},
		{
			name:   "single direct path with no inversion",
			ticker: constants.BITCOIN_USD,
			config: dydxtypes.ExchangeConfigJson{
				Exchanges: []dydxtypes.ExchangeMarketConfigJson{
					{
						ExchangeName: "CoinbasePro",
						Ticker:       "BTC-USD",
					},
				},
			},
			expectedPaths: mmtypes.Paths{
				Paths: []mmtypes.Path{
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     coinbaseapi.Name,
								CurrencyPair: constants.BITCOIN_USD.CurrencyPair,
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     coinbasews.Name,
								CurrencyPair: constants.BITCOIN_USD.CurrencyPair,
								Invert:       false,
							},
						},
					},
				},
			},
			expectedProviders: mmtypes.Providers{
				Providers: []mmtypes.ProviderConfig{
					{
						Name:           coinbaseapi.Name,
						OffChainTicker: "BTC-USD",
					},
					{
						Name:           coinbasews.Name,
						OffChainTicker: "BTC-USD",
					},
				},
			},
			expectedErr: false,
		},
		{
			name:   "single direct path with inversion",
			ticker: constants.USDT_USD,
			config: dydxtypes.ExchangeConfigJson{
				Exchanges: []dydxtypes.ExchangeMarketConfigJson{
					{
						ExchangeName: "Okx",
						Ticker:       "USDC-USDT",
						Invert:       true,
					},
				},
			},
			expectedPaths: mmtypes.Paths{
				Paths: []mmtypes.Path{
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     okx.Name,
								CurrencyPair: constants.USDT_USD.CurrencyPair,
								Invert:       true,
							},
						},
					},
				},
			},
			expectedProviders: mmtypes.Providers{
				Providers: []mmtypes.ProviderConfig{
					{
						Name:           okx.Name,
						OffChainTicker: "USDC-USDT",
					},
				},
			},
			expectedErr: false,
		},
		{
			name:   "single indirect path with an adjustable market",
			ticker: constants.BITCOIN_USD,
			config: dydxtypes.ExchangeConfigJson{
				Exchanges: []dydxtypes.ExchangeMarketConfigJson{
					{
						ExchangeName:   "Okx",
						Ticker:         "BTC-USDT",
						AdjustByMarket: "USDT-USD",
					},
				},
			},
			expectedPaths: mmtypes.Paths{
				Paths: []mmtypes.Path{
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     okx.Name,
								CurrencyPair: constants.BITCOIN_USD.CurrencyPair,
								Invert:       false,
							},
							{
								Provider:     mmtypes.IndexPrice,
								CurrencyPair: constants.USDT_USD.CurrencyPair,
								Invert:       false,
							},
						},
					},
				},
			},
			expectedProviders: mmtypes.Providers{
				Providers: []mmtypes.ProviderConfig{
					{
						Name:           okx.Name,
						OffChainTicker: "BTC-USDT",
					},
				},
			},
			expectedErr: false,
		},
		{
			name:   "single indirect path with an adjustable market and inversion that does not match the ticker",
			ticker: constants.USDT_USD,
			config: dydxtypes.ExchangeConfigJson{
				Exchanges: []dydxtypes.ExchangeMarketConfigJson{
					{
						ExchangeName:   "Kucoin",
						Ticker:         "BTC-USDT",
						AdjustByMarket: "BTC-USD",
						Invert:         true,
					},
				},
			},
			expectedPaths: mmtypes.Paths{
				Paths: []mmtypes.Path{
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     kucoin.Name,
								CurrencyPair: constants.BITCOIN_USD.CurrencyPair,
								Invert:       true,
							},
							{
								Provider:     mmtypes.IndexPrice,
								CurrencyPair: constants.BITCOIN_USD.CurrencyPair,
								Invert:       false,
							},
						},
					},
				},
			},
			expectedProviders: mmtypes.Providers{
				Providers: []mmtypes.ProviderConfig{},
			},
			expectedErr: false,
		},
		{
			name:   "invalid adjust by market",
			ticker: constants.BITCOIN_USD,
			config: dydxtypes.ExchangeConfigJson{
				Exchanges: []dydxtypes.ExchangeMarketConfigJson{
					{
						ExchangeName:   "CoinbasePro",
						Ticker:         "BTC-USDT",
						AdjustByMarket: "USDTUSD",
					},
				},
			},
			expectedPaths:     mmtypes.Paths{},
			expectedProviders: mmtypes.Providers{},
			expectedErr:       true,
		},
		{
			name:   "invalid exchange name",
			ticker: constants.BITCOIN_USD,
			config: dydxtypes.ExchangeConfigJson{
				Exchanges: []dydxtypes.ExchangeMarketConfigJson{
					{
						ExchangeName: "InvalidExchange",
						Ticker:       "BTC-USD",
					},
				},
			},
			expectedPaths:     mmtypes.Paths{},
			expectedProviders: mmtypes.Providers{},
			expectedErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			paths, providers, err := dydx.ConvertExchangeConfigJSON(tc.ticker, tc.config)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}

			require.Equal(t, len(tc.expectedPaths.Paths), len(paths.Paths))
			require.Equal(t, len(tc.expectedProviders.Providers), len(providers.Providers))

			if len(tc.expectedPaths.Paths) > 0 {
				require.Equal(t, tc.expectedPaths, paths)
			}
			if len(tc.expectedProviders.Providers) > 0 {
				require.Equal(t, tc.expectedProviders, providers)
			}
		})
	}
}
