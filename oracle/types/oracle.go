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

	// PriceAPIDataHandler is a type alias for the price API data handler. This
	// is responsible for parsing http responses and returning the resolved
	// and unresolved prices.
	PriceAPIDataHandler = apihandlers.APIDataHandler[mmtypes.Ticker, *big.Int]

	// PriceResponse is a type alias for the price response. A price response is
	// composed of a map of resolved prices and a map of unresolved prices. Resolved
	// prices are the prices that were successfully fetched from the API, while
	// unresolved prices are the prices that were not successfully fetched from the API.
	PriceResponse = providertypes.GetResponse[mmtypes.Ticker, *big.Int]

	// ResolvedPrices is a type alias for the resolved prices.
	ResolvedPrices = map[mmtypes.Ticker]providertypes.Result[*big.Int]

	// UnResolvedPrices is a type alias for the unresolved prices.
	UnResolvedPrices = map[mmtypes.Ticker]error

	// TickerPrices is a type alias for the map of prices. This is a map of tickers i.e.
	// BTC/USD, ETH/USD, etc. to their respective prices.
	TickerPrices = map[mmtypes.Ticker]*big.Int
)
