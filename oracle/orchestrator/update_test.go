package orchestrator_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/orchestrator"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/apis/binance"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	oraclefactory "github.com/skip-mev/slinky/providers/factories/oracle"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/websockets/okx"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

func TestUpdateWithMarketMap(t *testing.T) {
	t.Run("bad market map is reject", func(t *testing.T) {
		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfg,
			orchestrator.WithLogger(logger),
			orchestrator.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		require.NoError(t, o.Init())

		err = o.UpdateWithMarketMap(mmtypes.MarketMap{
			Markets: map[string]mmtypes.Market{
				"bad": {},
			},
		})
		require.Error(t, err)
	})

	t.Run("can update the orchestrator's market map and update the providers' market maps with no running providers", func(t *testing.T) {
		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfg,
			orchestrator.WithLogger(logger),
			orchestrator.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		require.NoError(t, o.Init())

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		// Update the orchestrator's market map.
		require.NoError(t, o.UpdateWithMarketMap(marketMap))

		providers = o.GetProviderState()

		// Check the state after the update.
		coinbaseState, ok := providers[coinbase.Name]
		require.True(t, ok)
		checkProviderState(t, []mmtypes.Ticker{constants.BITCOIN_USD, constants.ETHEREUM_USD}, coinbase.Name, providertypes.API, false, coinbaseState)

		okxState, ok := providers[okx.Name]
		require.True(t, ok)
		checkProviderState(t, []mmtypes.Ticker{constants.BITCOIN_USD, constants.ETHEREUM_USD}, okx.Name, providertypes.WebSockets, false, okxState)

		binanceState, ok := providers[binance.Name]
		require.True(t, ok)
		checkProviderState(t, nil, binance.Name, providertypes.API, false, binanceState)
	})

	t.Run("can update the orchestrator's market map and update the providers' market maps with running providers", func(t *testing.T) {
		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfg,
			orchestrator.WithLogger(logger),
			orchestrator.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		require.NoError(t, o.Init())

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		// Start the providers.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for _, state := range providers {
			go func(s orchestrator.ProviderState) {
				s.Provider.Start(ctx)
			}(state)
		}

		time.Sleep(1000 * time.Millisecond)

		// Update the orchestrator's market map.
		require.NoError(t, o.UpdateWithMarketMap(marketMap))

		time.Sleep(2000 * time.Millisecond)

		providers = o.GetProviderState()
		require.Len(t, providers, 3)

		// Check the state after the update.
		coinbaseState, ok := providers[coinbase.Name]
		require.True(t, ok)
		checkProviderState(t, []mmtypes.Ticker{constants.BITCOIN_USD, constants.ETHEREUM_USD}, coinbase.Name, providertypes.API, true, coinbaseState)

		okxState, ok := providers[okx.Name]
		require.True(t, ok)
		checkProviderState(t, []mmtypes.Ticker{constants.BITCOIN_USD, constants.ETHEREUM_USD}, okx.Name, providertypes.WebSockets, true, okxState)

		binanceState, ok := providers[binance.Name]
		require.True(t, ok)
		checkProviderState(t, nil, binance.Name, providertypes.API, true, binanceState)

		// Stop the providers.
		for _, state := range providers {
			state.Provider.Stop()
		}

		time.Sleep(2000 * time.Millisecond)

		// Ensure all providers are stopped.
		for _, state := range providers {
			require.Eventually(
				t,
				func() bool {
					return !state.Provider.IsRunning()
				},
				5*time.Second,
				500*time.Millisecond,
			)
		}
	})

	t.Run("can update the orchestrator's market map and update the providers' market maps with no tickers", func(t *testing.T) {
		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfg,
			orchestrator.WithLogger(logger),
			orchestrator.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		require.NoError(t, o.Init())

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		// Update the orchestrator's market map.
		require.NoError(t, o.UpdateWithMarketMap(mmtypes.MarketMap{}))

		providers = o.GetProviderState()

		// Check the state after the update.
		coinbaseState, ok := providers[coinbase.Name]
		require.True(t, ok)
		checkProviderState(t, nil, coinbase.Name, providertypes.API, false, coinbaseState)

		okxState, ok := providers[okx.Name]
		require.True(t, ok)
		checkProviderState(t, nil, okx.Name, providertypes.WebSockets, false, okxState)

		binanceState, ok := providers[binance.Name]
		require.True(t, ok)
		checkProviderState(t, nil, binance.Name, providertypes.API, false, binanceState)
	})

	t.Run("can update the orchestrator's market map and update the providers' market maps with no tickers and running providers", func(t *testing.T) {
		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfg,
			orchestrator.WithLogger(logger),
			orchestrator.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
			orchestrator.WithMarketMap(marketMap),
		)
		require.NoError(t, err)
		require.NoError(t, o.Init())

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		// Start the providers.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for _, state := range providers {
			go func(s orchestrator.ProviderState) {
				s.Provider.Start(ctx)
			}(state)
		}

		time.Sleep(1000 * time.Millisecond)

		// Update the orchestrator's market map.
		require.NoError(t, o.UpdateWithMarketMap(mmtypes.MarketMap{}))

		time.Sleep(2000 * time.Millisecond)

		providers = o.GetProviderState()
		require.Len(t, providers, 3)

		// Check the state after the update.
		coinbaseState, ok := providers[coinbase.Name]
		require.True(t, ok)
		checkProviderState(t, nil, coinbase.Name, providertypes.API, true, coinbaseState)

		okxState, ok := providers[okx.Name]
		require.True(t, ok)
		checkProviderState(t, nil, okx.Name, providertypes.WebSockets, true, okxState)

		binanceState, ok := providers[binance.Name]
		require.True(t, ok)
		checkProviderState(t, nil, binance.Name, providertypes.API, true, binanceState)

		// Stop the providers.
		for _, state := range providers {
			state.Provider.Stop()
		}

		time.Sleep(2000 * time.Millisecond)

		// Ensure all providers are stopped.
		for _, state := range providers {
			require.Eventually(
				t,
				func() bool {
					return !state.Provider.IsRunning()
				},
				5*time.Second,
				500*time.Millisecond,
			)
		}
	})
}

