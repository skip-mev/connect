package geckoterminal

import (
	"time"

	"github.com/skip-mev/connect/v2/oracle/config"
)

// NOTE: All documentation for this file can be located on the GeckoTerminal
// API specification:
//
// - https://api.geckoterminal.com/api/v2
// - https://www.geckoterminal.com/dex-api.

const (
	// Name is the name of the GeckoTerminal provider.
	Name = "gecko_terminal_api"

	// URL is the root URL for the GeckoTerminal API.
	ETH_URL = "https://api.geckoterminal.com/api/v2/simple/networks/eth/token_price/%s"

	// ExpectedResponseType is the expected attribute name for the response type in the
	// GeckoTerminal API response.
	ExpectedResponseType = "simple_token_price"
)

// DefaultETHAPIConfig is the default configuration for querying Ethereum mainnet tokens
// on the GeckoTerminal API.
var DefaultETHAPIConfig = config.APIConfig{
	Name:             Name,
	Atomic:           false,
	Enabled:          true,
	Timeout:          500 * time.Millisecond,
	Interval:         20 * time.Second,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       1,
	Endpoints:        []config.Endpoint{{URL: ETH_URL}},
}

type (
	// GeckoTerminalResponse is the expected response returned by the GeckoTerminal API.
	// The response is json formatted.
	// Response format:
	//
	// {
	// 	"data": {
	// 	  "id": "61ba8f36-7962-4a75-acc1-bdb07bb7eda5",
	// 	  "type": "simple_token_price",
	// 	  "attributes": {
	// 		"token_prices": {
	// 		  "0xaaee1a9723aadb7afa2810263653a34ba2c21c7a": "0.000000970708264000586"
	// 		}
	// 	  }
	// 	}
	// }.
	GeckoTerminalResponse struct { //nolint
		Data GeckoTerminalData `json:"data"`
	}

	// GeckoTerminalData is the data field in the GeckoTerminalResponse.
	GeckoTerminalData struct { //nolint
		ID         string                  `json:"id"`
		Type       string                  `json:"type"`
		Attributes GeckoTerminalAttributes `json:"attributes"`
	}

	// GeckoTerminalAttributes is the attributes field in the GeckoTerminalData.
	GeckoTerminalAttributes struct { //nolint
		TokenPrices map[string]string `json:"token_prices"`
	}
)
