package coinbase

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/constants"
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
				Ticker:         constants.ATOM_USD,
				OffChainTicker: "ATOM-USD",
			},
			"ATOM/USDC": {
				Ticker:         constants.ATOM_USDC,
				OffChainTicker: "ATOM-USDC",
			},
			"ATOM/USDT": {
				Ticker:         constants.ATOM_USDT,
				OffChainTicker: "ATOM-USDT",
			},
			"AVAX/USD": {
				Ticker:         constants.AVAX_USD,
				OffChainTicker: "AVAX-USD",
			},
			"AVAX/USDC": {
				Ticker:         constants.AVAX_USDC,
				OffChainTicker: "AVAX-USDC",
			},
			"AVAX/USDT": {
				Ticker:         constants.AVAX_USDT,
				OffChainTicker: "AVAX-USDT",
			},
			"BITCOIN/USD": {
				Ticker:         constants.BITCOIN_USD,
				OffChainTicker: "BTC-USD",
			},
			"BITCOIN/USDC": {
				Ticker:         constants.BITCOIN_USDC,
				OffChainTicker: "BTC-USDC",
			},
			"BITCOIN/USDT": {
				Ticker:         constants.BITCOIN_USDT,
				OffChainTicker: "BTC-USDT",
			},
			"CELESTIA/USD": {
				Ticker:         constants.CELESTIA_USD,
				OffChainTicker: "TIA-USD",
			},
			"CELESTIA/USDC": {
				Ticker:         constants.CELESTIA_USDC,
				OffChainTicker: "TIA-USDC",
			},
			"CELESTIA/USDT": {
				Ticker:         constants.CELESTIA_USDT,
				OffChainTicker: "TIA-USDT",
			},
			"DYDX/USD": {
				Ticker:         constants.DYDX_USD,
				OffChainTicker: "DYDX-USD",
			},
			"DYDX/USDC": {
				Ticker:         constants.DYDX_USDC,
				OffChainTicker: "DYDX-USDC",
			},
			"DYDX/USDT": {
				Ticker:         constants.ATOM_USDT,
				OffChainTicker: "DYDX-USDT",
			},
			"ETHEREUM/BITCOIN": {
				Ticker:         constants.ETHEREUM_BITCOIN,
				OffChainTicker: "ETH-BTC",
			},
			"ETHEREUM/USD": {
				Ticker:         constants.ETHEREUM_USD,
				OffChainTicker: "ETH-USD",
			},
			"ETHEREUM/USDC": {
				Ticker:         constants.ETHEREUM_USDC,
				OffChainTicker: "ETH-USDC",
			},
			"ETHEREUM/USDT": {
				Ticker:         constants.ETHEREUM_USDT,
				OffChainTicker: "ETH-USDT",
			},
			"OSMOSIS/USD": {
				Ticker:         constants.OSMOSIS_USD,
				OffChainTicker: "OSMO-USD",
			},
			"OSMOSIS/USDC": {
				Ticker:         constants.OSMOSIS_USDC,
				OffChainTicker: "OSMO-USDC",
			},
			"OSMOSIS/USDT": {
				Ticker:         constants.OSMOSIS_USDT,
				OffChainTicker: "OSMO-USDT",
			},
			"SOLANA/USD": {
				Ticker:         constants.SOLANA_USD,
				OffChainTicker: "SOL-USD",
			},
			"SOLANA/USDC": {
				Ticker:         constants.SOLANA_USDC,
				OffChainTicker: "SOL-USDC",
			},
			"SOLANA/USDT": {
				Ticker:         constants.SOLANA_USDT,
				OffChainTicker: "SOL-USDT",
			},
			"USDC/USD": {
				Ticker:         constants.USDC_USD,
				OffChainTicker: "USDC-USD",
			},
			"USDC/USDT": {
				Ticker:         constants.USDC_USDT,
				OffChainTicker: "USDC-USDT",
			},
			"USDT/USD": {
				Ticker:         constants.USDT_USD,
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
