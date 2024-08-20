package ethmulticlient

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/providers/base/api/metrics"

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
	apiMetrics  metrics.APIMetrics
	api         config.APIConfig
	redactedURL string

	// client is the underlying rpc client.
	client *rpc.Client
}

// NewGoEthereumClientImpl creates an EVMClient via a config.Endpoint. This
// includes optional authentication via a specified http header key and value.
func NewGoEthereumClientImpl(
	ctx context.Context,
	apiMetrics metrics.APIMetrics,
	api config.APIConfig,
	index int,
) (EVMClient, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is nil")
	}

	if apiMetrics == nil {
		return nil, fmt.Errorf("api metrics nil")
	}

	if err := api.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid api config: %w", err)
	}

	if !api.Enabled {
		return nil, fmt.Errorf("api config for %s is not enabled", api.Name)
	}

	// We need to have at least one endpoint to connect to.
	if len(api.Endpoints) < index {
		return nil, fmt.Errorf("expected endpoint at index %d, got %d endpoints", index, len(api.Endpoints))
	}

	var opts []rpc.ClientOption
	endpoint := api.Endpoints[index] // pin
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
		apiMetrics:  apiMetrics,
		api:         api,
		redactedURL: metrics.RedactedEndpointURL(index),
		client:      client,
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
func (c *GoEthereumClientImpl) BatchCallContext(ctx context.Context, calls []rpc.BatchElem) (err error) {
	start := time.Now()
	defer func() {
		c.apiMetrics.ObserveProviderResponseLatency(c.api.Name, c.redactedURL, time.Since(start))
	}()

	if err = c.client.BatchCallContext(ctx, calls); err != nil {
		c.apiMetrics.AddRPCStatusCode(c.api.Name, c.redactedURL, metrics.RPCCodeError)
		return
	}

	c.apiMetrics.AddRPCStatusCode(c.api.Name, c.redactedURL, metrics.RPCCodeOK)
	return
}
