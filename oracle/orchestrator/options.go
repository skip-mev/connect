package orchestrator

import (
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/types"
	mmclienttypes "github.com/skip-mev/slinky/service/clients/marketmap/types"
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

		m.logger = logger.With(zap.String("process", "provider_orchestrator"))
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

// WithPriceAPIQueryHandlerFactory sets the Price API query handler factory for the provider orchestrator.
// Specifically this is what is utilized to construct price providers that are API based.
func WithPriceAPIQueryHandlerFactory(factory types.PriceAPIQueryHandlerFactory) Option {
	return func(m *ProviderOrchestrator) {
		if factory == nil {
			panic("api query handler factory cannot be nil")
		}

		m.priceAPIFactory = factory
	}
}

// WithWebSocketQueryHandlerFactory sets the websocket query handler factory for the provider orchestrator.
// Specifically this is what is utilized to construct price providers that are websocket based.
func WithPriceWebSocketQueryHandlerFactory(factory types.PriceWebSocketQueryHandlerFactory) Option {
	return func(m *ProviderOrchestrator) {
		if factory == nil {
			panic("websocket query handler factory cannot be nil")
		}

		m.priceWSFactory = factory
	}
}

// WithMarketMapperFactory sets the market map factory for the provider orchestrator.
// Specifically this is what is utilized to construct market map providers.
func WithMarketMapperFactory(factory mmclienttypes.MarketMapFactory) Option {
	return func(m *ProviderOrchestrator) {
		if factory == nil {
			panic("market map factory cannot be nil")
		}

		m.marketMapperFactory = factory
	}
}
