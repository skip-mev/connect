package oracle_test

import (
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle"
	"github.com/skip-mev/connect/v2/oracle/config"
	oracletypes "github.com/skip-mev/connect/v2/oracle/types"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/providers/apis/binance"
	"github.com/skip-mev/connect/v2/providers/apis/coinbase"
	"github.com/skip-mev/connect/v2/providers/apis/dydx"
	"github.com/skip-mev/connect/v2/providers/base"
	"github.com/skip-mev/connect/v2/providers/base/api/handlers/mocks"
	apimetrics "github.com/skip-mev/connect/v2/providers/base/api/metrics"
	providermetrics "github.com/skip-mev/connect/v2/providers/base/metrics"
	"github.com/skip-mev/connect/v2/providers/static"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
	"github.com/skip-mev/connect/v2/providers/websockets/okx"
	mmclienttypes "github.com/skip-mev/connect/v2/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

var (
	btcusdtCP = connecttypes.NewCurrencyPair("BTC", "USDT")
	btcusdCP  = connecttypes.NewCurrencyPair("BTC", "USD")
	usdtusdCP = connecttypes.NewCurrencyPair("USDT", "USD")
	ethusdtCP = connecttypes.NewCurrencyPair("ETH", "USDT")
)

var (
	logger = zap.NewExample()

	oracleCfg = config.OracleConfig{
		Metrics: config.MetricsConfig{
			Enabled: false,
			Telemetry: config.TelemetryConfig{
				Disabled: true,
			},
		},
		UpdateInterval: 1500 * time.Millisecond,
		MaxPriceAge:    2 * time.Minute,
		Providers: map[string]config.ProviderConfig{
			binance.Name: { // Price API provider.
				Name: binance.Name,
				API:  binance.DefaultNonUSAPIConfig,
				Type: oracletypes.ConfigType,
			},
			coinbase.Name: { // Price API provider.
				Name: coinbase.Name,
				API:  coinbase.DefaultAPIConfig,
				Type: oracletypes.ConfigType,
			},
			okx.Name: { // Price WebSocket provider.
				Name:      okx.Name,
				WebSocket: okx.DefaultWebSocketConfig,
				Type:      oracletypes.ConfigType,
			},
		},
		Host: "localhost",
		Port: "8080",
	}

	oracleCfgWithMapper = config.OracleConfig{
		Metrics: config.MetricsConfig{
			Enabled: false,
			Telemetry: config.TelemetryConfig{
				Disabled: true,
			},
		},
		UpdateInterval: 1500 * time.Millisecond,
		MaxPriceAge:    2 * time.Minute,
		Providers: map[string]config.ProviderConfig{
			binance.Name: { // Price API provider.
				Name: binance.Name,
				API:  binance.DefaultNonUSAPIConfig,
				Type: oracletypes.ConfigType,
			},
			coinbase.Name: { // Price API provider.
				Name: coinbase.Name,
				API:  coinbase.DefaultAPIConfig,
				Type: oracletypes.ConfigType,
			},
			okx.Name: { // Price WebSocket provider.
				Name:      okx.Name,
				WebSocket: okx.DefaultWebSocketConfig,
				Type:      oracletypes.ConfigType,
			},
			// Market map provider.
			mapperCfg.Name: mapperCfg,
		},
		Host: "localhost",
		Port: "8080",
	}

	oracleCfgWithMockMapper = config.OracleConfig{
		Metrics: config.MetricsConfig{
			Enabled: false,
			Telemetry: config.TelemetryConfig{
				Disabled: true,
			},
		},
		UpdateInterval: 1500 * time.Millisecond,
		MaxPriceAge:    2 * time.Minute,
		Providers: map[string]config.ProviderConfig{
			binance.Name: { // Price API provider.
				Name: binance.Name,
				API:  binance.DefaultNonUSAPIConfig,
				Type: oracletypes.ConfigType,
			},
			coinbase.Name: { // Price API provider.
				Name: coinbase.Name,
				API:  coinbase.DefaultAPIConfig,
				Type: oracletypes.ConfigType,
			},
			okx.Name: { // Price WebSocket provider.
				Name:      okx.Name,
				WebSocket: okx.DefaultWebSocketConfig,
				Type:      oracletypes.ConfigType,
			},
			// Market map provider.
			mockMapperCfg.Name: mockMapperCfg,
		},
		Host: "localhost",
		Port: "8080",
	}

	oracleCfgWithOnlyMockMapper = config.OracleConfig{
		Metrics: config.MetricsConfig{
			Enabled: false,
			Telemetry: config.TelemetryConfig{
				Disabled: true,
			},
		},
		UpdateInterval: 1500 * time.Millisecond,
		MaxPriceAge:    2 * time.Minute,
		Providers: map[string]config.ProviderConfig{
			// Market map provider.
			mockMapperCfg.Name: mockMapperCfg,
		},
		Host: "localhost",
		Port: "8080",
	}

	mapperCfg = config.ProviderConfig{
		Name: dydx.Name,
		API:  dydx.DefaultAPIConfig,
		Type: mmclienttypes.ConfigType,
	}

	mockMapperCfg = config.ProviderConfig{
		Name: "mock-mapper",
		API: config.APIConfig{
			Enabled:          true,
			Timeout:          5 * time.Second,
			Interval:         1000 * time.Millisecond,
			ReconnectTimeout: 1000 * time.Millisecond,
			MaxQueries:       1,
			Atomic:           true,
			Endpoints:        []config.Endpoint{{URL: "http://test.com"}},
			Name:             "mock-mapper",
		},
		Type: mmclienttypes.ConfigType,
	}

	validMarketMapSubset = mmtypes.MarketMap{
		Markets: map[string]mmtypes.Market{
			ethusdtCP.String(): {
				Ticker: mmtypes.Ticker{
					CurrencyPair:     ethusdtCP,
					MinProviderCount: 1,
					Decimals:         8,
					Enabled:          true,
				},
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:           coinbase.Name,
						OffChainTicker: coinbaseethusd.GetOffChainTicker(),
					},
					{
						Name:           okx.Name,
						OffChainTicker: okxethusd.GetOffChainTicker(),
					},
				},
			},
		},
	}

	partialInvalidMarketMap = mmtypes.MarketMap{
		Markets: map[string]mmtypes.Market{
			btcusdCP.String(): {
				Ticker: mmtypes.Ticker{
					CurrencyPair:     btcusdCP,
					MinProviderCount: 1,
					Decimals:         8,
					Enabled:          true,
				},
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:            coinbase.Name,
						OffChainTicker:  coinbasebtcusd.GetOffChainTicker(),
						NormalizeByPair: &usdtusdCP,
					},
					{
						Name:            okx.Name,
						OffChainTicker:  okxbtcusd.GetOffChainTicker(),
						NormalizeByPair: &usdtusdCP,
					},
				},
			},
			ethusdtCP.String(): {
				Ticker: mmtypes.Ticker{
					CurrencyPair:     ethusdtCP,
					MinProviderCount: 1,
					Decimals:         8,
					Enabled:          true,
				},
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:           coinbase.Name,
						OffChainTicker: coinbaseethusd.GetOffChainTicker(),
					},
					{
						Name:           okx.Name,
						OffChainTicker: okxethusd.GetOffChainTicker(),
					},
				},
			},
		},
	}

	// Coinbase and OKX are supported by the marketmap.
	marketMap = mmtypes.MarketMap{
		Markets: map[string]mmtypes.Market{
			btcusdtCP.String(): {
				Ticker: mmtypes.Ticker{
					CurrencyPair:     btcusdtCP,
					MinProviderCount: 1,
					Decimals:         8,
					Enabled:          true,
				},
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:           coinbase.Name,
						OffChainTicker: coinbasebtcusd.GetOffChainTicker(),
					},
					{
						Name:           okx.Name,
						OffChainTicker: okxbtcusd.GetOffChainTicker(),
					},
				},
			},
			ethusdtCP.String(): {
				Ticker: mmtypes.Ticker{
					CurrencyPair:     ethusdtCP,
					MinProviderCount: 1,
					Decimals:         8,
					Enabled:          true,
				},
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:           coinbase.Name,
						OffChainTicker: coinbaseethusd.GetOffChainTicker(),
					},
					{
						Name:           okx.Name,
						OffChainTicker: okxethusd.GetOffChainTicker(),
					},
				},
			},
		},
	}
)
var _ oracle.PriceAggregator = &noOpPriceAggregator{}

