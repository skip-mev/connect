package marketmap

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
)

const (
	// Name is the name of the MarketMap provider.
	Name = "marketmap"
)

// DefaultAPIConfig returns the default configuration for the MarketMap API.
var DefaultAPIConfig = config.APIConfig{
	Name:       Name,
	Atomic:     true,
	Enabled:    true,
	Timeout:    500 * time.Millisecond,
	Interval:   1 * time.Second,
	MaxQueries: 1,
	URL:        "http://localhost:1317/slinky/marketmap/v1/marketmap",
}
