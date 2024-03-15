package coinbase

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
)

const (
	// The following URLs are used for the Coinbase Websocket feed. These can be found
	// in the Coinbase documentation at https://docs.cloud.coinbase.com/exchange/docs/websocket-overview.

	// Name is the name of the Coinbase provider.
	Name = "CoinbaseProWebSocket"

	// URL is the production Coinbase Websocket URL.
	URL = "wss://ws-feed.exchange.coinbase.com"

	// URL_SANDBOX is the sandbox Coinbase Websocket URL.
	URL_SANDBOX = "wss://ws-feed-public.sandbox.exchange.coinbase.com"
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
		Type:                          types.ConfigType,
	}

	// DefaultMarketConfig is the default market configuration for Coinbase.
	DefaultMarketConfig = types.TickerToProviderConfig{
		constants.APE_USD: {
			Name:           Name,
			OffChainTicker: "APE-USD",
		},
		constants.APE_USDT: {
			Name:           Name,
			OffChainTicker: "APE-USDT",
		},
		constants.APTOS_USD: {
			Name:           Name,
			OffChainTicker: "APT-USD",
		},
		constants.ARBITRUM_USD: {
			Name:           Name,
			OffChainTicker: "ARB-USD",
		},
		constants.ATOM_USD: {
			Name:           Name,
			OffChainTicker: "ATOM-USD",
		},
		constants.ATOM_USDT: {
			Name:           Name,
			OffChainTicker: "ATOM-USDT",
		},
		constants.AVAX_USD: {
			Name:           Name,
			OffChainTicker: "AVAX-USD",
		},
		constants.AVAX_USDT: {
			Name:           Name,
			OffChainTicker: "AVAX-USDT",
		},
		constants.BCH_USD: {
			Name:           Name,
			OffChainTicker: "BCH-USD",
		},
		constants.BITCOIN_USD: {
			Name:           Name,
			OffChainTicker: "BTC-USD",
		},
		constants.BITCOIN_USDT: {
			Name:           Name,
			OffChainTicker: "BTC-USDT",
		},
		constants.BLUR_USD: {
			Name:           Name,
			OffChainTicker: "BLUR-USD",
		},
		constants.CARDANO_USD: {
			Name:           Name,
			OffChainTicker: "ADA-USD",
		},
		constants.CELESTIA_USD: {
			Name:           Name,
			OffChainTicker: "TIA-USD",
		},
		constants.CELESTIA_USDC: {
			Name:           Name,
			OffChainTicker: "TIA-USDC",
		},
		constants.CELESTIA_USDT: {
			Name:           Name,
			OffChainTicker: "TIA-USDT",
		},
		constants.CHAINLINK_USD: {
			Name:           Name,
			OffChainTicker: "LINK-USD",
		},
		constants.COMPOUND_USD: {
			Name:           Name,
			OffChainTicker: "COMP-USD",
		},
		constants.CURVE_USD: {
			Name:           Name,
			OffChainTicker: "CRV-USD",
		},
		constants.DOGE_USD: {
			Name:           Name,
			OffChainTicker: "DOGE-USD",
		},
		constants.ETC_USD: {
			Name:           Name,
			OffChainTicker: "ETC-USD",
		},
		constants.ETHEREUM_BITCOIN: {
			Name:           Name,
			OffChainTicker: "ETH-BTC",
		},
		constants.ETHEREUM_USD: {
			Name:           Name,
			OffChainTicker: "ETH-USD",
		},
		constants.ETHEREUM_USDT: {
			Name:           Name,
			OffChainTicker: "ETH-USDT",
		},
		constants.FILECOIN_USD: {
			Name:           Name,
			OffChainTicker: "FIL-USD",
		},
		constants.LIDO_USD: {
			Name:           Name,
			OffChainTicker: "LDO-USD",
		},
		constants.LITECOIN_USD: {
			Name:           Name,
			OffChainTicker: "LTC-USD",
		},
		constants.MAKER_USD: {
			Name:           Name,
			OffChainTicker: "MKR-USD",
		},
		constants.NEAR_USD: {
			Name:           Name,
			OffChainTicker: "NEAR-USD",
		},
		constants.OPTIMISM_USD: {
			Name:           Name,
			OffChainTicker: "OP-USD",
		},
		constants.OSMOSIS_USD: {
			Name:           Name,
			OffChainTicker: "OSMO-USD",
		},
		constants.OSMOSIS_USDC: {
			Name:           Name,
			OffChainTicker: "OSMO-USDC",
		},
		constants.OSMOSIS_USDT: {
			Name:           Name,
			OffChainTicker: "OSMO-USDT",
		},
		constants.POLKADOT_USD: {
			Name:           Name,
			OffChainTicker: "DOT-USD",
		},
		constants.POLYGON_USD: {
			Name:           Name,
			OffChainTicker: "MATIC-USD",
		},
		constants.RIPPLE_USD: {
			Name:           Name,
			OffChainTicker: "XRP-USD",
		},
		constants.SEI_USD: {
			Name:           Name,
			OffChainTicker: "SEI-USD",
		},
		constants.SHIBA_USD: {
			Name:           Name,
			OffChainTicker: "SHIB-USD",
		},
		constants.SOLANA_USD: {
			Name:           Name,
			OffChainTicker: "SOL-USD",
		},
		constants.SOLANA_USDC: {
			Name:           Name,
			OffChainTicker: "SOL-USDC",
		},
		constants.SOLANA_USDT: {
			Name:           Name,
			OffChainTicker: "SOL-USDT",
		},
		constants.STELLAR_USD: {
			Name:           Name,
			OffChainTicker: "XLM-USD",
		},
		constants.SUI_USD: {
			Name:           Name,
			OffChainTicker: "SUI-USD",
		},
		constants.UNISWAP_USD: {
			Name:           Name,
			OffChainTicker: "UNI-USD",
		},
		constants.USDC_USDT: {
			Name:           Name,
			OffChainTicker: "USDC-USDT",
		},
		constants.USDT_USD: {
			Name:           Name,
			OffChainTicker: "USDT-USD",
		},
	}
)
