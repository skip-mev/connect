package orchestrator_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/orchestrator"
	oraclefactory "github.com/skip-mev/slinky/providers/factories/oracle"
)

func TestStart(t *testing.T) {
	t.Run("errors when init fails", func(t *testing.T) {
		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfg,
			orchestrator.WithLogger(logger),
		)
		require.NoError(t, err)

		err = o.Start(context.Background())
		require.Error(t, err)

		o.Stop()
	})

	t.Run("price providers with no market map", func(t *testing.T) {
		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfg,
			orchestrator.WithLogger(logger),
			orchestrator.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() {
			o.Start(ctx)
		}()

		time.Sleep(5 * time.Second)
		state := o.GetProviderState()
		require.Equal(t, len(state), len(oracleCfg.Providers))

		// Stop the provider orchestrator.
		o.Stop()
	})

	t.Run("price providers with market map", func(t *testing.T) {
		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfg,
			orchestrator.WithLogger(logger),
			orchestrator.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
			orchestrator.WithMarketMap(marketMap),
		)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() {
			o.Start(ctx)
		}()

		time.Sleep(5 * time.Second)
		state := o.GetProviderState()
		require.Equal(t, len(state), len(oracleCfg.Providers))

		// Stop the provider orchestrator.
		o.Stop()
	})

	t.Run("price providers with market map provider but price providers have no ids to fetch", func(t *testing.T) {
		o, err := orchestrator.NewProviderOrchestrator(
			oracleCfgWithMapper,
			orchestrator.WithLogger(logger),
			orchestrator.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			orchestrator.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
			orchestrator.WithMarketMapperFactory(oraclefactory.MarketMapProviderFactory),
		)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() {
			o.Start(ctx)
		}()

		time.Sleep(5 * time.Second)
		state := o.GetProviderState()
		require.Equal(t, len(state), len(oracleCfgWithMapper.Providers)-1)

		mapper := o.GetMarketMapProvider()
		require.NotNil(t, mapper)

		// Stop the provider orchestrator.
		o.Stop()
	})
}
