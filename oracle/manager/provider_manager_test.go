package manager_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/manager"
	"github.com/skip-mev/slinky/providers/apis/binance"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	oraclefactory "github.com/skip-mev/slinky/providers/factories/oracle"
	"github.com/skip-mev/slinky/providers/websockets/okx"
)

func TestInit(t *testing.T) {
	t.Run("creates all providers without a marketmap", func(t *testing.T) {
		manager, err := manager.NewProviderManager(
			oracleCfg,
			manager.WithLogger(logger),
			manager.WithAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			manager.WithWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)

		err = manager.Init()
		require.NoError(t, err)

		state := manager.GetProviderState()
		require.Equal(t, len(state), len(oracleCfg.Providers))

		for name, pState := range state {
			require.Equal(t, name, pState.Provider.Name())
			require.False(t, pState.Enabled)
			require.Equal(t, 0, len(pState.Market.GetTickers()))
		}
	})

	t.Run("creates some providers with a marketmap", func(t *testing.T) {
		manager, err := manager.NewProviderManager(
			oracleCfg,
			manager.WithLogger(logger),
			manager.WithMarketMap(marketMap),
			manager.WithAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			manager.WithWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)

		err = manager.Init()
		require.NoError(t, err)

		state := manager.GetProviderState()
		require.Equal(t, len(state), len(oracleCfg.Providers))

		coinbaseState, ok := state[coinbase.Name]
		require.True(t, ok)
		require.Equal(t, coinbaseState.Provider.Name(), coinbase.Name)
		require.True(t, coinbaseState.Enabled)
		require.Equal(t, 2, len(coinbaseState.Market.GetTickers()))

		okxState, ok := state[okx.Name]
		require.True(t, ok)
		require.Equal(t, okxState.Provider.Name(), okx.Name)
		require.True(t, okxState.Enabled)
		require.Equal(t, 2, len(okxState.Market.GetTickers()))

		// Ensure that the provider that is not supported by the marketmap is not enabled.
		binanceState, ok := state[binance.Name]
		require.True(t, ok)
		require.Equal(t, binanceState.Provider.Name(), binance.Name)
		require.False(t, binanceState.Enabled)
		require.Equal(t, 0, len(binanceState.Market.GetTickers()))
	})

	t.Run("errors when the API query handler factory is not set", func(t *testing.T) {
		manager, err := manager.NewProviderManager(
			oracleCfg,
			manager.WithLogger(logger),
			manager.WithMarketMap(marketMap),
			manager.WithWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)

		err = manager.Init()
		require.Error(t, err)
	})

	t.Run("errors when the WebSocket query handler factory is not set", func(t *testing.T) {
		manager, err := manager.NewProviderManager(
			oracleCfg,
			manager.WithLogger(logger),
			manager.WithMarketMap(marketMap),
			manager.WithAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
		)
		require.NoError(t, err)

		err = manager.Init()
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

		manager, err := manager.NewProviderManager(
			cfg,
			manager.WithLogger(logger),
			manager.WithMarketMap(marketMap),
			manager.WithAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			manager.WithWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)

		err = manager.Init()
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

		manager, err := manager.NewProviderManager(
			cfg,
			manager.WithLogger(logger),
			manager.WithMarketMap(marketMap),
			manager.WithAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			manager.WithWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)

		err = manager.Init()
		require.Error(t, err)
	})
}
