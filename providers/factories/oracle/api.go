package oracle

import (
	"context"
	"fmt"
	"github.com/skip-mev/slinky/providers/apis/defi/raydium"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/apis/binance"
	coinbaseapi "github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/apis/coingecko"
	"github.com/skip-mev/slinky/providers/apis/defi/uniswapv3"
	"github.com/skip-mev/slinky/providers/apis/geckoterminal"
	"github.com/skip-mev/slinky/providers/apis/kraken"
	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
	"github.com/skip-mev/slinky/providers/static"
	"github.com/skip-mev/slinky/providers/volatile"
)

// APIQueryHandlerFactory returns a sample implementation of the API query handler factory.
// Specifically, this factory function returns API query handlers that are used to fetch data from
// the price providers.
func APIQueryHandlerFactory(
	ctx context.Context,
	logger *zap.Logger,
	cfg config.ProviderConfig,
	metrics metrics.APIMetrics,
) (types.PriceAPIQueryHandler, error) {
	// Validate the provider config.
	err := cfg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	// Create the underlying client that will be used to fetch data from the API. This client
	// will limit the number of concurrent connections and uses the configured timeout to
	// ensure requests do not hang.
	client := &http.Client{
		Transport: &http.Transport{
			MaxConnsPerHost: cfg.API.MaxQueries,
			Proxy:           http.ProxyFromEnvironment,
		},
		Timeout: cfg.API.Timeout,
	}

	var (
		apiPriceFetcher types.PriceAPIFetcher
		apiDataHandler  types.PriceAPIDataHandler
	)

	requestHandler, err := apihandlers.NewRequestHandlerImpl(client)
	if err != nil {
		return nil, err
	}

	switch providerName := cfg.Name; {
	case providerName == binance.Name:
		apiDataHandler, err = binance.NewAPIHandler(cfg.API)
	case providerName == coinbaseapi.Name:
		apiDataHandler, err = coinbaseapi.NewAPIHandler(cfg.API)
	case providerName == coingecko.Name:
		apiDataHandler, err = coingecko.NewAPIHandler(cfg.API)
	case providerName == geckoterminal.Name:
		apiDataHandler, err = geckoterminal.NewAPIHandler(cfg.API)
	case providerName == kraken.Name:
		apiDataHandler, err = kraken.NewAPIHandler(cfg.API)
	case strings.HasPrefix(providerName, uniswapv3.BaseName):
		apiPriceFetcher, err = uniswapv3.NewPriceFetcher(ctx, logger, cfg.API)
	case providerName == static.Name:
		apiDataHandler = static.NewAPIHandler()
		requestHandler = static.NewStaticMockClient()
	case providerName == volatile.Name:
		apiDataHandler = volatile.NewAPIHandler()
		requestHandler = static.NewStaticMockClient()
	case providerName == raydium.Name:
		apiPriceFetcher, err = raydium.NewAPIPriceFetcher(
			cfg.API,
			logger,
		)
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
