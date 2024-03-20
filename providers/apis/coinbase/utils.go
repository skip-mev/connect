package coinbase

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
)

// NOTE: All documentation for this file can be located on the Coinbase
// API documentation: https://docs.cloud.coinbase.com/sign-in-with-coinbase/docs/api-prices#get-spot-price. This
// API does not require a subscription to use (i.e. No API key is required).

const (
	// Name is the name of the Coinbase provider.
	Name = "CoinbaseProAPI"

	// URL is the base URL of the Coinbase API. This includes the base and quote
	// currency pairs that need to be inserted into the URL.
	URL = "https://api.coinbase.com/v2/prices/%s/spot"
)

var (
	// DefaultAPIConfig is the default configuration for the Coinbase API.
	DefaultAPIConfig = config.APIConfig{
		Name:             Name,
		Atomic:           false,
		Enabled:          true,
		Timeout:          500 * time.Millisecond,
		Interval:         100 * time.Millisecond,
		ReconnectTimeout: 2000 * time.Millisecond,
		MaxQueries:       5,
		URL:              URL,
	}

	// DefaultMarketConfig is the default market configuration for Coinbase.
	DefaultMarketConfig = types.TickerToProviderConfig{
		constants.BITCOIN_USD: {
			Name:           Name,
			OffChainTicker: "BTC-USD",
		},
		constants.BITCOIN_USDT: {
			Name:           Name,
			OffChainTicker: "BTC-USDT",
		},
		constants.BITCOIN_USDC: {
			Name:           Name,
			OffChainTicker: "BTC-USDC",
		},
		constants.USDC_USD: {
			Name:           Name,
			OffChainTicker: "USDC-USD",
		},
		constants.USDC_USDT: {
			Name:           Name,
			OffChainTicker: "USDC-USDT",
		},
		constants.USDT_USD: {
			Name:           Name,
			OffChainTicker: "USDT-USD",
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
