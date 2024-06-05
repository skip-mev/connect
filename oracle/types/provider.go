package types

import (
	"fmt"
	"strings"
	"sync"
)

type (
	// ProviderTicker is the interface for the ticker that provider's utilize/return.
	ProviderTicker interface {
		fmt.Stringer

		// GetOffChainTicker returns the off-chain representation for the ticker.
		GetOffChainTicker() string
		// GetJSON returns additional JSON data for the ticker.
		GetJSON() string
	}

	// DefaultProviderTicker is a basic implementation of the ProviderTicker interface.
	// Provider's that utilize this implementation should be able to easily configure
	// custom json data for their tickers.
	DefaultProviderTicker struct {
		OffChainTicker string
		JSON           string
	}

	// ProviderTickers is a thread safe helper struct to manage a list of provider tickers.
	ProviderTickers struct {
		mut sync.Mutex

		cache map[string]ProviderTicker
	}
)

// NewProviderTicker returns a new provider ticker.
func NewProviderTicker(
	offChain, json string,
) ProviderTicker {
	return DefaultProviderTicker{
		OffChainTicker: offChain,
		JSON:           json,
	}
}

// GetOffChainTicker returns the off-chain representation for the ticker.
func (t DefaultProviderTicker) GetOffChainTicker() string {
	return t.OffChainTicker
}

// GetJSON returns additional JSON data for the ticker.
func (t DefaultProviderTicker) GetJSON() string {
	return t.JSON
}

// String returns the string representation of the provider ticker.
func (t DefaultProviderTicker) String() string {
	return t.OffChainTicker
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

// FromOffChainTicker returns the provider ticker from the off-chain ticker.
func (t *ProviderTickers) FromOffChainTicker(offChain string) (ProviderTicker, bool) {
	t.mut.Lock()
	defer t.mut.Unlock()

	ticker, ok := t.cache[offChain]
	return ticker, ok
}

// Add adds a provider ticker to the list of provider tickers.
func (t *ProviderTickers) Add(ticker ProviderTicker) {
	t.mut.Lock()
	defer t.mut.Unlock()

	t.cache[strings.ToLower(ticker.GetOffChainTicker())] = ticker
	t.cache[ticker.GetOffChainTicker()] = ticker
	t.cache[strings.ToUpper(ticker.GetOffChainTicker())] = ticker
}
