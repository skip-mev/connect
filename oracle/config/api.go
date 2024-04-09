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

	// MaxQueries is the maximum number of concurrent queries that the provider will make
	// within the interval.
	MaxQueries int `json:"maxQueries"`

	// Atomic is a flag that indicates whether the provider can fulfill its queries
	// in a single request.
	Atomic bool `json:"atomic"`

	// URL is the URL that is used to fetch data from the API.
	URL string `json:"url"`

	// Endpoints is a list of endpoints that the provider can query.
	Endpoints []Endpoint `json:"endpoints"`

	// BatchSize is the maximum number of IDs that the provider can query in a single
	// request. This parameter must be 0 for atomic providers. Otherwise, the effective
	// value will be max(1, BatchSize). Notice, if numCPs > batchSize * maxQueries then
	// some currency-pairs may not be fetched each interval.
	BatchSize int `json:"batchSize"`

	// Name is the name of the provider that corresponds to this config.
	Name string `json:"name"`
}

// Endpoint holds all data necessary for an API provider to connect to a given endpoint
// i.e URL, headers, authentication, etc.
type Endpoint struct {
	// URL is the URL that is used to fetch data from the API.
	URL string `json:"url"`

	// Authentication holds all data necessary for an API provider to authenticate with
	// an endpoint.
	Authentication Authentication `json:"authentication"`
}

// ValidateBasic performs basic validation of the API endpoint.
func (e Endpoint) ValidateBasic() error {
	if len(e.URL) == 0 {
		return fmt.Errorf("endpoint url cannot be empty")
	}

	return e.Authentication.ValidateBasic()
}

// Authentication holds all data necessary for an API provider to authenticate with an
// endpoint.
type Authentication struct {
	// Enabled is a flag that indicates whether the provider should authenticate with
	// the endpoint.
	Enabled bool `json:"enabled"`

	// HTTPHeaderAPIKey is the API-key that will be set under the X-Api-Key header
	APIKey string `json:"apiKey"`

	// APIKeyHeader is the header that will be used to set the API key.
	APIKeyHeader string `json:"apiKeyHeader"`
}

// ValidateBasic performs basic validation of the API authentication.
func (a Authentication) ValidateBasic() error {
	if !a.Enabled {
		return nil
	}

	if len(a.APIKey) == 0 {
		return fmt.Errorf("authentication http header api key cannot be empty")
	}

	if len(a.APIKeyHeader) == 0 {
		return fmt.Errorf("authentication api key header cannot be empty")
	}

	return nil
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
		return fmt.Errorf("provider url and endpoints cannot be empty")
	}

	if len(c.Name) == 0 {
		return fmt.Errorf("provider name cannot be empty")
	}

	if c.BatchSize > 0 && c.Atomic {
		return fmt.Errorf("batch size cannot be set for atomic providers")
	}

	for _, e := range c.Endpoints {
		if err := e.ValidateBasic(); err != nil {
			return err
		}
	}

	return nil
}
