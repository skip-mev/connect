package base

import (
	"go.uber.org/zap"

	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"
	wshandlers "github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

// ProviderOption is a function that can be used to modify a provider.
type ProviderOption[K providertypes.ResponseKey, V providertypes.ResponseValue] func(*Provider[K, V])

// WithLogger sets the logger for the provider.
func WithLogger[K providertypes.ResponseKey, V providertypes.ResponseValue](logger *zap.Logger) ProviderOption[K, V] {
	return func(p *Provider[K, V]) {
		if logger == nil {
			panic("cannot set nil logger")
		}

		p.logger = logger.With(zap.String("provider", p.cfg.Name))
	}
}

// WithAPIQueryHandler sets the APIQueryHandler for the provider. If your provider utilizes a
// API (HTTP) based provider, you should use this option to set the APIQueryHandler.
func WithAPIQueryHandler[K providertypes.ResponseKey, V providertypes.ResponseValue](api apihandlers.APIQueryHandler[K, V]) ProviderOption[K, V] {
	return func(p *Provider[K, V]) {
		if p.api != nil {
			panic("cannot set api query handler twice")
		}

		if api == nil {
			panic("cannot set nil api query handler")
		}

		p.api = api
	}
}

// WithWebSocketQueryHandler sets the WebSocketQueryHandler for the provider. If your provider
// utilizes a websocket based provider, you should use this option to set the WebSocketQueryHandler.
func WithWebSocketQueryHandler[K providertypes.ResponseKey, V providertypes.ResponseValue](ws wshandlers.WebSocketQueryHandler[K, V]) ProviderOption[K, V] {
	return func(p *Provider[K, V]) {
		if p.ws != nil {
			panic("cannot set web socket query handler twice")
		}

		if ws == nil {
			panic("cannot set nil web socket query handler")
		}

		p.ws = ws
	}
}

// WithIDs sets the IDs that the provider is responsible for fetching data for.
func WithIDs[K providertypes.ResponseKey, V providertypes.ResponseValue](ids []K) ProviderOption[K, V] {
	return func(p *Provider[K, V]) {
		if ids == nil {
			panic("cannot set nil ids")
		}

		p.ids = ids
	}
}

// WithMetrics sets the metrics implementation for the provider.
func WithMetrics[K providertypes.ResponseKey, V providertypes.ResponseValue](metrics providermetrics.ProviderMetrics) ProviderOption[K, V] {
	return func(p *Provider[K, V]) {
		if metrics == nil {
			panic("cannot set nil metrics")
		}

		p.metrics = metrics
	}
}
