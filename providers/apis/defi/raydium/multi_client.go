package raydium

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	slinkyhttp "github.com/skip-mev/slinky/pkg/http"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
)

// MultiJSONRPCClient is an implementation of the SolanaJSONRPCClient interface that delegates
// requests to multiple underlying clients, and aggregates over all provided responses.
type MultiJSONRPCClient struct {
	logger     *zap.Logger
	api        config.APIConfig
	apiMetrics metrics.APIMetrics

	// underlying clients
	clients []SolanaJSONRPCClient
}

func NewMultiJSONRPCClient(
	logger *zap.Logger,
	api config.APIConfig,
	apiMetrics metrics.APIMetrics,
	clients []SolanaJSONRPCClient,
) *MultiJSONRPCClient {
	return &MultiJSONRPCClient{
		logger:     logger,
		api:        api,
		apiMetrics: apiMetrics,
		clients:    clients,
	}
}

// NewMultiJSONRPCClientFromEndpoints creates a new MultiJSONRPCClient from a list of endpoints.
func NewMultiJSONRPCClientFromEndpoints(
	logger *zap.Logger,
	api config.APIConfig,
	apiMetrics metrics.APIMetrics,
) (*MultiJSONRPCClient, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	if apiMetrics == nil {
		return nil, fmt.Errorf("metrics cannot be nil")
	}

	if err := api.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid api config: %w", err)
	}

	if len(api.Endpoints) == 0 {
		return nil, fmt.Errorf("invalid endpoint: no endpoints provided")
	}

	var err error
	clients := make([]SolanaJSONRPCClient, len(api.Endpoints))
	for i := range api.Endpoints {
		clients[i], err = solanaClientFromEndpoint(api.Endpoints[i])
		if err != nil {
			return nil, fmt.Errorf("failed to create solana client from endpoint: %w", err)
		}
	}

	return NewMultiJSONRPCClient(
		logger,
		api,
		apiMetrics,
		clients,
	), nil
}

// solanaClientFromEndpoint creates a new SolanaJSONRPCClient from an endpoint.
func solanaClientFromEndpoint(endpoint config.Endpoint) (SolanaJSONRPCClient, error) {
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

// GetMultipleAccountsWithOpts delegates the request to all underlying clients and applies a filter
// to the responses.
func (c *MultiJSONRPCClient) GetMultipleAccountsWithOpts(
	ctx context.Context,
	accounts []solana.PublicKey,
	opts *rpc.GetMultipleAccountsOpts,
) (*rpc.GetMultipleAccountsResult, error) {
	// Create a channel to receive the responses from the underlying clients
	responsesCh := make(chan *rpc.GetMultipleAccountsResult, len(c.clients))

	// spawn a goroutine for each client to fetch the accounts
	var wg sync.WaitGroup
	wg.Add(len(c.clients))

	for i := range c.clients {
		url := c.api.Endpoints[i].URL
		go func(client SolanaJSONRPCClient) {
			// Observe the latency of the request.
			start := time.Now()
			defer func() {
				wg.Done()
				c.apiMetrics.ObserveProviderResponseLatency(c.api.Name, time.Since(start))
			}()

			resp, err := client.GetMultipleAccountsWithOpts(ctx, accounts, opts)
			if err != nil {
				c.apiMetrics.AddRPCStatusCode(c.api.Name, metrics.RPCCodeError)
				c.logger.Error("failed to fetch accounts", zap.String("url", url), zap.Error(err))
				return
			}

			c.apiMetrics.AddRPCStatusCode(c.api.Name, metrics.RPCCodeOK)
			responsesCh <- resp
			c.logger.Debug("successfully fetched accounts", zap.String("url", url))
		}(c.clients[i])
	}

	// close the channel once all responses are received, or the context is cancelled
	go func() {
		defer close(responsesCh)
		wg.Wait()
	}()

	responses := make([]*rpc.GetMultipleAccountsResult, 0, len(c.clients))
	for resp := range responsesCh {
		responses = append(responses, resp)
	}

	// filter the responses
	return filterAccountsResponses(responses)
}

// filterAccountsResponses chooses the rpc response with the highest slot number.
func filterAccountsResponses(responses []*rpc.GetMultipleAccountsResult) (*rpc.GetMultipleAccountsResult, error) {
	var (
		maxSlot uint64
		maxResp *rpc.GetMultipleAccountsResult
	)

	if len(responses) == 0 {
		return nil, fmt.Errorf("no responses to filter")
	}

	for _, resp := range responses {
		if resp.Context.Slot > maxSlot {
			maxSlot = resp.Context.Slot
			maxResp = resp
		}
	}

	return maxResp, nil
}
