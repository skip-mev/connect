package oracle

import (
	"context"
	"fmt"
	"sync"
	"time"

	"cosmossdk.io/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/service/metrics"
	"github.com/skip-mev/slinky/service/servers/oracle/types"
)

var _ OracleClient = (*GRPCClient)(nil)

// OracleClient defines the interface that will be utilized by the application
// to query the oracle service. This interface is meant to be implemented by
// the gRPC client that connects to the oracle service.
//
//go:generate mockery --name OracleClient --filename mock_oracle_client.go
type OracleClient interface { //nolint
	types.OracleClient

	// Start starts the oracle client.
	Start() error

	// Stop stops the oracle client.
	Stop() error
}

// GRPCClient defines an implementation of a gRPC oracle client. This client can
// be used in ABCI++ calls where the application wants the oracle process to be
// run out-of-process. The client must be started upon app construction and
// stopped upon app shutdown/cleanup.
type GRPCClient struct {
	// address of remote oracle server
	addr string
	// underlying oracle client
	client types.OracleClient
	// underlying grpc connection
	conn *grpc.ClientConn
	// timeout for the client, Price requests will block for this duration.
	timeout time.Duration
	// mutex to protect the client
	mtx sync.Mutex
}

// NewGRPCClient creates a new grpc client of the oracle service, given the
// address of the oracle server and a timeout for the client.
func NewGRPCClient(addr string, t time.Duration) OracleClient {
	return &GRPCClient{
		addr:    addr,
		timeout: t,
		mtx:     sync.Mutex{},
	}
}

// NewGRPCClientFromConfig creates a new grpc client of the oracle service, given the
// oracle app config.
func NewGRPCClientFromConfig(logger log.Logger, config config.AppConfig) OracleClient {
	if config.MetricsEnabled {
		return NewMetricsClient(
			logger,
			NewGRPCClient(config.OracleAddress, config.ClientTimeout),
			metrics.NewMetricsFromConfig(config),
		)
	}

	return NewGRPCClient(config.OracleAddress, config.ClientTimeout)
}

// Start starts the GRPC client. This method dials the remote oracle-service
// and errors if the connection fails.
func (c *GRPCClient) Start() error {
	conn, err := grpc.Dial(
		c.addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to dial oracle gRPC server: %w", err)
	}

	c.mtx.Lock()
	c.client = types.NewOracleClient(conn)
	c.conn = conn
	c.mtx.Unlock()

	return nil
}

// Stop stops the GRPC client. This method closes the connection to the remote.
func (c *GRPCClient) Stop() error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}

// Prices returns the prices from the remote oracle service. This method blocks for the timeout duration configured on the client,
// otherwise it returns the response from the remote oracle.
func (c *GRPCClient) Prices(ctx context.Context, req *types.QueryPricesRequest, _ ...grpc.CallOption) (*types.QueryPricesResponse, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	// set deadline on the context
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	if c.client == nil {
		return nil, fmt.Errorf("oracle client not started")
	}

	return c.client.Prices(ctx, req, grpc.WaitForReady(true))
}
