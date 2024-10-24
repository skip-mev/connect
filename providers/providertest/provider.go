package providertest

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle"
	oraclemetrics "github.com/skip-mev/slinky/oracle/metrics"
	oracletypes "github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/log"
	oraclemath "github.com/skip-mev/slinky/pkg/math/oracle"
	oraclefactory "github.com/skip-mev/slinky/providers/factories/oracle"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

type TestingOracle struct {
	Oracle *oracle.OracleImpl
	Logger *zap.Logger
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

func NewTestingOracle(ctx context.Context, providerNames []string, extraOpts ...oracle.Option) (TestingOracle, error) {
	logCfg := log.NewDefaultConfig()
	logCfg.StdOutLogLevel = "debug"
	logCfg.FileOutLogLevel = "debug"
	logger := log.NewLogger(logCfg)

	agg, err := oraclemath.NewIndexPriceAggregator(logger, mmtypes.MarketMap{}, oraclemetrics.NewNopMetrics())
	if err != nil {
		return TestingOracle{}, fmt.Errorf("failed to create oracle index price aggregator: %w", err)
	}

	cfg, err := OracleConfigForProvider(providerNames...)
	if err != nil {
		return TestingOracle{}, fmt.Errorf("failed to create oracle config: %w", err)
	}

	opts := []oracle.Option{
		oracle.WithLogger(logger),
		oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
		oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		oracle.WithMarketMapperFactory(oraclefactory.MarketMapProviderFactory),
	}
	opts = append(opts, extraOpts...)

	orc, err := oracle.New(
		cfg,
		agg,
		opts...,
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
		Logger: logger,
	}, nil
}

// Config is used to configure tests when calling RunMarket or RunMarketMap.
type Config struct {
	// TestDuration is the total duration of testing time.
	TestDuration time.Duration
	// PolInterval is the interval at which the provider will be queried for prices.
	PollInterval time.Duration
	// BurnInInterval is the amount of time to allow the provider to run before querying it.
	BurnInInterval time.Duration
}

func (c *Config) Validate() error {
	if c.TestDuration == 0 {
		return fmt.Errorf("test duration cannot be 0")
	}

	if c.PollInterval == 0 {
		return fmt.Errorf("poll interval cannot be 0")
	}

	if c.TestDuration/c.PollInterval < 1 {
		return fmt.Errorf("ratio of test duration to poll interval must be GTE 1")
	}

	return nil
}

// DefaultProviderTestConfig tests by:
// - allow the providers to run for 5 seconds
// - test for a total of 1 minute
// - poll each 5 seconds for prices.
func DefaultProviderTestConfig() Config {
	return Config{
		TestDuration:   1 * time.Minute,
		PollInterval:   5 * time.Second,
		BurnInInterval: 5 * time.Second,
	}
}

// PriceResults is a type alias for an array of PriceResult.
type PriceResults []PriceResult

// PriceResult is a snapshot of Prices results at a given time point when testing.
type PriceResult struct {
	Prices oracletypes.Prices
	Time   time.Time
}

func (o *TestingOracle) RunMarketMap(ctx context.Context, mm mmtypes.MarketMap, cfg Config) (PriceResults, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	err := o.UpdateMarketMap(mm)
	if err != nil {
		return nil, fmt.Errorf("failed to update oracle market map: %w", err)
	}

	expectedNumPrices := len(mm.Markets)
	if expectedNumPrices == 0 {
		return nil, fmt.Errorf("cannot test with empty market map")
	}

	go o.Start(ctx)
	time.Sleep(cfg.BurnInInterval)

	priceResults := make(PriceResults, 0, cfg.TestDuration/cfg.PollInterval)

	ticker := time.NewTicker(cfg.PollInterval)
	defer ticker.Stop()

	timer := time.NewTicker(cfg.TestDuration)
	defer timer.Stop()

	for {
		select {
		case <-ticker.C:
			prices := o.GetPrices()
			if len(prices) != expectedNumPrices {
				return nil, fmt.Errorf("expected %d prices, got %d", expectedNumPrices, len(prices))
			}

			priceResults = append(priceResults, PriceResult{
				Prices: prices,
				Time:   time.Now(),
			})

		case <-timer.C:
			o.Stop()

			// cleanup
			return priceResults, nil
		}
	}
}

func (o *TestingOracle) RunMarket(ctx context.Context, market mmtypes.Market, cfg Config) (PriceResults, error) {
	mm := mmtypes.MarketMap{
		Markets: map[string]mmtypes.Market{
			market.Ticker.String(): market,
		},
	}

	return o.RunMarketMap(ctx, mm, cfg)
}
