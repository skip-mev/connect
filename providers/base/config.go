package base

import (
	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	apihandler "github.com/skip-mev/connect/v2/providers/base/api/handlers"
	wshandlers "github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
)

// UpdateOption are the options that can be used to update the provider.
type UpdateOption[K providertypes.ResponseKey, V providertypes.ResponseValue] func(*Provider[K, V])

// WithNewIDs returns an option that sets the new set of IDs that the provider is responsible for fetching data for.
func WithNewIDs[K providertypes.ResponseKey, V providertypes.ResponseValue](ids []K) UpdateOption[K, V] {
	return func(p *Provider[K, V]) {
		p.setIDs(ids)
	}
}

// WithNewAPIHandler returns an option that sets the new API handler that the provider will use to fetch data.
func WithNewAPIHandler[K providertypes.ResponseKey, V providertypes.ResponseValue](
	apiHandler apihandler.APIQueryHandler[K, V],
) UpdateOption[K, V] {
	return func(p *Provider[K, V]) {
		p.setAPIHandler(apiHandler)
	}
}

// WithNewWebSocketHandler returns an option that sets the new WebSocket handler that the provider will use to fetch data.
func WithNewWebSocketHandler[K providertypes.ResponseKey, V providertypes.ResponseValue](
	wsHandler wshandlers.WebSocketQueryHandler[K, V],
) UpdateOption[K, V] {
	return func(p *Provider[K, V]) {
		p.setWebSocketHandler(wsHandler)
	}
}

// Update updates the provider with the given options.
func (p *Provider[K, V]) Update(opts ...UpdateOption[K, V]) {
	p.logger.Debug("updating provider")
	for _, opt := range opts {
		opt(p)
	}
	p.logger.Debug("provider updated")

	if _, cancel := p.getFetchCtx(); cancel != nil {
		p.logger.Debug("canceling fetch context; restarting provider")
		cancel()
	}
}

// SetIDs sets the set of IDs that the provider is responsible for fetching data for.
func (p *Provider[K, V]) setIDs(ids []K) {
	p.mu.Lock()
	p.ids = ids
	p.mu.Unlock()

	p.logger.Debug("set ids", zap.Any("ids", ids))
}

// GetIDs returns the set of IDs that the provider is responsible for fetching data for.
func (p *Provider[K, V]) GetIDs() []K {
	p.mu.Lock()
	defer p.mu.Unlock()

	ids := make([]K, len(p.ids))
	copy(ids, p.ids)

	return ids
}

// SetAPIHandler sets the API handler that the provider will use to fetch data.
func (p *Provider[K, V]) setAPIHandler(apiHandler apihandler.APIQueryHandler[K, V]) {
	if p.Type() != providertypes.API {
		panic("cannot set api handler for non-api provider")
	}

	p.mu.Lock()
	p.api = apiHandler
	p.mu.Unlock()

	p.logger.Debug("set api query handler")
}

// GetAPIHandler returns the API handler that the provider will use to fetch data.
func (p *Provider[K, V]) GetAPIHandler() apihandler.APIQueryHandler[K, V] {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.Type() != providertypes.API {
		panic("cannot get api handler for non-api provider")
	}

	return p.api
}

// SetWebSocketHandler sets the WebSocket handler that the provider will use to fetch data.
func (p *Provider[K, V]) setWebSocketHandler(wsHandler wshandlers.WebSocketQueryHandler[K, V]) {
	if p.Type() != providertypes.WebSockets {
		panic("cannot set websocket handler for non-websocket provider")
	}

	p.mu.Lock()
	p.ws = wsHandler
	p.mu.Unlock()

	p.logger.Debug("set websocket query handler")
}

// GetWebSocketHandler returns the WebSocket handler that the provider will use to fetch data.
func (p *Provider[K, V]) GetWebSocketHandler() wshandlers.WebSocketQueryHandler[K, V] {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.Type() != providertypes.WebSockets {
		panic("cannot get websocket handler for non-websocket provider")
	}

	return p.ws
}

// GetAPIConfig returns the API configuration for the provider.
func (p *Provider[K, V]) GetAPIConfig() config.APIConfig {
	return p.apiCfg
}
