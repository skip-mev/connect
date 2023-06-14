package coinbase

import (
	"github.com/cometbft/cometbft/libs/log"
	"github.com/skip-mev/slinky/oracle/types"
)

const (
	// Name is the name of the provider.
	Name = "coinbase"
)

var _ types.Provider = (*Provider)(nil)

// Provider implements the Provider interface for Coinbase. This provider
// is a very simple implementation that fetches spot prices from the Coinbase API.
type Provider struct {
	pairs  []types.CurrencyPair
	logger log.Logger
}

// NewProvider returns a new Coinbase provider.
//
// THIS PROVIDER SHOULD NOT BE USED IN PRODUCTION. IT IS ONLY MEANT FOR TESTING.
func NewProvider(logger log.Logger, pairs []types.CurrencyPair) *Provider {
	return &Provider{
		pairs:  pairs,
		logger: logger,
	}
}

// Name returns the name of the provider.
func (p *Provider) Name() string {
	return Name
}

// GetPrices returns the current set of prices for each of the currency pairs. The
// prices are fetched from the Coinbase API. The price is returned is the spot price
// for the given currency pair.
func (p *Provider) GetPrices() (map[types.CurrencyPair]types.TickerPrice, error) {
	resp := make(map[types.CurrencyPair]types.TickerPrice)

	for _, currencyPair := range p.pairs {
		spotPrice, err := getPriceForPair(currencyPair)
		if err != nil {
			p.logger.Error(
				p.Name(),
				"failed to get price for pair", currencyPair,
				"err", err,
			)
			continue
		}

		resp[currencyPair] = *spotPrice
	}

	return resp, nil
}

// SetPairs sets the currency pairs that the provider will fetch prices for.
func (p *Provider) SetPairs(pairs ...types.CurrencyPair) {
	p.pairs = pairs
}

// GetPairs returns the currency pairs that the provider is fetching prices for.
func (p *Provider) GetPairs() []types.CurrencyPair {
	return p.pairs
}
