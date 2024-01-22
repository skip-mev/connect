package oracle

import (
	"context"
	"time"

	"cosmossdk.io/log"
	"google.golang.org/grpc"

	"github.com/skip-mev/slinky/service/metrics"
	"github.com/skip-mev/slinky/service/servers/oracle/types"
)

var _ OracleClient = (*MetricsClient)(nil)

// MetricsClient is a wrapper around a Client implementation that exports system-level metrics about the client.
type MetricsClient struct {
	OracleClient

	logger  log.Logger
	metrics metrics.Metrics
}

func NewMetricsClient(logger log.Logger, client OracleClient, metrics metrics.Metrics) *MetricsClient {
	return &MetricsClient{
		logger:       logger,
		OracleClient: client,
		metrics:      metrics,
	}
}

func (m *MetricsClient) Prices(ctx context.Context, req *types.QueryPricesRequest, _ ...grpc.CallOption) (res *types.QueryPricesResponse, err error) {
	// measure the beginning of call
	start := time.Now()

	defer func() {
		// observe the duration of call
		m.metrics.ObserveOracleResponseLatency(time.Since(start))
		// observe error
		m.metrics.AddOracleResponse(metrics.StatusFromError(err))
	}()

	// call the underlying client
	res, err = m.OracleClient.Prices(ctx, req)

	m.logger.Debug("calling underlying client", "res", res, "err", err, "duration", time.Since(start))
	return
}
