package polymarket

import (
	"time"

	"github.com/skip-mev/connect/v2/oracle/config"
)

var DefaultAPIConfig = config.APIConfig{
	Name:             Name,
	Atomic:           false,
	Enabled:          true,
	Timeout:          3 * time.Second,
	Interval:         500 * time.Millisecond,
	ReconnectTimeout: 2 * time.Second,
	MaxQueries:       1,
	Endpoints:        []config.Endpoint{{URL: URL}},
}
