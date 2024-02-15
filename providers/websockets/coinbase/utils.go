package coinbase

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/constants"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

const (
	// The following URLs are used for the Coinbase Websocket feed. These can be found
	// in the Coinbase documentation at https://docs.cloud.coinbase.com/exchange/docs/websocket-overview.

	// Name is the name of the Coinbase provider.
	Name = "coinbase_websocket"

	// URL is the production Coinbase Websocket URL.
	URL = "wss://ws-feed.exchange.coinbase.com"

	// URL_SANDBOX is the sandbox Coinbase Websocket URL.
	URL_SANDBOX = "wss://ws-feed-public.sandbox.exchange.coinbase.com" //nolint
)

const (
	// The following websocket configuration values were taken from the Coinbase documentation
	// at https://docs.cloud.coinbase.com/exchange/docs/websocket-overview.

	// DefaultEnabledCompression is the default enabled compression for the Coinbase Websocket.
	// It is recommended to set this as true (please see the Coinbase documentation for more).
	DefaultEnabledCompression = false

	// DefaultWriteTimeout is the default write timeout for the Coinbase Websocket.
	// As recommended by Coinbase, this is set to 5 seconds.
	DefaultWriteTimeout = 5 * time.Second
)

var (
	// DefaultWebSocketConfig is the default configuration for the Coinbase Websocket.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Enabled:                       true,
		Name:                          Name,
		MaxBufferSize:                 config.DefaultMaxBufferSize,
		ReconnectionTimeout:           config.DefaultReconnectionTimeout,
		WSS:                           URL,
		ReadBufferSize:                config.DefaultReadBufferSize,
		WriteBufferSize:               config.DefaultWriteBufferSize,
		HandshakeTimeout:              config.DefaultHandshakeTimeout,
		EnableCompression:             DefaultEnabledCompression,
		WriteTimeout:                  DefaultWriteTimeout,
		ReadTimeout:                   config.DefaultReadTimeout,
		PingInterval:                  config.DefaultPingInterval,
		MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
		MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
	}

	// DefaultMarketConfig is the default market configuration for Coinbase.
	DefaultMarketConfig = mmtypes.MarketConfig{
		Name: Name,
		TickerConfigs: map[string]mmtypes.TickerConfig{
			"ATOM/USD": {
				Ticker:         constants.ATOM_USD,
				OffChainTicker: "ATOM-USD",
			},
			"ATOM/USDC": {
				Ticker:         constants.ATOM_USDC,
				OffChainTicker: "ATOM-USDC",
			},
			"ATOM/USDT": {
				Ticker:         constants.ATOM_USDT,
				OffChainTicker: "ATOM-USDT",
			},
			"AVAX/USD": {
				Ticker:         constants.AVAX_USD,
				OffChainTicker: "AVAX-USD",
			},
			"AVAX/USDC": {
				Ticker:         constants.AVAX_USDC,
				OffChainTicker: "AVAX-USDC",
			},
			"AVAX/USDT": {
				Ticker:         constants.AVAX_USDT,
				OffChainTicker: "AVAX-USDT",
			},
			"BITCOIN/USD": {
				Ticker:         constants.BITCOIN_USD,
				OffChainTicker: "BTC-USD",
			},
			"BITCOIN/USDC": {
				Ticker:         constants.BITCOIN_USDC,
				OffChainTicker: "BTC-USDC",
			},
			"BITCOIN/USDT": {
				Ticker:         constants.BITCOIN_USDT,
				OffChainTicker: "BTC-USDT",
			},
			"CELESTIA/USD": {
				Ticker:         constants.CELESTIA_USD,
				OffChainTicker: "TIA-USD",
			},
			"CELESTIA/USDC": {
				Ticker:         constants.CELESTIA_USDC,
				OffChainTicker: "TIA-USDC",
			},
			"CELESTIA/USDT": {
				Ticker:         constants.CELESTIA_USDT,
				OffChainTicker: "TIA-USDT",
			},
			"DYDX/USD": {
				Ticker:         constants.DYDX_USD,
				OffChainTicker: "DYDX-USD",
			},
			"DYDX/USDC": {
				Ticker:         constants.DYDX_USDC,
				OffChainTicker: "DYDX-USDC",
			},
			"DYDX/USDT": {
				Ticker:         constants.DYDX_USDT,
				OffChainTicker: "DYDX-USDT",
			},
			"ETHEREUM/BITCOIN": {
				Ticker:         constants.ETHEREUM_BITCOIN,
				OffChainTicker: "ETH-BTC",
			},
			"ETHEREUM/USD": {
				Ticker:         constants.ETHEREUM_USD,
				OffChainTicker: "ETH-USD",
			},
			"ETHEREUM/USDC": {
				Ticker:         constants.ETHEREUM_USDC,
				OffChainTicker: "ETH-USDC",
			},
			"ETHEREUM/USDT": {
				Ticker:         constants.ETHEREUM_USDT,
				OffChainTicker: "ETH-USDT",
			},
			"OSMOSIS/USD": {
				Ticker:         constants.OSMOSIS_USD,
				OffChainTicker: "OSMO-USD",
			},
			"OSMOSIS/USDC": {
				Ticker:         constants.OSMOSIS_USDC,
				OffChainTicker: "OSMO-USDC",
			},
			"OSMOSIS/USDT": {
				Ticker:         constants.OSMOSIS_USDT,
				OffChainTicker: "OSMO-USDT",
			},
			"SOLANA/USD": {
				Ticker:         constants.SOLANA_USD,
				OffChainTicker: "SOL-USD",
			},
			"SOLANA/USDC": {
				Ticker:         constants.SOLANA_USDC,
				OffChainTicker: "SOL-USDC",
			},
			"SOLANA/USDT": {
				Ticker:         constants.SOLANA_USDT,
				OffChainTicker: "SOL-USDT",
			},
			"USDC/USD": {
				Ticker:         constants.USDC_USD,
				OffChainTicker: "USDC-USD",
			},
			"USDC/USDT": {
				Ticker:         constants.USDC_USDT,
				OffChainTicker: "USDC-USDT",
			},
			"USDT/USD": {
				Ticker:         constants.USDT_USD,
				OffChainTicker: "USDT-USD",
			},
		},
	}
)
