package marketmap

import (
	"context"
	"math/big"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/providers/base"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

type GRPCClient struct {
	logger *zap.Logger
	mutex  sync.Mutex
	client mmtypes.QueryClient

	// cfg is the market map configuration. This is the latest market map configuration that the
	// marketmap client has fetched from the destination endpoint.
	cfg mmtypes.MarketMap

	// version is the version of the market map configuration. This is used to determine if the
	// market map configuration schema has changed.
	version uint64

	// lastUpdated is the block height at which the market map configuration was last updated.
	lastUpdated int64

	// apiHandlerFactory is responsible for constructing the api handler for a given set of
	// providers.
	apiHandlerFactory base.APIHandlerFactory

	// wsHandlerFactory is responsible for constructing the websocket handler for a given set of
	// providers.
	wsHandlerFactory base.WebSocketHandlerFactory

	// providers is the set of providers that the marketmap client is responsible for. Specifically,
	// anytime an update is received from the oracle, the marketmap client will update the providers
	// with the new configuration.
	providers map[string]base.Provider[oracletypes.CurrencyPair, *big.Int]

	// doneCh is the channel that is used to signal the marketmap client to stop.
	doneCh chan struct{}
}

// Start starts the marketmap client. The client blocks until the context is cancelled. At a high level,
// the marketmap client fetches the latest market map configurations from the destination endpoint and
// updates the providers and aggregation configurations.
func (c *GRPCClient) Start(ctx context.Context) error {
	ticker := time.NewTicker(c.cfg.UpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.doneCh:
			return nil
		case <-ticker.C:
			// Fetch the latest market map configurations from the destination endpoint.
			// Update the providers and aggregation configurations with the new configurations.
			marketUpdate, err := c.fetchMarketMap()
			if err != nil {
				c.logger.Error("failed to fetch market map", zap.Error(err))
				continue
			}

			// Check if the market update is different from the current market update. If it is, update the
			// providers and aggregation configurations with the new market update.
			if !c.isMarketUpdateDifferent(marketUpdate) {
				c.logger.Debug("market update is the same as the current market update")
				continue
			}

			// Update the providers and aggregation configurations with the new market update.
			if err := c.updateProviders(marketUpdate); err != nil {
				c.logger.Error("failed to update providers", zap.Error(err))
			}

			// Update the aggregation configurations with the new market update.
			if err := c.updateAggregations(marketUpdate); err != nil {
				c.logger.Error("failed to update aggregations", zap.Error(err))
				continue
			}
		}
	}
}

// Stop stops the marketmap client.
func (c *GRPCClient) Stop() {
	c.doneCh <- struct{}{}
}
