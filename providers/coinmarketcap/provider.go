package coinmarketcap

import (
	"context"
	"sync"

	"cosmossdk.io/log"
	"github.com/skip-mev/slinky/oracle/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// Name is the name of this provider
	Name = "coinmarketcap"
)

var _ types.Provider = (*Provider)(nil)

// Provider is the implementation of the oracle's Provider interface for coinmarketcap.
type Provider struct {
	// set of pairs to query prices for
	pairs []oracletypes.CurrencyPair

	// logger
	logger log.Logger

	// api-key is the api-key accompanying requests to the coinmarketcap api.
	apiKey string

	// TokenNameToMetadata is a map of currency pairs to their metadata.
	tokenNameToMetadata map[string]types.TokenMetadata
}

// NewProvider returns a new coinmarketcap provider. It uses the provided API-key in the header of outgoing requests to coinmarketcap's API.
func NewProvider(logger log.Logger, pairs []oracletypes.CurrencyPair, apiKey string, tokenNameToMetadata map[string]types.TokenMetadata) *Provider {
	return &Provider{
		pairs:               pairs,
		logger:              logger,
		apiKey:              apiKey,
		tokenNameToMetadata: tokenNameToMetadata,
	}
}

// GetPrices returns the current set of prices for each of the currency pairs. This method starts all
// price requests concurrently, and waits for them all to finish, or for the context to be cancelled,
// at which point it aggregates the responses and returns.
func (p *Provider) GetPrices(ctx context.Context) (map[oracletypes.CurrencyPair]types.QuotePrice, error) {
	type priceData struct {
		types.QuotePrice
		oracletypes.CurrencyPair
	}

	// create response channel
	resp := make(chan priceData, len(p.pairs))

	wg := sync.WaitGroup{}
	wg.Add(len(p.pairs))

	// fan-out requests to coinmarketcap api
	for _, currencyPair := range p.pairs {
		go func(pair oracletypes.CurrencyPair) {
			defer wg.Done()

			// get price
			qp, err := p.getPriceForPair(ctx, pair)
			if err != nil {
				p.logger.Error("failed to get price for pair", "provider", p.Name(), "pair", pair, "err", err)
			} else {
				p.logger.Info("Fetched price for pair", "pair", pair, "provider", p.Name())

				// send price to response channel
				resp <- priceData{
					qp,
					pair,
				}
			}
		}(currencyPair)
	}

	// close response channel when all requests have been processed, or if context is cancelled
	go func() {
		defer close(resp)

		select {
		case <-ctx.Done():
			return
		case <-finish(&wg):
			return
		}
	}()

	// fan-in
	prices := make(map[oracletypes.CurrencyPair]types.QuotePrice)
	for price := range resp {
		prices[price.CurrencyPair] = price.QuotePrice
	}

	return prices, nil
}

// Name returns the name of the provider.
func (p *Provider) Name() string {
	return Name
}

// SetPairs sets the currency pairs that the provider will fetch prices for.
func (p *Provider) SetPairs(pairs ...oracletypes.CurrencyPair) {
	p.pairs = pairs
}

// GetPairs returns the currency pairs that the provider is fetching prices for.
func (p *Provider) GetPairs() []oracletypes.CurrencyPair {
	return p.pairs
}

// getSymbolForPair returns the symbol for a currency pair.
func (p *Provider) getSymbolForTokenName(tokenName string) string {
	if metadata, ok := p.tokenNameToMetadata[tokenName]; ok {
		return metadata.Symbol
	}

	return tokenName
}

// finish takes a wait-group, and returns a channel that is sent on when the
// Waitgroup is finished.
func finish(wg *sync.WaitGroup) <-chan struct{} {
	ch := make(chan struct{})

	// non-blocing wait for waitgroup to finish, and return channel
	go func() {
		wg.Wait()
		close(ch)
	}()
	return ch
}
