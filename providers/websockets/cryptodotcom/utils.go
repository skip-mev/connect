package cryptodotcom

import (
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
)

const (
	// URL is the URL used to connect to the Crypto.com websocket API. This can be found here
	// https://exchange-docs.crypto.com/exchange/v1/rest-ws/index.html?javascript#websocket-root-endpoints
	// Note that Crypto.com offers a sandbox and production environment.

	// Name is the name of the Crypto.com provider.
	Name = "crypto_dot_com_ws"

	// URL_PROD is the URL used to connect to the Crypto.com production websocket API.
	URL_PROD = "wss://stream.crypto.com/exchange/v1/market"

	// URL_SANDBOX is the URL used to connect to the Crypto.com sandbox websocket API. This will
	// return static prices.
	URL_SANDBOX = "wss://uat-stream.3ona.co/exchange/v1/market"
)

var (
	// DefaultWebSocketConfig is the default configuration for the Crypto.com Websocket.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Name:                          Name,
		Enabled:                       true,
		MaxBufferSize:                 config.DefaultMaxBufferSize,
		ReconnectionTimeout:           config.DefaultReconnectionTimeout,
		WSS:                           URL_PROD,
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

	// DefaultMarketConfig is the default market configuration for Crypto.com.
	DefaultMarketConfig = types.TickerToProviderConfig{
		constants.ATOM_USD: {
			Name:           Name,
			OffChainTicker: "ATOMUSD-PERP",
		},
		constants.ATOM_USDT: {
			Name:           Name,
			OffChainTicker: "ATOM_USDT",
		},
		constants.AVAX_USD: {
			Name:           Name,
			OffChainTicker: "AVAXUSD-PERP",
		},
		constants.AVAX_USDT: {
			Name:           Name,
			OffChainTicker: "AVAX_USDT",
		},
		constants.BITCOIN_USD: {
			Name:           Name,
			OffChainTicker: "BTCUSD-PERP",
		},
		constants.BITCOIN_USDT: {
			Name:           Name,
			OffChainTicker: "BTC_USDT",
		},
		constants.CELESTIA_USD: {
			Name:           Name,
			OffChainTicker: "TIAUSD-PERP",
		},
		constants.CELESTIA_USDT: {
			Name:           Name,
			OffChainTicker: "TIA_USDT",
		},
		constants.DYDX_USD: {
			Name:           Name,
			OffChainTicker: "DYDXUSD-PERP",
		},
		constants.DYDX_USDT: {
			Name:           Name,
			OffChainTicker: "DYDX_USDT",
		},
		constants.ETHEREUM_BITCOIN: {
			Name:           Name,
			OffChainTicker: "ETH_BTC",
		},
		constants.ETHEREUM_USD: {
			Name:           Name,
			OffChainTicker: "ETHUSD-PERP",
		},
		constants.ETHEREUM_USDT: {
			Name:           Name,
			OffChainTicker: "ETH_USDT",
		},
		constants.OSMOSIS_USD: {
			Name:           Name,
			OffChainTicker: "OSMO_USD",
		},
		constants.SOLANA_USD: {
			Name:           Name,
			OffChainTicker: "SOLUSD-PERP",
		},
		constants.SOLANA_USDT: {
			Name:           Name,
			OffChainTicker: "SOL_USDT",
		},
		constants.USDT_USD: {
			Name:           Name,
			OffChainTicker: "USDT_USD",
		},
	}
)
