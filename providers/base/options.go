package base

import (
	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	apihandlers "github.com/skip-mev/connect/v2/providers/base/api/handlers"
	providermetrics "github.com/skip-mev/connect/v2/providers/base/metrics"
	wshandlers "github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
)

// ProviderOption is a function that can be used to modify a provider.
type ProviderOption[K providertypes.ResponseKey, V providertypes.ResponseValue] func(*Provider[K, V])

// WithName sets the name of the provider.
func WithName[K providertypes.ResponseKey, V providertypes.ResponseValue](name string) ProviderOption[K, V] {
	return func(p *Provider[K, V]) {
		if name == "" {
			panic("cannot set empty name")
		}

		p.name = name
	}
}

// WithLogger sets the logger for the provider.
func WithLogger[K providertypes.ResponseKey, V providertypes.ResponseValue](logger *zap.Logger) ProviderOption[K, V] {
	return func(p *Provider[K, V]) {
		if logger == nil {
			panic("cannot set nil logger")
		}

		p.logger = logger.With(zap.String("provider", p.name))
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

// WithAPIConfig sets the APIConfig for the provider.
func WithAPIConfig[K providertypes.ResponseKey, V providertypes.ResponseValue](cfg config.APIConfig) ProviderOption[K, V] {
	return func(p *Provider[K, V]) {
		if cfg.ValidateBasic() != nil {
			panic("invalid api config")
		}

		p.apiCfg = cfg
	}
}

// WithWebSocketQueryHandler sets the WebSocketQueryHandler for the provider. If your provider
// utilizes a websocket based provider, you should use this option to set the WebSocketQueryHandler.
func WithWebSocketQueryHandler[K providertypes.ResponseKey, V providertypes.ResponseValue](ws wshandlers.WebSocketQueryHandler[K, V]) ProviderOption[K, V] {
	return func(p *Provider[K, V]) {
		if p.ws != nil {
			panic("cannot set websocket query handler twice")
		}

		if ws == nil {
			panic("cannot set nil websocket query handler")
		}

		p.ws = ws
	}
}

// WithWebSocketConfig sets the WebSocketConfig for the provider.
func WithWebSocketConfig[K providertypes.ResponseKey, V providertypes.ResponseValue](cfg config.WebSocketConfig) ProviderOption[K, V] {
	return func(p *Provider[K, V]) {
		if cfg.ValidateBasic() != nil {
			panic("invalid websocket config")
		}

		p.wsCfg = cfg
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
