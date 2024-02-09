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
			"ATOM/USDT": {
				Ticker:       "ATOM_USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("ATOM", "USDT"),
			},
			"AVAX/USDT": {
				Ticker:       "AVAX_USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("AVAX", "USDT"),
			},
			"BITCOIN/USDT": {
				Ticker:       "BTC_USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("BITCOIN", "USDT"),
			},
			"CELESTIA/USDT": {
				Ticker:       "TIA_USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("CELESTIA", "USDT"),
			},
			"DYDX/USDT": {
				Ticker:       "DYDX_USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("DYDX", "USDT"),
			},
			"ETHEREUM/BITCOIN": {
				Ticker:       "ETH_BTC",
				CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "BITCOIN"),
			},
			"ETHEREUM/USDT": {
				Ticker:       "ETH_USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("ETHEREUM", "USDT"),
			},
			"SOLANA/USDC": {
				Ticker:       "SOL_USDC",
				CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USDC"),
			},
			"SOLANA/USDT": {
				Ticker:       "SOL_USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("SOLANA", "USDT"),
			},
			"USDC/USDT": {
				Ticker:       "USDC_USDT",
				CurrencyPair: slinkytypes.NewCurrencyPair("USDC", "USDT"),
			},
		},
	}
)
