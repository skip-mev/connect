package base

import (
	"context"

	"go.uber.org/zap"

	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

type (
	// The ConfigUpdater is an interface that be utilized to fetch the configuration
	// for a given provider. This can be used to update a provider's set of data
	// it is responsible for fetching.
	ConfigUpdater[K providertypes.ResponseKey] interface {
		// GetIDs is the channel that is used to update the set of IDs that the provider
		// will fetch data for. This blocks until there is a viable update.
		GetIDs() <-chan []K

		// UpdateIDs sets the set of IDs that the provider will fetch data for.
		UpdateIDs(ids []K)
	}

	// ConfigUpdaterImpl is a simple implementation of the ConfigClient interface. This
	// implementation inheritly blocks on all update operations.
	ConfigUpdaterImpl[K providertypes.ResponseKey] struct {
		// idsCh is the channel that is used to update the set of IDs that the provider
		// will fetch data for.
		idsCh chan []K
	}
)

var _ ConfigUpdater[oracletypes.CurrencyPair] = (*ConfigUpdaterImpl[oracletypes.CurrencyPair])(nil)

// NewConfigUpdater returns a new ConfigUpdaterImpl.
func NewConfigUpdater[K providertypes.ResponseKey]() *ConfigUpdaterImpl[K] {
	return &ConfigUpdaterImpl[K]{
		idsCh: make(chan []K),
	}
}

// GetIDs updates the set of IDs that the provider will fetch data for.
func (c *ConfigUpdaterImpl[K]) GetIDs() <-chan []K {
	return c.idsCh
}

// UpdateIDs sets the set of IDs that the provider will fetch data for.
func (c *ConfigUpdaterImpl[K]) UpdateIDs(ids []K) {
	c.idsCh <- ids
}

// listenConfigUpdater listens for updates from the config updater and updates
// the provider's internal configurations. This will trigger the provider to restart
// and is blocking until the context is cancelled.
func (p *Provider[K, V]) listenConfigUpdater(ctx context.Context) {
	if p.updater == nil {
		return
	}

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("stopping config client listener")
			return
		case ids := <-p.updater.GetIDs():
			p.logger.Debug("received new ids", zap.Any("ids", ids))
			p.SetIDs(ids)

			// Signal the provider to restart.
			p.restartCh <- struct{}{}
		}
	}
}
