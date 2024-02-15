package types

import (
	"math/big"

	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/types/factory"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

type (
	// PriceProviderFactory is a type alias for the price provider factory.
	PriceProviderFactory = factory.ProviderFactory[mmtypes.Ticker, *big.Int]

	// PriceProvider is a type alias for the price provider.
	PriceProvider = providertypes.Provider[mmtypes.Ticker, *big.Int]

	// PriceAPIDataHandler is a type alias for the price API data handler.
	PriceAPIDataHandler = apihandlers.APIDataHandler[mmtypes.Ticker, *big.Int]

	// PriceResponse is a type alias for the price response.
	PriceResponse = providertypes.GetResponse[mmtypes.Ticker, *big.Int]

	// ResolvedPrices is a type alias for the resolved prices.
	ResolvedPrices = map[mmtypes.Ticker]providertypes.Result[*big.Int]

	// UnResolvedPrices is a type alias for the unresolved prices.
	UnResolvedPrices = map[mmtypes.Ticker]error

	// TickerPrices is a type alias for the map of prices. This is a map of tickers i.e.
	// BTC/USD, ETH/USD, etc. to their respective prices.
	TickerPrices = map[mmtypes.Ticker]*big.Int
)