func TestUpdateProviderState(t *testing.T) {
	expectedTickers := []mmtypes.Ticker{constants.BITCOIN_USD, constants.ETHEREUM_USD}

	t.Run("can update a single api provider state with no configuration and non-running", func(t *testing.T) {
		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfg,
			orchestrator.WithLogger(logger),
			orchestrator.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
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
		checkProviderState(t, nil, coinbase.Name, providertypes.API, false, providerState)

		updatedState, err := o.UpdateProviderState(providerMarketMap, providerState)
		require.NoError(t, err)

		// Check the state after the update.
		checkProviderState(t, expectedTickers, coinbase.Name, providertypes.API, false, updatedState)

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
			orchestrator.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
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
		checkProviderState(t, nil, coinbase.Name, providertypes.API, true, providerState)

		updatedState, err := o.UpdateProviderState(providerMarketMap, providerState)
		require.NoError(t, err)

		time.Sleep(1000 * time.Millisecond)

		// Check the state after the update.
		checkProviderState(t, expectedTickers, coinbase.Name, providertypes.API, true, updatedState)

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
			orchestrator.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
			orchestrator.WithMarketMap(marketMap),
		)
		require.NoError(t, err)
		require.NoError(t, o.Init())

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		providerState, ok := providers[coinbase.Name]
		require.True(t, ok)

		// Check the state before any modifications are done.
		checkProviderState(t, expectedTickers, coinbase.Name, providertypes.API, false, providerState)

		pMarketMap := types.ProviderMarketMap{
			Name: coinbase.Name,
		}
		updatedState, err := o.UpdateProviderState(pMarketMap, providerState)
		require.NoError(t, err)

		time.Sleep(1000 * time.Millisecond)

		// Check the state after the update.
		checkProviderState(t, nil, coinbase.Name, providertypes.API, false, updatedState)

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
			orchestrator.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
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
		checkProviderState(t, expectedTickers, coinbase.Name, providertypes.API, true, providerState)

		pMarketMap := types.ProviderMarketMap{
			Name: coinbase.Name,
		}
		updatedState, err := o.UpdateProviderState(pMarketMap, providerState)
		require.NoError(t, err)

		time.Sleep(1000 * time.Millisecond)

		// Check the state after the update.
		checkProviderState(t, nil, coinbase.Name, providertypes.API, true, updatedState)

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
			orchestrator.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
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
		checkProviderState(t, nil, okx.Name, providertypes.WebSockets, false, providerState)

		updatedState, err := o.UpdateProviderState(providerMarketMap, providerState)
		require.NoError(t, err)

		// Check the state after the update.
		checkProviderState(t, expectedTickers, okx.Name, providertypes.WebSockets, false, updatedState)

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
			orchestrator.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
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
		checkProviderState(t, nil, okx.Name, providertypes.WebSockets, true, providerState)

		updatedState, err := o.UpdateProviderState(providerMarketMap, providerState)
		require.NoError(t, err)

		time.Sleep(3 * time.Millisecond)

		// Check the state after the update.
		checkProviderState(t, expectedTickers, okx.Name, providertypes.WebSockets, true, updatedState)

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
			orchestrator.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
			orchestrator.WithMarketMap(marketMap),
		)
		require.NoError(t, err)
		require.NoError(t, o.Init())

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		providerState, ok := providers[okx.Name]
		require.True(t, ok)

		// Check the state before any modifications are done.
		checkProviderState(t, expectedTickers, okx.Name, providertypes.WebSockets, false, providerState)

		pMarketMap := types.ProviderMarketMap{
			Name: okx.Name,
		}
		updatedState, err := o.UpdateProviderState(pMarketMap, providerState)
		require.NoError(t, err)

		time.Sleep(1000 * time.Millisecond)

		// Check the state after the update.
		checkProviderState(t, nil, okx.Name, providertypes.WebSockets, false, updatedState)

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
			orchestrator.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
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
		checkProviderState(t, expectedTickers, okx.Name, providertypes.WebSockets, true, providerState)

		pMarketMap := types.ProviderMarketMap{
			Name: okx.Name,
		}
		updatedState, err := o.UpdateProviderState(pMarketMap, providerState)
		require.NoError(t, err)

		time.Sleep(1000 * time.Millisecond)

		// Check the state after the update.
		checkProviderState(t, nil, okx.Name, providertypes.WebSockets, true, updatedState)

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
