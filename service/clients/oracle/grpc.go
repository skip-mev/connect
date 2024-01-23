package oracle

import (
	"context"
	"fmt"
	"net/url"
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

	// Start starts the oracle client. This should connect to the remote oracle
	// service and return an error if the connection fails.
	Start() error

	// Stop stops the oracle client.
	Stop() error
}

// GRPCClient defines an implementation of a gRPC oracle client. This client can
// be used in ABCI++ calls where the application wants the oracle process to be
// run out-of-process. The client must be started upon app construction and
// stopped upon app shutdown/cleanup.
type GRPCClient struct {
	logger log.Logger
	mtx    sync.Mutex

	// address of remote oracle server
	addr string
	// underlying oracle client
	client types.OracleClient
	// underlying grpc connection
	conn *grpc.ClientConn
	// timeout for the client, Price requests will block for this duration.
	timeout time.Duration
	// metrics contains the instrumentation for the oracle client
	metrics metrics.Metrics
}

// NewGRPCClientFromConfig creates a new grpc client of the oracle service with the given
// app configuration. This returns an error if the configuration is invalid.
func NewGRPCClientFromConfig(
	cfg config.AppConfig,
	logger log.Logger,
	metrics metrics.Metrics,
) (OracleClient, error) {
	if err := cfg.ValidateBasic(); err != nil {
		return nil, err
	}

	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	if metrics == nil {
		return nil, fmt.Errorf("metrics cannot be nil")
	}

	return &GRPCClient{
		logger:  logger,
		addr:    cfg.OracleAddress,
		timeout: cfg.ClientTimeout,
		metrics: metrics,
	}, nil
}

// NewGRPCClient creates a new grpc client of the oracle service with the given
// address and timeout.
func NewGRPCClient(
	logger log.Logger,
	addr string,
	timeout time.Duration,
	metrics metrics.Metrics,
) (OracleClient, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	if _, err := url.ParseRequestURI(addr); err != nil {
		return nil, fmt.Errorf("invalid oracle address: %w", err)
	}

	if metrics == nil {
		return nil, fmt.Errorf("metrics cannot be nil")
	}

	if timeout <= 0 {
		return nil, fmt.Errorf("timeout must be positive")
	}

	return &GRPCClient{
		logger:  logger,
		addr:    addr,
		timeout: timeout,
		metrics: metrics,
	}, nil
}

// Start starts the GRPC client. This method dials the remote oracle-service
// and errors if the connection fails.
func (c *GRPCClient) Start() error {
	c.logger.Info("starting oracle client", "addr", c.addr)

	conn, err := grpc.Dial(
		c.addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		c.logger.Error("failed to dial oracle gRPC server", "err", err)
		return fmt.Errorf("failed to dial oracle gRPC server: %w", err)
	}

	c.mtx.Lock()
	c.client = types.NewOracleClient(conn)
	c.conn = conn
	c.mtx.Unlock()

	c.logger.Info("oracle client started")
	return nil
}

// Stop stops the GRPC client. This method closes the connection to the remote.
func (c *GRPCClient) Stop() error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.logger.Info("stopping oracle client")
	if c.conn == nil {
		return nil
	}

	err := c.conn.Close()
	c.logger.Info("oracle client stopped", "err", err)
	return err
}

// Prices returns the prices from the remote oracle service. This method blocks for the timeout duration configured on the client,
// otherwise it returns the response from the remote oracle.
func (c *GRPCClient) Prices(
	ctx context.Context,
	req *types.QueryPricesRequest,
	_ ...grpc.CallOption,
) (resp *types.QueryPricesResponse, err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	start := time.Now()
	c.logger.Info("calling oracle client", "timestamp", start)

	defer func() {
		if err != nil {
			c.logger.Error("oracle client returned error", "err", err)
		} else {
			c.logger.Info("oracle client returned response", "resp", resp)
		}

		// Observe the duration of the call as well as the error.
		c.metrics.ObserveOracleResponseLatency(time.Since(start))
		c.metrics.AddOracleResponse(metrics.StatusFromError(err))
	}()

	// set deadline on the context
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	if c.client == nil {
		c.logger.Error("oracle client not started")
		return nil, fmt.Errorf("oracle client not started")
	}

	return c.client.Prices(ctx, req, grpc.WaitForReady(true))
}
