package providers

import (
	"context"
	"fmt"
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
				OffChainTicker: "USDC-USDT",
				Invert:         true,
			},
		},
	}

	marketMap = mmtypes.MarketMap{
		Markets: map[string]mmtypes.Market{
			usdtusd.Ticker.String(): usdtusd,
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

func NewTestingOracle(ctx context.Context) (TestingOracle, error) {
	cfg := log.NewDefaultConfig()
	cfg.StdOutLogLevel = "debug"
	logger := log.NewLogger(cfg)

	agg, err := oraclemath.NewIndexPriceAggregator(logger, mmtypes.MarketMap{}, oraclemetrics.NewNopMetrics())
	if err != nil {
		return TestingOracle{}, fmt.Errorf("failed to create oracle index price aggregator: %w", err)
	}

	orc, err := oracle.New(
		oracleCfg,
		agg,
		oracle.WithLogger(logger),
		oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
		oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
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
	TestDuration time.Duration
	PollInterval time.Duration
}

func (o *TestingOracle) RunMarketMap(ctx context.Context, mm mmtypes.MarketMap, cfg ProviderTestConfig) error {
	err := o.UpdateMarketMap(mm)
	if err != nil {
		return fmt.Errorf("failed to update oracle market map: %w", err)
	}

	go o.Start(ctx)

	ticker := time.NewTicker(cfg.PollInterval)
	defer ticker.Stop()

	timer := time.NewTicker(cfg.TestDuration)
	defer timer.Stop()

	for {
		select {
		case <-ticker.C:
			// do something
			hm := o.GetPrices()
			o.logger.Info("provider prices", zap.Any("prices", hm))

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
	p, err := NewTestingOracle(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = p.RunMarketMap(ctx, marketMap, ProviderTestConfig{
		TestDuration: 20 * time.Second,
		PollInterval: 1 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
}
