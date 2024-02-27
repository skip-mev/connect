package orchestrator_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/orchestrator"
	"github.com/skip-mev/slinky/providers/apis/binance"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	oraclefactory "github.com/skip-mev/slinky/providers/factories/oracle"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/websockets/okx"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

func TestInit(t *testing.T) {
	t.Run("creates all providers without a marketmap", func(t *testing.T) {
		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfg,
			orchestrator.WithLogger(logger),
			orchestrator.WithAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)

		err = o.Init()
		require.NoError(t, err)

		state := o.GetProviderState()
		require.Equal(t, len(state), len(oracleCfg.Providers))

		coinbaseState, ok := state[coinbase.Name]
		require.True(t, ok)
		checkProviderState(t, nil, coinbase.Name, false, providertypes.API, false, coinbaseState)

		okxState, ok := state[okx.Name]
		require.True(t, ok)
		checkProviderState(t, nil, okx.Name, false, providertypes.WebSockets, false, okxState)

		binanceState, ok := state[binance.Name]
		require.True(t, ok)
		checkProviderState(t, nil, binance.Name, false, providertypes.API, false, binanceState)
	})

	t.Run("creates some providers with a marketmap", func(t *testing.T) {
		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfg,
			orchestrator.WithLogger(logger),
			orchestrator.WithMarketMap(marketMap),
			orchestrator.WithAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)

		err = o.Init()
		require.NoError(t, err)

		state := o.GetProviderState()
		require.Equal(t, len(state), len(oracleCfg.Providers))

		expectedTickers := []mmtypes.Ticker{constants.BITCOIN_USD, constants.ETHEREUM_USD}

		coinbaseState, ok := state[coinbase.Name]
		require.True(t, ok)
		checkProviderState(t, expectedTickers, coinbase.Name, true, providertypes.API, false, coinbaseState)

		okxState, ok := state[okx.Name]
		require.True(t, ok)
		checkProviderState(t, expectedTickers, okx.Name, true, providertypes.WebSockets, false, okxState)

		// Ensure that the provider that is not supported by the marketmap is not enabled.
		binanceState, ok := state[binance.Name]
		require.True(t, ok)
		checkProviderState(t, nil, binance.Name, false, providertypes.API, false, binanceState)
	})

	t.Run("errors when the API query handler factory is not set", func(t *testing.T) {
		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfg,
			orchestrator.WithLogger(logger),
			orchestrator.WithMarketMap(marketMap),
			orchestrator.WithWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)

		err = o.Init()
		require.Error(t, err)
	})

	t.Run("errors when the WebSocket query handler factory is not set", func(t *testing.T) {
		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfg,
			orchestrator.WithLogger(logger),
			orchestrator.WithMarketMap(marketMap),
			orchestrator.WithAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
		)
		require.NoError(t, err)

		err = o.Init()
		require.Error(t, err)
	})

	t.Run("errors when a provider is not supported by the api query handler factory", func(t *testing.T) {
		cfg := oracleCfg
		cfg.Providers = append(cfg.Providers, config.ProviderConfig{
			Name: "unsupported",
			API: config.APIConfig{
				Enabled:    true,
				Timeout:    5,
				Interval:   5,
				MaxQueries: 5,
				URL:        "https://example.com",
				Name:       "unsupported",
			},
		})

		o, err := orchestrator.NewProviderOrchestrator(
			cfg,
			orchestrator.WithLogger(logger),
			orchestrator.WithMarketMap(marketMap),
			orchestrator.WithAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)

		err = o.Init()
		require.Error(t, err)
	})

	t.Run("errors when a provider is not supported by the web socket query handler factory", func(t *testing.T) {
		cfg := oracleCfg

		okx := okx.DefaultWebSocketConfig
		okx.Name = "unsupported"
		cfg.Providers = append(cfg.Providers, config.ProviderConfig{
			Name:      "unsupported",
			WebSocket: okx,
		})

		o, err := orchestrator.NewProviderOrchestrator(
			cfg,
			orchestrator.WithLogger(logger),
			orchestrator.WithMarketMap(marketMap),
			orchestrator.WithAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)

		err = o.Init()
		require.Error(t, err)
	})
}
