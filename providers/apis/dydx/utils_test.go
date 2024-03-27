package dydx_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	coinbaseapi "github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/apis/dydx"
	dydxtypes "github.com/skip-mev/slinky/providers/apis/dydx/types"
	coinbasews "github.com/skip-mev/slinky/providers/websockets/coinbase"
	"github.com/skip-mev/slinky/providers/websockets/okx"
	mmtypes "github.com/skip-mev/slinky/x/mm2/types"
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
				require.Equal(t, len(tc.expected.MarketMap.Markets), len(resp.MarketMap.Markets))
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
					Name:           coinbaseapi.Name,
					OffChainTicker: "BTC-USD",
				},
				{
					Name:           coinbasews.Name,
					OffChainTicker: "BTC-USD",
				},
			},
			expectedErr: false,
		},
		{
			name: "single direct provider with no inversion",
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
					Name:           coinbaseapi.Name,
					OffChainTicker: "BTC-USD",
				},
				{
					Name:           coinbasews.Name,
					OffChainTicker: "BTC-USD",
				},
			},
			expectedErr: false,
		},
		{
			name: "single direct provider with inversion",
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
			name: "single indirect provider with a normalize by market",
			config: dydxtypes.ExchangeConfigJson{
				Exchanges: []dydxtypes.ExchangeMarketConfigJson{
					{
						ExchangeName:   "Okx",
						Ticker:         "BTC-USDT",
						Invert:         true,
						AdjustByMarket: "USDT-USD",
					},
				},
			},
			expectedProviders: []mmtypes.ProviderConfig{
				{
					Name:           okx.Name,
					OffChainTicker: "BTC-USDT",
					Invert:         true,
					NormalizeByPair: &slinkytypes.CurrencyPair{
						Base:  "USDT",
						Quote: "USD",
					},
				},
			},
			expectedErr: false,
		},
		{
			name:              "No JSON returns empty provider config",
			config:            dydxtypes.ExchangeConfigJson{},
			expectedProviders: []mmtypes.ProviderConfig{},
			expectedErr:       false,
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
