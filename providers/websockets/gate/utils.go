package gate

import (
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// Name is the name of the Gate.io provider.
	Name = "gate.io"
	// URL is the base url of for the Gate.io websocket API.
	URL = "wss://api.gateio.ws/ws/v4/"
)

// DefaultWebSocketConfig is the default configuration for the Gate.io Websocket.
var (
	DefaultWebSocketConfig = config.WebSocketConfig{
		Name:                          Name,
		Enabled:                       true,
		MaxBufferSize:                 1000,
		ReconnectionTimeout:           10 * time.Second,
		WSS:                           URL,
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

	// DefaultMarketConfig is the default market configuration for Gate.io.
	DefaultMarketConfig = config.MarketConfig{
		Name: Name,
		CurrencyPairToMarketConfigs: map[string]config.CurrencyPairMarketConfig{
			"BITCOIN/USDT": {
				Ticker:       "BTC_USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("BITCOIN", "USDT"),
			},
			"ETHEREUM/USDT": {
				Ticker:       "ETH_USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "USDT"),
			},
			"ATOM/USDT": {
				Ticker:       "ATOM_USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("ATOM", "USDT"),
			},
			"SOLANA/USDT": {
				Ticker:       "SOL_USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USDT"),
			},
			"SOLANA/USDC": {
				Ticker:       "SOL_USDC",
				CurrencyPair: oracletypes.NewCurrencyPair("SOLANA", "USDC"),
			},
			"CELESTIA/USDT": {
				Ticker:       "TIA_USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("CELESTIA", "USDT"),
			},
			"AVAX/USDT": {
				Ticker:       "AVAX_USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("AVAX", "USDT"),
			},
			"DYDX/USDT": {
				Ticker:       "DYDX_USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("DYDX", "USDT"),
			},
			"ETHEREUM/BITCOIN": {
				Ticker:       "ETH_BTC",
				CurrencyPair: oracletypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
			},
			"USDC/USDT": {
				Ticker:       "USDC_USDT",
				CurrencyPair: oracletypes.NewCurrencyPair("USDC", "USDT"),
			},
			"USDT/USD": {
				Ticker:       "USDT_USD",
				CurrencyPair: oracletypes.NewCurrencyPair("USDT", "USD"),
			},
		},
	}
)
