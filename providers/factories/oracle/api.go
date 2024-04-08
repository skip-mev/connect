package oracle

import (
	"fmt"
	"net/http"

	"github.com/skip-mev/slinky/providers/apis/defi/raydium"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/providers/apis/binance"
	coinbaseapi "github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/apis/coingecko"
	"github.com/skip-mev/slinky/providers/apis/geckoterminal"
	"github.com/skip-mev/slinky/providers/apis/kraken"
	"github.com/skip-mev/slinky/providers/apis/uniswapv3"
	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
	"github.com/skip-mev/slinky/providers/static"
	"github.com/skip-mev/slinky/providers/volatile"
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
		apiPriceFetcher types.PriceAPIFetcher
		apiDataHandler  types.PriceAPIDataHandler
	)

	requestHandler, err := apihandlers.NewRequestHandlerImpl(client)
	if err != nil {
		return nil, err
	}

	switch cfg.Name {
	case binance.Name:
		apiDataHandler, err = binance.NewAPIHandler(marketMap, cfg.API)
	case coinbaseapi.Name:
		apiDataHandler, err = coinbaseapi.NewAPIHandler(marketMap, cfg.API)
	case coingecko.Name:
		apiDataHandler, err = coingecko.NewAPIHandler(marketMap, cfg.API)
	case geckoterminal.Name:
		apiDataHandler, err = geckoterminal.NewAPIHandler(marketMap, cfg.API)
	case kraken.Name:
		apiDataHandler, err = kraken.NewAPIHandler(marketMap, cfg.API)
	case uniswapv3.Name:
		var ethClient uniswapv3.EVMClient
		ethClient, err = uniswapv3.NewGoEthereumClientImpl(cfg.API.URL)
		if err != nil {
			return nil, err
		}

		apiPriceFetcher, err = uniswapv3.NewPriceFetcher(logger, metrics, cfg.API, ethClient)
	case static.Name:
		apiDataHandler, err = static.NewAPIHandler(marketMap)
		if err != nil {
			return nil, err
		}

		requestHandler = static.NewStaticMockClient()
	case volatile.Name:
		apiDataHandler, err = volatile.NewAPIHandler(marketMap)
		if err != nil {
			return nil, err
		}

		requestHandler = static.NewStaticMockClient()
	case raydium.Name:
		apiPriceFetcher, err = raydium.NewAPIPriceFetcher(
			marketMap,
			cfg.API,
			logger,
		)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown provider: %s", cfg.Name)
	}
	if err != nil {
		return nil, err
	}

	// if no apiPriceFetcher has been created yet, create a default REST API price fetcher.
	if apiPriceFetcher == nil {
		apiPriceFetcher, err = apihandlers.NewRestAPIFetcher(
			requestHandler,
			apiDataHandler,
			metrics,
			cfg.API,
			logger,
		)
		if err != nil {
			return nil, err
		}
	}

	// Create the API query handler which encapsulates all of the fetching and parsing logic.
	return types.NewPriceAPIQueryHandlerWithFetcher(
		logger,
		cfg.API,
		apiPriceFetcher,
		metrics,
	)
}
