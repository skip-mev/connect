package oracle_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle"
	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/apis/binance"
	"github.com/skip-mev/connect/v2/providers/apis/coinbase"
	oraclefactory "github.com/skip-mev/connect/v2/providers/factories/oracle"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
	"github.com/skip-mev/connect/v2/providers/websockets/okx"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

func TestUpdateWithMarketMap(t *testing.T) {
	t.Run("bad market map is rejected", func(t *testing.T) {
		orc, err := oracle.New(
			oracleCfg,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		o := orc.(*oracle.OracleImpl)
		require.NoError(t, o.Init(context.Background()))

		err = o.UpdateMarketMap(mmtypes.MarketMap{
			Markets: map[string]mmtypes.Market{
				"bad": {},
			},
		})
		require.Error(t, err)

		o.Stop()
	})

	t.Run("can update the oracle's market map and update the providers' market maps with no running providers", func(t *testing.T) {
		orc, err := oracle.New(
			oracleCfg,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		o := orc.(*oracle.OracleImpl)
		require.NoError(t, o.Init(context.TODO()))

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		// Update the oracle's market map.
		require.NoError(t, o.UpdateMarketMap(marketMap))

		providers = o.GetProviderState()

		cbTickers, err := types.ProviderTickersFromMarketMap(coinbase.Name, marketMap)
		require.NoError(t, err)

		// Check the state after the update.
		coinbaseState, ok := providers[coinbase.Name]
		require.True(t, ok)
		checkProviderState(
			t,
			cbTickers,
			coinbase.Name,
			providertypes.API,
			false,
			coinbaseState,
		)

		okxTickers, err := types.ProviderTickersFromMarketMap(okx.Name, marketMap)
		require.NoError(t, err)

		okxState, ok := providers[okx.Name]
		require.True(t, ok)
		checkProviderState(
			t,
			okxTickers,
			okx.Name,
			providertypes.WebSockets,
			false,
			okxState,
		)

		binanceState, ok := providers[binance.Name]
		require.True(t, ok)
		checkProviderState(t, nil, binance.Name, providertypes.API, false, binanceState)

		o.Stop()
	})

	t.Run("can update the oracle's market map and update the providers' market maps with running providers", func(t *testing.T) {
		orc, err := oracle.New(
			oracleCfg,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		o := orc.(*oracle.OracleImpl)

		// Start the providers.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		go func() {
			require.ErrorIs(t, o.Start(ctx), context.Canceled)
		}()

		time.Sleep(2 * time.Second)

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		// Update the oracle's market map.
		require.NoError(t, o.UpdateMarketMap(marketMap))

		time.Sleep(2 * time.Second)

		providers = o.GetProviderState()
		require.Len(t, providers, 3)

		cbTickers, err := types.ProviderTickersFromMarketMap(coinbase.Name, marketMap)
		require.NoError(t, err)

		// Check the state after the update.
		coinbaseState, ok := providers[coinbase.Name]
		require.True(t, ok)
		checkProviderState(
			t,
			cbTickers,
			coinbase.Name,
			providertypes.API,
			true,
			coinbaseState,
		)

		okxTickers, err := types.ProviderTickersFromMarketMap(okx.Name, marketMap)
		require.NoError(t, err)

		okxState, ok := providers[okx.Name]
		require.True(t, ok)
		checkProviderState(
			t,
			okxTickers,
			okx.Name,
			providertypes.WebSockets,
			true,
			okxState,
		)

		binanceState, ok := providers[binance.Name]
		require.True(t, ok)
		checkProviderState(
			t,
			nil,
			binance.Name,
			providertypes.API,
			false,
			binanceState,
		)

		// Stop the providers.
		o.Stop()
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

	t.Run("can update the oracle's market map and update the providers' market maps with no tickers", func(t *testing.T) {
		orc, err := oracle.New(
			oracleCfg,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		o := orc.(*oracle.OracleImpl)
		require.NoError(t, o.Init(context.Background()))

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		// Update the oracle's market map.
		require.NoError(t, o.UpdateMarketMap(mmtypes.MarketMap{}))

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

		o.Stop()
	})

	t.Run("can update the oracle's market map and update the providers' market maps with no tickers and running providers", func(t *testing.T) {
		orc, err := oracle.New(
			oracleCfg,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
			oracle.WithMarketMap(marketMap),
		)
		require.NoError(t, err)
		o := orc.(*oracle.OracleImpl)

		// Start the providers.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		go func() {
			require.ErrorIs(t, o.Start(ctx), context.Canceled)
		}()

		time.Sleep(2 * time.Second)
		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		// Update the oracle's market map.
		require.NoError(t, o.UpdateMarketMap(mmtypes.MarketMap{}))

		time.Sleep(2 * time.Second)

		providers = o.GetProviderState()
		require.Len(t, providers, 3)

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

		// Stop the providers.
		o.Stop()
	})
}

func TestUpdateProviderState(t *testing.T) {
	t.Run("can update a single api provider state with no configuration and non-running", func(t *testing.T) {
		orc, err := oracle.New(
			oracleCfg,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		o := orc.(*oracle.OracleImpl)
		require.NoError(t, o.Init(context.TODO()))

		tickers, err := types.ProviderTickersFromMarketMap(coinbase.Name, marketMap)
		require.NoError(t, err)

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		providerState, ok := providers[coinbase.Name]
		require.True(t, ok)

		// Check the state before any modifications are done.
		checkProviderState(t, nil, coinbase.Name, providertypes.API, false, providerState)

		updatedState, err := o.UpdateProviderState(tickers, providerState)
		require.NoError(t, err)

		// Check the state after the update.
		checkProviderState(t, tickers, coinbase.Name, providertypes.API, false, updatedState)
	})

	t.Run("can update a single api provider state with no configuration and running", func(t *testing.T) {
		orc, err := oracle.New(
			oracleCfg,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		o := orc.(*oracle.OracleImpl)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Start the provider.
		go func() {
			require.ErrorIs(t, o.Start(ctx), context.Canceled)
		}()

		time.Sleep(500 * time.Millisecond)

		tickers, err := types.ProviderTickersFromMarketMap(coinbase.Name, marketMap)
		require.NoError(t, err)

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		providerState, ok := providers[coinbase.Name]
		require.True(t, ok)

		// Check the state before any modifications are done.
		checkProviderState(t, nil, coinbase.Name, providertypes.API, false, providerState)

		updatedState, err := o.UpdateProviderState(tickers, providerState)
		require.NoError(t, err)

		time.Sleep(1000 * time.Millisecond)

		// Check the state after the update.
		checkProviderState(t, tickers, coinbase.Name, providertypes.API, true, updatedState)

		o.Stop()
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
		orc, err := oracle.New(
			oracleCfg,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
			oracle.WithMarketMap(marketMap),
		)
		require.NoError(t, err)
		o := orc.(*oracle.OracleImpl)
		require.NoError(t, o.Init(context.TODO()))

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		providerState, ok := providers[coinbase.Name]
		require.True(t, ok)

		tickers, err := types.ProviderTickersFromMarketMap(coinbase.Name, marketMap)
		require.NoError(t, err)

		// Check the state before any modifications are done.
		checkProviderState(t, tickers, coinbase.Name, providertypes.API, false, providerState)

		updatedState, err := o.UpdateProviderState(nil, providerState)
		require.NoError(t, err)

		time.Sleep(1000 * time.Millisecond)

		// Check the state after the update.
		checkProviderState(t, nil, coinbase.Name, providertypes.API, false, updatedState)

		o.Stop()
	})

	t.Run("can update a single api provider state removing all tickers on a running provider", func(t *testing.T) {
		orc, err := oracle.New(
			oracleCfg,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
			oracle.WithMarketMap(marketMap),
		)
		require.NoError(t, err)
		o := orc.(*oracle.OracleImpl)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// Start the provider.
		go func() {
			require.ErrorIs(t, o.Start(ctx), context.Canceled)
		}()

		time.Sleep(1000 * time.Millisecond)
		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		providerState, ok := providers[coinbase.Name]
		require.True(t, ok)

		tickers, err := types.ProviderTickersFromMarketMap(coinbase.Name, marketMap)
		require.NoError(t, err)

		// Check the state before any modifications are done.
		checkProviderState(t, tickers, coinbase.Name, providertypes.API, true, providerState)
		updatedState, err := o.UpdateProviderState(nil, providerState)
		require.NoError(t, err)

		time.Sleep(1000 * time.Millisecond)

		// Check the state after the update.
		checkProviderState(t, nil, coinbase.Name, providertypes.API, false, updatedState)

		o.Stop()
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
		orc, err := oracle.New(
			oracleCfg,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		o := orc.(*oracle.OracleImpl)
		require.NoError(t, o.Init(context.Background()))

		tickers, err := types.ProviderTickersFromMarketMap(coinbase.Name, marketMap)
		require.NoError(t, err)

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		providerState, ok := providers[okx.Name]
		require.True(t, ok)

		// Check the state before any modifications are done.
		checkProviderState(t, nil, okx.Name, providertypes.WebSockets, false, providerState)

		updatedState, err := o.UpdateProviderState(tickers, providerState)
		require.NoError(t, err)

		// Check the state after the update.
		checkProviderState(t, tickers, okx.Name, providertypes.WebSockets, false, updatedState)

		o.Stop()
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
		orc, err := oracle.New(
			oracleCfg,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		o := orc.(*oracle.OracleImpl)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Start the provider.
		go func() {
			require.ErrorIs(t, o.Start(ctx), context.Canceled)
		}()

		time.Sleep(3 * time.Millisecond)
		tickers, err := types.ProviderTickersFromMarketMap(okx.Name, marketMap)
		require.NoError(t, err)

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		providerState, ok := providers[okx.Name]
		require.True(t, ok)

		// Check the state before any modifications are done.
		checkProviderState(t, nil, okx.Name, providertypes.WebSockets, false, providerState)

		updatedState, err := o.UpdateProviderState(tickers, providerState)
		require.NoError(t, err)

		time.Sleep(3 * time.Millisecond)

		// Check the state after the update.
		checkProviderState(t, tickers, okx.Name, providertypes.WebSockets, true, updatedState)

		o.Stop()
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
		orc, err := oracle.New(
			oracleCfg,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
			oracle.WithMarketMap(marketMap),
		)
		require.NoError(t, err)
		o := orc.(*oracle.OracleImpl)

		require.NoError(t, o.Init(context.TODO()))

		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		providerState, ok := providers[okx.Name]
		require.True(t, ok)

		tickers, err := types.ProviderTickersFromMarketMap(okx.Name, marketMap)
		require.NoError(t, err)

		// Check the state before any modifications are done.
		checkProviderState(t, tickers, okx.Name, providertypes.WebSockets, false, providerState)
		updatedState, err := o.UpdateProviderState(nil, providerState)
		require.NoError(t, err)

		time.Sleep(1000 * time.Millisecond)

		// Check the state after the update.
		checkProviderState(t, nil, okx.Name, providertypes.WebSockets, false, updatedState)

		o.Stop()
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
		orc, err := oracle.New(
			oracleCfg,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
			oracle.WithMarketMap(marketMap),
		)
		require.NoError(t, err)
		o := orc.(*oracle.OracleImpl)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// Start the provider.
		go func() {
			require.ErrorIs(t, o.Start(ctx), context.Canceled)
		}()

		time.Sleep(1000 * time.Millisecond)
		providers := o.GetProviderState()
		require.Len(t, providers, 3)

		providerState, ok := providers[okx.Name]
		require.True(t, ok)

		tickers, err := types.ProviderTickersFromMarketMap(okx.Name, marketMap)
		require.NoError(t, err)

		// Check the state before any modifications are done.
		checkProviderState(t, tickers, okx.Name, providertypes.WebSockets, true, providerState)
		updatedState, err := o.UpdateProviderState(nil, providerState)
		require.NoError(t, err)

		time.Sleep(1000 * time.Millisecond)

		// Check the state after the update.
		checkProviderState(t, nil, okx.Name, providertypes.WebSockets, false, updatedState)

		o.Stop()
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
