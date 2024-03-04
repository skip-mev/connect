//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/apis/binance"
	coinbaseapi "github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/apis/coingecko"
	"github.com/skip-mev/slinky/providers/websockets/bitfinex"
	"github.com/skip-mev/slinky/providers/websockets/bitstamp"
	"github.com/skip-mev/slinky/providers/websockets/bybit"
	coinbasews "github.com/skip-mev/slinky/providers/websockets/coinbase"
	"github.com/skip-mev/slinky/providers/websockets/cryptodotcom"
	"github.com/skip-mev/slinky/providers/websockets/gate"
	"github.com/skip-mev/slinky/providers/websockets/huobi"
	"github.com/skip-mev/slinky/providers/websockets/kraken"
	"github.com/skip-mev/slinky/providers/websockets/kucoin"
	"github.com/skip-mev/slinky/providers/websockets/mexc"
	"github.com/skip-mev/slinky/providers/websockets/okx"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	// oracleCfgPath is the path to write the oracle config file to. By default, this
	// will write the oracle config file to the local directory.
	oracleCfgPath = flag.String(
		"oracle-config-path",
		"oracle.json",
		"path to write the oracle config file to",
	)

	// marketCfgPath is the path to write the market config file to. By default, this
	// will write the market config file to the local directory.
	marketCfgPath = flag.String(
		"market-config-path",
		"market.json",
		"path to write the market config file to",
	)

	// providersOpt defines an optional list of providers to include in the local
	// market config file. If this is not provided, all providers will be included.
	providersOpt = flag.String(
		"providers",
		"",
		"optional list of providers to include in the local market config file",
	)

	// tickersOpt defines an optional list of tickers to include in the local market
	// config file. If this is not provided, all tickers will be included.
	tickersOpt = flag.String(
		"tickers",
		"",
		"optional list of tickers to include in the local market config file",
	)

	// usePathsOpt defines an optional flag to include the conversion paths in the
	// local market config file. If this is not provided, the conversion paths will
	// not be included.
	usePathsOpt = flag.Bool(
		"use-paths",
		false,
		"optional flag to include the conversion paths in the local market config file",
	)

	// TickerPaths defines a map of tickers to the corresponding conversion markets
	// that should be utilized to determine a final price.
	TickerPaths = map[mmtypes.Ticker]mmtypes.Paths{
		constants.ATOM_USD:     constants.ATOM_USD_PATHS,
		constants.AVAX_USD:     constants.AVAX_USD_PATHS,
		constants.BITCOIN_USD:  constants.BITCOIN_USD_PATHS,
		constants.CELESTIA_USD: constants.CELESTIA_USD_PATHS,
		constants.DYDX_USD:     constants.DYDX_USD_PATHS,
		constants.ETHEREUM_USD: constants.ETHEREUM_USD_PATHS,
		constants.OSMOSIS_USD:  constants.OSMOSIS_USD_PATHS,
		constants.SOLANA_USD:   constants.SOLANA_USD_PATHS,
	}

	// ProviderToMarkets defines a map of provider names to their respective market
	// configurations. This is used to generate the local market config file.
	ProviderToMarkets = map[string]types.TickerToProviderConfig{
		// -----------------------------------------------------------	//
		// ---------------------Start API Providers--------------------	//
		// -----------------------------------------------------------	//
		binance.Name:     binance.DefaultNonUSMarketConfig,
		coinbaseapi.Name: coinbaseapi.DefaultMarketConfig,
		coingecko.Name:   coingecko.DefaultMarketConfig,
		// // -----------------------------------------------------------	//
		// // ---------------------Start WebSocket Providers--------------	//
		// // -----------------------------------------------------------	//
		bitfinex.Name:     bitfinex.DefaultMarketConfig,
		bitstamp.Name:     bitstamp.DefaultMarketConfig,
		bybit.Name:        bybit.DefaultMarketConfig,
		coinbasews.Name:   coinbasews.DefaultMarketConfig,
		cryptodotcom.Name: cryptodotcom.DefaultMarketConfig,
		gate.Name:         gate.DefaultMarketConfig,
		huobi.Name:        huobi.DefaultMarketConfig,
		kraken.Name:       kraken.DefaultMarketConfig,
		kucoin.Name:       kucoin.DefaultMarketConfig,
		mexc.Name:         mexc.DefaultMarketConfig,
		okx.Name:          okx.DefaultMarketConfig,
	}

	// LocalConfig defines a readable config for local development. Any changes to this
	// file should be reflected in oracle.json. To update the oracle.json file, run
	// `make update-local-config`. This will update any changes to the oracle.json file
	// as they are made to this file.
	LocalOracleConfig = config.OracleConfig{
		Production: false,
		// -----------------------------------------------------------	//
		// ----------------------Metrics Config-----------------------	//
		// -----------------------------------------------------------	//
		Metrics: config.MetricsConfig{
			Enabled:                 true,
			PrometheusServerAddress: "0.0.0.0:8002",
		},
		UpdateInterval: 1500 * time.Millisecond,
		MaxPriceAge:    2 * time.Minute,
		Providers: []config.ProviderConfig{
			// -----------------------------------------------------------	//
			// ---------------------Start API Providers--------------------	//
			// -----------------------------------------------------------	//
			//
			// NOTE: Some of the provider's are only capable of fetching data for a subset of
			// all currency pairs. Before adding a new market to the oracle, ensure that
			// the provider supports fetching data for the currency pair.
			{
				Name: binance.Name,
				API:  binance.DefaultNonUSAPIConfig,
			},
			{
				Name: coinbaseapi.Name,
				API:  coinbaseapi.DefaultAPIConfig,
			},
			{
				Name: coingecko.Name,
				API:  coingecko.DefaultAPIConfig,
			},
			// -----------------------------------------------------------	//
			// ---------------------Start WebSocket Providers--------------	//
			// -----------------------------------------------------------	//
			//
			// NOTE: Some of the provider's are only capable of fetching data for a subset of
			// all currency pairs. Before adding a new market to the oracle, ensure that
			// the provider supports fetching data for the currency pair.
			{
				Name:      bitfinex.Name,
				WebSocket: bitfinex.DefaultWebSocketConfig,
			},
			{
				Name:      bitstamp.Name,
				WebSocket: bitstamp.DefaultWebSocketConfig,
			},
			{
				Name:      bybit.Name,
				WebSocket: bybit.DefaultWebSocketConfig,
			},
			{
				Name:      coinbasews.Name,
				WebSocket: coinbasews.DefaultWebSocketConfig,
			},
			{
				Name:      cryptodotcom.Name,
				WebSocket: cryptodotcom.DefaultWebSocketConfig,
			},
			{
				Name:      gate.Name,
				WebSocket: gate.DefaultWebSocketConfig,
			},
			{
				Name:      huobi.Name,
				WebSocket: huobi.DefaultWebSocketConfig,
			},
			{
				Name:      kraken.Name,
				WebSocket: kraken.DefaultWebSocketConfig,
			},
			{
				Name:      kucoin.Name,
				WebSocket: kucoin.DefaultWebSocketConfig,
				API:       kucoin.DefaultAPIConfig,
			},
			{
				Name:      mexc.Name,
				WebSocket: mexc.DefaultWebSocketConfig,
			},
			{
				Name:      okx.Name,
				WebSocket: okx.DefaultWebSocketConfig,
			},
		},
	}
)

