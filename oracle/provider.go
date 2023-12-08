package oracle

import (
	"math/big"

	"cosmossdk.io/log"
	"golang.org/x/net/context"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/x/oracle/types"
)

// Provider defines an interface an exchange price provider must implement.
//
//go:generate mockery --name Provider --filename mock_provider.go
type Provider interface {
	// Name returns the name of the provider.
	Name() string

	// GetPrices returns the aggregated prices based on the provided currency pairs.
	GetPrices(context.Context) (map[types.CurrencyPair]*big.Int, error)

	// SetPairs sets the pairs that the provider should fetch prices for.
	SetPairs(...types.CurrencyPair)

	// GetPairs returns the pairs that the provider is fetching prices for.
	GetPairs() []types.CurrencyPair
}

// ProviderFactory inputs the oracle configuration and returns a set of providers. Developers
// can implement their own provider factory to create their own providers.
type ProviderFactory func(log.Logger, config.OracleConfig) ([]Provider, error)
