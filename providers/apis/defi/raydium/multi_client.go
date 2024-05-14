package raydium

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
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
) SolanaJSONRPCClient {
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
) (SolanaJSONRPCClient, error) {
	if len(api.Endpoints) == 0 {
		return nil, fmt.Errorf("no endpoints provided")
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
		logger.With(zap.String("multi_client", Name)),
		api,
		apiMetrics,
		clients,
	), nil
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
