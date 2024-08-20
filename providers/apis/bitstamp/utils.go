package bitstamp

import (
	"time"

	"github.com/skip-mev/connect/v2/oracle/config"
)

// NOTE: All documentation for this file can be located on the Bitstamp GitHub
// API documentation: https://www.bitstamp.net/api/v2/ticker/. This
// API does not require a subscription to use (i.e. No API key is required).

const (
	// Name is the name of the Bitstamp provider.
	Name = "bitstamp_api"

	// URL is the base URL of the Bitstamp API. The URL returns all prices, some
	// of which are not needed.
	URL = "https://www.bitstamp.net/api/v2/ticker/"
)

// DefaultAPIConfig is the default configuration for the Bitstamp API.
var DefaultAPIConfig = config.APIConfig{
	Name:             Name,
	Atomic:           true,
	Enabled:          true,
	Timeout:          3000 * time.Millisecond,
	Interval:         3000 * time.Millisecond,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       1,
	Endpoints:        []config.Endpoint{{URL: URL}},
}

// MarketTickerResponse is the expected response returned by the Bitstamp API.
//
// ex.
//
// [
//
//	{
//		"ask": "2211.00",
//		"bid": "2188.97",
//		"high": "2811.00",
//		"last": "2211.00",
//		"low": "2188.97",
//		"open": "2211.00",
//		"open_24": "2211.00",
//		"pair": "BTC/USD",
//		"percent_change_24": "13.57",
//		"side": "0",
//		"timestamp": "1643640186",
//		"volume": "213.26801100",
//		"vwap": "2189.80"
//	}
//
// ]
//
// ref: https://www.bitstamp.net/api/v2/ticker/
type MarketTickerResponse []MarketTickerData

// MarketTickerData is the data returned by the Bitstamp API.
type MarketTickerData struct {
	Last string `json:"last"`
	Pair string `json:"pair"`
}
