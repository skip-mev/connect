package orchestrator_test

import (
	"context"
	"testing"
	"time"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/orchestrator"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	oraclefactory "github.com/skip-mev/slinky/providers/factories/oracle"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/websockets/okx"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	"github.com/stretchr/testify/require"
)

func TestUpdateProviderState(t *testing.T) {
	expectedTickers := []mmtypes.Ticker{constants.BITCOIN_USD, constants.ETHEREUM_USD}

	t.Run("can update a single api provider state with no configuration and non-running", func(t *testing.T) {
		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfg,
			orchestrator.WithLogger(logger),
			orchestrator.WithAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		require.NoError(t, o.Init())

		providerMarketMap, err := types.ProviderMarketMapFromMarketMap(coinbase.Name, marketMap)
		require.NoError(t, err)

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		providerState, ok := providers[coinbase.Name]
		require.True(t, ok)

		// Check the state before any modifications are done.
		checkProviderState(t, nil, coinbase.Name, false, providertypes.API, false, providerState)

		updatedState, err := o.UpdateProviderState(providerMarketMap, providerState)
		require.NoError(t, err)

		// Check the state after the update.
		checkProviderState(t, expectedTickers, coinbase.Name, true, providertypes.API, false, updatedState)

		updatedState.Provider.Stop()
		time.Sleep(2000 * time.Millisecond)
		require.Eventually(
			t,
			func() bool {
				return !updatedState.Provider.IsRunning()
			},
			5*time.Second,
			500*time.Millisecond,
		)
	})

	t.Run("can update a single api provider state with no configuration and running", func(t *testing.T) {
		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfg,
			orchestrator.WithLogger(logger),
			orchestrator.WithAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		require.NoError(t, o.Init())

		providerMarketMap, err := types.ProviderMarketMapFromMarketMap(coinbase.Name, marketMap)
		require.NoError(t, err)

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		providerState, ok := providers[coinbase.Name]
		require.True(t, ok)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Start the provider.
		go func() {
			providerState.Provider.Start(ctx)
		}()

		time.Sleep(500 * time.Millisecond)

		// Check the state before any modifications are done.
		checkProviderState(t, nil, coinbase.Name, false, providertypes.API, true, providerState)

		updatedState, err := o.UpdateProviderState(providerMarketMap, providerState)
		require.NoError(t, err)

		time.Sleep(1000 * time.Millisecond)

		// Check the state after the update.
		checkProviderState(t, expectedTickers, coinbase.Name, true, providertypes.API, true, updatedState)

		updatedState.Provider.Stop()
		time.Sleep(2000 * time.Millisecond)
		require.Eventually(
			t,
			func() bool {
				return !updatedState.Provider.IsRunning()
			},
			5*time.Second,
			500*time.Millisecond,
		)
	})

	t.Run("can update a single api provider state removing all tickers on a non-running provider", func(t *testing.T) {
		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfg,
			orchestrator.WithLogger(logger),
			orchestrator.WithAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
			orchestrator.WithMarketMap(marketMap),
		)
		require.NoError(t, err)
		require.NoError(t, o.Init())

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		providerState, ok := providers[coinbase.Name]
		require.True(t, ok)

		// Check the state before any modifications are done.
		checkProviderState(t, expectedTickers, coinbase.Name, true, providertypes.API, false, providerState)

		pMarketMap := types.ProviderMarketMap{
			Name: coinbase.Name,
		}
		updatedState, err := o.UpdateProviderState(pMarketMap, providerState)
		require.NoError(t, err)

		time.Sleep(1000 * time.Millisecond)

		// Check the state after the update.
		checkProviderState(t, nil, coinbase.Name, false, providertypes.API, false, updatedState)

		updatedState.Provider.Stop()
		time.Sleep(2000 * time.Millisecond)
		require.Eventually(
			t,
			func() bool {
				return !updatedState.Provider.IsRunning()
			},
			5*time.Second,
			500*time.Millisecond,
		)
	})

	t.Run("can update a single api provider state removing all tickers on a running provider", func(t *testing.T) {
		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfg,
			orchestrator.WithLogger(logger),
			orchestrator.WithAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
			orchestrator.WithMarketMap(marketMap),
		)
		require.NoError(t, err)
		require.NoError(t, o.Init())

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		providerState, ok := providers[coinbase.Name]
		require.True(t, ok)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// Start the provider.
		go func() {
			providerState.Provider.Start(ctx)
		}()

		time.Sleep(1000 * time.Millisecond)

		// Check the state before any modifications are done.
		checkProviderState(t, expectedTickers, coinbase.Name, true, providertypes.API, true, providerState)

		pMarketMap := types.ProviderMarketMap{
			Name: coinbase.Name,
		}
		updatedState, err := o.UpdateProviderState(pMarketMap, providerState)
		require.NoError(t, err)

		time.Sleep(1000 * time.Millisecond)

		// Check the state after the update.
		checkProviderState(t, nil, coinbase.Name, false, providertypes.API, true, updatedState)

		updatedState.Provider.Stop()
		time.Sleep(2000 * time.Millisecond)
		require.Eventually(
			t,
			func() bool {
				return !updatedState.Provider.IsRunning()
			},
			5*time.Second,
			500*time.Millisecond,
		)
	})

	t.Run("can update a single websocket provider state with no configuration and non-running", func(t *testing.T) {
		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfg,
			orchestrator.WithLogger(logger),
			orchestrator.WithAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		require.NoError(t, o.Init())

		providerMarketMap, err := types.ProviderMarketMapFromMarketMap(okx.Name, marketMap)
		require.NoError(t, err)

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		providerState, ok := providers[okx.Name]
		require.True(t, ok)

		// Check the state before any modifications are done.
		checkProviderState(t, nil, okx.Name, false, providertypes.WebSockets, false, providerState)

		updatedState, err := o.UpdateProviderState(providerMarketMap, providerState)
		require.NoError(t, err)

		// Check the state after the update.
		checkProviderState(t, expectedTickers, okx.Name, true, providertypes.WebSockets, false, updatedState)

		updatedState.Provider.Stop()
		time.Sleep(2000 * time.Millisecond)
		require.Eventually(
			t,
			func() bool {
				return !updatedState.Provider.IsRunning()
			},
			5*time.Second,
			500*time.Millisecond,
		)
	})

	t.Run("can update a single websocket provider state with no configuration and running", func(t *testing.T) {
		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfg,
			orchestrator.WithLogger(logger),
			orchestrator.WithAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		require.NoError(t, o.Init())

		providerMarketMap, err := types.ProviderMarketMapFromMarketMap(okx.Name, marketMap)
		require.NoError(t, err)

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		providerState, ok := providers[okx.Name]
		require.True(t, ok)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Start the provider.
		go func() {
			providerState.Provider.Start(ctx)
		}()

		time.Sleep(3 * time.Millisecond)

		// Check the state before any modifications are done.
		checkProviderState(t, nil, okx.Name, false, providertypes.WebSockets, true, providerState)

		updatedState, err := o.UpdateProviderState(providerMarketMap, providerState)
		require.NoError(t, err)

		time.Sleep(3 * time.Millisecond)

		// Check the state after the update.
		checkProviderState(t, expectedTickers, okx.Name, true, providertypes.WebSockets, true, updatedState)

		updatedState.Provider.Stop()
		time.Sleep(2000 * time.Millisecond)
		require.Eventually(
			t,
			func() bool {
				return !updatedState.Provider.IsRunning()
			},
			10*time.Second,
			500*time.Millisecond,
		)
	})

	t.Run("can update a single websocket provider state removing all tickers on a non-running provider", func(t *testing.T) {
		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfg,
			orchestrator.WithLogger(logger),
			orchestrator.WithAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
			orchestrator.WithMarketMap(marketMap),
		)
		require.NoError(t, err)
		require.NoError(t, o.Init())

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		providerState, ok := providers[okx.Name]
		require.True(t, ok)

		// Check the state before any modifications are done.
		checkProviderState(t, expectedTickers, okx.Name, true, providertypes.WebSockets, false, providerState)

		pMarketMap := types.ProviderMarketMap{
			Name: okx.Name,
		}
		updatedState, err := o.UpdateProviderState(pMarketMap, providerState)
		require.NoError(t, err)

		time.Sleep(1000 * time.Millisecond)

		// Check the state after the update.
		checkProviderState(t, nil, okx.Name, false, providertypes.WebSockets, false, updatedState)

		updatedState.Provider.Stop()
		time.Sleep(2000 * time.Millisecond)
		require.Eventually(
			t,
			func() bool {
				return !updatedState.Provider.IsRunning()
			},
			5*time.Second,
			500*time.Millisecond,
		)
	})

	t.Run("can update a single websocket provider state removing all tickers on a running provider", func(t *testing.T) {
		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfg,
			orchestrator.WithLogger(logger),
			orchestrator.WithAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
			orchestrator.WithMarketMap(marketMap),
		)
		require.NoError(t, err)
		require.NoError(t, o.Init())

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		providerState, ok := providers[okx.Name]
		require.True(t, ok)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// Start the provider.
		go func() {
			providerState.Provider.Start(ctx)
		}()

		time.Sleep(1000 * time.Millisecond)

		// Check the state before any modifications are done.
		checkProviderState(t, expectedTickers, okx.Name, true, providertypes.WebSockets, true, providerState)

		pMarketMap := types.ProviderMarketMap{
			Name: okx.Name,
		}
		updatedState, err := o.UpdateProviderState(pMarketMap, providerState)
		require.NoError(t, err)

		time.Sleep(1000 * time.Millisecond)

		// Check the state after the update.
		checkProviderState(t, nil, okx.Name, false, providertypes.WebSockets, true, updatedState)

		updatedState.Provider.Stop()
		time.Sleep(2000 * time.Millisecond)
		require.Eventually(
			t,
			func() bool {
				return !updatedState.Provider.IsRunning()
			},
			10*time.Second,
			500*time.Millisecond,
		)
	})
}
