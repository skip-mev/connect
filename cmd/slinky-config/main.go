package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/providers"
	"github.com/skip-mev/slinky/providers/apis/binance"
	coinbaseapi "github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/apis/coingecko"
	raydium "github.com/skip-mev/slinky/providers/apis/defi/raydium"
	"github.com/skip-mev/slinky/providers/apis/defi/uniswapv3"
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
	// raydiumPairsFixture is the path to the fixture file containing all raydium markets.
	raydiumPairsFixture = "./cmd/slinky-config/fixtures/raydium_pairs.json"
)

var (
	rootCmd = &cobra.Command{
		Use:   "slinky-config",
		Short: "Create configuration required for running slinky.",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			// Create the oracle config that contains all providers that are supported.
			if err := createOracleConfig(); err != nil {
				panic(err)
			}

			// Create the market map that contains all tickers and providers that are
			// supported.
			if err := createMarketMap(); err != nil {
				panic(err)
			}
		},
	}

	// oracleCfgPath is the path to write the oracle config file to. By default, this
	// will write the oracle config file to the local directory.
	oracleCfgPath string
	// marketCfgPath is the path to write the market config file to. By default, this
	// will write the market config file to the local directory.
	marketCfgPath string
	// chain defines the chain that we expect the oracle to be running on.
	chain string
	// dydxResearchJSONMarketMap determines whether we want to fetch the dydx market-map
	// from the chain, or the dydx research JSON file.
	dydxResearchJSONMarketMap bool
	// nodeURL is the URL of the validator. This is required if running the oracle with the dydx market map provider.
	nodeURL string
	// host is the oracle / prometheus server host.
	host string
	// pricesPort is the port that the oracle will make prices available on.
	pricesPort string
	// prometheusPort is the port that prometheus will make metrics available on.
	prometheusPort string
	// disabledMetrics is a flag that disables the prometheus server.
	disabledMetrics bool
	// debug is a flag that enables debug mode. Specifically, all logging will be
	// in debug mode.
	debug bool
	// updateInterval is the interval at which the oracle will update the prices.
	updateInterval time.Duration
	// maxPriceAge is the maximum age of a price that the oracle will accept.
	maxPriceAge time.Duration
	// raydium-enabled determine whether or not the raydium defi provider will be configured.
	raydiumEnabled bool
	// solana node url is the solana node that the raydium provider will connect to.
	solanaNodeURLs []string
	// uniswapv3-enabled determines whether or not the uniswapv3 defi provider will be configured.
	uniswapv3Enabled bool
	// ethNodeURLs is the set of ethereum nodes evm providers will connect to.
	ethNodeURLs []string
	// providerMarkets is the set of providers to output config for. It is used for testing.
	providerMarkets []string
	// ProviderToMarkets defines a map of provider names to their respective market
	// configurations. This is used to generate the local market config file.
	ProviderToMarkets = map[string]types.CurrencyPairsToProviderTickers{
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
		// // -----------------------------------------------------------	//
		// // ---------------------Start Defi Providers-------------------	//
		// // -----------------------------------------------------------	//
		uniswapv3.ProviderNames[constants.ETHEREUM]: uniswapv3.DefaultETHMarketConfig,
	}

	// LocalOracleConfig defines a readable config for local development. Any changes to this
	// file should be reflected in oracle.json. To update the oracle.json file, run
	// `make update-local-config`. This will update any changes to the oracle.json file
	// as they are made to this file.
	LocalOracleConfig = config.OracleConfig{
		// -----------------------------------------------------------	//
		// ----------------------Metrics Config-----------------------	//
		// -----------------------------------------------------------	//
		Metrics:        config.MetricsConfig{},
		UpdateInterval: 250 * time.Millisecond,
		MaxPriceAge:    2 * time.Minute,
		Providers:      providers.ProviderDefaults,
	}
)

