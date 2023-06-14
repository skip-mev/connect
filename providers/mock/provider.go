package mock

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/skip-mev/slinky/oracle/types"
)

var _ types.Provider = (*MockProvider)(nil)

type (
	// MockProvider defines a mocked exchange rate provider using fixed exchange
	// rates.
	MockProvider struct {
		exchangeRates map[string]types.TickerPrice
		currencyPairs []types.CurrencyPair
	}

	// FailingMockProvider defines a mocked exchange rate provider that always
	// fails when fetching prices.
	FailingMockProvider struct {
		*MockProvider
	}

	// TimeoutMockProvider defines a mocked exchange rate provider that always
	// times out when fetching prices.
	TimeoutMockProvider struct {
		*MockProvider
	}
)

// NewMockProvider returns a new mock provider. The mock provider
// will always return the same static data. Meant to be used for testing.
func NewMockProvider() *MockProvider {
	return &MockProvider{
		exchangeRates: map[string]types.TickerPrice{
			"ATOM/USDC": {Price: sdk.MustNewDecFromStr("11.34")},
			"ATOM/USDT": {Price: sdk.MustNewDecFromStr("11.36")},
			"ATOM/USD":  {Price: sdk.MustNewDecFromStr("11.35")},
			"OSMO/USDC": {Price: sdk.MustNewDecFromStr("1.34")},
			"OSMO/USDT": {Price: sdk.MustNewDecFromStr("1.36")},
			"OSMO/USD":  {Price: sdk.MustNewDecFromStr("1.35")},
			"WETH/USDC": {Price: sdk.MustNewDecFromStr("1560.34")},
			"WETH/USDT": {Price: sdk.MustNewDecFromStr("1560.36")},
			"WETH/USD":  {Price: sdk.MustNewDecFromStr("1560.35")},
			"BTC/USD":   {Price: sdk.MustNewDecFromStr("50000.00")},
		},
		currencyPairs: []types.CurrencyPair{
			{Base: "ATOM", Quote: "USDC"},
			{Base: "ATOM", Quote: "USDT"},
			{Base: "ATOM", Quote: "USD"},
			{Base: "OSMO", Quote: "USDC"},
			{Base: "OSMO", Quote: "USDT"},
			{Base: "OSMO", Quote: "USD"},
			{Base: "WETH", Quote: "USDC"},
			{Base: "WETH", Quote: "USDT"},
			{Base: "WETH", Quote: "USD"},
			{Base: "BTC", Quote: "USD"},
		},
	}
}

// Name returns the name of the mock provider.
func (p MockProvider) Name() string {
	return "mock-provider"
}

// GetPrices returns the mocked exchange rates.
func (p MockProvider) GetPrices() (map[string]types.TickerPrice, error) {
	return p.exchangeRates, nil
}

// SetPairs is a no-op for the mock provider.
func (p MockProvider) SetPairs(pairs ...types.CurrencyPair) {}

// GetPairs is a no-op for the mock provider.
func (p MockProvider) GetPairs() []types.CurrencyPair {
	return p.currencyPairs
}

// NewFailingMockProvider returns a new failing mock provider.
func NewFailingMockProvider() *FailingMockProvider {
	return &FailingMockProvider{
		MockProvider: NewMockProvider(),
	}
}

// Name returns the name of the failing mock provider.
func (p FailingMockProvider) Name() string {
	return "failing-mock-provider"
}

// GetPrices always fails for the failing mock provider.
func (p FailingMockProvider) GetPrices() (map[string]types.TickerPrice, error) {
	panic("mock provider always fails")
}

// NewTimeoutMockProvider returns a new timeout mock provider.
func NewTimeoutMockProvider() *TimeoutMockProvider {
	return &TimeoutMockProvider{
		MockProvider: NewMockProvider(),
	}
}

// Name returns the name of the timeout mock provider.
func (p TimeoutMockProvider) Name() string {
	return "timeout-mock-provider"
}

// GetPrices always times out for the timeout mock provider.
func (p TimeoutMockProvider) GetPrices() (map[string]types.TickerPrice, error) {
	time.Sleep(1000 * time.Second)

	panic("mock provider should always times out")
}
