package orchestrator

import (
	"fmt"
	"math/big"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	oracletypes "github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base"
	mmclienttypes "github.com/skip-mev/slinky/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	"go.uber.org/zap"
)

// Init initializes the all of the providers that are configured via the oracle config.
// Specifically, this will:
//
// 1. This will initialize the provider.
// 2. Create the provider specific market map, if configured with a marketmap.
// 3. Enable the provider if the provider is included in the oracle config and marketmap.
func (o *ProviderOrchestrator) Init() error {
	o.mut.Lock()
	defer o.mut.Unlock()

	for _, providerCfg := range o.cfg.Providers {
		// Initialize the provider.
		if err := o.createProvider(providerCfg); err != nil {
			o.logger.Error("failed to create provider state", zap.Error(err))
			return err
		}
	}

	return nil
}

// CreateProviderState creates a provider state for the given provider. This constructs the
// query handler, based on the provider's type and configuration. The provider state is then
// enabled/disabled based on whether the provider is configured to support any of the tickers.
func (o *ProviderOrchestrator) createProvider(
	cfg config.ProviderConfig,
) error {
	switch cfg.Type {
	case oracletypes.ConfigType:
		return o.createPriceProvider(cfg)
	case mmclienttypes.ConfigType:
		return o.createMarketMapProvider(cfg)
	default:
		return fmt.Errorf("unknown provider type: %s", cfg.Type)
	}
}

func (o *ProviderOrchestrator) createPriceProvider(cfg config.ProviderConfig) error {
	// Create the provider market map. This creates the tickers the provider is configured to
	// support.
	market, err := types.ProviderMarketMapFromMarketMap(cfg.Name, o.marketMap)
	if err != nil {
		return fmt.Errorf("failed to create %s's provider market map: %w", cfg.Name, err)
	}

	// Select the query handler based on the provider's configuration.
	var provider *types.PriceProvider
	switch {
	case cfg.API.Enabled:
		queryHandler, err := o.createAPIQueryHandler(cfg, market)
		if err != nil {
			return fmt.Errorf("failed to create %s's api query handler: %w", cfg.Name, err)
		}

		provider, err = types.NewPriceProvider(
			base.WithName[mmtypes.Ticker, *big.Int](cfg.Name),
			base.WithLogger[mmtypes.Ticker, *big.Int](o.logger),
			base.WithAPIQueryHandler(queryHandler),
			base.WithAPIConfig[mmtypes.Ticker, *big.Int](cfg.API),
			base.WithIDs[mmtypes.Ticker, *big.Int](market.GetTickers()),
			base.WithMetrics[mmtypes.Ticker, *big.Int](o.providerMetrics),
		)
		if err != nil {
			return fmt.Errorf("failed to create %s's provider: %w", cfg.Name, err)
		}
	case cfg.WebSocket.Enabled:
		queryHandler, err := o.createWebSocketQueryHandler(cfg, market)
		if err != nil {
			return fmt.Errorf("failed to create %s's web socket query handler: %w", cfg.Name, err)
		}

		provider, err = types.NewPriceProvider(
			base.WithName[mmtypes.Ticker, *big.Int](cfg.Name),
			base.WithLogger[mmtypes.Ticker, *big.Int](o.logger),
			base.WithWebSocketQueryHandler(queryHandler),
			base.WithWebSocketConfig[mmtypes.Ticker, *big.Int](cfg.WebSocket),
			base.WithIDs[mmtypes.Ticker, *big.Int](market.GetTickers()),
			base.WithMetrics[mmtypes.Ticker, *big.Int](o.providerMetrics),
		)
		if err != nil {
			return fmt.Errorf("failed to create %s's provider: %w", cfg.Name, err)
		}
	default:
		return fmt.Errorf("provider %s has no enabled query handlers", cfg.Name)
	}

	state := ProviderState{
		Provider: provider,
		Market:   market,
		Enabled:  len(market.GetTickers()) > 0,
		Cfg:      cfg,
	}

	// Add the provider to the orchestrator.
	o.providers[cfg.Name] = state

	o.logger.Info(
		"created price provider state",
		zap.String("provider", cfg.Name),
		zap.Bool("enabled", state.Enabled),
		zap.Int("num_tickers", len(state.Market.GetTickers())),
	)
	return nil
}

// createAPIQueryHandler creates a new API query handler for the given provider configuration.
func (o *ProviderOrchestrator) createAPIQueryHandler(
	cfg config.ProviderConfig,
	market types.ProviderMarketMap,
) (types.PriceAPIQueryHandler, error) {
	if o.priceAPIFactory == nil {
		return nil, fmt.Errorf("cannot create provider; api query handler factory is not set")
	}

	return o.priceAPIFactory(o.logger, cfg, o.apiMetrics, market)
}

// createWebSocketQueryHandler creates a new web socket query handler for the given provider configuration.
func (o *ProviderOrchestrator) createWebSocketQueryHandler(
	cfg config.ProviderConfig,
	market types.ProviderMarketMap,
) (types.PriceWebSocketQueryHandler, error) {
	if o.priceWSFactory == nil {
		return nil, fmt.Errorf("cannot create provider; web socket query handler factory is not set")
	}

	return o.priceWSFactory(o.logger, cfg, o.wsMetrics, market)
}

// CreateMarketMapProvider creates a new market map provider for the given provider configuration.
func (o *ProviderOrchestrator) createMarketMapProvider(cfg config.ProviderConfig) error {
	if o.marketMapperFactory == nil {
		return fmt.Errorf("cannot create market map provider; market map factory is not set")
	}

	mapper, err := o.marketMapperFactory(
		o.logger,
		o.providerMetrics,
		o.apiMetrics,
		cfg,
	)
	if err != nil {
		return fmt.Errorf("failed to create market map provider: %w", err)
	}

	state := MapperState{
		Mapper:   mapper,
		Interval: cfg.API.Interval,
	}

	o.mapper = state
	o.logger.Info(
		"created market map provider",
		zap.String("provider", cfg.Name),
		zap.Duration("interval", state.Interval),
	)
	return nil
}
