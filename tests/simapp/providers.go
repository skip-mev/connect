package simapp

import (
	"fmt"
	"math/big"
	"net/http"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/pkg/math"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/providers/apis/binance"
	coinbaseapi "github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/apis/coingecko"
	"github.com/skip-mev/slinky/providers/base"
	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
	apimetrics "github.com/skip-mev/slinky/providers/base/api/metrics"
	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"
	wshandlers "github.com/skip-mev/slinky/providers/base/websocket/handlers"
	wsmetrics "github.com/skip-mev/slinky/providers/base/websocket/metrics"
	"github.com/skip-mev/slinky/providers/static"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/providers/websockets/bitfinex"
	"github.com/skip-mev/slinky/providers/websockets/bitstamp"
	"github.com/skip-mev/slinky/providers/websockets/bybit"
	coinbasews "github.com/skip-mev/slinky/providers/websockets/coinbase"
	"github.com/skip-mev/slinky/providers/websockets/cryptodotcom"
	"github.com/skip-mev/slinky/providers/websockets/gate"
	"github.com/skip-mev/slinky/providers/websockets/huobi"
	"github.com/skip-mev/slinky/providers/websockets/kraken"
	"github.com/skip-mev/slinky/providers/websockets/kucoin"
	"github.com/skip-mev/slinky/providers/websockets/mexc"
	"github.com/skip-mev/slinky/providers/websockets/okx"
)

