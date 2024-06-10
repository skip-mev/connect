package oracle

import (
	"github.com/skip-mev/slinky/oracle/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// PriceAggregator is an interface for aggregating prices from multiple providers. Implementations of PriceAggregator
// should be made safe for concurrent use.
//
//go:generate mockery --name PriceAggregator
type PriceAggregator interface {
	SetProviderPrices(provider string, prices types.Prices)
	UpdateMarketMap(mmtypes.MarketMap)
	AggregatePrices()
	GetPrices() types.Prices
	Reset()
}
