package coinbase

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// NOTE: All documentation for this file can be located on the Coinbase
// API documentation: https://docs.cloud.coinbase.com/sign-in-with-coinbase/docs/api-prices#get-spot-price. This
// API does not require a subscription to use (i.e. No API key is required).

const (
	// Name is the name of the Coinbase provider.
	Name = "coinbase"

	// URL is the base URL of the Coinbase API. This includes the base and quote
	// currency pairs that need to be inserted into the URL.
	URL = "https://api.coinbase.com/v2/prices/%s/spot"
)

var (
	// DefaultAPIConfig is the default configuration for the Coinbase API.
	DefaultAPIConfig = config.APIConfig{
		Name:       Name,
		Atomic:     false,
		Enabled:    true,
		Timeout:    500 * time.Millisecond,
		Interval:   1 * time.Second,
		MaxQueries: 5,
		URL:        URL,
	}

	// DefaultMarketConfig is the default market configuration for Coinbase.
	DefaultMarketConfig = mmtypes.MarketConfig{
		Name: Name,
		TickerConfigs: map[string]mmtypes.TickerConfig{
			"ATOM/USD": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("ATOM", "USD"),
					Decimals:     8,
				},
				OffChainTicker: "ATOM-USD",
			},
			"ATOM/USDC": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("ATOM", "USDC"),
					Decimals:     8,
				},
				OffChainTicker: "ATOM-USDC",
			},
			"ATOM/USDT": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("ATOM", "USDT"),
					Decimals:     8,
				},
				OffChainTicker: "ATOM-USDT",
			},
			"AVAX/USD": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("AVAX", "USD"),
					Decimals:     8,
				},
				OffChainTicker: "AVAX-USD",
			},
			"AVAX/USDC": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("AVAX", "USDC"),
					Decimals:     8,
				},
				OffChainTicker: "AVAX-USDC",
			},
			"AVAX/USDT": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("AVAX", "USDT"),
					Decimals:     8,
				},
				OffChainTicker: "AVAX-USDT",
			},
			"BITCOIN/USD": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USD"),
					Decimals:     8,
				},
				OffChainTicker: "BTC-USD",
			},
			"BITCOIN/USDC": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDC"),
					Decimals:     8,
				},
				OffChainTicker: "BTC-USDC",
			},
			"BITCOIN/USDT": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDT"),
					Decimals:     8,
				},
				OffChainTicker: "BTC-USDT",
			},
			"CELESTIA/USD": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("CELESTIA", "USD"),
					Decimals:     8,
				},
				OffChainTicker: "TIA-USD",
			},
			"CELESTIA/USDC": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("CELESTIA", "USDC"),
					Decimals:     8,
				},
				OffChainTicker: "TIA-USDC",
			},
			"CELESTIA/USDT": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("CELESTIA", "USDT"),
					Decimals:     8,
				},
				OffChainTicker: "TIA-USDT",
			},
			"DYDX/USD": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("DYDX", "USD"),
					Decimals:     8,
				},
				OffChainTicker: "DYDX-USD",
			},
			"DYDX/USDC": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("DYDX", "USDC"),
					Decimals:     8,
				},
				OffChainTicker: "DYDX-USDC",
			},
			"DYDX/USDT": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("DYDX", "USDT"),
					Decimals:     8,
				},
				OffChainTicker: "DYDX-USDT",
			},
			"ETHEREUM/BITCOIN": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
					Decimals:     8,
				},
				OffChainTicker: "ETH-BTC",
			},
			"ETHEREUM/USD": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USD"),
					Decimals:     8,
				},
				OffChainTicker: "ETH-USD",
			},
			"ETHEREUM/USDC": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USDC"),
					Decimals:     8,
				},
				OffChainTicker: "ETH-USDC",
			},
			"ETHEREUM/USDT": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USDT"),
					Decimals:     8,
				},
				OffChainTicker: "ETH-USDT",
			},
			"OSMOSIS/USD": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("OSMOSIS", "USD"),
					Decimals:     8,
				},
				OffChainTicker: "OSMO-USD",
			},
			"OSMOSIS/USDC": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("OSMOSIS", "USDC"),
					Decimals:     8,
				},
				OffChainTicker: "OSMO-USDC",
			},
			"OSMOSIS/USDT": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("OSMOSIS", "USDT"),
					Decimals:     8,
				},
				OffChainTicker: "OSMO-USDT",
			},
			"SOLANA/USD": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USD"),
					Decimals:     8,
				},
				OffChainTicker: "SOL-USD",
			},
			"SOLANA/USDC": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USDC"),
					Decimals:     8,
				},
				OffChainTicker: "SOL-USDC",
			},
			"SOLANA/USDT": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USDT"),
					Decimals:     8,
				},
				OffChainTicker: "SOL-USDT",
			},
			"USDC/USD": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("USDC", "USD"),
					Decimals:     8,
				},
				OffChainTicker: "USDC-USD",
			},
			"USDC/USDT": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("USDC", "USDT"),
					Decimals:     8,
				},
				OffChainTicker: "USDC-USDT",
			},
			"USDT/USD": {
				Ticker: mmtypes.Ticker{
					CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
					Decimals:     8,
				},
				OffChainTicker: "USDT-USD",
			},
		},
	}
)

type (
	// CoinBaseResponse is the expected response returned by the Coinbase API.
	// The response is json formatted.
	// Response format:
	//
	//	{
	//	  "data": {
	//	    "amount": "1020.25",
	//	    "currency": "USD"
	//	  }
	//	}
	CoinBaseResponse struct { //nolint
		Data CoinBaseData `json:"data"`
	}

	// CoinBaseData is the data returned by the Coinbase API.
	CoinBaseData struct { //nolint
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	}
)
