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

	// LocalMarketConfig defines a readable config for local development. Any changes to this
	// file should be reflected in market.json. To update the market.json file, run
	// `make update-local-config`. This will update any changes to the market.json file
	// as they are made to this file.
	LocalMarketConfig = mmtypes.AggregateMarketConfig{
		MarketConfigs: map[string]mmtypes.MarketConfig{
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
		},
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

	if err := LocalMarketConfig.ValidateBasic(); err != nil {
		fmt.Fprintf(os.Stderr, "error validating local market config: %v\n", err)
		return
	}

	// Encode the local market config file.
	encoder = json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(LocalMarketConfig); err != nil {
		panic(err)
	}
}
