package bitfinex

import (
	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// Name is the name of the BitFinex provider.
	Name = "bitfinex"

	// URLProd is the public BitFinex Websocket URL.
	URLProd = "wss://api-pub.bitfinex.com/ws/2"

	MaxSubscriptionsPerConnection = 30
)

var (
	// DefaultWebSocketConfig is the default configuration for the BitFinex Websocket.
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
		PingInterval:                  config.DefaultPingInterval,
		MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
		MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
	}

	// DefaultMarketConfig is the default market configuration for BitFinex.
	DefaultMarketConfig = config.MarketConfig{
		Name: Name,
		CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
			"BITCOIN/USD": {
				Ticker:       "BTCUSD",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			"CELESTIA/USD": {
				Ticker:       "TIAUSD",
				CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USD"),
			},
			"ETHEREUM/BITCOIN": {
				Ticker:       "ETHBTC",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
			},
			"ETHEREUM/USD": {
				Ticker:       "ETHUSD",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			"SOLANA/USD": {
				Ticker:       "SOLUSD",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD"),
			},
		},
	}
)
