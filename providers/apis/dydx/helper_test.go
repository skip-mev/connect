package dydx_test

import (
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/providers/apis/kraken"
	"github.com/skip-mev/connect/v2/providers/websockets/binance"
	"github.com/skip-mev/connect/v2/providers/websockets/bybit"
	coinbasews "github.com/skip-mev/connect/v2/providers/websockets/coinbase"
	"github.com/skip-mev/connect/v2/providers/websockets/huobi"
	"github.com/skip-mev/connect/v2/providers/websockets/kucoin"
	"github.com/skip-mev/connect/v2/providers/websockets/mexc"
	"github.com/skip-mev/connect/v2/providers/websockets/okx"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

const dYdXResponseValid = `
{
	"market_params": [
	  {
		"id": 0,
		"pair": "BTC-USD",
		"exponent": -5,
		"min_exchanges": 3,
		"min_price_change_ppm": 1000,
		"exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"BTCUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Bybit\",\"ticker\":\"BTCUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"BTC-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"btcusdt\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"XXBTZUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"BTC-USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Mexc\",\"ticker\":\"BTC_USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"BTC-USDT\",\"adjustByMarket\":\"USDT-USD\"}]}"
	  },
	  {
		"id": 1,
		"pair": "ETH-USD",
		"exponent": -6,
		"min_exchanges": 3,
		"min_price_change_ppm": 1000,
		"exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"ETHUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Bybit\",\"ticker\":\"ETHUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"ETH-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"ethusdt\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"XETHZUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"ETH-USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Mexc\",\"ticker\":\"ETH_USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"ETH-USDT\",\"adjustByMarket\":\"USDT-USD\"}]}"
	  },
	  {
		"id": 1000000,
		"pair": "USDT-USD",
		"exponent": -9,
		"min_exchanges": 3,
		"min_price_change_ppm": 1000,
		"exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"USDCUSDT\",\"invert\":true},{\"exchangeName\":\"Bybit\",\"ticker\":\"USDCUSDT\",\"invert\":true},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"USDT-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"ethusdt\",\"adjustByMarket\":\"ETH-USD\",\"invert\":true},{\"exchangeName\":\"Kraken\",\"ticker\":\"USDTZUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"BTC-USDT\",\"adjustByMarket\":\"BTC-USD\",\"invert\":true},{\"exchangeName\":\"Okx\",\"ticker\":\"USDC-USDT\",\"invert\":true}]}"
	  }
	],
	"pagination": {
	  "next_key": null,
	  "total": "58"
	}
}
`

const dYdXResponseInvalid = `
{
	"market_params": [
	  {
		"id": 0,
		"pair": "BTC-USD",
		"exponent": -5,
		"min_exchanges": 3,
		"min_price_change_ppm": 1000,
		"exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"BTCUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Bybit\",\"ticker\":\"BTCUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"BTC-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"btcusdt\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"XXBTZUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"BTC-USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Mexc\",\"ticker\":\"BTC_USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"BTC-USDT\",\"adjustByMarket\":\"USDT-USD\"}]}"
	  },
	  {
		"id": 1,
		"pair": "ETH-USD",
		"exponent": -6,
		"min_exchanges": 3,
		"min_price_change_ppm": 1000,
		"exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"ETHUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Bybit\",\"ticker\":\"ETHUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"ETH-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"ethusdt\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"XETHZUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"ETH-USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Mexc\",\"ticker\":\"ETH_USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"ETH-USDT\",\"adjustByMarket\":\"USDT-USD\"}]}"
	  },
	],
	"pagination": {
	  "next_key": null,
	  "total": "58"
	}
}
`

var usdtusd = connecttypes.NewCurrencyPair("USDT", "USD")

var convertedResponse = mmtypes.MarketMapResponse{
	MarketMap: mmtypes.MarketMap{
		Markets: map[string]mmtypes.Market{
			"BTC/USD": {
				Ticker: mmtypes.NewTicker("BTC", "USD", 5, 3, true),
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:            binance.Name,
						OffChainTicker:  "BTCUSDT",
						NormalizeByPair: &usdtusd,
					},
					{
						Name:            bybit.Name,
						OffChainTicker:  "BTCUSDT",
						NormalizeByPair: &usdtusd,
					},
					{
						Name:           coinbasews.Name,
						OffChainTicker: "BTC-USD",
					},
					{
						Name:            huobi.Name,
						OffChainTicker:  "btcusdt",
						NormalizeByPair: &usdtusd,
					},
					{
						Name:           kraken.Name,
						OffChainTicker: "XXBTZUSD",
					},
					{
						Name:            kucoin.Name,
						OffChainTicker:  "BTC-USDT",
						NormalizeByPair: &usdtusd,
					},
					{
						Name:            mexc.Name,
						OffChainTicker:  "BTCUSDT",
						NormalizeByPair: &usdtusd,
					},
					{
						Name:            okx.Name,
						OffChainTicker:  "BTC-USDT",
						NormalizeByPair: &usdtusd,
					},
				},
			},
			"ETH/USD": {
				Ticker: mmtypes.NewTicker("ETH", "USD", 6, 3, true),
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:            binance.Name,
						OffChainTicker:  "ETHUSDT",
						NormalizeByPair: &usdtusd,
					},
					{
						Name:            bybit.Name,
						OffChainTicker:  "ETHUSDT",
						NormalizeByPair: &usdtusd,
					},
					{
						Name:           coinbasews.Name,
						OffChainTicker: "ETH-USD",
					},
					{
						Name:            huobi.Name,
						OffChainTicker:  "ethusdt",
						NormalizeByPair: &usdtusd,
					},
					{
						Name:           kraken.Name,
						OffChainTicker: "XETHZUSD",
					},
					{
						Name:            kucoin.Name,
						OffChainTicker:  "ETH-USDT",
						NormalizeByPair: &usdtusd,
					},
					{
						Name:            mexc.Name,
						OffChainTicker:  "ETHUSDT",
						NormalizeByPair: &usdtusd,
					},
					{
						Name:            okx.Name,
						OffChainTicker:  "ETH-USDT",
						NormalizeByPair: &usdtusd,
					},
				},
			},
			"USDT/USD": {
				Ticker: mmtypes.NewTicker("USDT", "USD", 9, 3, true),
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:           binance.Name,
						OffChainTicker: "USDCUSDT",
						Invert:         true,
					},
					{
						Name:           bybit.Name,
						OffChainTicker: "USDCUSDT",
						Invert:         true,
					},
					{
						Name:           coinbasews.Name,
						OffChainTicker: "USDT-USD",
					},
					{
						Name:            huobi.Name,
						OffChainTicker:  "ethusdt",
						NormalizeByPair: &connecttypes.CurrencyPair{Base: "ETH", Quote: "USD"},
						Invert:          true,
					},
					{
						Name:           kraken.Name,
						OffChainTicker: "USDTZUSD",
					},
					{
						Name:            kucoin.Name,
						OffChainTicker:  "BTC-USDT",
						NormalizeByPair: &connecttypes.CurrencyPair{Base: "BTC", Quote: "USD"},
						Invert:          true,
					},
					{
						Name:           okx.Name,
						OffChainTicker: "USDC-USDT",
						Invert:         true,
					},
				},
			},
		},
	},
}
