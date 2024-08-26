package raydium

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/providers/apis/defi/types"
	"github.com/skip-mev/connect/v2/providers/base/api/metrics"
)

// MultiJSONRPCClient is an implementation of the SolanaJSONRPCClient interface that delegates
// requests to multiple underlying clients, and aggregates over all provided responses.
type MultiJSONRPCClient struct {
	logger     *zap.Logger
	api        config.APIConfig
	apiMetrics metrics.APIMetrics

	// underlying clients
	clients []SolanaJSONRPCClient

	blockAgeChecker types.BlockAgeChecker
}

// NewMultiJSONRPCClient returns a new MultiJSONRPCClient.
func NewMultiJSONRPCClient(
	logger *zap.Logger,
	api config.APIConfig,
	apiMetrics metrics.APIMetrics,
	clients []SolanaJSONRPCClient,
) SolanaJSONRPCClient {
	return &MultiJSONRPCClient{
		logger:          logger,
		api:             api,
		apiMetrics:      apiMetrics,
		clients:         clients,
		blockAgeChecker: types.NewBlockAgeChecker(api.MaxBlockHeightAge),
	}
}

// NewMultiJSONRPCClientFromEndpoints creates a new MultiJSONRPCClient from a list of endpoints.
func NewMultiJSONRPCClientFromEndpoints(
	logger *zap.Logger,
	api config.APIConfig,
	apiMetrics metrics.APIMetrics,
) (SolanaJSONRPCClient, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	if err := api.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid api config: %w", err)
	}

	if !api.Enabled {
		return nil, fmt.Errorf("api is not enabled")
	}

	if api.Name != Name {
		return nil, fmt.Errorf("invalid api name; expected %s, got %s", Name, api.Name)
	}

	if len(api.Endpoints) == 0 {
		return nil, fmt.Errorf("no endpoints provided")
	}

	if apiMetrics == nil {
		return nil, fmt.Errorf("metrics is nil")
	}

	var err error
	clients := make([]SolanaJSONRPCClient, len(api.Endpoints))
	for i := range api.Endpoints {
		clients[i], err = solanaClientFromEndpoint(api.Endpoints[i])
		if err != nil {
			return nil, fmt.Errorf("failed to create solana client from endpoint: %w", err)
		}
	}

	return &MultiJSONRPCClient{
		logger:          logger.With(zap.String("multi_client", Name)),
		api:             api,
		apiMetrics:      apiMetrics,
		clients:         clients,
		blockAgeChecker: types.NewBlockAgeChecker(api.MaxBlockHeightAge),
	}, nil
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
		// Pin
		url := c.api.Endpoints[i].URL
		index := i
		go func(client SolanaJSONRPCClient) {
			// Observe the latency of the request.
			start := time.Now()
			defer func() {
				wg.Done()
				c.apiMetrics.ObserveProviderResponseLatency(c.api.Name, metrics.RedactedEndpointURL(index), time.Since(start))
			}()

			resp, err := client.GetMultipleAccountsWithOpts(ctx, accounts, opts)
			if err != nil {
				c.apiMetrics.AddRPCStatusCode(c.api.Name, metrics.RedactedEndpointURL(index), metrics.RPCCodeError)
				c.logger.Error("failed to fetch accounts", zap.String("url", url), zap.Error(err))
				return
			}

			c.apiMetrics.AddRPCStatusCode(c.api.Name, metrics.RedactedEndpointURL(index), metrics.RPCCodeOK)
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
	return c.filterAccountsResponses(responses)
}

// filterAccountsResponses chooses the rpc response with the highest slot number.
func (c *MultiJSONRPCClient) filterAccountsResponses(responses []*rpc.GetMultipleAccountsResult) (*rpc.GetMultipleAccountsResult, error) {
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

	// check the block height (slot)
	if valid := c.blockAgeChecker.IsHeightValid(maxSlot); !valid {
		return nil, fmt.Errorf("height %d is stale and older than %d", maxSlot, c.api.MaxBlockHeightAge)
	}

	return maxResp, nil
}