type noOpPriceAggregator struct{}

func (n noOpPriceAggregator) SetProviderPrices(_ string, _ oracletypes.Prices) {
}

func (n noOpPriceAggregator) UpdateMarketMap(_ mmtypes.MarketMap) {
}

func (n noOpPriceAggregator) AggregatePrices() {
}

func (n noOpPriceAggregator) GetPrices() oracletypes.Prices {
	return oracletypes.Prices{}
}

func (n noOpPriceAggregator) Reset() {
}

func checkProviderState(
	t *testing.T,
	expectedTickers []oracletypes.ProviderTicker,
	expectedName string,
	expectedType providertypes.ProviderType,
	isRunning bool,
	state oracle.ProviderState,
) {
	t.Helper()

	// Ensure that the provider is enabled and supports the expected tickers.
	provider := state.Provider
	require.Equal(t, expectedName, provider.Name())
	require.Equal(t, expectedType, provider.Type())

	ids := provider.GetIDs()
	require.Equal(t, len(expectedTickers), len(ids))
	seenTickers := make(map[oracletypes.ProviderTicker]bool)
	for _, id := range ids {
		seenTickers[id] = true
	}
	for _, ticker := range expectedTickers {
		require.True(t, seenTickers[ticker])
	}

	// Check the market map.
	require.Equal(t, len(expectedTickers), len(provider.GetIDs()))
	for _, ticker := range provider.GetIDs() {
		require.True(t, seenTickers[ticker])
	}

	// Ensure that the provider is running/no-running.
	require.Equal(t, isRunning, provider.IsRunning())
}

