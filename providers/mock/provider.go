package mock

import (
	"time"

	"github.com/holiman/uint256"
	"github.com/skip-mev/slinky/oracle/types"
)

var _ types.Provider = (*NormalMockProvider)(nil)

type (
	// NormalMockProvider defines a mocked exchange rate provider using fixed exchange
	// rates.
	NormalMockProvider struct {
		exchangeRates map[types.CurrencyPair]types.QuotePrice
		currencyPairs []types.CurrencyPair
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
		exchangeRates: map[types.CurrencyPair]types.QuotePrice{
			types.NewCurrencyPair("COSMOS", "USDC", 6):   {Price: uint256.NewInt(1134)},
			types.NewCurrencyPair("COSMOS", "USDT", 6):   {Price: uint256.NewInt(1135)},
			types.NewCurrencyPair("COSMOS", "USD", 6):    {Price: uint256.NewInt(1136)},
			types.NewCurrencyPair("OSMOSIS", "USDC", 6):  {Price: uint256.NewInt(1137)},
			types.NewCurrencyPair("OSMOSIS", "USDT", 6):  {Price: uint256.NewInt(1138)},
			types.NewCurrencyPair("OSMOSIS", "USD", 6):   {Price: uint256.NewInt(1139)},
			types.NewCurrencyPair("ETHEREUM", "USDC", 6): {Price: uint256.NewInt(1140)},
			types.NewCurrencyPair("ETHEREUM", "USDT", 6): {Price: uint256.NewInt(1141)},
			types.NewCurrencyPair("ETHEREUM", "USD", 6):  {Price: uint256.NewInt(1142)},
			types.NewCurrencyPair("BITCOIN", "USD", 6):   {Price: uint256.NewInt(1143)},
		},
		currencyPairs: []types.CurrencyPair{
			types.NewCurrencyPair("COSMOS", "USDC", 6),
			types.NewCurrencyPair("COSMOS", "USDT", 6),
			types.NewCurrencyPair("COSMOS", "USD", 6),
			types.NewCurrencyPair("OSMOSIS", "USDC", 6),
			types.NewCurrencyPair("OSMOSIS", "USDT", 6),
			types.NewCurrencyPair("OSMOSIS", "USD", 6),
			types.NewCurrencyPair("ETHEREUM", "USDC", 6),
			types.NewCurrencyPair("ETHEREUM", "USDT", 6),
			types.NewCurrencyPair("ETHEREUM", "USD", 6),
			types.NewCurrencyPair("BITCOIN", "USD", 6),
		},
	}
}

// Name returns the name of the mock provider.
func (p NormalMockProvider) Name() string {
	return "mock-provider"
}

// GetPrices returns the mocked exchange rates.
func (p NormalMockProvider) GetPrices() (map[types.CurrencyPair]types.QuotePrice, error) {
	return p.exchangeRates, nil
}

// SetPairs is a no-op for the mock provider.
func (p NormalMockProvider) SetPairs(_ ...types.CurrencyPair) {}

// GetPairs is a no-op for the mock provider.
func (p NormalMockProvider) GetPairs() []types.CurrencyPair {
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
func (p FailingMockProvider) GetPrices() (map[types.CurrencyPair]types.QuotePrice, error) {
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
func (p TimeoutMockProvider) GetPrices() (map[types.CurrencyPair]types.QuotePrice, error) {
	time.Sleep(1*time.Second + p.timeout)

	panic("mock provider should always times out")
}
