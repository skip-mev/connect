package oracle

import (
	"fmt"
	"math/big"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base"
	apimetrics "github.com/skip-mev/slinky/providers/base/api/metrics"
	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"
	wsmetrics "github.com/skip-mev/slinky/providers/base/websocket/metrics"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/types/factory"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// DefaultProviderFactory is a sample implementation of the provider factory. This provider
// factory function returns providers that are API & websocket based.
type DefaultProviderFactory struct {
	logger *zap.Logger

	// apiFactory is the factory function that creates API query handlers.
	apiFactory factory.APIQueryHandlerFactory[oracletypes.CurrencyPair, *big.Int]

	// wsFactory is the factory function that creates websocket query handlers.
	wsFactory factory.WebSocketQueryHandlerFactory[oracletypes.CurrencyPair, *big.Int]
}

// NewDefaultProviderFactory returns a new instance of the default provider factory.
func NewDefaultProviderFactory(
	logger *zap.Logger,
	apiFactory factory.APIQueryHandlerFactory[oracletypes.CurrencyPair, *big.Int],
	wsFactory factory.WebSocketQueryHandlerFactory[oracletypes.CurrencyPair, *big.Int],
) (*DefaultProviderFactory, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	if apiFactory == nil {
		return nil, fmt.Errorf("apiFactory cannot be nil")
	}

	if wsFactory == nil {
		return nil, fmt.Errorf("wsFactory cannot be nil")
	}

	return &DefaultProviderFactory{
		logger:     logger,
		apiFactory: apiFactory,
		wsFactory:  wsFactory,
	}, nil
}

// DefaultProviderFactory returns a sample implementation of the provider factory. This provider
// factory function returns providers that are API & websocket based.
func (f *DefaultProviderFactory) Factory() factory.ProviderFactory[oracletypes.CurrencyPair, *big.Int] {
	return func(cfg config.OracleConfig) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
		if err := cfg.ValidateBasic(); err != nil {
			return nil, err
		}

		// Create the metrics that are used by the providers.
		wsMetrics := wsmetrics.NewWebSocketMetricsFromConfig(cfg.Metrics)
		apiMetrics := apimetrics.NewAPIMetricsFromConfig(cfg.Metrics)
		providerMetrics := providermetrics.NewProviderMetricsFromConfig(cfg.Metrics)

		// Create the providers.
		providers := make([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], 0)
		for _, p := range cfg.Providers {
			switch {
			case p.API.Enabled:
				queryHandler, err := f.apiFactory(f.logger, p, apiMetrics)
				if err != nil {
					return nil, err
				}

				// Create the provider.
				provider, err := base.NewProvider[oracletypes.CurrencyPair, *big.Int](
					base.WithName[oracletypes.CurrencyPair, *big.Int](p.Name),
					base.WithLogger[oracletypes.CurrencyPair, *big.Int](f.logger),
					base.WithAPIQueryHandler(queryHandler),
					base.WithAPIConfig[oracletypes.CurrencyPair, *big.Int](p.API),
					base.WithIDs[oracletypes.CurrencyPair, *big.Int](cfg.Market.GetCurrencyPairs()),
					base.WithMetrics[oracletypes.CurrencyPair, *big.Int](providerMetrics),
				)
				if err != nil {
					return nil, err
				}

				providers = append(providers, provider)
			case p.WebSocket.Enabled:
				// Create the websocket query handler which encapsulates all fetching and parsing logic.
				queryHandler, err := f.wsFactory(f.logger, p, wsMetrics)
				if err != nil {
					return nil, err
				}

				// Create the provider.
				provider, err := base.NewProvider[oracletypes.CurrencyPair, *big.Int](
					base.WithName[oracletypes.CurrencyPair, *big.Int](p.Name),
					base.WithLogger[oracletypes.CurrencyPair, *big.Int](f.logger),
					base.WithWebSocketQueryHandler(queryHandler),
					base.WithWebSocketConfig[oracletypes.CurrencyPair, *big.Int](p.WebSocket),
					base.WithIDs[oracletypes.CurrencyPair, *big.Int](cfg.Market.GetCurrencyPairs()),
					base.WithMetrics[oracletypes.CurrencyPair, *big.Int](providerMetrics),
				)
				if err != nil {
					return nil, err
				}

				providers = append(providers, provider)
			default:
				f.logger.Info("unknown provider type", zap.String("provider", p.Name))
				return nil, fmt.Errorf("unknown provider type: %s", p.Name)
			}
		}

		return providers, nil
	}
}
