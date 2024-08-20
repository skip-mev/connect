package coinbase

import (
	"time"

	"github.com/skip-mev/connect/v2/oracle/config"
)

// NOTE: All documentation for this file can be located on the Coinbase
// API documentation: https://docs.cloud.coinbase.com/sign-in-with-coinbase/docs/api-prices#get-spot-price. This
// API does not require a subscription to use (i.e. No API key is required).

const (
	// Name is the name of the Coinbase provider.
	Name = "coinbase_api"

	// URL is the base URL of the Coinbase API. This includes the base and quote
	// currency pairs that need to be inserted into the URL.
	URL = "https://api.coinbase.com/v2/prices/%s/spot"
)

// DefaultAPIConfig is the default configuration for the Coinbase API.
var DefaultAPIConfig = config.APIConfig{
	Name:             Name,
	Atomic:           false,
	Enabled:          true,
	Timeout:          3000 * time.Millisecond,
	Interval:         100 * time.Millisecond,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       1,
	Endpoints:        []config.Endpoint{{URL: URL}},
}

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
