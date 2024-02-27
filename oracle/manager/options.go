package manager

import (
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/types"
	apimetrics "github.com/skip-mev/slinky/providers/base/api/metrics"
	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"
	wsmetrics "github.com/skip-mev/slinky/providers/base/websocket/metrics"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// WithLogger sets the logger for the provider manager.
func WithLogger(logger *zap.Logger) Option {
	return func(m *ProviderManager) {
		if logger == nil {
			panic("logger cannot be nil")
		}

		m.logger = logger
	}
}

// WithMarketMap sets the market map for the provider manager.
func WithMarketMap(marketMap mmtypes.MarketMap) Option {
	return func(m *ProviderManager) {
		if err := marketMap.ValidateBasic(); err != nil {
			panic(err)
		}

		m.marketMap = marketMap
	}
}

// WithAPIQueryHandlerFactory sets the API query handler factory for the provider manager.
func WithAPIQueryHandlerFactory(factory types.PriceAPIQueryHandlerFactory) Option {
	return func(m *ProviderManager) {
		if factory == nil {
			panic("api query handler factory cannot be nil")
		}

		m.apiQueryHandlerFactory = factory
	}
}

// WithWebSocketQueryHandlerFactory sets the websocket query handler factory for the provider manager.
func WithWebSocketQueryHandlerFactory(factory types.PriceWebSocketQueryHandlerFactory) Option {
	return func(m *ProviderManager) {
		if factory == nil {
			panic("websocket query handler factory cannot be nil")
		}

		m.webSocketQueryHandlerFactory = factory
	}
}

// WithWebSocketMetrics sets the websocket metrics for the provider manager.
func WithWebSocketMetrics(metrics wsmetrics.WebSocketMetrics) Option {
	return func(m *ProviderManager) {
		if metrics == nil {
			panic("websocket metrics cannot be nil")
		}

		m.wsMetrics = metrics
	}
}

// WithAPIMetrics sets the API metrics for the provider manager.
func WithAPIMetrics(metrics apimetrics.APIMetrics) Option {
	return func(m *ProviderManager) {
		if metrics == nil {
			panic("api metrics cannot be nil")
		}

		m.apiMetrics = metrics
	}
}

// WithProviderMetrics sets the provider metrics for the provider manager.
func WithProviderMetrics(metrics providermetrics.ProviderMetrics) Option {
	return func(m *ProviderManager) {
		if metrics == nil {
			panic("provider metrics cannot be nil")
		}

		m.providerMetrics = metrics
	}
}
