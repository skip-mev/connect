package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/apis/binance"
	coinbaseapi "github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/apis/coingecko"
	"github.com/skip-mev/slinky/providers/apis/dydx"
	"github.com/skip-mev/slinky/providers/apis/geckoterminal"
	krakenapi "github.com/skip-mev/slinky/providers/apis/kraken"
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
	mmclienttypes "github.com/skip-mev/slinky/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

const (
	dYdXChain = "dydx"
)

var (
	// oracleCfgPath is the path to write the oracle config file to. By default, this
	// will write the oracle config file to the local directory.
	oracleCfgPath = flag.String(
		"oracle-config-path",
		"oracle.json",
		"path to write the oracle config file to. this file is required to run the oracle.",
	)

	// marketCfgPath is the path to write the market config file to. By default, this
	// will write the market config file to the local directory.
	marketCfgPath = flag.String(
		"market-config-path",
		"market.json",
		"path to write the market config file to. this file is required to run the oracle.",
	)

	// chain defines the chain that we expect the oracle to be running on. Currently only
	// supprting dYdX.
	chain = flag.String(
		"chain",
		"",
		"chain that we expect the oracle to be running on. ex dydx.",
	)

	// nodeURL is the URL of the validator. This is required if running the oracle with a market map provider.
	nodeURL = flag.String(
		"node-http-url",
		"",
		"URL of the dYdX node. this is required if running the oracle on the dYdX chain. this is the http address of the dYdX node.",
	)

	// host is the oracle / prometheus server host.
	host = flag.String(
		"host",
		"0.0.0.0",
		"host is the oracle / prometheus server host.",
	)

	// pricesPort is the port that the oracle will make prices available on.
	pricesPort = flag.String(
		"port",
		"8080",
		"port that the oracle will make prices available on. to query prices after starting the oracle, use the following command: curl http://<host>:<port>/slinky/oracle/v1/prices",
	)

	// prometheusPort is the port that prometheus will make metrics available on.
	prometheusPort = flag.String(
		"prometheus-port",
		"8002",
		"port that the prometheus server will listen on. to query prometheus metrics after starting the oracle, use the following command: curl http://<host>:<port>/metrics",
	)

	// disabledMetrics is a flag that disables the prometheus server.
	disabledMetrics = flag.Bool(
		"disable-metrics",
		false,
		"flag that disables the prometheus server. if this is enabled the prometheus port must be specified. to query prometheus metrics after starting the oracle, use the following command: curl http://<host>:<port>/metrics",
	)

	// debug is a flag that enables debug mode. Specifically, all logging will be
	// in debug mode.
	debug = flag.Bool(
		"debug-mode",
		false,
		"flag that enables debug mode. specifically the side-car will run in debug mode. this is useful for local development / debugging.",
	)

	// updateInterval is the interval at which the oracle will update the prices.
	updateInterval = flag.Duration(
		"update-interval",
		1500*time.Millisecond,
		"interval at which the oracle will update the prices. this should be set to the interval desired by the chain.",
	)

	// maxPriceAge is the maximum age of a price that the oracle will accept.
	maxPriceAge = flag.Duration(
		"max-price-age",
		2*time.Minute,
		"maximum age of a price that the oracle will accept. this should be set to the maximum age desired by the chain.",
	)

	// ProviderToMarkets defines a map of provider names to their respective market
	// configurations. This is used to generate the local market config file.
	ProviderToMarkets = map[string]types.TickerToProviderConfig{
		// -----------------------------------------------------------	//
		// ---------------------Start API Providers--------------------	//
		// -----------------------------------------------------------	//
		binance.Name:       binance.DefaultNonUSMarketConfig,
		coinbaseapi.Name:   coinbaseapi.DefaultMarketConfig,
		coingecko.Name:     coingecko.DefaultMarketConfig,
		geckoterminal.Name: geckoterminal.DefaultETHMarketConfig,
		krakenapi.Name:     krakenapi.DefaultMarketConfig,
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
		Production: true,
		// -----------------------------------------------------------	//
		// ----------------------Metrics Config-----------------------	//
		// -----------------------------------------------------------	//
		Metrics:        config.MetricsConfig{},
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
				Type: types.ConfigType,
			},
			{
				Name: coinbaseapi.Name,
				API:  coinbaseapi.DefaultAPIConfig,
				Type: types.ConfigType,
			},
			{
				Name: coingecko.Name,
				API:  coingecko.DefaultAPIConfig,
				Type: types.ConfigType,
			},
			{
				Name: geckoterminal.Name,
				API:  geckoterminal.DefaultETHAPIConfig,
				Type: types.ConfigType,
			},
			{
				Name: krakenapi.Name,
				API:  krakenapi.DefaultAPIConfig,
				Type: types.ConfigType,
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
				Type:      types.ConfigType,
			},
			{
				Name:      bitstamp.Name,
				WebSocket: bitstamp.DefaultWebSocketConfig,
				Type:      types.ConfigType,
			},
			{
				Name:      bybit.Name,
				WebSocket: bybit.DefaultWebSocketConfig,
				Type:      types.ConfigType,
			},
			{
				Name:      coinbasews.Name,
				WebSocket: coinbasews.DefaultWebSocketConfig,
				Type:      types.ConfigType,
			},
			{
				Name:      cryptodotcom.Name,
				WebSocket: cryptodotcom.DefaultWebSocketConfig,
				Type:      types.ConfigType,
			},
			{
				Name:      gate.Name,
				WebSocket: gate.DefaultWebSocketConfig,
				Type:      types.ConfigType,
			},
			{
				Name:      huobi.Name,
				WebSocket: huobi.DefaultWebSocketConfig,
				Type:      types.ConfigType,
			},
			{
				Name:      kraken.Name,
				WebSocket: kraken.DefaultWebSocketConfig,
				Type:      types.ConfigType,
			},
			{
				Name:      kucoin.Name,
				WebSocket: kucoin.DefaultWebSocketConfig,
				API:       kucoin.DefaultAPIConfig,
				Type:      types.ConfigType,
			},
			{
				Name:      mexc.Name,
				WebSocket: mexc.DefaultWebSocketConfig,
				Type:      types.ConfigType,
			},
			{
				Name:      okx.Name,
				WebSocket: okx.DefaultWebSocketConfig,
				Type:      types.ConfigType,
			},
		},
	}
)

