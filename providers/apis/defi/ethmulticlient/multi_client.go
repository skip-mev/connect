package ethmulticlient

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
)

// MultiRPCClient implements the EVMClient interface by calling multiple underlying EVMClients and choosing
// the best response.
type MultiRPCClient struct {
	logger *zap.Logger

	// underlying clients
	clients []EVMClient
}

// NewMultiRPCClient returns a new MultiRPCClient.
func NewMultiRPCClient(
	logger *zap.Logger,
	clients []EVMClient,
) *MultiRPCClient {
	return &MultiRPCClient{
		logger:  logger,
		clients: clients,
	}
}

// NewMultiRPCClientFromEndpoints creates a MultiRPCClient from config endpoints.
func NewMultiRPCClientFromEndpoints(
	ctx context.Context,
	logger *zap.Logger,
	api config.APIConfig,
	apiMetrics metrics.APIMetrics,
) (*MultiRPCClient, error) {
	clients := make([]EVMClient, len(api.Endpoints))
	for i, endpoint := range api.Endpoints {
		var err error
		clients[i], err = NewGoEthereumClientImplFromEndpoint(ctx, apiMetrics, endpoint, api.Name)
		if err != nil {
			logger.Error(
				"endpoint failed to construct client",
				zap.String("endpoint.URL", endpoint.URL),
			)
			return nil, fmt.Errorf("failed to create eth client from endpoint: %w", err)
		}
	}

	return NewMultiRPCClient(logger, clients), nil
}

// BatchCallContext injects a call to eth_blockNumber, and makes batch calls to the underlying EVMClients.
// It returns the first response it sees from a node which has the greatest height.
// An error is returned only when all clients fail.
func (m *MultiRPCClient) BatchCallContext(ctx context.Context, batchElems []rpc.BatchElem) error {
	if len(batchElems) == 0 {
		m.logger.Debug("BatchCallContext called with 0 elems")
		return nil
	}

	req := make([]rpc.BatchElem, len(batchElems)+1)
	copy(req, batchElems)
	blockNumReqIndex := len(batchElems)
	req[blockNumReqIndex] = EthBlockNumberBatchElem()
	errs := fmt.Errorf("all eth client requests failed")

	// TODO(david): consider parallelizing these requests.
	var maxHeight uint64
	for _, client := range m.clients {
		err := client.BatchCallContext(ctx, req)
		if err != nil || req[blockNumReqIndex].Result == "" || req[blockNumReqIndex].Error != nil {
			errs = fmt.Errorf("%w: endpoint request failed: %w, %w", errs, err, req[blockNumReqIndex].Error)
			m.logger.Debug(
				"endpoint request failed",
				zap.Error(err),
				zap.Any("result", req[blockNumReqIndex].Result),
				zap.Error(req[blockNumReqIndex].Error),
			)
			continue
		}

		r, ok := req[blockNumReqIndex].Result.(*string)
		if !ok {
			errs = fmt.Errorf("%w: result from eth_blockNumber was not a string", errs)
			m.logger.Debug(
				"result from eth_blockNumber was not a string",
			)
			continue
		}

		newHeight, err := hexutil.DecodeUint64(*r)
		if err != nil {
			errs = fmt.Errorf("%w: could not decode hex eth height: %w", errs, err)
			m.logger.Debug(
				"could not decode hex eth height",
				zap.Error(err),
			)
			continue
		}
		m.logger.Debug(
			"got height for eth batch request",
			zap.Uint64("height", newHeight),
		)

		if newHeight > maxHeight {
			m.logger.Debug("new max eth height seen",
				zap.Uint64("prev_height", maxHeight),
				zap.Uint64("new_height", newHeight),
			)

			maxHeight = newHeight
			copy(batchElems, req[:blockNumReqIndex])
		}
	}

	if maxHeight == 0 {
		return errs
	}
	return nil
}