func init() {
	rootCmd.Flags().StringVarP(
		&oracleCfgPath,
		"oracle-config-path",
		"",
		"oracle.json",
		"Path to write the oracle config file to. This file is required to run the oracle.",
	)
	rootCmd.Flags().StringVarP(
		&marketCfgPath,
		"market-config-path",
		"",
		"market.json",
		"Path to write the market config file to. This file is required to run the oracle.",
	)
	rootCmd.Flags().StringVarP(
		&chain,
		"chain",
		"",
		"",
		"Chain that we expect the oracle to be running on {dydx, \"\"}. This should only be specified if required by the chain.",
	)

	rootCmd.Flags().BoolVarP(
		&dydxResearchJSONMarketMap,
		"dydx-research-json-market-map",
		"",
		false,
		"Use the dydx-research json to configure markets alongside the chain.",
	)

	rootCmd.Flags().StringVarP(
		&nodeURL,
		"node-http-url",
		"",
		"",
		"Http endpoint of the cosmos sdk node corresponding to the chain (typically localhost:1317 or a remote API). This should only be specified if required by the chain.",
	)
	rootCmd.Flags().StringVarP(
		&host,
		"host",
		"",
		"0.0.0.0",
		"Host is the oracle / prometheus server host.",
	)
	rootCmd.Flags().StringVarP(
		&pricesPort,
		"port",
		"",
		"8080",
		"Port that the oracle will make prices available on. To query prices after starting the oracle, use the following command: curl http://<host>:<port>/slinky/oracle/v1/prices",
	)
	rootCmd.Flags().StringVarP(
		&prometheusPort,
		"prometheus-port",
		"",
		"8002",
		"Port that the prometheus server will listen on. To query prometheus metrics after starting the oracle, use the following command: curl http://<host>:<port>/metrics",
	)
	rootCmd.Flags().BoolVarP(
		&disabledMetrics,
		"disable-metrics",
		"",
		false,
		"Flag that disables the prometheus server. If this is enabled the prometheus port must be specified. To query prometheus metrics after starting the oracle, use the following command: curl http://<host>:<port>/metrics",
	)
	rootCmd.Flags().BoolVarP(
		&debug,
		"debug-mode",
		"",
		false,
		"Flag that enables debug mode for the side-car. This is useful for local development / debugging.",
	)
	rootCmd.Flags().DurationVarP(
		&updateInterval,
		"update-interval",
		"",
		250*time.Millisecond,
		"Interval at which the oracle will update the prices. This should be set to the interval desired by the chain.",
	)
	rootCmd.Flags().DurationVarP(
		&maxPriceAge,
		"max-price-age",
		"",
		2*time.Minute,
		"Maximum age of a price that the oracle will accept. This should be set to the maximum age desired by the chain.",
	)
	rootCmd.Flags().BoolVarP(
		&raydiumEnabled,
		"raydium-enabled",
		"",
		false,
		"whether or not to enable raydium support",
	)
	rootCmd.Flags().StringSliceVarP(
		&solanaNodeURLs,
		"solana-node-endpoint",
		"",
		nil,
		"The HTTP endpoints of the solana node endpoint the raydium provider will be configured to use. If multiple are given they must be comma delimited",
	)
	rootCmd.Flags().BoolVarP(
		&uniswapv3Enabled,
		"uniswapv3-enabled",
		"",
		false,
		"whether or not to enable uniswapv3 support",
	)
	rootCmd.Flags().StringSliceVarP(
		&ethNodeURLs,
		"eth-node-endpoint",
		"",
		nil,
		"The HTTP endpoints of the eth node endpoint the eth providers will be configured to use. If multiple are given they must be comma delimited",
	)
	rootCmd.Flags().StringSliceVarP(
		&providerMarkets,
		"provider-markets",
		"",
		nil,
		"The set of providers to add markets for.",
	)
}

// main executes a simple script that encodes the local config file to the local
// directory.
func main() {
	rootCmd.Execute()
}

func configureDYDXProviders() error {
	// Filter out the providers that are not supported by the dYdX chain.
	validProviders := make(map[string]struct{})
	for _, slinkyProvider := range dydx.ProviderMapping {
		validProviders[slinkyProvider] = struct{}{}
	}

	ps := make([]config.ProviderConfig, 0)
	for _, provider := range LocalOracleConfig.Providers {
		if _, ok := validProviders[provider.Name]; ok {
			ps = append(ps, provider)
		}
	}

	// if we want to use the research json
	var marketMapProviderConfig config.APIConfig
	name := dydx.Name
	if dydxResearchJSONMarketMap {
		marketMapProviderConfig = dydx.DefaultResearchAPIConfig
		name = dydx.ResearchAPIHandlerName
	} else {
		if len(nodeURL) == 0 {
			return fmt.Errorf("dYdX node URL is required; please specify your dYdX node URL using the --node-http-url flag (ex. --node-http-url http://localhost:1317)")
		}
		marketMapProviderConfig = dydx.DefaultAPIConfig
		marketMapProviderConfig.URL = nodeURL
	}

	// Add the dYdX market map provider to the list of providers.
	ps = append(ps, config.ProviderConfig{
		Name: name,
		API:  marketMapProviderConfig,
		Type: mmclienttypes.ConfigType,
	})
	LocalOracleConfig.Providers = ps
	return nil
}

// createOracleConfig creates an oracle config given all local provider configurations.
func createOracleConfig() error {
	// If the providers is not empty, filter the providers to include only the
	// providers that are specified.
	if strings.ToLower(chain) == constants.DYDX {
		if err := configureDYDXProviders(); err != nil {
			return err
		}
	}

	// add raydium provider to the list of providers if enabled
	if raydiumEnabled {
		cfg := raydium.DefaultAPIConfig
		for _, node := range solanaNodeURLs {
			cfg.Endpoints = append(cfg.Endpoints, config.Endpoint{
				URL: node,
			})
		}

		LocalOracleConfig.Providers = append(LocalOracleConfig.Providers, config.ProviderConfig{
			Name: raydium.Name,
			API:  cfg,
			Type: types.ConfigType,
		})
	}
	if uniswapv3Enabled {
		cfg := uniswapv3.DefaultETHAPIConfig
		for _, node := range ethNodeURLs {
			cfg.Endpoints = append(cfg.Endpoints, config.Endpoint{
				URL: node,
			})
		}
	}

	// Set the host and port for the oracle.
	LocalOracleConfig.Host = host
	LocalOracleConfig.Port = pricesPort

	// Set the prometheus server address for the oracle.
	if !disabledMetrics {
		LocalOracleConfig.Metrics.Enabled = true
		LocalOracleConfig.Metrics.PrometheusServerAddress = fmt.Sprintf("%s:%s", host, prometheusPort)
	}

	// Set the update interval for the oracle.
	LocalOracleConfig.UpdateInterval = updateInterval
	LocalOracleConfig.MaxPriceAge = maxPriceAge

	if err := LocalOracleConfig.ValidateBasic(); err != nil {
		fmt.Fprintf(os.Stderr, "error validating local config: %v\n", err)
		return err
	}

	f, err := os.Create(oracleCfgPath)
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

	fmt.Fprintf(os.Stdout, "successfully created oracle config file at %s\n", oracleCfgPath)
	return nil
}

