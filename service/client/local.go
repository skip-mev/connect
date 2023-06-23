package client

import (
	"context"
	"errors"
	"time"

	"github.com/skip-mev/slinky/service"
	"github.com/skip-mev/slinky/service/types"
)

var _ service.OracleService = (*LocalClient)(nil)

// LocalClient defines an implementation of a local, i.e. in-process, oracle client.
// This client can be used in ABCI++ calls where the application wants the oracle
// process to be run in-process. The client must be started upon app construction
// and stopped upon app shutdown/cleanup.
type LocalClient struct {
	// oracle is the underlying oracle implementation
	oracle types.Oracle
	// timeout is the timeout for the oracle
	timeout time.Duration
}

// NewLocalClient returns a new instance of the LocalClient, given an implementation of the Oracle interface and a timeout.
// Requests to the client will timeout after timeout or if the provided context is cancelled.
func NewLocalClient(o types.Oracle, t time.Duration) *LocalClient {
	return &LocalClient{
		oracle:  o,
		timeout: t,
	}
}

// Prices returns the current prices from the oracle. This method blocks until the oracle returns a response or the context is cancelled (via timeout or by caller).
// This method errors if the oracle is not running, or if the request is nil.
func (c *LocalClient) Prices(ctx context.Context, req *service.QueryPricesRequest) (*service.QueryPricesResponse, error) {
	// check that the request is non-nil
	if req == nil {
		return nil, types.ErrorNilRequest
	}

	// check that oracle is running
	if !c.oracle.IsRunning() {
		return nil, types.ErrorOracleNotRunning
	}

	resCh := make(chan *service.QueryPricesResponse)

	// set a deadline on the context
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// run the request in a goroutine, to unblock server + ctx cancellation
	go func() {
		// get the prices
		prices := c.oracle.GetPrices()

		// get timestamp
		timestamp := c.oracle.GetLastSyncTime()
		resCh <- &service.QueryPricesResponse{
			Prices:    types.ToReqPrices(prices),
			Timestamp: timestamp,
		}
	}()

	// defer to context closure
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	case resp := <-resCh:
		return resp, nil
	}
}

// Note: Start(ctx) is a blocking call, so the caller will need to run it in a
// goroutine. This method blocks until the underlying oracle is stopped.
func (c *LocalClient) Start(ctx context.Context) error {
	if c.oracle.IsRunning() {
		return errors.New("oracle already running")
	}

	return c.oracle.Start(ctx)
}

// Stop stops the underlying oracle. This method blocks until the underlying oracle is stopped.
func (c *LocalClient) Stop(_ context.Context) error {
	c.oracle.Stop()
	return nil
}
