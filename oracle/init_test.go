package oracle_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle"
	"github.com/skip-mev/connect/v2/oracle/config"
	oracletypes "github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/apis/binance"
	"github.com/skip-mev/connect/v2/providers/apis/coinbase"
	oraclefactory "github.com/skip-mev/connect/v2/providers/factories/oracle"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
	"github.com/skip-mev/connect/v2/providers/websockets/okx"
)

var (
	coinbasebtcusd = oracletypes.DefaultProviderTicker{
		OffChainTicker: "BTCUSD",
	}
	coinbaseethusd = oracletypes.DefaultProviderTicker{
		OffChainTicker: "ETHUSD",
	}
	okxbtcusd = oracletypes.DefaultProviderTicker{
		OffChainTicker: "BTC-USD",
	}
	okxethusd = oracletypes.DefaultProviderTicker{
		OffChainTicker: "ETH-USD",
	}
)

func TestInit(t *testing.T) {
	t.Run("creates all providers without a marketmap", func(t *testing.T) {
		orc, err := oracle.New(
			oracleCfg,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		o := orc.(*oracle.OracleImpl)

		err = o.Init(context.TODO())
		require.NoError(t, err)

		state := o.GetProviderState()
		require.Equal(t, len(state), len(oracleCfg.Providers))

		coinbaseState, ok := state[coinbase.Name]
		require.True(t, ok)
		checkProviderState(t, nil, coinbase.Name, providertypes.API, false, coinbaseState)

		okxState, ok := state[okx.Name]
		require.True(t, ok)
		checkProviderState(t, nil, okx.Name, providertypes.WebSockets, false, okxState)

		binanceState, ok := state[binance.Name]
		require.True(t, ok)
		checkProviderState(t, nil, binance.Name, providertypes.API, false, binanceState)
	})

	t.Run("creates some providers with a marketmap", func(t *testing.T) {
		orc, err := oracle.New(
			oracleCfg,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithMarketMap(marketMap),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		o := orc.(*oracle.OracleImpl)

		err = o.Init(context.TODO())
		require.NoError(t, err)

		state := o.GetProviderState()
		require.Equal(t, len(state), len(oracleCfg.Providers))

		coinbaseState, ok := state[coinbase.Name]
		require.True(t, ok)
		checkProviderState(
			t,
			[]oracletypes.ProviderTicker{
				coinbasebtcusd,
				coinbaseethusd,
			},
			coinbase.Name,
			providertypes.API,
			false,
			coinbaseState,
		)

		okxState, ok := state[okx.Name]
		require.True(t, ok)
		checkProviderState(
			t,
			[]oracletypes.ProviderTicker{
				okxbtcusd,
				okxethusd,
			},
			okx.Name,
			providertypes.WebSockets,
			false,
			okxState,
		)

		// Ensure that the provider that is not supported by the marketmap is not enabled.
		binanceState, ok := state[binance.Name]
		require.True(t, ok)
		checkProviderState(
			t,
			nil,
			binance.Name,
			providertypes.API,
			false,
			binanceState,
		)
	})

	t.Run("errors when the API query handler factory is not set", func(t *testing.T) {
		orc, err := oracle.New(
			oracleCfg,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithMarketMap(marketMap),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		o := orc.(*oracle.OracleImpl)

		err = o.Init(context.TODO())
		require.Error(t, err)
	})

	t.Run("errors when the WebSocket query handler factory is not set", func(t *testing.T) {
		orc, err := oracle.New(
			oracleCfg,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithMarketMap(marketMap),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
		)
		require.NoError(t, err)
		o := orc.(*oracle.OracleImpl)

		require.Error(t, o.Init(context.TODO()))
	})

	t.Run("errors when a provider is not supported by the api query handler factory", func(t *testing.T) {
		cfg := copyConfig(oracleCfg)

		cfg.Providers["unsupported"] = config.ProviderConfig{
			Name: "unsupported",
			API: config.APIConfig{
				Enabled:          true,
				Timeout:          5,
				Interval:         5,
				MaxQueries:       5,
				ReconnectTimeout: 5 * time.Second,
				Endpoints:        []config.Endpoint{{URL: "http://test.com"}},
				Name:             "unsupported",
			},
			Type: oracletypes.ConfigType,
		}

		orc, err := oracle.New(
			cfg,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithMarketMap(marketMap),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		o := orc.(*oracle.OracleImpl)

		require.Error(t, o.Init(context.TODO()))
	})

	t.Run("errors when a provider is not supported by the web socket query handler factory", func(t *testing.T) {
		cfg := copyConfig(oracleCfg)

		okxCfg := okx.DefaultWebSocketConfig
		okxCfg.Name = "unsupported"
		cfg.Providers["unsupported"] = config.ProviderConfig{
			Name:      "unsupported",
			WebSocket: okxCfg,
			Type:      oracletypes.ConfigType,
		}

		orc, err := oracle.New(
			cfg,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithMarketMap(marketMap),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		o := orc.(*oracle.OracleImpl)

		require.Error(t, o.Init(context.TODO()))
	})

	t.Run("creates a marketmap provider with price providers", func(t *testing.T) {
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

		err = o.Init(context.TODO())
		require.NoError(t, err)

		mapper := o.GetMarketMapProvider()
		require.NotNil(t, mapper)
	})

	t.Run("errors when the market map factory is not set", func(t *testing.T) {
		orc, err := oracle.New(
			oracleCfgWithMapper,
			noOpPriceAggregator{},
			oracle.WithLogger(logger),
			oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
			oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		)
		require.NoError(t, err)
		o := orc.(*oracle.OracleImpl)

		require.Error(t, o.Init(context.TODO()))
	})
}
