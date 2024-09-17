package dydx_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle/constants"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/providers/apis/defi/raydium"
	"github.com/skip-mev/connect/v2/providers/apis/defi/uniswapv3"
	"github.com/skip-mev/connect/v2/providers/apis/dydx"
	dydxtypes "github.com/skip-mev/connect/v2/providers/apis/dydx/types"
	coinbasews "github.com/skip-mev/connect/v2/providers/websockets/coinbase"
	"github.com/skip-mev/connect/v2/providers/websockets/kucoin"
	"github.com/skip-mev/connect/v2/providers/websockets/mexc"
	"github.com/skip-mev/connect/v2/providers/websockets/okx"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

func TestConvertMarketParamsToMarketMap(t *testing.T) {
	testCases := []struct {
		name     string
		params   dydxtypes.QueryAllMarketParamsResponse
		expected mmtypes.MarketMapResponse
		err      bool
	}{
		{
			name:   "empty market params",
			params: dydxtypes.QueryAllMarketParamsResponse{},
			expected: mmtypes.MarketMapResponse{
				MarketMap: mmtypes.MarketMap{
					Markets: make(map[string]mmtypes.Market),
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
				require.Equal(t, tc.expected, resp)
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
				CurrencyPair:     connecttypes.NewCurrencyPair("BTC", "USD"),
				Decimals:         8,
				MinProviderCount: 3,
				Enabled:          true,
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
		config            dydxtypes.ExchangeConfigJson
		expectedProviders []mmtypes.ProviderConfig
		expectedErr       bool
	}{
		{
			name: "handles duplicate configs",
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
			expectedProviders: []mmtypes.ProviderConfig{
				{
					Name:           coinbasews.Name,
					OffChainTicker: "BTC-USD",
				},
			},
			expectedErr: false,
		},
		{
			name: "single direct path with no inversion",
			config: dydxtypes.ExchangeConfigJson{
				Exchanges: []dydxtypes.ExchangeMarketConfigJson{
					{
						ExchangeName: "CoinbasePro",
						Ticker:       "BTC-USD",
					},
				},
			},
			expectedProviders: []mmtypes.ProviderConfig{
				{
					Name:           coinbasews.Name,
					OffChainTicker: "BTC-USD",
				},
			},
			expectedErr: false,
		},
		{
			name: "single direct path with inversion",
			config: dydxtypes.ExchangeConfigJson{
				Exchanges: []dydxtypes.ExchangeMarketConfigJson{
					{
						ExchangeName: "Okx",
						Ticker:       "USDC-USDT",
						Invert:       true,
					},
				},
			},
			expectedProviders: []mmtypes.ProviderConfig{
				{
					Name:           okx.Name,
					OffChainTicker: "USDC-USDT",
					Invert:         true,
				},
			},
			expectedErr: false,
		},
		{
			name: "single indirect path with an adjustable market",
			config: dydxtypes.ExchangeConfigJson{
				Exchanges: []dydxtypes.ExchangeMarketConfigJson{
					{
						ExchangeName:   "Okx",
						Ticker:         "BTC-USDT",
						AdjustByMarket: "USDT-USD",
					},
				},
			},
			expectedProviders: []mmtypes.ProviderConfig{
				{
					Name:           okx.Name,
					OffChainTicker: "BTC-USDT",
					NormalizeByPair: &connecttypes.CurrencyPair{
						Base:  "USDT",
						Quote: "USD",
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "single indirect path with an adjustable market and inversion that does not match the ticker",
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
			expectedProviders: []mmtypes.ProviderConfig{
				{
					Name:           kucoin.Name,
					OffChainTicker: "BTC-USDT",
					NormalizeByPair: &connecttypes.CurrencyPair{
						Base:  "BTC",
						Quote: "USD",
					},
					Invert: true,
				},
			},
			expectedErr: false,
		},
		{
			name: "invalid adjust by market",
			config: dydxtypes.ExchangeConfigJson{
				Exchanges: []dydxtypes.ExchangeMarketConfigJson{
					{
						ExchangeName:   "CoinbasePro",
						Ticker:         "BTC-USDT",
						AdjustByMarket: "USDTUSD",
					},
				},
			},
			expectedProviders: []mmtypes.ProviderConfig{},
			expectedErr:       true,
		},
		{
			name: "invalid exchange name - should ignore",
			config: dydxtypes.ExchangeConfigJson{
				Exchanges: []dydxtypes.ExchangeMarketConfigJson{
					{
						ExchangeName: "InvalidExchange",
						Ticker:       "BTC-USD",
					},
					{
						ExchangeName: "CoinbasePro",
						Ticker:       "BTC-USD",
					},
				},
			},
			expectedProviders: []mmtypes.ProviderConfig{
				{
					Name:           coinbasews.Name,
					OffChainTicker: "BTC-USD",
				},
			},
			expectedErr: false,
		},
		{
			name: "exchange that includes a denom that needs to be converted",
			config: dydxtypes.ExchangeConfigJson{
				Exchanges: []dydxtypes.ExchangeMarketConfigJson{
					{
						ExchangeName:   "Mexc",
						Ticker:         "ETH_USDT",
						AdjustByMarket: "USDT-USD",
					},
				},
			},
			expectedProviders: []mmtypes.ProviderConfig{
				{
					Name:           mexc.Name,
					OffChainTicker: "ETHUSDT",
					NormalizeByPair: &connecttypes.CurrencyPair{
						Base:  "USDT",
						Quote: "USD",
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "raydium exchange config",
			config: dydxtypes.ExchangeConfigJson{
				Exchanges: []dydxtypes.ExchangeMarketConfigJson{
					{
						ExchangeName: "Raydium",
						Ticker:       "SMOLE-SOL-VDZ9kwvKRbqhNdsoRZyLVzAAQMbGY9akHbtM6YugViS-8-HiLcngHP5y1Jno53tuuNeFHKWhyyZp3XuxtKPszD6rG2-9-FeKBjZ5rBvHPyppHf11qjYxwaQuiympppCTQ5pC6om3F-5EgCcjkuE42YyTZY4QG8qTioUwNh6agTvJuNRyEqcqV1",
					},
				},
			},
			expectedProviders: []mmtypes.ProviderConfig{
				{
					Name:           raydium.Name,
					OffChainTicker: "SMOLE/SOL",
					Metadata_JSON:  "{\"base_token_vault\":{\"token_vault_address\":\"VDZ9kwvKRbqhNdsoRZyLVzAAQMbGY9akHbtM6YugViS\",\"token_decimals\":8},\"quote_token_vault\":{\"token_vault_address\":\"HiLcngHP5y1Jno53tuuNeFHKWhyyZp3XuxtKPszD6rG2\",\"token_decimals\":9},\"amm_info_address\":\"5EgCcjkuE42YyTZY4QG8qTioUwNh6agTvJuNRyEqcqV1\",\"open_orders_address\":\"FeKBjZ5rBvHPyppHf11qjYxwaQuiympppCTQ5pC6om3F\"}",
				},
			},
			expectedErr: false,
		},
		{
			name: "uniswapv3-ethereum exchange config",
			config: dydxtypes.ExchangeConfigJson{
				Exchanges: []dydxtypes.ExchangeMarketConfigJson{
					{
						ExchangeName:   "UniswapV3-Ethereum",
						Ticker:         "0x0c30062368eEfB96bF3AdE1218E685306b8E89Fa-8-18",
						AdjustByMarket: "ETH-USD",
						Invert:         false,
					},
				},
			},
			expectedProviders: []mmtypes.ProviderConfig{
				{
					Name:           uniswapv3.ProviderNames[constants.ETHEREUM],
					OffChainTicker: "0x0c30062368eEfB96bF3AdE1218E685306b8E89Fa-8-18",
					Metadata_JSON:  "{\"address\":\"0x0c30062368eEfB96bF3AdE1218E685306b8E89Fa\",\"base_decimals\":8,\"quote_decimals\":18,\"invert\":false}",
					NormalizeByPair: &connecttypes.CurrencyPair{
						Base:  "ETH",
						Quote: "USD",
					},
				},
			},
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			providers, err := dydx.ConvertExchangeConfigJSON(tc.config)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}

			require.Equal(t, len(tc.expectedProviders), len(providers))

			if len(tc.expectedProviders) > 0 {
				require.Equal(t, tc.expectedProviders, providers)
			}
		})
	}
}

func TestExtractMetadata(t *testing.T) {
	testcases := []struct {
		name             string
		providerName     string
		cfg              dydxtypes.ExchangeMarketConfigJson
		expectedMetadata string
		expectedErr      bool
	}{
		{
			name:             "non-raydium provider",
			providerName:     kucoin.Name,
			cfg:              dydxtypes.ExchangeMarketConfigJson{Ticker: "BTC-USDT"},
			expectedMetadata: "",
			expectedErr:      false,
		},
		{
			name:             "raydium provider w/o additional metadata in ticker",
			providerName:     raydium.Name,
			cfg:              dydxtypes.ExchangeMarketConfigJson{Ticker: "BTC-USDT"},
			expectedMetadata: "",
			expectedErr:      true,
		},
		{
			name:             "raydium provider w/ non-solana base token",
			providerName:     raydium.Name,
			cfg:              dydxtypes.ExchangeMarketConfigJson{Ticker: "SMOLE-SOL-abc-6-def-7"},
			expectedMetadata: "",
			expectedErr:      true,
		},
		{
			name:             "raydium provider w/ non-solana quote token",
			providerName:     raydium.Name,
			cfg:              dydxtypes.ExchangeMarketConfigJson{Ticker: "SMOLE-SOL-VDZ9kwvKRbqhNdsoRZyLVzAAQMbGY9akHbtM6YugViS-6-def-7"},
			expectedMetadata: "",
			expectedErr:      true,
		},
		{
			name:         "raydium provider w/ incorrect base decimals",
			providerName: raydium.Name,
			cfg:          dydxtypes.ExchangeMarketConfigJson{Ticker: "SMOLE-SOL-VDZ9kwvKRbqhNdsoRZyLVzAAQMbGY9akHbtM6YugViS-a-HiLcngHP5y1Jno53tuuNeFHKWhyyZp3XuxtKPszD6rG2-7"},
			expectedErr:  true,
		},
		{
			name:         "raydium provider w/ incorrect base decimals",
			providerName: raydium.Name,
			cfg:          dydxtypes.ExchangeMarketConfigJson{Ticker: "SMOLE-SOL-VDZ9kwvKRbqhNdsoRZyLVzAAQMbGY9akHbtM6YugViS-8-HiLcngHP5y1Jno53tuuNeFHKWhyyZp3XuxtKPszD6rG2-a"},
			expectedErr:  true,
		},
		{
			name:         "raydium provider w/ incorrect open-orders account",
			providerName: raydium.Name,
			cfg:          dydxtypes.ExchangeMarketConfigJson{Ticker: "SMOLE-SOL-VDZ9kwvKRbqhNdsoRZyLVzAAQMbGY9akHbtM6YugViS-8-HiLcngHP5y1Jno53tuuNeFHKWhyyZp3XuxtKPszD6rG2-9-a-5EgCcjkuE42YyTZY4QG8qTioUwNh6agTvJuNRyEqcqV1"},
			expectedErr:  true,
		},
		{
			name:         "raydium provider w/ incorrect ammId account",
			providerName: raydium.Name,
			cfg:          dydxtypes.ExchangeMarketConfigJson{Ticker: "SMOLE-SOL-VDZ9kwvKRbqhNdsoRZyLVzAAQMbGY9akHbtM6YugViS-8-HiLcngHP5y1Jno53tuuNeFHKWhyyZp3XuxtKPszD6rG2-9-FeKBjZ5rBvHPyppHf11qjYxwaQuiympppCTQ5pC6om3F-a"},
			expectedErr:  true,
		},
		{
			name:             "raydium provider w/ correct metadata",
			providerName:     raydium.Name,
			cfg:              dydxtypes.ExchangeMarketConfigJson{Ticker: "SMOLE-SOL-VDZ9kwvKRbqhNdsoRZyLVzAAQMbGY9akHbtM6YugViS-8-HiLcngHP5y1Jno53tuuNeFHKWhyyZp3XuxtKPszD6rG2-9-FeKBjZ5rBvHPyppHf11qjYxwaQuiympppCTQ5pC6om3F-5EgCcjkuE42YyTZY4QG8qTioUwNh6agTvJuNRyEqcqV1"},
			expectedMetadata: "{\"base_token_vault\":{\"token_vault_address\":\"VDZ9kwvKRbqhNdsoRZyLVzAAQMbGY9akHbtM6YugViS\",\"token_decimals\":8},\"quote_token_vault\":{\"token_vault_address\":\"HiLcngHP5y1Jno53tuuNeFHKWhyyZp3XuxtKPszD6rG2\",\"token_decimals\":9},\"amm_info_address\":\"5EgCcjkuE42YyTZY4QG8qTioUwNh6agTvJuNRyEqcqV1\",\"open_orders_address\":\"FeKBjZ5rBvHPyppHf11qjYxwaQuiympppCTQ5pC6om3F\"}",
			expectedErr:      false,
		},
		{
			name:         "invalid exchange",
			providerName: "foobar",
			expectedErr:  false,
		},
		{
			name:         "uniswapv3-ethereum invalid field number",
			providerName: uniswapv3.ProviderNames[constants.ETHEREUM],
			cfg:          dydxtypes.ExchangeMarketConfigJson{Ticker: "0xabc123-abc"},
			expectedErr:  true,
		},
		{
			name:         "uniswapv3-ethereum invalid base decimals",
			providerName: uniswapv3.ProviderNames[constants.ETHEREUM],
			cfg:          dydxtypes.ExchangeMarketConfigJson{Ticker: "0xabc123-abc-12"},
			expectedErr:  true,
		},
		{
			name:         "uniswapv3-ethereum invalid quote decimals",
			providerName: uniswapv3.ProviderNames[constants.ETHEREUM],
			cfg:          dydxtypes.ExchangeMarketConfigJson{Ticker: "0xabc123-8-abc"},
			expectedErr:  true,
		},
		{
			name:         "uniswapv3-ethereum invalid pool address",
			providerName: uniswapv3.ProviderNames[constants.ETHEREUM],
			cfg:          dydxtypes.ExchangeMarketConfigJson{Ticker: "zzzzzz-8-18"},
			expectedErr:  true,
		},
		{
			name:         "uniswapv3-ethereum valid config",
			providerName: uniswapv3.ProviderNames[constants.ETHEREUM],
			cfg: dydxtypes.ExchangeMarketConfigJson{
				Ticker: "0xE7F6720C1F546217081667A5ab7fEbB688036856-8-18",
				Invert: true,
			},
			expectedMetadata: "{\"address\":\"0xE7F6720C1F546217081667A5ab7fEbB688036856\",\"base_decimals\":8,\"quote_decimals\":18,\"invert\":true}",
			expectedErr:      false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			metadata, err := dydx.ExtractMetadata(tc.providerName, tc.cfg)
			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedMetadata, metadata)
			}
		})
	}
}
