package oracle

import (
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/metrics"
	"github.com/skip-mev/slinky/oracle/types"
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
func WithMetrics(metrics metrics.Metrics) Option {
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
		o.metrics = metrics.NewMetricsFromConfig(config)
	}
}

// WithAggregateFunction sets the aggregate function on the Oracle.
func WithAggregateFunction(fn aggregator.AggregateFn[string, types.TickerPrices]) Option {
	return func(o *OracleImpl) {
		if fn == nil {
			panic("cannot set aggregate function on nil aggregator")
		}

		o.priceAggregator = aggregator.NewDataAggregator[string, types.TickerPrices](
			aggregator.WithAggregateFn(fn),
		)
	}
}

// WithDataAggregator sets the data aggregator on the Oracle.
func WithDataAggregator(agg *aggregator.DataAggregator[string, types.TickerPrices]) Option {
	return func(o *OracleImpl) {
		if agg == nil {
			panic("cannot set nil aggregator")
		}

		o.priceAggregator = agg
	}
}

// WithProviders sets the providers on the Oracle.
func WithProviders(providers []types.PriceProvider) Option {
	return func(o *OracleImpl) {
		o.providers = providers
	}
}
