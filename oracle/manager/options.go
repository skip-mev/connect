package manager

import (
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// Option is a functional option for the market map state.
type Option func(*ProviderManager)

// WithLogger sets the logger for the provider manager.
func WithLogger(logger *zap.Logger) Option {
	return func(m *ProviderManager) {
		if logger == nil {
			panic("logger cannot be nil")
		}

		m.logger = logger.With(zap.String("process", "provider manager"))
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
