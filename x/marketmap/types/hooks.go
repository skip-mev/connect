package types

// MarketMapHooks is the interface that defines the hooks that can be integrated by other modules.
type MarketMapHooks interface {
}

var _ MarketMapHooks = MultiMarketMapHooks{}

// MultiMarketMapHooks defines an array of MarketMapHooks which can be executed in sequence.
type MultiMarketMapHooks []MarketMapHooks
