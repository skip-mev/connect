//go:build ignore
// +build ignore

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/skip-mev/slinky/oracle/config"
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
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var (
	// oracleCfgPath is the path to write the oracle config file to. By default, this
	// will write the oracle config file to the local directory.
	oracleCfgPath = flag.String("oracle-config-path", "oracle.toml", "path to write the oracle config file to")

	// LocalConfig defines a readable config for local development. Any changes to this
	// file should be reflected in oracle.toml. To update the oracle.toml file, run
	// `make update-local-config`. This will update any changes to the oracle.toml file
	// as they are made to this file.
	LocalConfig = config.OracleConfig{
		// -----------------------------------------------------------	//
		// -----------------Aggregate Market Config-------------------	//
		// -----------------------------------------------------------	//
		Market: config.AggregateMarketConfig{
			Feeds:           Feeds,
			AggregatedFeeds: AggregatedFeeds,
		},
		Production: false,
		// -----------------------------------------------------------	//
		// ----------------------Metrics Config-----------------------	//
		// -----------------------------------------------------------	//
		Metrics: config.MetricsConfig{
			Enabled:                 true,
			PrometheusServerAddress: "0.0.0.0:8002",
		},
		UpdateInterval: 1 * time.Second,
		Providers: []config.ProviderConfig{
			// -----------------------------------------------------------	//
			// ---------------------Start API Providers--------------------	//
			// -----------------------------------------------------------	//
			//
			// NOTE: Some of the provider's are only capable of fetching data for a subset of
			// all currency pairs. Before adding a new market to the oracle, ensure that
			// the provider supports fetching data for the currency pair.
			{
				Name:   binance.Name,
				API:    binance.DefaultUSAPIConfig,
				Market: binance.DefaultUSMarketConfig,
			},
			{
				Name:   coinbaseapi.Name,
				API:    coinbaseapi.DefaultAPIConfig,
				Market: coinbaseapi.DefaultMarketConfig,
			},
			{
				Name:   coingecko.Name,
				API:    coingecko.DefaultAPIConfig,
				Market: coingecko.DefaultMarketConfig,
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
				Market:    bitfinex.DefaultMarketConfig,
			},
			{
				Name:      bitstamp.Name,
				WebSocket: bitstamp.DefaultWebSocketConfig,
				Market:    bitstamp.DefaultMarketConfig,
			},
			{
				Name:      bybit.Name,
				WebSocket: bybit.DefaultWebSocketConfig,
				Market:    bybit.DefaultMarketConfig,
			},
			{
				Name:      coinbasews.Name,
				WebSocket: coinbasews.DefaultWebSocketConfig,
				Market:    coinbasews.DefaultMarketConfig,
			},
			{
				Name:      cryptodotcom.Name,
				WebSocket: cryptodotcom.DefaultWebSocketConfig,
				Market:    cryptodotcom.DefaultMarketConfig,
			},
			{
				Name:      gate.Name,
				WebSocket: gate.DefaultWebSocketConfig,
				Market:    gate.DefaultMarketConfig,
			},
			{
				Name:      huobi.Name,
				WebSocket: huobi.DefaultWebSocketConfig,
				Market:    huobi.DefaultMarketConfig,
			},
			{
				Name:      kraken.Name,
				WebSocket: kraken.DefaultWebSocketConfig,
				Market:    kraken.DefaultMarketConfig,
			},
			{
				Name:      kucoin.Name,
				WebSocket: kucoin.DefaultWebSocketConfig,
				API:       kucoin.DefaultAPIConfig,
				Market:    kucoin.DefaultMarketConfig,
			},
			{
				Name:      mexc.Name,
				WebSocket: mexc.DefaultWebSocketConfig,
				Market:    mexc.DefaultMarketConfig,
			},
			{
				Name:      okx.Name,
				WebSocket: okx.DefaultWebSocketConfig,
				Market:    okx.DefaultMarketConfig,
			},
		},
	}

	// Feeds is a map of all of the price feeds that the oracle will fetch prices for.
	Feeds = map[string]config.FeedConfig{
		"ATOM/USD": {
			CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USD"),
		},
		"ATOM/USDC": {
			CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USDC"),
		},
		"ATOM/USDT": {
			CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USDT"),
		},
		"AVAX/USD": {
			CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USD"),
		},
		"AVAX/USDC": {
			CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USDC"),
		},
		"AVAX/USDT": {
			CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USDT"),
		},
		"BITCOIN/USD": {
			CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
		},
		"BITCOIN/USDC": {
			CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDC"),
		},
		"BITCOIN/USDT": {
			CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
		},
		"CELESTIA/USD": {
			CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USD"),
		},
		"CELESTIA/USDC": {
			CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USDC"),
		},
		"CELESTIA/USDT": {
			CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USDT"),
		},
		"DYDX/USD": {
			CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USD"),
		},
		"DYDX/USDC": {
			CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USDC"),
		},
		"DYDX/USDT": {
			CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USDT"),
		},
		"ETHEREUM/BITCOIN": {
			CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
		},
		"ETHEREUM/USD": {
			CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
		},
		"ETHEREUM/USDC": {
			CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDC"),
		},
		"ETHEREUM/USDT": {
			CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
		},
		"OSMOSIS/USD": {
			CurrencyPair: oracletypes.NewCurrencyPair("OSMOSIS", "USD"),
		},
		"OSMOSIS/USDC": {
			CurrencyPair: oracletypes.NewCurrencyPair("OSMOSIS", "USDC"),
		},
		"OSMOSIS/USDT": {
			CurrencyPair: oracletypes.NewCurrencyPair("OSMOSIS", "USDT"),
		},
		"SOLANA/USD": {
			CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD"),
		},
		"SOLANA/USDC": {
			CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USDC"),
		},
		"SOLANA/USDT": {
			CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USDT"),
		},
		"USDC/USD": {
			CurrencyPair: oracletypes.NewCurrencyPair("USDC", "USD"),
		},
		"USDC/USDT": {
			CurrencyPair: oracletypes.NewCurrencyPair("USDC", "USDT"),
		},
		"USDT/USD": {
			CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
		},
	}

	// AggregatedFeeds is a map of all of the conversion markets that will be used to convert
	// all of the price feeds into a common set of currency pairs.
	AggregatedFeeds = map[string]config.AggregateFeedConfig{
		"BITCOIN/USD": {
			CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			Conversions: []config.Conversions{
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
						Invert:       false,
					},
					{
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDC"),
						Invert:       false,
					},
					{
						CurrencyPair: oracletypes.NewCurrencyPair("USDC", "USD"),
						Invert:       false,
					},
				},
			},
		},
		"ETHEREUM/USD": {
			CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
			Conversions: []config.Conversions{
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
						Invert:       false,
					},
					{
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDC"),
						Invert:       false,
					},
					{
						CurrencyPair: oracletypes.NewCurrencyPair("USDC", "USD"),
						Invert:       false,
					},
				},
			},
		},
		"ATOM/USD": {
			CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USD"),
			Conversions: []config.Conversions{
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USDT"),
						Invert:       false,
					},
					{
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USDC"),
						Invert:       false,
					},
					{
						CurrencyPair: oracletypes.NewCurrencyPair("USDC", "USD"),
						Invert:       false,
					},
				},
			},
		},
		"SOLANA/USD": {
			CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD"),
			Conversions: []config.Conversions{
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USDT"),
						Invert:       false,
					},
					{
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USDC"),
						Invert:       false,
					},
					{
						CurrencyPair: oracletypes.NewCurrencyPair("USDC", "USD"),
						Invert:       false,
					},
				},
			},
		},
		"CELESTIA/USD": {
			CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USD"),
			Conversions: []config.Conversions{
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USDT"),
						Invert:       false,
					},
					{
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USDC"),
						Invert:       false,
					},
					{
						CurrencyPair: oracletypes.NewCurrencyPair("USDC", "USD"),
						Invert:       false,
					},
				},
			},
		},
		"AVAX/USD": {
			CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USD"),
			Conversions: []config.Conversions{
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USDT"),
						Invert:       false,
					},
					{
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USDC"),
						Invert:       false,
					},
					{
						CurrencyPair: oracletypes.NewCurrencyPair("USDC", "USD"),
						Invert:       false,
					},
				},
			},
		},
		"DYDX/USD": {
			CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USD"),
			Conversions: []config.Conversions{
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USDT"),
						Invert:       false,
					},
					{
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USDC"),
						Invert:       false,
					},
					{
						CurrencyPair: oracletypes.NewCurrencyPair("USDC", "USD"),
						Invert:       false,
					},
				},
			},
		},
		"ETHEREUM/BITCOIN": {
			CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
			Conversions: []config.Conversions{
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
						Invert:       false,
					},
				},
			},
		},
		"OSMOSIS/USD": {
			CurrencyPair: oracletypes.NewCurrencyPair("OSMOSIS", "USD"),
			Conversions: []config.Conversions{
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("OSMOSIS", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("OSMOSIS", "USDT"),
						Invert:       false,
					},
					{
						CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("OSMOSIS", "USDC"),
						Invert:       false,
					},
					{
						CurrencyPair: oracletypes.NewCurrencyPair("USDC", "USD"),
						Invert:       false,
					},
				},
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

	if err := LocalConfig.ValidateBasic(); err != nil {
		fmt.Fprintf(os.Stderr, "error validating local config: %v\n", err)
		return
	}

	// Encode the local config file.
	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(LocalConfig); err != nil {
		panic(err)
	}
}
