package bitstamp

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
)

const (
	// Name is the name of the bitstamp provider.
	Name = "bitstamp_ws"

	// WSS is the bitstamp websocket address.
	WSS = "wss://ws.bitstamp.net"

	// DefaultPingInterval is the default ping interval for the bitstamp websocket.
	DefaultPingInterval = 10 * time.Second
)

var (
	// DefaultWebSocketConfig returns the default websocket config for bitstamp.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Enabled:                       true,
		Name:                          Name,
		MaxBufferSize:                 config.DefaultMaxBufferSize,
		ReconnectionTimeout:           config.DefaultReconnectionTimeout,
		WSS:                           WSS,
		ReadBufferSize:                config.DefaultReadBufferSize,
		WriteBufferSize:               config.DefaultWriteBufferSize,
		HandshakeTimeout:              config.DefaultHandshakeTimeout,
		EnableCompression:             config.DefaultEnableCompression,
		WriteTimeout:                  config.DefaultWriteTimeout,
		ReadTimeout:                   config.DefaultReadTimeout,
		PingInterval:                  DefaultPingInterval,
		MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
		MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
	}

	// DefaultProviderConfig returns the default market config for bitstamp.
	DefaultProviderConfig = types.TickerToProviderConfig{
		constants.AVAX_USD: {
			Name:           Name,
			OffChainTicker: "avaxusd",
		},
		constants.BITCOIN_USD: {
			Name:           Name,
			OffChainTicker: "btcusd",
		},
		constants.BITCOIN_USDC: {
			Name:           Name,
			OffChainTicker: "btcusdc",
		},
		constants.BITCOIN_USDT: {
			Name:           Name,
			OffChainTicker: "btcusdt",
		},
		constants.ETHEREUM_BITCOIN: {
			Name:           Name,
			OffChainTicker: "ethbtc",
		},
		constants.ETHEREUM_USD: {
			Name:           Name,
			OffChainTicker: "ethusd",
		},
		constants.SOLANA_USD: {
			Name:           Name,
			OffChainTicker: "solusd",
		},
		constants.USDC_USDT: {
			Name:           Name,
			OffChainTicker: "usdcusdt",
		},
		constants.USDT_USD: {
			Name:           Name,
			OffChainTicker: "usdtusd",
		},
	}
)
