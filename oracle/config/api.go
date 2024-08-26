package config

import (
	"fmt"
	"time"
)

// APIConfig defines a config for an API based data provider.
type APIConfig struct {
	// Enabled indicates if the provider is enabled.
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

	// Endpoints is a list of endpoints that the provider can query.
	Endpoints []Endpoint `json:"endpoints"`

	// BatchSize is the maximum number of IDs that the provider can query in a single
	// request. This parameter must be 0 for atomic providers. Otherwise, the effective
	// value will be max(1, BatchSize). Notice, if numCPs > batchSize * maxQueries then
	// some currency-pairs may not be fetched each interval.
	BatchSize int `json:"batchSize"`

	// Name is the name of the provider that corresponds to this config.
	Name string `json:"name"`

	// MaxBlockHeightAge is the oldest an update from an on-chain data source can be without having its
	// block height incremented.  In the case where a data source has exceeded this limit and the block
	// height is not increasing, price reporting will be skipped until the block height increases.
	MaxBlockHeightAge time.Duration `json:"maxBlockHeightAge"`
}

// Endpoint holds all data necessary for an API provider to connect to a given endpoint
// i.e. URL, headers, authentication, etc.
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
	// HTTPHeaderAPIKey is the API-key that will be set under the X-Api-Key header
	APIKey string `json:"apiKey"`

	// APIKeyHeader is the header that will be used to set the API key.
	APIKeyHeader string `json:"apiKeyHeader"`
}

// Enabled returns true if the authentication is enabled.
func (a Authentication) Enabled() bool {
	return a.APIKey != "" && a.APIKeyHeader != ""
}

// ValidateBasic performs basic validation of the API authentication. Specifically, the APIKey + APIKeyHeader
// must be set atomically.
func (a Authentication) ValidateBasic() error {
	if a.APIKey != "" && a.APIKeyHeader == "" {
		return fmt.Errorf("api key header cannot be empty when api key is set")
	}

	if a.APIKey == "" && a.APIKeyHeader != "" {
		return fmt.Errorf("api key cannot be empty when api key header is set")
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

	if len(c.Name) == 0 {
		return fmt.Errorf("provider name cannot be empty")
	}

	if c.BatchSize > 0 && c.Atomic {
		return fmt.Errorf("batch size cannot be set for atomic providers")
	}

	if len(c.Endpoints) == 0 {
		return fmt.Errorf("endpoints cannot be empty")
	}

	for _, e := range c.Endpoints {
		if err := e.ValidateBasic(); err != nil {
			return err
		}
	}

	if c.MaxBlockHeightAge < 0 {
		return fmt.Errorf("max_block_height_age cannot be negative")
	}

	return nil
}
