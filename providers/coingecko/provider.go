package coingecko

import (
	"context"
	"strings"

	"cosmossdk.io/log"
	"github.com/skip-mev/slinky/oracle/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// Name is the name of the provider.
	Name = "coingecko"
)

var _ types.Provider = (*Provider)(nil)

// Provider implements the Provider interface for CoinGecko. This provider
// is a very simple implementation that fetches prices from the CoinGecko API.
type Provider struct {
	pairs  []oracletypes.CurrencyPair
	logger log.Logger

	// bases is a list of base currencies that the provider should fetch
	// prices for.
	bases string

	// quotes is a list of quote currencies that the provider should fetch
	// prices for.
	quotes string
}

// NewProvider returns a new CoinGecko provider.
func NewProvider(logger log.Logger, pairs []oracletypes.CurrencyPair) *Provider {
	bases, quotes := getUniqueBaseAndQuoteDenoms(pairs)

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
func (p *Provider) GetPrices(ctx context.Context) (map[oracletypes.CurrencyPair]types.QuotePrice, error) {
	return p.getPrices(ctx)
}

// SetPairs sets the currency pairs that the provider will fetch prices for.
func (p *Provider) SetPairs(pairs ...oracletypes.CurrencyPair) {
	bases, quotes := getUniqueBaseAndQuoteDenoms(pairs)
	p.bases = strings.Join(bases, ",")
	p.quotes = strings.Join(quotes, ",")

	p.pairs = pairs
}

// GetPairs returns the currency pairs that the provider is fetching prices for.
func (p *Provider) GetPairs() []oracletypes.CurrencyPair {
	return p.pairs
}
