package raydium

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
	"github.com/skip-mev/slinky/oracle/config"
	slinkyhttp "github.com/skip-mev/slinky/pkg/http"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
)

// JSONRPCClient is an implementation of the Solana JSON RPC client with
// additional functionality for metrics and logging.
type JSONRPCClient struct {
	api        config.APIConfig
	apiMetrics metrics.APIMetrics

	// client is the underlying solana-go JSON-RPC client.
	client *rpc.Client
}

// NewJSONRPCClient returns a new JSONRPCClient.
func NewJSONRPCClient(
	api config.APIConfig,
	apiMetrics metrics.APIMetrics,
) (SolanaJSONRPCClient, error) {
	var (
		client *rpc.Client
		err    error
	)
	switch {
	case len(api.Endpoints) == 1:
		client, err = solanaClientFromEndpoint(api.Endpoints[0])
	case len(api.URL) > 0:
		client = rpc.New(api.URL)
	default:
		return nil, fmt.Errorf("no valid endpoints or url were provided")
	}
	if err != nil {
		return nil, err
	}

	return &JSONRPCClient{
		api:        api,
		apiMetrics: apiMetrics,
		client:     client,
	}, nil
}

// GetMultipleAccountsWithOpts is a wrapper around the solana-go GetMultipleAccountsWithOpts method.
func (c *JSONRPCClient) GetMultipleAccountsWithOpts(
	ctx context.Context,
	accounts []solana.PublicKey,
	opts *rpc.GetMultipleAccountsOpts,
) (out *rpc.GetMultipleAccountsResult, err error) {
	start := time.Now()
	defer func() {
		c.apiMetrics.ObserveProviderResponseLatency(c.api.Name, time.Since(start))
	}()

	out, err = c.client.GetMultipleAccountsWithOpts(ctx, accounts, opts)
	if err != nil {
		c.apiMetrics.AddRPCStatusCode(c.api.Name, metrics.RPCCodeError)
		return
	}

	c.apiMetrics.AddRPCStatusCode(c.api.Name, metrics.RPCCodeOK)
	return
}

// solanaClientFromEndpoint creates a new SolanaJSONRPCClient from an endpoint.
func solanaClientFromEndpoint(endpoint config.Endpoint) (*rpc.Client, error) {
	// if authentication is enabled
	if endpoint.Authentication.Enabled() {
		transport := slinkyhttp.NewRoundTripperWithHeaders(map[string]string{
			endpoint.Authentication.APIKeyHeader: endpoint.Authentication.APIKey,
		}, http.DefaultTransport)

		client := rpc.NewWithCustomRPCClient(jsonrpc.NewClientWithOpts(endpoint.URL, &jsonrpc.RPCClientOpts{
			HTTPClient: &http.Client{
				Transport: transport,
			},
		}))

		return client, nil
	}
	return rpc.New(endpoint.URL), nil
}
