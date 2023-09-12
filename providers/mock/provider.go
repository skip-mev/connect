package mock

import (
	"context"
	"strconv"
	"time"

	"github.com/holiman/uint256"
	"github.com/skip-mev/slinky/oracle/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var _ types.Provider = (*StaticMockProvider)(nil)

type (
	// StaticMockProvider defines a mocked exchange rate provider using fixed exchange
	// rates.
	StaticMockProvider struct {
		exchangeRates map[oracletypes.CurrencyPair]types.QuotePrice
		currencyPairs []oracletypes.CurrencyPair
	}

	// FailingMockProvider defines a mocked exchange rate provider that always
	// fails when fetching prices.
	FailingMockProvider struct {
		*StaticMockProvider
	}

	// TimeoutMockProvider defines a mocked exchange rate provider that always
	// times out when fetching prices.
	TimeoutMockProvider struct {
		*StaticMockProvider
		timeout time.Duration
	}
)

// NewStaticMockProvider returns a new mock provider. The mock provider
// will always return the same static data. Meant to be used for testing.
func NewStaticMockProvider() *StaticMockProvider {
	return &StaticMockProvider{
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

// NewStaticMockProviderFromConfig constructs a new static mock provider from the config
// Notice this method expects the TokenNameToSymbol map to be populated w/ entries of the form
// CurrencyPair.ToString(): uint256.NewInt(price)
func NewStaticMockProviderFromConfig(config types.ProviderConfig) *StaticMockProvider {
	s := StaticMockProvider{
		exchangeRates: make(map[oracletypes.CurrencyPair]types.QuotePrice),
		currencyPairs: make([]oracletypes.CurrencyPair, 0),
	}

	for cpString, metadata := range config.TokenNameToMetadata {
		cp, err := oracletypes.CurrencyPairFromString(cpString)
		if err != nil {
			continue
		}

		priceString := metadata.Symbol
		priceInt, err := strconv.Atoi(priceString)
		if err != nil {
			continue
		}

		s.exchangeRates[cp] = types.QuotePrice{Price: uint256.NewInt(uint64(priceInt))}
		s.currencyPairs = append(s.currencyPairs, cp)
	}

	return &s
}

// Name returns the name of the mock provider.
func (p StaticMockProvider) Name() string {
	return "static-mock-provider"
}

// GetPrices returns the mocked exchange rates.
func (p StaticMockProvider) GetPrices(_ context.Context) (map[oracletypes.CurrencyPair]types.QuotePrice, error) {
	return p.exchangeRates, nil
}

// SetPairs is a no-op for the mock provider.
func (p StaticMockProvider) SetPairs(_ ...oracletypes.CurrencyPair) {}

// GetPairs is a no-op for the mock provider.
func (p StaticMockProvider) GetPairs() []oracletypes.CurrencyPair {
	return p.currencyPairs
}

var _ types.Provider = (*FailingMockProvider)(nil)

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
func (p FailingMockProvider) GetPrices(_ context.Context) (map[oracletypes.CurrencyPair]types.QuotePrice, error) {
	panic("mock provider always fails")
}

var _ types.Provider = (*TimeoutMockProvider)(nil)

// NewTimeoutMockProvider returns a new timeout mock provider.
func NewTimeoutMockProvider(timeout time.Duration) *TimeoutMockProvider {
	return &TimeoutMockProvider{
		StaticMockProvider: NewStaticMockProvider(),
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
