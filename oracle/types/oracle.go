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
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// ConfigType is the type of the API/WebSocket configuration.
const ConfigType = "price_provider"

type (
	// PriceProviderFactory is a type alias for the price provider factory. This
	// specifically only returns price providers that implement the provider interface
	// and the additional base provider methods.
	PriceProviderFactory = factory.BaseProviderFactory[mmtypes.Ticker, *big.Int]

	// PriceProviderFactory is a type alias for the price provider factory. This
	// specifically only returns price providers that implement the provider interface.
	PriceProviderFactoryI = factory.ProviderFactory[mmtypes.Ticker, *big.Int]

	// PriceProvider is a type alias for the base price provider. This specifically
	// implements the provider interface for the price provider along with the
	// additional base provider methods.
	PriceProvider = base.Provider[mmtypes.Ticker, *big.Int]

	// PriceProviderI is a type alias for the price provider. This specifically
	// implements the provider interface for the price provider.
	PriceProviderI = providertypes.Provider[mmtypes.Ticker, *big.Int]

	// PriceAPIQueryHandlerFactory is a type alias for the price API query handler factory.
	// This is responsible for creating the API query handler for the price provider.
	PriceAPIQueryHandlerFactory = func(
		logger *zap.Logger,
		cfg config.ProviderConfig,
		apiMetrics apimetrics.APIMetrics,
		pMarketMap ProviderMarketMap,
	) (PriceAPIQueryHandler, error)

	// PriceAPIFetcher is a type alias for the price API fetcher. This is responsible
	// for fetching the prices from the price provider.
	PriceAPIFetcher = apihandlers.APIFetcher[mmtypes.Ticker, *big.Int]

	// PriceAPIDataHandler is a type alias for the price API data handler. This
	// is responsible for parsing http responses and returning the resolved
	// and unresolved prices.
	PriceAPIDataHandler = apihandlers.APIDataHandler[mmtypes.Ticker, *big.Int]

	// PriceAPIQueryHandler is a type alias for the price API query handler. This
	// is responsible for building the API query for the price provider and
	// returning the resolved and unresolved prices.
	PriceAPIQueryHandler = apihandlers.APIQueryHandler[mmtypes.Ticker, *big.Int]

	// APIPriceFetcher is a type alias for the API price fetcher. This is responsible
	// for fetching the prices from the price provider using the API.
	APIPriceFetcher = apihandlers.APIPriceFetcher[mmtypes.Ticker, *big.Int]

	// PriceWebSocketQueryHandlerFactory is a type alias for the price web socket query handler factory.
	// This is responsible for creating the web socket query handler for the price provider.
	PriceWebSocketQueryHandlerFactory = func(
		logger *zap.Logger,
		cfg config.ProviderConfig,
		wsMetrics wsmetrics.WebSocketMetrics,
		pMarketMap ProviderMarketMap,
	) (PriceWebSocketQueryHandler, error)

	// PriceWebSocketDataHandler is a type alias for the price web socket data handler.
	// This is responsible for parsing web socket messages and returning the resolved
	// and unresolved prices.
	PriceWebSocketDataHandler = wshandlers.WebSocketDataHandler[mmtypes.Ticker, *big.Int]

	// PriceWebSocketQueryHandler is a type alias for the price web socket query handler.
	// This is responsible for building the web socket query for the price provider and
	// returning the resolved and unresolved prices.
	PriceWebSocketQueryHandler = wshandlers.WebSocketQueryHandler[mmtypes.Ticker, *big.Int]

	// PriceResponse is a type alias for the price response. A price response is
	// composed of a map of resolved prices and a map of unresolved prices. Resolved
	// prices are the prices that were successfully fetched from the API, while
	// unresolved prices are the prices that were not successfully fetched from the API.
	PriceResponse = providertypes.GetResponse[mmtypes.Ticker, *big.Int]

	// ResolvedPrices is a type alias for the resolved prices.
	ResolvedPrices = map[mmtypes.Ticker]providertypes.ResolvedResult[*big.Int]

	// UnResolvedPrices is a type alias for the unresolved prices.
	UnResolvedPrices = map[mmtypes.Ticker]providertypes.UnresolvedResult

	// TickerPrices is a type alias for the map of prices. This is a map of tickers i.e.
	// BTC/USD, ETH/USD, etc. to their respective prices.
	TickerPrices = map[mmtypes.Ticker]*big.Int

	// PriceAggregator is a type alias for the price aggregator. This is responsible for
	// aggregating the resolved prices from the price providers.
	PriceAggregator = aggregator.Aggregator[string, TickerPrices]

	// PriceAggregationFn is a type alias for the price aggregation function. This function
	// is used to aggregate the resolved prices from the price providers.
	PriceAggregationFn = aggregator.AggregateFn[string, TickerPrices]

	// AggregatedProviderPrices is a type alias for the aggregated provider prices. This is
	// a map of provider names to their respective ticker prices.
	AggregatedProviderPrices = aggregator.AggregatedProviderData[string, TickerPrices]
)

var (
	// NewPriceResult is a function alias for the new price result.
	NewPriceResult = providertypes.NewResult[*big.Int]

	// NewPricesResponse is a function alias for the new price response.
	NewPriceResponse = providertypes.NewGetResponse[mmtypes.Ticker, *big.Int]

	// NewPriceResponseWithErr is a function alias for the new price response with errors.
	NewPriceResponseWithErr = providertypes.NewGetResponseWithErr[mmtypes.Ticker, *big.Int]

	// NewPriceProvider is a function alias for the new price provider.
	NewPriceProvider = base.NewProvider[mmtypes.Ticker, *big.Int]

	// NewPriceAPIQueryHandler is a function alias for the new API query handler meant to be
	// used by the price providers.
	NewPriceAPIQueryHandler = apihandlers.NewAPIQueryHandler[mmtypes.Ticker, *big.Int]

	// NewPriceAPIQueryHandlerWithFetcher is a function alias for the new API query handler with fetcher.
	NewPriceAPIQueryHandlerWithFetcher = apihandlers.NewAPIQueryHandlerWithFetcher[mmtypes.Ticker, *big.Int]

	// NewPriceWebSocketQueryHandler is a function alias for the new web socket query handler meant to be
	// used by the price providers.
	NewPriceWebSocketQueryHandler = wshandlers.NewWebSocketQueryHandler[mmtypes.Ticker, *big.Int]

	// NewPriceAggregator is a function alias for the new price aggregator.
	NewPriceAggregator = aggregator.NewDataAggregator[string, TickerPrices]
)
