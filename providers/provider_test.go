package providers

import (
	"context"
	"fmt"
	cmdconfig "github.com/skip-mev/connect/v2/cmd/connect/config"
	"github.com/skip-mev/connect/v2/cmd/constants"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/pkg/log"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"

	"github.com/skip-mev/connect/v2/oracle"
	"github.com/skip-mev/connect/v2/oracle/config"
	oraclemetrics "github.com/skip-mev/connect/v2/oracle/metrics"
	oracletypes "github.com/skip-mev/connect/v2/oracle/types"
	oraclemath "github.com/skip-mev/connect/v2/pkg/math/oracle"
	"github.com/skip-mev/connect/v2/providers/apis/binance"
	"github.com/skip-mev/connect/v2/providers/apis/coinbase"
	oraclefactory "github.com/skip-mev/connect/v2/providers/factories/oracle"
	"github.com/skip-mev/connect/v2/providers/websockets/okx"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

var oracleCfg = config.OracleConfig{
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

var (
	usdtusd = mmtypes.Market{
		Ticker: mmtypes.Ticker{
			CurrencyPair: connecttypes.CurrencyPair{
				Base:  "USDT",
				Quote: "USD",
			},
			Decimals:         8,
			MinProviderCount: 1,
			Enabled:          true,
		},
		ProviderConfigs: []mmtypes.ProviderConfig{
			{
				Name:           "okx_ws",
				OffChainTicker: "USDC{-USDT",
				Invert:         true,
			},
		},
	}
)

type TestingOracle struct {
	Oracle *oracle.OracleImpl
	logger *zap.Logger
}

func (o *TestingOracle) Start(ctx context.Context) error {
	return o.Oracle.Start(ctx)
}

func (o *TestingOracle) Stop() {
	o.Oracle.Stop()
}

func (o *TestingOracle) GetPrices() oracletypes.Prices {
	return o.Oracle.GetPrices()
}

func (o *TestingOracle) UpdateMarketMap(mm mmtypes.MarketMap) error {
	return o.Oracle.UpdateMarketMap(mm)
}

func FilterMarketMapToProviders(mm mmtypes.MarketMap) map[string]mmtypes.MarketMap {
	m := make(map[string]mmtypes.MarketMap)

	for _, market := range mm.Markets {
		// check each provider config
		for _, pc := range market.ProviderConfigs {
			// create a market from the given provider config
			isolatedMarket := mmtypes.Market{
				Ticker: market.Ticker,
				ProviderConfigs: []mmtypes.ProviderConfig{
					pc,
				},
			}

			// init mm if necessary
			if _, found := m[pc.Name]; !found {
				m[pc.Name] = mmtypes.MarketMap{
					Markets: map[string]mmtypes.Market{
						isolatedMarket.Ticker.String(): isolatedMarket,
					},
				}
				// otherwise insert
			} else {
				m[pc.Name].Markets[isolatedMarket.Ticker.String()] = isolatedMarket
			}
		}
	}

	return m
}

func OracleConfigForProvider(providerNames ...string) (config.OracleConfig, error) {
	cfg := config.OracleConfig{
		UpdateInterval: cmdconfig.DefaultUpdateInterval,
		MaxPriceAge:    cmdconfig.DefaultMaxPriceAge,
		Metrics: config.MetricsConfig{
			Enabled: false,
			Telemetry: config.TelemetryConfig{
				Disabled: true,
			},
		},
		Providers: make(map[string]config.ProviderConfig),
		Host:      cmdconfig.DefaultHost,
		Port:      cmdconfig.DefaultPort,
	}

	for _, provider := range append(constants.Providers, constants.AlternativeMarketMapProviders...) {
		for _, providerName := range providerNames {
			if provider.Name == providerName {
				cfg.Providers[provider.Name] = provider
			}
		}
	}

	if err := cfg.ValidateBasic(); err != nil {

		return cfg, fmt.Errorf("default oracle config is invalid: %s", err)
	}

	return cfg, nil
}

func NewTestingOracle(ctx context.Context, providerNames ...string) (TestingOracle, error) {
	logCfg := log.NewDefaultConfig()
	logCfg.StdOutLogLevel = "debug"
	logger := log.NewLogger(logCfg)

	agg, err := oraclemath.NewIndexPriceAggregator(logger, mmtypes.MarketMap{}, oraclemetrics.NewNopMetrics())
	if err != nil {
		return TestingOracle{}, fmt.Errorf("failed to create oracle index price aggregator: %w", err)
	}

	cfg, err := OracleConfigForProvider(providerNames...)
	if err != nil {
		return TestingOracle{}, fmt.Errorf("failed to create oracle config: %w", err)
	}

	orc, err := oracle.New(
		cfg,
		agg,
		oracle.WithLogger(logger),
		oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
		oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		oracle.WithMarketMapperFactory(oraclefactory.MarketMapProviderFactory),
	)
	if err != nil {
		return TestingOracle{}, fmt.Errorf("failed to create oracle: %w", err)
	}

	o := orc.(*oracle.OracleImpl)
	err = o.Init(ctx)
	if err != nil {
		return TestingOracle{}, err
	}

	return TestingOracle{
		Oracle: o,
		logger: logger,
	}, nil
}

type ProviderTestConfig struct {
	TestDuration   time.Duration
	PollInterval   time.Duration
	BurnInInterval time.Duration
}

func (o *TestingOracle) RunMarketMap(ctx context.Context, mm mmtypes.MarketMap, cfg ProviderTestConfig) error {
	err := o.UpdateMarketMap(mm)
	if err != nil {
		return fmt.Errorf("failed to update oracle market map: %w", err)
	}

	expectedNumPrices := len(mm.Markets)
	if expectedNumPrices == 0 {
		return fmt.Errorf("cannot test with empty market map")
	}

	go o.Start(ctx)
	time.Sleep(cfg.BurnInInterval)

	ticker := time.NewTicker(cfg.PollInterval)
	defer ticker.Stop()

	timer := time.NewTicker(cfg.TestDuration)
	defer timer.Stop()

	for {
		select {
		case <-ticker.C:
			prices := o.GetPrices()
			if len(prices) != expectedNumPrices {
				return fmt.Errorf("expected %d prices, got %d", expectedNumPrices, len(prices))
			}
			o.logger.Info("provider prices", zap.Any("prices", prices))

		case <-timer.C:
			o.Stop()

			// cleanup
			return nil
		}
	}
}

func (o *TestingOracle) RunMarket(ctx context.Context, market mmtypes.Market, cfg ProviderTestConfig) error {
	mm := mmtypes.MarketMap{
		Markets: map[string]mmtypes.Market{
			market.Ticker.String(): market,
		},
	}

	return o.RunMarketMap(ctx, mm, cfg)
}

func TestProvider(t *testing.T) {
	ctx := context.Background()
	p, err := NewTestingOracle(ctx, okx.Name)
	if err != nil {
		t.Fatal(err)
	}

	err = p.RunMarket(ctx, usdtusd, ProviderTestConfig{
		TestDuration:   20 * time.Second,
		PollInterval:   1 * time.Second,
		BurnInInterval: 2 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
}
