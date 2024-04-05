package types

import (
	"fmt"
)

// DefaultTickerDecimals is the number of decimal places every single price
// is scaled to before being sent for aggregation.
const DefaultTickerDecimals = 18

// ProviderTicker is the interface for the ticker that provider's utilize/return.
type (
	ProviderTicker interface {
		fmt.Stringer

		// Provider returns the provider for the ticker.
		Provider() string
		// OnChainTicker returns the on-chain representation for the ticker.
		OnChainTicker() string
		// OffChainTicker returns the off-chain representation for the ticker.
		OffChainTicker() string
		// Decimals returns the number of decimals for the ticker.
		Decimals() uint64
		// JSON returns additional JSON data for the ticker.
		JSON() string
	}

	// DefaultProviderTicker is a basic implementation of the ProviderTicker interface.
	DefaultProviderTicker struct {
		provider string
		onChain  string
		offChain string
		decimals uint64
		json     string
	}

	// ProviderTickers is a type alias for a list of provider tickers.
	ProviderTickers struct {
		tickers []ProviderTicker
		cache   map[string]ProviderTicker
	}
)

// NewProviderTicker returns a new provider ticker.
func NewProviderTicker(
	provider,
	onChain,
	json,
	offChain string,
	decimals uint64,
) ProviderTicker {
	return DefaultProviderTicker{
		provider: provider,
		onChain:  onChain,
		offChain: offChain,
		decimals: decimals,
		json:     json,
	}
}

// Provider returns the provider for the ticker.
func (t DefaultProviderTicker) Provider() string {
	return t.provider
}

// OnChainTicker returns the on-chain representation for the ticker.
func (t DefaultProviderTicker) OnChainTicker() string {
	return t.onChain
}

// OffChainTicker returns the off-chain representation for the ticker.
func (t DefaultProviderTicker) OffChainTicker() string {
	return t.offChain
}

// Decimals returns the number of decimals for the ticker.
func (t DefaultProviderTicker) Decimals() uint64 {
	return t.decimals
}

// JSON returns additional JSON data for the ticker.
func (t DefaultProviderTicker) JSON() string {
	return t.json
}

// String returns the string representation of the provider ticker.
func (t DefaultProviderTicker) String() string {
	return fmt.Sprintf(
		"provider: %s, on-chain-ticker: %s, off-chain-ticker: %s, decimals: %d",
		t.provider,
		t.onChain,
		t.offChain,
		t.decimals,
	)
}

// NewProviderTickers returns a new list of provider tickers.
func NewProviderTickers(tickers ...ProviderTicker) ProviderTickers {
	cache := make(map[string]ProviderTicker)
	for _, ticker := range tickers {
		cache[ticker.OffChainTicker()] = ticker
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
	t.cache[ticker.OffChainTicker()] = ticker
}

// Reset resets the provider tickers.
func (t *ProviderTickers) Reset() {
	t.cache = make(map[string]ProviderTicker)
}
