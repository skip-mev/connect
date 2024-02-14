package oracle

import (
	"fmt"
	"math/big"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
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

	// apiFactory is the factory function that creates API query handlers.
	apiFactory factory.APIQueryHandlerFactory[slinkytypes.CurrencyPair, *big.Int]
	// wsFactory is the factory function that creates websocket query handlers.
	wsFactory factory.WebSocketQueryHandlerFactory[slinkytypes.CurrencyPair, *big.Int]
	// marketMap contains the entire set of price feeds that the oracle will fetch prices for.
	marketMap mmtypes.AggregateMarketConfig
}

// NewDefaultProviderFactory returns a new instance of the default provider factory.
func NewDefaultProviderFactory(
	logger *zap.Logger,
	apiFactory factory.APIQueryHandlerFactory[slinkytypes.CurrencyPair, *big.Int],
	wsFactory factory.WebSocketQueryHandlerFactory[slinkytypes.CurrencyPair, *big.Int],
	marketmap mmtypes.AggregateMarketConfig,
) (*DefaultOracleProviderFactory, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	if apiFactory == nil {
		return nil, fmt.Errorf("apiFactory cannot be nil")
	}

	if wsFactory == nil {
		return nil, fmt.Errorf("wsFactory cannot be nil")
	}

	if err := marketmap.ValidateBasic(); err != nil {
		return nil, err
	}

	return &DefaultOracleProviderFactory{
		logger:     logger,
		apiFactory: apiFactory,
		wsFactory:  wsFactory,
		marketMap:  marketmap,
	}, nil
}

// Factory returns a factory function that creates providers based on the oracle configuration.
func (f *DefaultOracleProviderFactory) Factory() factory.ProviderFactory[slinkytypes.CurrencyPair, *big.Int] {
	return func(cfg config.OracleConfig) ([]providertypes.Provider[slinkytypes.CurrencyPair, *big.Int], error) {
		if err := cfg.ValidateBasic(); err != nil {
			return nil, err
		}

		// Create the metrics that are used by the providers.
		wsMetrics := wsmetrics.NewWebSocketMetricsFromConfig(cfg.Metrics)
		apiMetrics := apimetrics.NewAPIMetricsFromConfig(cfg.Metrics)
		providerMetrics := providermetrics.NewProviderMetricsFromConfig(cfg.Metrics)

		// Create the providers.
		providers := make([]providertypes.Provider[slinkytypes.CurrencyPair, *big.Int], len(cfg.Providers))
		for i, p := range cfg.Providers {
			// Get the market configuration for the provider.
			market, ok := f.marketMap.MarketConfigs[p.Name]
			if !ok {
				f.logger.Info("market config not found", zap.String("provider", p.Name))
				continue
			}

			switch {
			case p.API.Enabled:
				queryHandler, err := f.apiFactory(f.logger, p, apiMetrics)
				if err != nil {
					return nil, err
				}

				// Create the provider.
				provider, err := base.NewProvider[slinkytypes.CurrencyPair, *big.Int](
					base.WithName[slinkytypes.CurrencyPair, *big.Int](p.Name),
					base.WithLogger[slinkytypes.CurrencyPair, *big.Int](f.logger),
					base.WithAPIQueryHandler(queryHandler),
					base.WithAPIConfig[slinkytypes.CurrencyPair, *big.Int](p.API),
					base.WithIDs[slinkytypes.CurrencyPair, *big.Int](cfg.Market.GetCurrencyPairs()),
					base.WithMetrics[slinkytypes.CurrencyPair, *big.Int](providerMetrics),
				)
				if err != nil {
					return nil, err
				}

				providers[i] = provider
			case p.WebSocket.Enabled:
				// Create the websocket query handler which encapsulates all fetching and parsing logic.
				queryHandler, err := f.wsFactory(f.logger, p, wsMetrics)
				if err != nil {
					return nil, err
				}

				// Create the provider.
				provider, err := base.NewProvider[slinkytypes.CurrencyPair, *big.Int](
					base.WithName[slinkytypes.CurrencyPair, *big.Int](p.Name),
					base.WithLogger[slinkytypes.CurrencyPair, *big.Int](f.logger),
					base.WithWebSocketQueryHandler(queryHandler),
					base.WithWebSocketConfig[slinkytypes.CurrencyPair, *big.Int](p.WebSocket),
					base.WithIDs[slinkytypes.CurrencyPair, *big.Int](cfg.Market.GetCurrencyPairs()),
					base.WithMetrics[slinkytypes.CurrencyPair, *big.Int](providerMetrics),
				)
				if err != nil {
					return nil, err
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
