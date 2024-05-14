package ethmulticlient

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/api/metrics"

	"github.com/ethereum/go-ethereum/rpc"
)

// EVMClient is an interface that abstracts the evm client.
//
//go:generate mockery --name EVMClient
type EVMClient interface {
	// BatchCallContext is a batch call to an EVM.
	BatchCallContext(ctx context.Context, calls []rpc.BatchElem) error
}

var _ EVMClient = (*GoEthereumClientImpl)(nil)

// GoEthereumClientImpl is a go-ethereum client implementation using the go-ethereum RPC
// library.
type GoEthereumClientImpl struct {
	apiMetrics metrics.APIMetrics
	api        config.APIConfig

	// client is the underlying rpc client.
	client *rpc.Client
}

// NewGoEthereumClientImplFromURL returns a new go-ethereum client. This is the default
// implementation that connects to an ethereum node via rpc.
func NewGoEthereumClientImplFromURL(
	ctx context.Context,
	apiMetrics metrics.APIMetrics,
	api config.APIConfig,
) (EVMClient, error) {
	if apiMetrics == nil {
		return nil, fmt.Errorf("metrics cannot be nil")
	}

	client, err := rpc.DialOptions(ctx, api.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to dial go ethereum client: %w", err)
	}

	return &GoEthereumClientImpl{
		apiMetrics: apiMetrics,
		api:        api,
		client:     client,
	}, nil
}

// NewGoEthereumClientImplFromEndpoint creates an EVMClient via a config.Endpoint. This includes optional
// authentication via a specified http header key and value.
func NewGoEthereumClientImplFromEndpoint(
	ctx context.Context,
	apiMetrics metrics.APIMetrics,
	api config.APIConfig,
) (EVMClient, error) {
	if apiMetrics == nil {
		return nil, fmt.Errorf("apiMetrics cannot be nil")
	}

	if err := api.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid api config: %w", err)
	}

	if len(api.Endpoints) != 1 {
		return nil, fmt.Errorf("expected exactly one endpoint, got %d", len(api.Endpoints))
	}

	var opts []rpc.ClientOption
	endpoint := api.Endpoints[0]
	if endpoint.Authentication.Enabled() {
		opts = append(opts, rpc.WithHTTPAuth(func(h http.Header) error {
			h.Set(endpoint.Authentication.APIKeyHeader, endpoint.Authentication.APIKey)
			return nil
		}))
	}

	client, err := rpc.DialOptions(ctx, endpoint.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial go ethereum client: %w", err)
	}

	return &GoEthereumClientImpl{
		apiMetrics: apiMetrics,
		api:        api,
		client:     client,
	}, nil
}

// BatchCallContext sends all given requests as a single batch and waits for the server
// to return a response for all of them. The wait duration is bounded by the context's deadline.
//
// In contrast to CallContext, BatchCallContext only returns errors that have occurred while
// sending the request. Any error specific to a request is reported through the Error field of
// the corresponding BatchElem.
//
// Note that batch calls may not be executed atomically on the server side.
func (c *GoEthereumClientImpl) BatchCallContext(ctx context.Context, calls []rpc.BatchElem) error {
	start := time.Now()
	defer func() {
		c.apiMetrics.ObserveProviderResponseLatency(c.api.Name, time.Since(start))
	}()

	if err := c.client.BatchCallContext(ctx, calls); err != nil {
		c.apiMetrics.AddRPCStatusCode(c.api.Name, metrics.RPCCodeError)
		return fmt.Errorf("failed to batch call: %w", err)
	}

	c.apiMetrics.AddRPCStatusCode(c.api.Name, metrics.RPCCodeOK)
	return nil
}
