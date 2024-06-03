//go:build ignore
// +build ignore

package main

import (
	"flag"
	"fmt"

	"github.com/skip-mev/slinky/cmd/constants"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	useCore    = flag.Bool("use-core", false, "use core markets")
	useRaydium = flag.Bool("use-raydium", false, "use raydium markets")
	tempFile   = flag.String("temp-file", "markets.json", "temporary file")
)

func main() {
	// Based on the flags, we determine what market.json to configure. By default, we use Core markets.
	// If the user specifies a different market.json, we use that instead.
	flag.Parse()

	marketMap := mmtypes.MarketMap{
		Markets: make(map[string]mmtypes.Market),
	}

	if *useCore {
		fmt.Fprintf(flag.CommandLine.Output(), "Using core markets\n")
		marketMap = mergeMarketMaps(marketMap, constants.CoreMarketMap)
	}

	if *useRaydium {
		fmt.Fprintf(flag.CommandLine.Output(), "Using raydium markets\n")
		marketMap = mergeMarketMaps(marketMap, constants.RaydiumMarketMap)
	}

	// Write the market map to the temporary file.
	if *tempFile == "" {
		fmt.Fprintf(flag.CommandLine.Output(), "temp file cannot be empty\n")
		panic("temp file cannot be empty")
	}

	if err := mmtypes.WriteMarketMapToFile(marketMap, *tempFile); err != nil {
		fmt.Fprintf(flag.CommandLine.Output(), "failed to write market map to file: %s\n", err)
		panic(err)
	}
}

// mergeMarketMaps merges the two market maps together. If a market already exists in one of the maps, we
// merge based on the provider set.
func mergeMarketMaps(this, other mmtypes.MarketMap) mmtypes.MarketMap {
	for name, market := range other.Markets {
		// If the market does not exist in this map, we add it.
		if _, ok := this.Markets[name]; !ok {
			this.Markets[name] = market
			continue
		}

		seen := make(map[string]struct{})
		for _, provider := range market.ProviderConfigs {
			key := providerConfigToKey(provider)
			seen[key] = struct{}{}
		}

		for _, provider := range this.Markets[name].ProviderConfigs {
			key := providerConfigToKey(provider)
			if _, ok := seen[key]; !ok {
				market.ProviderConfigs = append(market.ProviderConfigs, provider)
			}
		}

		this.Markets[name] = market
	}

	return this
}

func providerConfigToKey(cfg mmtypes.ProviderConfig) string {
	return cfg.Name + cfg.OffChainTicker
}
