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
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
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
		Production: true,
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
			CurrencyPair: slinkytypes.NewCurrencyPair("ATOM", "USD"),
		},
		"ATOM/USDC": {
			CurrencyPair: slinkytypes.NewCurrencyPair("ATOM", "USDC"),
		},
		"ATOM/USDT": {
			CurrencyPair: slinkytypes.NewCurrencyPair("ATOM", "USDT"),
		},
		"AVAX/USD": {
			CurrencyPair: slinkytypes.NewCurrencyPair("AVAX", "USD"),
		},
		"AVAX/USDC": {
			CurrencyPair: slinkytypes.NewCurrencyPair("AVAX", "USDC"),
		},
		"AVAX/USDT": {
			CurrencyPair: slinkytypes.NewCurrencyPair("AVAX", "USDT"),
		},
		"BITCOIN/USD": {
			CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
		},
		"BITCOIN/USDC": {
			CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDC"),
		},
		"BITCOIN/USDT": {
			CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDT"),
		},
		"CELESTIA/USD": {
			CurrencyPair: slinkytypes.NewCurrencyPair("CELESTIA", "USD"),
		},
		"CELESTIA/USDC": {
			CurrencyPair: slinkytypes.NewCurrencyPair("CELESTIA", "USDC"),
		},
		"CELESTIA/USDT": {
			CurrencyPair: slinkytypes.NewCurrencyPair("CELESTIA", "USDT"),
		},
		"DYDX/USD": {
			CurrencyPair: slinkytypes.NewCurrencyPair("DYDX", "USD"),
		},
		"DYDX/USDC": {
			CurrencyPair: slinkytypes.NewCurrencyPair("DYDX", "USDC"),
		},
		"DYDX/USDT": {
			CurrencyPair: slinkytypes.NewCurrencyPair("DYDX", "USDT"),
		},
		"ETHEREUM/BITCOIN": {
			CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
		},
		"ETHEREUM/USD": {
			CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
		},
		"ETHEREUM/USDC": {
			CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USDC"),
		},
		"ETHEREUM/USDT": {
			CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USDT"),
		},
		"OSMOSIS/USD": {
			CurrencyPair: slinkytypes.NewCurrencyPair("OSMOSIS", "USD"),
		},
		"OSMOSIS/USDC": {
			CurrencyPair: slinkytypes.NewCurrencyPair("OSMOSIS", "USDC"),
		},
		"OSMOSIS/USDT": {
			CurrencyPair: slinkytypes.NewCurrencyPair("OSMOSIS", "USDT"),
		},
		"SOLANA/USD": {
			CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USD"),
		},
		"SOLANA/USDC": {
			CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USDC"),
		},
		"SOLANA/USDT": {
			CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USDT"),
		},
		"USDC/USD": {
			CurrencyPair: slinkytypes.NewCurrencyPair("USDC", "USD"),
		},
		"USDC/USDT": {
			CurrencyPair: slinkytypes.NewCurrencyPair("USDC", "USDT"),
		},
		"USDT/USD": {
			CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
		},
	}

	// AggregatedFeeds is a map of all of the conversion markets that will be used to convert
	// all of the price feeds into a common set of currency pairs.
	AggregatedFeeds = map[string]config.AggregateFeedConfig{
		"ATOM/USD": {
			CurrencyPair: slinkytypes.NewCurrencyPair("ATOM", "USD"),
			Conversions: []config.Conversions{
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("ATOM", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("ATOM", "USDC"),
						Invert:       false,
					},
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("USDC", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("ATOM", "USDT"),
						Invert:       false,
					},
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
						Invert:       false,
					},
				},
			},
		},
		"AVAX/USD": {
			CurrencyPair: slinkytypes.NewCurrencyPair("AVAX", "USD"),
			Conversions: []config.Conversions{
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("AVAX", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("AVAX", "USDC"),
						Invert:       false,
					},
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("USDC", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("AVAX", "USDT"),
						Invert:       false,
					},
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
						Invert:       false,
					},
				},
			},
		},
		"BITCOIN/USD": {
			CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
			Conversions: []config.Conversions{
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDC"),
						Invert:       false,
					},
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("USDC", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDT"),
						Invert:       false,
					},
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
						Invert:       false,
					},
				},
			},
		},
		"CELESTIA/USD": {
			CurrencyPair: slinkytypes.NewCurrencyPair("CELESTIA", "USD"),
			Conversions: []config.Conversions{
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("CELESTIA", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("CELESTIA", "USDC"),
						Invert:       false,
					},
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("USDC", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("CELESTIA", "USDT"),
						Invert:       false,
					},
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
						Invert:       false,
					},
				},
			},
		},
		"DYDX/USD": {
			CurrencyPair: slinkytypes.NewCurrencyPair("DYDX", "USD"),
			Conversions: []config.Conversions{
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("DYDX", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("DYDX", "USDC"),
						Invert:       false,
					},
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("USDC", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("DYDX", "USDT"),
						Invert:       false,
					},
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
						Invert:       false,
					},
				},
			},
		},
		"ETHEREUM/BITCOIN": {
			CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
			Conversions: []config.Conversions{
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
						Invert:       false,
					},
				},
			},
		},
		"ETHEREUM/USD": {
			CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
			Conversions: []config.Conversions{
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USDC"),
						Invert:       false,
					},
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("USDC", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USDT"),
						Invert:       false,
					},
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
						Invert:       false,
					},
				},
			},
		},
		"OSMOSIS/USD": {
			CurrencyPair: slinkytypes.NewCurrencyPair("OSMOSIS", "USD"),
			Conversions: []config.Conversions{
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("OSMOSIS", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("OSMOSIS", "USDC"),
						Invert:       false,
					},
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("USDC", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("OSMOSIS", "USDT"),
						Invert:       false,
					},
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
						Invert:       false,
					},
				},
			},
		},
		"SOLANA/USD": {
			CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USD"),
			Conversions: []config.Conversions{
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USDC"),
						Invert:       false,
					},
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("USDC", "USD"),
						Invert:       false,
					},
				},
				{
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USDT"),
						Invert:       false,
					},
					{
						CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
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