// main executes a simple script that encodes the local config file to the local
// directory.
func main() {
	flag.Parse()

	// Create the oracle config that contains all of the providers that are supported.
	if err := createOracleConfig(); err != nil {
		panic(err)
	}

	// Create the market map that contains all of the tickers and providers that are
	// supported.
	if err := createMarketMap(); err != nil {
		panic(err)
	}
}

// createOracleConfig creates an oracle config given all of the local provider configurations. If the
// the chain is set to dYdX, the dYdX node URL is required. We do so to ensure that the oracle is
// always started using the oracle config that is expected to be stored by the chain.
func createOracleConfig() error {
	// If the providers is not empty, filter the providers to include only the
	// providers that are specified.
	if *chain == dYdXChain {
		// Filter out the providers that are not supported by the dYdX chain.
		validProviders := make(map[string]struct{})
		for _, providers := range dydx.ProviderMapping {
			for _, provider := range providers {
				validProviders[provider] = struct{}{}
			}
		}

		ps := make([]config.ProviderConfig, 0)
		for _, provider := range LocalOracleConfig.Providers {
			if _, ok := validProviders[provider.Name]; ok {
				ps = append(ps, provider)
			}
		}

		if len(*nodeURL) == 0 {
			return fmt.Errorf("dYdX node URL is required; please specify your dYdX node URL using the --node-http-url flag (ex. --node-http-url http://localhost:1317)")
		}
		apiCfg := dydx.DefaultAPIConfig
		apiCfg.URL = *nodeURL

		// Add the dYdX market map provider to the list of providers.
		ps = append(ps, config.ProviderConfig{
			Name: dydx.Name,
			API:  apiCfg,
			Type: mmclienttypes.ConfigType,
		})
		LocalOracleConfig.Providers = ps
	}

	// Set the host and port for the oracle.
	LocalOracleConfig.Host = *host
	LocalOracleConfig.Port = *pricesPort

	// Set the prometheus server address for the oracle.
	if !*disabledMetrics {
		LocalOracleConfig.Metrics.Enabled = true
		LocalOracleConfig.Metrics.PrometheusServerAddress = fmt.Sprintf("%s:%s", *host, *prometheusPort)
	}

	// Set the update interval for the oracle.
	LocalOracleConfig.UpdateInterval = *updateInterval
	LocalOracleConfig.MaxPriceAge = *maxPriceAge

	if *debug {
		LocalOracleConfig.Production = false
	}

	if err := LocalOracleConfig.ValidateBasic(); err != nil {
		fmt.Fprintf(os.Stderr, "error validating local config: %v\n", err)
		return err
	}

	f, err := os.Create(*oracleCfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating local config file: %v\n", err)
	}
	defer f.Close()

	// Encode the local oracle config file.
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(LocalOracleConfig); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "successfully created oracle config file at %s\n", *oracleCfgPath)
	return nil
}

// createMarketMap creates a market map given all of the local market configurations for
// each provider as well as the custom conversion markets. We do so to ensure that the
// oracle is always started using the market map that is expected to be stored by the
// market map module.
func createMarketMap() error {
	if *chain == dYdXChain {
		fmt.Fprintf(
			os.Stderr,
			"dYdX chain requires the use of a predetermined market map. please use the market map provided by the Skip/dYdX team or the default market map provided in /config/dydx/market.json",
		)
		return nil
	}

	var (
		// Tickers defines a map of tickers to their respective ticker configurations. This
		// contains all of the tickers that are supported by the oracle.
		tickers = make(map[string]mmtypes.Ticker)
		// TickersToProviders defines a map of tickers to their respective providers. This
		// contains all of the providers that are supported per ticker.
		tickersToProviders = make(map[string]mmtypes.Providers)
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
			} else if t != ticker {
				return fmt.Errorf("ticker %s already exists with different configuration for provider %s", tickerStr, name)
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

	// Create a new market map from the provider to market map.
	marketMap := mmtypes.MarketMap{
		Tickers:   tickers,
		Providers: tickersToProviders,
	}

	// Validate the market map.
	if err := marketMap.ValidateBasic(); err != nil {
		return fmt.Errorf("error validating the market map: %w", err)
	}

	// Open the local market config file. This will overwrite any changes made to the
	// local market config file.
	f, err := os.Create(*marketCfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating local market config file: %v\n", err)
		return err
	}
	defer f.Close()

	// Validate the market map.
	if err := marketMap.ValidateBasic(); err != nil {
		fmt.Fprintf(os.Stderr, "error validating the market map: %v\n", err)
		return err
	}

	// Encode the local market config file.
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(marketMap); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "successfully created market config file at %s\n", *marketCfgPath)
	return nil
}
