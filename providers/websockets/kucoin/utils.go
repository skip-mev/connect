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
	DefaultMarketConfig = types.TickerToProviderConfig{
		constants.APE_USDC: {
			Name:           Name,
			OffChainTicker: "APE-USDC",
		},
		constants.APE_USDT: {
			Name:           Name,
			OffChainTicker: "APE-USDT",
		},
		constants.APTOS_USDT: {
			Name:           Name,
			OffChainTicker: "APT-USDT",
		},
		constants.ARBITRUM_USDT: {
			Name:           Name,
			OffChainTicker: "ARB-USDT",
		},
		constants.ATOM_USDC: {
			Name:           Name,
			OffChainTicker: "ATOM-USDC",
		},
		constants.ATOM_USDT: {
			Name:           Name,
			OffChainTicker: "ATOM-USDT",
		},
		constants.AVAX_USDC: {
			Name:           Name,
			OffChainTicker: "AVAX-USDC",
		},
		constants.AVAX_USDT: {
			Name:           Name,
			OffChainTicker: "AVAX-USDT",
		},
		constants.BCH_USDT: {
			Name:           Name,
			OffChainTicker: "BCH-USDT",
		},
		constants.BITCOIN_USDC: {
			Name:           Name,
			OffChainTicker: "BTC-USDC",
		},
		constants.BITCOIN_USDT: {
			Name:           Name,
			OffChainTicker: "BTC-USDT",
		},
		constants.BLUR_USDT: {
			Name:           Name,
			OffChainTicker: "BLUR-USDT",
		},
		constants.CARDANO_USDC: {
			Name:           Name,
			OffChainTicker: "ADA-USDC",
		},
		constants.CARDANO_USDT: {
			Name:           Name,
			OffChainTicker: "ADA-USDT",
		},
		constants.CELESTIA_USDT: {
			Name:           Name,
			OffChainTicker: "TIA-USDT",
		},
		constants.CHAINLINK_USDT: {
			Name:           Name,
			OffChainTicker: "LINK-USDT",
		},
		constants.CURVE_USDT: {
			Name:           Name,
			OffChainTicker: "CRV-USDT",
		},
		constants.DOGE_USDT: {
			Name:           Name,
			OffChainTicker: "DOGE-USDT",
		},
		constants.DYDX_USDT: {
			Name:           Name,
			OffChainTicker: "DYDX-USDT",
		},
		constants.ETC_USDT: {
			Name:           Name,
			OffChainTicker: "ETC-USDT",
		},
		constants.ETHEREUM_BITCOIN: {
			Name:           Name,
			OffChainTicker: "ETH-BTC",
		},
		constants.ETHEREUM_USDC: {
			Name:           Name,
			OffChainTicker: "ETH-USDC",
		},
		constants.ETHEREUM_USDT: {
			Name:           Name,
			OffChainTicker: "ETH-USDT",
		},
		constants.LIDO_USDT: {
			Name:           Name,
			OffChainTicker: "LDO-USDT",
		},
		constants.LITECOIN_USDT: {
			Name:           Name,
			OffChainTicker: "LTC-USDT",
		},
		constants.MAKER_USDT: {
			Name:           Name,
			OffChainTicker: "MKR-USDT",
		},
		constants.NEAR_USDT: {
			Name:           Name,
			OffChainTicker: "NEAR-USDT",
		},
		constants.OPTIMISM_USDT: {
			Name:           Name,
			OffChainTicker: "OP-USDT",
		},
		constants.OSMOSIS_USDT: {
			Name:           Name,
			OffChainTicker: "OSMO-USDT",
		},
		constants.PEPE_USDT: {
			Name:           Name,
			OffChainTicker: "PEPE-USDT",
		},
		constants.POLKADOT_USDT: {
			Name:           Name,
			OffChainTicker: "DOT-USDT",
		},
		constants.POLYGON_USDT: {
			Name:           Name,
			OffChainTicker: "MATIC-USDT",
		},
		constants.RIPPLE_USDT: {
			Name:           Name,
			OffChainTicker: "XRP-USDT",
		},
		constants.SEI_USDT: {
			Name:           Name,
			OffChainTicker: "SEI-USDT",
		},
		constants.SHIBA_USDT: {
			Name:           Name,
			OffChainTicker: "SHIB-USDT",
		},
		constants.SOLANA_USDC: {
			Name:           Name,
			OffChainTicker: "SOL-USDC",
		},
		constants.SOLANA_USDT: {
			Name:           Name,
			OffChainTicker: "SOL-USDT",
		},
		constants.STELLAR_USDT: {
			Name:           Name,
			OffChainTicker: "XLM-USDT",
		},
		constants.SUI_USDT: {
			Name:           Name,
			OffChainTicker: "SUI-USDT",
		},
		constants.TRON_USDT: {
			Name:           Name,
			OffChainTicker: "TRX-USDT",
		},
		constants.UNISWAP_USDT: {
			Name:           Name,
			OffChainTicker: "UNI-USDT",
		},
		constants.USDC_USDT: {
			Name:           Name,
			OffChainTicker: "USDC-USDT",
		},
		constants.WORLD_USDT: {
			Name:           Name,
			OffChainTicker: "WLD-USDT",
		},
	}
)
