package marketmap

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// MarketMapClient is a wrapper around the x/marketmap QueryClient.
type MarketMapClient struct {
	mmtypes.QueryClient

	metrics metrics.APIMetrics
	config  config.APIConfig
}

// NewGRPCClient returns a new GRPC client for MarketMap module.
func NewMarketMapClient(
	config config.APIConfig,
	metrics metrics.APIMetrics,
) (mmtypes.QueryClient, error) {
	// TODO: Do we want to ignore proxy settings?
	conn, err := grpc.Dial(
		config.URL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &MarketMapClient{
		QueryClient: mmtypes.NewQueryClient(conn),
		metrics:     metrics,
		config:      config,
	}, nil
}

// MarketMap wraps the MarketMapClient's query with additional metrics.
func (c *MarketMapClient) MarketMap(
	ctx context.Context,
	req *mmtypes.MarketMapRequest,
	_ ...grpc.CallOption,
) (*mmtypes.MarketMapResponse, error) {
	start := time.Now()
	defer func() {
		c.metrics.ObserveProviderResponseLatency(c.config.Name, c.config.URL, time.Since(start))
	}()

	resp, err := c.QueryClient.MarketMap(ctx, req)
	if err != nil {
		c.metrics.AddRPCStatusCode(c.config.Name, c.config.URL, metrics.RPCCodeOK)
		return resp, err
	}

	c.metrics.AddRPCStatusCode(c.config.Name, c.config.URL, metrics.RPCCodeOK)
	return resp, nil
}
