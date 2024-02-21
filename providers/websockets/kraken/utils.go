package kraken

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
)

const (
	// URL is the websocket URL for Kraken. You can find the documentation here:
	// https://docs.kraken.com/websockets/. Kraken provides an authenticated and
	// unauthenticated websocket. The URLs defined below are all unauthenticated.

	// Name is the name of the Kraken provider.
	Name = "kraken"

	// URL is the production websocket URL for Kraken.
	URL = "wss://ws.kraken.com"

	// URL_BETA is the demo websocket URL for Kraken.
	URL_BETA = "wss://beta-ws.kraken.com"
)

var (
	// DefaultWebSocketConfig is the default configuration for the Kraken Websocket.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Name:                          Name,
		Enabled:                       true,
		MaxBufferSize:                 1000,
		ReconnectionTimeout:           10 * time.Second,
		WSS:                           URL,
		ReadBufferSize:                config.DefaultReadBufferSize,
		WriteBufferSize:               config.DefaultWriteBufferSize,
		HandshakeTimeout:              config.DefaultHandshakeTimeout,
		EnableCompression:             config.DefaultEnableCompression,
		ReadTimeout:                   config.DefaultReadTimeout,
		WriteTimeout:                  config.DefaultWriteTimeout,
		PingInterval:                  config.DefaultPingInterval,
		MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
		MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
	}

	// DefaultMarketConfig is the default market configuration for Kraken.
	DefaultMarketConfig = types.TickerToProviderConfig{
		constants.ATOM_USD: {
			Name:           Name,
			OffChainTicker: "ATOM/USD",
		},
		constants.AVAX_USD: {
			Name:           Name,
			OffChainTicker: "AVAX/USD",
		},
		constants.AVAX_USDT: {
			Name:           Name,
			OffChainTicker: "AVAX/USDT",
		},
		constants.BITCOIN_USD: {
			Name:           Name,
			OffChainTicker: "BITCOIN/USD",
		},
		constants.BITCOIN_USDC: {
			Name:           Name,
			OffChainTicker: "BITCOIN/USDC",
		},
		constants.BITCOIN_USDT: {
			Name:           Name,
			OffChainTicker: "BITCOIN/USDT",
		},
		constants.CELESTIA_USD: {
			Name:           Name,
			OffChainTicker: "CELESTIA/USD",
		},
		constants.DYDX_USD: {
			Name:           Name,
			OffChainTicker: "DYDX/USD",
		},
		constants.ETHEREUM_BITCOIN: {
			Name:           Name,
			OffChainTicker: "ETHEREUM/BITCOIN",
		},
		constants.ETHEREUM_USD: {
			Name:           Name,
			OffChainTicker: "ETHEREUM/USD",
		},
		constants.ETHEREUM_USDC: {
			Name:           Name,
			OffChainTicker: "ETHEREUM/USDC",
		},
		constants.ETHEREUM_USDT: {
			Name:           Name,
			OffChainTicker: "ETHEREUM/USDT",
		},
		constants.SOLANA_USD: {
			Name:           Name,
			OffChainTicker: "SOLANA/USD",
		},
		constants.SOLANA_USDT: {
			Name:           Name,
			OffChainTicker: "SOLANA/USDT",
		},
		constants.USDC_USD: {
			Name:           Name,
			OffChainTicker: "USDC/USD",
		},
		constants.USDC_USDT: {
			Name:           Name,
			OffChainTicker: "USDC/USDT",
		},
		constants.USDT_USD: {
			Name:           Name,
			OffChainTicker: "USDT/USD",
		},
	}
)