// main executes a simple script that encodes the local config file to the local
// directory.
func main() {
	flag.Parse()

	// Open the local config file. This will overwrite any changes made to the
	// local config file.
	f, err := os.Create(*oracleCfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating local config file: %v\n", err)
	}
	defer f.Close()

	// If the providers is not empty, filter the providers to include only the
	// providers that are specified.
	providersFlag := *providersOpt
	if providersFlag != "" {
		ps := make([]config.ProviderConfig, 0)
		for _, provider := range LocalOracleConfig.Providers {
			if strings.Contains(providersFlag, provider.Name) {
				ps = append(ps, provider)
			}
		}

		LocalOracleConfig.Providers = ps
	}

	if err := LocalOracleConfig.ValidateBasic(); err != nil {
		fmt.Fprintf(os.Stderr, "error validating local config: %v\n", err)
		return
	}

	// Encode the local oracle config file.
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(LocalOracleConfig); err != nil {
		panic(err)
	}

	// Open the local market config file. This will overwrite any changes made to the
	// local market config file.
	f, err = os.Create(*marketCfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating local market config file: %v\n", err)
		return
	}
	defer f.Close()

	marketMap, err := createMarketMap()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating the local market map")
		return
	}

	// If the tickers is not empty, filter the tickers to include only the tickers
	// that are specified.
	tickersFlag := *tickersOpt
	if tickersFlag != "" {
		for _, ticker := range marketMap.Tickers {
			if !strings.Contains(tickersFlag, ticker.String()) {
				delete(marketMap.Providers, ticker.String())
				delete(marketMap.Paths, ticker.String())
				delete(marketMap.Tickers, ticker.String())
			}
		}

		if len(providersFlag) > 0 {
			for ticker, providers := range marketMap.Providers {
				ps := make([]mmtypes.ProviderConfig, 0)

				// Filter the providers to include only the providers that are specified.
				for _, provider := range providers.Providers {
					if strings.Contains(providersFlag, provider.Name) {
						ps = append(ps, provider)
					}
				}

				providers.Providers = ps
				marketMap.Providers[ticker] = providers

				// If there are no providers for a given ticker, remove the ticker from the
				// market map.
				if len(providers.Providers) == 0 {
					delete(marketMap.Providers, ticker)
					delete(marketMap.Paths, ticker)
					delete(marketMap.Tickers, ticker)
				}
			}
		}
	}

	// Validate the market map.
	if err := marketMap.ValidateBasic(); err != nil {
		fmt.Fprintf(os.Stderr, "error validating the market map: %v\n", err)
		return
	}

	// Encode the local market config file.
	encoder = json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(marketMap); err != nil {
		panic(err)
	}
}

