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
	if err := api.ValidateBasic(); err != nil {
		return nil, err
	}

	if api.Name != Name {
		return nil, fmt.Errorf("invalid api name; expected %s, got %s", Name, api.Name)
	}

	if !api.Enabled {
		return nil, fmt.Errorf("api is not enabled")
	}

	if apiMetrics == nil {
		return nil, fmt.Errorf("metrics is required")
	}

	// TODO: Do we want to ignore proxy settings?
	conn, err := grpc.NewClient(
		api.URL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
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
