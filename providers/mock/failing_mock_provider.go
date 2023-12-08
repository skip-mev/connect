package mock

import (
	"context"
	"math/big"

	"github.com/skip-mev/slinky/oracle"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var _ oracle.Provider = (*FailingMockProvider)(nil)

type (
	// FailingMockProvider defines a mocked exchange rate provider that always
	// fails when fetching prices.
	FailingMockProvider struct {
		*StaticMockProvider
	}
)

// NewFailingMockProvider returns a new failing mock provider.
func NewFailingMockProvider() *FailingMockProvider {
	return &FailingMockProvider{
		StaticMockProvider: NewStaticMockProvider(),
	}
}

// Name returns the name of the failing mock provider.
func (p FailingMockProvider) Name() string {
	return "failing-mock-provider"
}

// GetPrices always fails for the failing mock provider.
func (p FailingMockProvider) GetPrices(_ context.Context) (map[oracletypes.CurrencyPair]*big.Int, error) {
	panic("mock provider always fails")
}
