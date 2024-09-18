package dydx_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/skip-mev/connect/v2/providers/apis/coinmarketcap"
	dydxtypes "github.com/skip-mev/connect/v2/providers/apis/dydx/types"
	"github.com/skip-mev/connect/v2/providers/base/testutils"
	"github.com/skip-mev/connect/v2/providers/websockets/binance"
	"github.com/skip-mev/connect/v2/providers/websockets/coinbase"
	"github.com/skip-mev/connect/v2/providers/websockets/gate"
	"github.com/skip-mev/connect/v2/providers/websockets/kucoin"
	"github.com/skip-mev/connect/v2/providers/websockets/mexc"
	"github.com/skip-mev/connect/v2/providers/websockets/okx"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/providers/apis/dydx"
	"github.com/skip-mev/connect/v2/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

func TestNewResearchAPIHandler(t *testing.T) {
	t.Run("fail if the name is incorrect", func(t *testing.T) {
		_, err := dydx.NewResearchAPIHandler(zap.NewNop(), config.APIConfig{
			Name: "incorrect",
		})
		require.Error(t, err)
	})

	t.Run("fail if the api is not enabled", func(t *testing.T) {
		_, err := dydx.NewResearchAPIHandler(zap.NewNop(), config.APIConfig{
			Name:    dydx.ResearchAPIHandlerName,
			Enabled: false,
		})
		require.Error(t, err)
	})

	t.Run("test failure of api-config validation", func(t *testing.T) {
		cfg := dydx.DefaultResearchAPIConfig
		cfg.Endpoints = []config.Endpoint{
			{
				URL: "",
			},
		}

		_, err := dydx.NewResearchAPIHandler(zap.NewNop(), cfg)
		require.Error(t, err)
	})

	t.Run("test failure if no endpoint is given", func(t *testing.T) {
		cfg := dydx.DefaultResearchAPIConfig
		cfg.Endpoints = nil

		_, err := dydx.NewResearchAPIHandler(zap.NewNop(), cfg)
		require.Error(t, err)
	})

	t.Run("test success", func(t *testing.T) {
		_, err := dydx.NewResearchAPIHandler(zap.NewNop(), dydx.DefaultResearchAPIConfig)
		require.NoError(t, err)
	})
}

// TestCreateURL tests that:
//   - If no chain in the given chains are dydx - fail
//   - If one chain in the given chains is dydx - return the first endpoint configured
func TestCreateURLResearchHandler(t *testing.T) {
	ah, err := dydx.NewResearchAPIHandler(
		zap.NewNop(),
		dydx.DefaultResearchAPIConfig,
	)
	require.NoError(t, err)

	t.Run("non-dydx chains", func(t *testing.T) {
		chains := []types.Chain{
			{
				ChainID: "osmosis",
			},
		}

		url, err := ah.CreateURL(chains)
		require.Error(t, err)
		require.Empty(t, url)
	})
	t.Run("multiple chains w/ a dydx chain", func(t *testing.T) {
		chains := []types.Chain{
			{
				ChainID: "osmosis",
			},
			{
				ChainID: dydx.ChainID,
			},
		}

		url, err := ah.CreateURL(chains)
		require.NoError(t, err)
		require.Equal(t, dydx.DefaultResearchAPIConfig.Endpoints[1].URL, url)
	})
}

