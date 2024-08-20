package oracle_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle"
	oraclefactory "github.com/skip-mev/connect/v2/providers/factories/oracle"
)

func TestStart(t *testing.T) {
	t.Run("errors when init fails", func(t *testing.T) {
		o, err := oracle.New(
			oracleCfg,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
		)
		require.NoError(t, err)

		err = o.Start(context.Background())
		require.Error(t, err)

		o.Stop()
	})

	t.Run("price providers with no market map", func(t *testing.T) {
		orc, err := oracle.New(
			oracleCfg,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		o := orc.(*oracle.OracleImpl)

		go func() {
			err := o.Start(context.Background())
			if !errors.Is(err, context.Canceled) {
				t.Errorf("Start() should have returned context.Canceled error")
			}
		}()

		time.Sleep(5 * time.Second)
		state := o.GetProviderState()
		require.Equal(t, len(state), len(oracleCfg.Providers))

		// Stop the oracle.
		o.Stop()
	})

	t.Run("price providers with market map", func(t *testing.T) {
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

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() {
			err := o.Start(ctx)
			if !errors.Is(err, context.Canceled) {
				t.Errorf("Start() should have returned context.Canceled error")
			}
		}()

		time.Sleep(5 * time.Second)
		state := o.GetProviderState()
		require.Equal(t, len(state), len(oracleCfg.Providers))

		// Stop the oracle.
		o.Stop()
	})

	t.Run("price providers with market map provider but price providers have no ids to fetch", func(t *testing.T) {
		orc, err := oracle.New(
			oracleCfgWithMapper,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
			oracle.WithMarketMapperFactory(oraclefactory.MarketMapProviderFactory),
		)
		require.NoError(t, err)
		o := orc.(*oracle.OracleImpl)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() {
			err := o.Start(ctx)
			if !errors.Is(err, context.Canceled) {
				t.Errorf("Start() should have returned context.Canceled error")
			}
		}()

		time.Sleep(5 * time.Second)
		state := o.GetProviderState()
		require.Equal(t, len(state), len(oracleCfgWithMapper.Providers)-1)

		mapper := o.GetMarketMapProvider()
		require.NotNil(t, mapper)

		// Stop the oracle.
		o.Stop()
	})
}
