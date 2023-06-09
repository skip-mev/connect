package client

import (
	"context"

	"github.com/skip-mev/slinky/service"
)

var _ service.OracleService = (*GRPCClient)(nil)

// GRPCClient defines an implementation of a gRPC oracle client. This client can
// be used in ABCI++ calls where the application wants the oracle process to be
// run out-of-process. The client must be started upon app construction and
// stopped upon app shutdown/cleanup.
type GRPCClient struct{}

func NewGRPCClient() *GRPCClient {
	return &GRPCClient{}
}

func (c *GRPCClient) Prices(_ context.Context, req *service.QueryPricesRequest) (*service.QueryPricesResponse, error) {
	panic("not implemented")
}

// Note: Start(ctx) is a blocking call, so the caller will need to run it in a
// goroutine.
func (c *GRPCClient) Start(ctx context.Context) error {
	panic("not implemented")
}

func (c *GRPCClient) Stop(ctx context.Context) error {
	panic("not implemented")
}
