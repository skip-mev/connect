package kucoin

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/constants"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

const (
	// Name is the name of the KuCoin provider.
	Name = "kucoin"

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
	DefaultMarketConfig = mmtypes.MarketConfig{
		Name: Name,
		TickerConfigs: map[string]mmtypes.TickerConfig{
			"ATOM/USDC": {
				Ticker:         constants.ATOM_USDC,
				OffChainTicker: "ATOM-USDC",
			},
			"ATOM/USDT": {
				Ticker:         constants.ATOM_USDT,
				OffChainTicker: "ATOM-USDT",
			},
			"AVAX/USDC": {
				Ticker:         constants.AVAX_USDC,
				OffChainTicker: "AVAX-USDC",
			},
			"AVAX/USDT": {
				Ticker:         constants.AVAX_USDT,
				OffChainTicker: "AVAX-USDT",
			},
			"BITCOIN/USDC": {
				Ticker:         constants.BITCOIN_USDC,
				OffChainTicker: "BTC-USDC",
			},
			"BITCOIN/USDT": {
				Ticker:         constants.BITCOIN_USDT,
				OffChainTicker: "BTC-USDT",
			},
			"CELESTIA/USDT": {
				Ticker:         constants.CELESTIA_USDT,
				OffChainTicker: "TIA-USDT",
			},
			"DYDX/USDT": {
				Ticker:         constants.DYDX_USDT,
				OffChainTicker: "DYDX-USDT",
			},
			"ETHEREUM/BITCOIN": {
				Ticker:         constants.ETHEREUM_BITCOIN,
				OffChainTicker: "ETH-BTC",
			},
			"ETHEREUM/USDC": {
				Ticker:         constants.ETHEREUM_USDC,
				OffChainTicker: "ETH-USDC",
			},
			"ETHEREUM/USDT": {
				Ticker:         constants.ETHEREUM_USDT,
				OffChainTicker: "ETH-USDT",
			},
			"OSMOSIS/USDT": {
				Ticker:         constants.OSMOSIS_USDT,
				OffChainTicker: "OSMO-USDT",
			},
			"SOLANA/USDC": {
				Ticker:         constants.SOLANA_USDC,
				OffChainTicker: "SOL-USDC",
			},
			"SOLANA/USDT": {
				Ticker:         constants.SOLANA_USDT,
				OffChainTicker: "SOL-USDT",
			},
			"USDC/USDT": {
				Ticker:         constants.USDC_USDT,
				OffChainTicker: "USDC-USDT",
			},
		},
	}
)
