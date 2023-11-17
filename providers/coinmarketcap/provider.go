package coinmarketcap

import (
	"context"
	"fmt"
	"sync"

	"cosmossdk.io/log"

	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// Name is the name of this provider
	Name = "coinmarketcap"
)

var _ oracle.Provider = (*Provider)(nil)

// Provider is the implementation of the oracle's Provider interface for coinmarketcap.
type Provider struct {
	logger log.Logger

	// pairs is a list of currency pairs that the provider should fetch
	// prices for.
	pairs []oracletypes.CurrencyPair

	// config is the coinmarketcap config.
	config Config
}

// NewProvider returns a new coinmarketcap provider. It uses the provided API-key in the header of outgoing requests to coinmarketcap's API.
func NewProvider(logger log.Logger, pairs []oracletypes.CurrencyPair, providerConfig config.ProviderConfig) (*Provider, error) {
	if providerConfig.Name != Name {
		return nil, fmt.Errorf("expected provider config name %s, got %s", Name, providerConfig.Name)
	}

	config, err := ReadCoinMarketCapConfigFromFile(providerConfig.Path)
	if err != nil {
		return nil, err
	}

	logger = logger.With("provider", Name)
	logger.Info("creating new coinmarketcap provider", "pairs", pairs, "config", config)

	return &Provider{
		pairs:  pairs,
		logger: logger,
		config: config,
	}, nil
}

// GetPrices returns the current set of prices for each of the currency pairs. This method starts all
// price requests concurrently, and waits for them all to finish, or for the context to be cancelled,
// at which point it aggregates the responses and returns.
func (p *Provider) GetPrices(ctx context.Context) (map[oracletypes.CurrencyPair]aggregator.QuotePrice, error) {
	type priceData struct {
		aggregator.QuotePrice
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
	prices := make(map[oracletypes.CurrencyPair]aggregator.QuotePrice)
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
