package constants

import (
	"github.com/skip-mev/slinky/oracle/config"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/apis/defi/raydium"
	"github.com/skip-mev/slinky/providers/apis/defi/uniswapv3"
	"github.com/skip-mev/slinky/providers/apis/marketmaps/dydx"
	"github.com/skip-mev/slinky/providers/apis/marketmaps/marketmap"
	binanceapi "github.com/skip-mev/slinky/providers/apis/prices/binance"
	coinbaseapi "github.com/skip-mev/slinky/providers/apis/prices/coinbase"
	bitstampapi "github.com/skip-mev/slinky/providers/apis/prices/bitstamp"
	"github.com/skip-mev/slinky/providers/apis/prices/coingecko"
	"github.com/skip-mev/slinky/providers/apis/prices/coinmarketcap"
	krakenapi "github.com/skip-mev/slinky/providers/apis/prices/kraken"
	"github.com/skip-mev/slinky/providers/volatile"
	binancews "github.com/skip-mev/slinky/providers/websockets/prices/binance"
	"github.com/skip-mev/slinky/providers/websockets/prices/bitfinex"
	"github.com/skip-mev/slinky/providers/websockets/prices/bitstamp"
	"github.com/skip-mev/slinky/providers/websockets/prices/bybit"
	"github.com/skip-mev/slinky/providers/websockets/prices/coinbase"
	"github.com/skip-mev/slinky/providers/websockets/prices/cryptodotcom"
	"github.com/skip-mev/slinky/providers/websockets/prices/gate"
	"github.com/skip-mev/slinky/providers/websockets/prices/huobi"
	"github.com/skip-mev/slinky/providers/websockets/prices/kraken"
	"github.com/skip-mev/slinky/providers/websockets/prices/kucoin"
	"github.com/skip-mev/slinky/providers/websockets/prices/mexc"
	"github.com/skip-mev/slinky/providers/websockets/prices/okx"
	mmtypes "github.com/skip-mev/slinky/service/clients/marketmap/types"
)

var (
	Providers = []config.ProviderConfig{
		// DEFI providers
		{
			Name: raydium.Name,
			API:  raydium.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: uniswapv3.ProviderNames[constants.ETHEREUM],
			API:  uniswapv3.DefaultETHAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: uniswapv3.ProviderNames[constants.BASE],
			API:  uniswapv3.DefaultBaseAPIConfig,
			Type: types.ConfigType,
		},

		// Exchange API providers
		{
			Name: binanceapi.Name,
			API:  binanceapi.DefaultNonUSAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: bitstampapi.Name,
			API:  bitstampapi.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: coinbaseapi.Name,
			API:  coinbaseapi.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: coingecko.Name,
			API:  coingecko.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: coinmarketcap.Name,
			API:  coinmarketcap.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: krakenapi.Name,
			API:  krakenapi.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		{
			Name: volatile.Name,
			API:  volatile.DefaultAPIConfig,
			Type: types.ConfigType,
		},
		// Exchange WebSocket providers
		{
			Name:      binancews.Name,
			WebSocket: binancews.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      bitfinex.Name,
			WebSocket: bitfinex.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      bitstamp.Name,
			WebSocket: bitstamp.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      bybit.Name,
			WebSocket: bybit.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      coinbase.Name,
			WebSocket: coinbase.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      cryptodotcom.Name,
			WebSocket: cryptodotcom.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      gate.Name,
			WebSocket: gate.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      huobi.Name,
			WebSocket: huobi.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      kraken.Name,
			WebSocket: kraken.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      kucoin.Name,
			WebSocket: kucoin.DefaultWebSocketConfig,
			API:       kucoin.DefaultAPIConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      mexc.Name,
			WebSocket: mexc.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},
		{
			Name:      okx.Name,
			WebSocket: okx.DefaultWebSocketConfig,
			Type:      types.ConfigType,
		},

		// MarketMap provider
		{
			Name: marketmap.Name,
			API:  marketmap.DefaultAPIConfig,
			Type: mmtypes.ConfigType,
		},
	}

	AlternativeMarketMapProviders = []config.ProviderConfig{
		{
			Name: dydx.Name,
			API:  dydx.DefaultAPIConfig,
			Type: mmtypes.ConfigType,
		},
		{
			Name: dydx.ResearchAPIHandlerName,
			API:  dydx.DefaultResearchAPIConfig,
			Type: mmtypes.ConfigType,
		},
		{
			Name: dydx.ResearchCMCAPIHandlerName,
			API:  dydx.DefaultResearchCMCAPIConfig,
			Type: mmtypes.ConfigType,
		},
	}

	MarketMapProviderNames = map[string]struct{}{
		dydx.Name:                      {},
		dydx.ResearchAPIHandlerName:    {},
		dydx.ResearchCMCAPIHandlerName: {},
		marketmap.Name:                 {},
	}
)
