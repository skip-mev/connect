package client

import (
	"context"
	"time"

	"cosmossdk.io/log"
	oracle_metrics "github.com/skip-mev/slinky/oracle/metrics"
	"github.com/skip-mev/slinky/service"
	"github.com/skip-mev/slinky/service/metrics"
)

var _ service.OracleService = (*MetricsClient)(nil)

// MetricsClient is a wrapper around a Client implementation that exports system-level metrics about the client.
type MetricsClient struct {
	logger log.Logger
	service.OracleService
	metrics metrics.Metrics
}

func NewMetricsClient(logger log.Logger, client service.OracleService, metrics metrics.Metrics) *MetricsClient {
	return &MetricsClient{
		logger:        logger,
		OracleService: client,
		metrics:       metrics,
	}
}

func (m *MetricsClient) Prices(ctx context.Context, req *service.QueryPricesRequest) (res *service.QueryPricesResponse, err error) {
	// measure the beginning of call
	start := time.Now()

	defer func() {
		// observe the duration of call
		m.metrics.ObserveOracleResponseLatency(time.Since(start))
		// observe error
		m.metrics.AddOracleResponse(oracle_metrics.StatusFromError(err))
	}()

	// call the underlying client
	res, err = m.OracleService.Prices(ctx, req)

	m.logger.Debug("calling underlying client", "res", res, "err", err, "duration", time.Since(start))
	return
}
