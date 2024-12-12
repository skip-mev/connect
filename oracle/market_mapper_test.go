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

	"github.com/skip-mev/slinky/oracle"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	oraclefactory "github.com/skip-mev/slinky/providers/factories/oracle"
	"github.com/skip-mev/slinky/providers/providertest"
	mmclienttypes "github.com/skip-mev/slinky/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	btcusdt = mmtypes.Market{
		Ticker: mmtypes.Ticker{
			CurrencyPair: slinkytypes.CurrencyPair{
				Base:  "BITCOIN",
				Quote: "USDT",
			},
			Decimals:         8,
			MinProviderCount: 1,
		},
		ProviderConfigs: []mmtypes.ProviderConfig{
			{
				Name:           "kucoin",
				OffChainTicker: "btc-usdt",
			},
		},
	}

	usdtusd = mmtypes.Market{
		Ticker: mmtypes.Ticker{
			CurrencyPair: slinkytypes.CurrencyPair{
				Base:  "USDT",
				Quote: "USD",
			},
			Decimals:         8,
			MinProviderCount: 1,
		},
		ProviderConfigs: []mmtypes.ProviderConfig{
			{
				Name:           "kucoin",
				OffChainTicker: "usdt-usd",
			},
		},
	}

	usdcusd = mmtypes.Market{
		Ticker: mmtypes.Ticker{
			CurrencyPair: slinkytypes.CurrencyPair{
				Base:  "USDC",
				Quote: "USD",
			},
			Decimals:         8,
			MinProviderCount: 1,
		},
		ProviderConfigs: []mmtypes.ProviderConfig{
			{
				Name:           "kucoin",
				OffChainTicker: "usdc-usd",
			},
		},
	}

	ethusdt = mmtypes.Market{
		Ticker: mmtypes.Ticker{
			CurrencyPair: slinkytypes.CurrencyPair{
				Base:  "ETHEREUM",
				Quote: "USDT",
			},
			Decimals:         8,
			MinProviderCount: 1,
		},
		ProviderConfigs: []mmtypes.ProviderConfig{
			{
				Name:           "kucoin",
				OffChainTicker: "eth-usdt",
				// include a normalize pair
				NormalizeByPair: &usdcusd.Ticker.CurrencyPair,
			},
		},
	}

	marketsMap = map[string]mmtypes.Market{
		btcusdt.Ticker.String(): btcusdt,
		usdcusd.Ticker.String(): usdcusd,
		usdtusd.Ticker.String(): usdtusd,
		ethusdt.Ticker.String(): ethusdt,
	}

	// invalid because we are excluding the usdcusd pair which
	// is used as a normalization in ethusdt.
	marketsMapInvalid = map[string]mmtypes.Market{
		btcusdt.Ticker.String(): btcusdt,
		usdtusd.Ticker.String(): usdtusd,
		ethusdt.Ticker.String(): ethusdt,
	}

	// remove the ethusdt which was requiring a normalization pair that wasn't in the map.
	marketsMapValidSubset = map[string]mmtypes.Market{
		btcusdt.Ticker.String(): btcusdt,
		usdtusd.Ticker.String(): usdtusd,
	}
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

func TestOracleImpl_IsMarketMapValidUpdated(t *testing.T) {
	tests := []struct {
		name             string
		resp             *mmtypes.MarketMapResponse
		initialMarketMap mmtypes.MarketMap
		lastUpdated      uint64
		wantMM           mmtypes.MarketMap
		wantUpdated      bool
		wantErr          bool
	}{
		{
			name:    "error on nil response",
			wantErr: true,
		},
		{
			name:        "do nothing on empty response - no initial state",
			resp:        &mmtypes.MarketMapResponse{},
			wantErr:     false,
			wantUpdated: false,
		},
		{
			name: "response is empty - initial state - update to empty",
			initialMarketMap: mmtypes.MarketMap{
				Markets: marketsMap,
			},
			resp:        &mmtypes.MarketMapResponse{},
			wantErr:     false,
			wantUpdated: true,
			wantMM: mmtypes.MarketMap{
				Markets: map[string]mmtypes.Market{},
			},
		},
		{
			name: "response is equal - initial state - do nothing",
			initialMarketMap: mmtypes.MarketMap{
				Markets: marketsMap,
			},
			resp: &mmtypes.MarketMapResponse{
				MarketMap: mmtypes.MarketMap{
					Markets: marketsMap,
				},
			},
			wantErr:     false,
			wantUpdated: false,
		},
		{
			name: "response is invalid - initial state - update to valid",
			initialMarketMap: mmtypes.MarketMap{
				Markets: marketsMap,
			},
			resp: &mmtypes.MarketMapResponse{
				MarketMap: mmtypes.MarketMap{
					Markets: marketsMapInvalid,
				},
			},
			wantErr:     false,
			wantUpdated: true,
			wantMM: mmtypes.MarketMap{
				Markets: marketsMapValidSubset,
			},
		},
		{
			name: "response is invalid - no initial state - update to valid",
			resp: &mmtypes.MarketMapResponse{
				MarketMap: mmtypes.MarketMap{
					Markets: marketsMapInvalid,
				},
			},
			wantErr:     false,
			wantUpdated: true,
			wantMM: mmtypes.MarketMap{
				Markets: marketsMapValidSubset,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			o, err := providertest.NewTestingOracle(ctx, []string{},
				oracle.WithLastUpdated(tt.lastUpdated),
				oracle.WithMarketMap(tt.initialMarketMap),
			)
			require.NoError(t, err)
			gotMM, isUpdated, err := o.Oracle.IsMarketMapValidUpdated(tt.resp)

			if tt.wantErr {
				require.Error(t, err)
				require.False(t, isUpdated)
				return
			}

			require.NoError(t, err)

			if tt.wantUpdated {
				require.Equal(t, tt.wantMM, gotMM)
				require.Equal(t, tt.wantUpdated, isUpdated)
				return
			}

			require.Equal(t, tt.wantUpdated, isUpdated)
		})
	}
}
