package coingecko

import (
	"context"
	"fmt"
	"strings"

	"cosmossdk.io/log"
	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// Name is the name of the provider.
	Name = "coingecko"
)

var _ oracle.Provider = (*Provider)(nil)

// Provider implements the Provider interface for CoinGecko. This provider
// is a very simple implementation that fetches prices from the CoinGecko API.
type Provider struct {
	logger log.Logger

	// pairs is a list of currency pairs that the provider should fetch
	// prices for.
	pairs []oracletypes.CurrencyPair

	// bases is a list of base currencies that the provider should fetch
	// prices for.
	bases string

	// quotes is a list of quote currencies that the provider should fetch
	// prices for.
	quotes string

	// config is the CoinGecko config.
	config Config
}

// NewProvider returns a new CoinGecko provider.
func NewProvider(logger log.Logger, pairs []oracletypes.CurrencyPair, providerConfig config.ProviderConfig) (*Provider, error) {
	if providerConfig.Name != Name {
		return nil, fmt.Errorf("expected provider config name %s, got %s", Name, providerConfig.Name)
	}

	config, err := ReadCoinGeckoConfigFromFile(providerConfig.Path)
	if err != nil {
		return nil, err
	}

	bases, quotes := getUniqueBaseAndQuoteDenoms(pairs)

	logger = logger.With("provider", Name)
	logger.Info("creating new coingecko provider", "pairs", pairs, "config", config)

	return &Provider{
		pairs:  pairs,
		logger: logger,
		bases:  strings.Join(bases, ","),
		quotes: strings.Join(quotes, ","),
		config: config,
	}, nil
}

// Name returns the name of the provider.
func (p *Provider) Name() string {
	return Name
}

// GetPrices returns the current set of prices for each of the currency pairs.
func (p *Provider) GetPrices(ctx context.Context) (map[oracletypes.CurrencyPair]aggregator.QuotePrice, error) {
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
