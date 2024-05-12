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
	Name = "coinbase_ws"

	Type = types.ConfigType

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
	}

	DefaultProviderConfig = config.ProviderConfig{
		Name:      Name,
		WebSocket: DefaultWebSocketConfig,
		Type:      Type,
	}

	// DefaultMarketConfig is the default market configuration for Coinbase.
	DefaultMarketConfig = types.CurrencyPairsToProviderTickers{
		constants.APE_USD: {
			OffChainTicker: "APE-USD",
		},
		constants.APE_USDT: {
			OffChainTicker: "APE-USDT",
		},
		constants.APTOS_USD: {
			OffChainTicker: "APT-USD",
		},
		constants.ARBITRUM_USD: {
			OffChainTicker: "ARB-USD",
		},
		constants.ATOM_USD: {
			OffChainTicker: "ATOM-USD",
		},
		constants.ATOM_USDT: {
			OffChainTicker: "ATOM-USDT",
		},
		constants.AVAX_USD: {
			OffChainTicker: "AVAX-USD",
		},
		constants.AVAX_USDT: {
			OffChainTicker: "AVAX-USDT",
		},
		constants.BCH_USD: {
			OffChainTicker: "BCH-USD",
		},
		constants.BITCOIN_USD: {
			OffChainTicker: "BTC-USD",
		},
		constants.BITCOIN_USDT: {
			OffChainTicker: "BTC-USDT",
		},
		constants.BLUR_USD: {
			OffChainTicker: "BLUR-USD",
		},
		constants.CARDANO_USD: {
			OffChainTicker: "ADA-USD",
		},
		constants.CELESTIA_USD: {
			OffChainTicker: "TIA-USD",
		},
		constants.CELESTIA_USDC: {
			OffChainTicker: "TIA-USDC",
		},
		constants.CELESTIA_USDT: {
			OffChainTicker: "TIA-USDT",
		},
		constants.CHAINLINK_USD: {
			OffChainTicker: "LINK-USD",
		},
		constants.COMPOUND_USD: {
			OffChainTicker: "COMP-USD",
		},
		constants.CURVE_USD: {
			OffChainTicker: "CRV-USD",
		},
		constants.DOGE_USD: {
			OffChainTicker: "DOGE-USD",
		},
		constants.ETC_USD: {
			OffChainTicker: "ETC-USD",
		},
		constants.ETHEREUM_BITCOIN: {
			OffChainTicker: "ETH-BTC",
		},
		constants.ETHEREUM_USD: {
			OffChainTicker: "ETH-USD",
		},
		constants.ETHEREUM_USDT: {
			OffChainTicker: "ETH-USDT",
		},
		constants.FILECOIN_USD: {
			OffChainTicker: "FIL-USD",
		},
		constants.LIDO_USD: {
			OffChainTicker: "LDO-USD",
		},
		constants.LITECOIN_USD: {
			OffChainTicker: "LTC-USD",
		},
		constants.MAKER_USD: {
			OffChainTicker: "MKR-USD",
		},
		constants.NEAR_USD: {
			OffChainTicker: "NEAR-USD",
		},
		constants.OPTIMISM_USD: {
			OffChainTicker: "OP-USD",
		},
		constants.OSMOSIS_USD: {
			OffChainTicker: "OSMO-USD",
		},
		constants.OSMOSIS_USDC: {
			OffChainTicker: "OSMO-USDC",
		},
		constants.OSMOSIS_USDT: {
			OffChainTicker: "OSMO-USDT",
		},
		constants.POLKADOT_USD: {
			OffChainTicker: "DOT-USD",
		},
		constants.POLYGON_USD: {
			OffChainTicker: "MATIC-USD",
		},
		constants.RIPPLE_USD: {
			OffChainTicker: "XRP-USD",
		},
		constants.SEI_USD: {
			OffChainTicker: "SEI-USD",
		},
		constants.SHIBA_USD: {
			OffChainTicker: "SHIB-USD",
		},
		constants.SOLANA_USD: {
			OffChainTicker: "SOL-USD",
		},
		constants.SOLANA_USDC: {
			OffChainTicker: "SOL-USDC",
		},
		constants.SOLANA_USDT: {
			OffChainTicker: "SOL-USDT",
		},
		constants.STELLAR_USD: {
			OffChainTicker: "XLM-USD",
		},
		constants.SUI_USD: {
			OffChainTicker: "SUI-USD",
		},
		constants.UNISWAP_USD: {
			OffChainTicker: "UNI-USD",
		},
		constants.USDC_USDT: {
			OffChainTicker: "USDC-USDT",
		},
		constants.USDT_USD: {
			OffChainTicker: "USDT-USD",
		},
	}
)
