package oracle

import (
	"fmt"
	"math/big"
	"net/http"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/pkg/math"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
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
)

// WebSocketQueryHandlerFactory returns a sample implementation of the websocket query handler
// factory. Specifically, this factory function returns websocket query handlers that are used to
// fetch data from the price providers.
func WebSocketQueryHandlerFactory() factory.WebSocketQueryHandlerFactory[slinkytypes.CurrencyPair, *big.Int] {
	return func(logger *zap.Logger, cfg config.ProviderConfig, wsMetrics wsmetrics.WebSocketMetrics) (wshandlers.WebSocketQueryHandler[slinkytypes.CurrencyPair, *big.Int], error) {
		// Validate the provider config.
		err := cfg.ValidateBasic()
		if err != nil {
			return nil, err
		}

		// Create the underlying client that can be utilized by websocket providers that need to
		// interact with an API.
		cps := cfg.Market.GetCurrencyPairs()
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
		return wshandlers.NewWebSocketQueryHandler[slinkytypes.CurrencyPair, *big.Int](
			logger,
			cfg.WebSocket,
			wsDataHandler,
			connHandler,
			wsMetrics,
		)
	}
}
