package orchestrator

import (
	"context"
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
