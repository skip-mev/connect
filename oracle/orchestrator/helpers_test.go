package orchestrator_test

import (
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/orchestrator"
	oracletypes "github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/apis/binance"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/apis/dydx"
	"github.com/skip-mev/slinky/providers/base"
	"github.com/skip-mev/slinky/providers/base/api/handlers/mocks"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"
	"github.com/skip-mev/slinky/providers/static"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/websockets/okx"
	mmclienttypes "github.com/skip-mev/slinky/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	logger = zap.NewExample()

	oracleCfg = config.OracleConfig{
		Production: true,
		Metrics: config.MetricsConfig{
			Enabled: false,
		},
		UpdateInterval: 1500 * time.Millisecond,
		MaxPriceAge:    2 * time.Minute,
		Providers: []config.ProviderConfig{
			{ // Price API provider.
				Name: binance.Name,
				API:  binance.DefaultUSAPIConfig,
				Type: oracletypes.ConfigType,
			},
			{ // Price API provider.
				Name: coinbase.Name,
				API:  coinbase.DefaultAPIConfig,
				Type: oracletypes.ConfigType,
			},
			{ // Price WebSocket provider.
				Name:      okx.Name,
				WebSocket: okx.DefaultWebSocketConfig,
				Type:      oracletypes.ConfigType,
			},
		},
	}

	oracleCfgWithMapper = config.OracleConfig{
		Production: true,
		Metrics: config.MetricsConfig{
			Enabled: false,
		},
		UpdateInterval: 1500 * time.Millisecond,
		MaxPriceAge:    2 * time.Minute,
		Providers: []config.ProviderConfig{
			{ // Price API provider.
				Name: binance.Name,
				API:  binance.DefaultUSAPIConfig,
				Type: oracletypes.ConfigType,
			},
			{ // Price API provider.
				Name: coinbase.Name,
				API:  coinbase.DefaultAPIConfig,
				Type: oracletypes.ConfigType,
			},
			{ // Price WebSocket provider.
				Name:      okx.Name,
				WebSocket: okx.DefaultWebSocketConfig,
				Type:      oracletypes.ConfigType,
			},
			// Market map provider.
			mapperCfg,
		},
	}

	mapperCfg = config.ProviderConfig{
		Name: dydx.Name,
		API:  dydx.DefaultAPIConfig,
		Type: mmclienttypes.ConfigType,
	}

	// Coinbase and OKX are supported by the marketmap.
	marketMap = mmtypes.MarketMap{
		Tickers: map[string]mmtypes.Ticker{
			constants.BITCOIN_USD.String():  constants.BITCOIN_USD,
			constants.ETHEREUM_USD.String(): constants.ETHEREUM_USD,
		},
		Providers: map[string]mmtypes.Providers{
			constants.BITCOIN_USD.String(): {
				Providers: []mmtypes.ProviderConfig{
					coinbase.DefaultMarketConfig[constants.BITCOIN_USD],
					okx.DefaultMarketConfig[constants.BITCOIN_USD],
				},
			},
			constants.ETHEREUM_USD.String(): {
				Providers: []mmtypes.ProviderConfig{
					coinbase.DefaultMarketConfig[constants.ETHEREUM_USD],
					okx.DefaultMarketConfig[constants.ETHEREUM_USD],
				},
			},
		},
	}
)

func checkProviderState(
	t *testing.T,
	expectedTickers []mmtypes.Ticker,
	expectedName string,
	expectedType providertypes.ProviderType,
	isRunning bool,
	state orchestrator.ProviderState,
) {
	t.Helper()

	// Ensure that the provider is enabled and supports the expected tickers.
	provider := state.Provider
	require.Equal(t, expectedName, provider.Name())
	require.Equal(t, expectedType, provider.Type())

	ids := provider.GetIDs()
	require.Equal(t, len(expectedTickers), len(ids))
	seenTickers := make(map[mmtypes.Ticker]bool)
	for _, id := range ids {
		seenTickers[id] = true
	}
	for _, ticker := range expectedTickers {
		require.True(t, seenTickers[ticker])
	}

	// Check the market map.
	require.Equal(t, len(expectedTickers), len(provider.GetIDs()))
	for _, ticker := range provider.GetIDs() {
		require.True(t, seenTickers[ticker])
	}

	// Ensure that the provider is running/no-running.
	require.Equal(t, isRunning, provider.IsRunning())
}

func createTestMarketMapProvider(
	t *testing.T,
	responses []mmclienttypes.MarketMapResponse,
	timeout time.Duration,
	ids []mmclienttypes.Chain,
) (*mocks.APIDataHandler[mmclienttypes.Chain, *mmtypes.GetMarketMapResponse], *mmclienttypes.MarketMapProvider) {
	t.Helper()

	// Create a market map api handler.
	handler := mocks.NewAPIDataHandler[mmclienttypes.Chain, *mmtypes.GetMarketMapResponse](
		t,
	)

	queryHandler, err := mmclienttypes.NewMarketMapAPIQueryHandler(
		logger,
		mapperCfg.API,
		static.NewStaticMockClient(),
		handler,
		metrics.NewNopAPIMetrics(),
	)
	require.NoError(t, err)

	provider, err := mmclienttypes.NewMarketMapProvider(
		base.WithName[mmclienttypes.Chain, *mmtypes.GetMarketMapResponse](mapperCfg.Name),
		base.WithLogger[mmclienttypes.Chain, *mmtypes.GetMarketMapResponse](logger),
		base.WithAPIQueryHandler(queryHandler),
		base.WithAPIConfig[mmclienttypes.Chain, *mmtypes.GetMarketMapResponse](mapperCfg.API),
		base.WithMetrics[mmclienttypes.Chain, *mmtypes.GetMarketMapResponse](providermetrics.NewNopProviderMetrics()),
	)
	require.NoError(t, err)

	return handler, provider
}
