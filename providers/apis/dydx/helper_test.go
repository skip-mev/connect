package dydx_test

import (
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/providers/apis/binance"
	coinbaseapi "github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/apis/kraken"
	"github.com/skip-mev/slinky/providers/websockets/bybit"
	coinbasews "github.com/skip-mev/slinky/providers/websockets/coinbase"
	"github.com/skip-mev/slinky/providers/websockets/huobi"
	"github.com/skip-mev/slinky/providers/websockets/kucoin"
	"github.com/skip-mev/slinky/providers/websockets/mexc"
	"github.com/skip-mev/slinky/providers/websockets/okx"
	mmtypes "github.com/skip-mev/slinky/x/mm2/types"
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

var convertedResponse = mmtypes.MarketMapResponse{
	MarketMap: mmtypes.MarketMap{
		Markets: map[string]mmtypes.Market{
			"BTC/USD": {
				Ticker: mmtypes.Ticker{
					CurrencyPair:     slinkytypes.NewCurrencyPair("BTC", "USD"),
					Decimals:         5,
					MinProviderCount: 3,
				},
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:            binance.Name,
						OffChainTicker:  "BTCUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"},
						Invert:          false,
					},
					{
						Name:            bybit.Name,
						OffChainTicker:  "BTCUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"},
						Invert:          false,
					},
					{
						Name:           coinbaseapi.Name,
						OffChainTicker: "BTC-USD",
					},
					{
						Name:           coinbasews.Name,
						OffChainTicker: "BTC-USD",
					},
					{
						Name:            huobi.Name,
						OffChainTicker:  "btcusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"},
						Invert:          false,
					},
					{
						Name:           kraken.Name,
						OffChainTicker: "XXBTZUSD",
					},
					{
						Name:            kucoin.Name,
						OffChainTicker:  "BTC-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"},
						Invert:          false,
					},
					{
						Name:            mexc.Name,
						OffChainTicker:  "BTCUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"},
						Invert:          false,
					},
					{
						Name:            okx.Name,
						OffChainTicker:  "BTC-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"},
						Invert:          false,
					},
				},
			},
			"ETH/USD": {
				Ticker: mmtypes.Ticker{
					CurrencyPair:     slinkytypes.NewCurrencyPair("ETH", "USD"),
					Decimals:         6,
					MinProviderCount: 3,
				},
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:            binance.Name,
						OffChainTicker:  "ETHUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"},
						Invert:          false,
					},
					{
						Name:            bybit.Name,
						OffChainTicker:  "ETHUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"},
						Invert:          false,
					},
					{
						Name:           coinbaseapi.Name,
						OffChainTicker: "ETH-USD",
					},
					{
						Name:           coinbasews.Name,
						OffChainTicker: "ETH-USD",
					},
					{
						Name:            huobi.Name,
						OffChainTicker:  "ethusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"},
						Invert:          false,
					},
					{
						Name:           kraken.Name,
						OffChainTicker: "XETHZUSD",
					},
					{
						Name:            kucoin.Name,
						OffChainTicker:  "ETH-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"},
						Invert:          false,
					},
					{
						Name:            mexc.Name,
						OffChainTicker:  "ETHUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"},
						Invert:          false,
					},
					{
						Name:            okx.Name,
						OffChainTicker:  "ETH-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"},
						Invert:          false,
					},
				},
			},
			"USDT/USD": {
				Ticker: mmtypes.Ticker{
					CurrencyPair:     slinkytypes.NewCurrencyPair("USDT", "USD"),
					Decimals:         9,
					MinProviderCount: 3,
				},
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:           binance.Name,
						OffChainTicker: "USDCUSDT",
					},
					{
						Name:           bybit.Name,
						OffChainTicker: "USDCUSDT",
					},
					{
						Name:           coinbaseapi.Name,
						OffChainTicker: "USDT-USD",
					},
					{
						Name:           coinbasews.Name,
						OffChainTicker: "USDT-USD",
					},
					{
						Name:            huobi.Name,
						OffChainTicker:  "ethusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "ETH", Quote: "USD"},
						Invert:          true,
					},
					{
						Name:           kraken.Name,
						OffChainTicker: "USDTZUSD",
					},
					{
						Name:            kucoin.Name,
						OffChainTicker:  "BTC-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "BTC", Quote: "USD"},
						Invert:          true,
					},
					{
						Name:           okx.Name,
						OffChainTicker: "USDC-USDT",
					},
				},
			},
		},
	},
}
