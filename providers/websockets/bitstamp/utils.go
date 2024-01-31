package bitstamp

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// Name is the name of the bitstamp provider.
	Name = "bitstamp"

	// WSS is the bitstamp websocket address.
	WSS = "wss://ws.bitstamp.net"
)

var (
	// DefaultWebSocketConfig returns the default websocket config for bitstamp.
	DefaultWebSocketConfig = config.WebSocketConfig{
		Enabled:             true,
		Name:                Name,
		MaxBufferSize:       config.DefaultMaxBufferSize,
		ReconnectionTimeout: config.DefaultReconnectionTimeout,
		WSS:                 WSS,
		ReadBufferSize:      config.DefaultReadBufferSize,
		WriteBufferSize:     config.DefaultWriteBufferSize,
		HandshakeTimeout:    config.DefaultHandshakeTimeout,
		EnableCompression:   config.DefaultEnableCompression,
		WriteTimeout:        config.DefaultWriteTimeout,
		ReadTimeout:         config.DefaultReadTimeout,
		PingInterval:        10 * time.Second,
		MaxReadErrorCount:   config.DefaultMaxReadErrorCount,
	}

	// DefaultMarketConfig returns the default market config for bitstamp.
	DefaultMarketConfig = config.MarketConfig{
		Name: Name,
		CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
			"BITCOIN/USD": {
				Ticker:       "btcusd",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USD"),
			},
			"ETHEREUM/USD": {
				Ticker:       "ethusd",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
			},
			"SOLANA/USD": {
				Ticker:       "solusd",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USD"),
			},
			"AVAX/USD": {
				Ticker:       "avaxusd",
				CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USD"),
			},
			"DYDX/USD": {
				Ticker:       "dydxusd",
				CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USD"),
			},
			"ETHEREUM/BITCOIN": {
				Ticker:       "ethbtc",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
			},
		},
	}
)
