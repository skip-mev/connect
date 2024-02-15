package oracle

import (
	"fmt"
	"math/big"
	"net/http"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/providers/apis/binance"
	coinbaseapi "github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/apis/coingecko"
	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
	"github.com/skip-mev/slinky/providers/static"
	"github.com/skip-mev/slinky/providers/types/factory"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// APIQueryHandlerFactory returns a sample implementation of the API query handler factory.
// Specifically, this factory function returns API query handlers that are used to fetch data from
// the price providers.
func APIQueryHandlerFactory(aggConfig mmtypes.AggregateMarketConfig) factory.APIQueryHandlerFactory[mmtypes.Ticker, *big.Int] {
	return func(logger *zap.Logger, cfg config.ProviderConfig, metrics metrics.APIMetrics) (apihandlers.APIQueryHandler[mmtypes.Ticker, *big.Int], error) {
		// If the API is not enabled, return an error.
		if !cfg.API.Enabled {
			return nil, fmt.Errorf("API for provider %s is not enabled", cfg.Name)
		}

		// Validate the provider config.
		if err := cfg.ValidateBasic(); err != nil {
			return nil, err
		}

		// Ensure the market config is valid.
		if err := aggConfig.ValidateBasic(); err != nil {
			return nil, err
		}

		// Ensure that the market configuration is supported by the provider.
		market, ok := aggConfig.MarketConfigs[cfg.Name]
		if !ok {
			return nil, fmt.Errorf("provider %s is not supported by the market config", cfg.Name)
		}

		// Create the underlying client that will be used to fetch data from the API. This client
		// will limit the number of concurrent connections and uses the configured timeout to
		// ensure requests do not hang.
		maxCons := math.Min(len(market.TickerConfigs), cfg.API.MaxQueries)
		client := &http.Client{
			Transport: &http.Transport{MaxConnsPerHost: maxCons},
			Timeout:   cfg.API.Timeout,
		}

		var (
			apiDataHandler apihandlers.APIDataHandler[mmtypes.Ticker, *big.Int]
			requestHandler apihandlers.RequestHandler
			err            error
		)

		switch cfg.Name {
		case binance.Name:
			apiDataHandler, err = binance.NewAPIHandler(market, cfg.API)
		case coinbaseapi.Name:
			apiDataHandler, err = coinbaseapi.NewAPIHandler(market, cfg.API)
		case coingecko.Name:
			apiDataHandler, err = coingecko.NewAPIHandler(market, cfg.API)
		case static.Name:
			apiDataHandler, err = static.NewAPIHandler(market)
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
		return apihandlers.NewAPIQueryHandler[mmtypes.Ticker, *big.Int](
			logger,
			cfg.API,
			requestHandler,
			apiDataHandler,
			metrics,
		)
	}
}