func createTestMarketMapProvider(
	t *testing.T,
	ids []mmclienttypes.Chain,
) (*mocks.APIDataHandler[mmclienttypes.Chain, *mmtypes.MarketMapResponse], *mmclienttypes.MarketMapProvider) {
	t.Helper()

	// Create a market map api handler.
	handler := mocks.NewAPIDataHandler[mmclienttypes.Chain, *mmtypes.MarketMapResponse](
		t,
	)

	queryHandler, err := mmclienttypes.NewMarketMapAPIQueryHandler(
		logger,
		mockMapperCfg.API,
		static.NewStaticMockClient(),
		handler,
		apimetrics.NewNopAPIMetrics(),
	)
	require.NoError(t, err)

	var provider *mmclienttypes.MarketMapProvider
	if len(ids) != 0 {
		provider, err = mmclienttypes.NewMarketMapProvider(
			base.WithName[mmclienttypes.Chain, *mmtypes.MarketMapResponse](mockMapperCfg.Name),
			base.WithLogger[mmclienttypes.Chain, *mmtypes.MarketMapResponse](logger),
			base.WithAPIQueryHandler(queryHandler),
			base.WithAPIConfig[mmclienttypes.Chain, *mmtypes.MarketMapResponse](mockMapperCfg.API),
			base.WithMetrics[mmclienttypes.Chain, *mmtypes.MarketMapResponse](providermetrics.NewNopProviderMetrics()),
			base.WithIDs[mmclienttypes.Chain, *mmtypes.MarketMapResponse](ids),
		)
	} else {
		provider, err = mmclienttypes.NewMarketMapProvider(
			base.WithName[mmclienttypes.Chain, *mmtypes.MarketMapResponse](mockMapperCfg.Name),
			base.WithLogger[mmclienttypes.Chain, *mmtypes.MarketMapResponse](logger),
			base.WithAPIQueryHandler(queryHandler),
			base.WithAPIConfig[mmclienttypes.Chain, *mmtypes.MarketMapResponse](mockMapperCfg.API),
			base.WithMetrics[mmclienttypes.Chain, *mmtypes.MarketMapResponse](providermetrics.NewNopProviderMetrics()),
		)
	}
	require.NoError(t, err)

	return handler, provider
}

func marketMapperFactory(
	t *testing.T,
	ids []mmclienttypes.Chain,
) (*mocks.APIDataHandler[mmclienttypes.Chain, *mmtypes.MarketMapResponse], mmclienttypes.MarketMapFactory) {
	t.Helper()

	handler, provider := createTestMarketMapProvider(
		t,
		ids,
	)

	return handler, func(
		*zap.Logger,
		providermetrics.ProviderMetrics,
		apimetrics.APIMetrics,
		config.ProviderConfig,
	) (*mmclienttypes.MarketMapProvider, error) {
		return provider, nil
	}
}

func copyConfig(cfg config.OracleConfig) config.OracleConfig {
	// copy providers map
	newCfg := cfg

	newCfg.Providers = make(map[string]config.ProviderConfig)
	for k, v := range cfg.Providers {
		newCfg.Providers[k] = v
	}

	return newCfg
}