func TestParseResponseResearchAPI(t *testing.T) {
	ah, err := dydx.NewResearchAPIHandler(
		zap.NewNop(),
		dydx.DefaultResearchAPIConfig,
	)
	require.NoError(t, err)

	t.Run("fail if none of the chains given are dydx", func(t *testing.T) {
		chains := []types.Chain{
			{
				ChainID: "osmosis",
			},
		}

		resp := ah.ParseResponse(chains, &http.Response{})
		// expect a failure response for each chain
		require.Len(t, resp.UnResolved, 1)
		require.Len(t, resp.Resolved, 0)

		require.Error(t, resp.UnResolved[chains[0]])
	})

	t.Run("failing to parse ResearchJSON response", func(t *testing.T) {
		chains := []types.Chain{
			{
				ChainID: dydx.ChainID,
			},
		}

		resp := ah.ParseResponse(chains, &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("")),
		})

		require.Len(t, resp.UnResolved, 1)
		require.Len(t, resp.Resolved, 0)

		require.Error(t, resp.UnResolved[chains[0]])
	})

	t.Run("failing to convert ResearchJSON response into QueryAllMarketsParams", func(t *testing.T) {
		chains := []types.Chain{
			{
				ChainID: dydx.ChainID,
			},
		}

		resp := ah.ParseResponse(chains, &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(bytes.NewBufferString(`{
				"1INCH": {
				}
			}`)),
		})

		require.Len(t, resp.UnResolved, 1)
		require.Len(t, resp.Resolved, 0)

		require.Error(t, resp.UnResolved[chains[0]])
	})

	t.Run("success", func(t *testing.T) {
		chains := []types.Chain{
			{
				ChainID: dydx.ChainID,
			},
		}

		researchJSON := dydxtypes.ResearchJSON{
			"1INCH": {
				ResearchJSONMarketParam: dydxtypes.ResearchJSONMarketParam{
					ID:                0,
					Pair:              "1INCH-USD",
					Exponent:          -10.0,
					MinPriceChangePpm: 4000,
					MinExchanges:      3,
					ExchangeConfigJSON: []dydxtypes.ExchangeMarketConfigJson{
						{
							ExchangeName: "Binance",
							Ticker:       "1INCHUSDT",
						},
						{
							ExchangeName: "CoinbasePro",
							Ticker:       "1INCH-USD",
						},
						{
							ExchangeName: "Gate",
							Ticker:       "1INCH_USDT",
						},
						{
							ExchangeName: "Kucoin",
							Ticker:       "1INCH-USDT",
						},
						{
							ExchangeName: "Mexc",
							Ticker:       "1INCH_USDT",
						},
						{
							ExchangeName: "Okx",
							Ticker:       "1INCH-USDT",
						},
					},
				},
			},
		}
		bz, err := json.Marshal(researchJSON)
		require.NoError(t, err)

		resp := ah.ParseResponse(chains, testutils.CreateResponseFromJSON(string(bz)))

		require.Len(t, resp.UnResolved, 0)
		require.Len(t, resp.Resolved, 1)

		mm := resp.Resolved[chains[0]].Value.MarketMap
		require.Len(t, mm.Markets, 1)

		// index by the pair
		market, ok := mm.Markets["1INCH/USD"]
		require.True(t, ok)

		// check the ticker
		expectedTicker := mmtypes.Ticker{
			CurrencyPair:     connecttypes.NewCurrencyPair("1INCH", "USD"),
			Decimals:         10,
			MinProviderCount: 3,
			Enabled:          true,
		}
		require.Equal(t, expectedTicker, market.Ticker)

		// check each provider
		expectedProviders := map[string]mmtypes.ProviderConfig{
			binance.Name: {
				Name:           binance.Name,
				OffChainTicker: "1INCHUSDT",
			},
			coinbase.Name: {
				Name:           coinbase.Name,
				OffChainTicker: "1INCH-USD",
			},
			gate.Name: {
				Name:           gate.Name,
				OffChainTicker: "1INCH_USDT",
			},
			kucoin.Name: {
				Name:           kucoin.Name,
				OffChainTicker: "1INCH-USDT",
			},
			mexc.Name: {
				Name:           mexc.Name,
				OffChainTicker: "1INCHUSDT",
			},
			okx.Name: {
				Name:           okx.Name,
				OffChainTicker: "1INCH-USDT",
			},
		}

		for _, provider := range market.ProviderConfigs {
			expectedProvider, ok := expectedProviders[provider.Name]
			require.True(t, ok)
			require.Equal(t, expectedProvider, provider)
		}
	})
}

func TestParseResponseResearchCMCAPI(t *testing.T) {
	ah, err := dydx.NewResearchAPIHandler(
		zap.NewNop(),
		dydx.DefaultResearchCMCAPIConfig,
	)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		chains := []types.Chain{
			{
				ChainID: dydx.ChainID,
			},
		}

		researchJSON := dydxtypes.ResearchJSON{
			"1INCH": {
				ResearchJSONMarketParam: dydxtypes.ResearchJSONMarketParam{
					ID:                0,
					Pair:              "1INCH-USD",
					Exponent:          -10.0,
					MinPriceChangePpm: 4000,
					MinExchanges:      3,
					ExchangeConfigJSON: []dydxtypes.ExchangeMarketConfigJson{
						{
							ExchangeName: "Binance",
							Ticker:       "1INCHUSDT",
						},
						{
							ExchangeName: "CoinbasePro",
							Ticker:       "1INCH-USD",
						},
						{
							ExchangeName: "Gate",
							Ticker:       "1INCH_USDT",
						},
						{
							ExchangeName: "Kucoin",
							Ticker:       "1INCH-USDT",
						},
						{
							ExchangeName: "Mexc",
							Ticker:       "1INCH_USDT",
						},
						{
							ExchangeName: "Okx",
							Ticker:       "1INCH-USDT",
						},
					},
				},
				MetaData: dydxtypes.MetaData{
					CMCID: 1,
				},
			},
		}

		bz, err := json.Marshal(researchJSON)
		require.NoError(t, err)

		resp := ah.ParseResponse(chains, testutils.CreateResponseFromJSON(string(bz)))

		require.Len(t, resp.UnResolved, 0)
		require.Len(t, resp.Resolved, 1)

		mm := resp.Resolved[chains[0]].Value.MarketMap
		require.Len(t, mm.Markets, 1)

		// index by the pair
		market, ok := mm.Markets["1INCH/USD"]
		require.True(t, ok)

		// check the ticker
		expectedTicker := mmtypes.Ticker{
			CurrencyPair:     connecttypes.NewCurrencyPair("1INCH", "USD"),
			Decimals:         10,
			MinProviderCount: 1,
			Enabled:          true,
		}
		require.Equal(t, expectedTicker, market.Ticker)

		// check each provider
		expectedProviders := map[string]mmtypes.ProviderConfig{
			coinmarketcap.Name: {
				Name:           coinmarketcap.Name,
				OffChainTicker: "1",
			},
		}

		for _, provider := range market.ProviderConfigs {
			expectedProvider, ok := expectedProviders[provider.Name]
			require.True(t, ok)
			require.Equal(t, expectedProvider, provider)
		}
	})
}
