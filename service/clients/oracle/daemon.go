package oracle

import (
	"context"
	"fmt"
	"sync"
	"time"

	"cosmossdk.io/log"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/service/servers/oracle/types"
	"google.golang.org/grpc"
)

var _ OracleClient = (*PriceDaemon)(nil)

type PriceDaemon struct {
	mutex  sync.Mutex
	logger log.Logger

	// config is the configuration of the daemon.
	config config.AppConfig
	// client is the underlying oracle client used to fetch prices.
	OracleClient
	// latestResponse is the latest price response fetched by the daemon.
	latestResponse *types.QueryPricesResponse
	// timestamp is the time at which the latest response was fetched.
	timestamp time.Time
	// doneCh is a channel that is closed when the daemon is stopped.
	doneCh chan error
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
		doneCh:       make(chan error),
	}, nil
}

// Start starts the price daemon. This method will block until the daemon is stopped.
func (d *PriceDaemon) Start(ctx context.Context) error {
	if err := d.OracleClient.Start(ctx); err != nil {
		return err
	}

	ticker := time.NewTicker(d.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			d.logger.Info("stopping price daemon")
			return nil
		case err := <-d.doneCh:
			return err
		case <-ticker.C:
			d.fetchPrices()
		}
	}
}

// fetchPrices fetches the latest prices from the oracle client.
func (d *PriceDaemon) fetchPrices() {
	d.logger.Debug("fetching prices")

	fetchCtx, cancel := context.WithTimeout(context.Background(), d.config.ClientTimeout)
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

	d.mutex.Lock()
	d.latestResponse = resp
	d.timestamp = time.Now()
	d.mutex.Unlock()
}

// Prices returns the latest price response fetched by the daemon. If the latest response
// is too stale, an error is returned.
func (d *PriceDaemon) Prices(
	_ context.Context,
	_ *types.QueryPricesRequest,
	_ ...grpc.CallOption,
) (*types.QueryPricesResponse, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.latestResponse == nil {
		return nil, fmt.Errorf("no prices fetched by price daemon yet")
	}

	if time.Since(d.timestamp) > d.config.MaxAge {
		return nil, fmt.Errorf("latest prices from the price daemon are too stale; last fetched at %s", d.timestamp.Format(time.RFC3339))
	}

	return d.latestResponse, nil
}

// Stop stops the price daemon.
func (d *PriceDaemon) Stop() error {
	err := d.OracleClient.Stop()
	d.doneCh <- err
	return err
}
