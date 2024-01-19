package main

import (
	"fmt"
	"os"

	"time"

	"github.com/BurntSushi/toml"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/apis/binanceus"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/websockets/cryptodotcom"
	"github.com/skip-mev/slinky/providers/websockets/okx"

	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// LocalConfig defines a readable config for local development. Any changes to this
// file should be reflected in oracle.toml. To update the oracle.toml file, run
// `make update-local-config`. This will update any changes to the oracle.toml file
// as they are made to this file.
var LocalConfig = config.OracleConfig{
	UpdateInterval: 1 * time.Second,
	Providers: []config.ProviderConfig{
		{
			// -----------------------------------------------------------	//
			// ---------------------Start API Providers--------------------	//
			// -----------------------------------------------------------	//
			Name: coinbase.Name,
			API: config.APIConfig{
				Atomic:     false,
				Enabled:    true,
				Timeout:    500 * time.Millisecond,
				Interval:   1 * time.Second,
				MaxQueries: 5,
			},
			MarketConfig: config.MarketConfig{
				Name: coinbase.Name,
				CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
					"BITCOIN/USD": {
						Ticker:       "BTC-USD",
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
					},
					"ETHEREUM/USD": {
						Ticker:       "ETH-USD",
						CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
					},
					"ATOM/USD": {
						Ticker:       "ATOM-USD",
						CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USD"),
					},
					"SOLANA/USD": {
						Ticker:       "SOL-USD",
						CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD"),
					},
					"CELESTIA/USD": {
						Ticker:       "TIA-USD",
						CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USD"),
					},
					"AVAX/USD": {
						Ticker:       "AVAX-USD",
						CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USD"),
					},
					"DYDX/USD": {
						Ticker:       "DYDX-USD",
						CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USD"),
					},
					"ETHEREUM/BITCOIN": {
						Ticker:       "ETH-BTC",
						CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
					},
					"OSMOSIS/USD": {
						Ticker:       "OSMO-USD",
						CurrencyPair: oracletypes.NewCurrencyPair("OSMOSIS", "USD"),
					},
				},
			},
		},
		{
			Name: binanceus.Name,
			API: config.APIConfig{
				Atomic:     true,
				Enabled:    true,
				Timeout:    500 * time.Millisecond,
				Interval:   1 * time.Second,
				MaxQueries: 1,
			},
			MarketConfig: config.MarketConfig{
				Name: "binanceus",
				CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
					"BITCOIN/USD": {
						Ticker:       "BTCUSDT",
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
					},
					"ETHEREUM/USD": {
						Ticker:       "ETHUSDT",
						CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
					},
					"ATOM/USD": {
						Ticker:       "ATOMUSDT",
						CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USD"),
					},
					"SOLANA/USD": {
						Ticker:       "SOLUSDT",
						CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD"),
					},
					"CELESTIA/USD": {
						Ticker:       "TIAUSDT",
						CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USD"),
					},
					"AVAX/USD": {
						Ticker:       "AVAXUSDT",
						CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USD"),
					},
					"DYDX/USD": {
						Ticker:       "DYDXUSDT",
						CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USD"),
					},
					"ETHEREUM/BITCOIN": {
						Ticker:       "ETHBTC",
						CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
					},
					"OSMOSIS/USD": {
						Ticker:       "OSMOUSDT",
						CurrencyPair: oracletypes.NewCurrencyPair("OSMOSIS", "USD"),
					},
				},
			},
		},
		// -----------------------------------------------------------	//
		// ---------------------Start WebSocket Providers--------------	//
		// -----------------------------------------------------------	//
		{
			Name: cryptodotcom.Name,
			WebSocket: config.WebSocketConfig{
				Enabled:             true,
				MaxBufferSize:       1000,
				ReconnectionTimeout: 5 * time.Second,
				WSS:                 cryptodotcom.ProductionURL,
			},
			MarketConfig: config.MarketConfig{
				Name: cryptodotcom.Name,
				CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
					"BITCOIN/USD": {
						Ticker:       "BTCUSD-PERP",
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
					},
					"ETHEREUM/USD": {
						Ticker:       "ETHUSD-PERP",
						CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
					},
					"ATOM/USD": {
						Ticker:       "ATOMUSD-PERP",
						CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USD"),
					},
					"SOLANA/USD": {
						Ticker:       "SOLUSD-PERP",
						CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD"),
					},
					"CELESTIA/USD": {
						Ticker:       "TIAUSD-PERP",
						CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USD"),
					},
					"AVAX/USD": {
						Ticker:       "AVAXUSD-PERP",
						CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USD"),
					},
					"DYDX/USD": {
						Ticker:       "DYDXUSD-PERP",
						CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USD"),
					},
					"ETHEREUM/BITCOIN": {
						Ticker:       "ETH_BTC",
						CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
					},
					"OSMOSIS/USD": {
						Ticker:       "OSMOUSD-PERP",
						CurrencyPair: oracletypes.NewCurrencyPair("OSMOSIS", "USD"),
					},
				},
			},
		},
		{
			Name: okx.Name,
			WebSocket: config.WebSocketConfig{
				Enabled:             true,
				MaxBufferSize:       1000,
				ReconnectionTimeout: 10 * time.Second,
				WSS:                 okx.ProductionURL,
			},
			MarketConfig: config.MarketConfig{
				Name: okx.Name,
				CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
					"BITCOIN/USD": {
						Ticker:       "BTC-USD",
						CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
					},
					"ETHEREUM/USD": {
						Ticker:       "ETH-USD",
						CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
					},
					"ATOM/USD": {
						Ticker:       "ATOM-USD",
						CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USD"),
					},
					"SOLANA/USD": {
						Ticker:       "SOL-USD",
						CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD"),
					},
					"CELESTIA/USD": {
						Ticker:       "TIA-USD",
						CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USD"),
					},
					"AVAX/USD": {
						Ticker:       "AVAX-USD",
						CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USD"),
					},
					"DYDX/USD": {
						Ticker:       "DYDX-USD",
						CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USD"),
					},
					"ETHEREUM/BITCOIN": {
						Ticker:       "ETH-BTC",
						CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
					},
					"OSMOSIS/USD": {
						Ticker:       "OSMO-USD",
						CurrencyPair: oracletypes.NewCurrencyPair("OSMOSIS", "USD"),
					},
				},
			},
		},
	},
	CurrencyPairs: []oracletypes.CurrencyPair{
		oracletypes.NewCurrencyPair("BITCOIN", "USD"),
		oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
		oracletypes.NewCurrencyPair("ATOM", "USD"),
		oracletypes.NewCurrencyPair("SOLANA", "USD"),
		oracletypes.NewCurrencyPair("CELESTIA", "USD"),
		oracletypes.NewCurrencyPair("AVAX", "USD"),
		oracletypes.NewCurrencyPair("DYDX", "USD"),
		oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
		oracletypes.NewCurrencyPair("OSMOSIS", "USD"),
	},
	Production: false,
	Metrics: config.MetricsConfig{
		Enabled:                 true,
		PrometheusServerAddress: "localhost:8000",
	},
}

const (
	// localPath is the path to the local config file.
	localPath = "./config/local/oracle.toml"
)

// main executes a simple script that encodes the local config file to the local
// directory.
func main() {
	// Open the local config file. This will overwrite any changes made to the
	// local config file.
	f, err := os.Create(localPath)
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
