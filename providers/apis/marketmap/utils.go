package marketmap

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/service/clients/marketmap/types"
)

const (
	// Name is the name of the MarketMap provider.
	Name = "marketmap_api"
	Type = types.ConfigType
)

// DefaultAPIConfig returns the default configuration for the MarketMap API.
var DefaultAPIConfig = config.APIConfig{
	Name:             Name,
	Atomic:           true,
	Enabled:          true,
	Timeout:          20 * time.Second,
	Interval:         10 * time.Second,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       1,
	URL:              "localhost:9090",
}

var DefaultProviderConfig = config.ProviderConfig{
	Name: Name,
	API:  DefaultAPIConfig,
	Type: Type,
}
