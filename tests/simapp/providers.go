package simapp

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/skip-mev/slinky/providers/websockets/bybit"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/providers/apis/binanceus"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/apis/coingecko"
	"github.com/skip-mev/slinky/providers/base"
	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
	apimetrics "github.com/skip-mev/slinky/providers/base/api/metrics"
	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"
	wshandlers "github.com/skip-mev/slinky/providers/base/websocket/handlers"
	wsmetrics "github.com/skip-mev/slinky/providers/base/websocket/metrics"
	"github.com/skip-mev/slinky/providers/static"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/websockets/cryptodotcom"
	"github.com/skip-mev/slinky/providers/websockets/okx"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// DefaultProviderFactory returns a sample implementation of the provider factory. This provider
// factory function returns providers the are API & web socket based.
func DefaultProviderFactory() providertypes.ProviderFactory[oracletypes.CurrencyPair, *big.Int] {
	return func(logger *zap.Logger, oracleCfg config.OracleConfig, metricsCfg config.OracleMetricsConfig) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
		if err := oracleCfg.ValidateBasic(); err != nil {
			return nil, err
		}

		cps := oracleCfg.CurrencyPairs

		// Create the metrics that are used by the providers.
		mWebSocket := wsmetrics.NewWebSocketMetricsFromConfig(metricsCfg)
		mAPI := apimetrics.NewAPIMetricsFromConfig(metricsCfg)
		mProviders := providermetrics.NewProviderMetricsFromConfig(metricsCfg)

		providers := make([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], 0)
		for _, p := range oracleCfg.Providers {
			switch {
			case p.API.Enabled:
				provider, err := apiProviderFromProviderConfig(logger, p, cps, mAPI, mProviders)
				if err != nil {
					return nil, err
				}

				providers = append(providers, provider)
			case p.WebSocket.Enabled:
				provider, err := webSocketProviderFromProviderConfig(logger, p, cps, mWebSocket, mProviders)
				if err != nil {
					return nil, err
				}

				providers = append(providers, provider)
			default:
				logger.Info("unknown provider type", zap.String("provider", p.Name))
				return nil, fmt.Errorf("unknown provider type: %s", p.Name)
			}
		}

		return providers, nil
	}
}

// apiProviderFromProviderConfig returns an API provider from a provider config. These providers are
// NOT production ready and are only meant for testing purposes.
func apiProviderFromProviderConfig(
	logger *zap.Logger,
	cfg config.ProviderConfig,
	cps []oracletypes.CurrencyPair,
	mAPI apimetrics.APIMetrics,
	mProvider providermetrics.ProviderMetrics,
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
		apiDataHandler apihandlers.APIDataHandler[oracletypes.CurrencyPair, *big.Int]
		requestHandler apihandlers.RequestHandler
	)

	switch cfg.Name {
	case coingecko.Name:
		apiDataHandler, err = coingecko.NewCoinGeckoAPIHandler(cfg)
	case coinbase.Name:
		apiDataHandler, err = coinbase.NewCoinBaseAPIHandler(cfg)
	case binanceus.Name:
		apiDataHandler, err = binanceus.NewBinanceUSAPIHandler(cfg)
	case static.Name:
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

	// If a custom request handler is not provided, create a new default one.
	if requestHandler == nil {
		requestHandler = apihandlers.NewRequestHandlerImpl(client)
	}

	// Create the API query handler which encapsulates all of the fetching and parsing logic.
	apiQueryHandler, err := apihandlers.NewAPIQueryHandler[oracletypes.CurrencyPair, *big.Int](
		logger,
		requestHandler,
		apiDataHandler,
		mAPI,
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
		base.WithMetrics[oracletypes.CurrencyPair, *big.Int](mProvider),
	)
}

// webSocketProviderFromProviderConfig returns a websocket provider from a provider config. These providers are
// NOT production ready and are only meant for testing purposes.
func webSocketProviderFromProviderConfig(
	logger *zap.Logger,
	cfg config.ProviderConfig,
	cps []oracletypes.CurrencyPair,
	wsMetrics wsmetrics.WebSocketMetrics,
	pMetrics providermetrics.ProviderMetrics,
) (providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
	// Validate the provider config.
	err := cfg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	var (
		wsDataHandler wshandlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int]
		connHandler   wshandlers.WebSocketConnHandler
	)

	switch cfg.Name {
	case cryptodotcom.Name:
		wsDataHandler, err = cryptodotcom.NewWebSocketDataHandlerFromConfig(logger, cfg)
	case okx.Name:
		wsDataHandler, err = okx.NewWebSocketDataHandlerFromConfig(logger, cfg)
	case bybit.Name:
		wsDataHandler, err = bybit.NewWebSocketDataHandlerFromConfig(logger, cfg)
		connHandler = bybit.NewWebSocketHandler(logger)
	default:
		return nil, fmt.Errorf("unknown provider: %s", cfg.Name)
	}
	if err != nil {
		return nil, err
	}

	// If a custom request handler is not provided, create a new default one.
	if connHandler == nil {
		connHandler = wshandlers.NewWebSocketHandlerImpl()
	}

	// Create the web socket query handler which encapsulates all of the fetching and parsing logic.
	wsQueryHandler, err := wshandlers.NewWebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int](
		logger,
		wsDataHandler,
		connHandler,
		wsMetrics,
	)
	if err != nil {
		return nil, err
	}

	// Create the provider.
	return base.NewProvider[oracletypes.CurrencyPair, *big.Int](
		cfg,
		base.WithLogger[oracletypes.CurrencyPair, *big.Int](logger),
		base.WithWebSocketQueryHandler(wsQueryHandler),
		base.WithIDs[oracletypes.CurrencyPair, *big.Int](cps),
		base.WithMetrics[oracletypes.CurrencyPair, *big.Int](pMetrics),
	)
}
