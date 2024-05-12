package providers

import (
	"github.com/skip-mev/slinky/oracle/config"
	binanceapi "github.com/skip-mev/slinky/providers/apis/binance"
	coinbaseapi "github.com/skip-mev/slinky/providers/apis/coinbase"
	coingeckoapi "github.com/skip-mev/slinky/providers/apis/coingecko"
	"github.com/skip-mev/slinky/providers/apis/defi/raydium"
	"github.com/skip-mev/slinky/providers/apis/defi/uniswapv3"
	geckoterminalapi "github.com/skip-mev/slinky/providers/apis/geckoterminal"
	krakenapi "github.com/skip-mev/slinky/providers/apis/kraken"
	"github.com/skip-mev/slinky/providers/apis/marketmap"
	"github.com/skip-mev/slinky/providers/volatile"
	bitfinexws "github.com/skip-mev/slinky/providers/websockets/bitfinex"
	bitstampws "github.com/skip-mev/slinky/providers/websockets/bitstamp"
	bybitws "github.com/skip-mev/slinky/providers/websockets/bybit"
	coinbasews "github.com/skip-mev/slinky/providers/websockets/coinbase"
	cryptodotcomws "github.com/skip-mev/slinky/providers/websockets/cryptodotcom"
	gatews "github.com/skip-mev/slinky/providers/websockets/gate"
	huobiws "github.com/skip-mev/slinky/providers/websockets/huobi"
	krakenws "github.com/skip-mev/slinky/providers/websockets/kraken"
	kucoinws "github.com/skip-mev/slinky/providers/websockets/kucoin"
	mexcws "github.com/skip-mev/slinky/providers/websockets/mexc"
	okxws "github.com/skip-mev/slinky/providers/websockets/okx"
)

var ProviderDefaults = []config.ProviderConfig{
	// API Providers
	binanceapi.DefaultProviderConfig,
	coinbaseapi.DefaultProviderConfig,
	coingeckoapi.DefaultProviderConfig,
	geckoterminalapi.DefaultProviderConfig,
	krakenapi.DefaultProviderConfig,
	// Websocket Providers
	bitfinexws.DefaultProviderConfig,
	bitstampws.DefaultProviderConfig,
	bybitws.DefaultProviderConfig,
	coinbasews.DefaultProviderConfig,
	cryptodotcomws.DefaultProviderConfig,
	gatews.DefaultProviderConfig,
	huobiws.DefaultProviderConfig,
	krakenws.DefaultProviderConfig,
	kucoinws.DefaultProviderConfig,
	mexcws.DefaultProviderConfig,
	okxws.DefaultProviderConfig,
	// Defi Providers
	raydium.DefaultProviderConfig,
	uniswapv3.DefaultETHProviderConfig,
	// MM Provider
	marketmap.DefaultProviderConfig,
	// Test Providers
	volatile.DefaultProviderConfig,
}
