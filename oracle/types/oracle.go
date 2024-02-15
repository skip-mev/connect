package types

import (
	"math/big"

	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/types/factory"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

type (
	// PriceProviderFactory is a type alias for the price provider factory.
	PriceProviderFactory = factory.ProviderFactory[mmtypes.Ticker, *big.Int]

	// PriceProvider is a type alias for the price provider.
	PriceProvider = providertypes.Provider[mmtypes.Ticker, *big.Int]

	// TickerPrices is a type alias for the map of prices. This is a map of tickers i.e.
	// BTC/USD, ETH/USD, etc. to their respective prices.
	TickerPrices = map[mmtypes.Ticker]*big.Int
)
