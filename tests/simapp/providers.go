package simapp

import (
	"fmt"
	"math/big"
	"net/http"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/providers/base"
	"github.com/skip-mev/slinky/providers/base/api/handlers"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
	"github.com/skip-mev/slinky/providers/coinbase"
	"github.com/skip-mev/slinky/providers/coingecko"
	"github.com/skip-mev/slinky/providers/static"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// DefaultAPIProviderFactory returns a sample implementation of the provider factory. This provider
// factory function only returns providers the are API based.
func DefaultAPIProviderFactory() providertypes.ProviderFactory[oracletypes.CurrencyPair, *big.Int] {
	return func(logger *zap.Logger, oracleCfg config.OracleConfig, metricsCfg config.OracleMetricsConfig) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
		if err := oracleCfg.ValidateBasic(); err != nil {
			return nil, err
		}

		m := metrics.NewAPIMetricsFromConfig(metricsCfg)
		cps := oracleCfg.CurrencyPairs

		var (
			err       error
			providers = make([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], len(oracleCfg.Providers))
		)
		for i, p := range oracleCfg.Providers {
			if providers[i], err = providerFromProviderConfig(logger, p, cps, m); err != nil {
				return nil, err
			}
		}

		return providers, nil
	}
}

// providerFromProviderConfig returns a provider from a provider config. These providers are
// NOT production ready and are only meant for testing purposes.
func providerFromProviderConfig(
	logger *zap.Logger,
	cfg config.ProviderConfig,
	cps []oracletypes.CurrencyPair,
	m metrics.APIMetrics,
) (providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
	// Validate the provider config.
	err := cfg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	// Create the underlying client that will be used to fetch data from the API. This client
	// will limit the number of concurrent connections and uses the configured timeout to
	// ensure requests do not hang.
	maxCons := math.Min(len(cps), cfg.API.MaxQueries)
	client := &http.Client{
		Transport: &http.Transport{MaxConnsPerHost: maxCons},
		Timeout:   cfg.API.Timeout,
	}

	var (
		apiDataHandler handlers.APIDataHandler[oracletypes.CurrencyPair, *big.Int]
		requestHandler handlers.RequestHandler
	)

	switch cfg.Name {
	case "coingecko":
		apiDataHandler, err = coingecko.NewCoinGeckoAPIHandler(cfg)
	case "coinbase":
		apiDataHandler, err = coinbase.NewCoinBaseAPIHandler(cfg)
	case "static-mock-provider":
		apiDataHandler, err = static.NewStaticMockAPIHandler(cfg)
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

	if apiDataHandler == nil {
		return nil, fmt.Errorf("failed to create api data handler for provider %s", cfg.Name)
	}

	// If a custom request handler is not provided, create a new default one.
	if requestHandler == nil {
		requestHandler = handlers.NewRequestHandlerImpl(client)
	}

	// Create the API query handler which encapsulates all of the fetching and parsing logic.
	apiQueryHandler, err := handlers.NewAPIQueryHandler[oracletypes.CurrencyPair, *big.Int](
		logger,
		requestHandler,
		apiDataHandler,
		m,
	)
	if err != nil {
		return nil, err
	}

	// Create the provider.
	return base.NewProvider[oracletypes.CurrencyPair, *big.Int](
		cfg,
		base.WithLogger[oracletypes.CurrencyPair, *big.Int](logger),
		base.WithAPIQueryHandler(apiQueryHandler),
		base.WithIDs[oracletypes.CurrencyPair, *big.Int](cps),
	)
}
