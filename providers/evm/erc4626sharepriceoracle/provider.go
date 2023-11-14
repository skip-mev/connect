package erc4626sharepriceoracle

import (
	"context"
	"fmt"
	"sync"

	"cosmossdk.io/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/evm"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// Name is the name of this provider
	Name = "erc4626-share-price-oracle"
)

var _ oracle.Provider = (*Provider)(nil)

type (
	// Provider is the implementation of the oracle's Provider interface for instances of the
	// ERC4626SharePriceOracle.sol contract.
	Provider struct {
		logger log.Logger

		// pairs is a list of currency pairs that the provider should fetch
		// prices for.
		pairs []oracletypes.CurrencyPair

		// config is the ERC4626SharePriceOracle config.
		config evm.Config

		// rpcEndpoint is the endpoint of the ethereum rpc node to use for querying. This
		// is used to make RPC calls to the Ethereum node with a configured API key.
		rpcEndpoint string
	}
)

// NewProvider returns a new ERC4626SharePriceOracle provider. It uses the provided API-key to
// make RPC calls to Alchemy. Note that only the Quote denom is used; the Quote/Base pair is
// naturally determined by the contract address, so be sure the configured addresses are
// correct.
func NewProvider(logger log.Logger, pairs []oracletypes.CurrencyPair, providerConfig config.ProviderConfig) (*Provider, error) {
	if providerConfig.Name != Name {
		return nil, fmt.Errorf("expected provider config name to be %s, got %s", Name, providerConfig.Name)
	}

	config, err := evm.ReadEVMConfigFromFile(providerConfig.Path)
	if err != nil {
		return nil, err
	}

	provider := &Provider{}
	for _, pair := range pairs {
		if metadata, ok := config.TokenNameToMetadata[pair.Quote]; ok {
			if !common.IsHexAddress(metadata.Symbol) {
				return nil, fmt.Errorf("invalid contract address: %s", metadata.Symbol)
			}

			provider.pairs = append(provider.pairs, pair)
		}
	}

	logger = logger.With("provider", Name)
	logger.Info("creating new erc4626-share-price-oracle provider", "pairs", pairs, "config", config)

	provider.logger = logger
	provider.rpcEndpoint = getRPCEndpoint(config)
	provider.config = config

	return provider, nil
}

// Name returns the name of this provider.
func (p *Provider) Name() string {
	return Name
}

// GetPrices returns the prices of the given pairs.
func (p *Provider) GetPrices(ctx context.Context) (map[oracletypes.CurrencyPair]aggregator.QuotePrice, error) {
	type priceData struct {
		aggregator.QuotePrice
		oracletypes.CurrencyPair
	}

	// create response channel
	resp := make(chan priceData, len(p.pairs))

	wg := sync.WaitGroup{}
	wg.Add(len(p.pairs))

	// fan-out requests to RPC provider
	for _, currencyPair := range p.pairs {
		go func(pair oracletypes.CurrencyPair) {
			defer wg.Done()

			// get price
			qp, err := p.getPriceForPair(pair)
			if err != nil {
				p.logger.Error("failed to get price for pair", "provider", p.Name(), "pair", pair, "err", err)
			} else {
				p.logger.Info("fetched price for pair", "pair", pair, "provider", p.Name())

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

// GetPairs returns the pairs this provider can provide prices for.
func (p *Provider) GetPairs() []oracletypes.CurrencyPair {
	return p.pairs
}

// SetPairs sets the pairs this provider can provide prices for. This method will map new pairs
// to an empty string in the contract address mapping. Be sure that pairs added have
// corresponding contract addresses in their config metadata.
func (p *Provider) SetPairs(pairs ...oracletypes.CurrencyPair) {
	p.pairs = pairs
}

// getPairContractAddress gets the contract address for the pair.
func (p *Provider) getPairContractAddress(pair oracletypes.CurrencyPair) (string, bool) {
	metadata, found := p.config.TokenNameToMetadata[pair.Quote]
	if found {
		return metadata.Symbol, found
	}

	return "", found
}

// finish takes a wait-group, and returns a channel that is sent on when the
// Waitgroup is finished.
func finish(wg *sync.WaitGroup) <-chan struct{} {
	ch := make(chan struct{})

	// non-blocking wait for waitgroup to finish, and return channel
	go func() {
		wg.Wait()
		close(ch)
	}()
	return ch
}
