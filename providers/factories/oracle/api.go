package oracle

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/providers/apis/binance"
	coinbaseapi "github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/apis/coingecko"
	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
	"github.com/skip-mev/slinky/providers/static"
	"github.com/skip-mev/slinky/providers/types/factory"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"go.uber.org/zap"
)

// DefaultAPIQueryHandlerFactory returns a sample implementation of the API query handler factory.
func DefaultAPIQueryHandlerFactory() factory.APIQueryHandlerFactory[oracletypes.CurrencyPair, *big.Int] {
	return func(logger *zap.Logger, cfg config.ProviderConfig, metrics metrics.APIMetrics) (apihandlers.APIQueryHandler[oracletypes.CurrencyPair, *big.Int], error) {
		// Validate the provider config.
		err := cfg.ValidateBasic()
		if err != nil {
			return nil, err
		}

		// Create the underlying client that will be used to fetch data from the API. This client
		// will limit the number of concurrent connections and uses the configured timeout to
		// ensure requests do not hang.
		cps := cfg.Market.GetCurrencyPairs()
		maxCons := math.Min(len(cps), cfg.API.MaxQueries)
		client := &http.Client{
			Transport: &http.Transport{MaxConnsPerHost: maxCons},
			Timeout:   cfg.API.Timeout,
		}

		var (
			apiDataHandler apihandlers.APIDataHandler[oracletypes.CurrencyPair, *big.Int]
			requestHandler apihandlers.RequestHandler
		)

		switch cfg.Name {
		case binance.Name:
			apiDataHandler, err = binance.NewAPIHandler(cfg)
		case coinbaseapi.Name:
			apiDataHandler, err = coinbaseapi.NewAPIHandler(cfg)
		case coingecko.Name:
			apiDataHandler, err = coingecko.NewAPIHandler(cfg)
		case static.Name:
			apiDataHandler, err = static.NewAPIHandler(cfg)
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
		return apihandlers.NewAPIQueryHandler[oracletypes.CurrencyPair, *big.Int](
			logger,
			cfg.API,
			requestHandler,
			apiDataHandler,
			metrics,
		)
	}
}
