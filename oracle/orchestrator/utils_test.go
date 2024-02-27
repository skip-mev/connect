package orchestrator_test

import (
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/orchestrator"
	"github.com/skip-mev/slinky/providers/apis/binance"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/websockets/okx"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	"github.com/stretchr/testify/require"
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
			{
				Name: binance.Name,
				API:  binance.DefaultUSAPIConfig,
			},
			{
				Name: coinbase.Name,
				API:  coinbase.DefaultAPIConfig,
			},
			{
				Name:      okx.Name,
				WebSocket: okx.DefaultWebSocketConfig,
			},
		},
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
	enabled bool,
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
	market := state.Market
	require.Equal(t, len(expectedTickers), len(market.GetTickers()))
	for _, ticker := range market.GetTickers() {
		require.True(t, seenTickers[ticker])
	}

	// Ensure that the provider is enabled/disabled.
	require.Equal(t, enabled, state.Enabled)

	// Ensure that the provider is running/no-running.
	require.Equal(t, isRunning, provider.IsRunning())

	// By default, all providers should be configured with a config updater.
	require.NotNil(t, provider.GetConfigUpdater())
}