// createMarketMap creates a market map given all local market configurations for
// each provider as well as the custom conversion markets. We do so to ensure that the
// oracle is always started using the market map that is expected to be stored by the
// market map module.
func createMarketMap() error {
	if strings.ToLower(chain) == constants.DYDX {
		fmt.Fprintf(
			os.Stderr,
			"dYdX chain requires the use of a predetermined market map. please use the market map provided by the Skip/dYdX team or the default market map provided in /config/dydx/market.json",
		)
		return nil
	}

	// Tickers defines a map of tickers to their respective ticker configurations. This
	// contains all tickers that are supported by the oracle.
	marketMap := mmtypes.MarketMap{
		Markets: make(map[string]mmtypes.Market),
	}

	if providerMarkets != nil {
		pruned := make(map[string]types.CurrencyPairsToProviderTickers)
		for _, provider := range providerMarkets {
			val, ok := ProviderToMarkets[provider]
			if ok {
				pruned[provider] = val
			}
		}
		ProviderToMarkets = pruned
	}

	// if raydium is enabled, configure the raydium markets based on the local raydium_pairs fixture
	if raydiumEnabled {
		ProviderToMarkets = addRaydiumMarkets(ProviderToMarkets)
	}

	// Iterate through all provider ticker configurations and update the
	// tickers and tickers to providers maps.
	for provider, providerConfig := range ProviderToMarkets {
		for cp, config := range providerConfig {
			ticker := mmtypes.Ticker{
				CurrencyPair:     cp,
				Decimals:         18,
				MinProviderCount: 1,
				Enabled:          true,
			}

			// Add the ticker to the tickers map iff the ticker does not already exist.
			if _, ok := marketMap.Markets[ticker.String()]; !ok {
				marketMap.Markets[ticker.String()] = mmtypes.Market{
					Ticker:          ticker,
					ProviderConfigs: make([]mmtypes.ProviderConfig, 0),
				}
			}

			market := marketMap.Markets[ticker.String()]
			market.ProviderConfigs = append(market.ProviderConfigs, mmtypes.ProviderConfig{
				Name:           provider,
				OffChainTicker: config.OffChainTicker,
				Metadata_JSON:  config.JSON,
			})
			marketMap.Markets[ticker.String()] = market
		}
	}

	// Validate the market map.
	if err := marketMap.ValidateBasic(); err != nil {
		return fmt.Errorf("error validating the market map: %w", err)
	}

	// Open the local market config file. This will overwrite any changes made to the
	// local market config file.
	f, err := os.Create(marketCfgPath)
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

	fmt.Fprintf(os.Stdout, "successfully created market config file at %s\n", marketCfgPath)
	return nil
}

type TickerMetaData struct {
	Cp             slinkytypes.CurrencyPair `json:"currency_pair"`
	TickerMetaData raydium.TickerMetadata   `json:"ticker_metadata"`
}

func addRaydiumMarkets(providerToMarkets map[string]types.CurrencyPairsToProviderTickers) map[string]types.CurrencyPairsToProviderTickers {
	// read the raydium_pairs fixture
	if !raydiumEnabled {
		return providerToMarkets
	}

	// read the raydium_pairs fixture
	file, err := os.Open(raydiumPairsFixture)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading raydium_pairs fixture: %v\n", err)
		return providerToMarkets
	}
	defer file.Close()

	bz, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading raydium_pairs fixture: %v\n", err)
		return providerToMarkets
	}

	var raydiumPairs []TickerMetaData
	if err := json.Unmarshal(bz, &raydiumPairs); err != nil {
		fmt.Fprintf(os.Stderr, "error unmarshalling raydium_pairs fixture: %v\n", err)
		return providerToMarkets
	}

	// add the raydium markets to the provider to markets map
	providerToMarkets[raydium.Name] = make(types.CurrencyPairsToProviderTickers)
	for _, pair := range raydiumPairs {
		providerToMarkets[raydium.Name][pair.Cp] = types.DefaultProviderTicker{
			OffChainTicker: pair.Cp.String(),
			JSON:           marshalToJSONString(pair.TickerMetaData),
		}
	}

	return providerToMarkets
}

func marshalToJSONString(obj interface{}) string {
	bz, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(bz)
}
