package raydium

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
	"go.uber.org/zap"

	oracleconfig "github.com/skip-mev/slinky/oracle/config"
	slinkyhttp "github.com/skip-mev/slinky/pkg/http"
)

const (
	// authenticationHeaderField is the header field used to store api-key for authentication.
	authenticationHeaderField = "X-Api-Key"
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
func NewMultiJSONRPCClientFromEndpoints(endpoints []oracleconfig.Endpoint, logger *zap.Logger) (*MultiJSONRPCClient, error) {
	clients := make([]SolanaJSONRPCClient, len(endpoints))

	var err error
	for i := range endpoints {
		clients[i], err = solanaClientFromEndpoint(endpoints[i])
		if err != nil {
			return nil, fmt.Errorf("failed to create solana client from endpoint: %w", err)
		}
	}

	return NewMultiJSONRPCClient(clients, logger), nil
}

// solanaClientFromEndpoint creates a new SolanaJSONRPCClient from an endpoint.
func solanaClientFromEndpoint(endpoint oracleconfig.Endpoint) (SolanaJSONRPCClient, error) {
	// if authentication is enabled
	if endpoint.Authentication.Enabled {
		transport := slinkyhttp.NewRoundTripperWithHeaders(map[string]string{
			authenticationHeaderField: endpoint.Authentication.HTTPHeaderAPIKey,
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
