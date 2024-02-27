package orchestrator

import (
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// Option is a functional option for the market map state.
type Option func(*ProviderOrchestrator)

// WithLogger sets the logger for the provider orchestrator.
func WithLogger(logger *zap.Logger) Option {
	return func(m *ProviderOrchestrator) {
		if logger == nil {
			panic("logger cannot be nil")
		}

		m.logger = logger.With(zap.String("process", "provider orchestrator"))
	}
}

// WithMarketMap sets the market map for the provider orchestrator.
func WithMarketMap(marketMap mmtypes.MarketMap) Option {
	return func(m *ProviderOrchestrator) {
		if err := marketMap.ValidateBasic(); err != nil {
			panic(err)
		}

		m.marketMap = marketMap
	}
}

// WithAPIQueryHandlerFactory sets the API query handler factory for the provider orchestrator.
func WithAPIQueryHandlerFactory(factory types.PriceAPIQueryHandlerFactory) Option {
	return func(m *ProviderOrchestrator) {
		if factory == nil {
			panic("api query handler factory cannot be nil")
		}

		m.apiQueryHandlerFactory = factory
	}
}

// WithWebSocketQueryHandlerFactory sets the websocket query handler factory for the provider orchestrator.
func WithWebSocketQueryHandlerFactory(factory types.PriceWebSocketQueryHandlerFactory) Option {
	return func(m *ProviderOrchestrator) {
		if factory == nil {
			panic("websocket query handler factory cannot be nil")
		}

		m.webSocketQueryHandlerFactory = factory
	}
}
