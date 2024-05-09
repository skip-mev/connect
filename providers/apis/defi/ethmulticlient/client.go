package ethmulticlient

import (
	"context"
	"fmt"
	"net/http"

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
	rpcMetrics metrics.APIMetrics
	api        config.APIConfig

	// client is the underlying rpc client.
	client *rpc.Client
}

// NewGoEthereumClientImplFromURL returns a new go-ethereum client. This is the default
// implementation that connects to an ethereum node via rpc.
func NewGoEthereumClientImplFromURL(
	ctx context.Context,
	rpcMetrics metrics.APIMetrics,
	api config.APIConfig,
) (EVMClient, error) {
	client, err := rpc.DialOptions(ctx, api.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to dial go ethereum client: %w", err)
	}

	return &GoEthereumClientImpl{
		rpcMetrics: rpcMetrics,
		api:        api,
		client:     client,
	}, nil
}

// NewGoEthereumClientImplFromEndpoint creates an EVMClient via a config.Endpoint. This includes optional
// authentication via a specified http header key and value.
func NewGoEthereumClientImplFromEndpoint(
	ctx context.Context,
	rpcMetrics metrics.APIMetrics,
	api config.APIConfig,
) (EVMClient, error) {
	var opts []rpc.ClientOption

	if len(api.Endpoints) != 1 {
		return nil, fmt.Errorf("no endpoints provided to create go-ethereum client")
	}

	// fail if we have an invalid endpoint
	endpoint := api.Endpoints[0]
	if err := endpoint.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid endpoint %v: %w", endpoint, err)
	}
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
		rpcMetrics: rpcMetrics,
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
	if err := c.client.BatchCallContext(ctx, calls); err != nil {
		return fmt.Errorf("failed to batch call: %w", err)
	}

	return nil
}
