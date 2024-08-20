package oracle

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"cosmossdk.io/log"
	"google.golang.org/grpc"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/service/servers/oracle/types"
)

var _ OracleClient = (*PriceDaemon)(nil)

type PriceDaemon struct {
	logger log.Logger

	// isRunning is an atomic boolean that indicates whether the daemon is running.
	isRunning atomic.Bool
	// config is the configuration of the daemon.
	config config.AppConfig
	// client is the underlying oracle client used to fetch prices.
	OracleClient
	// latestResponse is the latest price response fetched by the daemon.
	resp ThreadSafeResponse
	// doneCh is a channel that is closed when the daemon is stopped.
	doneCh chan struct{}
}

// NewPriceDaemon creates a new price daemon with the given configuration.
func NewPriceDaemon(
	logger log.Logger,
	cfg config.AppConfig,
	client OracleClient,
) (*PriceDaemon, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	if err := cfg.ValidateBasic(); err != nil {
		return nil, err
	}

	if client == nil {
		return nil, fmt.Errorf("oracle client cannot be nil")
	}

	return &PriceDaemon{
		logger:       logger.With("process", "price_daemon"),
		config:       cfg,
		OracleClient: client,
		doneCh:       make(chan struct{}),
	}, nil
}

// Start starts the price daemon. This method will block until the daemon is stopped.
func (d *PriceDaemon) Start(ctx context.Context) error {
	if err := d.OracleClient.Start(ctx); err != nil {
		return err
	}
	defer d.OracleClient.Stop()

	ticker := time.NewTicker(d.config.Interval)
	defer ticker.Stop()

	d.logger.Info("starting price daemon")
	d.isRunning.Store(true)
	defer d.isRunning.Store(false)

	for {
		select {
		case <-ctx.Done():
			d.logger.Info("stopping price daemon from context")
			return ctx.Err()
		case <-d.doneCh:
			d.logger.Info("price daemon stopped")
			return nil
		case <-ticker.C:
			d.fetchPrices(ctx)
		}
	}
}

// fetchPrices fetches the latest prices from the oracle client.
func (d *PriceDaemon) fetchPrices(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			d.logger.Error("recovered from panic", "err", r)
		}
	}()

	d.logger.Debug("fetching prices")

	fetchCtx, cancel := context.WithTimeout(ctx, d.config.ClientTimeout)
	defer cancel()

	resp, err := d.OracleClient.Prices(fetchCtx, &types.QueryPricesRequest{})
	if err != nil {
		d.logger.Error(
			"failed to fetch prices from sidecar",
			"err", err,
			"address", d.config.OracleAddress,
		)

		return
	}

	ts := time.Now()
	d.logger.Debug("fetched prices", "timestamp", ts, "prices", resp.Prices)
	d.resp.Update(resp)
}

// Prices returns the latest price response fetched by the daemon. If the latest response
// is too stale, an error is returned.
func (d *PriceDaemon) Prices(
	_ context.Context,
	_ *types.QueryPricesRequest,
	_ ...grpc.CallOption,
) (*types.QueryPricesResponse, error) {
	latest, ts := d.resp.Get()
	if latest == nil {
		d.logger.Error("no prices fetched by price daemon yet")
		return nil, fmt.Errorf("no prices fetched by price daemon yet")
	}

	if time.Since(ts) > d.config.PriceTTL {
		d.logger.Error(
			"latest prices from the price daemon are too stale",
			"last_fetched_at", ts.String(),
			"diff", time.Since(ts).String(),
			"ttl", d.config.PriceTTL.String(),
		)

		return nil, fmt.Errorf(
			"latest prices from the price daemon are too stale; last fetched at %s; diff %s ago",
			ts.Format(time.RFC3339),
			time.Since(ts).String(),
		)
	}

	return latest, nil
}

// Stop stops the price daemon.
func (d *PriceDaemon) Stop() error {
	if d.isRunning.Load() {
		d.doneCh <- struct{}{}
		close(d.doneCh)
	}

	return nil
}

// ThreadSafeResponse is a thread-safe wrapper around a QueryPricesResponse.
type ThreadSafeResponse struct {
	sync.Mutex

	resp      *types.QueryPricesResponse
	timestamp time.Time
}

// NewThreadSafeResponse creates a new thread-safe response.
func NewThreadSafeResponse() *ThreadSafeResponse {
	return &ThreadSafeResponse{
		resp:      nil,
		timestamp: time.Time{},
	}
}

// Update updates the response and timestamp of the thread-safe response.
func (r *ThreadSafeResponse) Update(resp *types.QueryPricesResponse) {
	r.Lock()
	defer r.Unlock()

	r.resp = resp
	r.timestamp = time.Now()
}

// Get returns the response and timestamp of the thread-safe response.
func (r *ThreadSafeResponse) Get() (*types.QueryPricesResponse, time.Time) {
	r.Lock()
	defer r.Unlock()

	return r.resp, r.timestamp
}
