package config

import (
	"fmt"
	"time"
)

// APIConfig defines a config for an API based data provider.
type APIConfig struct {
	// Enabled is a flag that indicates whether the provider is API based.
	Enabled bool `json:"enabled"`

	// Timeout is the amount of time the provider should wait for a response from
	// its API before timing out.
	Timeout time.Duration `json:"timeout"`

	// Interval is the interval at which the provider should update the prices.
	Interval time.Duration `json:"interval"`

	// ReconnectTimeout is the amount of time the provider should wait before
	// reconnecting to the API.
	ReconnectTimeout time.Duration `json:"reconnectTimeout"`

	// MaxQueries is the maximum number of queries that the provider will make
	// within the interval. If the provider makes more queries than this, it will
	// stop making queries until the next interval.
	MaxQueries int `json:"maxQueries"`

	// Atomic is a flag that indicates whether the provider can fulfill its queries
	// in a single request.
	Atomic bool `json:"atomic"`

	// URL is the URL that is used to fetch data from the API.
	URL string `json:"url"`

	// Endpoints is a list of endpoints that the provider can query.
	Endpoints []Endpoint `json:"endpoints"`

	// Name is the name of the provider that corresponds to this config.
	Name string `json:"name"`
}

// Endpoint holds all data necessary for an API provider to connect to a given endpoint
// i.e URL, headers, authentication, etc
type Endpoint struct {
	// URL is the URL that is used to fetch data from the API.
	URL string `json:"url"`
}

// ValidateBasic performs basic validation of the API config.
func (c *APIConfig) ValidateBasic() error {
	if !c.Enabled {
		return nil
	}

	if c.MaxQueries < 1 {
		return fmt.Errorf("api max queries must be greater than 0")
	}

	if c.Interval <= 0 || c.Timeout <= 0 || c.ReconnectTimeout <= 0 {
		return fmt.Errorf("provider interval, timeout and reconnect timeout must be strictly positive")
	}

	if len(c.URL) == 0 && len(c.Endpoints) == 0 {
		return fmt.Errorf("provider url cannot be empty")
	}

	if len(c.Name) == 0 {
		return fmt.Errorf("provider name cannot be empty")
	}

	return nil
}
