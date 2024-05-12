package kucoin

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
)

const (
	// Name is the name of the KuCoin provider.
	Name = "kucoin_ws"

	Type = types.ConfigType

	// WSSEndpoint contains the endpoint format for Kucoin websocket API. Specifically
	// this inputs the dynamically generated token from the user and the endpoint.
	WSSEndpoint = "%s?token=%s"

	// WSS is the websocket URL for Kucoin. Note that this may change as the URL is
	// dynamically generated. A token is required to connect to the websocket feed.
	WSS = "wss://ws-api-spot.kucoin.com/"

	// URL is the Kucoin websocket URL. This URL specifically points to the public
	// spot and maring REST API.
	URL = "https://api.kucoin.com"

	// DefaultPingInterval is the default ping interval for the KuCoin websocket.
	DefaultPingInterval = 10 * time.Second
)

var (
	// DefaultWebSocketConfig defines the default websocket config for Kucoin.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Enabled:                       true,
		MaxBufferSize:                 config.DefaultMaxBufferSize,
		ReconnectionTimeout:           config.DefaultReconnectionTimeout,
		WSS:                           WSS, // Note that this may change as the URL is dynamically generated.
		Name:                          Name,
		ReadBufferSize:                config.DefaultReadBufferSize,
		WriteBufferSize:               config.DefaultWriteBufferSize,
		HandshakeTimeout:              config.DefaultHandshakeTimeout,
		EnableCompression:             config.DefaultEnableCompression,
		ReadTimeout:                   config.DefaultReadTimeout,
		WriteTimeout:                  config.DefaultWriteTimeout,
		PingInterval:                  DefaultPingInterval,
		MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
		MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
	}

	DefaultProviderConfig = config.ProviderConfig{
		Name:      Name,
		WebSocket: DefaultWebSocketConfig,
		Type:      Type,
	}

	// DefaultAPIConfig defines the default API config for KuCoin. This is
	// only utilized on the initial connection to the websocket feed.
	DefaultAPIConfig = config.APIConfig{
		Enabled:    false,
		Timeout:    5 * time.Second, // KuCoin recommends a timeout of 5 seconds.
		Interval:   1 * time.Minute, // This is not used.
		MaxQueries: 1,               // This is not used.
		URL:        URL,
		Name:       Name,
	}

	// DefaultMarketConfig defines the default market config for Kucoin.
	DefaultMarketConfig = types.CurrencyPairsToProviderTickers{
		constants.APE_USDC: {
			OffChainTicker: "APE-USDC",
		},
		constants.APE_USDT: {
			OffChainTicker: "APE-USDT",
		},
		constants.APTOS_USDT: {
			OffChainTicker: "APT-USDT",
		},
		constants.ARBITRUM_USDT: {
			OffChainTicker: "ARB-USDT",
		},
		constants.ATOM_USDC: {
			OffChainTicker: "ATOM-USDC",
		},
		constants.ATOM_USDT: {
			OffChainTicker: "ATOM-USDT",
		},
		constants.AVAX_USDC: {
			OffChainTicker: "AVAX-USDC",
		},
		constants.AVAX_USDT: {
			OffChainTicker: "AVAX-USDT",
		},
		constants.BCH_USDT: {
			OffChainTicker: "BCH-USDT",
		},
		constants.BITCOIN_USDC: {
			OffChainTicker: "BTC-USDC",
		},
		constants.BITCOIN_USDT: {
			OffChainTicker: "BTC-USDT",
		},
		constants.BLUR_USDT: {
			OffChainTicker: "BLUR-USDT",
		},
		constants.CARDANO_USDC: {
			OffChainTicker: "ADA-USDC",
		},
		constants.CARDANO_USDT: {
			OffChainTicker: "ADA-USDT",
		},
		constants.CELESTIA_USDT: {
			OffChainTicker: "TIA-USDT",
		},
		constants.CHAINLINK_USDT: {
			OffChainTicker: "LINK-USDT",
		},
		constants.CURVE_USDT: {
			OffChainTicker: "CRV-USDT",
		},
		constants.DOGE_USDT: {
			OffChainTicker: "DOGE-USDT",
		},
		constants.DYDX_USDT: {
			OffChainTicker: "DYDX-USDT",
		},
		constants.ETC_USDT: {
			OffChainTicker: "ETC-USDT",
		},
		constants.ETHEREUM_BITCOIN: {
			OffChainTicker: "ETH-BTC",
		},
		constants.ETHEREUM_USDC: {
			OffChainTicker: "ETH-USDC",
		},
		constants.ETHEREUM_USDT: {
			OffChainTicker: "ETH-USDT",
		},
		constants.LIDO_USDT: {
			OffChainTicker: "LDO-USDT",
		},
		constants.LITECOIN_USDT: {
			OffChainTicker: "LTC-USDT",
		},
		constants.MAKER_USDT: {
			OffChainTicker: "MKR-USDT",
		},
		constants.NEAR_USDT: {
			OffChainTicker: "NEAR-USDT",
		},
		constants.OPTIMISM_USDT: {
			OffChainTicker: "OP-USDT",
		},
		constants.OSMOSIS_USDT: {
			OffChainTicker: "OSMO-USDT",
		},
		constants.PEPE_USDT: {
			OffChainTicker: "PEPE-USDT",
		},
		constants.POLKADOT_USDT: {
			OffChainTicker: "DOT-USDT",
		},
		constants.POLYGON_USDT: {
			OffChainTicker: "MATIC-USDT",
		},
		constants.RIPPLE_USDT: {
			OffChainTicker: "XRP-USDT",
		},
		constants.SEI_USDT: {
			OffChainTicker: "SEI-USDT",
		},
		constants.SHIBA_USDT: {
			OffChainTicker: "SHIB-USDT",
		},
		constants.SOLANA_USDC: {
			OffChainTicker: "SOL-USDC",
		},
		constants.SOLANA_USDT: {
			OffChainTicker: "SOL-USDT",
		},
		constants.STELLAR_USDT: {
			OffChainTicker: "XLM-USDT",
		},
		constants.SUI_USDT: {
			OffChainTicker: "SUI-USDT",
		},
		constants.TRON_USDT: {
			OffChainTicker: "TRX-USDT",
		},
		constants.UNISWAP_USDT: {
			OffChainTicker: "UNI-USDT",
		},
		constants.USDC_USDT: {
			OffChainTicker: "USDC-USDT",
		},
		constants.WORLD_USDT: {
			OffChainTicker: "WLD-USDT",
		},
	}
)
