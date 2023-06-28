package mock

import (
	"context"
	"time"

	"github.com/holiman/uint256"
	"github.com/skip-mev/slinky/oracle/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var _ types.Provider = (*NormalMockProvider)(nil)

type (
	// NormalMockProvider defines a mocked exchange rate provider using fixed exchange
	// rates.
	NormalMockProvider struct {
		exchangeRates map[oracletypes.CurrencyPair]types.QuotePrice
		currencyPairs []oracletypes.CurrencyPair
	}

	// FailingMockProvider defines a mocked exchange rate provider that always
	// fails when fetching prices.
	FailingMockProvider struct {
		*NormalMockProvider
	}

	// TimeoutMockProvider defines a mocked exchange rate provider that always
	// times out when fetching prices.
	TimeoutMockProvider struct {
		*NormalMockProvider
		timeout time.Duration
	}
)

// NewMockProvider returns a new mock provider. The mock provider
// will always return the same static data. Meant to be used for testing.
func NewMockProvider() *NormalMockProvider {
	return &NormalMockProvider{
		exchangeRates: map[oracletypes.CurrencyPair]types.QuotePrice{
			oracletypes.NewCurrencyPair("COSMOS", "USDC"):   {Price: uint256.NewInt(1134)},
			oracletypes.NewCurrencyPair("COSMOS", "USDT"):   {Price: uint256.NewInt(1135)},
			oracletypes.NewCurrencyPair("COSMOS", "USD"):    {Price: uint256.NewInt(1136)},
			oracletypes.NewCurrencyPair("OSMOSIS", "USDC"):  {Price: uint256.NewInt(1137)},
			oracletypes.NewCurrencyPair("OSMOSIS", "USDT"):  {Price: uint256.NewInt(1138)},
			oracletypes.NewCurrencyPair("OSMOSIS", "USD"):   {Price: uint256.NewInt(1139)},
			oracletypes.NewCurrencyPair("ETHEREUM", "USDC"): {Price: uint256.NewInt(1140)},
			oracletypes.NewCurrencyPair("ETHEREUM", "USDT"): {Price: uint256.NewInt(1141)},
			oracletypes.NewCurrencyPair("ETHEREUM", "USD"):  {Price: uint256.NewInt(1142)},
			oracletypes.NewCurrencyPair("BITCOIN", "USD"):   {Price: uint256.NewInt(1143)},
		},
		currencyPairs: []oracletypes.CurrencyPair{
			oracletypes.NewCurrencyPair("COSMOS", "USDC"),
			oracletypes.NewCurrencyPair("COSMOS", "USDT"),
			oracletypes.NewCurrencyPair("COSMOS", "USD"),
			oracletypes.NewCurrencyPair("OSMOSIS", "USDC"),
			oracletypes.NewCurrencyPair("OSMOSIS", "USDT"),
			oracletypes.NewCurrencyPair("OSMOSIS", "USD"),
			oracletypes.NewCurrencyPair("ETHEREUM", "USDC"),
			oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
			oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
			oracletypes.NewCurrencyPair("BITCOIN", "USD"),
		},
	}
}

// Name returns the name of the mock provider.
func (p NormalMockProvider) Name() string {
	return "mock-provider"
}

// GetPrices returns the mocked exchange rates.
func (p NormalMockProvider) GetPrices(_ context.Context) (map[oracletypes.CurrencyPair]types.QuotePrice, error) {
	return p.exchangeRates, nil
}

// SetPairs is a no-op for the mock provider.
func (p NormalMockProvider) SetPairs(_ ...oracletypes.CurrencyPair) {}

// GetPairs is a no-op for the mock provider.
func (p NormalMockProvider) GetPairs() []oracletypes.CurrencyPair {
	return p.currencyPairs
}

var _ types.Provider = (*FailingMockProvider)(nil)

// NewFailingMockProvider returns a new failing mock provider.
func NewFailingMockProvider() *FailingMockProvider {
	return &FailingMockProvider{
		NormalMockProvider: NewMockProvider(),
	}
}

// Name returns the name of the failing mock provider.
func (p FailingMockProvider) Name() string {
	return "failing-mock-provider"
}

// GetPrices always fails for the failing mock provider.
func (p FailingMockProvider) GetPrices(_ context.Context) (map[oracletypes.CurrencyPair]types.QuotePrice, error) {
	panic("mock provider always fails")
}

var _ types.Provider = (*TimeoutMockProvider)(nil)

// NewTimeoutMockProvider returns a new timeout mock provider.
func NewTimeoutMockProvider(timeout time.Duration) *TimeoutMockProvider {
	return &TimeoutMockProvider{
		NormalMockProvider: NewMockProvider(),
		timeout:            timeout,
	}
}

// Name returns the name of the timeout mock provider.
func (p TimeoutMockProvider) Name() string {
	return "timeout-mock-provider"
}

// GetPrices always times out for the timeout mock provider.
func (p TimeoutMockProvider) GetPrices(_ context.Context) (map[oracletypes.CurrencyPair]types.QuotePrice, error) {
	time.Sleep(1*time.Second + p.timeout)

	panic("mock provider should always times out")
}
