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
)

var (
	// DefaultWebSocketConfig is the default configuration for the ByBit Websocket.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Name:                Name,
		Enabled:             true,
		MaxBufferSize:       1000,
		ReconnectionTimeout: config.DefaultReconnectionTimeout,
		WSS:                 URLProd,
		ReadBufferSize:      config.DefaultReadBufferSize,
		WriteBufferSize:     config.DefaultWriteBufferSize,
		HandshakeTimeout:    config.DefaultHandshakeTimeout,
		EnableCompression:   config.DefaultEnableCompression,
		ReadTimeout:         config.DefaultReadTimeout,
		WriteTimeout:        config.DefaultWriteTimeout,
		PingInterval:        15 * time.Second,
	}

	// DefaultMarketConfig is the default market configuration for ByBit.
	DefaultMarketConfig = config.MarketConfig{
		Name: Name,
		CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
			"BITCOIN/USD": {
				Ticker:       "BTCUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			"ETHEREUM/USD": {
				Ticker:       "ETHUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			"ATOM/USD": {
				Ticker:       "ATOMUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USD"),
			},
			"SOLANA/USD": {
				Ticker:       "SOLUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD"),
			},
			"AVAX/USD": {
				Ticker:       "AVAXUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USD"),
			},
			"DYDX/USD": {
				Ticker:       "DYDXUSDT",
				CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USD"),
			},
		},
	}
)
