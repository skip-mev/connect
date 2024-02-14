package base

import (
	"context"

	apihandler "github.com/skip-mev/slinky/providers/base/api/handlers"
	wshandlers "github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"go.uber.org/zap"
)

// listenOnConfigUpdater listens for updates from the config updater and updates the
// provider's internal configurations. This will trigger the provider to restart
// and is blocking until the context is cancelled.
func (p *Provider[K, V]) listenOnConfigUpdater(ctx context.Context) {
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
		case apiHandler := <-p.updater.GetAPIHandler():
			p.logger.Debug("received new api handler")
			p.SetAPIHandler(apiHandler)
		case wsHandler := <-p.updater.GetWebSocketHandler():
			p.logger.Debug("received new websocket handler")
			p.SetWebSocketHandler(wsHandler)
		}

		p.restartCh <- struct{}{}
	}
}

// GetConfigUpdater returns the config updater for the provider.
func (p *Provider[K, V]) GetConfigUpdater() ConfigUpdater[K, V] {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.updater
}

// SetIDs sets the set of IDs that the provider is responsible for fetching data for.
func (p *Provider[K, V]) SetIDs(ids []K) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.ids = ids
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
func (p *Provider[K, V]) SetAPIHandler(apiHandler apihandler.APIQueryHandler[K, V]) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.Type() != providertypes.API {
		panic("cannot set api handler for non-api provider")
	}

	p.api = apiHandler
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
func (p *Provider[K, V]) SetWebSocketHandler(wsHandler wshandlers.WebSocketQueryHandler[K, V]) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.Type() != providertypes.WebSockets {
		panic("cannot set websocket handler for non-websocket provider")
	}

	p.ws = wsHandler
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
