package simapp

import (
	"fmt"

	"cosmossdk.io/log"

	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/evm/erc4626"
	"github.com/skip-mev/slinky/providers/evm/erc4626sharepriceoracle"
	"github.com/skip-mev/slinky/providers/mock"
	"github.com/skip-mev/slinky/x/oracle/types"
)

// DefaultProviderFactory returns a sample implementation of the provider factory.
func DefaultProviderFactory() oracle.ProviderFactory {
	return func(logger log.Logger, oracleCfg config.OracleConfig) ([]oracle.Provider, error) {
		providers := make([]oracle.Provider, len(oracleCfg.Providers))

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
func providerFromProviderConfig(logger log.Logger, cps []types.CurrencyPair, cfg config.ProviderConfig) (oracle.Provider, error) {
	switch cfg.Name {
	// TODO: Uncomment this when the coingecko API is fixed.
	// case "coingecko":
	// 	return coingecko.NewProvider(logger, cps, cfg)
	// TODO: Uncomment this when the coinbase API is fixed.
	// case "coinbase":
	// 	return coinbase.NewProvider(logger, cps, cfg)
	// TODO: Uncomment this when the coinmarketcap API is fixed.
	// case "coinmarketcap":
	// 	return coinmarketcap.NewProvider(logger, cps, cfg)
	case "erc4626":
		return erc4626.NewProvider(logger, cps, cfg)
	case "erc4626-share-price-oracle":
		return erc4626sharepriceoracle.NewProvider(logger, cps, cfg)
	case "failing-mock-provider":
		// This will always panic whenever GetPrices is called
		return mock.NewFailingMockProvider(), nil
	case "static-mock-provider":
		// This will always return the same price.
		return mock.NewStaticMockProviderFromConfig(cfg)
	default:
		return nil, fmt.Errorf("unknown provider: %s", cfg.Name)
	}
}
