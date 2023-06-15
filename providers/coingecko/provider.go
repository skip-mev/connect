package coingecko

import (
	"strings"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/skip-mev/slinky/oracle/types"
)

const (
	// Name is the name of the provider.
	Name = "coingecko"
)

var _ types.Provider = (*Provider)(nil)

// Provider implements the Provider interface for CoinGecko. This provider
// is a very simple implementation that fetches prices from the CoinGecko API.
type Provider struct {
	pairs  []types.CurrencyPair
	logger log.Logger

	// bases is a list of base currencies that the provider should fetch
	// prices for.
	bases string

	// quotes is a list of quote currencies that the provider should fetch
	// prices for.
	quotes string
}

// NewProvider returns a new CoinGecko provider.
func NewProvider(logger log.Logger, pairs []types.CurrencyPair) *Provider {
	seenQuotes := make(map[string]struct{})
	quotes := make([]string, 0)

	seenBases := make(map[string]struct{})
	bases := make([]string, 0)

	for _, pair := range pairs {
		if _, ok := seenQuotes[pair.Quote]; !ok {
			seenQuotes[pair.Quote] = struct{}{}
			quotes = append(quotes, pair.Quote)
		}

		if _, ok := seenBases[pair.Base]; !ok {
			seenBases[pair.Base] = struct{}{}
			bases = append(bases, pair.Base)
		}
	}

	return &Provider{
		pairs:  pairs,
		logger: logger,
		bases:  strings.Join(bases, ","),
		quotes: strings.Join(quotes, ","),
	}
}

// Name returns the name of the provider.
func (p *Provider) Name() string {
	return Name
}

// GetPrices returns the current set of prices for each of the currency pairs.
func (p *Provider) GetPrices() (map[types.CurrencyPair]types.TickerPrice, error) {
	return p.getPrices()
}

// SetPairs sets the currency pairs that the provider will fetch prices for.
func (p *Provider) SetPairs(pairs ...types.CurrencyPair) {
	p.pairs = pairs
}

// GetPairs returns the currency pairs that the provider is fetching prices for.
func (p *Provider) GetPairs() []types.CurrencyPair {
	return p.pairs
}
