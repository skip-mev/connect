package oracle

import (
	"go.uber.org/zap"

	oraclemetrics "github.com/skip-mev/slinky/oracle/metrics"
)

// Option is a function that can be used to configure an Oracle.
type Option func(*OracleImpl)

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
