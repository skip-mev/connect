package oracle

import (
	"context"
	"time"

	"github.com/skip-mev/connect/v2/oracle/types"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

// Oracle defines the expected interface for an oracle. It is consumed by the oracle server.
//
//go:generate mockery --name Oracle --filename mock_oracle.go
type Oracle interface {
	IsRunning() bool
	GetLastSyncTime() time.Time
	GetPrices() types.Prices
	GetMarketMap() mmtypes.MarketMap
	Start(ctx context.Context) error
	Stop()
}

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

// generalProvider is an interface for a provider that implements the base provider.
type generalProvider interface {
	// Start starts the provider.
	Start(ctx context.Context) error
	// Name is the provider's name.
	Name() string
}
