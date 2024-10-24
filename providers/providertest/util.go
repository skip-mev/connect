package providertest

import (
	"fmt"

	"github.com/skip-mev/slinky/cmd/constants"
	cmdconfig "github.com/skip-mev/slinky/cmd/slinky/config"
	"github.com/skip-mev/slinky/oracle/config"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

func FilterMarketMapToProviders(mm mmtypes.MarketMap) map[string]mmtypes.MarketMap {
	m := make(map[string]mmtypes.MarketMap)

	for _, market := range mm.Markets {
		// check each provider config
		for _, pc := range market.ProviderConfigs {
			// remove normalizations to isolate markets
			pc.NormalizeByPair = nil

			// create a market from the given provider config
			isolatedMarket := mmtypes.Market{
				Ticker: market.Ticker,
				ProviderConfigs: []mmtypes.ProviderConfig{
					pc,
				},
			}

			// always enable and set minprovider count to 1 so that it can be run isolated
			isolatedMarket.Ticker.Enabled = true
			isolatedMarket.Ticker.MinProviderCount = 1

			// init mm if necessary
			if _, found := m[pc.Name]; !found {
				m[pc.Name] = mmtypes.MarketMap{
					Markets: map[string]mmtypes.Market{
						isolatedMarket.Ticker.String(): isolatedMarket,
					},
				}
				// otherwise insert
			} else {
				m[pc.Name].Markets[isolatedMarket.Ticker.String()] = isolatedMarket
			}
		}
	}

	return m
}

func OracleConfigForProvider(providerNames ...string) (config.OracleConfig, error) {
	cfg := config.OracleConfig{
		UpdateInterval: cmdconfig.DefaultUpdateInterval,
		MaxPriceAge:    cmdconfig.DefaultMaxPriceAge,
		Metrics: config.MetricsConfig{
			Enabled: false,
		},
		Providers: make(map[string]config.ProviderConfig),
		Host:      cmdconfig.DefaultHost,
		Port:      cmdconfig.DefaultPort,
	}

	for _, provider := range append(constants.Providers, constants.AlternativeMarketMapProviders...) {
		for _, providerName := range providerNames {
			if provider.Name == providerName {
				cfg.Providers[provider.Name] = provider
			}
		}
	}

	if err := cfg.ValidateBasic(); err != nil {
		return cfg, fmt.Errorf("default oracle config is invalid: %w", err)
	}

	return cfg, nil
}
