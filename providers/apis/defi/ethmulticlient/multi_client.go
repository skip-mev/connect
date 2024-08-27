package ethmulticlient

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/skip-mev/connect/v2/providers/apis/defi/types"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/providers/base/api/metrics"
)

// MultiRPCClient implements the EVMClient interface by calling multiple underlying EVMClients and choosing
// the best response. Specifically, it calls eth_blockNumber on each client and chooses the response with the
// highest block number.
type MultiRPCClient struct {
	logger *zap.Logger
	api    config.APIConfig

	// underlying clients
	clients []EVMClient

	blockAgeChecker types.BlockAgeChecker
}

// NewMultiRPCClient returns a new MultiRPCClient.
func NewMultiRPCClient(
	logger *zap.Logger,
	api config.APIConfig,
	clients []EVMClient,
) EVMClient {
	return &MultiRPCClient{
		logger:          logger,
		clients:         clients,
		api:             api,
		blockAgeChecker: types.NewBlockAgeChecker(api.MaxBlockHeightAge),
	}
}

// NewMultiRPCClientFromEndpoints creates a MultiRPCClient from config endpoints.
func NewMultiRPCClientFromEndpoints(
	ctx context.Context,
	logger *zap.Logger,
	api config.APIConfig,
	apiMetrics metrics.APIMetrics,
) (EVMClient, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is nil")
	}

	if logger == nil {
		return nil, fmt.Errorf("logger is nil")
	}

	if err := api.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid api config: %w", err)
	}

	if !api.Enabled {
		return nil, fmt.Errorf("api config for %s is not enabled", api.Name)
	}

	if len(api.Endpoints) == 0 {
		return nil, fmt.Errorf("no endpoints provided")
	}

	clients := make([]EVMClient, len(api.Endpoints))
	for i, endpoint := range api.Endpoints {
		// Pin the endpoint directly into a copy of the config.
		var err error
		clients[i], err = NewGoEthereumClientImpl(ctx, apiMetrics, api, i)
		if err != nil {
			logger.Error(
				"endpoint failed to construct client",
				zap.String("endpoint.URL", endpoint.URL),
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to create eth client from endpoint: %w", err)
		}
	}

	return &MultiRPCClient{
		logger:          logger.With(zap.String("multi_client", api.Name)),
		api:             api,
		clients:         clients,
		blockAgeChecker: types.NewBlockAgeChecker(api.MaxBlockHeightAge),
	}, nil
}

// define a result struct that go routines will populate and append to a slice when they complete their request.
type result struct {
	height  uint64
	results []rpc.BatchElem
	err     error
}

// BatchCallContext injects a call to eth_blockNumber, and makes batch calls to the underlying EVMClients.
// It returns the response that has the greatest height from the eth_blockNumber call. An error is returned
// only when no client was able to successfully provide a height or errored when sending the BatchCall.
func (m *MultiRPCClient) BatchCallContext(ctx context.Context, batchElems []rpc.BatchElem) error {
	if len(batchElems) == 0 {
		m.logger.Debug("BatchCallContext called with 0 elems")
		return nil
	}

	results := make([]result, len(m.clients))

	wg := new(sync.WaitGroup)
	// this is the index of where we will have an eth_blockNumber call.
	blockNumReqIndex := len(batchElems)
	// for each client, spin up a go routine that executes a BatchCall.
	for clientIdx, client := range m.clients {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			url := m.api.Endpoints[i].URL

			// append an eth_blockNumber call to the requests. we do this because we want the greatest height results only.
			req := make([]rpc.BatchElem, len(batchElems)+1)
			copy(req, batchElems)
			req[blockNumReqIndex] = EthBlockNumberBatchElem()

			err := client.BatchCallContext(ctx, req)

			// if there was an error, or if the block_num request didn't have result / errored
			// we log the error and append to error slice.
			if err != nil || req[blockNumReqIndex].Result == "" || req[blockNumReqIndex].Error != nil {
				resultErr := fmt.Errorf("endpoint request failed: %w, %w", err, req[blockNumReqIndex].Error)
				results[i] = result{0, nil, resultErr}
				m.logger.Debug(
					"endpoint request failed",
					zap.Error(err),
					zap.Any("result", req[blockNumReqIndex].Result),
					zap.Error(req[blockNumReqIndex].Error),
					zap.String("url", url),
				)
				return
			}
			// the batch call succeeded, and the eth_blockNumber call had results.\
			// try to get the block number.
			r, ok := req[blockNumReqIndex].Result.(*string)
			if !ok {
				resultErr := fmt.Errorf("result from eth_blockNumber was not a string")
				results[i] = result{0, nil, resultErr}
				m.logger.Debug(
					"result from eth_blockNumber was not a string",
					zap.String("url", url),
				)
				return
			}

			// decode the new height
			height, err := hexutil.DecodeUint64(*r)
			if err != nil { // if we can't decode the height, log an error.
				resultErr := fmt.Errorf("could not decode hex eth height: %w", err)
				results[i] = result{0, nil, resultErr}
				m.logger.Debug(
					"could not decode hex eth height",
					zap.String("url", url),
					zap.Error(err),
				)
				return
			}
			m.logger.Debug(
				"got height for eth batch request",
				zap.Uint64("height", height),
				zap.String("url", url),
			)
			// append the results, minus the appended eth_blockNumber request.
			results[i] = result{height, req[:blockNumReqIndex], nil}
		}(clientIdx)
	}
	wg.Wait()

	filtered, err := m.filterResponses(results)
	if err != nil {
		return fmt.Errorf("error filtering responses: %w", err)
	}

	// copy the results from the results that had the largest height.
	copy(batchElems, filtered)
	return nil
}

// filterAccountsResponses chooses the rpc response with the highest block number.
func (m *MultiRPCClient) filterResponses(responses []result) ([]rpc.BatchElem, error) {
	// see which of the results had the largest height, and store the index of that result.
	var (
		maxHeight      uint64
		maxHeightIndex int
		errs           = make([]error, len(responses))
	)
	for i, res := range responses {
		errs[i] = res.err
		if res.height > maxHeight {
			maxHeight = res.height
			maxHeightIndex = i
		}
	}
	// maxHeight being 0 means there were no results. something bad happened. return all the errors.
	if maxHeight == 0 {
		err := errors.Join(errs...)
		if err != nil {
			return nil, err
		}
		// this should never happen... but who knows. maybe something terrible happened.
		return nil, errors.New("no errors were encountered, however no go routine was able to report a height")

	}

	// check the block height
	if valid := m.blockAgeChecker.IsHeightValid(maxHeight); !valid {
		return nil, fmt.Errorf("height %d is stale and older than %d", maxHeight, m.api.MaxBlockHeightAge)
	}

	return responses[maxHeightIndex].results, nil
}
