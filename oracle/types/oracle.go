package types

import (
	"context"
	"math/big"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/providers/base"
	apihandlers "github.com/skip-mev/connect/v2/providers/base/api/handlers"
	apimetrics "github.com/skip-mev/connect/v2/providers/base/api/metrics"
	wshandlers "github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
	wsmetrics "github.com/skip-mev/connect/v2/providers/base/websocket/metrics"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
)

// ConfigType is the type of the API/WebSocket configuration.
const ConfigType = "price_provider"

type (
	// PriceProvider is a type alias for the base price provider. This specifically
	// implements the provider interface for the price provider along with the
	// additional base provider methods.
	PriceProvider = base.Provider[ProviderTicker, *big.Float]

	// PriceAPIQueryHandlerFactory is a type alias for the price API query handler factory.
	// This is responsible for creating the API query handler for the price provider.
	PriceAPIQueryHandlerFactory = func(
		ctx context.Context,
		logger *zap.Logger,
		cfg config.ProviderConfig,
		apiMetrics apimetrics.APIMetrics,
	) (PriceAPIQueryHandler, error)

	// PriceAPIFetcher is a type alias for the price API fetcher. This is responsible
	// for fetching the prices from the price provider.
	PriceAPIFetcher = apihandlers.APIFetcher[ProviderTicker, *big.Float]

	// PriceAPIDataHandler is a type alias for the price API data handler. This
	// is responsible for parsing http responses and returning the resolved
	// and unresolved prices.
	PriceAPIDataHandler = apihandlers.APIDataHandler[ProviderTicker, *big.Float]

	// PriceAPIQueryHandler is a type alias for the price API query handler. This
	// is responsible for building the API query for the price provider and
	// returning the resolved and unresolved prices.
	PriceAPIQueryHandler = apihandlers.APIQueryHandler[ProviderTicker, *big.Float]

	// PriceWebSocketQueryHandlerFactory is a type alias for the price web socket query handler factory.
	// This is responsible for creating the web socket query handler for the price provider.
	PriceWebSocketQueryHandlerFactory = func(
		ctx context.Context,
		logger *zap.Logger,
		cfg config.ProviderConfig,
		wsMetrics wsmetrics.WebSocketMetrics,
	) (PriceWebSocketQueryHandler, error)

	// PriceWebSocketDataHandler is a type alias for the price web socket data handler.
	// This is responsible for parsing web socket messages and returning the resolved
	// and unresolved prices.
	PriceWebSocketDataHandler = wshandlers.WebSocketDataHandler[ProviderTicker, *big.Float]

	// PriceWebSocketQueryHandler is a type alias for the price web socket query handler.
	// This is responsible for building the web socket query for the price provider and
	// returning the resolved and unresolved prices.
	PriceWebSocketQueryHandler = wshandlers.WebSocketQueryHandler[ProviderTicker, *big.Float]

	// PriceResponse is a type alias for the price response. A price response is
	// composed of a map of resolved prices and a map of unresolved prices. Resolved
	// prices are the prices that were successfully fetched from the API, while
	// unresolved prices are the prices that were not successfully fetched from the API.
	PriceResponse = providertypes.GetResponse[ProviderTicker, *big.Float]

	// ResolvedPrices is a type alias for the resolved prices.
	ResolvedPrices = map[ProviderTicker]providertypes.ResolvedResult[*big.Float]

	// UnResolvedPrices is a type alias for the unresolved prices.
	UnResolvedPrices = map[ProviderTicker]providertypes.UnresolvedResult

	// Prices is a type alias for a map of ticker to a price.
	Prices = map[string]*big.Float
)

var (
	// NewPriceResult is a function alias for the new price result.
	NewPriceResult = providertypes.NewResult[*big.Float]

	// NewPriceResultWithCode is a function alias for the new price result with code.
	NewPriceResultWithCode = providertypes.NewResultWithCode[*big.Float]

	// NewPriceResponse is a function alias for the new price response.
	NewPriceResponse = providertypes.NewGetResponse[ProviderTicker, *big.Float]

	// NewPriceResponseWithErr is a function alias for the new price response with errors.
	NewPriceResponseWithErr = providertypes.NewGetResponseWithErr[ProviderTicker, *big.Float]

	// NewPriceProvider is a function alias for the new price provider.
	NewPriceProvider = base.NewProvider[ProviderTicker, *big.Float]

	// NewPriceAPIQueryHandlerWithFetcher is a function alias for the new API query handler with fetcher.
	NewPriceAPIQueryHandlerWithFetcher = apihandlers.NewAPIQueryHandlerWithFetcher[ProviderTicker, *big.Float]

	// NewPriceWebSocketQueryHandler is a function alias for the new web socket query handler meant to be
	// used by the price providers.
	NewPriceWebSocketQueryHandler = wshandlers.NewWebSocketQueryHandler[ProviderTicker, *big.Float]
)
