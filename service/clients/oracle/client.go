package oracle

import (
	"context"
	"fmt"
	"sync"
	"time"

	"cosmossdk.io/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/skip-mev/connect/v2/oracle/config"
	connectgrpc "github.com/skip-mev/connect/v2/pkg/grpc"
	"github.com/skip-mev/connect/v2/service/metrics"
	"github.com/skip-mev/connect/v2/service/servers/oracle/types"
)

var _ OracleClient = (*GRPCClient)(nil)

// GRPCClient defines an implementation of a gRPC oracle client. This client can
// be used in ABCI++ calls where the application wants the oracle process to be
// run out-of-process. The client must be started upon app construction and
// stopped upon app shutdown/cleanup.
type GRPCClient struct {
	logger log.Logger
	mutex  sync.Mutex

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
	// blockingDial is a parameter which determines whether the client should block on dialing the server
	blockingDial bool
}

// NewClientFromConfig creates a new grpc client of the oracle service with the given
// app configuration. This returns an error if the configuration is invalid.
func NewClientFromConfig(
	cfg config.AppConfig,
	logger log.Logger,
	metrics metrics.Metrics,
	opts ...Option,
) (OracleClient, error) {
	if err := cfg.ValidateBasic(); err != nil {
		return nil, err
	}

	if !cfg.Enabled {
		return &NoOpClient{}, nil
	}

	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	if metrics == nil {
		return nil, fmt.Errorf("metrics cannot be nil")
	}

	return NewClient(logger, cfg.OracleAddress, cfg.ClientTimeout, metrics, opts...)
}

// NewPriceDaemonClientFromConfig creates a new grpc client of the oracle service with the given
// app configuration. This returns an error if the configuration is invalid. Specifically, this
// client will be a daemon client that has prices available in constant time.
func NewPriceDaemonClientFromConfig(
	cfg config.AppConfig,
	logger log.Logger,
	metrics metrics.Metrics,
	opts ...Option,
) (OracleClient, error) {
	if !cfg.Enabled {
		return &NoOpClient{}, nil
	}

	client, err := NewClientFromConfig(cfg, logger, metrics, opts...)
	if err != nil {
		return nil, err
	}

	return NewPriceDaemon(logger, cfg, client)
}

// NewClient creates a new grpc client of the oracle service with the given
// address and timeout.
func NewClient(
	logger log.Logger,
	addr string,
	timeout time.Duration,
	metrics metrics.Metrics,
	opts ...Option,
) (OracleClient, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	if metrics == nil {
		return nil, fmt.Errorf("metrics cannot be nil")
	}

	if timeout <= 0 {
		return nil, fmt.Errorf("timeout must be positive")
	}

	client := &GRPCClient{
		logger:  logger,
		addr:    addr,
		timeout: timeout,
		metrics: metrics,
	}

	// apply options
	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

// Start starts the GRPC client. This method dials the remote oracle-service
// and errors if the connection fails. This method may block (depending on the blockingDial option).
func (c *GRPCClient) Start(ctx context.Context) error {
	c.logger.Info("starting oracle client", "addr", c.addr)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// dial the client, but defer to context closure, if necessary
	var (
		conn *grpc.ClientConn
		err  error
		done = make(chan struct{})
	)
	go func() {
		defer close(done)
		conn, err = connectgrpc.NewClient(c.addr, opts...)

		// attempt to connect + wait for change in connection state
		if c.blockingDial {
			// connect
			conn.Connect()

			if err == nil {
				conn.WaitForStateChange(ctx, connectivity.Ready)
			}
		}
	}()

	// wait for either the context to close or the dial to complete
	select {
	case <-ctx.Done():
		err = fmt.Errorf("context closed before oracle client could start: %w", ctx.Err())
	case <-done:
	}
	if err != nil {
		c.logger.Error("failed to dial oracle gRPC server", "err", err)
		return fmt.Errorf("failed to dial oracle gRPC server: %w", err)
	}

	c.mutex.Lock()
	c.client = types.NewOracleClient(conn)
	c.conn = conn
	c.mutex.Unlock()

	c.logger.Info("oracle client started")

	return nil
}

// Stop stops the GRPC client. This method closes the connection to the remote.
func (c *GRPCClient) Stop() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

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
	c.mutex.Lock()
	defer c.mutex.Unlock()

	start := time.Now()
	defer func() {
		// Observe the duration of the call as well as the error.
		c.metrics.ObserveOracleResponseLatency(time.Since(start))
		c.metrics.AddOracleResponse(metrics.StatusFromError(err))
	}()

	// set deadline on the context
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	if c.client == nil {
		return nil, fmt.Errorf("oracle client not started")
	}

	return c.client.Prices(ctx, req, grpc.WaitForReady(true))
}

func (c *GRPCClient) MarketMap(ctx context.Context, req *types.QueryMarketMapRequest, _ ...grpc.CallOption) (res *types.QueryMarketMapResponse, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	start := time.Now()
	defer func() {
		// Observe the duration of the call as well as the error.
		c.metrics.ObserveOracleResponseLatency(time.Since(start))
		c.metrics.AddOracleResponse(metrics.StatusFromError(err))
	}()

	// set deadline on the context
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	if c.client == nil {
		return nil, fmt.Errorf("oracle client not started")
	}

	return c.client.MarketMap(ctx, req, grpc.WaitForReady(true))
}

// Version returns the version of the oracle service.
func (c *GRPCClient) Version(ctx context.Context, req *types.QueryVersionRequest, _ ...grpc.CallOption) (res *types.QueryVersionResponse, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	start := time.Now()
	defer func() {
		// Observe the duration of the call as well as the error.
		c.metrics.ObserveOracleResponseLatency(time.Since(start))
		c.metrics.AddOracleResponse(metrics.StatusFromError(err))
	}()

	// set deadline on the context
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	if c.client == nil {
		return nil, fmt.Errorf("oracle client not started")
	}

	return c.client.Version(ctx, req, grpc.WaitForReady(true))
}
