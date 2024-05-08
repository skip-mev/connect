package marketmap

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
)

const (
	// Name is the name of the MarketMap provider.
	Name = "marketmap_api"
)

// DefaultAPIConfig returns the default configuration for the MarketMap API.
var DefaultAPIConfig = config.APIConfig{
	Name:             Name,
	Atomic:           true,
	Enabled:          true,
	Timeout:          500 * time.Millisecond,
	Interval:         1 * time.Second,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       1,
	URL:              "localhost:9090",
}
