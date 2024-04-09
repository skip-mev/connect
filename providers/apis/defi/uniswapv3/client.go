package uniswapv3

import (
	"context"

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

// GoEthereumClient is a go-ethereum client implementation using the go-ethereum RPC
// library.
type GoEthereumClientImpl struct {
	client *rpc.Client
}

// NewGoEthereumClientImpl returns a new go-ethereum client. This is the default
// implementation that connects to an ethereum node via rpc.
func NewGoEthereumClientImpl(url string) (EVMClient, error) {
	return rpc.Dial(url)
}

// BatchCallContext sends all given requests as a single batch and waits for the server
// to return a response for all of them. The wait duration is bounded by the context's deadline.
//
// In contrast to CallContext, BatchCallContext only returns errors that have occurred while
// sending the request. Any error specific to a request is reported through the Error field of
// the corresponding BatchElem.
//
// Note that batch calls may not be executed atomically on the server side.
func (c *GoEthereumClientImpl) BatchCallContext(ctx context.Context, calls []rpc.BatchElem) error {
	return c.client.BatchCallContext(ctx, calls)
}
