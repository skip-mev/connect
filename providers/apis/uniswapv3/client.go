package uniswapv3

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
)

// EVMClient is an interface that abstracts the go-ethereum client.
//
//go:generate mockery --name EVMClient
type EVMClient interface {
	// BatchCallContext is a batch call to the ethereum network.
	BatchCallContext(ctx context.Context, calls []rpc.BatchElem) error
}

var _ EVMClient = (*GoEthereumClientImpl)(nil)

// GoEthereumClient is a go-ethereum client.
type GoEthereumClientImpl struct {
	client *rpc.Client
}

// NewGoEthereumClientImpl returns a new go-ethereum client. This is the default
// implementation that connects to an ethereum node via rpc.
func NewGoEthereumClientImpl(url string) (EVMClient, error) {
	return rpc.Dial(url)
}

// BatchCallContext is a batch call to the ethereum network.
func (c *GoEthereumClientImpl) BatchCallContext(ctx context.Context, calls []rpc.BatchElem) error {
	return c.client.BatchCallContext(ctx, calls)
}
