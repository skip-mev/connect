//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
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
		binance.Name:     binance.DefaultUSMarketConfig,
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
	// `make update-local-config`. This will update any changes to the oracle.toml file
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
				API:  binance.DefaultUSAPIConfig,
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

	// Iterate through all of the ticker paths and update the optional ticker paths map.
	for ticker, paths := range TickerPaths {
		optionalTickerPaths[ticker.String()] = paths
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
