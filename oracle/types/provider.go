package types

import (
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
	GetPrices(context.Context) (map[types.CurrencyPair]QuotePrice, error)

	// SetPairs sets the pairs that the provider should fetch prices for.
	SetPairs(...types.CurrencyPair)

	// GetPairs returns the pairs that the provider is fetching prices for.
	GetPairs() []types.CurrencyPair
}

type ProviderConfig struct {
	// Name identifies which provider this config is for
	Name string `mapstructure:"name"`

	// Apikey is the api-key accompanying requests to the provider's API.
	Apikey string `mapstructure:"apikey"`

	// TokenNameToSymbol is a map of token names to their symbols, i.e how each token in the CurrencyPair should
	// map to the token references in the queried provider's API.
	TokenNameToSymbol map[string]string `mapstructure:"token_name_to_symbol"`
}
