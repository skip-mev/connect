package coinbase

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
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
	DefaultMarketConfig = config.MarketConfig{
		Name: Name,
		CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
			"BITCOIN/USD": {
				Ticker:       "BTC-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			"BITCOIN/USDT": {
				Ticker:       "BTC-USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
			},
			"BITCOIN/USDC": {
				Ticker:       "BTC-USDC",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDC"),
			},
			"ETHEREUM/USD": {
				Ticker:       "ETH-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			"ETHEREUM/USDT": {
				Ticker:       "ETH-USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
			},
			"ETHEREUM/USDC": {
				Ticker:       "ETH-USDC",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDC"),
			},
			"ATOM/USD": {
				Ticker:       "ATOM-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USD"),
			},
			"ATOM/USDT": {
				Ticker:       "ATOM-USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USDT"),
			},
			"ATOM/USDC": {
				Ticker:       "ATOM-USDC",
				CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USDC"),
			},
			"SOLANA/USD": {
				Ticker:       "SOL-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD"),
			},
			"SOLANA/USDT": {
				Ticker:       "SOL-USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USDT"),
			},
			"SOLANA/USDC": {
				Ticker:       "SOL-USDC",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USDC"),
			},
			"CELESTIA/USD": {
				Ticker:       "TIA-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USD"),
			},
			"CELESTIA/USDT": {
				Ticker:       "TIA-USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USDT"),
			},
			"CELESTIA/USDC": {
				Ticker:       "TIA-USDC",
				CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USDC"),
			},
			"AVAX/USD": {
				Ticker:       "AVAX-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USD"),
			},
			"AVAX/USDT": {
				Ticker:       "AVAX-USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USDT"),
			},
			"AVAX/USDC": {
				Ticker:       "AVAX-USDC",
				CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USDC"),
			},
			"DYDX/USD": {
				Ticker:       "DYDX-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USD"),
			},
			"DYDX/USDT": {
				Ticker:       "DYDX-USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USDT"),
			},
			"DYDX/USDC": {
				Ticker:       "DYDX-USDC",
				CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USDC"),
			},
			"ETHEREUM/BITCOIN": {
				Ticker:       "ETH-BTC",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
			},
			"OSMOSIS/USD": {
				Ticker:       "OSMO-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("OSMOSIS", "USD"),
			},
			"OSMOSIS/USDT": {
				Ticker:       "OSMO-USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("OSMOSIS", "USDT"),
			},
			"OSMOSIS/USDC": {
				Ticker:       "OSMO-USDC",
				CurrencyPair: oracletypes.NewCurrencyPair("OSMOSIS", "USDC"),
			},
			"USDT/USD": {
				Ticker:       "USDT-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
			},
			"USDC/USD": {
				Ticker:       "USDC-USD",
				CurrencyPair: oracletypes.NewCurrencyPair("USDC", "USD"),
			},
			"USDC/USDT": {
				Ticker:       "USDC-USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("USDC", "USDT"),
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
