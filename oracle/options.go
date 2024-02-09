package oracle

import (
	"math/big"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/metrics"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	providertypes "github.com/skip-mev/slinky/providers/types"
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
func WithAggregateFunction(fn aggregator.AggregateFn[string, map[slinkytypes.CurrencyPair]*big.Int]) Option {
	return func(o *OracleImpl) {
		if fn == nil {
			panic("cannot set aggregate function on nil aggregator")
		}

		o.priceAggregator = aggregator.NewDataAggregator[string, map[slinkytypes.CurrencyPair]*big.Int](
			aggregator.WithAggregateFn(fn),
		)
	}
}

// WithDataAggregator sets the data aggregator on the Oracle.
func WithDataAggregator(agg *aggregator.DataAggregator[string, map[slinkytypes.CurrencyPair]*big.Int]) Option {
	return func(o *OracleImpl) {
		if agg == nil {
			panic("cannot set nil aggregator")
		}

		o.priceAggregator = agg
	}
}

// WithProviders sets the providers on the Oracle.
func WithProviders(providers []providertypes.Provider[slinkytypes.CurrencyPair, *big.Int]) Option {
	return func(o *OracleImpl) {
		o.providers = providers
	}
}
