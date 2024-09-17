package oracle_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle"
	oraclefactory "github.com/skip-mev/connect/v2/providers/factories/oracle"
	mmclienttypes "github.com/skip-mev/connect/v2/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

func TestListenForMarketMapUpdates(t *testing.T) {
	t.Run("mapper has no chain IDs to fetch should not update the oracle", func(t *testing.T) {
		handler, factory := marketMapperFactory(t, nil)
		handler.On("CreateURL", mock.Anything).Return("", fmt.Errorf("no ids")).Maybe()

		o, err := oracle.New(
			oracleCfgWithOnlyMockMapper,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithMarketMapperFactory(factory),
		)
		require.NoError(t, err)
		current := o.GetMarketMap()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			err := o.Start(ctx)
			if !errors.Is(err, context.Canceled) {
				t.Errorf("Start() should have returned context.Canceled error")
			}
		}()

		// Wait for the oracle to start.
		time.Sleep(1000 * time.Millisecond)

		// The oracle should not have been updated.
		require.Equal(t, current, o.GetMarketMap())

		// Stop the oracle.
		cancel()
		o.Stop()
	})

	t.Run("mapper is responsible for more than one chain", func(t *testing.T) {
		handler, factory := marketMapperFactory(t, []mmclienttypes.Chain{{ChainID: "eth"}, {ChainID: "bsc"}})
		handler.On("CreateURL", mock.Anything).Return("", fmt.Errorf("too many")).Maybe()

		o, err := oracle.New(
			oracleCfgWithOnlyMockMapper,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithMarketMapperFactory(factory),
		)
		require.NoError(t, err)
		current := o.GetMarketMap()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() {
			err := o.Start(ctx)
			if !errors.Is(err, context.Canceled) {
				t.Errorf("Start() should have returned context.Canceled error")
			}
		}()

		// Wait for the oracle to start.
		time.Sleep(1000 * time.Millisecond)

		// The oracle should not have been updated.
		require.Equal(t, current, o.GetMarketMap())

		// Stop the oracle.
		cancel()
		o.Stop()
	})

	t.Run("mapper has a single chain ID but fails to get a any response for the chain", func(t *testing.T) {
		handler, factory := marketMapperFactory(t, []mmclienttypes.Chain{{ChainID: "dYdX"}})
		handler.On("CreateURL", mock.Anything).Return("", fmt.Errorf("failed to create url")).Maybe()

		o, err := oracle.New(
			oracleCfgWithOnlyMockMapper,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithMarketMapperFactory(factory),
		)
		require.NoError(t, err)
		current := o.GetMarketMap()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() {
			err := o.Start(ctx)
			if !errors.Is(err, context.Canceled) {
				t.Errorf("Start() should have returned context.Canceled error")
			}
		}()

		// Wait for the oracle to start.
		time.Sleep(2000 * time.Millisecond)

		// The oracle should not have been updated.
		require.Equal(t, current, o.GetMarketMap())

		// Stop the oracle.
		cancel()
		o.Stop()
	})

	t.Run("mapper gets response but it is the same as the current market map", func(t *testing.T) {
		chains := []mmclienttypes.Chain{{ChainID: "dYdX"}}
		handler, factory := marketMapperFactory(t, chains)
		handler.On("CreateURL", mock.Anything).Return("", nil).Maybe()

		resolved := make(mmclienttypes.ResolvedMarketMap)
		resp := mmtypes.MarketMapResponse{
			MarketMap: marketMap,
		}
		resolved[chains[0]] = mmclienttypes.NewMarketMapResult(&resp, time.Now())
		handler.On("ParseResponse", mock.Anything, mock.Anything).Return(mmclienttypes.NewMarketMapResponse(resolved, nil)).Maybe()

		o, err := oracle.New(
			oracleCfgWithOnlyMockMapper,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithMarketMapperFactory(factory),
			oracle.WithMarketMap(marketMap),
		)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			err := o.Start(ctx)
			if !errors.Is(err, context.Canceled) {
				t.Errorf("Start() should have returned context.Canceled error")
			}
		}()

		// Wait for the oracle to start.
		time.Sleep(2000 * time.Millisecond)

		// The oracle should not have been updated.
		require.Equal(t, marketMap, o.GetMarketMap())

		// Stop the oracle.
		cancel()
		o.Stop()
	})

	t.Run("mapper gets response and it is different from the current market map", func(t *testing.T) {
		chains := []mmclienttypes.Chain{{ChainID: "dYdX"}}
		handler, factory := marketMapperFactory(t, chains)
		handler.On("CreateURL", mock.Anything).Return("", nil).Maybe()

		resolved := make(mmclienttypes.ResolvedMarketMap)
		resp := mmtypes.MarketMapResponse{
			MarketMap: marketMap,
		}
		resolved[chains[0]] = mmclienttypes.NewMarketMapResult(&resp, time.Now())
		handler.On("ParseResponse", mock.Anything, mock.Anything).Return(mmclienttypes.NewMarketMapResponse(resolved, nil)).Maybe()

		o, err := oracle.New(
			oracleCfgWithOnlyMockMapper,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithMarketMapperFactory(factory),
		)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			err := o.Start(ctx)
			if !errors.Is(err, context.Canceled) {
				t.Errorf("Start() should have returned context.Canceled error")
			}
		}()

		// Wait for the oracle to start.
		time.Sleep(2000 * time.Millisecond)

		// The oracle should not have been updated.
		require.Equal(t, marketMap, o.GetMarketMap())

		// Stop the oracle.
		cancel()
		o.Stop()
	})

	t.Run("can update providers with a new market map", func(t *testing.T) {
		chains := []mmclienttypes.Chain{{ChainID: "dYdX"}}
		handler, factory := marketMapperFactory(t, chains)
		handler.On("CreateURL", mock.Anything).Return("", nil).Maybe()

		resolved := make(mmclienttypes.ResolvedMarketMap)
		resp := mmtypes.MarketMapResponse{
			MarketMap: marketMap,
		}
		resolved[chains[0]] = mmclienttypes.NewMarketMapResult(&resp, time.Now())
		handler.On("ParseResponse", mock.Anything, mock.Anything).Return(mmclienttypes.NewMarketMapResponse(resolved, nil)).Maybe()

		o, err := oracle.New(
			oracleCfgWithMockMapper,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithMarketMapperFactory(factory),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			err := o.Start(ctx)
			if !errors.Is(err, context.Canceled) {
				t.Errorf("Start() should have returned context.Canceled error")
			}
		}()

		// Wait for the oracle to start.
		time.Sleep(5000 * time.Millisecond)

		// The oracle should have been updated.
		require.Equal(t, marketMap, o.GetMarketMap())

		// Stop the oracle.
		cancel()
		o.Stop()
	})

	t.Run("can update providers with a new market map and write the updated market map", func(t *testing.T) {
		chains := []mmclienttypes.Chain{{ChainID: "dYdX"}}
		handler, factory := marketMapperFactory(t, chains)
		handler.On("CreateURL", mock.Anything).Return("", nil).Maybe()

		resolved := make(mmclienttypes.ResolvedMarketMap)
		resp := mmtypes.MarketMapResponse{
			MarketMap: marketMap,
		}
		resolved[chains[0]] = mmclienttypes.NewMarketMapResult(&resp, time.Now())
		handler.On("ParseResponse", mock.Anything, mock.Anything).Return(mmclienttypes.NewMarketMapResponse(resolved, nil)).Maybe()

		path := "test.json"
		o, err := oracle.New(
			oracleCfgWithMockMapper,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithMarketMapperFactory(factory),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
			oracle.WithWriteTo(path),
		)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			err := o.Start(ctx)
			if !errors.Is(err, context.Canceled) {
				t.Errorf("Start() should have returned context.Canceled error")
			}
		}()

		// Wait for the oracle to start.
		time.Sleep(5000 * time.Millisecond)

		// The oracle should have been updated.
		require.Equal(t, marketMap, o.GetMarketMap())

		// Stop the oracle.
		cancel()
		o.Stop()

		// Check that the market map was written to the path.
		mm, err := mmtypes.ReadMarketMapFromFile(path)
		require.NoError(t, err)
		require.Equal(t, o.GetMarketMap(), mm)

		// Clean up the file.
		require.NoError(t, os.Remove(path))
	})
	t.Run("can update providers with a new market map and handle partially invalid state", func(t *testing.T) {
		chains := []mmclienttypes.Chain{{ChainID: "dYdX"}}
		handler, factory := marketMapperFactory(t, chains)
		handler.On("CreateURL", mock.Anything).Return("", nil).Maybe()

		resolved := make(mmclienttypes.ResolvedMarketMap)
		resp := mmtypes.MarketMapResponse{
			MarketMap: partialInvalidMarketMap,
		}
		resolved[chains[0]] = mmclienttypes.NewMarketMapResult(&resp, time.Now())
		handler.On("ParseResponse", mock.Anything, mock.Anything).Return(mmclienttypes.NewMarketMapResponse(resolved, nil)).Maybe()

		o, err := oracle.New(
			oracleCfgWithMockMapper,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithMarketMapperFactory(factory),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			err := o.Start(ctx)
			if !errors.Is(err, context.Canceled) {
				t.Errorf("Start() should have returned context.Canceled error")
			}
		}()

		// Wait for the oracle to start.
		time.Sleep(5000 * time.Millisecond)

		// The oracle should have been updated.
		require.Equal(t, validMarketMapSubset, o.GetMarketMap())

		// Stop the oracle.
		cancel()
		o.Stop()
	})
}
