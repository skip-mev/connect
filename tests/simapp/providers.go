package simapp

import (
	"fmt"
	"math/big"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/coinbase"
	"github.com/skip-mev/slinky/providers/coingecko"
	"github.com/skip-mev/slinky/providers/coinmarketcap"
	"github.com/skip-mev/slinky/providers/evm/erc4626"
	"github.com/skip-mev/slinky/providers/evm/erc4626sharepriceoracle"
	"github.com/skip-mev/slinky/providers/static"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// DefaultProviderFactory returns a sample implementation of the provider factory.
func DefaultProviderFactory() providertypes.ProviderFactory[oracletypes.CurrencyPair, *big.Int] {
	return func(logger *zap.Logger, oracleCfg config.OracleConfig) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
		providers := make([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], len(oracleCfg.Providers))

		var err error
		for i, p := range oracleCfg.Providers {
			if providers[i], err = providerFromProviderConfig(logger, oracleCfg.CurrencyPairs, p); err != nil {
				return nil, err
			}
		}

		return providers, nil
	}
}

// providerFromProviderConfig returns a provider from a provider config. These providers are
// NOT production ready and are only meant for testing purposes.
func providerFromProviderConfig(
	logger *zap.Logger,
	cps []oracletypes.CurrencyPair,
	cfg config.ProviderConfig,
) (providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
	switch cfg.Name {
	case "coingecko":
		return coingecko.NewProvider(logger, cps, cfg)
	case "coinbase":
		return coinbase.NewProvider(logger, cps, cfg)
	case "coinmarketcap":
		return coinmarketcap.NewProvider(logger, cps, cfg)
	case "erc4626":
		return erc4626.NewProvider(logger, cps, cfg)
	case "erc4626-share-price-oracle":
		return erc4626sharepriceoracle.NewProvider(logger, cps, cfg)
	case "static-mock-provider":
		// This will always return the same price.
		return static.NewProvider(logger, cps, cfg)
	default:
		return nil, fmt.Errorf("unknown provider: %s", cfg.Name)
	}
}
