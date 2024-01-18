package oracle

import (
	"math/big"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/metrics"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// OracleOption is a function that can be used to configure an Oracle.
type OracleOption func(*Oracle) //nolint

// WithLogger sets the logger on the Oracle.
func WithLogger(logger *zap.Logger) OracleOption {
	return func(o *Oracle) {
		if logger == nil {
			panic("cannot set nil logger")
		}

		o.logger = logger.With(zap.String("process", "oracle"))
	}
}

// WithMetrics sets the metrics on the Oracle.
func WithMetrics(metrics metrics.Metrics) OracleOption {
	return func(o *Oracle) {
		if metrics == nil {
			panic("cannot set nil metrics")
		}

		o.metrics = metrics
	}
}

// WithMetricsConfig sets the metrics on the oracle from the given config.
func WithMetricsConfig(config config.OracleMetricsConfig) OracleOption {
	return func(o *Oracle) {
		o.metrics = metrics.NewMetricsFromConfig(config)
	}
}

// WithAggregateFunction sets the aggregate function on the Oracle.
func WithAggregateFunction(fn aggregator.AggregateFn[string, map[oracletypes.CurrencyPair]*big.Int]) OracleOption {
	return func(o *Oracle) {
		if fn == nil {
			panic("cannot set aggregate function on nil aggregator")
		}

		o.priceAggregator = aggregator.NewDataAggregator[string, map[oracletypes.CurrencyPair]*big.Int](
			aggregator.WithAggregateFn(fn),
		)
	}
}

// WithDataAggregator sets the data aggregator on the Oracle.
func WithDataAggregator(agg *aggregator.DataAggregator[string, map[oracletypes.CurrencyPair]*big.Int]) OracleOption {
	return func(o *Oracle) {
		if agg == nil {
			panic("cannot set nil aggregator")
		}

		o.priceAggregator = agg
	}
}

// WithProviders sets the providers on the Oracle.
func WithProviders(providers []providertypes.Provider[oracletypes.CurrencyPair, *big.Int]) OracleOption {
	return func(o *Oracle) {
		o.providers = providers
	}
}
