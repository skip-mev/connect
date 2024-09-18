package raydium

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/jsonrpc"

	"github.com/skip-mev/connect/v2/oracle/config"
	connecthttp "github.com/skip-mev/connect/v2/pkg/http"
	"github.com/skip-mev/connect/v2/providers/base/api/metrics"
)

// JSONRPCClient is an implementation of the Solana JSON RPC client with
// additional functionality for metrics and logging.
type JSONRPCClient struct {
	api         config.APIConfig
	apiMetrics  metrics.APIMetrics
	redactedURL string

	// client is the underlying solana-go JSON-RPC client.
	client *rpc.Client
}

// NewJSONRPCClient returns a new JSONRPCClient.
func NewJSONRPCClient(
	api config.APIConfig,
	apiMetrics metrics.APIMetrics,
) (SolanaJSONRPCClient, error) {
	if err := api.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid api config: %w", err)
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

	var (
		client      *rpc.Client
		redactedURL string
		err         error
	)
	switch {
	case len(api.Endpoints) == 1:
		client, err = solanaClientFromEndpoint(api.Endpoints[0])
		redactedURL = metrics.RedactedEndpointURL(0)
	default:
		return nil, fmt.Errorf("no valid endpoints or url were provided")
	}
	if err != nil {
		return nil, err
	}

	return &JSONRPCClient{
		api:         api,
		apiMetrics:  apiMetrics,
		redactedURL: redactedURL,
		client:      client,
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
		c.apiMetrics.ObserveProviderResponseLatency(c.api.Name, c.redactedURL, time.Since(start))
	}()

	out, err = c.client.GetMultipleAccountsWithOpts(ctx, accounts, opts)
	if err != nil {
		c.apiMetrics.AddRPCStatusCode(c.api.Name, c.redactedURL, metrics.RPCCodeError)
		return
	}

	c.apiMetrics.AddRPCStatusCode(c.api.Name, c.redactedURL, metrics.RPCCodeOK)
	return
}

// solanaClientFromEndpoint creates a new SolanaJSONRPCClient from an endpoint.
func solanaClientFromEndpoint(endpoint config.Endpoint) (*rpc.Client, error) {
	opts := []connecthttp.HeaderOption{
		connecthttp.WithConnectVersionUserAgent(),
	}

	// if authentication is enabled
	if endpoint.Authentication.Enabled() {
		// add authentication header
		opts = append(opts, connecthttp.WithAuthentication(endpoint.Authentication.APIKeyHeader, endpoint.Authentication.APIKey))
	}

	transport := connecthttp.NewRoundTripperWithHeaders(http.DefaultTransport, opts...)

	client := rpc.NewWithCustomRPCClient(jsonrpc.NewClientWithOpts(endpoint.URL, &jsonrpc.RPCClientOpts{
		HTTPClient: &http.Client{
			Transport: transport,
		},
	}))

	return client, nil
}
