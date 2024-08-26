package types

import (
	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/providers/base"
	apihandlers "github.com/skip-mev/connect/v2/providers/base/api/handlers"
	apimetrics "github.com/skip-mev/connect/v2/providers/base/api/metrics"
	providermetrics "github.com/skip-mev/connect/v2/providers/base/metrics"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

// ConfigType is the type of the API/WebSocket configuration.
const ConfigType = "market_map_provider"

// Chain is a type alias for the market map key. This is used to uniquely
// identify a market map.
type Chain struct {
	// ChainID is the chain ID of the market map.
	ChainID string
}

// String returns the string representation of the Chain schema.
func (mms Chain) String() string {
	return mms.ChainID
}

type (
	// MarketMapProvider is a type alias for the market map provider.
	MarketMapProvider = base.Provider[Chain, *mmtypes.MarketMapResponse]

	// MarketMapFactory is a type alias for the market map factory.
	MarketMapFactory = func(
		logger *zap.Logger,
		providerMetrics providermetrics.ProviderMetrics,
		apiMetrics apimetrics.APIMetrics,
		cfg config.ProviderConfig,
	) (*MarketMapProvider, error)

	// MarketMapAPIDataHandler is a type alias for the market map API data handler. This
	// is responsible for parsing http responses and returning the resolved and unresolved
	// market map data.
	MarketMapAPIDataHandler = apihandlers.APIDataHandler[Chain, *mmtypes.MarketMapResponse]

	// MarketMapFetcher is a type alias for the market map fetcher. This is responsible for
	// fetching the market map data.
	MarketMapFetcher = apihandlers.APIFetcher[Chain, *mmtypes.MarketMapResponse]

	// MarketMapResponse is a type alias for the market map response. This is used to
	// represent the resolved and unresolved market map data.
	MarketMapResponse = providertypes.GetResponse[Chain, *mmtypes.MarketMapResponse]

	// MarketMapResult is a type alias for the market map result.
	MarketMapResult = providertypes.ResolvedResult[*mmtypes.MarketMapResponse]

	// ResolvedMarketMap is a type alias for the resolved market map.
	ResolvedMarketMap = map[Chain]MarketMapResult

	// UnResolvedMarketMap is a type alias for the unresolved market map.
	UnResolvedMarketMap = map[Chain]providertypes.UnresolvedResult
)

var (
	// NewMarketMapResult is a function alias for the new market map result.
	NewMarketMapResult = providertypes.NewResult[*mmtypes.MarketMapResponse]

	// NewMarketMapResponse is a function alias for the new market map response.
	NewMarketMapResponse = providertypes.NewGetResponse[Chain, *mmtypes.MarketMapResponse]

	// NewMarketMapResponseWithErr returns a new market map response with an error.
	NewMarketMapResponseWithErr = providertypes.NewGetResponseWithErr[Chain, *mmtypes.MarketMapResponse]

	// NewMarketMapProvider is a function alias for the new market map provider.
	NewMarketMapProvider = base.NewProvider[Chain, *mmtypes.MarketMapResponse]

	// NewMarketMapAPIQueryHandler is a function alias for the new market map API query handler.
	NewMarketMapAPIQueryHandler = apihandlers.NewAPIQueryHandler[Chain, *mmtypes.MarketMapResponse]

	// NewMarketMapAPIQueryHandlerWithMarketMapFetcher is a function alias for the new market map API query handler with market map fetcher.
	NewMarketMapAPIQueryHandlerWithMarketMapFetcher = apihandlers.NewAPIQueryHandlerWithFetcher[Chain, *mmtypes.MarketMapResponse]
)
