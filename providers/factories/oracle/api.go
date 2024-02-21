package oracle

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/providers/apis/binance"
	coinbaseapi "github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/apis/coingecko"
	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
	"github.com/skip-mev/slinky/providers/static"
)

// APIQueryHandlerFactory returns a sample implementation of the API query handler factory.
// Specifically, this factory function returns API query handlers that are used to fetch data from
// the price providers.
func APIQueryHandlerFactory(
	logger *zap.Logger,
	cfg config.ProviderConfig,
	metrics metrics.APIMetrics,
	marketMap types.ProviderMarketMap,
) (types.PriceAPIQueryHandler, error) {
	// Validate the provider config.
	err := cfg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	// Create the underlying client that will be used to fetch data from the API. This client
	// will limit the number of concurrent connections and uses the configured timeout to
	// ensure requests do not hang.
	tickers := marketMap.GetTickers()
	maxCons := math.Min(len(tickers), cfg.API.MaxQueries)
	client := &http.Client{
		Transport: &http.Transport{MaxConnsPerHost: maxCons},
		Timeout:   cfg.API.Timeout,
	}

	var (
		apiDataHandler types.PriceAPIDataHandler
		requestHandler apihandlers.RequestHandler
	)

	switch cfg.Name {
	case binance.Name:
		apiDataHandler, err = binance.NewAPIHandler(marketMap, cfg.API)
	case coinbaseapi.Name:
		apiDataHandler, err = coinbaseapi.NewAPIHandler(marketMap, cfg.API)
	case coingecko.Name:
		apiDataHandler, err = coingecko.NewAPIHandler(marketMap, cfg.API)
	case static.Name:
		apiDataHandler, err = static.NewAPIHandler(marketMap)
		if err != nil {
			return nil, err
		}

		requestHandler = static.NewStaticMockClient()
	default:
		return nil, fmt.Errorf("unknown provider: %s", cfg.Name)
	}
	if err != nil {
		return nil, err
	}

	// If a custom request handler is not provided, create a new default one.
	if requestHandler == nil {
		requestHandler, err = apihandlers.NewRequestHandlerImpl(client)
		if err != nil {
			return nil, err
		}
	}

	// Create the API query handler which encapsulates all of the fetching and parsing logic.
	return types.NewPriceAPIQueryHandler(
		logger,
		cfg.API,
		requestHandler,
		apiDataHandler,
		metrics,
	)
}
