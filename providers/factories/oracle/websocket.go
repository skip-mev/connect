package oracle

import (
	"fmt"
	"math/big"
	"net/http"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/pkg/math"
	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
	wshandlers "github.com/skip-mev/slinky/providers/base/websocket/handlers"
	wsmetrics "github.com/skip-mev/slinky/providers/base/websocket/metrics"
	"github.com/skip-mev/slinky/providers/types/factory"
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
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// WebSocketQueryHandlerFactory returns a sample implementation of the websocket query handler
// factory. Specifically, this factory function returns websocket query handlers that are used to
// fetch data from the price providers.
func WebSocketQueryHandlerFactory(
	aggConfig mmtypes.AggregateMarketConfig,
) factory.WebSocketQueryHandlerFactory[mmtypes.Ticker, *big.Int] {
	return func(
		logger *zap.Logger,
		cfg config.ProviderConfig,
		wsMetrics wsmetrics.WebSocketMetrics,
	) (wshandlers.WebSocketQueryHandler[mmtypes.Ticker, *big.Int], error) {
		// If the websocket is not enabled, return an error.
		if !cfg.WebSocket.Enabled {
			return nil, fmt.Errorf("websocket for provider %s is not enabled", cfg.Name)
		}

		// Validate the websocket provider config.
		if err := cfg.WebSocket.ValidateBasic(); err != nil {
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

		// Create the underlying client that can be utilized by websocket providers that need to
		// interact with an API.
		maxCons := math.Min(len(market.TickerConfigs), cfg.API.MaxQueries)
		client := &http.Client{
			Transport: &http.Transport{MaxConnsPerHost: maxCons},
			Timeout:   cfg.API.Timeout,
		}

		var (
			requestHandler apihandlers.RequestHandler
			wsDataHandler  wshandlers.WebSocketDataHandler[mmtypes.Ticker, *big.Int]
			connHandler    wshandlers.WebSocketConnHandler
			err            error
		)

		logger = logger.With(zap.String("websocket_data_handler", cfg.Name))
		switch cfg.Name {
		case bitfinex.Name:
			wsDataHandler, err = bitfinex.NewWebSocketDataHandler(logger, market, cfg.WebSocket)
		case bitstamp.Name:
			wsDataHandler, err = bitstamp.NewWebSocketDataHandler(logger, market, cfg.WebSocket)
		case bybit.Name:
			wsDataHandler, err = bybit.NewWebSocketDataHandler(logger, market, cfg.WebSocket)
		case coinbasews.Name:
			wsDataHandler, err = coinbasews.NewWebSocketDataHandler(logger, market, cfg.WebSocket)
		case cryptodotcom.Name:
			wsDataHandler, err = cryptodotcom.NewWebSocketDataHandler(logger, market, cfg.WebSocket)
		case gate.Name:
			wsDataHandler, err = gate.NewWebSocketDataHandler(logger, market, cfg.WebSocket)
		case huobi.Name:
			wsDataHandler, err = huobi.NewWebSocketDataHandler(logger, market, cfg.WebSocket)
		case kraken.Name:
			wsDataHandler, err = kraken.NewWebSocketDataHandler(logger, market, cfg.WebSocket)
		case kucoin.Name:
			// Create the KuCoin websocket data handler.
			wsDataHandler, err = kucoin.NewWebSocketDataHandler(logger, market, cfg.WebSocket)
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
			wsDataHandler, err = mexc.NewWebSocketDataHandler(logger, market, cfg.WebSocket)
		case okx.Name:
			wsDataHandler, err = okx.NewWebSocketDataHandler(logger, market, cfg.WebSocket)
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
		return wshandlers.NewWebSocketQueryHandler[mmtypes.Ticker, *big.Int](
			logger,
			cfg.WebSocket,
			wsDataHandler,
			connHandler,
			wsMetrics,
		)
	}
}
