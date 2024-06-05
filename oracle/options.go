package oracle

import (
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	oraclemetrics "github.com/skip-mev/slinky/oracle/metrics"
	"github.com/skip-mev/slinky/oracle/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// Option is a function that can be used to configure an Oracle.
type Option func(*OracleImpl)

// WithUpdateInterval sets the update interval on the Oracle.
func WithUpdateInterval(updateInterval time.Duration) Option {
	return func(o *OracleImpl) {
		if updateInterval <= 0 {
			panic("update interval must be positive")
		}

		o.updateInterval = updateInterval
	}
}

// WithMaxCacheAge sets the max cache age on the Oracle.
func WithMaxCacheAge(maxCacheAge time.Duration) Option {
	return func(o *OracleImpl) {
		if maxCacheAge <= 0 {
			panic("max cache age must be positive")
		}

		o.maxCacheAge = maxCacheAge
	}
}

// WithLogger sets the logger on the Oracle.
func WithLogger(logger *zap.Logger) Option {
	return func(o *OracleImpl) {
		if logger == nil {
			panic("cannot set nil logger")
		}

		o.logger = logger.With(zap.String("process", "oracle"))
	}
}

// WithMetrics sets the metrics on the Oracle.
func WithMetrics(metrics oraclemetrics.Metrics) Option {
	return func(o *OracleImpl) {
		if metrics == nil {
			panic("cannot set nil metrics")
		}

		o.metrics = metrics
	}
}

// WithMetricsConfig sets the metrics on the oracle from the given config.
func WithMetricsConfig(config config.MetricsConfig) Option {
	return func(o *OracleImpl) {
		o.metrics = oraclemetrics.NewMetricsFromConfig(config)
	}
}

// WithPriceAggregator sets the data aggregator on the Oracle.
func WithPriceAggregator(agg PriceAggregator) Option {
	return func(o *OracleImpl) {
		if agg == nil {
			panic("cannot set nil aggregator")
		}

		o.priceAggregator = agg
	}
}

// WithMarketMapGetter sets a getter function for the latest market map on the Oracle.
func WithMarketMapGetter(fn func() mmtypes.MarketMap) Option {
	return func(o *OracleImpl) {
		o.marketMapGetter = fn
	}
}

// WithProviders sets the providers on the Oracle.
func WithProviders(providers []*types.PriceProvider) Option {
	return func(o *OracleImpl) {
		o.providers = providers
	}
}
