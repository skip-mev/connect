package client

import (
	"context"
	"fmt"

	"github.com/skip-mev/slinky/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ service.OracleService = (*GRPCClient)(nil)

// GRPCClient defines an implementation of a gRPC oracle client. This client can
// be used in ABCI++ calls where the application wants the oracle process to be
// run out-of-process. The client must be started upon app construction and
// stopped upon app shutdown/cleanup.
type GRPCClient struct {
	addr   string
	client service.OracleClient
	conn   *grpc.ClientConn
}

func NewGRPCClient(addr string) *GRPCClient {
	return &GRPCClient{
		addr: addr,
	}
}

func (c *GRPCClient) Start(ctx context.Context) error {
	conn, err := grpc.Dial(
		c.addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to dial oracle gRPC server: %w", err)
	}

	c.client = service.NewOracleClient(conn)
	c.conn = conn

	return nil
}

func (c *GRPCClient) Stop(ctx context.Context) error {
	return c.conn.Close()
}

func (c *GRPCClient) Prices(ctx context.Context, req *service.QueryPricesRequest) (*service.QueryPricesResponse, error) {
	return c.client.Prices(ctx, req, grpc.WaitForReady(true))
}
