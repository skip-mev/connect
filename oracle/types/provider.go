package types

import (
	"time"

	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/x/oracle/types"
	"golang.org/x/net/context"
)

// Provider defines an interface an exchange price provider must implement.
//
//go:generate mockery --name Provider --filename mock_provider.go
type Provider interface {
	// Name returns the name of the provider.
	Name() string

	// GetPrices returns the aggregated prices based on the provided currency pairs.
	GetPrices(context.Context) (map[types.CurrencyPair]aggregator.QuotePrice, error)

	// SetPairs sets the pairs that the provider should fetch prices for.
	SetPairs(...types.CurrencyPair)

	// GetPairs returns the pairs that the provider is fetching prices for.
	GetPairs() []types.CurrencyPair
}

type ProviderConfig struct {
	// Name identifies which provider this config is for
	Name string `mapstructure:"name" toml:"name"`

	// Apikey is the api-key accompanying requests to the provider's API.
	Apikey string `mapstructure:"apikey" toml:"apikey"`

	// TokenNameToMetadata is a map of token names to their metadata, i.e how each token in the CurrencyPair should
	// map to the token references in the queried provider's API, and required decimals for on-chain data sources.
	TokenNameToMetadata map[string]TokenMetadata `mapstructure:"token_name_to_metadata" toml:"token_name_to_metadata"`

	// PairToContractAddress is a map of pairs to the address of a contract that provides the price for that pair
	// for providers that use on-chain data. The key of the outer map is the quote token, the key of the inner map
	// is the base token.
	PairToContractAddress map[string]map[string]string `mapstructure:"pair_to_contract_address"`

	// ProviderTimeout is the maximum amount of time to wait for a response from the provider.
	ProviderTimeout time.Duration `mapstructure:"provider_timeout"`
}
