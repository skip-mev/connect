package oracle

import (
	"context"
	"fmt"
	"math/big"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/base"
	mmclienttypes "github.com/skip-mev/connect/v2/service/clients/marketmap/types"
)

// Init initializes the all providers that are configured via the oracle config.
func (o *OracleImpl) Init(ctx context.Context) error {
	o.mut.Lock()
	defer o.mut.Unlock()

	for _, cfg := range o.cfg.Providers {
		// Initialize the provider.
		var err error
		switch cfg.Type {
		case types.ConfigType:
			err = o.createPriceProvider(ctx, cfg)
		case mmclienttypes.ConfigType:
			err = o.createMarketMapProvider(cfg)
		default:
			err = fmt.Errorf("unknown provider type: %s", cfg.Type)
		}

		if err != nil {
			o.logger.Error(
				"failed to initialize provider",
				zap.String("provider", cfg.Name),
				zap.Error(err),
			)

			return fmt.Errorf("failed to initialize %s provider: %w", cfg.Name, err)
		}
	}

	return nil
}

// createPriceProvider creates a new price provider for the given provider configuration.
func (o *OracleImpl) createPriceProvider(ctx context.Context, cfg config.ProviderConfig) error {
	// Create the provider market map. This creates the tickers the provider is configured to
	// support.
	tickers, err := types.ProviderTickersFromMarketMap(cfg.Name, o.marketMap)
	if err != nil {
		return fmt.Errorf("failed to create %s's provider market map: %w", cfg.Name, err)
	}

	// Select the query handler based on the provider's configuration.
	var provider *types.PriceProvider
	switch {
	case cfg.API.Enabled:
		queryHandler, err := o.createAPIQueryHandler(ctx, cfg)
		if err != nil {
			return fmt.Errorf("failed to create %s's api query handler: %w", cfg.Name, err)
		}

		provider, err = types.NewPriceProvider(
			base.WithName[types.ProviderTicker, *big.Float](cfg.Name),
			base.WithLogger[types.ProviderTicker, *big.Float](o.logger),
			base.WithAPIQueryHandler(queryHandler),
			base.WithAPIConfig[types.ProviderTicker, *big.Float](cfg.API),
			base.WithIDs[types.ProviderTicker, *big.Float](tickers),
			base.WithMetrics[types.ProviderTicker, *big.Float](o.providerMetrics),
		)
		if err != nil {
			return fmt.Errorf("failed to create %s's provider: %w", cfg.Name, err)
		}
	case cfg.WebSocket.Enabled:
		queryHandler, err := o.createWebSocketQueryHandler(ctx, cfg)
		if err != nil {
			return fmt.Errorf("failed to create %s's web socket query handler: %w", cfg.Name, err)
		}

		provider, err = types.NewPriceProvider(
			base.WithName[types.ProviderTicker, *big.Float](cfg.Name),
			base.WithLogger[types.ProviderTicker, *big.Float](o.logger),
			base.WithWebSocketQueryHandler(queryHandler),
			base.WithWebSocketConfig[types.ProviderTicker, *big.Float](cfg.WebSocket),
			base.WithIDs[types.ProviderTicker, *big.Float](tickers),
			base.WithMetrics[types.ProviderTicker, *big.Float](o.providerMetrics),
		)
		if err != nil {
			return fmt.Errorf("failed to create %s's provider: %w", cfg.Name, err)
		}
	default:
		return fmt.Errorf("provider %s has no enabled query handlers", cfg.Name)
	}

	state := ProviderState{
		Provider: provider,
		Cfg:      cfg,
	}

	// Add the provider to the oracle.
	o.priceProviders[provider.Name()] = state

	// Add the provider name to the message here since we want these to ignore log sampling limits
	o.logger.Info(
		fmt.Sprintf("created %s provider state", provider.Name()),
		zap.String("provider", provider.Name()),
		zap.Int("num_tickers", len(provider.GetIDs())),
	)
	return nil
}

// createAPIQueryHandler creates a new API query handler for the given provider configuration.
func (o *OracleImpl) createAPIQueryHandler(
	ctx context.Context,
	cfg config.ProviderConfig,
) (types.PriceAPIQueryHandler, error) {
	if o.priceAPIFactory == nil {
		return nil, fmt.Errorf("cannot create provider; api query handler factory is not set")
	}

	return o.priceAPIFactory(ctx, o.logger, cfg, o.apiMetrics)
}

// createWebSocketQueryHandler creates a new web socket query handler for the given provider configuration.
func (o *OracleImpl) createWebSocketQueryHandler(
	ctx context.Context,
	cfg config.ProviderConfig,
) (types.PriceWebSocketQueryHandler, error) {
	if o.priceWSFactory == nil {
		return nil, fmt.Errorf("cannot create provider; web socket query handler factory is not set")
	}

	return o.priceWSFactory(ctx, o.logger, cfg, o.wsMetrics)
}

// createMarketMapProvider creates a new market map provider for the given provider configuration.
func (o *OracleImpl) createMarketMapProvider(cfg config.ProviderConfig) error {
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
		return fmt.Errorf("failed to create market map provider (%s): %w", cfg.Name, err)
	}

	o.mmProvider = mapper
	o.logger.Info(
		"created market map provider",
		zap.String("provider", mapper.Name()),
	)
	return nil
}
