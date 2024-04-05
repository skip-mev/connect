package raydium

import (
	"context"
	"fmt"
	"sync"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"go.uber.org/zap"

	oracleconfig "github.com/skip-mev/slinky/oracle/config"
)

// MultiJSONRPCClient is an implementation of the SolanaJSONRPCClient interface that delegates
// requests to multiple underlying clients, and aggregates over all provided responses.
type MultiJSONRPCClient struct {
	// underlying clients
	clients []SolanaJSONRPCClient

	// logger
	logger *zap.Logger
}

func NewMultiJSONRPCClient(clients []SolanaJSONRPCClient, logger *zap.Logger) *MultiJSONRPCClient {
	return &MultiJSONRPCClient{
		clients: clients,
		logger:  logger,
	}
}

// NewMultiJSONRPCClientFromEndpoints creates a new MultiJSONRPCClient from a list of endpoints.
func NewMultiJSONRPCClientFromEndpoints(endpoints []oracleconfig.Endpoint, logger *zap.Logger) *MultiJSONRPCClient {
	clients := make([]SolanaJSONRPCClient, len(endpoints))
	for i := range endpoints {
		client := rpc.New(endpoints[i].URL)
		clients[i] = client
	}
	return NewMultiJSONRPCClient(clients, logger)
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
		go func(client SolanaJSONRPCClient) {
			defer wg.Done()
			resp, err := client.GetMultipleAccountsWithOpts(ctx, accounts, opts)
			if err != nil {
				c.logger.Error("failed to fetch accounts", zap.Error(err))
				return
			}
			responsesCh <- resp
		}(c.clients[i])
	}

	// close the channel once all responses are received, or the context is cancelled
	go func() {
		select {
		case <-ctx.Done():
			c.logger.Error("context cancelled")
		case <-channelForWaitGroup(&wg):
		}
		close(responsesCh)
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

// channelForWaitGroup returns a channel that is closed when a waitgroup is done.
func channelForWaitGroup(wg *sync.WaitGroup) chan struct{} {
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()
	return ch
}
