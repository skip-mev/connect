package base

import (
	providertypes "github.com/skip-mev/slinky/providers/types"
)

type (
	// The ConfigUpdater is an interface that be utilized to fetch the configuration
	// for a given provider. This allows the provider to be updated asynchronously.
	ConfigUpdater[K providertypes.ResponseKey] interface {
		// GetIDs is the channel that is used to update the set of IDs that the provider
		// will fetch data for. This blocks until there is a viable update.
		GetIDs() <-chan []K

		// UpdateIDs sets the set of IDs that the provider will fetch data for.
		UpdateIDs(ids []K)
	}

	// ConfigUpdaterImpl is a simple implementation of the ConfigUpdater interface. This
	// implementation blocks on all receive operations.
	ConfigUpdaterImpl[K providertypes.ResponseKey] struct {
		// idsCh is the channel that is used to update the set of IDs that the provider
		// will fetch data for.
		idsCh chan []K
	}
)

// NewConfigUpdater returns a new ConfigUpdaterImpl.
func NewConfigUpdater[K providertypes.ResponseKey]() ConfigUpdater[K] {
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
