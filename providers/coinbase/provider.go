package coinbase

import (
	"context"
	"fmt"

	"cosmossdk.io/log"
	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// Name is the name of the provider.
	Name = "coinbase"
)

var _ oracle.Provider = (*Provider)(nil)

// Provider implements the Provider interface for Coinbase. This provider
// is a very simple implementation that fetches spot prices from the Coinbase API.
type Provider struct {
	logger log.Logger

	// pairs is a list of currency pairs that the provider should fetch
	// prices for.
	pairs []oracletypes.CurrencyPair

	// config is the Coinbase config.
	config Config
}

// NewProvider returns a new Coinbase provider.
//
// THIS PROVIDER SHOULD NOT BE USED IN PRODUCTION. IT IS ONLY MEANT FOR TESTING.
func NewProvider(logger log.Logger, pairs []oracletypes.CurrencyPair, providerConfig config.ProviderConfig) (*Provider, error) {
	if providerConfig.Name != Name {
		return nil, fmt.Errorf("expected provider config name %s, got %s", Name, providerConfig.Name)
	}

	config, err := ReadCoinbaseConfigFromFile(providerConfig.Path)
	if err != nil {
		return nil, err
	}

	logger = logger.With("provider", Name)
	logger.Info("creating new coinbase provider", "pairs", pairs, "config", config)

	return &Provider{
		pairs:  pairs,
		logger: logger,
		config: config,
	}, nil
}

// Name returns the name of the provider.
func (p *Provider) Name() string {
	return Name
}

// GetPrices returns the current set of prices for each of the currency pairs.
func (p *Provider) GetPrices(ctx context.Context) (map[oracletypes.CurrencyPair]aggregator.QuotePrice, error) {
	resp := make(map[oracletypes.CurrencyPair]aggregator.QuotePrice)

	for _, currencyPair := range p.pairs {
		p.logger.Info("fetching price for pair", currencyPair.ToString())

		spotPrice, err := p.getPriceForPair(ctx, currencyPair)
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
func (p *Provider) SetPairs(pairs ...oracletypes.CurrencyPair) {
	p.pairs = pairs
}

// GetPairs returns the currency pairs that the provider is fetching prices for.
func (p *Provider) GetPairs() []oracletypes.CurrencyPair {
	return p.pairs
}
