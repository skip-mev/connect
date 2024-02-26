package types

import (
	"fmt"

	"github.com/skip-mev/slinky/providers/base"
	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// Chain is a type alias for the market map key. This is used to uniquely
// identify a market map.
type Chain struct {
	// ChainID is the chain ID of the market map.
	ChainID string
}

// String returns the string representation of the Chain schema.
func (mms Chain) String() string {
	return fmt.Sprintf("ChainID: %s", mms.ChainID)
}

type (
	// MarketMapProvider is a type alias for the market map provider.
	MarketMapProvider = providertypes.Provider[Chain, *mmtypes.GetMarketMapResponse]

	// MarketMapAPIQueryHandler is a type alias for the market map API query handler. This
	// is responsible for querying the market map API and returning the resolved and unresolved
	// market map data.
	MarketMapAPIQueryHandler = apihandlers.APIQueryHandler[Chain, *mmtypes.GetMarketMapResponse]

	// MarketMapAPIDataHandler is a type alias for the market map API data handler. This
	// is responsible for parsing http responses and returning the resolved and unresolved
	// market map data.
	MarketMapAPIDataHandler = apihandlers.APIDataHandler[Chain, *mmtypes.GetMarketMapResponse]

	// MarketMapResponse is a type alias for the market map response. This is used to
	// represent the resolved and unresolved market map data.
	MarketMapResponse = providertypes.GetResponse[Chain, *mmtypes.GetMarketMapResponse]

	// MarketMapResult is a type alias for the market map result.
	MarketMapResult = providertypes.Result[*mmtypes.GetMarketMapResponse]

	// ResolvedMarketMap is a type alias for the resolved market map.
	ResolvedMarketMap = map[Chain]MarketMapResult

	// UnResolvedMarketMap is a type alias for the unresolved market map.
	UnResolvedMarketMap = map[Chain]error
)

var (
	// NewMarketMapResult is a function alias for the new market map result.
	NewMarketMapResult = providertypes.NewResult[*mmtypes.GetMarketMapResponse]

	// NewMarketMapResponse is a function alias for the new market map response.
	NewMarketMapResponse = providertypes.NewGetResponse[Chain, *mmtypes.GetMarketMapResponse]

	// NewMarketMapResponseWithErr returns a new market map response with an error.
	NewMarketMapResponseWithErr = providertypes.NewGetResponseWithErr[Chain, *mmtypes.GetMarketMapResponse]

	// NewMarketMapProvider is a function alias for the new market map provider.
	NewMarketMapProvider = base.NewProvider[Chain, *mmtypes.GetMarketMapResponse]

	// NewMarketMapAPIQueryHandler is a function alias for the new market map API query handler.
	NewMarketMapAPIQueryHandler = apihandlers.NewAPIQueryHandler[Chain, *mmtypes.GetMarketMapResponse]
)
