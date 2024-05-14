package marketmap

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// MarketMapClient is a wrapper around the x/marketmap QueryClient.
type MarketMapClient struct { //nolint
	mmtypes.QueryClient

	// metrics is the metrics collector for the MarketMapClient.
	metrics metrics.APIMetrics
	// api is the APIConfig for the MarketMapClient.
	api config.APIConfig
}

// NewGRPCClient returns a new GRPC client for MarketMap module.
func NewMarketMapClient(
	api config.APIConfig,
	metrics metrics.APIMetrics,
) (mmtypes.QueryClient, error) {
	if err := api.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid api config: %w", err)
	}

	if metrics == nil {
		return nil, fmt.Errorf("metrics is nil")
	}

	// TODO: Do we want to ignore proxy settings?
	conn, err := grpc.Dial(
		api.URL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &MarketMapClient{
		QueryClient: mmtypes.NewQueryClient(conn),
		metrics:     metrics,
		api:         api,
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
		c.metrics.ObserveProviderResponseLatency(c.api.Name, time.Since(start))
	}()

	resp, err := c.QueryClient.MarketMap(ctx, req)
	if err != nil {
		c.metrics.AddRPCStatusCode(c.api.Name, metrics.RPCCodeOK)
		return resp, err
	}

	c.metrics.AddRPCStatusCode(c.api.Name, metrics.RPCCodeOK)
	return resp, nil
}