// createMarketMap creates a market map given all of the local market configurations for
// each provider as well as the custom conversion markets. We do so to ensure that the
// oracle is always started using the market map that is expected to be stored by the
// market map module.
func createMarketMap() (mmtypes.MarketMap, error) {
	var (
		// Tickers defines a map of tickers to their respective ticker configurations. This
		// contains all of the tickers that are supported by the oracle.
		tickers = make(map[string]mmtypes.Ticker)
		// TickersToProviders defines a map of tickers to their respective providers. This
		// contains all of the providers that are supported per ticker.
		tickersToProviders = make(map[string]mmtypes.Providers)
		// OptionalTickerPaths defines a map of tickers to their respective conversion markets
		// that should be utilized to determine a final price. Not that this is optional as the
		// aggregation function utilized by the oracle may not require conversion markets to be
		// specified.
		optionalTickerPaths = make(map[string]mmtypes.Paths)
	)

	// Iterate through all of the provider ticker configurations and update the
	// tickers and tickers to providers maps.
	for name, providerConfig := range ProviderToMarkets {
		for ticker, config := range providerConfig {
			tickerStr := ticker.String()

			// Add the ticker to the tickers map iff the ticker does not already exist. If the
			// ticker already exists, ensure that the ticker configuration is the same.
			if t, ok := tickers[tickerStr]; !ok {
				tickers[tickerStr] = ticker
			} else {
				if t != ticker {
					return mmtypes.MarketMap{},
						fmt.Errorf("ticker %s already exists with different configuration for provider %s", tickerStr, name)
				}
			}

			// Instantiate the providers for a given ticker.
			if _, ok := tickersToProviders[tickerStr]; !ok {
				tickersToProviders[tickerStr] = mmtypes.Providers{}
			}

			// Add the provider to the tickers to providers map.
			providers := tickersToProviders[tickerStr].Providers
			providers = append(providers, config)
			tickersToProviders[tickerStr] = mmtypes.Providers{Providers: providers}
		}
	}
	if *usePathsOpt {
		for ticker, paths := range TickerPaths {
			optionalTickerPaths[ticker.String()] = paths
		}
	}

	// Create a new market map from the provider to market map.
	marketMap := mmtypes.MarketMap{
		Tickers:   tickers,
		Providers: tickersToProviders,
		Paths:     optionalTickerPaths,
	}

	// Validate the market map.
	if err := marketMap.ValidateBasic(); err != nil {
		return mmtypes.MarketMap{}, fmt.Errorf("error validating the market map: %w", err)
	}

	return marketMap, nil
}
