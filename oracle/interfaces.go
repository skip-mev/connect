package oracle

import "github.com/skip-mev/slinky/oracle/types"

// PriceAggregator is an interface for aggregating prices from multiple providers.
//
//go:generate mockery --name PriceAggregator
type PriceAggregator interface {
	SetProviderPrices(provider string, prices types.Prices)
	AggregatePrices()
	GetPrices() types.Prices
	Reset()
}
