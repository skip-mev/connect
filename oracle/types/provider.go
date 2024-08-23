package types

import (
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	providertypes "github.com/skip-mev/connect/v2/providers/types"
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
		cache[strings.ToLower(ticker.GetOffChainTicker())] = ticker
		cache[ticker.GetOffChainTicker()] = ticker
		cache[strings.ToUpper(ticker.GetOffChainTicker())] = ticker
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

// NoPriceChangeResponse is used to handle a message that indicates that the price has not changed.
// In particular, this will update the base provider with the ResponseCodeUnchanged code for all tickers.
func (t *ProviderTickers) NoPriceChangeResponse() PriceResponse {
	resolved := make(ResolvedPrices)
	seen := make(map[ProviderTicker]struct{})
	for _, ticker := range t.cache {
		if _, ok := seen[ticker]; ok {
			continue
		}

		resolved[ticker] = NewPriceResultWithCode(
			big.NewFloat(0),
			time.Now().UTC(),
			providertypes.ResponseCodeUnchanged,
		)
		seen[ticker] = struct{}{}
	}
	return NewPriceResponse(resolved, nil)
}
