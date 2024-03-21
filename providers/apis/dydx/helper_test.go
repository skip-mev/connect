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
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
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

var convertedResponse = mmtypes.GetMarketMapResponse{
	MarketMap: mmtypes.MarketMap{
		Tickers: map[string]mmtypes.Ticker{
			"BTC/USD": {
				CurrencyPair:     slinkytypes.NewCurrencyPair("BTC", "USD"),
				Decimals:         5,
				MinProviderCount: 3,
			},
			"ETH/USD": {
				CurrencyPair:     slinkytypes.NewCurrencyPair("ETH", "USD"),
				Decimals:         6,
				MinProviderCount: 3,
			},
			"USDT/USD": {
				CurrencyPair:     slinkytypes.NewCurrencyPair("USDT", "USD"),
				Decimals:         9,
				MinProviderCount: 3,
			},
		},
		Providers: map[string]mmtypes.Providers{
			"BTC/USD": {
				Providers: []mmtypes.ProviderConfig{
					{
						Name:           binance.Name,
						OffChainTicker: "BTCUSDT",
					},
					{
						Name:           bybit.Name,
						OffChainTicker: "BTCUSDT",
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
						Name:           huobi.Name,
						OffChainTicker: "btcusdt",
					},
					{
						Name:           kraken.Name,
						OffChainTicker: "XXBTZUSD",
					},
					{
						Name:           kucoin.Name,
						OffChainTicker: "BTC-USDT",
					},
					{
						Name:           mexc.Name,
						OffChainTicker: "BTCUSDT",
					},
					{
						Name:           okx.Name,
						OffChainTicker: "BTC-USDT",
					},
				},
			},
			"ETH/USD": {
				Providers: []mmtypes.ProviderConfig{
					{
						Name:           binance.Name,
						OffChainTicker: "ETHUSDT",
					},
					{
						Name:           bybit.Name,
						OffChainTicker: "ETHUSDT",
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
						Name:           huobi.Name,
						OffChainTicker: "ethusdt",
					},
					{
						Name:           kraken.Name,
						OffChainTicker: "XETHZUSD",
					},
					{
						Name:           kucoin.Name,
						OffChainTicker: "ETH-USDT",
					},
					{
						Name:           mexc.Name,
						OffChainTicker: "ETHUSDT",
					},
					{
						Name:           okx.Name,
						OffChainTicker: "ETH-USDT",
					},
				},
			},
			"USDT/USD": {
				Providers: []mmtypes.ProviderConfig{
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
						Name:           kraken.Name,
						OffChainTicker: "USDTZUSD",
					},
					{
						Name:           okx.Name,
						OffChainTicker: "USDC-USDT",
					},
				},
			},
		},
		Paths: map[string]mmtypes.Paths{
			"BTC/USD": {
				Paths: []mmtypes.Path{
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     binance.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
								Invert:       false,
							},
							{
								Provider:     mmtypes.IndexPrice,
								CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     bybit.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
								Invert:       false,
							},
							{
								Provider:     mmtypes.IndexPrice,
								CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     coinbaseapi.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     coinbasews.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     huobi.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
								Invert:       false,
							},
							{
								Provider:     mmtypes.IndexPrice,
								CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     kraken.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     kucoin.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
								Invert:       false,
							},
							{
								Provider:     mmtypes.IndexPrice,
								CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     mexc.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
								Invert:       false,
							},
							{
								Provider:     mmtypes.IndexPrice,
								CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     okx.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
								Invert:       false,
							},
							{
								Provider:     mmtypes.IndexPrice,
								CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
								Invert:       false,
							},
						},
					},
				},
			},
			"ETH/USD": {
				Paths: []mmtypes.Path{
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     binance.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
								Invert:       false,
							},
							{
								Provider:     mmtypes.IndexPrice,
								CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     bybit.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
								Invert:       false,
							},
							{
								Provider:     mmtypes.IndexPrice,
								CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     coinbaseapi.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     coinbasews.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     huobi.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
								Invert:       false,
							},
							{
								Provider:     mmtypes.IndexPrice,
								CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     kraken.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     kucoin.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
								Invert:       false,
							},
							{
								Provider:     mmtypes.IndexPrice,
								CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     mexc.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
								Invert:       false,
							},
							{
								Provider:     mmtypes.IndexPrice,
								CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     okx.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
								Invert:       false,
							},
							{
								Provider:     mmtypes.IndexPrice,
								CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
								Invert:       false,
							},
						},
					},
				},
			},
			"USDT/USD": {
				Paths: []mmtypes.Path{
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     binance.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
								Invert:       true,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     bybit.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
								Invert:       true,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     coinbaseapi.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     coinbasews.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     huobi.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
								Invert:       true,
							},
							{
								Provider:     mmtypes.IndexPrice,
								CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     kraken.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     kucoin.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
								Invert:       true,
							},
							{
								Provider:     mmtypes.IndexPrice,
								CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
								Invert:       false,
							},
						},
					},
					{
						Operations: []mmtypes.Operation{
							{
								Provider:     okx.Name,
								CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
								Invert:       true,
							},
						},
					},
				},
			},
		},
		AggregationType: mmtypes.AggregationType_INDEX_PRICE_AGGREGATION,
	},
}
