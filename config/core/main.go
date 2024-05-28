package main

import (
	"encoding/json"
	"os"

	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

func main() {
	cfg, err := mmtypes.ReadMarketMapFromFile("./config/core/market.json")
	if err != nil {
		panic(err)
	}

	for k, market := range cfg.Markets {
		// Do something with the market
		cfgs := market.ProviderConfigs
		for _, provider := range cfgs {
			if provider.Name == "binance_ws" {
				apiCfg := mmtypes.ProviderConfig{
					Name:            "binance_api",
					OffChainTicker:  provider.OffChainTicker,
					Metadata_JSON:   provider.Metadata_JSON,
					NormalizeByPair: provider.NormalizeByPair,
				}
				cfgs = append(cfgs, apiCfg)
			}
		}

		market.ProviderConfigs = cfgs
		cfg.Markets[k] = market
	}

	//Write json file
	f, err := os.Create("output.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cfg); err != nil {
		panic(err)
	}
}