// DefaultProviderFactory returns a sample implementation of the provider factory. This provider
// factory function returns providers that are API & websocket based.
func DefaultProviderFactory() providertypes.ProviderFactory[slinkytypes.CurrencyPair, *big.Int] {
	return func(logger *zap.Logger, cfg config.OracleConfig) ([]providertypes.Provider[slinkytypes.CurrencyPair, *big.Int], error) {
		if err := cfg.ValidateBasic(); err != nil {
			return nil, err
		}

		cps := cfg.Market.GetCurrencyPairs()

		// Create the metrics that are used by the providers.
		mWebSocket := wsmetrics.NewWebSocketMetricsFromConfig(cfg.Metrics)
		mAPI := apimetrics.NewAPIMetricsFromConfig(cfg.Metrics)
		mProviders := providermetrics.NewProviderMetricsFromConfig(cfg.Metrics)

		// Create the providers.
		providers := make([]providertypes.Provider[slinkytypes.CurrencyPair, *big.Int], 0)
		for _, p := range cfg.Providers {
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
	cps []slinkytypes.CurrencyPair,
	mAPI apimetrics.APIMetrics,
	mProvider providermetrics.ProviderMetrics,
) (providertypes.Provider[slinkytypes.CurrencyPair, *big.Int], error) {
	// Validate the provider config.
	err := cfg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	// Filter the currency pairs to only include the ones that are configured in the provider
	// config.
	filteredCPs, err := filterForConfiguredCurrencyPairs(logger, cps, cfg)
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
		apiDataHandler apihandlers.APIDataHandler[slinkytypes.CurrencyPair, *big.Int]
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
	apiQueryHandler, err := apihandlers.NewAPIQueryHandler[slinkytypes.CurrencyPair, *big.Int](
		logger,
		cfg.API,
		requestHandler,
		apiDataHandler,
		mAPI,
	)
	if err != nil {
		return nil, err
	}

	// Create the provider.
	return base.NewProvider[slinkytypes.CurrencyPair, *big.Int](
		base.WithName[slinkytypes.CurrencyPair, *big.Int](cfg.Name),
		base.WithLogger[slinkytypes.CurrencyPair, *big.Int](logger),
		base.WithAPIQueryHandler(apiQueryHandler),
		base.WithAPIConfig[slinkytypes.CurrencyPair, *big.Int](cfg.API),
		base.WithIDs[slinkytypes.CurrencyPair, *big.Int](filteredCPs),
		base.WithMetrics[slinkytypes.CurrencyPair, *big.Int](mProvider),
	)
}

// webSocketProviderFromProviderConfig returns a websocket provider from a provider config. These providers are
// NOT production ready and are only meant for testing purposes.
func webSocketProviderFromProviderConfig(
	logger *zap.Logger,
	cfg config.ProviderConfig,
	cps []slinkytypes.CurrencyPair,
	wsMetrics wsmetrics.WebSocketMetrics,
	pMetrics providermetrics.ProviderMetrics,
) (providertypes.Provider[slinkytypes.CurrencyPair, *big.Int], error) {
	// Validate the provider config.
	err := cfg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	// Filter the currency pairs to only include the ones that are configured in the provider
	// config.
	filteredCPs, err := filterForConfiguredCurrencyPairs(logger, cps, cfg)
	if err != nil {
		return nil, err
	}

	// Create the underlying client that can be utilized by websocket providers that need to
	// interact with an API.
	maxCons := math.Min(len(cps), cfg.API.MaxQueries)
	client := &http.Client{
		Transport: &http.Transport{MaxConnsPerHost: maxCons},
		Timeout:   cfg.API.Timeout,
	}

	var (
		requestHandler apihandlers.RequestHandler
		wsDataHandler  wshandlers.WebSocketDataHandler[slinkytypes.CurrencyPair, *big.Int]
		connHandler    wshandlers.WebSocketConnHandler
	)

	switch cfg.Name {
	case bitfinex.Name:
		wsDataHandler, err = bitfinex.NewWebSocketDataHandler(logger, cfg)
	case bitstamp.Name:
		wsDataHandler, err = bitstamp.NewWebSocketDataHandler(logger, cfg)
	case bybit.Name:
		wsDataHandler, err = bybit.NewWebSocketDataHandler(logger, cfg)
	case coinbasews.Name:
		wsDataHandler, err = coinbasews.NewWebSocketDataHandler(logger, cfg)
	case cryptodotcom.Name:
		wsDataHandler, err = cryptodotcom.NewWebSocketDataHandler(logger, cfg)
	case gate.Name:
		wsDataHandler, err = gate.NewWebSocketDataHandler(logger, cfg)
	case huobi.Name:
		wsDataHandler, err = huobi.NewWebSocketDataHandler(logger, cfg)
	case kraken.Name:
		wsDataHandler, err = kraken.NewWebSocketDataHandler(logger, cfg)
	case kucoin.Name:
		// Create the KuCoin websocket data handler.
		wsDataHandler, err = kucoin.NewWebSocketDataHandler(logger, cfg)
		if err != nil {
			return nil, err
		}

		// The request handler requires POST requests when first establishing the connection.
		requestHandler, err = apihandlers.NewRequestHandlerImpl(
			client,
			apihandlers.WithHTTPMethod(http.MethodPost),
		)
		if err != nil {
			return nil, err
		}

		// Create the KuCoin websocket connection handler.
		connHandler, err = wshandlers.NewWebSocketHandlerImpl(
			cfg.WebSocket,
			wshandlers.WithPreDialHook(kucoin.PreDialHook(cfg.API, requestHandler)),
		)
	case mexc.Name:
		wsDataHandler, err = mexc.NewWebSocketDataHandler(logger, cfg)
	case okx.Name:
		wsDataHandler, err = okx.NewWebSocketDataHandler(logger, cfg)
	default:
		return nil, fmt.Errorf("unknown provider: %s", cfg.Name)
	}
	if err != nil {
		return nil, err
	}

	// If a custom request handler is not provided, create a new default one.
	if connHandler == nil {
		connHandler, err = wshandlers.NewWebSocketHandlerImpl(cfg.WebSocket)
		if err != nil {
			return nil, err
		}
	}

	// Create the websocket query handler which encapsulates all fetching and parsing logic.
	wsQueryHandler, err := wshandlers.NewWebSocketQueryHandler[slinkytypes.CurrencyPair, *big.Int](
		logger,
		cfg.WebSocket,
		wsDataHandler,
		connHandler,
		wsMetrics,
	)
	if err != nil {
		return nil, err
	}

	// Create the provider.
	return base.NewProvider[slinkytypes.CurrencyPair, *big.Int](
		base.WithName[slinkytypes.CurrencyPair, *big.Int](cfg.Name),
		base.WithLogger[slinkytypes.CurrencyPair, *big.Int](logger),
		base.WithWebSocketQueryHandler(wsQueryHandler),
		base.WithWebSocketConfig[slinkytypes.CurrencyPair, *big.Int](cfg.WebSocket),
		base.WithIDs[slinkytypes.CurrencyPair, *big.Int](filteredCPs),
		base.WithMetrics[slinkytypes.CurrencyPair, *big.Int](pMetrics),
	)
}

// filterForConfiguredCurrencyPairs returns the set of currency pairs that are configured in the
// providers config.
func filterForConfiguredCurrencyPairs(
	logger *zap.Logger,
	cps []slinkytypes.CurrencyPair,
	cfg config.ProviderConfig,
) ([]slinkytypes.CurrencyPair, error) {
	filteredCps := make([]slinkytypes.CurrencyPair, 0)

	for _, cp := range cps {
		if _, ok := cfg.Market.CurrencyPairToMarketConfigs[cp.String()]; ok {
			logger.Debug("provider supports currency pair", zap.String("currency_pair", cp.String()), zap.String("provider", cfg.Name))
			filteredCps = append(filteredCps, cp)
		} else {
			logger.Debug("provider does not support currency pair", zap.String("currency_pair", cp.String()), zap.String("provider", cfg.Name))
		}
	}

	if len(filteredCps) == 0 {
		return nil, fmt.Errorf("no currency pairs supported by provider: %s", cfg.Name)
	}

	return filteredCps, nil
}
