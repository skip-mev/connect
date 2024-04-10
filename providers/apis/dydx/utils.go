package dydx

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
)

const (
	// Name is the name of the MarketMap provider.
	Name = "dydx_api"

	// ChainID is the chain ID for the dYdX market map provider.
	ChainID = "dydx-node"

	// Endpoint is the endpoint for the dYdX market map API.
	Endpoint = "%s/dydxprotocol/prices/params/market?limit=10000"

	// Delimeter is the delimeter used to separate the base and quote assets in a pair.
	Delimeter = "-"
)

// DefaultAPIConfig returns the default configuration for the dYdX market map API.
var DefaultAPIConfig = config.APIConfig{
	Name:             Name,
	Atomic:           true,
	Enabled:          true,
	Timeout:          20 * time.Second, // Set a high timeout to account for slow API responses in the case where many markets are queried.
	Interval:         10 * time.Second,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       1,
	URL:              "localhost:1317",
}
