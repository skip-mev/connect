package constants

import (
	"github.com/skip-mev/slinky/oracle/config"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/apis/dydx"
	"github.com/skip-mev/slinky/providers/apis/marketmap"
	"github.com/skip-mev/slinky/providers/websockets/kucoin"
	mmtypes "github.com/skip-mev/slinky/service/clients/marketmap/types"
)

var (
	Providers = []config.ProviderConfig{
		// // DEFI providers
		// {
		// 	Name: raydium.Name,
		// 	API:  raydium.DefaultAPIConfig,
		// 	Type: types.ConfigType,
		// },
		// {
		// 	Name: uniswapv3.ProviderNames[constants.ETHEREUM],
		// 	API:  uniswapv3.DefaultETHAPIConfig,
		// 	Type: types.ConfigType,
		// },

		// // Exchange providers
		// {
		// 	Name: coinbaseapi.Name,
		// 	API:  coinbaseapi.DefaultAPIConfig,
		// 	Type: types.ConfigType,
		// },
		// {
		// 	Name: binanceapi.Name,
		// 	API:  binanceapi.DefaultNonUSAPIConfig,
		// 	Type: types.ConfigType,
		// },
		// {
		// 	Name: krakenapi.Name,
		// 	API:  krakenapi.DefaultAPIConfig,
		// 	Type: types.ConfigType,
		// },
		// {
		// 	Name: volatile.Name,
		// 	API:  volatile.DefaultAPIConfig,
		// 	Type: types.ConfigType,
		// },
		// {
		// 	Name:      bitfinex.Name,
		// 	WebSocket: bitfinex.DefaultWebSocketConfig,
		// 	Type:      types.ConfigType,
		// },
		// {
		// 	Name:      bitstamp.Name,
		// 	WebSocket: bitstamp.DefaultWebSocketConfig,
		// 	Type:      types.ConfigType,
		// },
		// {
		// 	Name:      bybit.Name,
		// 	WebSocket: bybit.DefaultWebSocketConfig,
		// 	Type:      types.ConfigType,
		// },
		// {
		// 	Name:      coinbase.Name,
		// 	WebSocket: coinbase.DefaultWebSocketConfig,
		// 	Type:      types.ConfigType,
		// },
		// {
		// 	Name:      cryptodotcom.Name,
		// 	WebSocket: cryptodotcom.DefaultWebSocketConfig,
		// 	Type:      types.ConfigType,
		// },
		// {
		// 	Name:      gate.Name,
		// 	WebSocket: gate.DefaultWebSocketConfig,
		// 	Type:      types.ConfigType,
		// },
		// {
		// 	Name:      huobi.Name,
		// 	WebSocket: huobi.DefaultWebSocketConfig,
		// 	Type:      types.ConfigType,
		// },
		// {
		// 	Name:      kraken.Name,
		// 	WebSocket: kraken.DefaultWebSocketConfig,
		// 	Type:      types.ConfigType,
		// },
		{
			Name:      kucoin.Name,
			WebSocket: kucoin.DefaultWebSocketConfig,
			API:       kucoin.DefaultAPIConfig,
			Type:      types.ConfigType,
		},
		// {
		// 	Name:      mexc.Name,
		// 	WebSocket: mexc.DefaultWebSocketConfig,
		// 	Type:      types.ConfigType,
		// },
		// {
		// 	Name:      okx.Name,
		// 	WebSocket: okx.DefaultWebSocketConfig,
		// 	Type:      types.ConfigType,
		// },

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
	}

	MarketMapProviderNames = map[string]struct{}{
		dydx.Name:                   {},
		dydx.ResearchAPIHandlerName: {},
		marketmap.Name:              {},
	}
)
