package orchestrator

import (
	"context"
	"fmt"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	mmclienttypes "github.com/skip-mev/slinky/service/clients/marketmap/types"
)

type (
	// ProviderState is the state of a provider. This includes the provider implementation,
	// the provider specific market map, and whether the provider is enabled.
	ProviderState struct {
		// Provider is the price provider implementation.
		Provider *types.PriceProvider
		// Market is the market map view for the provider.
		Market types.ProviderMarketMap
		// Enabled is a flag that indicates whether the provider is enabled. A provider
		// is enabled iff it is configured with a market map and the market map has
		// tickers.
		Enabled bool
		// Cfg is the provider configuration.
		Cfg config.ProviderConfig
	}

	// MapperState is the state of the market map provider.
	MapperState struct {
		// Mapper is the market map provider implementation.
		Mapper mmclienttypes.MarketMapProvider
		// Interval is the interval at which the orchestrator will check for the latest
		// market map data.
		Interval time.Duration
	}

	// GeneralProvider is a interface for the general provider.
	GeneralProvider interface {
		// Start starts the provider.
		Start(ctx context.Context) error
		// Stop stops the provider.
		Name() string
	}
)

// setMainCtx sets the main context for the provider orchestrator.
func (o *ProviderOrchestrator) setMainCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	o.mut.Lock()
	defer o.mut.Unlock()

	o.mainCtx, o.mainCancel = context.WithCancel(ctx)
	return o.mainCtx, o.mainCancel
}

// getMainCtx returns the main context for the provider orchestrator.
func (o *ProviderOrchestrator) getMainCtx() (context.Context, context.CancelFunc) {
	o.mut.Lock()
	defer o.mut.Unlock()

	return o.mainCtx, o.mainCancel
}

// createAPIQueryHandler creates a new API query handler for the given provider configuration.
func (o *ProviderOrchestrator) createAPIQueryHandler(
	cfg config.ProviderConfig,
	market types.ProviderMarketMap,
) (types.PriceAPIQueryHandler, error) {
	if o.apiQueryHandlerFactory == nil {
		return nil, fmt.Errorf("cannot create provider; api query handler factory is not set")
	}

	return o.apiQueryHandlerFactory(o.logger, cfg, o.apiMetrics, market)
}

// createWebSocketQueryHandler creates a new web socket query handler for the given provider configuration.
func (o *ProviderOrchestrator) createWebSocketQueryHandler(
	cfg config.ProviderConfig,
	market types.ProviderMarketMap,
) (types.PriceWebSocketQueryHandler, error) {
	if o.webSocketQueryHandlerFactory == nil {
		return nil, fmt.Errorf("cannot create provider; web socket query handler factory is not set")
	}

	return o.webSocketQueryHandlerFactory(o.logger, cfg, o.wsMetrics, market)
}
