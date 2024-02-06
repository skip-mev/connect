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

var oracleCfgPath = flag.String("oracle-config-path", "oracle.toml", "path to write the oracle config file to")

// LocalConfig defines a readable config for local development. Any changes to this
// file should be reflected in oracle.toml. To update the oracle.toml file, run
// `make update-local-config`. This will update any changes to the oracle.toml file
// as they are made to this file.
var LocalConfig = config.OracleConfig{
	// -----------------------------------------------------------	//
	// -----------------Aggregate Market Config-------------------	//
	// -----------------------------------------------------------	//
	Market: config.AggregateMarketConfig{
		Feeds: map[string]config.FeedConfig{
			"BITCOIN/USD": {
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
			},
			"ETHEREUM/USD": {
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD", oracletypes.DefaultDecimals),
			},
			"ATOM/USD": {
				CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USD", oracletypes.DefaultDecimals),
			},
			"SOLANA/USD": {
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD", oracletypes.DefaultDecimals),
			},
			"CELESTIA/USD": {
				CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USD", oracletypes.DefaultDecimals),
			},
			"AVAX/USD": {
				CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USD", oracletypes.DefaultDecimals),
			},
			"DYDX/USD": {
				CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USD", oracletypes.DefaultDecimals),
			},
			"ETHEREUM/BITCOIN": {
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN", oracletypes.DefaultDecimals),
			},
			"OSMOSIS/USD": {
				CurrencyPair: oracletypes.NewCurrencyPair("OSMOSIS", "USD", oracletypes.DefaultDecimals),
			},
		},
		AggregatedFeeds: map[string][][]config.Conversion{
			"BITCOIN/USD": {
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD", oracletypes.DefaultDecimals),
						Invert:       false,
					},
				},
			},
			"ETHEREUM/USD": {
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD", oracletypes.DefaultDecimals),
						Invert:       false,
					},
				},
			},
			"ATOM/USD": {
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USD", oracletypes.DefaultDecimals),
						Invert:       false,
					},
				},
			},
			"SOLANA/USD": {
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD", oracletypes.DefaultDecimals),
						Invert:       false,
					},
				},
			},
			"CELESTIA/USD": {
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USD", oracletypes.DefaultDecimals),
						Invert:       false,
					},
				},
			},
			"AVAX/USD": {
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USD", oracletypes.DefaultDecimals),
						Invert:       false,
					},
				},
			},
			"DYDX/USD": {
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USD", oracletypes.DefaultDecimals),
						Invert:       false,
					},
				},
			},
			"ETHEREUM/BITCOIN": {
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN", oracletypes.DefaultDecimals),
						Invert:       false,
					},
				},
			},
			"OSMOSIS/USD": {
				{
					{
						CurrencyPair: oracletypes.NewCurrencyPair("OSMOSIS", "USD", oracletypes.DefaultDecimals),
						Invert:       false,
					},
				},
			},
		},
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
			Market: binance.DefaultMarketConfig,
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

	// Encode the local config file.
	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(LocalConfig); err != nil {
		panic(err)
	}
}
