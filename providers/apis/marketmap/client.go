package marketmap

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/skip-mev/connect/v2/oracle/config"
	connectgrpc "github.com/skip-mev/connect/v2/pkg/grpc"
	"github.com/skip-mev/connect/v2/providers/base/api/metrics"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

// MarketMapClient is a wrapper around the x/marketmap QueryClient.
type MarketMapClient struct { //nolint
	mmtypes.QueryClient

	// apiMetrics is the metrics collector for the MarketMapClient.
	apiMetrics metrics.APIMetrics
	// api is the APIConfig for the MarketMapClient.
	api config.APIConfig
}

// NewGRPCClient returns a new GRPC client for MarketMap module.
func NewGRPCClient(
	api config.APIConfig,
	apiMetrics metrics.APIMetrics,
) (mmtypes.QueryClient, error) {
	if api.Name != Name {
		return nil, fmt.Errorf("invalid api name; expected %s, got %s", Name, api.Name)
	}

	// TODO: Do we want to ignore proxy settings?
	conn, err := connectgrpc.NewClient(
		api.Endpoints[0].URL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithNoProxy(),
	)
	if err != nil {
		return nil, err
	}

	return NewGRPCClientWithConn(
		conn,
		api,
		apiMetrics,
	)
}

// NewGRPCClientWithConn returns a new GRPC client for MarketMap module with the given connection.
func NewGRPCClientWithConn(
	conn *grpc.ClientConn,
	api config.APIConfig,
	apiMetrics metrics.APIMetrics,
) (mmtypes.QueryClient, error) {
	if conn == nil {
		return nil, fmt.Errorf("connection is required but got nil")
	}
	if err := api.ValidateBasic(); err != nil {
		return nil, err
	}

	if !api.Enabled {
		return nil, fmt.Errorf("api is not enabled")
	}

	if apiMetrics == nil {
		return nil, fmt.Errorf("metrics is required")
	}

	return &MarketMapClient{
		QueryClient: mmtypes.NewQueryClient(conn),
		apiMetrics:  apiMetrics,
		api:         api,
	}, nil
}

// MarketMap wraps the MarketMapClient's query with additional metrics.
func (c *MarketMapClient) MarketMap(
	ctx context.Context,
	req *mmtypes.MarketMapRequest,
	_ ...grpc.CallOption,
) (resp *mmtypes.MarketMapResponse, err error) {
	start := time.Now()
	defer func() {
		c.apiMetrics.ObserveProviderResponseLatency(c.api.Name, metrics.RedactedURL, time.Since(start))
	}()

	resp, err = c.QueryClient.MarketMap(ctx, req)
	if err != nil {
		c.apiMetrics.AddRPCStatusCode(c.api.Name, metrics.RedactedURL, metrics.RPCCodeError)
		return
	}

	c.apiMetrics.AddRPCStatusCode(c.api.Name, metrics.RedactedURL, metrics.RPCCodeOK)
	return
}
