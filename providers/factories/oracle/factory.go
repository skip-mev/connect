package oracle

import (
	"fmt"
	"math/big"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base"
	apimetrics "github.com/skip-mev/slinky/providers/base/api/metrics"
	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"
	wsmetrics "github.com/skip-mev/slinky/providers/base/websocket/metrics"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/types/factory"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// DefaultOracleProviderFactory is a sample implementation of the provider factory. This provider
// factory function returns providers that are API & websocket based.
type DefaultOracleProviderFactory struct {
	logger *zap.Logger

	// marketMap is the market map that is used to configure the providers.
	marketMap mmtypes.MarketMap
}

// NewDefaultProviderFactory returns a new instance of the default provider factory.
func NewDefaultProviderFactory(
	logger *zap.Logger,
	marketmap mmtypes.MarketMap,
) (*DefaultOracleProviderFactory, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	if err := marketmap.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid market map: %w", err)
	}

	return &DefaultOracleProviderFactory{
		logger:    logger,
		marketMap: marketmap,
	}, nil
}

// Factory returns a factory function that creates providers based on the oracle configuration.
func (f *DefaultOracleProviderFactory) Factory() factory.ProviderFactory[mmtypes.Ticker, *big.Int] {
	return func(cfg config.OracleConfig) ([]providertypes.Provider[mmtypes.Ticker, *big.Int], error) {
		if err := cfg.ValidateBasic(); err != nil {
			return nil, err
		}

		// Create the metrics that are used by the providers.
		wsMetrics := wsmetrics.NewWebSocketMetricsFromConfig(cfg.Metrics)
		apiMetrics := apimetrics.NewAPIMetricsFromConfig(cfg.Metrics)
		providerMetrics := providermetrics.NewProviderMetricsFromConfig(cfg.Metrics)

		// Create the providers.
		providers := make([]providertypes.Provider[mmtypes.Ticker, *big.Int], len(cfg.Providers))
		for i, p := range cfg.Providers {
			// Create the providers relative market map.
			providerMarketMap, err := types.ProviderMarketMapFromMarketMap(p.Name, f.marketMap)
			if err != nil {
				return nil, fmt.Errorf("failed to create %s's provider market map: %w", p.Name, err)
			}

			switch {
			case p.API.Enabled:
				queryHandler, err := APIQueryHandlerFactory(f.logger, p, apiMetrics, providerMarketMap)
				if err != nil {
					return nil, fmt.Errorf("failed to create %s's API query handler: %w", p.Name, err)
				}

				// Create the provider.
				provider, err := types.NewPriceProvider(
					base.WithName[mmtypes.Ticker, *big.Int](p.Name),
					base.WithLogger[mmtypes.Ticker, *big.Int](f.logger),
					base.WithAPIQueryHandler(queryHandler),
					base.WithAPIConfig[mmtypes.Ticker, *big.Int](p.API),
					base.WithIDs[mmtypes.Ticker, *big.Int](providerMarketMap.GetTickers()),
					base.WithMetrics[mmtypes.Ticker, *big.Int](providerMetrics),
				)
				if err != nil {
					return nil, fmt.Errorf("failed to create %s's provider: %w", p.Name, err)
				}

				providers[i] = provider
			case p.WebSocket.Enabled:
				queryHandler, err := WebSocketQueryHandlerFactory(f.logger, p, wsMetrics, providerMarketMap)
				if err != nil {
					return nil, fmt.Errorf("failed to create %s's web socket query handler: %w", p.Name, err)
				}

				// Create the provider.
				provider, err := base.NewProvider[mmtypes.Ticker, *big.Int](
					base.WithName[mmtypes.Ticker, *big.Int](p.Name),
					base.WithLogger[mmtypes.Ticker, *big.Int](f.logger),
					base.WithWebSocketQueryHandler(queryHandler),
					base.WithWebSocketConfig[mmtypes.Ticker, *big.Int](p.WebSocket),
					base.WithIDs[mmtypes.Ticker, *big.Int](providerMarketMap.GetTickers()),
					base.WithMetrics[mmtypes.Ticker, *big.Int](providerMetrics),
				)
				if err != nil {
					return nil, fmt.Errorf("failed to create %s's provider: %w", p.Name, err)
				}

				providers[i] = provider
			default:
				f.logger.Info("unknown provider type", zap.String("provider", p.Name))
				return nil, fmt.Errorf("unknown provider type: %s", p.Name)
			}
		}

		return providers, nil
	}
}
