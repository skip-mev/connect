package types

import (
	"math/big"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base"
	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
	apimetrics "github.com/skip-mev/slinky/providers/base/api/metrics"
	wshandlers "github.com/skip-mev/slinky/providers/base/websocket/handlers"
	wsmetrics "github.com/skip-mev/slinky/providers/base/websocket/metrics"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/types/factory"
)

// ConfigType is the type of the API/WebSocket configuration.
const ConfigType = "price_provider"

type (
	// PriceProviderFactory is a type alias for the price provider factory. This
	// specifically only returns price providers that implement the provider interface
	// and the additional base provider methods.
	PriceProviderFactory = factory.BaseProviderFactory[ProviderTicker, *big.Float]

	// PriceProviderFactory is a type alias for the price provider factory. This
	// specifically only returns price providers that implement the provider interface.
	PriceProviderFactoryI = factory.ProviderFactory[ProviderTicker, *big.Float]

	// PriceProvider is a type alias for the base price provider. This specifically
	// implements the provider interface for the price provider along with the
	// additional base provider methods.
	PriceProvider = base.Provider[ProviderTicker, *big.Float]

	// PriceProviderI is a type alias for the price provider. This specifically
	// implements the provider interface for the price provider.
	PriceProviderI = providertypes.Provider[ProviderTicker, *big.Float]

	// PriceAPIQueryHandlerFactory is a type alias for the price API query handler factory.
	// This is responsible for creating the API query handler for the price provider.
	PriceAPIQueryHandlerFactory = func(
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

	// TickerPrices is a type alias for the map of prices. This is a map of tickers i.e.
	// BTC/USD, ETH/USD, etc. to their respective prices.
	TickerPrices = map[ProviderTicker]*big.Float

	// AggregatorPrices is a type alias for a map of off-chain ticker to the price.
	AggregatorPrices = map[string]*big.Float

	// PriceAggregator is a type alias for the price aggregator. This is responsible for
	// aggregating the resolved prices from the price providers.
	PriceAggregator = aggregator.Aggregator[string, AggregatorPrices]

	// PriceAggregationFn is a type alias for the price aggregation function. This function
	// is used to aggregate the resolved prices from the price providers.
	PriceAggregationFn = aggregator.AggregateFn[string, AggregatorPrices]

	// AggregatedProviderPrices is a type alias for the aggregated provider prices. This is
	// a map of provider names to their respective ticker prices.
	AggregatedProviderPrices = aggregator.AggregatedProviderData[string, AggregatorPrices]
)

var (
	// NewPriceResult is a function alias for the new price result.
	NewPriceResult = providertypes.NewResult[*big.Float]

	// NewPricesResponse is a function alias for the new price response.
	NewPriceResponse = providertypes.NewGetResponse[ProviderTicker, *big.Float]

	// NewPriceResponseWithErr is a function alias for the new price response with errors.
	NewPriceResponseWithErr = providertypes.NewGetResponseWithErr[ProviderTicker, *big.Float]

	// NewPriceProvider is a function alias for the new price provider.
	NewPriceProvider = base.NewProvider[ProviderTicker, *big.Float]

	// NewPriceAPIQueryHandler is a function alias for the new API query handler meant to be
	// used by the price providers.
	NewPriceAPIQueryHandler = apihandlers.NewAPIQueryHandler[ProviderTicker, *big.Float]

	// NewPriceAPIQueryHandlerWithFetcher is a function alias for the new API query handler with fetcher.
	NewPriceAPIQueryHandlerWithFetcher = apihandlers.NewAPIQueryHandlerWithFetcher[ProviderTicker, *big.Float]

	// NewPriceWebSocketQueryHandler is a function alias for the new web socket query handler meant to be
	// used by the price providers.
	NewPriceWebSocketQueryHandler = wshandlers.NewWebSocketQueryHandler[ProviderTicker, *big.Float]

	// NewPriceAggregator is a function alias for the new price aggregator.
	NewPriceAggregator = aggregator.NewDataAggregator[string, AggregatorPrices]
)
