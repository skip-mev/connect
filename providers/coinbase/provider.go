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

type (
	Provider struct {
		log   log.Logger
		pairs []types.CurrencyPair
	}
)

func NewProvider(logger log.Logger) *Provider {
	return &Provider{
		log: logger,
	}
}

func (p *Provider) Name() string {
	return Name
}

func (p *Provider) GetPrices() (map[string]types.TickerPrice, error) {
	resp := make(map[string]types.TickerPrice)

	for _, currencyPair := range p.pairs {
		spotPrice, err := getPriceForPair(currencyPair)
		if err != nil {
			return nil, err
		}

		if spotPrice != nil {
			resp[currencyPair.String()] = *spotPrice
		}
	}

	return resp, nil
}

func (p *Provider) SetPairs(pairs ...types.CurrencyPair) {
	p.pairs = pairs
}

func (p *Provider) GetPairs() []types.CurrencyPair {
	return p.pairs
}
