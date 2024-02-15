package base

import (
	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
	wshandlers "github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

type (
	// The ConfigUpdater is an interface that be utilized to fetch the configuration
	// for a given provider. This allows the provider to be updated asynchronously.
	ConfigUpdater[K providertypes.ResponseKey, V providertypes.ResponseValue] interface {
		// GetIDs is the channel that is used to update the set of IDs that the provider
		// will fetch data for. This blocks until there is a viable update.
		GetIDs() <-chan []K

		// UpdateIDs sets the set of IDs that the provider will fetch data for.
		UpdateIDs(ids []K)

		// GetAPIHandler is the channel that is used to update the API handler that the provider
		// will use to fetch data. This blocks until there is a viable update.
		GetAPIHandler() <-chan apihandlers.APIQueryHandler[K, V]

		// UpdateAPIHandler sets the API handler that the provider will use to fetch data.
		UpdateAPIHandler(apiHandler apihandlers.APIQueryHandler[K, V])

		// GetWSHandler is the channel that is used to update the WebSocket handler that the provider
		// will use to fetch data. This blocks until there is a viable update.
		GetWebSocketHandler() <-chan wshandlers.WebSocketQueryHandler[K, V]

		// UpdateWSHandler sets the WebSocket handler that the provider will use to fetch data.
		UpdateWebSocketHandler(wsHandler wshandlers.WebSocketQueryHandler[K, V])
	}

	// ConfigUpdaterImpl is a simple implementation of the ConfigUpdater interface. This
	// implementation blocks on all receive operations. All of the channels are buffered
	// with a size of 1 to ensure that the provider can update the configuration without
	// blocking.
	ConfigUpdaterImpl[K providertypes.ResponseKey, V providertypes.ResponseValue] struct {
		// idsCh is the channel that is used to update the set of IDs that the provider
		// will fetch data for.
		idsCh chan []K

		// apiHandlerCh is the channel that is used to update the API handler that the provider
		// will use to fetch data.
		apiHandlerCh chan apihandlers.APIQueryHandler[K, V]

		// wsHandlerCh is the channel that is used to update the WebSocket handler that the provider
		// will use to fetch data.
		wsHandlerCh chan wshandlers.WebSocketQueryHandler[K, V]
	}
)

// NewConfigUpdater returns a new ConfigUpdaterImpl.
func NewConfigUpdater[K providertypes.ResponseKey, V providertypes.ResponseValue]() ConfigUpdater[K, V] {
	return &ConfigUpdaterImpl[K, V]{
		idsCh:        make(chan []K, 1),
		apiHandlerCh: make(chan apihandlers.APIQueryHandler[K, V], 1),
		wsHandlerCh:  make(chan wshandlers.WebSocketQueryHandler[K, V], 1),
	}
}

// GetIDs updates the set of IDs that the provider will fetch data for.
func (c *ConfigUpdaterImpl[K, V]) GetIDs() <-chan []K {
	return c.idsCh
}

// UpdateIDs sets the set of IDs that the provider will fetch data for.
func (c *ConfigUpdaterImpl[K, V]) UpdateIDs(ids []K) {
	c.idsCh <- ids
}

// GetAPIHandler updates the API handler that the provider will use to fetch data.
func (c *ConfigUpdaterImpl[K, V]) GetAPIHandler() <-chan apihandlers.APIQueryHandler[K, V] {
	return c.apiHandlerCh
}

// UpdateAPIHandler sets the API handler that the provider will use to fetch data.
func (c *ConfigUpdaterImpl[K, V]) UpdateAPIHandler(apiHandler apihandlers.APIQueryHandler[K, V]) {
	c.apiHandlerCh <- apiHandler
}

// GetWebSocketHandler updates the WebSocket handler that the provider will use to fetch data.
func (c *ConfigUpdaterImpl[K, V]) GetWebSocketHandler() <-chan wshandlers.WebSocketQueryHandler[K, V] {
	return c.wsHandlerCh
}

// UpdateWebSocketHandler sets the WebSocket handler that the provider will use to fetch data.
func (c *ConfigUpdaterImpl[K, V]) UpdateWebSocketHandler(wsHandler wshandlers.WebSocketQueryHandler[K, V]) {
	c.wsHandlerCh <- wsHandler
}
