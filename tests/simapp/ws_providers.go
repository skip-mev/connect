package simapp

import (
	"fmt"
	"math/big"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/websockets/cryptodotcom"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"go.uber.org/zap"
)

// DefaultWebSocketProviderFactory returns a sample implementation of the provider factory. This provider
// factory function only returns providers the are web socket based.
func DefaultWebSocketProviderFactory() providertypes.ProviderFactory[oracletypes.CurrencyPair, *big.Int] {
	return func(logger *zap.Logger, oracleCfg config.OracleConfig, metricsCfg config.OracleMetricsConfig) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
		if err := oracleCfg.ValidateBasic(); err != nil {
			return nil, err
		}

		cps := oracleCfg.CurrencyPairs

		providers := make([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], 0)
		for _, p := range oracleCfg.Providers {
			// Skip providers that are not web socket based.
			if !p.WebSocket.Enabled {
				continue
			}

			provider, err := webSocketProviderFromProviderConfig(logger, p, cps)
			if err != nil {
				return nil, err
			}

			providers = append(providers, provider)
		}

		return providers, nil
	}
}

// webSocketProviderFromProviderConfig returns a provider from a provider config. These providers are
// NOT production ready and are only meant for testing purposes.
func webSocketProviderFromProviderConfig(
	logger *zap.Logger,
	cfg config.ProviderConfig,
	cps []oracletypes.CurrencyPair,
) (providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
	// Validate the provider config.
	err := cfg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	var (
		wsDataHandler handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int]
		connHandler   handlers.WebSocketConnHandler
	)

	switch cfg.Name {
	case cryptodotcom.Name:
		wsDataHandler, err = cryptodotcom.NewWebSocketDataHandlerFromConfig(logger, cfg)
	default:
		return nil, fmt.Errorf("unknown provider: %s", cfg.Name)
	}
	if err != nil {
		return nil, err
	}

	// If a custom request handler is not provided, create a new default one.
	if connHandler == nil {
		connHandler = handlers.NewWebSocketHandlerImpl()
	}

	// Create the API query handler which encapsulates all of the fetching and parsing logic.
	apiQueryHandler, err := handlers.NewWebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int](
		logger,
		wsDataHandler,
		connHandler,
	)
	if err != nil {
		return nil, err
	}

	// Create the provider.
	return base.NewProvider[oracletypes.CurrencyPair, *big.Int](
		cfg,
		base.WithLogger[oracletypes.CurrencyPair, *big.Int](logger),
		base.WithWebSocketQueryHandler(apiQueryHandler),
		base.WithIDs[oracletypes.CurrencyPair, *big.Int](cps),
	)
}
