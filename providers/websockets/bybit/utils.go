package bybit

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// ByBit provides a few different URLs for its Websocket API. The URLs can be found
	// in the documentation here: https://bybit-exchange.github.io/docs/v5/ws/connect
	// The two production URLs are defined in ProductionURL and TestnetURL.

	// Name is the name of the ByBit provider.
	Name = "bybit"

	// URLProd is the public ByBit Websocket URL.
	URLProd = "wss://stream.bybit.com/v5/public/spot"

	// URLTest is the public testnet ByBit Websocket URL.
	URLTest = "wss://stream-testnet.bybit.com/v5/public/spot"

	// DefaultPingInterval is the default ping interval for the ByBit websocket.
	DefaultPingInterval = 15 * time.Second
)

var (
	// DefaultWebSocketConfig is the default configuration for the ByBit Websocket.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Name:                          Name,
		Enabled:                       true,
		MaxBufferSize:                 1000,
		ReconnectionTimeout:           config.DefaultReconnectionTimeout,
		WSS:                           URLProd,
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

	// DefaultMarketConfig is the default market configuration for ByBit.
	DefaultMarketConfig = config.MarketConfig{
		Name: Name,
		CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
			"ATOM/USDT": {
				Ticker:       "ATOMUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USDT"),
			},
			"AVAX/USDC": {
				Ticker:       "AVAXUSDC",
				CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USDC"),
			},
			"AVAX/USDT": {
				Ticker:       "AVAXUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USDT"),
			},
			"BITCOIN/USDC": {
				Ticker:       "BTCUSDC",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDC"),
			},
			"BITCOIN/USDT": {
				Ticker:       "BTCUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
			},
			"DYDX/USDT": {
				Ticker:       "DYDXUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USDT"),
			},
			"ETHEREUM/USDC": {
				Ticker:       "ETHUSDC",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDC"),
			},
			"ETHEREUM/USDT": {
				Ticker:       "ETHUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
			},
			"SOLANA/USDC": {
				Ticker:       "SOLUSDC",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USDC"),
			},
			"SOLANA/USDT": {
				Ticker:       "SOLUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USDT"),
			},
			"USDC/USDT": {
				Ticker:       "USDCUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("USDC", "USDT"),
			},
		},
	}
)
