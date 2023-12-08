package mock

import (
	"context"
	"math/big"
	"time"

	"github.com/skip-mev/slinky/oracle"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

type (
	// TimeoutMockProvider defines a mocked exchange rate provider that always
	// times out when fetching prices.
	TimeoutMockProvider struct {
		*StaticMockProvider
		timeout time.Duration
	}
)

var _ oracle.Provider = (*TimeoutMockProvider)(nil)

// NewTimeoutMockProvider returns a new timeout mock provider.
func NewTimeoutMockProvider(timeout time.Duration) (*TimeoutMockProvider, error) {
	return &TimeoutMockProvider{
		StaticMockProvider: NewStaticMockProvider(),
		timeout:            timeout,
	}, nil
}

// Name returns the name of the timeout mock provider.
func (p TimeoutMockProvider) Name() string {
	return "timeout-mock-provider"
}

// GetPrices always times out for the timeout mock provider.
func (p TimeoutMockProvider) GetPrices(_ context.Context) (map[oracletypes.CurrencyPair]*big.Int, error) {
	time.Sleep(1*time.Second + p.timeout)

	panic("mock provider should always times out")
}
