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
		exchangeRates map[types.CurrencyPair]types.TickerPrice
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
		timeout time.Duration
	}
)

// NewMockProvider returns a new mock provider. The mock provider
// will always return the same static data. Meant to be used for testing.
func NewMockProvider() *MockProvider {
	return &MockProvider{
		exchangeRates: map[types.CurrencyPair]types.TickerPrice{
			{Base: "ATOM", Quote: "USDC"}: {Price: sdk.MustNewDecFromStr("11.34")},
			{Base: "ATOM", Quote: "USDT"}: {Price: sdk.MustNewDecFromStr("11.36")},
			{Base: "ATOM", Quote: "USD"}:  {Price: sdk.MustNewDecFromStr("11.35")},
			{Base: "OSMO", Quote: "USDC"}: {Price: sdk.MustNewDecFromStr("1.34")},
			{Base: "OSMO", Quote: "USDT"}: {Price: sdk.MustNewDecFromStr("1.36")},
			{Base: "OSMO", Quote: "USD"}:  {Price: sdk.MustNewDecFromStr("1.35")},
			{Base: "WETH", Quote: "USDC"}: {Price: sdk.MustNewDecFromStr("1560.34")},
			{Base: "WETH", Quote: "USDT"}: {Price: sdk.MustNewDecFromStr("1560.36")},
			{Base: "WETH", Quote: "USD"}:  {Price: sdk.MustNewDecFromStr("1560.35")},
			{Base: "BTC", Quote: "USD"}:   {Price: sdk.MustNewDecFromStr("50000.00")},
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
func (p MockProvider) GetPrices() (map[types.CurrencyPair]types.TickerPrice, error) {
	return p.exchangeRates, nil
}

// SetPairs is a no-op for the mock provider.
func (p MockProvider) SetPairs(pairs ...types.CurrencyPair) {}

// GetPairs is a no-op for the mock provider.
func (p MockProvider) GetPairs() []types.CurrencyPair {
	return p.currencyPairs
}

var _ types.Provider = (*FailingMockProvider)(nil)

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
func (p FailingMockProvider) GetPrices() (map[types.CurrencyPair]types.TickerPrice, error) {
	panic("mock provider always fails")
}

var _ types.Provider = (*TimeoutMockProvider)(nil)

// NewTimeoutMockProvider returns a new timeout mock provider.
func NewTimeoutMockProvider(timeout time.Duration) *TimeoutMockProvider {
	return &TimeoutMockProvider{
		MockProvider: NewMockProvider(),
		timeout:      timeout,
	}
}

// Name returns the name of the timeout mock provider.
func (p TimeoutMockProvider) Name() string {
	return "timeout-mock-provider"
}

// GetPrices always times out for the timeout mock provider.
func (p TimeoutMockProvider) GetPrices() (map[types.CurrencyPair]types.TickerPrice, error) {
	time.Sleep(1*time.Second + p.timeout)

	panic("mock provider should always times out")
}
