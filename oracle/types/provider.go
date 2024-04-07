package types

import (
	"fmt"
)

// DefaultTickerDecimals is the number of decimal places every single price
// is scaled to before being sent for aggregation.
const DefaultTickerDecimals = 18

type (
	// ProviderTicker is the interface for the ticker that provider's utilize/return.
	ProviderTicker interface {
		fmt.Stringer

		// GetProvider returns the provider for the ticker.
		GetProvider() string
		// GetOffChainTicker returns the off-chain representation for the ticker.
		GetOffChainTicker() string
		// GetDecimals returns the number of decimals for the ticker.
		GetDecimals() uint64
		// GetJSON returns additional JSON data for the ticker.
		GetJSON() string
	}

	// DefaultProviderTicker is a basic implementation of the ProviderTicker interface.
	// Provider's that utilize this implementation should be able to easily configure
	// custom json data for their tickers.
	DefaultProviderTicker struct {
		Name           string
		OffChainTicker string
		Decimals       uint64
		JSON           string
	}

	// ProviderTickers is helper struct to manage a list of provider tickers.
	ProviderTickers struct {
		cache map[string]ProviderTicker
	}
)

// NewProviderTicker returns a new provider ticker.
func NewProviderTicker(
	provider, offChain, json string,
	decimals uint64,
) ProviderTicker {
	return DefaultProviderTicker{
		Name:           provider,
		OffChainTicker: offChain,
		Decimals:       decimals,
		JSON:           json,
	}
}

// Provider returns the provider for the ticker.
func (t DefaultProviderTicker) GetProvider() string {
	return t.Name
}

// OffChainTicker returns the off-chain representation for the ticker.
func (t DefaultProviderTicker) GetOffChainTicker() string {
	return t.OffChainTicker
}

// Decimals returns the number of decimals for the ticker.
func (t DefaultProviderTicker) GetDecimals() uint64 {
	return t.Decimals
}

// JSON returns additional JSON data for the ticker.
func (t DefaultProviderTicker) GetJSON() string {
	return t.JSON
}

// String returns the string representation of the provider ticker.
func (t DefaultProviderTicker) String() string {
	return fmt.Sprintf(
		"provider: %s, off-chain-ticker: %s, decimals: %d, json: %s",
		t.Name,
		t.OffChainTicker,
		t.Decimals,
		t.JSON,
	)
}

// NewProviderTickers returns a new list of provider tickers.
func NewProviderTickers(tickers ...ProviderTicker) ProviderTickers {
	cache := make(map[string]ProviderTicker)
	for _, ticker := range tickers {
		cache[ticker.GetOffChainTicker()] = ticker
	}
	return ProviderTickers{
		cache: cache,
	}
}

// FromOffChain returns the provider ticker from the off-chain ticker.
func (t ProviderTickers) FromOffChain(offChain string) (ProviderTicker, bool) {
	ticker, ok := t.cache[offChain]
	return ticker, ok
}

// Add adds a provider ticker to the list of provider tickers.
func (t *ProviderTickers) Add(ticker ProviderTicker) {
	t.cache[ticker.GetOffChainTicker()] = ticker
}

// Reset resets the provider tickers.
func (t *ProviderTickers) Reset() {
	t.cache = make(map[string]ProviderTicker)
}
